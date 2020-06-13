package mutex

import (
	"encoding/json"
	"goBully/pkg"

	"github.com/sirupsen/logrus"
)

/*
receivedReplyOkMessage - receive a reply-ok message
TODO: detailed description
*/
func receivedReplyOkMessage(mutexMessage MessageMutexEntity) {
	logrus.Infof("[mutex_replyok.receivedReplyOkMessage] message received")

	// natify waiting channel
	for _, replyOkWaitingRoom := range replyOkwaitingList {
		// chack how many users still need to answer
		requestsNeeded := len(replyOkWaitingRoom.replyOkReceivingList)
		for _, userRequestChannel := range replyOkWaitingRoom.replyOkReceivingList {
			if userRequestChannel.user.Endpoint == mutexMessage.User {
				// notify wairing user serivce to stop sending heartbeats
				replyOkWaitingRoom.allReplyOkReceived <- ReplyOKMessage
				// last user answered - recevied all necessary reply-ok messages
				if requestsNeeded <= 1 {
					replyOkWaitingRoom.allReplyOkReceived <- ReplyOKMessage
				}
				return
			}
		}
	}

	logrus.Warnf("[mutex_replyok.receivedReplyOkMessage] reply-ok could not be matched")
}

/*
sendReplyOkMessage sending reply-ok message to user mutex endpoint
*/
func sendReplyOkMessage(endpoint string) {
	logrus.Infof("[mutex_replyok.SendReplyOkMessage] ")
	// std mutex reply-ok message
	mutexMessage := MessageMutexEntity{
		Msg:   ReplyOKMessage,
		Time:  clock,
		Reply: mutexYourReply,
		User:  mutexYourUser,
	}
	payload, err := json.Marshal(&mutexMessage)
	if err != nil {
		logrus.Fatalf("[mutex_replyok.sendReplyOkMessage] Error Marshal mutexMessage")
	}

	logrus.Infof("[mutex_replyok.SendReplyOkMessage] sending message")
	_, err = pkg.RequestPOST(endpoint, string(payload))
	if err != nil {
		logrus.Fatalf("[mutex_replyok.sendReplyOkMessage] Error sending RequestPOST to: %s", endpoint)
	}
}
