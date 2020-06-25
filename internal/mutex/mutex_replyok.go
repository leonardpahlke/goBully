package mutex

import (
	"encoding/json"
	"goBully/internal/identity"
	"goBully/pkg"

	"github.com/sirupsen/logrus"
)

/*
receivedReplyOkMessage - receive a reply-ok message
1. get first waiting channel that fits the mutexMessage.Endpoint
2. check how many users still need to answer
3. notify waiting task to stop sending heartbeats
4. if last user answered - send message through channel allReplyOkReceived to notify waiting task
*/
func receivedReplyOkMessage(mutexMessage MessageMutexEntity) {
	logrus.Infof("[mutex_replyok.receivedReplyOkMessage] message received")

	// 1. get first waiting channel that fits the mutexMessage.Endpoint
	for _, replyOkWaitingRoom := range replyOkwaitingList {
		// 2. check how many users still need to answer
		requestsNeeded := len(replyOkWaitingRoom.replyOkReceivingList)

		for _, userRequestChannel := range replyOkWaitingRoom.replyOkReceivingList {

			if userRequestChannel.user.Endpoint == mutexMessage.User {
				// 3. notify waiting task to stop sending heartbeats
				userRequestChannel.channel <- ReplyOKMessage

				// 4. if last user answered - send message through channel allReplyOkReceived to notify waiting task
				if requestsNeeded <= 1 {
					replyOkWaitingRoom.allReplyOkReceived <- ReplyOKMessage
				}
				return
			}
		}
	}
	logrus.Warnf("[mutex_replyok.receivedReplyOkMessage] reply-ok could not be matched to a waiting task")
}

/*
sendReplyOkMessage sending reply-ok message to user mutex endpoint
1. create a reply-ok mutexMessage
2. send reply-ok message to user mutex endpoint
*/
func sendReplyOkMessage(endpoint string) {
	logrus.Infof("[mutex_replyok.SendReplyOkMessage] ")

	// 1. create a reply-ok mutexMessage
	mutexMessage := getMutexMessage(ReplyOKMessage)
	payload, err := json.Marshal(&mutexMessage)
	if err != nil {
		logrus.Fatalf("[mutex_replyok.sendReplyOkMessage] Error Marshal mutexMessage")
	}

	// 2. send reply-ok message to user mutex endpoint
	logrus.Infof("[mutex_replyok.SendReplyOkMessage] sending message")
	_, err = pkg.RequestPOST(endpoint, string(payload))
	if err != nil {
		logrus.Fatalf("[mutex_replyok.sendReplyOkMessage] Error sending RequestPOST to: %s", endpoint)
	}
}

// --------------------
// HELPER METHODS

/*
getWaitingTaskInformation remove first task found
1. loop over all replyOkwaitingList entries
2. loop over each replyOkwaitingList.replyOkReceivingList entries
3. if user waiting entry found
4. create a new list without the user waiting entry
5. set new information in replyOkwaitingList
*/
func rmWaitingTaskInformation(user identity.InformationUserDTO) {
	// 1. loop over all replyOkwaitingList entries
	for i, replyOkWaitingRoom := range replyOkwaitingList {
		// 2. loop over each replyOkwaitingList.replyOkReceivingList entries

		for j, userRequestChannel := range replyOkWaitingRoom.replyOkReceivingList {
			// 3. if user waiting entry found

			if userRequestChannel.user.Endpoint == user.Endpoint {
				// 4. create a new list without the user waiting entry
				replyOkWaitingRoom.replyOkReceivingList = rmWaitingEntry(j, replyOkWaitingRoom.replyOkReceivingList)

				// 5. set new information in replyOkwaitingList
				replyOkwaitingList[i] = replyOkWaitingRoom
			}
		}
	}
}

/*
rmWaitingEntry delete a waiting user task from userReponseChannelList
*/
func rmWaitingEntry(i int, userRepChannels []userReponseChannel) []userReponseChannel {
	userRepChannels[i] = userRepChannels[len(userRepChannels)-1]
	userRepChannels = userRepChannels[:len(userRepChannels)-1]
	logrus.Infof("[mutex_replyok.rmWaitingEntry] user rm from reply-ok waiting list")
	return userRepChannels
}
