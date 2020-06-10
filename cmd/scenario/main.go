package main

import (
	"github.com/sirupsen/logrus"
	"goBully/pkg"
)

/*
1. connect container 1 with 3
2. connect container 2 with 1, 3
 */
func main() {
	const container1Endpoint = "0.0.0.0:8080"
	//const container2Endpoint = "localhost:8081"
	const container3Endpoint = "0.0.0.0:8082"

	// 1. send triggerRegisterToService
	// container 1 sends post to container 3
	_, err := pkg.RequestPOST("http://" + container1Endpoint + "/sendregister/ " + container3Endpoint, "")
	if err != nil {
		logrus.Fatalf("[scenario.main] Error marshal newUser with error %s", err)
	}
}
