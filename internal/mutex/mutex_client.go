package mutex

import (
	"encoding/json"
	"goBully/internal/identity"
	"goBully/pkg"
	"time"

	"github.com/sirupsen/logrus"
)

/*
PUBLIC vals
*/

// API Endpoints

// RouteMutexMessage api endpoint to send mutex messages
const RouteMutexMessage = "/mutex"

// RouteMutexState api endpoint to request current user mutex state
const RouteMutexState = "/mutexstate"

// Mutex States

// StateReleased user is in idle state
const StateReleased = "released"

// StateWanting user wants to enter critical section and is waiting for reply-ok messages
const StateWanting = "wanting"

// StateHeld user is currently in the critical section
const StateHeld = "held"

// Mutex Messages

// RequestMessage this message is send if the user would like to enter the critical section
const RequestMessage = "request"

// ReplyOKMessage this message is send to a 'request' message if it's ok for the user the reqesting user enters the critical section
const ReplyOKMessage = "reply-ok"

/*
PRIVATE vals
*/

// Config

// waitingTime how long to wait until a user sends a response back
const waitingTime = time.Second * 10

// mutexYourReply mutex state static reply response
var mutexYourReply = identity.YourUserInformation.Endpoint

// mutexYourUser mutex state static user response
var mutexYourUser = RouteMutexMessage

// clock mutex internal local lamport clock
var clock int32 = 0

// state current user mutex state
var state = StateReleased

/*
METHODS
*/

/*
ReceiveMutexMessage - receive a mutex message from a user
respond with request or reply-ok
*/
func ReceiveMutexMessage(mutexMessage MessageMutexEntity) {
	logrus.Infof("[mutex_client.ReceiveMutexMessage] received mutex message from user %s", mutexMessage.User)
	receiveMutexMessage(mutexMessage)
}

/*
RequestMutexState - return local mutex a state
*/
func RequestMutexState() StateMutexEntity {
	logrus.Infof("[mutex_client.RequestMutexState] received mutex state request")
	return StateMutexEntity{
		State: state,
		Time:  clock,
	}
}

/*
RequestCriticalArea - try to enter restricted area (your initiative)
*/
func RequestCriticalArea() {
	requestCriticalArea()
	// you are now in the critical section until you invoke LeaveCriticalSection()
}

/*
LeaveCriticalSection - execute this method if you are finished with critical area work
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
func RequestUserState(userEndpoint string, userMutexStateEndpoint string) StateMutexEntity {
	res, err := pkg.RequestGET(userEndpoint + userMutexStateEndpoint)
	if err != nil {
		logrus.Fatalf("[mutex_client.RequestUserState] Error request with error %s", err)
	}
	var stateMutexDTO StateMutexEntity
	err = json.Unmarshal(res, &stateMutexDTO)
	if err != nil {
		logrus.Fatalf("[mutex_client.RequestUserState] Error Unmarshal stateMutexDTO with error %s", err)
	}
	return stateMutexDTO
}

/*
Public STRUCTS
*/

// MessageMutexEntity mutex message
// swagger:model
type MessageMutexEntity struct {
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

// StateMutexEntity mutex state
// swagger:model
type StateMutexEntity struct {
	// current state: released, wanting or held
	// required: true
	State string `json:"state"`
	// the current lamport clock
	// required: true
	Time int32 `json:"time"`
}
