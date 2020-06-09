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
	id "goBully/internal/identity"
)

const DefaultSuccessMessage = "successful operation"
const DefaultErrorMessage = "error in operation"
// const DefaultNotAvailableMessage = "operation not available"

func StartAPI(port string) {
	// create api server - gin framework
	r := gin.New()

	// USER
	// new identity register information
	r.GET("/users", adapterUsersInfo)

	// REGISTER
	// new identity register information
	r.POST(RegisterRoute, adapterRegisterService)
	// trigger identity register
	r.POST(SendRegisterRoute + "/:ip" , adapterTriggerRegisterToService)
	// trigger identity unregister from other identity services
	r.POST(UnRegisterRoute, adapterUnRegisterFromService)
	// trigger identity unregister from other identity services
	r.POST(SendUnRegisterRoute, adapterSendUnRegisterToServices)

	// ELECTION
	// election algorithm endpoint
	r.POST(election.RouteElection, adapterElectionMessage)
	// start election algorithm endpoint
	r.POST(election.StartRouteElection, adapterStartElectionMessage)

	// start api server
	err := r.Run(":" + port)
	if err != nil {
		logrus.Fatalf("[api.StartAPI] Error running server with error %s", err)
	}
}

// swagger:operation POST /register user registerService
// Get registered user information's and coordinator
// ---
// consumes:
// - application/json
// produces:
// - application/json
// responses:
//  '200':
//    description: successful operation
//    schema:
//      $ref: "#/definitions/InformationUserInfoDTO"
//  '404':
//    description: error in operation
//  '403':
//    description: operation not available
func adapterUsersInfo(c *gin.Context) {
	var informationUserInfoDTO = id.InformationUserInfoDTO{
		Users:       id.Users,
		Coordinator: election.CoordinatorUserId,
	}
	// return all registered users and coordinator information
	c.JSON(200, informationUserInfoDTO)
}


// swagger:operation POST /register register registerService
// Register User information to api
// ---
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
// - in: body
//   name: api
//   description: send register information to get in the network
//   required: true
//   schema:
//     "$ref": "#/definitions/RegisterInfoDTO"
// responses:
//  '200':
//    description: successful operation
//    schema:
//      $ref: "#/definitions/RegisterResponseDTO"
//  '404':
//    description: error in operation
//  '403':
//    description: operation not available
func adapterRegisterService(c *gin.Context) {
	var serviceRegisterInfo RegisterInfoDTO
	err := c.BindJSON(&serviceRegisterInfo)
	if err != nil {
		logrus.Fatalf("[api.adapterRegisterService] Error marshal serviceRegisterInfo with error %s", err)
	}
	serviceRegisterResponse := receiveServiceRegister(serviceRegisterInfo)
	// return all registered users to new identity
	c.JSON(200, serviceRegisterResponse)
}


// swagger:operation POST /sendregister register triggerRegisterToService
// User sends register request to another user and kick off election to get the new coordinator
// ---
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
// - in: body
//   name: electionInformation
//   description: start election algorithm - to get a coordinator
//   required: true
//   schema:
//     "$ref": "#/definitions/InformationElectionDTO"
// - in: query
//   type: number
//   name: ip
//   description: trigger registration, api sends registration message to other
//   required: true
// responses:
//  '200':
//    description: successful operation
//  '404':
//    description: error in operation
//  '403':
//    description: operation not available
func adapterTriggerRegisterToService(c *gin.Context) {
	// send post request to other endpoint to trigger connection cycle
	var informationElectionDTO election.InformationElectionDTO
	err := c.BindJSON(&informationElectionDTO)
	if err != nil {
		logrus.Fatalf("[api.adapterTriggerRegisterToService] Error marshal informationElectionDTO with error %s", err)
	}
	ip, _ := c.Params.Get("ip")
	msg := registerToService(ip, informationElectionDTO)
	// response check only if request was success full and has no further impact
	c.String(200, msg)
}

// swagger:operation POST /unregister register unregisterFromService
// unregister api from your user list
// ---
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
// - in: body
//   name: api
//   description: some api is unregistering from all users, remove user from active users
//   required: true
//   schema:
//     "$ref": "#/definitions/InformationUserDTO"
// responses:
//  '200':
//    description: successful operation
//  '404':
//    description: error in operation
//  '403':
//    description: operation not available
func adapterUnRegisterFromService(c *gin.Context) {
	// send post request to other endpoint to trigger connection cycle
	var informationUserDTO id.InformationUserDTO
	err := c.BindJSON(&informationUserDTO)
	if err != nil {
		logrus.Fatalf("[api.adapterUnRegisterFromService] Error marshal informationUserDTO with error %s", err)
	}
	success := unregisterUserFromYourUserList(informationUserDTO)
	if success {
		c.String(200, DefaultSuccessMessage)
	} else {
		c.String(404, DefaultErrorMessage)
	}
}

// swagger:operation POST /sendunregister register sendUnregisterToServices
// unregister yourself from other user api user lists
// ---
// consumes:
// - application/json
// produces:
// - application/json
// responses:
//  '200':
//    description: successful operation
//  '404':
//    description: error in operation
//  '403':
//    description: operation not available
func adapterSendUnRegisterToServices(c *gin.Context) {
	// trigger method to send all unregister messages to users
	success := sendUnregisterUserFromYourUserList()
	if success {
		c.String(200, DefaultSuccessMessage)
	} else {
		c.String(404, DefaultErrorMessage)
	}
	// c.String(403, DefaultNotAvailableMessage)
}

// swagger:operation POST /election election electionMessage
// handle election algorithm state
// ---
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
// - in: body
//   name: election
//   description: election algorithm - get a coordinator
//   required: true
//   schema:
//     "$ref": "#/definitions/InformationElectionDTO"
// responses:
//  '200':
//    description: successful operation
//    schema:
//      $ref: "#/definitions/InformationElectionDTO"
//  '404':
//    description: error in operation
//  '403':
//    description: operation not available
func adapterElectionMessage(c *gin.Context) {
	var electionInformation election.InformationElectionDTO
	err := c.BindJSON(&electionInformation)
	if err != nil {
		logrus.Fatalf("[api.adapterElectionMessage] Error marshal electionInformation with error %s", err)
	}
	electionInformationResponse := election.ReceiveMessage(electionInformation)
	c.JSON(200, electionInformationResponse)
}

// swagger:operation POST /startelection election startElectionMessage
// execute election algorithm
// ---
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
// - in: body
//   name: election
//   description: start election algorithm - to get a coordinator
//   required: true
//   schema:
//     "$ref": "#/definitions/InformationElectionDTO"
// responses:
//  '200':
//    description: successful operation
//    schema:
//      $ref: "#/definitions/InformationElectionDTO"
//  '404':
//    description: error in operation
//  '403':
//    description: operation not available
func adapterStartElectionMessage(c *gin.Context) {
	var electionInformation election.InformationElectionDTO
	err := c.BindJSON(&electionInformation)
	if err != nil {
		logrus.Fatalf("[api.adapterStartElectionMessage] Error marshal electionInformation with error %s", err)
	}
	electionInformationResponse := election.StartElectionAlgorithm(electionInformation)
	c.JSON(200, electionInformationResponse)
}