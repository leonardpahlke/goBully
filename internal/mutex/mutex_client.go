package mutex

import "github.com/sirupsen/logrus"

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
const ReplyMessage = "reply-ok"

// PRIVATE
// local lamport clock
var clock int32 = 0
// local mutex state
var state = StateHeld

/*
ReceiveMutexMessage - receive a mutex message from a user
respond with request or reply-ok
 */
func ReceiveMutexMessage(mutexMessage MessageMutexDTO) MessageMutexDTO {
	logrus.Infof("[mutex.ReceiveMutexMessage] received mutex message from user %s", mutexMessage.User)
	mutexMessageResponse := receiveMutexMessage(mutexMessage)
	return mutexMessageResponse
}

/*
RequestMutexState - return local mutex a state
*/
func RequestMutexState() StateMutexDTO {
	logrus.Infof("[mutex.RequestMutexState] received mutex state request")
	return StateMutexDTO {
		State: state,
		Time:  clock,
	}
}

/*
try to enter restricted area (your initiative)
 */
func ApplyEnterRestrictedArea() {
	// TODO
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
