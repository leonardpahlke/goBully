package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"goBully/internal/election"
	id "goBully/internal/identity"
)

// REGISTER

// swagger:operation POST /register register registerService
// Register User information to api
// ---
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
// - in: body
//   name: register
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
	var serviceRegisterInfo id.RegisterInfoDTO
	err := c.BindJSON(&serviceRegisterInfo)
	if err != nil {
		logrus.Fatalf("[api.adapterRegisterService] Error marshal serviceRegisterInfo with error %s", err)
	}
	serviceRegisterResponse := id.ReceiveServiceRegister(serviceRegisterInfo)
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
// - in: query
//   type: number
//   name: sendregister
//   description: trigger registration, api sends registration message to other
//   required: true
// responses:
//  '200':
//    description: successful operation
//  '404':
//    description: error in operation
//  '403':
//    description: operation not available
func adapterSendRegisterToService(c *gin.Context) {
	// send post request to other endpoint to trigger connection cycle
	userEndpoint, _ := c.Params.Get("userEndpoint")
	logrus.Infof("[api.adapterSendRegisterToService] Received userEndpoint: %s", userEndpoint)
	msg := id.RegisterToService(userEndpoint)

	logrus.Infof("[api.adapterSendRegisterToService] register response received, message: %s - starting election, ...", msg)
	// start election to find a coordinator
	election.StartElectionAlgorithm(election.DummyElectionInfoDTO())

	// response check only if request was success full and has no further impact
	c.String(200, msg)
}

// UNREGISTER

// swagger:operation POST /unregister register unregisterFromService
// unregister api from your user list
// ---
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
// - in: body
//   name: unregister
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
	success := id.UnregisterUserFromYourUserList(informationUserDTO)
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
	success := id.SendUnregisterUserFromYourUserList()
	if success {
		c.String(200, DefaultSuccessMessage)
	} else {
		c.String(404, DefaultErrorMessage)
	}
	// c.String(403, DefaultNotAvailableMessage)
}
