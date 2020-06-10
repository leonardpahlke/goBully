// Package api gBully API
//
// This project implements the bully algorithm with docker containers.
// Several containers are served, each of which is accessible with a rest API.
// For more information, see the code comments
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     Host: localhost
//     BasePath: /
//     Version: 0.2.0
//     License: Apache 2.0 http://www.apache.org/licenses/LICENSE-2.0.html
//
//     Consumes:
//     - application/json
//     - application/xml
//
//     Produces:
//     - application/json
//     - application/xml
//
// swagger:meta
package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"goBully/internal/election"
	"time"
)

const DefaultSuccessMessage = "successful operation"
const DefaultErrorMessage = "error in operation"
// const DefaultNotAvailableMessage = "operation not available"

func StartAPI(port string) {
	// create api server - gin framework
	r := gin.New()

	// REST_USER
	// new identity register information
	r.GET("/users", adapterUsersInfo)

	// REST_REGISTER
	// new identity register information
	r.POST(RegisterRoute, adapterRegisterService)
	// trigger identity register
	r.POST(SendRegisterRoute + "/:userEndpoint" , adapterSendRegisterToService)
	// trigger identity unregister from other identity services
	r.POST(UnRegisterRoute, adapterUnRegisterFromService)
	// trigger identity unregister from other identity services
	r.POST(SendUnRegisterRoute, adapterSendUnRegisterToServices)

	// REST_ELECTION
	// election algorithm endpoint
	r.POST(election.RouteElection, adapterElectionMessage)
	// start election algorithm endpoint
	r.POST(election.StartRouteElection, adapterStartElectionMessage)
	// start test election with static input
	r.POST(election.StartStaticRouteElection, adapterStartStaticElectionMessage)

	// start api server
	err := r.Run(":" + port)
	if err != nil {
		logrus.Fatalf("[api.StartAPI] Error running server with error %s", err)
	}
}

func ConnectToService(connectTo string) {
	time.Sleep(2 * time.Second)
	logrus.Infof("[api.ConnectToService] Connect to service %s", connectTo)
	msg := registerToService(connectTo)
	logrus.Infof("[api.ConnectToService] Connect to service with message: %s", msg)
}
