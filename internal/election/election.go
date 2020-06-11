package election

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"goBully/internal/identity"
	"goBully/pkg"
	"time"
)

/* METHODS overview:
	- receiveMessage()            // get a message from a api (election, coordinator)
	ElectionMessage
	- receiveMessageElection()    // handle incoming election message
	- sendMessageElection()       // send a election message to another user
    CoordinatorMessage
	- receiveMessageCoordinator() // set local coordinator reference with incoming details
	- sendMessagesCoordinator()   // send coordinator messages to other users
 */

/*
receiveMessage POST (Hero <- Hero) - receive message
 */
func receiveMessage(electionInformation InformationElectionDTO) InformationElectionDTO {
	// response is set in receiveMessageCoordinator && receiveMessageElection
	var electionInformationResponse InformationElectionDTO
	switch electionInformation.Payload {
		case MessageCoordinator: receiveMessageCoordinator(electionInformation, &electionInformationResponse)
		case MessageElection: receiveMessageElection(electionInformation, &electionInformationResponse)
		default: logrus.Warningf("[election.receiveMessage] message: %s, could not be identified - abroad ", electionInformation.Algorithm)
	}
	return electionInformationResponse
}

// ELECTION

/*
receiveMessageElection - election message received
---
receiveMessageElection(InformationElectionDTO)
1. filter users to send election messages to (UserID > YourID)
2. if |filtered users| <= 0
   	YES: 2.1 you have the highest ID and win - send coordinatorMessages - exit
   	NO : 2.2 transform message and create POST payload
		 2.3 add user information to local callbackList
         2.4 GO - sendMessageElection(callbackResponse, msgPayload)
            2.4.1 send POST request to client
            2.4.2 if response is OK add client to client who have responded responded
         2.5 wait a few seconds (enough time users can answer request)
         2.6 Sort users who have called back and who are not
         2.7 if |answered users| <= 0
			2.7.1 YES: send coordinatorMessages
		 2.8 remove all users how didn't answered from userList
         2.9 clear callback list
3. send response back (answer)
*/
func receiveMessageElection(electionInformation InformationElectionDTO, electionInformationResponse *InformationElectionDTO) {
	logrus.Infof("[election.receiveMessageElection] election message received")
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
		logrus.Infof("[election.receiveMessageElection] no users found with a higher userId")
		sendMessagesCoordinator(electionInformation)
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
			logrus.Fatalf("[election.receiveMessageElection] Error marshal electionCoordinatorMessage with error %s", err)
		}
		var callbacks []identity.InformationUserDTO
		var didCallBackUsers []identity.InformationUserDTO // store all users who have replied
		for _, user := range selectedUsers {
			// 2.3 add user information to local callbackList
			callbacks = append(callbacks, user)
			// 2.4 GO - sendMessageElection()
			go sendMessageElection(&didCallBackUsers, &user, payload)
		}
		logrus.Infof("[election.receiveMessageElection] election messages send, waiting " + WaitingTime.String() + " seconds for a response")
		// 2.5 wait a few seconds (enough time users can answer request)
		time.Sleep(WaitingTime)
		// 2.6 Sort users who have called back and who are not
		if len(callbacks) != len(didCallBackUsers) {
			for _, user := range callbacks {
				if identity.ContainsUser(didCallBackUsers, user) {
					logrus.Warnf("[election.receiveMessageElection] user %s did not call back", user.UserId)
					identity.DeleteUser(user)
				}
			}
		}
		// 2.7 if |answered users| <= 0
		if len(didCallBackUsers) <= 0 {
			// 2.7.1 send coordinatorMessages
			sendMessagesCoordinator(electionInformation)
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
sendMessageElection POST (Hero -> Hero)
ALGORITHM - OVERVIEW
2.4.1 send POST request to client
2.4.2 if response is OK add client to client who have responded responded
 */
func sendMessageElection(didCallback *[]identity.InformationUserDTO, userInfoCallback *identity.InformationUserDTO, msgPayload []byte) {
	// 2.4.1 send POST request to client
	logrus.Info("[election.sendMessageElection] send election message to identity: " + userInfoCallback.UserId)
	res, err := pkg.RequestPOST(userInfoCallback.Endpoint + RouteElection, string(msgPayload))
	if err != nil {
		logrus.Fatalf("[election.sendMessageElection] Error send post request with error %s", err)
	}
	// 2.4.2 if response is OK add client to client who have responded responded
	var electionAnswerResponse InformationElectionDTO
	err = json.Unmarshal(res, &electionAnswerResponse)
	if err != nil {
		logrus.Fatalf("[election.sendMessageElection] Error Unmarshal electionAnswerResponse with error %s", err)
	}
	logrus.Infof("[election.sendMessageElection] response received, user: %s", electionAnswerResponse.User)
	if electionAnswerResponse.Payload == MessageAnswer {
		// add client to clients how replied
		*didCallback = append(*didCallback, *userInfoCallback)
	}
}

// COORDINATOR

/*
CoordinatorUserId message received - new CoordinatorUserId found
 */
func receiveMessageCoordinator(electionInformation InformationElectionDTO, electionInformationResponse *InformationElectionDTO) {
	// set CoordinatorUserId to local coordinator information
	logrus.Infof("[election.receiveMessageCoordinator] new CoordinatorUserId set")
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
sendMessagesCoordinator POST (Hero -> Hero)
1. create coordinator message
2. send all users the coordinator message
*/
func sendMessagesCoordinator(electionInformation InformationElectionDTO) {
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
		logrus.Fatalf("[election.sendMessagesCoordinator] Error marshal electionCoordinatorMessage with error %s", err)
	}
	// 2. send all users the coordinator message
	for _, user := range identity.Users {
		_, err := pkg.RequestPOST(user.Endpoint + RouteElection, string(payload))
		if err != nil {
			logrus.Fatalf("[election.sendMessagesCoordinator] Error sending post request to identity with error %s", err)
		}
	}
	logrus.Info("[election.sendMessagesCoordinator] CoordinatorUserId message send to users")
}
