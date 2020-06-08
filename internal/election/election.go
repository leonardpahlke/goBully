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
const waitingTime = time.Second * 2

// message types
const coordinatorMessage = "CoordinatorUserId"
const answerMessage = "answer"
const electionMessage = "election"

// store service callback here (empty array)
var callbacks []callbackResponse

/* METHODS overview:
	- ReceiveMessage()             // get a message from a service (election, coordinator)
	- messageReceivedElection()    // handle incoming election message
	- sendElectionMessage()        // send a election message to another user
      ---------------------
	- messageReceivedCoordinator() // set local coordinator reference with incoming details
	- sendCoordinatorMessages()    // send coordinator messages to other users
 */

/*
ReceiveMessage POST (Hero <- Hero) - receive message
 */
func ReceiveMessage(electionInformation InformationElectionDTO) InformationElectionDTO {
	// response is set in messageReceivedCoordinator && messageReceivedElection
	var electionInformationResponse InformationElectionDTO

	switch electionInformation.Message {
		case coordinatorMessage: messageReceivedCoordinator(electionInformation, &electionInformationResponse)
		case electionMessage: messageReceivedElection(electionInformation, &electionInformationResponse)
		default: fmt.Printf("[election.ReceiveMessage] message: %s, could not get parsed - abroad ", electionInformation)
	}
	return electionInformationResponse
}

// ELECTION

/*
election message received
---
messageReceivedElection(InformationElectionDTO)
1. filter users to send election messages to (UserID > YourID)
2. if |filtered users| <= 0
   	YES: 2.1 you have the highest ID and win - send coordinatorMessages - exit
   	NO : 2.2 transform message and create POST payload
		 2.3 add callback information to local callbackList
         2.4 GO - sendElectionMessage(callbackResponse, msgPayload)
            2.4.1 send POST request to client
            2.4.2 if response is OK check client callback
         2.5 wait a few seconds (enough time users can answer request)
         2.6 Sort users who have called back and who are not
         2.7 if |answered users| <= 0
			2.7.1 YES: send coordinatorMessages
		 2.8 remove all users how didn't answered from userList
         2.9 clear callback list
3. send response back (answer)
*/
func messageReceivedElection(electionInformation InformationElectionDTO, electionInformationResponse *InformationElectionDTO) {
	logrus.Infof("[election.messageReceivedElection] election notification received, filter users")
	// 1. filter users to send election messages to (UserID > YourID)
	var selectedUsers []id.InformationUserDTO
	for _, user := range id.Users {
		if user.UserId > id.YourUserInformation.UserId {
			selectedUsers = append(selectedUsers, user)
		}
	}
	// 2. if filtered users <= 0
	if len(selectedUsers) == 0 {
		// 2.1 you have the highest ID and win - send coordinatorMessages - exit
		logrus.Infof("[election.messageReceivedElection] no users found with a higher userId")
		sendCoordinatorMessages(electionInformation)
	} else {
		// 2.2 transform message and create POST payload
		myElectionInformation := InformationElectionDTO{
			Algorithm: electionInformation.Algorithm,
			Payload:   electionMessage,
			User: 	   id.YourUserInformation.UserId,
			Job:       electionInformation.Job,
			Message:   "election in progress please answer me",
		}
		payload, err := json.Marshal(myElectionInformation)
		if err != nil {
			logrus.Fatalf("[election.sendElectionMessage] Error marshal electionCoordinatorMessage with error %s", err)
		}
		for _, user := range selectedUsers {
			userCallback := callbackResponse{
				userInfo:   user,
				calledBack: false,
			}
			callbacks = append(callbacks, userCallback)
			// 2.3 GO - sendElectionMessage()
			go sendElectionMessage(&userCallback, payload)
		}
		logrus.Infof("[election.messageReceivedElection] election messages send, waiting " + waitingTime.String() + " seconds for a response")
		// 2.4 wait a few seconds (enough time users can answer request)
		time.Sleep(waitingTime)
		// 2.6 Sort users who have called back and who are not
		var didCallBackUsers []id.InformationUserDTO   // all users who have replied
		var didntCallBackUsers []id.InformationUserDTO // all users who have not replied
		for _, userCallback := range callbacks {
			if userCallback.calledBack {
				didCallBackUsers = append(didCallBackUsers, userCallback.userInfo)
			} else {
				didntCallBackUsers = append(didntCallBackUsers, userCallback.userInfo)
			}
		}
		// 2.7 if |answered users| <= 0
		if len(didCallBackUsers) <= 0 {
			// 2.7.1 send coordinatorMessages
			sendCoordinatorMessages(electionInformation)
		}
		// 2.8 remove all users how didn't answered from userList
		for _, user := range didntCallBackUsers {
			id.DeleteUser(user)
		}
		// 2.9 clear callback list
		callbacks = []callbackResponse{}
	}
	// 3. send response back (answer)
	*electionInformationResponse = InformationElectionDTO{
		Algorithm: electionInformation.Algorithm,
		Payload:   answerMessage,
		User:      id.YourUserInformation.UserId,
		Job:       electionInformation.Job,
		Message:   "election message send to the others",
	}
}

/*
sendElectionMessage POST (Hero -> Hero)
ALGORITHM - OVERVIEW
2.4.1 send POST request to client
2.4.2 if response is OK check client callback
 */
func sendElectionMessage(userCallback *callbackResponse, msgPayload []byte) {
	// 2.4.1 send POST request to client
	logrus.Info("[election.sendElectionMessage] send election message to identity: " + userCallback.userInfo.UserId)
	res, err := pkg.RequestPOST(userCallback.userInfo.Endpoint + RouteElection, string(msgPayload), "")
	if err != nil {
		logrus.Fatalf("[election.sendElectionMessage] Error send post request with error %s", err)
	}
	// 2.4.2 if response is OK check client callback
	var electionAnswerResponse InformationElectionDTO
	err = json.Unmarshal(res, &electionAnswerResponse)
	if err != nil {
		logrus.Fatalf("[election.sendElectionMessage] Error Unmarshal electionAnswerResponse with error %s", err)
	}
	if electionAnswerResponse.Payload == answerMessage {
		// check client callback
		*userCallback = callbackResponse{
			userInfo: userCallback.userInfo,
			calledBack: true,
		}
	}
}

// COORDINATOR

/*
CoordinatorUserId message received - new CoordinatorUserId found
 */
func messageReceivedCoordinator(electionInformation InformationElectionDTO, electionInformationResponse *InformationElectionDTO) {
	// set CoordinatorUserId to local coordinator information
	logrus.Infof("[election.messageReceivedCoordinator] new CoordinatorUserId set")
	CoordinatorUserId = electionInformation.User
	*electionInformationResponse = InformationElectionDTO{
		Algorithm: electionInformation.Algorithm,
		Payload:   answerMessage,
		User:      id.YourUserInformation.UserId,
		Job:       electionInformation.Job,
		Message:   "OK",
	}
}

/*
sendCoordinatorMessages POST (Hero -> Hero)
---
1. create coordinator message
2. send all users the coordinator message
*/
func sendCoordinatorMessages(electionInformation InformationElectionDTO) {
	// 1. create coordinator message
	electionCoordinatorMessage := InformationElectionDTO{
		Algorithm: electionInformation.Algorithm,
		Payload:   coordinatorMessage,
		Job:       electionInformation.Job,
		User:      id.YourUserInformation.UserId,
		Message:   "new elected coordinator found",
	}
	payload, err := json.Marshal(electionCoordinatorMessage)
	if err != nil {
		logrus.Fatalf("[election.sendCoordinatorMessages] Error marshal electionCoordinatorMessage with error %s", err)
	}
	// 2. send all users the coordinator message
	for _, user := range id.Users {
		_, err := pkg.RequestPOST(user.Endpoint + RouteElection, string(payload), "")
		if err != nil {
			logrus.Fatalf("[election.sendCoordinatorMessages] Error sending post request to identity with error %s", err)
		}
	}
	logrus.Info("[election.sendCoordinatorMessages] CoordinatorUserId message send to users")
}

/* STRUCT */
// control callbacks after sending an election message
type callbackResponse struct {
	userInfo   id.InformationUserDTO // username as an identifier
	calledBack bool                  // tells if a identity send a message back
}