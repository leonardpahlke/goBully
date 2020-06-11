package mutex

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"goBully/internal/identity"
	"goBully/pkg"
	"time"
)

/*
requestCriticalArea - tell all users that this user wants to enter the critical section
1. set state to 'wanting'
2. increment clock, you are about to send messages
3. create a request mutex message
4. GO - send all users the mutex message
5. wait for all users to reply-ok to your request
6. enterCriticalSection() - and leave critical section if this method returns
*/
func requestCriticalArea() {
	// 1. set state to 'wanting'
	state = StateWanting
	logrus.Infof("[mutex_send.requestCriticalArea] starting, state: %s", state)
	// 2. increment clock, you are about to send messages
	incrementClock(clock)
	// 3. create a request mutex message
	var mutexRequestMessage = MessageMutexDTO{
		Msg:   RequestMessage,
		Time:  clock,
		Reply: mutexYourReply,
		User:  mutexYourUser,
	}
	payload, err := json.Marshal(mutexRequestMessage)
	if err != nil {
		logrus.Fatalf("[mutex_send.requestCriticalArea] Error marshal mutexMessage with error %s", err)
	}
	payloadString := string(payload)
	// 4. send all users the mutex message
	for _, user := range identity.Users {
		go sendRequestToUser(user, payloadString)
	}
	// 5. wait for all users to 'reply-ok' to your request
	mutexReceivedAllRequestsMessage := <- mutexReceivedAllRequests
	logrus.Infof("[mutex_send.requestCriticalArea] received all requests, you can now enter the critical area - %s", mutexReceivedAllRequestsMessage)
	// 6. enterCriticalSection() - and leave critical section if this method returns
	enterCriticalSection()
	// exec leaveCriticalSection() to leave
}

/*
sendRequestToUser - send request message to a user
1. create channel and add it to mutexWaitingRequests
2. GO - checkClientIfResponded() start listening and asking back for user availability
3. send POST to user and wait for reply-ok answer
4. receive answer message
5. check if answer message is reply-ok message
6. send message through channel that user responded -> no need to listen anymore
7. add waiting request to responded requests
8. check if all users responded
*/
func sendRequestToUser(user identity.InformationUserDTO, payloadString string) {
	// 1. create channel and add it to mutexWaitingRequests
	waitingRequestChannel := make(chan string)
	waitingRequest := channelUserRequest{
		userEndpoint:     user.Endpoint,
		channel:          waitingRequestChannel,
		user: 		      user,
		sendHealthChecks: true,
	}
	mutexWaitingRequests = append(mutexWaitingRequests, waitingRequest)
	// 2. start listening and asking back for user availability
	go checkClientIfResponded(waitingRequest)
	// 3. send POST to user and wait for reply-ok answer
	res, err := pkg.RequestPOST(user.Endpoint + RouteMutexMessage, payloadString)
	if err != nil {
		logrus.Fatalf("[mutex_send.sendRequestToUser] Error sending POST with error %s", err)
	}
	// 4. receive answer message
	logrus.Infof("[mutex_send.sendRequestToUser] send request messages to user %s", user.UserId)
	var mutexAnswerMessage MessageMutexDTO
	err = json.Unmarshal(res, &mutexAnswerMessage)
	if err != nil {
		logrus.Fatalf("[mutex_send.sendRequestToUser] Error sending POST with error %s", err)
	}
	// 5. check if answer message is reply-ok message
	if mutexAnswerMessage.Msg != ReplyOKMessage {
		logrus.Fatalf("[mutex_send.sendRequestToUser] Error sending POST with error %s", err)
	}
	// 6. send message through channel that user responded -> no need to listen anymore
	waitingRequestChannel <- ReplyOKMessage
	logrus.Infof("[mutex_send.sendRequestToUser] reply-ok message received from user %s", mutexAnswerMessage.User)
	// 7. add waiting request to responded requests
	mutexReceivedRequests = append(mutexReceivedRequests, waitingRequest)
	// 8. check if all users responded
	checkIfAllUsersResponded()
}

/*
checkClientIfResponded - listen if client reply-ok'ed and check with him back if not
1. GO - clientHealthCheck()
2. receiving message
3. if message is not reply-ok
3.1 abroad health checks, user answered
4. if message is something else
4.1 ping user
4.2 wait some time
4.3 if answer
4.4 YES: loopback to 2
4.5.1 NO: delete user
4.5.2 send reply-ok message to waitingRequestChannel
4.5.3 stop hearth beat
*/
func checkClientIfResponded(waitingForUserResponseObj channelUserRequest) {
	logrus.Infof("[mutex_send.checkClientIfResponded] listen if user answers %s", waitingForUserResponseObj.userEndpoint)
	// 1. clientHealthCheck() - sends beats through channel
	go clientHealthCheck(waitingForUserResponseObj)
	for true {
		// 2. receiving message
		msg := <- waitingForUserResponseObj.channel
		// 3. if message is reply-ok
		if msg == ReplyOKMessage {
			// 3.1 abroad health checks, user answered
			break
		} else {
			logrus.Warnf("[mutex_send.checkClientIfResponded] received message: %s", msg)
			var mutexUserStatusResponse StateMutexDTO
			// 4.1 ping user
			pingUser(waitingForUserResponseObj.userEndpoint, RouteMutexState, &mutexUserStatusResponse)
			// 4.2 wait some time
			time.Sleep(waitingTime)
			// 4.3 if answer
			if !checkIfStateObjectIsEmpty(mutexUserStatusResponse) {
				// 4.5.1 NO: delete user
				logrus.Warnf("[mutex_send.checkClientIfResponded] user: %s, did not respond", waitingForUserResponseObj.userEndpoint)
				identity.DeleteUser(waitingForUserResponseObj.user)
				// 4.5.2 send reply-ok message to waitingRequestChannel
				// break waiting -- and send a artificial reply-ok, remove user because of inactivity
				waitingForUserResponseObj.channel <- ReplyOKMessage
				// 4.5.3 stop hearth beat, update waitingForUserResponseObj
				removeChannelUserRequest(waitingForUserResponseObj, mutexWaitingRequests)
				waitingForUserResponseObj = channelUserRequest{
					userEndpoint:     waitingForUserResponseObj.userEndpoint,
					user:             waitingForUserResponseObj.user,
					channel:          waitingForUserResponseObj.channel,
					sendHealthChecks: false,
				}
				mutexWaitingRequests = append(mutexWaitingRequests, waitingForUserResponseObj)
				break
			} else {
				logrus.Infof("[mutex_send.checkClientIfResponded] client is alive and in state: %s", mutexUserStatusResponse.State)
			}
		}
	}
}

/*
clientHealthCheck - send health check to the client after a period of time
1. start loop
2. wait some time
3. check if user answered
4. YES, break, return
5. NO, send none message through waitingRequestChannel
*/
func clientHealthCheck(waitingForUserResponseObj channelUserRequest) {
	logrus.Infof("[mutex_send.clientHealthCheck] listen if user answers %s", waitingForUserResponseObj.userEndpoint)
	// 1. start loop
	for true {
		// 2. wait some time
		time.Sleep(waitingTime)
		// 3. check if user answered
		// 4. update waitingForUserResponseObj to know if
		waitingForUserResponseObj = getWaitingForUserResponseObj(waitingForUserResponseObj)
		if waitingForUserResponseObj.sendHealthChecks {
			// 4. YES, break, return
			break
		} else {
			// 5. NO, send none message through waitingRequestChannel
			waitingForUserResponseObj.channel <- "none"
		}
	}
}

// --------------------
// HELPER METHODS

/*
pingUser - request user mutex state to tell if he is alive
*/
func pingUser(userEndpoint string, mutexStateEndpoint string, mutexUserStatus *StateMutexDTO) {
	*mutexUserStatus = RequestUserState(userEndpoint, mutexStateEndpoint)
}

/*
getWaitingForUserResponseObj - return waitingForUserResponseObj from list mutexWaitingRequests; which might changed
 */
func getWaitingForUserResponseObj(waitingForUserResponseObj channelUserRequest) channelUserRequest {
	for _, waitingRequest := range mutexWaitingRequests {
		if waitingRequest.userEndpoint == waitingForUserResponseObj.userEndpoint {
			return waitingRequest
		}
	}
	logrus.Fatalf("[mutex_send.getWaitingForUserResponseObj] all users answered, send notification to waiting task")
	return waitingForUserResponseObj
}

/*
checkIfAllUsersResponded - check if all users responded with reply-ok after requesting critical section
1. clean types mutexReceivedRequests, mutexWaitingRequests
2. notify channel mutexReceivedAllRequests
3. clean up both list to reply to
*/
func checkIfAllUsersResponded() {
	if len(mutexWaitingRequests) == len(mutexReceivedRequests) {
		// 1. clean types mutexReceivedRequests, mutexWaitingRequests
		mutexReceivedRequests = []channelUserRequest{}
		mutexWaitingRequests = []channelUserRequest{}
		logrus.Infof("[mutex_send.checkIfAllUsersResponded] all users answered, send notification to waiting task")
		// 2. notify channel mutexReceivedAllRequests
		mutexReceivedAllRequests <- "good to go"
		// 3. clean up both list to reply to
		mutexWaitingRequests = []channelUserRequest{}
		mutexReceivedRequests = []channelUserRequest{}
	}
}

/*
- TODO description
 */
func checkIfStateObjectIsEmpty(mutexUserState StateMutexDTO) bool {
	var emptyMutexUserState StateMutexDTO
	return mutexUserState == emptyMutexUserState
}
