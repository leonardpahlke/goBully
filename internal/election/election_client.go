package election

import "github.com/sirupsen/logrus"

// Public function to interact with election

// API Endpoints
const RouteElection = "/election"
const StartRouteElection = "/startelection"

// TODO - print info logs - maybe also in a config file defined
// var Verbose = true

// current CoordinatorUserId
var CoordinatorUserId = ""

/*
start election algorithm (your initiative)
 */
func StartElectionAlgorithm(informationElectionDTO InformationElectionDTO) InformationElectionDTO {
	logrus.Infof("[election.StartElectionAlgorithm] starting..")
	response := ReceiveMessage(informationElectionDTO)
	return response
}

// election state information
// swagger:model
type InformationElectionDTO struct {
	// name of the algorithm used
	// required: true
	Algorithm string            `json:"algorithm"`
	// the payload for the current state of the algorithm
	// required: true
	Payload   string            `json:"payload"`
	// uri of the identity sending this request
	// required: true
	User      string            `json:"identity"`
	// jon information in InformationJobDTO
	// required: true
	Job       InformationJobDTO `json:"job"`
	// something you want to tell the other one
	// required: true
	Message   string            `json:"message"`
}

// election job details
// swagger:model
type InformationJobDTO struct {
	// some identity chosen by the initiator to identify this request
	// required: true
	Id       string `json:"id"`
	// uri to the task to accomplish
	// required: true
	Task     string `json:"task"`
	// uri or url to resource where actions are required
	// required: true
	Resource string `json:"resource"`
	// method to take – if already known
	// required: true
	Method   string `json:"method"`
	// data to use/post for the task
	// required: true
	Data     string `json:"data"`
	// an url where the initiator can be reached with the results/token
	// required: true
	Callback string `json:"callback"`
	// something you want to tell the other one
	// required: true
	Message  string `json:"message"`
}