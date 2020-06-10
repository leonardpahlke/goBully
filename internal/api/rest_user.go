package api

import (
	"github.com/gin-gonic/gin"
	"goBully/internal/election"
	id "goBully/internal/identity"
)

// swagger:operation GET /users user users
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
