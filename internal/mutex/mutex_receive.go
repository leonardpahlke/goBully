package mutex

import (
	"github.com/sirupsen/logrus"
	"goBully/internal/identity"
)

/*
receivedRequestMessage - received a 'request' message
- in critical section (stage held) -> store request and send reply-ok as soon as leaving critical section
- in idle (stage released) -> send reply-ok
- waiting (stage wanting) -> compare clocks and store request if yours is lower (if you are not the one waiting)
*/
func receivedRequestMessage(mutexMessage MessageMutexDTO, mutexResponseMessage *MessageMutexDTO) {
	logrus.Infof("[mutex_receive.receivedRequestMessage] state: %s", state)

	switch state {
	case StateReleased: *mutexResponseMessage = getReplyOkMessage()
	case StateHeld: *mutexResponseMessage = receiveMessageHeld(mutexMessage)
	case StateWanting: *mutexResponseMessage = receiveMessageWanting(mutexMessage)
	default: logrus.Fatalf("[mutex_receive.receivedRequestMessage] state: %s, could not be identified", state)
	}

	logrus.Infof("[mutex_receive.receivedRequestMessage] return response, current state: %s", state)
}

/*
receiveMessageHeld - received a request message, your state: held
store request and send reply-ok as soon as leaving critical section
*/
func receiveMessageHeld(mutexMessage MessageMutexDTO) MessageMutexDTO {
	waitingForSendingAnswerBack(mutexMessage)
	return getReplyOkMessage()
}

/*
receiveMessageWanting - received a request message, your state: wanting
compare clocks and store request if yours is lower (if you are not the one waiting)
*/
func receiveMessageWanting(mutexMessage MessageMutexDTO) MessageMutexDTO {
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
	return getReplyOkMessage()
}

// --------------------
// HELPER METHODS

/*
waitingForSendingAnswerBack - add request to list of requests to answer back to
*/
func waitingForSendingAnswerBack(mutexMessage MessageMutexDTO) {
	requestChannel := make(chan string)
	requestChannelInfo := channelUserRequest{
		userEndpoint: mutexMessage.User,
		channel:      requestChannel,
	}
	mutexSendRequests = append(mutexSendRequests, requestChannelInfo)
	// wait until it is allowed to send a reply-ok
	logrus.Infof("[mutex_receive.waitingForSendingAnswerBack] wait until it is allowed to send a reply-ok to: %s", mutexMessage.User)
	msg := <- requestChannel
	logrus.Infof("[mutex_receive.waitingForSendingAnswerBack] received %s from channel", msg)
	// remove requestChannelInfo form mutexSendRequests
	mutexSendRequests = removeChannelUserRequest(requestChannelInfo, mutexSendRequests)
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
