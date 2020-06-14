package mutex

import (
	"goBully/internal/identity"

	"github.com/sirupsen/logrus"
)

/*
receivedRequestMessage - received a 'request' message
- in idle (stage: released) -> send reply-ok
- in critical section (stage: held) -> store request and send reply-ok as soon as leaving critical section
- waiting (stage: wanting) -> compare clocks and store request if yours is lower (if you are not the one waiting)
*/
func receivedRequestMessage(mutexMessage MessageMutexEntity) {
	logrus.Infof("[mutex_request.receivedRequestMessage] state: %s", state)

	switch state {
	case StateReleased:
		sendReplyOkMessage(mutexMessage.User + mutexMessage.Reply)
	case StateHeld:
		receiveMessageHeld(mutexMessage)
	case StateWanting:
		receiveMessageWanting(mutexMessage)
	default:
		logrus.Fatalf("[mutex_request.receivedRequestMessage] state: %s, could not be identified", state)
	}
}

/*
receiveMessageHeld - received a request message, your state: held
store request and send reply-ok as soon as leaving critical section
*/
func receiveMessageHeld(mutexMessage MessageMutexEntity) {
	waitingForSendingAnswerBack(mutexMessage)
	sendReplyOkMessage(mutexMessage.User + mutexMessage.Reply)
}

/*
receiveMessageWanting - received a request message, your state: wanting
compare clocks and store request if yours is lower (if you are not the one waiting)
*/
func receiveMessageWanting(mutexMessage MessageMutexEntity) {
	// you are in the waiting state and therefore sending a response-ok message to your request
	if mutexMessage.User != identity.YourUserInformation.Endpoint {
		// higher clock goes first
		if mutexMessage.Time < clock {
			logrus.Infof("[mutex_request.receiveMessageWanting] my clock is > mutexMessage.Time")
			waitingForSendingAnswerBack(mutexMessage)
		}
		// else return reply-ok
	}
	// else return reply-ok
	sendReplyOkMessage(mutexMessage.User + mutexMessage.Reply)
}

// --------------------
// HELPER METHODS

/*
waitingForSendingAnswerBack - add request to list of requests to answer back to
1. add request to replyOkSendingList
2. wait until it is allowed to send a reply-ok
3. remove requestChannelInfo form replyOkSendingList
*/
func waitingForSendingAnswerBack(mutexMessage MessageMutexEntity) {
	requestChannel := make(chan string)
	userSendingChan := userSendingChannel{
		userEndpoint: mutexMessage.User,
		channel:      requestChannel,
	}

	// 1. add request to replyOkSendingList
	replyOkSendingList = append(replyOkSendingList, userSendingChan)

	// 2. wait until it is allowed to send a reply-ok
	logrus.Infof("[mutex_request.waitingForSendingAnswerBack] wait until it is allowed to send a reply-ok to: %s", mutexMessage.User)
	msg := <-requestChannel
	logrus.Infof("[mutex_request.waitingForSendingAnswerBack] received: %s from channel", msg)

	// 3. remove requestChannelInfo form replyOkSendingList
	rmReplyOkSendingUser(userSendingChan)

}

/*
STRUCTS
*/

// channel to manage user reponses
type userSendingChannel struct {
	userEndpoint string
	channel      chan string
}
