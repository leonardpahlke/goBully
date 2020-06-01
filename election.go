package goBully

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

// message types
const CoordinatorMessage = "coordinator"
const AnswerMessage = "answer"
const ElectionMessage = "election"

// callback types
const CallbackAboard = "aboard"        // coordinator found
const CallbackMessageReceived  = "msg" // service answered

// control callbacks after sending an election message
type callbackResponse struct {
	userID string               // username as an identifier
	callbackChannel chan string // channel notify after receiving a message
	calledBack bool 		    // tells if a user send a message back
}

// current coordinator
var coordinator = ""

// store service callback here (empty array)
var callbacks = []callbackResponse[]

/** METHODS:
	- receiveMessage()         // get a message from a service (election, answer, coordinator)
	- receiveCallback()        // get a response after sending a election message
	- sendAnswerMessage()      // callback to service
	- sendElectionMessage()    // send a service an election message and wait for response
	- sendCoordinatorMessages() // send a service that you are the coordinator now
      ---------------------
	- electionMessage()        // handle election message
	- coordinatorMessage()     // handle coordinator message
 */

// receiveMessage POST (Hero <- Hero) - receive message
func receiveMessage(electionInformationString []byte) {
	var electionInformation ElectionInformation
	err := json.Unmarshal(electionInformationString, &electionInformation)
	if err != nil {
		logrus.Fatalf("[receiveMessage] Error unmarshal election message with error %s", err)
	}
	switch electionInformation.Message {
	case CoordinatorMessage: coordinatorMessage(electionInformation)
	case ElectionMessage: electionMessage(electionInformation)
	default: fmt.Printf("[receiveMessage] message: %s, could not get parsed - abroad ", electionInformation)
	}
}

// receiveCallback POST (Hero <- Hero) - receive callback message
// get a response back from a service after sending a election message
func receiveCallback(electionCallbackString []byte) {
	var electionCallback ElectionCallbackInformation
	err := json.Unmarshal(electionCallbackString, &electionCallback)
	if err != nil {
		logrus.Fatalf("[receiveCallback] Error unmarshal election callback message with error %s", err)
	}
	// find callback type in var callbacks
	for _, elem := range callbacks {
		if elem.userID == electionCallback.User {
			// check if message is ok and set a bool // ok := bool
			// set var calledBack to ok
			elem.calledBack = true
			// send CallbackMessageReceived through channel
			elem.callbackChannel <- CallbackMessageReceived
			logrus.Infof("[receiveCallback] User %s callback received", elem.userID)
		}
	}
}


// sendAnswerMessage POST (Hero -> Hero) (CALLBACK) TODO
func sendAnswerMessage() {
	// send an AnswerMessage to the endpoint
}

// sendElectionMessage POST (Hero -> Hero) TODO
func sendElectionMessage() {
	// store user as a new entry in callbacks
	// send an ElectionMessage to the endpoint
}

// sendCoordinatorMessages POST (Hero -> Hero) TODO
func sendCoordinatorMessages() {
	// get all users and send a everybody CoordinatorMessage
	// send an CoordinatorMessage to the endpoint
}

// ------------------------------ HANDLE MESSAGES ------------------------------

// election message received TODO
func electionMessage(electionInformation ElectionInformation) {
	// get all users
	// filter user after userID > yours
	var selectedUsers []UserInformation[]
	for _, user := range Users {
		if user.UserID > YourUserInformation.UserID {
			selectedUsers = append(selectedUsers, user)
		}
	}
	// if filtered list is empty - you have the highest ID and win
	if len(selectedUsers) == 0 {
		sendCoordinatorMessages()
	}

	// send user election message
	// wait period of time
	// check if a service calledback yet
	// 	YES: - aboard other channels, - clear list (other service will take lead)
	// 	NO : send coordinator message
}

// coordinator message received - new coordinator found
func coordinatorMessage(electionInformation ElectionInformation) {
	// close all running elections
	fmt.Println("coordinator found")
	for _, elem := range callbacks {
		elem.callbackChannel <- CallbackAboard // tell election services to abroad process
		logrus.Infof("[coordinatorMessage] %s told to abroad election process", elem.userID)
	}
	// set coordinator
	logrus.Infof("[coordinatorMessage] new coordinator set")
	coordinator = electionInformation.Message // TODO check
}