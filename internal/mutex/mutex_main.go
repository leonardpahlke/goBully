package mutex

import (
	"github.com/sirupsen/logrus"
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

var mutexCriticalSection = make(chan string)

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

/*
enterCriticalSection - enter critical section
1. update state to 'held'
 */
func enterCriticalSection() {
	state = StateHeld
	for true {
		stateChange := <- mutexCriticalSection
		if stateChange == StateReleased {
			logrus.Infof("[mutex.enterCriticalSection] 'released' state received, return")
			break
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
	logrus.Infof("[mutex.leaveCriticalSection] leave, state: %s", state)
	mutexCriticalSection <- state
	// 2. notify critical section user is no longer in it
	for _, mutexSendRequest := range mutexSendRequests {
		mutexSendRequest.channel <- "leaveCriticalSection"
	}
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
removeChannelUserRequest - TODO description
 */
func removeChannelUserRequest(channelReq channelUserRequest, channelReqs []channelUserRequest) []channelUserRequest {
	for i, req := range channelReqs {
		if req.userEndpoint == channelReq.userEndpoint {
			// delete identity from the list
			channelReqs[i] = channelReqs[len(channelReqs)-1]
			channelReqs = channelReqs[:len(channelReqs)-1]
			logrus.Infof("[mutex.removeChannelUserRequest] channel req deleted %s", channelReq.userEndpoint)
			return channelReqs
		}
	}
	logrus.Warningf("[mutex.removeChannelUserRequest] channel req could not be found %s", channelReq.userEndpoint)
	return channelReqs
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
	sendHealthChecks bool
}
