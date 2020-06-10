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
	case ReplyOKMessage: receivedReplyMessage(mutexMessage, &mutexMessageResponse)
	default: logrus.Warningf("[mutex.receiveMutexMessage] message: %s, could not be identified", mutexMessage.Msg)
	}
	// send a reply-ok message - increase clock
	incrementClock(mutexMessage.Time)
	logrus.Infof("[mutex.receiveMutexMessage] send response message")
	return mutexMessageResponse
}

/*
receivedRequestMessage - received a 'request' message
- in critical section (stage held) -> store request and send reply-ok as soon as leaving critical section
- in idle (stage released) -> send reply-ok
- waiting (stage wanting) -> compare clocks and store request if yours is lower (if you are not the one waiting)
 */
func receivedRequestMessage(mutexMessage MessageMutexDTO, mutexResponseMessage *MessageMutexDTO) {
	switch state {

	// STATE RELEASED
	case StateReleased: {
		// send reply-ok
		logrus.Infof("[mutex.receivedRequestMessage] state: %s identified, setting mutexResponseMessage", StateReleased)
		*mutexResponseMessage = getReplyOkMessage()
	}

	// STATE HELD
	case StateHeld: {
		// store request and send reply-ok as soon as leaving critical section
		logrus.Infof("[mutex.receivedRequestMessage] state: %s identified, add mutexResponseMessage to mutexSendRequests list", StateHeld)
		channel := make(chan string)
		mutexSendRequests = append(mutexSendRequests, channelUserRequest{
			userEndpoint: mutexMessage.User,
			channel:      channel,
		})
		// wait until it is allowed to send a reply-ok
		logrus.Infof("[mutex.receivedRequestMessage] wait until it is allowed to send a reply-ok to: %s", mutexMessage.User)
		msg := <- channel
		logrus.Infof("[mutex.receivedRequestMessage] received %s from channel", msg)
		*mutexResponseMessage = getReplyOkMessage()
	}

	// STATE WANTING
	case StateWanting: {
		// compare clocks and store request if yours is lower (if you are not the one waiting)
		logrus.Infof("[mutex.receivedRequestMessage] state: %s identified, setting mutexResponseMessage", StateWanting)
		// you are in the waiting state and therefore sending a response-ok message to your request
		if mutexMessage.User == identity.YourUserInformation.Endpoint {
			*mutexResponseMessage = getReplyOkMessage()
		} else {
			// higher clock goes first
			if mutexMessage.Time > clock {
				*mutexResponseMessage = getReplyOkMessage()
			} else {
				channel := make(chan string)
				mutexSendRequests = append(mutexSendRequests, channelUserRequest{
					userEndpoint: mutexMessage.User,
					channel:      channel,
				})
				// wait until it is allowed to send a reply-ok
				logrus.Infof("[mutex.receivedRequestMessage] wait until it is allowed to send a reply-ok to: %s", mutexMessage.User)
				msg := <- channel
				logrus.Infof("[mutex.receivedRequestMessage] received %s from channel", msg)
				*mutexResponseMessage = getReplyOkMessage()
			}
		}
	}
	default: logrus.Fatalf("[mutex.receivedRequestMessage] state: %s, could not be identified", state)
	}
}

/*
getReplyOkMessage - Helper Method
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
receivedReplyMessage - received a 'reply-ok' message
1. notify channel in list
- go all required reply-ok messages -> you may enter the critical area
- else wait
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
	// TODO
}

/*
checkIfAllUsersResponded - Helper Method
1. clean types mutexReceivedRequests, mutexWaitingRequests
2. notify channel mutexReceivedAllRequests
 */
func checkIfAllUsersResponded() {
	if len(mutexWaitingRequests) == len(mutexReceivedRequests) {
		// 1. clean types mutexReceivedRequests, mutexWaitingRequests
		mutexReceivedRequests = []channelUserRequest{}
		mutexWaitingRequests = []channelUserRequest{}
		// 2. notify channel mutexReceivedAllRequests
		logrus.Infof("[mutex.checkIfAllUsersResponded] all users answered, send notification to waiting task")
		mutexReceivedAllRequests <- "good to go" // TODO maybe replace this with a more interesting message
	}
}

/*
requestCriticalArea - tell all users that this user wants to enter the critical section
1. increment clock
2. create mutex message
3. send all users the mutex message
4. TODO
5. wait for all users to reply-ok to your request
6. update state to 'held'
 */
func requestCriticalArea() {
	logrus.Infof("[mutex.requestCriticalArea] starting")
	// 1. increment clock
	incrementClock(clock)
	// 2. create mutex message
	var mutexMessage = MessageMutexDTO{
		Msg:   RequestMessage,
		Time:  clock,
		Reply: mutexYourReply,
		User:  mutexYourUser,
	}
	payload, err := json.Marshal(mutexMessage)
	if err != nil {
		logrus.Fatalf("[mutex.requestCriticalArea] Error marshal mutexMessage with error %s", err)
	}
	payloadString := string(payload)
	for _, user := range identity.Users {
		go sendRequestToUser(user, payloadString)
	}
	// 5. wait for all users to 'reply-ok' to your request
	mutexReceivedAllRequestsMessage := <- mutexReceivedAllRequests
	// you can now enter the critical area
	logrus.Infof("[mutex.requestCriticalArea] received all requests - %s", mutexReceivedAllRequestsMessage)
	// 6. update state to 'held'
	state = StateHeld
	// 7. clear receiveRequest channels
}

/*
sendRequestToUser - send request message to a user
1. create channel and add it to mutexWaitingRequests
2. prepare
 */
func sendRequestToUser(user identity.InformationUserDTO, payloadString string) {
	// mutexReceivedRequests
	channel := make(chan string)
	mutexWaitingRequests = append(mutexWaitingRequests, channelUserRequest{
		userEndpoint: user.Endpoint,
		channel:      channel,
	})
	// 3. send all users the mutex message
	res, err := pkg.RequestPOST(user.Endpoint + RouteMutexMessage, payloadString) // TODO listen if user answers (like in election)
	if err != nil {
		logrus.Warningf("[mutex.sendRequestToUser] Error sending POST with error %s", err)
	}
	logrus.Infof("[mutex.sendRequestToUser] send request messages to user %s", user.UserId)
	var mutexMessage MessageMutexDTO
	err = json.Unmarshal(res, &mutexMessage)
	if err != nil {
		logrus.Fatalf("[mutex.sendRequestToUser] Error sending POST with error %s", err)
	}
	receiveMutexMessage(mutexMessage)
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

// message to channel to identify users
type channelUserRequest struct {
	userEndpoint string
	channel chan string
}
