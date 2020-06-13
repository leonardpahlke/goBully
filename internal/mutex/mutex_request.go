package mutex

import (
	"goBully/internal/identity"

	"github.com/sirupsen/logrus"
)

/*
receivedRequestMessage - received a 'request' message
- in critical section (stage held) -> store request and send reply-ok as soon as leaving critical section
- in idle (stage released) -> send reply-ok
- waiting (stage wanting) -> compare clocks and store request if yours is lower (if you are not the one waiting)
*/
func receivedRequestMessage(mutexMessage MessageMutexEntity) {
	logrus.Infof("[mutex_receive.receivedRequestMessage] state: %s", state)

	switch state {
	case StateReleased:
		sendReplyOkMessage(mutexMessage.User + mutexMessage.Reply)
	case StateHeld:
		receiveMessageHeld(mutexMessage)
	case StateWanting:
		receiveMessageWanting(mutexMessage)
	default:
		logrus.Fatalf("[mutex_receive.receivedRequestMessage] state: %s, could not be identified", state)
	}

	logrus.Infof("[mutex_receive.receivedRequestMessage] finished processing request message")
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
			logrus.Infof("[mutex_receive.receiveMessageWanting] my clock is > mutexMessage.Time")
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
*/
func waitingForSendingAnswerBack(mutexMessage MessageMutexEntity) {
	requestChannel := make(chan string)
	userSendingChan := userSendingChannel{
		userEndpoint: mutexMessage.User,
		channel:      requestChannel,
	}
	replyOkSendingList = append(replyOkSendingList, userSendingChan)
	// wait until it is allowed to send a reply-ok
	logrus.Infof("[mutex_receive.waitingForSendingAnswerBack] wait until it is allowed to send a reply-ok to: %s", mutexMessage.User)
	msg := <-requestChannel
	logrus.Infof("[mutex_receive.waitingForSendingAnswerBack] received: %s from channel", msg)
	// remove requestChannelInfo form mutexSendRequests
	rmReplyOkSendingUser(userSendingChan)

}

/*
getReplyOkMessage - return reply-ok message
*/
func getReplyOkMessage() MessageMutexEntity {
	return MessageMutexEntity{
		Msg:   ReplyOKMessage,
		Time:  clock,
		Reply: mutexYourReply,
		User:  mutexYourUser,
	}
}

/*
STRUCTS
*/

// channel to manage user reponses
type userSendingChannel struct {
	userEndpoint string
	channel      chan string
}
