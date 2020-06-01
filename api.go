package goBully

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

// TODO add api framework

const registerEndpoint = "register"

// receiveServiceRegister - a service sends a registration message to you - notify all other services in the system
func receiveServiceRegister() {
	logrus.Info("[receiveServiceRegister] register information received")
	newUser := UserInformation{ // info given (exec input)
		UserID:           "sample",
		CallbackEndpoint: "callback",
		Endpoint:         "endpoint",
	}
	payload, err := json.Marshal(newUser)
	if err != nil {
		logrus.Fatalf("[receiveServiceRegister] Error marshal newUser with error %s", err)
	}
	for _, user := range Users {
		if user.UserID != YourUserInformation.UserID {
			_, err := RequestPOST(user.Endpoint, payload, "")
			if err != nil {
				logrus.Fatalf("[receiveServiceRegister] Error sending post request with error %s", err)
			}
		}
	}
	logrus.Info("[receiveServiceRegister] register information send to services")
}

// registerToService - send a registration message containing user details to an another endpoint
func registerToService() {
	endpoint := "http://localhost:8080" // info given (exec input)
	// send YourUserInformation as a payload to the service to get your identification
	payload, err := json.Marshal(YourUserInformation)
	if err != nil {
		logrus.Fatalf("[registerToService] Error marshal newUser with error %s", err)
	}
	_, err := RequestPOST(endpoint + "/" + registerEndpoint, payload, "")
	if err != nil {
		logrus.Fatalf("[registerToService] Error sending post request with error %s", err)
	}
	logrus.Info("[registerToService] register information send")
}