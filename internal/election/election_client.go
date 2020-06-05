package election

// Public function to interact with election

// API Endpoints
const RouteElection = "/election"
// TODO - print info logs - maybe also in a config file defined
// var Verbose = true
// current CoordinatorUserId
var CoordinatorUserId = ""

/*
start election algorithm (your initiative)
 */
func StartElectionAlgorithm() {
	// TODO
}

/*
pass election message into logic and handle display
 */
func HandleElectionMessage() {
	// TODO
}

/* STRUCT */
// election state information
type InformationElection struct {
	Algorithm string         `json:"algorithm"` // name of the algorithm used
	Payload   string         `json:"payload"`   // the payload for the current state of the algorithm
	User      string         `json:"identity"`      // uri of the identity sending this request
	Job       InformationJob `json:"job"`
	Message   string         `json:"message"`   // something you want to tell the other one
}
// election job details
type InformationJob struct {
	Id       string `json:"id"`       // some identity choosen by the initiator to identify this request
	Task     string `json:"task"`     // uri to the task to accomplish
	Resource string `json:"resource"` // uri or url to resource where actions are required
	Method   string `json:"method"`   // method to take – if already known
	Data     string `json:"data"`     // data to use/post for the task
	Callback string `json:"callback"` // an url where the initiator can be reached with the results/token
	Message  string `json:"message"`  // something you want to tell the other one
}