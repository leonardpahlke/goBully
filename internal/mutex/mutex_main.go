package mutex

import (
	"github.com/sirupsen/logrus"
)

// store all users where I still neeed to respond with reply-ok to
var replyOkSendingList = []userSendingChannel{}

// list of
var replyOkwaitingList = []responseChannel{}

// channel to notify criticalSection to check state
var mutexCriticalSection = make(chan string)

/*
receiveMutexMessage - map logic after message
*/
func receiveMutexMessage(mutexMessage MessageMutexEntity) {
	logrus.Infof("[mutex_main.receiveMutexMessage] received message")

	switch mutexMessage.Msg {
	case RequestMessage:
		receivedRequestMessage(mutexMessage)
	case ReplyOKMessage:
		receivedReplyOkMessage(mutexMessage)
	default:
		logrus.Fatalf("[mutex_main.receiveMutexMessage] message: %s, is not could not a request message", mutexMessage.Msg)
	}
	// completed processing request message
	incrementClock(mutexMessage.Time)
}

/*
enterCriticalSection - enter critical section
1. update state to 'held'
2. wait to get notified to
*/
func enterCriticalSection() {
	state = StateHeld
	for true {
		stateChange := <-mutexCriticalSection
		if stateChange == StateReleased {
			logrus.Infof("[mutex_main.enterCriticalSection] 'released' state received, return")
			break
		} else {
			logrus.Warnf("[mutex_main.enterCriticalSection] state: %s, not 'released' yet", stateChange)
		}
	}
}

/*
leaveCriticalSection - enter critical section
1. update state to 'released'
2. notify critical section user is no longer in it
*/
func leaveCriticalSection() {
	// 1. update state to 'released'
	state = StateReleased
	logrus.Infof("[mutex_main.leaveCriticalSection] leave, state: %s", state)
	mutexCriticalSection <- state
	// 2. notify critical section user is no longer in it
	for _, replyokSendInvokeChannel := range replyOkSendingList {
		replyokSendInvokeChannel.channel <- ReplyOKMessage
	}
}

// --------------------
// HELPER METHODS

/*
removeChannelUserRequest - remove a channelUserRequest from a list (mutexSendRequests, mutexWaitingRequests, mutexReceivedRequests)
*/
func rmReplyOkSendingUser(userSendingChan userSendingChannel) {
	for i, userChan := range replyOkSendingList {
		if userChan.userEndpoint == userSendingChan.userEndpoint {
			// delete identity from the list
			replyOkSendingList[i] = replyOkSendingList[len(replyOkSendingList)-1]
			replyOkSendingList = replyOkSendingList[:len(replyOkSendingList)-1]
			logrus.Infof("[mutex_main.removeChannelUserRequest] user rm from reply-ok waiting list %s", userChan.userEndpoint)
			return
		}
	}
	logrus.Warnf("[mutex_main.removeChannelUserRequest]  user could not be found deleted")
}

/*
incrementClock - increase local lamport lock
*/
func incrementClock(i int32) int32 {
	clock = max(clock, i)
	logrus.Infof("[mutex_main.incrementClock] increasing local clock to %s", clock)
	return clock
}

/*
max - simple max function with int32 types
*/
func max(i int32, j int32) int32 {
	if i > j {
		return i
	}
	return j
}
