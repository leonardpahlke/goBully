package api

import (
	"goBully/internal/election"
	"goBully/internal/identity"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

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
	var electionInformation election.InformationElectionEntity
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
//   name: startelection
//   description: start election algorithm - to get a coordinator
//   required: true
//   schema:
//     "$ref": "#/definitions/InputInformationElectionDTO"
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
	var electionInformation election.InputInformationElectionEntity
	err := c.BindJSON(&electionInformation)
	if err != nil {
		logrus.Fatalf("[api.adapterStartElectionMessage] Error marshal electionInformation with error %s", err)
	}
	electionInfoResponse := election.StartElectionAlgorithm(election.TransformInputInfoElectionDTO(electionInformation))
	c.JSON(200, electionInfoResponse)
}

// swagger:operation POST /startstaticelection election startStaticElectionMessage
// execute election algorithm with preset input
// ---
// consumes:
// - application/json
// produces:
// - application/json
// responses:
//  '200':
//    description: successful operation
//    schema:
//      $ref: "#/definitions/InformationElectionDTO"
//  '404':
//    description: error in operation
//  '403':
//    description: operation not available
func adapterStartStaticElectionMessage(c *gin.Context) {
	var electionInformation = election.InformationElectionEntity{
		Algorithm: election.Algorithm,
		Payload:   election.MessageElection,
		User:      identity.YourUserInformation.UserId,
		Job:       election.InformationJobEntity{},
		Message:   "origin adapterStartStaticElectionMessage",
	}
	electionInformationResponse := election.StartElectionAlgorithm(electionInformation)
	c.JSON(200, electionInformationResponse)
}
