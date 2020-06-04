package internal

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

// API Endpoints
const ElectionRoute = "/election"

// TODO (use) - time to wait until a service is declared as unreachable
const WaitingTime = time.Second * 5

// message types
const CoordinatorMessage = "coordinator"
const AnswerMessage = "answer"
const ElectionMessage = "election"

// callback types
const ServiceRespondAboard = "aboard"       // coordinator found
const ServiceRespondMessageReceived = "msg" // service answered

// current coordinator userId
var coordinator = ""

// store service callback here (empty array)
var callbacks []CallbackResponse

/** METHODS:
	- ReceiveMessage()             // get a message from a service (election, answer, coordinator)
	- SendAnswerMessage()          // callback to service
	- SendElectionMessage()        // send a service an election message and wait for response
	- SendCoordinatorMessages()    // send a service that you are the coordinator now
      ---------------------
	- ElectionMessageReceived()    // handle election message
	- CoordinatorMessageReceived() // handle coordinator message
 */

// ReceiveMessage POST (Hero <- Hero) - receive message
func ReceiveMessage(electionInformationString []byte) {
	var electionInformation ElectionInformation
	err := json.Unmarshal(electionInformationString, &electionInformation)
	if err != nil {
		logrus.Fatalf("[election.ReceiveMessage] Error unmarshal election message with error %s", err)
	}

	switch electionInformation.Message {
	case AnswerMessage: AnswerMessageReceived(electionInformation)
	case CoordinatorMessage: CoordinatorMessageReceived(electionInformation)
	case ElectionMessage: ElectionMessageReceived(electionInformation)
	default: fmt.Printf("[election.ReceiveMessage] message: %s, could not get parsed - abroad ", electionInformation)
	}
	// TODO what to return? is this the default way to send AnswerMessages? - and other less use full ones
}

// SendElectionMessage POST (Hero -> Hero) TODO
func SendElectionMessage(electionInformation ElectionInformation, user UserInformation) {
	myElectionInformation := ElectionInformation{
		Algorithm: electionInformation.Algorithm,
		Payload:   ElectionMessage,
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
	res, err := RequestPOST(user.Endpoint + ElectionRoute, string(payload), "") // TODO wait some time and trigger channel

	// check if user answered and delete user from callbacks if so
	// otherwise delete user form user list and notify others
	// TODO TODO TODO
	// wait period of time
	// check if a service called back yet
	// 	YES: - aboard other channels, - clear list (other service will take lead)
	// 	NO : send coordinator message

	//ElectionInformation{
	//	Algorithm: electionInformation.Algorithm,
	//	Payload:   AnswerMessage,
	//	User:      YourUserInformation.UserID,
	//	Job:       electionInformation.Job,
	//	Message:   "election message received, I will take over " + YourUserInformation.UserID,
	//}

	if err != nil {
		logrus.Fatalf("[election.SendElectionMessage] Error send post request with error %s", err)
	}
	var electionInfoResponse ElectionInformation
	err = json.Unmarshal(res, electionInfoResponse)
	if err != nil {
		logrus.Fatalf("[election.SendElectionMessage] Error Unmarshal electionInfoResponse with error %s", err)
	}
}

// SendCoordinatorMessages POST (Hero -> Hero)
func SendCoordinatorMessages(electionInformation ElectionInformation) {
	// get all users and send a everybody CoordinatorMessageReceived
	electionCoordinatorMessage := ElectionInformation{
		Algorithm: electionInformation.Algorithm,
		Payload:   CoordinatorMessage,
		Job:       electionInformation.Job,
		Message:   YourUserInformation.UserID, // TODO check if this is the right spot - later
	}
	payload, err := json.Marshal(electionCoordinatorMessage)
	if err != nil {
		logrus.Fatalf("[election.SendCoordinatorMessages] Error marshal electionCoordinatorMessage with error %s", err)
	}
	// send CoordinatorMessageReceived to users
	for _, user := range Users {
		_, err := RequestPOST(user.Endpoint + ElectionRoute, string(payload), "")
		if err != nil {
			logrus.Fatalf("[election.SendCoordinatorMessages] Error sending post request to user with error %s", err)
		}
	}
	logrus.Info("[election.SendCoordinatorMessages] Coordinator message send to users")
}

// ------------------------------ HANDLE MESSAGES ------------------------------

// AnswerMessageReceived POST (Hero <- Hero) - receive callback message
// get a response back from a service after sending a election message
func AnswerMessageReceived(electionInformation ElectionInformation) {
	// find callback type in var callbacks
	for _, elem := range callbacks {
		if elem.userID == electionInformation.User {
			// check if message is ok and set a bool // ok := bool
			// set var calledBack to ok
			elem.calledBack = true
			// send CallbackMessageReceived through channel
			elem.callbackChannel <- ServiceRespondMessageReceived
			logrus.Infof("[election.AnswerMessageReceived] User %s callback received", elem.userID)
		}
	}
}

// election message received TODO
func ElectionMessageReceived(electionInformation ElectionInformation) {
	logrus.Infof("[election.ElectionMessageReceived] election notification received, filter users")
	// filter user after userID > yours
	var selectedUsers []UserInformation
	for _, user := range Users {
		if user.UserID > YourUserInformation.UserID {
			selectedUsers = append(selectedUsers, user)
		}
	}
	// if filtered list is empty - you have the highest ID and win
	if len(selectedUsers) == 0 {
		logrus.Infof("[election.ElectionMessageReceived] no users found with a higher userId")
		SendCoordinatorMessages(electionInformation)
	} else {
		for _, user := range selectedUsers {
			go SendElectionMessage(electionInformation, user)
		}
		logrus.Infof("[election.ElectionMessageReceived] election messages send")
	}
}

// coordinator message received - new coordinator found
func CoordinatorMessageReceived(electionInformation ElectionInformation) {
	// close all running elections
	for _, elem := range callbacks {
		elem.callbackChannel <- ServiceRespondAboard // tell election services to abroad process
		logrus.Infof("[election.CoordinatorMessageReceived] %s told to abroad election process", elem.userID)
	}
	// set coordinator
	logrus.Infof("[election.CoordinatorMessageReceived] new coordinator set")
	coordinator = electionInformation.Message // TODO check if reference is correct
}

// STRUCT'S
type ElectionInformation struct {
	Algorithm string         `json:"algorithm"` // name of the algorithm used
	Payload   string         `json:"payload"`   // the payload for the current state of the algorithm
	User      string         `json:"user"`  // uri of the user sending this request
	Job       JobInformation `json:"job"`
	Message   string         `json:"message"`   // something you want to tell the other one
}

type JobInformation struct {
	Id       string `json:"id"`       // some identity choosen by the initiator to identify this request
	Task     string `json:"task"`     // uri to the task to accomplish
	Resource string `json:"resource"` // uri or url to resource where actions are required
	Method   string `json:"method"`   // method to take – if already known
	Data     string `json:"data"`     // data to use/post for the task
	Callback string `json:"callback"` // an url where the initiator can be reached with the results/token
	Message  string `json:"message"`  // something you want to tell the other one
}

// control callbacks after sending an election message
type CallbackResponse struct {
	userID string               // username as an identifier
	callbackChannel chan string // channel notify after receiving a message
	calledBack bool 		    // tells if a user send a message back
}