package mutex

import (
	"encoding/json"
	"goBully/internal/identity"
	"goBully/pkg"
	"time"

	"github.com/sirupsen/logrus"
)

/*
requestCriticalArea - tell all users that this user wants to enter the critical section
1. set state to 'wanting'
2. increment clock, you are about to send mutex-messages
3. create a request mutex-message
4. create a response channel for every user (including yourself)
5. create new object to manage responses of this request
6. add new requestResponseChannel to replyOkwaitingList
7. GO - send all users the request mutex-message
8. wait for all users to reply-ok to your request
9. remove the waiting reponses object from the list
10. enterCriticalSection()
*/
func requestCriticalArea() {
	// 1. set state to 'wanting'
	state = StateWanting
	logrus.Infof("[mutex_send.requestCriticalArea] starting, state: %s", state)
	// 2. increment clock, you are about to send mutex-messages
	incrementClock(clock)
	// 3. create a request mutex-message
	var mutexRequestMessage = MessageMutexEntity{
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

	// 4. create a response channel for every user (including yourself)
	var userReponseChannels []userReponseChannel
	for _, user := range identity.Users {
		userReponseChannels = append(userReponseChannels, userReponseChannel{
			user:    user,
			channel: make(chan string),
		})
	}

	// 5. create new object to manage responses of this request
	requestResponseChannel := responseChannel{
		replyOkReceivingList: userReponseChannels,
		allReplyOkReceived:   make(chan string),
	}

	// 6. add new requestResponseChannel to replyOkwaitingList
	replyOkwaitingList = append(replyOkwaitingList, requestResponseChannel)

	// 7. GO - send all users the request mutex-message
	for _, userChannel := range userReponseChannels {
		go sendRequestToUser(userChannel, payloadString)
	}
	// 8. wait for all users to reply-ok to your request
	receivedAllRequestsMessage := <-requestResponseChannel.allReplyOkReceived

	logrus.Infof("[mutex_send.requestCriticalArea] received all requests, you can now enter the critical area - %s", receivedAllRequestsMessage)
	// 9. remove the waiting reponses object from the list
	removeFirstWaitingTask()

	// 10. enterCriticalSection()
	enterCriticalSection()
	// exec LeaveCriticalSection() to leave
}

/*
sendRequestToUser - send request message to a user
1. send POST to user and wait for reply-ok answer
2. start checking if user answered
*/
func sendRequestToUser(userResChannel userReponseChannel, payloadString string) {
	// 1. send POST to user and wait for reply-ok answer
	_, err := pkg.RequestPOST(userResChannel.user.Endpoint+RouteMutexMessage, payloadString)
	if err != nil {
		logrus.Fatalf("[mutex_send.sendRequestToUser] Error sending POST with error %s", err)
	}
	logrus.Infof("[mutex_send.sendRequestToUser] send request messages to user %s", userResChannel.user.UserId)
	// 2. start checking if user answered
	checkClientIfResponded(userResChannel)
}

/*
checkClientIfResponded - listen if client reply-ok'ed and check with him back if not
1. GO - clientHealthCheck()
2. receiving message
3. if message is reply-ok, return
4 ping user
5 wait some time
6 if answered: loopback to 2
7 if not answered, delete user
*/
func checkClientIfResponded(userResChannel userReponseChannel) {
	logrus.Infof("[mutex_send.checkClientIfResponded] listen if user answers %s", userResChannel.user.Endpoint)
	// 1. clientHealthCheck() - sends beats through channel
	go clientHealthCheck(userResChannel)
	for true {
		// 2. receiving message
		msg := <-userResChannel.channel
		// 3. if message is reply-ok, return
		if msg == ReplyOKMessage {
			return
		}

		// user still did not answered - pinging user to check if he is alive

		logrus.Warnf("[mutex_send.checkClientIfResponded] received message: %s", msg)
		var mutexUserStatusResponse StateMutexEntity
		// 4 ping user
		pingUser(userResChannel.user.Endpoint, RouteMutexState, &mutexUserStatusResponse)
		// 5 wait some time
		time.Sleep(waitingTime)
		if checkIfStateObjectIsEmpty(mutexUserStatusResponse) {
			// 6 if answered: loopback to 2
			logrus.Infof("[mutex_send.checkClientIfResponded] client is alive and in state: %s", mutexUserStatusResponse.State)
		} else {
			logrus.Warnf("[mutex_send.checkClientIfResponded] user: %s, did not respond", userResChannel.user.Endpoint)
			// 7 delete user - inactive
			identity.DeleteUser(userResChannel.user)
			return
		}
	}
}

/*
clientHealthCheck - send health check to the client after a period of time
1. start loop
2. wait some time
3. NO, send none message through userResChannel.channel
!!! this method returns if checkClientIfResponded() returns!!!
*/
func clientHealthCheck(userResChannel userReponseChannel) {
	logrus.Infof("[mutex_send.clientHealthCheck] listen if user answers endpoint: %s", userResChannel.user.Endpoint)
	// 1. start loop
	for true {
		// 2. wait some time
		time.Sleep(waitingTime)
		// 3. NO, send none message through userResChannel.channel
		userResChannel.channel <- "none"
	}
}

// --------------------
// HELPER METHODS

/*
pingUser - request user mutex state to tell if he is alive
*/
func pingUser(userEndpoint string, mutexStateEndpoint string, mutexUserStatus *StateMutexEntity) {
	*mutexUserStatus = RequestUserState(userEndpoint, mutexStateEndpoint)
}

/*
checkIfStateObjectIsEmpty - return whether the StateMutexDTO is empty
*/
func checkIfStateObjectIsEmpty(mutexUserState StateMutexEntity) bool {
	var emptyMutexUserState StateMutexEntity
	return mutexUserState == emptyMutexUserState
}

/*
removeFirstWaitingTask - remove first task from waiting list
*/
func removeFirstWaitingTask() {
	i := 0
	copy(replyOkwaitingList[i:], replyOkwaitingList[i+1:])              // Shift a[i+1:] left one index.
	replyOkwaitingList[len(replyOkwaitingList)-1] = responseChannel{}   // Erase last element (write zero value).
	replyOkwaitingList = replyOkwaitingList[:len(replyOkwaitingList)-1] // Truncate slice.
}

/*
STRUCTS
*/

// message to channel to identify users
type responseChannel struct {
	replyOkReceivingList []userReponseChannel
	allReplyOkReceived   chan string
}

// channel to manage user reponses
type userReponseChannel struct {
	user    identity.InformationUserDTO
	channel chan string
}
