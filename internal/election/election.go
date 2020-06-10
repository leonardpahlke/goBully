package election

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"goBully/internal/identity"
	"goBully/pkg"
	"time"
)

// TODO think about a config file
const waitingTime = time.Second * 3

/* METHODS overview:
	- receiveMessage()             // get a message from a api (election, coordinator)
	- messageReceivedElection()    // handle incoming election message
	- sendElectionMessage()        // send a election message to another user
      ---------------------
	- messageReceivedCoordinator() // set local coordinator reference with incoming details
	- sendCoordinatorMessages()    // send coordinator messages to other users
 */

/*
receiveMessage POST (Hero <- Hero) - receive message
 */
func receiveMessage(electionInformation InformationElectionDTO) InformationElectionDTO {
	// response is set in messageReceivedCoordinator && messageReceivedElection
	var electionInformationResponse InformationElectionDTO

	switch electionInformation.Payload {
		case MessageCoordinator: messageReceivedCoordinator(electionInformation, &electionInformationResponse)
		case MessageElection: messageReceivedElection(electionInformation, &electionInformationResponse)
		default: fmt.Printf("[election.receiveMessage] message: %s, could not get parsed - abroad ", electionInformation.Algorithm)
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
	logrus.Infof("[election.messageReceivedElection] election message received")
	// 1. filter users to send election messages to (UserID > YourID)
	var selectedUsers []identity.InformationUserDTO
	for _, user := range identity.Users {
		if user.UserId > identity.YourUserInformation.UserId {
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
			Payload:   MessageElection,
			User:      identity.YourUserInformation.UserId,
			Job:       electionInformation.Job,
			Message:   "election in progress please answer me",
		}
		payload, err := json.Marshal(myElectionInformation)
		if err != nil {
			logrus.Fatalf("[election.messageReceivedElection] Error marshal electionCoordinatorMessage with error %s", err)
		}
		// store api callback here (empty array)
		var callbacks []identity.InformationUserDTO
		var didCallBackUsers []identity.InformationUserDTO // store all users who have replied
		for _, user := range selectedUsers {
			callbacks = append(callbacks, user)
			// 2.3 GO - sendElectionMessage()
			go sendElectionMessage(&didCallBackUsers, &user, payload)
		}
		logrus.Infof("[election.messageReceivedElection] election messages send, waiting " + waitingTime.String() + " seconds for a response")
		// 2.4 wait a few seconds (enough time users can answer request)
		time.Sleep(waitingTime)
		// 2.6 Sort users who have called back and who are not
		if len(callbacks) != len(didCallBackUsers) {
			for _, user := range callbacks {
				if identity.ContainsUser(didCallBackUsers, user) {
					logrus.Warnf("[election.messageReceivedElection] user %s did not call back", user.UserId)
					identity.DeleteUser(user)
				}
			}
		}
		// 2.7 if |answered users| <= 0
		if len(didCallBackUsers) <= 0 {
			// 2.7.1 send coordinatorMessages
			sendCoordinatorMessages(electionInformation)
		}
		// 2.8 remove all users how didn't answered from userList
		// 2.9 clear callback list
		callbacks = []identity.InformationUserDTO{}
	}
	// 3. send response back (answer)
	*electionInformationResponse = InformationElectionDTO{
		Algorithm: electionInformation.Algorithm,
		Payload:   MessageAnswer,
		User:      identity.YourUserInformation.UserId,
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
func sendElectionMessage(didCallback *[]identity.InformationUserDTO, userInfoCallback *identity.InformationUserDTO, msgPayload []byte) {
	// 2.4.1 send POST request to client
	logrus.Info("[election.sendElectionMessage] send election message to identity: " + userInfoCallback.UserId)
	res, err := pkg.RequestPOST(userInfoCallback.Endpoint + RouteElection, string(msgPayload))
	if err != nil {
		logrus.Fatalf("[election.sendElectionMessage] Error send post request with error %s", err)
	}
	// 2.4.2 if response is OK check client callback
	var electionAnswerResponse InformationElectionDTO
	err = json.Unmarshal(res, &electionAnswerResponse)
	if err != nil {
		logrus.Fatalf("[election.sendElectionMessage] Error Unmarshal electionAnswerResponse with error %s", err)
	}
	logrus.Infof("[election.sendElectionMessage] response received, user: %s", electionAnswerResponse.User)
	if electionAnswerResponse.Payload == MessageAnswer {
		// add client to clients how replied
		*didCallback = append(*didCallback, *userInfoCallback)
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
		Payload:   MessageAnswer,
		User:      identity.YourUserInformation.UserId,
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
		Payload:   MessageCoordinator,
		Job:       electionInformation.Job,
		User:      identity.YourUserInformation.UserId,
		Message:   "new elected coordinator found",
	}
	payload, err := json.Marshal(electionCoordinatorMessage)
	if err != nil {
		logrus.Fatalf("[election.sendCoordinatorMessages] Error marshal electionCoordinatorMessage with error %s", err)
	}
	// 2. send all users the coordinator message
	for _, user := range identity.Users {
		_, err := pkg.RequestPOST(user.Endpoint + RouteElection, string(payload))
		if err != nil {
			logrus.Fatalf("[election.sendCoordinatorMessages] Error sending post request to identity with error %s", err)
		}
	}
	logrus.Info("[election.sendCoordinatorMessages] CoordinatorUserId message send to users")
}
