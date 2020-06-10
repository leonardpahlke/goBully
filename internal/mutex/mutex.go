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
var mutexSendRequests []string
// all requests where you waiting for an answer (string - ENDPOINT)
var mutexReceiveRequests []string

/*
receiveMutexMessage - get a mutex message and map logic
 */
func receiveMutexMessage(mutexMessage MessageMutexDTO) MessageMutexDTO {
	// received message -> increment time
	incrementClock(mutexMessage.Time)

	// response is set in receivedRequestMessage && receivedReplyMessage
	var mutexMessageResponse MessageMutexDTO

	switch mutexMessage.Msg {
	case RequestMessage: receivedRequestMessage(mutexMessage, &mutexMessageResponse)
	case ReplyOKMessage: receivedReplyMessage(mutexMessage, &mutexMessageResponse)
	default: logrus.Warningf("[mutex.receiveMutexMessage] message: %s, could not be identified", mutexMessage.Msg)
	}
	incrementClock(mutexMessage.Time)
	return mutexMessageResponse
}

/*
receivedRequestMessage - received a 'request' message
- in critical section (stage held) -> store request and send reply-ok as soon as leaving critical section
- in idle (stage released) -> send reply-ok
- waiting (stage wanting) -> compare clocks and store request if yours is lower (if you are not the one waiting)
 */
func receivedRequestMessage(mutexMessage MessageMutexDTO, mutexResponseMessage *MessageMutexDTO) {
	var replyOkMessage = MessageMutexDTO{
		Msg:   ReplyOKMessage,
		Time:  clock,
		Reply: mutexYourReply,
		User:  mutexYourUser,
	}
	switch state {
	// send reply-ok
	case StateReleased: {
		logrus.Infof("[mutex.receivedRequestMessage] state: %s identified, setting mutexResponseMessage", StateReleased)
		*mutexResponseMessage = replyOkMessage
	}
	// store request and send reply-ok as soon as leaving critical section
	case StateHeld: {
		logrus.Infof("[mutex.receivedRequestMessage] state: %s identified, add mutexResponseMessage to mutexSendRequests list", StateHeld)
		mutexSendRequests = append(mutexSendRequests, mutexMessage.User)
		// TODO stall ?
	}
	// compare clocks and store request if yours is lower (if you are not the one waiting)
	case StateWanting: {
		logrus.Infof("[mutex.receivedRequestMessage] state: %s identified, setting mutexResponseMessage", StateWanting)
		// you are in the waiting state and therefore sending a response-ok message to your request
		if mutexMessage.User == identity.YourUserInformation.Endpoint {
			*mutexResponseMessage = replyOkMessage
		} else {
			// higher clock goes first
			if mutexMessage.Time > clock {
				*mutexResponseMessage = replyOkMessage
			} else {
				mutexSendRequests = append(mutexSendRequests, mutexMessage.User)
				// TODO stall ?
			}
		}
	}
	default: logrus.Panicf("[mutex.receivedRequestMessage] state: %s, could not be identified", state)
	}
}

/*
receivedReplyMessage - received a 'reply-ok' message
1. remove
- go all required reply-ok messages -> you may enter the critical area
- else wait
*/
func receivedReplyMessage(mutexMessage MessageMutexDTO, mutexResponseMessage *MessageMutexDTO) {
	// TODO
}

/*
requestCriticalArea - tell all users that this user wants to enter the critical section
1. increment clock
2. create mutex message
3. send all users the mutex message
4. TODO
5. update state
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
	for _, user := range identity.Users {
		mutexReceiveRequests = append(mutexReceiveRequests, user.Endpoint)
		// 3. send all users the mutex message
		_, err := pkg.RequestPOST(user.Endpoint + RouteMutexMessage, string(payload)) // TODO listen if user answers (like in election)
		if err != nil {
			logrus.Fatalf("[mutex.sendMessage] Error sending POST with error %s", err)
		}
		logrus.Infof("[mutex.requestCriticalArea] send request messages to user %s", user.UserId)
	}
	// TODO ...
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
