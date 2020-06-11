package mutex

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"goBully/internal/identity"
	"goBully/pkg"
)

/* METHODS overview:
	- receiveMutexMessage()       // get a message from a api (election, coordinator)
	- receivedRequestMessage()    // handle incoming request message
	- receivedReplyMessage()      // handle incoming reply message
      ---------------------
	TODO
*/

// all requests that are currently on hold and shall receive a reply-ok answer (string - ENDPOINT)
var mutexSendRequests []channelUserRequest
// all requests where you waiting for an answer (string - ENDPOINT)
var mutexWaitingRequests []channelUserRequest
var mutexReceivedRequests []channelUserRequest
// send message through this channel if you received all requests
var mutexReceivedAllRequests = make(chan string)

/*
receiveMutexMessage - map logic after message
 */
func receiveMutexMessage(mutexMessage MessageMutexDTO) MessageMutexDTO {
	logrus.Infof("[mutex.receiveMutexMessage] received message")
	// received a request message
	incrementClock(mutexMessage.Time)

	// response is set in receivedRequestMessage && receivedReplyMessage
	var mutexMessageResponse MessageMutexDTO

	switch mutexMessage.Msg {
	case RequestMessage: receivedRequestMessage(mutexMessage, &mutexMessageResponse)
	// TODO case not necessary? - fallback
	case ReplyOKMessage: receivedReplyMessage(mutexMessage, &mutexMessageResponse)
	default: logrus.Warningf("[mutex.receiveMutexMessage] message: %s, could not be identified", mutexMessage.Msg)
	}
	// send a reply-ok message - increase clock
	incrementClock(mutexMessage.Time)
	logrus.Infof("[mutex.receiveMutexMessage] send response message")
	return mutexMessageResponse
}

// --------------------
// RECEIVE REQUEST

/*
receivedRequestMessage - received a 'request' message
- in critical section (stage held) -> store request and send reply-ok as soon as leaving critical section
- in idle (stage released) -> send reply-ok
- waiting (stage wanting) -> compare clocks and store request if yours is lower (if you are not the one waiting)
 */
func receivedRequestMessage(mutexMessage MessageMutexDTO, mutexResponseMessage *MessageMutexDTO) {
	logrus.Infof("[mutex.receivedRequestMessage] state: %s", state)

	switch state {
	case StateReleased: *mutexResponseMessage = getReplyOkMessage()
	case StateHeld: *mutexResponseMessage = receiveMessageHeld(mutexMessage)
	case StateWanting: *mutexResponseMessage = receiveMessageWanting(mutexMessage)
	default: logrus.Fatalf("[mutex.receivedRequestMessage] state: %s, could not be identified", state)
	}

	logrus.Infof("[mutex.receivedRequestMessage] return response, current state: %s", state)
}

/*
receiveMessageHeld - received a request message, state: held
store request and send reply-ok as soon as leaving critical section
*/
func receiveMessageHeld(mutexMessage MessageMutexDTO) MessageMutexDTO {
	requestChannel := make(chan string)
	mutexSendRequests = append(mutexSendRequests, channelUserRequest{
		userEndpoint: mutexMessage.User,
		channel:      requestChannel,
	})
	// wait until it is allowed to send a reply-ok
	logrus.Infof("[mutex.receiveMessageHeld] wait until it is allowed to send a reply-ok to: %s", mutexMessage.User)
	msg := <- requestChannel
	logrus.Infof("[mutex.receiveMessageHeld] received %s from channel", msg)
	return getReplyOkMessage()
}

/*
receiveMessageWanting - received a request message, state: wanting
compare clocks and store request if yours is lower (if you are not the one waiting)
*/
func receiveMessageWanting(mutexMessage MessageMutexDTO) MessageMutexDTO {
	// you are in the waiting state and therefore sending a response-ok message to your request
	if mutexMessage.User != identity.YourUserInformation.Endpoint {
		// higher clock goes first
		if mutexMessage.Time < clock {
			requestChannel := make(chan string)
			mutexSendRequests = append(mutexSendRequests, channelUserRequest{
				userEndpoint: mutexMessage.User,
				channel:      requestChannel,
			})
			// wait until it is allowed to send a reply-ok
			logrus.Infof("[mutex.receiveMessageWanting] wait until it is allowed to send a reply-ok to: %s", mutexMessage.User)
			msg := <- requestChannel
			logrus.Infof("[mutex.receiveMessageWanting] received %s from channel", msg)
		}
		// else return reply-ok
	}
	// else return reply-ok
	return getReplyOkMessage()
}

// --------------------
// RECEIVE REPLY-OK TODO method not necessary
/*
receivedReplyMessage - received a 'reply-ok' message
1. notify channel in list
- go all required reply-ok messages -> you may enter the critical area
- else wait
TODO method not necessary
*/
func receivedReplyMessage(mutexMessage MessageMutexDTO, mutexResponseMessage *MessageMutexDTO) {
	logrus.Infof("[mutex.receivedReplyMessage] user: %s", mutexMessage.User)
	for _, userCallback := range mutexWaitingRequests {
		if userCallback.userEndpoint == mutexMessage.User {
			userCallback.channel <- mutexMessage.Reply // TODO maybe send something else?
			logrus.Infof("[mutex.receivedReplyMessage] send message through user channel: %s", mutexMessage.User)
			mutexReceivedRequests = append(mutexReceivedRequests, userCallback)
			checkIfAllUsersResponded()
			break
		}
	}
	// ...
}

// --------------------
// SEND REQUEST MESSAGE

/*
requestCriticalArea - tell all users that this user wants to enter the critical section
1. set state to 'wanting'
2. increment clock, you are about to send messages
3. create a request mutex message
4. GO - send all users the mutex message
5. wait for all users to reply-ok to your request
6. enterCriticalSection() - and leave critical section if this method returns
7. leaveCriticalSection()
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
	// 7. leaveCriticalSection()
	leaveCriticalSection()
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
	go checkClientIfResponded(waitingRequestChannel)
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
4.2 loopback to 2.
 */
func checkClientIfResponded(waitingChannel chan string) {
	logrus.Infof("[mutex.checkClientIfResponded] listen if user answers %s", mutexAnswerMessage.User)
	// 1. clientHealthCheck()
	go clientHealthCheck(waitingChannel)
	for true {
		// 2. receiving message
		msg := <- waitingRequestChannel
		// 3. if message is reply-ok
		if msg == ReplyOKMessage {
			break
		}
	}
	// TODO
}

/*
clientHealthCheck - send health check to the client after a period of time
 */
func clientHealthCheck(waitingChannel chan string) {
	logrus.Infof("[mutex.clientHealthCheck] listen if user answers %s", mutexAnswerMessage.User)
	// TODO
}

/*
enterCriticalSection - enter critical section
1. update state to 'held'
 */
func enterCriticalSection() {
	state = StateHeld
	logrus.Infof("[mutex.enterCriticalSection] enter, state: %s", state)
	// TODO do something
	// leave critical section if this method returns
}

/*
leaveCriticalSection - enter critical section
1. update state to 'released'
2. notify waiting users
*/
func leaveCriticalSection() {
	state = StateReleased
	logrus.Infof("[mutex.leaveCriticalSection] leave, state: %s", state)
	// TODO
}

// --------------------
// HELPER METHODS

/*
checkIfAllUsersResponded - check if all users responded with reply-ok after requesting critical section
1. clean types mutexReceivedRequests, mutexWaitingRequests
2. notify channel mutexReceivedAllRequests
TODO
*/
func checkIfAllUsersResponded() {
	if len(mutexWaitingRequests) == len(mutexReceivedRequests) {
		// 1. clean types mutexReceivedRequests, mutexWaitingRequests
		mutexReceivedRequests = []channelUserRequest{}
		mutexWaitingRequests = []channelUserRequest{}
		logrus.Infof("[mutex.checkIfAllUsersResponded] all users answered, send notification to waiting task")
		// 2. notify channel mutexReceivedAllRequests
		mutexReceivedAllRequests <- "good to go" // TODO maybe replace this with a more interesting message
	}
	// TODO check if this is ok
}

/*
getReplyOkMessage - return reply-ok message
*/
func getReplyOkMessage() MessageMutexDTO {
	return MessageMutexDTO{
		Msg:   ReplyOKMessage,
		Time:  clock,
		Reply: mutexYourReply,
		User:  mutexYourUser,
	}
}

/*
incrementClock - increase local lamport lock
 */
func incrementClock(i int32) int32 {
	clock = max(clock, i)
	logrus.Infof("[mutex.incrementClock] increasing local clock to %s", clock)
	return clock
}

/*
max - simple max function with int32 types
 */
func max(i int32, j int32) int32 {
	if i > j {
		return i
	} else {
		return j
	}
}

// --------------------
// PRIVATE TYPES

// message to channel to identify users
type channelUserRequest struct {
	userEndpoint string
	channel chan string
}
