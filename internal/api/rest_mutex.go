package api

import (
	"goBully/internal/mutex"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /mutex mutex mutexMessage
// handle mutex message
// ---
// consumes:
// - application/json
// produces:
// - application/json
// parameters:
// - in: body
//   name: mutex
//   description: mutex message information
//   required: true
//   schema:
//     "$ref": "#/definitions/MessageMutexDTO"
// responses:
//  '200':
//    description: successful operation
//    schema:
//      $ref: "#/definitions/MessageMutexDTO"
//  '404':
//    description: error in operation
//  '403':
//    description: operation not available
func adapterMutexMessage(c *gin.Context) {
	var mutexMessage mutex.MessageMutexDTO
	err := c.BindJSON(&mutexMessage)
	if err != nil {
		logrus.Fatalf("[api.adapterMutexMessage] Error marshal mutexMessage with error %s", err)
	}
	mutexResponse := mutex.ReceiveMutexMessage(mutexMessage)
	c.JSON(200, mutexResponse)
}

// swagger:operation GET /mutexstate mutex mutexStateRequest
// handle mutex a state request message
// ---
// consumes:
// - application/json
// produces:
// - application/json
// responses:
//  '200':
//    description: successful operation
//    schema:
//      $ref: "#/definitions/StateMutexDTO"
//  '404':
//    description: error in operation
//  '403':
//    description: operation not available
func adapterMutexStateMessage(c *gin.Context) {
	mutexState := mutex.RequestMutexState()
	c.JSON(200, mutexState)
}
