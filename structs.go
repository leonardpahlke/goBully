package goBully

// ELECTION
type ElectionInformation struct {
	Algorithm string                 `json:"algorithm"` // name of the algorithm used
	Payload   string                 `json:"payload"`   // the payload for the current state of the algorithm
	Callback  string                 `json:"callback"`  // uri of the user sending this request
	Job       ElectionJobInformation `json:"job"`
	Message   string                 `json:"message"` // something you want to tell the other one
}

type ElectionJobInformation struct {
	Id       string `json:"id"`       // some identity choosen by the initiator to identify this request
	Task     string `json:"task"`     // uri to the task to accomplish
	Resource string `json:"resource"` // uri or url to resource where actions are required
	Method   string `json:"method"`   // method to take – if already known
	Data     string `json:"data"`     // data to use/post for the task
	Callback string `json:"callback"` // an url where the initiator can be reached with the results/token
	Message  string `json:"message"`  // something you want to tell the other one
}

type ElectionCallbackInformation struct {
	Algorithm string                 `json:"algorithm"` // name of the algorithm used
	Payload   string                 `json:"payload"`   // the payload for the current state of the algorithm
	User      string                 `json:"user"`      // uri of the user sending this request
	Job       ElectionJobInformation `json:"job"`
	Message   string                 `json:"message"` // something you want to tell the other one
}

// USER (only for this scenario)
type UserInformation struct {
	UserID string `json:"userID"`
	CallbackEndpoint string `json:"callbackEndpoint"`
	Endpoint string `json:"endpoint"`
}