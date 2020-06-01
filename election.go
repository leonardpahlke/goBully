package goBully

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

const electionEndpoint = "election"
const waitingTime = time.Second * 5 // time to wait until a service is declared as unreachable

// message types
const CoordinatorMessage = "coordinator"
const AnswerMessage = "answer"
const ElectionMessage = "election"

// callback types
const CallbackAboard = "aboard"        // coordinator found
const CallbackMessageReceived  = "msg" // service answered

// current coordinator
var coordinator = ""

// store service callback here (empty array)
var callbacks []CallbackResponse

/** METHODS:
	- ReceiveMessage()         // get a message from a service (election, answer, coordinator)
	- ReceiveCallback()        // get a response after sending a election message
	- SendAnswerMessage()      // callback to service
	- SendElectionMessage()    // send a service an election message and wait for response
	- SendCoordinatorMessages() // send a service that you are the coordinator now
      ---------------------
	- ElectionMessageReceived()        // handle election message
	- CoordinatorMessageReceived()     // handle coordinator message
 */

// ReceiveMessage POST (Hero <- Hero) - receive message
func ReceiveMessage(electionInformationString []byte) {
	var electionInformation ElectionInformation
	err := json.Unmarshal(electionInformationString, &electionInformation)
	if err != nil {
		logrus.Fatalf("[election.ReceiveMessage] Error unmarshal election message with error %s", err)
	}
	switch electionInformation.Message {
	case CoordinatorMessage: CoordinatorMessageReceived(electionInformation)
	case ElectionMessage: ElectionMessageReceived(electionInformation)
	default: fmt.Printf("[election.ReceiveMessage] message: %s, could not get parsed - abroad ", electionInformation)
	}
}

// ReceiveCallback POST (Hero <- Hero) - receive callback message
// get a response back from a service after sending a election message
func ReceiveCallback(electionCallbackString []byte) {
	var electionCallback ElectionCallbackInformation
	err := json.Unmarshal(electionCallbackString, &electionCallback)
	if err != nil {
		logrus.Fatalf("[election.ReceiveCallback] Error unmarshal election callback message with error %s", err)
	}
	// find callback type in var callbacks
	for _, elem := range callbacks {
		if elem.userID == electionCallback.User {
			// check if message is ok and set a bool // ok := bool
			// set var calledBack to ok
			elem.calledBack = true
			// send CallbackMessageReceived through channel
			elem.callbackChannel <- CallbackMessageReceived
			logrus.Infof("[election.ReceiveCallback] User %s callback received", elem.userID)
		}
	}
}

// SendAnswerMessage POST (Hero -> Hero) (CALLBACK)
func SendAnswerMessage(electionCallbackInformation ElectionInformation) {
	electionInformation := ElectionCallbackInformation{
		Algorithm: electionCallbackInformation.Algorithm,
		Payload:   AnswerMessage,
		User:      YourUserInformation.UserID,
		Job:       electionCallbackInformation.Job,
		Message:   "election message received, I will take over " + YourUserInformation.UserID,
	}
	// send an AnswerMessage to the endpoint
	payload, err := json.Marshal(electionInformation)
	if err != nil {
		logrus.Fatalf("[election.SendAnswerMessage] Error marshal electionInformation with error %s", err)
	}
	// TODO maybe wait if user answer's and remove him from list (notify others that he left)
	_, err = RequestPOST(electionCallbackInformation.Callback, string(payload), "")
	if err != nil {
		logrus.Fatalf("[election.SendAnswerMessage] Error sending post request to user with error %s", err)
	}
	logrus.Info("[election.SendAnswerMessage] Answer message send to user")
}

// SendElectionMessage POST (Hero -> Hero) TODO
func SendElectionMessage(electionInformation ElectionInformation, user UserInformation) {
	myElectionInformation := ElectionInformation{
		Algorithm: electionInformation.Algorithm,
		Payload:   ElectionMessage,
		Callback:  YourUserInformation.CallbackEndpoint,
		Job:       electionInformation.Job,
		Message:   "election in progress please answer me",
	}
	payload, err := json.Marshal(myElectionInformation)
	if err != nil {
		logrus.Fatalf("[election.SendElectionMessage] Error marshal electionCoordinatorMessage with error %s", err)
	}
	// store user as a new entry in callbacks
	callbacks = append(callbacks, CallbackResponse{
		userID:          user.UserID,
		callbackChannel: make(chan string),
		calledBack:      false,
	})
	// send ElectionMessageReceived to the endpoint
	logrus.Info("[election.SendElectionMessage] send election message to user: " + user.UserID)
	reponse, err := RequestPOST(user.Endpoint + "/" + electionEndpoint, string(payload), "") // TODO response? callback?
	// check if user answered and delete user from callbacks if so
	// otherwise delete user form user list and notify others
}

// SendCoordinatorMessages POST (Hero -> Hero)
func SendCoordinatorMessages(electionInformation ElectionInformation) {
	// get all users and send a everybody CoordinatorMessageReceived
	electionCoordinatorMessage := ElectionInformation{
		Algorithm: electionInformation.Algorithm,
		Payload:   CoordinatorMessage,
		Callback:  YourUserInformation.CallbackEndpoint,
		Job:       electionInformation.Job,
		Message:   YourUserInformation.UserID, // TODO check if this is the right spot - later
	}
	payload, err := json.Marshal(electionCoordinatorMessage)
	if err != nil {
		logrus.Fatalf("[election.SendCoordinatorMessages] Error marshal electionCoordinatorMessage with error %s", err)
	}
	// send CoordinatorMessageReceived to users
	for _, user := range Users {
		_, err := RequestPOST(user.Endpoint + "/" + electionEndpoint, string(payload), "")
		if err != nil {
			logrus.Fatalf("[election.SendCoordinatorMessages] Error sending post request to user with error %s", err)
		}
	}
	logrus.Info("[election.SendCoordinatorMessages] Coordinator message send to users")
}

// ------------------------------ HANDLE MESSAGES ------------------------------

// election message received TODO
func ElectionMessageReceived(electionInformation ElectionInformation) {
	// get all users
	// filter user after userID > yours
	var selectedUsers []UserInformation
	for _, user := range Users {
		if user.UserID > YourUserInformation.UserID {
			selectedUsers = append(selectedUsers, user)
		}
	}
	// if filtered list is empty - you have the highest ID and win
	if len(selectedUsers) == 0 {
		SendCoordinatorMessages(electionInformation)
	}

	// send user election message
	// wait period of time
	// check if a service calledback yet
	// 	YES: - aboard other channels, - clear list (other service will take lead)
	// 	NO : send coordinator message
}

// coordinator message received - new coordinator found
func CoordinatorMessageReceived(electionInformation ElectionInformation) {
	// close all running elections
	for _, elem := range callbacks {
		elem.callbackChannel <- CallbackAboard // tell election services to abroad process
		logrus.Infof("[election.CoordinatorMessageReceived] %s told to abroad election process", elem.userID)
	}
	// set coordinator
	logrus.Infof("[election.CoordinatorMessageReceived] new coordinator set")
	coordinator = electionInformation.Message // TODO check if reference is correct
}