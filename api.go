package goBully

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

// TODO add api framework

const registerEndpoint = "register"

// ReceiveServiceRegister - a service sends a registration message to you - notify all other services in the system
func ReceiveServiceRegister() {
	logrus.Info("[api.ReceiveServiceRegister] register information received")
	newUser := UserInformation{ // info given (exec input)
		UserID:           "sample",
		CallbackEndpoint: "callback",
		Endpoint:         "endpoint",
	}
	payload, err := json.Marshal(newUser)
	if err != nil {
		logrus.Fatalf("[api.ReceiveServiceRegister] Error marshal newUser with error %s", err)
	}
	for _, user := range Users {
		if user.UserID != YourUserInformation.UserID {
			_, err := RequestPOST(user.Endpoint, string(payload), "")
			if err != nil {
				logrus.Fatalf("[api.ReceiveServiceRegister] Error sending post request with error %s", err)
			}
		}
	}
	logrus.Info("[api.ReceiveServiceRegister] register information send to services")
	// TODO return userList
	// TODO maybe send entire Users list (not only the new one -- in case of meh)
}

// RegisterToService - send a registration message containing user details to an another endpoint
func RegisterToService() {
	endpoint := "http://localhost:8080" // info given (exec input)
	// send YourUserInformation as a payload to the service to get your identification
	payload, err := json.Marshal(YourUserInformation)
	if err != nil {
		logrus.Fatalf("[api.RegisterToService] Error marshal newUser with error %s", err)
	}
	_, err = RequestPOST(endpoint + "/" + registerEndpoint, string(payload), "")
	if err != nil {
		logrus.Fatalf("[api.RegisterToService] Error sending post request with error %s", err)
	}
	logrus.Info("[api.RegisterToService] register information send")
}