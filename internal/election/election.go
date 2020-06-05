package election

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	id "gobully/internal/identity"
	"gobully/pkg"
	"time"
)

// TODO think about a config file
const waitingTime = time.Second * 5

// message types
const coordinatorMessage = "CoordinatorUserId"
const answerMessage = "answer"
const electionMessage = "election"

// callback types
const serviceRespondAboard = "aboard"       // CoordinatorUserId found
const serviceRespondMessageReceived = "msg" // service answered

// store service callback here (empty array)
var callbacks []callbackResponse

/* METHODS overview:
	- receiveMessage()             // get a message from a service (election, answer, CoordinatorUserId)
	- sendElectionMessage()        // send a service an election message and wait for response
	- sendCoordinatorMessages()    // send a service that you are the CoordinatorUserId now
      ---------------------
	- messageReceivedAnswer()      // handle answer message
	- messageReceivedElection()    // handle election message
	- messageReceivedCoordinator() // handle coordinator message
 */

/*
receiveMessage POST (Hero <- Hero) - receive message
 */
func receiveMessage(electionInformationString []byte) {
	var electionInformation InformationElection
	err := json.Unmarshal(electionInformationString, &electionInformation)
	if err != nil {
		logrus.Fatalf("[election.receiveMessage] Error unmarshal election message with error %s", err)
	}

	switch electionInformation.Message {
	case answerMessage: messageReceivedAnswer(electionInformation)
	case coordinatorMessage: messageReceivedCoordinator(electionInformation)
	case electionMessage: messageReceivedElection(electionInformation)
	default: fmt.Printf("[election.receiveMessage] message: %s, could not get parsed - abroad ", electionInformation)
	}
	// TODO what to return? is this the default way to send AnswerMessages? - and other less use full ones
}

/*
sendElectionMessage POST (Hero -> Hero) TODO
 */
func sendElectionMessage(electionInformation InformationElection, user id.InformationUser) {
	myElectionInformation := InformationElection{
		Algorithm: electionInformation.Algorithm,
		Payload:   electionMessage,
		Job:       electionInformation.Job,
		Message:   "election in progress please answer me",
	}
	payload, err := json.Marshal(myElectionInformation)
	if err != nil {
		logrus.Fatalf("[election.sendElectionMessage] Error marshal electionCoordinatorMessage with error %s", err)
	}
	// store identity as a new entry in callbacks
	callbacks = append(callbacks, callbackResponse{
		userID:          user.UserId,
		callbackChannel: make(chan string),
		calledBack:      false,
	})
	// send messageReceivedElection to the endpoint
	logrus.Info("[election.sendElectionMessage] send election message to identity: " + user.UserId)
	res, err := pkg.RequestPOST(user.Endpoint +RouteElection, string(payload), "") // TODO wait some time and trigger channel

	// check if identity answered and delete identity from callbacks if so
	// otherwise delete identity form identity list and notify others
	// TODO
	// wait period of time
	// check if a service called back yet
	// 	YES: - aboard other channels, - clear list (other service will take lead)
	// 	NO : send CoordinatorUserId message

	//InformationElection{
	//	Algorithm: electionInformation.Algorithm,
	//	Payload:   answerMessage,
	//	User:      YourUserInformation.UserId,
	//	Job:       electionInformation.Job,
	//	Message:   "election message received, I will take over " + YourUserInformation.UserId,
	//}

	if err != nil {
		logrus.Fatalf("[election.sendElectionMessage] Error send post request with error %s", err)
	}
	var electionInfoResponse InformationElection
	err = json.Unmarshal(res, electionInfoResponse)
	if err != nil {
		logrus.Fatalf("[election.sendElectionMessage] Error Unmarshal electionInfoResponse with error %s", err)
	}
}

/*
sendCoordinatorMessages POST (Hero -> Hero)
 */
func sendCoordinatorMessages(electionInformation InformationElection) {
	// get all users and send a everybody messageReceivedCoordinator
	electionCoordinatorMessage := InformationElection{
		Algorithm: electionInformation.Algorithm,
		Payload:   coordinatorMessage,
		Job:       electionInformation.Job,
		Message:   id.YourUserInformation.UserId, // TODO check if this is the right spot - later
	}
	payload, err := json.Marshal(electionCoordinatorMessage)
	if err != nil {
		logrus.Fatalf("[election.sendCoordinatorMessages] Error marshal electionCoordinatorMessage with error %s", err)
	}
	// send messageReceivedCoordinator to users
	for _, user := range id.Users {
		_, err := pkg.RequestPOST(user.Endpoint +RouteElection, string(payload), "")
		if err != nil {
			logrus.Fatalf("[election.sendCoordinatorMessages] Error sending post request to identity with error %s", err)
		}
	}
	logrus.Info("[election.sendCoordinatorMessages] CoordinatorUserId message send to users")
}

// ------------------------------ HANDLE MESSAGES ------------------------------

/*
messageReceivedAnswer POST (Hero <- Hero) - receive callback message
get a response back from a service after sending a election message
 */
func messageReceivedAnswer(electionInformation InformationElection) {
	// find callback type in var callbacks
	for _, elem := range callbacks {
		if elem.userID == electionInformation.User {
			// check if message is ok and set a bool // ok := bool
			// set var calledBack to ok
			elem.calledBack = true
			// send CallbackMessageReceived through channel
			elem.callbackChannel <- serviceRespondMessageReceived
			logrus.Infof("[election.messageReceivedAnswer] User %s callback received", elem.userID)
		}
	}
}

/*
election message received
 */
func messageReceivedElection(electionInformation InformationElection) {
	logrus.Infof("[election.messageReceivedElection] election notification received, filter users")
	// filter identity after userID > yours
	var selectedUsers []id.InformationUser
	for _, user := range id.Users {
		if user.UserId > id.YourUserInformation.UserId {
			selectedUsers = append(selectedUsers, user)
		}
	}
	// if filtered list is empty - you have the highest ID and win
	if len(selectedUsers) == 0 {
		logrus.Infof("[election.messageReceivedElection] no users found with a higher userId")
		sendCoordinatorMessages(electionInformation)
	} else {
		for _, user := range selectedUsers {
			go sendElectionMessage(electionInformation, user)
		}
		logrus.Infof("[election.messageReceivedElection] election messages send")
	}
}

/*
CoordinatorUserId message received - new CoordinatorUserId found
 */
func messageReceivedCoordinator(electionInformation InformationElection) {
	// close all running elections
	for _, elem := range callbacks {
		elem.callbackChannel <- serviceRespondAboard // tell election services to abroad process
		logrus.Infof("[election.messageReceivedCoordinator] %s told to abroad election process", elem.userID)
	}
	// set CoordinatorUserId
	logrus.Infof("[election.messageReceivedCoordinator] new CoordinatorUserId set")
	CoordinatorUserId = electionInformation.Message // TODO check if reference is correct
}

/* STRUCT */
// control callbacks after sending an election message
type callbackResponse struct {
	userID string               // username as an identifier
	callbackChannel chan string // channel notify after receiving a message
	calledBack bool 		    // tells if a identity send a message back
}