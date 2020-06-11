package mutex

import (
	"github.com/sirupsen/logrus"
	"goBully/internal/identity"
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
	logrus.Infof("[mutex_main.receiveMutexMessage] received message")
	// received a request message
	incrementClock(mutexMessage.Time)

	// response is set in receivedRequestMessage && receivedReplyMessage
	var mutexMessageResponse MessageMutexDTO

	switch mutexMessage.Msg {
	case RequestMessage: receivedRequestMessage(mutexMessage, &mutexMessageResponse)
	// TODO case not necessary? - fallback
	case ReplyOKMessage: receivedReplyMessage(mutexMessage, &mutexMessageResponse)
	default: logrus.Warningf("[mutex_main.receiveMutexMessage] message: %s, could not be identified", mutexMessage.Msg)
	}
	// send a reply-ok message - increase clock
	incrementClock(mutexMessage.Time)
	logrus.Infof("[mutex_main.receiveMutexMessage] send response message")
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
			logrus.Infof("[mutex_main.enterCriticalSection] 'released' state received, return")
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
	logrus.Infof("[mutex_main.leaveCriticalSection] leave, state: %s", state)
	mutexCriticalSection <- state
	// 2. notify critical section user is no longer in it
	for _, mutexSendRequest := range mutexSendRequests {
		mutexSendRequest.channel <- "leaveCriticalSection"
	}
}

// --------------------
// HELPER METHODS

/*
removeChannelUserRequest - TODO description
 */
func removeChannelUserRequest(channelReq channelUserRequest, channelReqs []channelUserRequest) []channelUserRequest {
	for i, req := range channelReqs {
		if req.userEndpoint == channelReq.userEndpoint {
			// delete identity from the list
			channelReqs[i] = channelReqs[len(channelReqs)-1]
			channelReqs = channelReqs[:len(channelReqs)-1]
			logrus.Infof("[mutex_main.removeChannelUserRequest] channel req deleted %s", channelReq.userEndpoint)
			return channelReqs
		}
	}
	logrus.Warningf("[mutex_main.removeChannelUserRequest] channel req could not be found %s", channelReq.userEndpoint)
	return channelReqs
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
	} else {
		return j
	}
}

// --------------------
// PRIVATE TYPES

// message to channel to identify users
type channelUserRequest struct {
	userEndpoint string
	user identity.InformationUserDTO
	channel chan string
	sendHealthChecks bool
}
