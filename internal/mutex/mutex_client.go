package mutex

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"goBully/internal/identity"
	"goBully/pkg"
	"time"
)

// TODO enhancement - config file

// PUBLIC
// API Endpoints
const RouteMutexMessage = "/mutex"
const RouteMutexState = "/mutexstate"

// states a client can be in regarding entering the mutex
const StateReleased = "released"
const StateWanting = "wanting"
const StateHeld = "held"

// messages send across clients
const RequestMessage = "request"
const ReplyOKMessage = "reply-ok"

// static mutex val's
var mutexYourReply = identity.YourUserInformation.Endpoint
var mutexYourUser = RouteMutexMessage

// waiting time to send health checks
const waitingTime = time.Second * 3

// PRIVATE
// local lamport clock
var clock int32 = 0
// local mutex state
var state = StateReleased

/*
ReceiveMutexMessage - receive a mutex message from a user
respond with request or reply-ok
 */
func ReceiveMutexMessage(mutexMessage MessageMutexDTO) MessageMutexDTO {
	logrus.Infof("[mutex_client.ReceiveMutexMessage] received mutex message from user %s", mutexMessage.User)
	mutexMessageResponse := receiveMutexMessage(mutexMessage)
	return mutexMessageResponse
}

/*
RequestMutexState - return local mutex a state
*/
func RequestMutexState() StateMutexDTO {
	logrus.Infof("[mutex_client.RequestMutexState] received mutex state request")
	return StateMutexDTO {
		State: state,
		Time:  clock,
	}
}

/*
ApplyEnterRestrictedArea - try to enter restricted area (your initiative)
TODO
 */
func ApplyEnterRestrictedArea() {
	requestCriticalArea()
}

/*
- TDOO
 */
func LeaveCriticalSection() {
	if state == StateHeld {
		leaveCriticalSection()
	} else {
		logrus.Infof("[mutex_client.LeaveCriticalSection] requesting to leave critical section but you are currently in state: %s", state)
	}
}

/*
RequestUserState - request user state information
 */
func RequestUserState(userEndpoint string, userMutexStateEndpoint string) StateMutexDTO {
	res, err := pkg.RequestGET(userEndpoint + userMutexStateEndpoint)
	if err != nil {
		logrus.Fatalf("[mutex_client.RequestUserState] Error request with error %s", err)
	}
	var stateMutexDTO StateMutexDTO
	err = json.Unmarshal(res, &stateMutexDTO)
	if err != nil {
		logrus.Fatalf("[mutex_client.RequestUserState] Error Unmarshal stateMutexDTO with error %s", err)
	}
	return stateMutexDTO
}

// mutex message
// swagger:model
type MessageMutexDTO struct {
	// message, reply-ok or request
	// required: true
	Msg string `json:"msg"`
	// the current lamport clock
	// required: true
	Time int32 `json:"time"`
	// url to the endpoint where responses shall be send
	// required: true
	Reply string `json:"reply"`
	// url to the user sending the message
	// required: true
	User string `json:"user"`
}

// mutex state
// swagger:model
type StateMutexDTO struct {
	// current state: released, wanting or held
	// required: true
	State string `json:"state"`
	// the current lamport clock
	// required: true
	Time int32 `json:"time"`
}
