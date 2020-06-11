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
	logrus.Infof("[mutex.requestCriticalArea] starting, state: %s", state)
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
		logrus.Fatalf("[mutex.requestCriticalArea] Error marshal mutexMessage with error %s", err)
	}
	payloadString := string(payload)
	// 4. send all users the mutex message
	for _, user := range identity.Users {
		go sendRequestToUser(user, payloadString)
	}
	// 5. wait for all users to 'reply-ok' to your request
	mutexReceivedAllRequestsMessage := <- mutexReceivedAllRequests
	logrus.Infof("[mutex.requestCriticalArea] received all requests, you can now enter the critical area - %s", mutexReceivedAllRequestsMessage)
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
		userEndpoint: user.Endpoint,
		channel:      waitingRequestChannel,
	}
	mutexWaitingRequests = append(mutexWaitingRequests, waitingRequest)
	// 2. start listening and asking back for user availability
	go checkClientIfResponded(waitingRequest)
	// 3. send POST to user and wait for reply-ok answer
	res, err := pkg.RequestPOST(user.Endpoint + RouteMutexMessage, payloadString)
	if err != nil {
		logrus.Fatalf("[mutex.sendRequestToUser] Error sending POST with error %s", err)
	}
	// 4. receive answer message
	logrus.Infof("[mutex.sendRequestToUser] send request messages to user %s", user.UserId)
	var mutexAnswerMessage MessageMutexDTO
	err = json.Unmarshal(res, &mutexAnswerMessage)
	if err != nil {
		logrus.Fatalf("[mutex.sendRequestToUser] Error sending POST with error %s", err)
	}
	// 5. check if answer message is reply-ok message
	if mutexAnswerMessage.Msg != ReplyOKMessage {
		logrus.Fatalf("[mutex.sendRequestToUser] Error sending POST with error %s", err)
	}
	// 6. send message through channel that user responded -> no need to listen anymore
	waitingRequestChannel <- ReplyOKMessage
	logrus.Infof("[mutex.sendRequestToUser] reply-ok message received from user %s", mutexAnswerMessage.User)
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
*/
func checkClientIfResponded(userEndpoint string, userMutexStateEndpoint string, waitingResponse channelUserRequest) {
	logrus.Infof("[mutex.checkClientIfResponded] listen if user answers %s", waitingResponse.userEndpoint)
	// 1. clientHealthCheck()
	go clientHealthCheck(waitingResponse)
	for true {
		// 2. receiving message
		msg := <- waitingRequestChannel
		// 3. if message is reply-ok
		if msg == ReplyOKMessage {
			// 3.1 abroad health checks, user answered
			break
		} else {
			logrus.Warnf("[mutex.checkClientIfResponded] received message %s", msg)
			var mutexUserStatus StateMutexDTO
			// 4.1 ping user
			pingUser(userEndpoint, userMutexStateEndpoint, &mutexUserStatus)
			// 4.2 wait some time
			time.Sleep(waitingTime)
			// 4.3 if answer
			// 4.4 YES: loopback to 2
			// 4.5.1 NO: delete user
			// 4.5.2 NO: send reply-ok message to waitingRequestChannel

			// TODO
			// 4.2 loopback to 2.
		}
	}
	// TODO
}

/*
pingUser - request user mutex state to tell if he is alive
*/
func pingUser(userEndpoint string, mutexStateEndpoint string, mutexUserStatus *StateMutexDTO) {
	*mutexUserStatus = RequestUserState(userEndpoint, mutexStateEndpoint)
}

/*
clientHealthCheck - send health check to the client after a period of time
1. start loop
2. wait some time
3. check if user answered
4. YES, break, return
5. NO, send none message through waitingRequestChannel
*/
func clientHealthCheck(waitingResponse channelUserRequest) {
	logrus.Infof("[mutex.clientHealthCheck] listen if user answers %s", mutexAnswerMessage.User)
	// 1. start loop
	for true {
		// 2. wait some time
		time.Sleep(waitingTime)
		// 3. check if user answered
		if waitingResponse.sendHealthChecks { // TODO update condition
			// 4. YES, break, return
			break
		} else {
			// 5. NO, send none message through waitingRequestChannel
			waitingResponse.channel <- "none"
		}
	}
}
