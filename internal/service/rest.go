package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gobully/internal/election"
)

func StartAPI(endpoint string, port string) {
	// create api server - gin framework
	r := gin.New()

	// API ENDPOINTS
	// new identity register information
	r.POST(RegisterRoute, HandleRegisterService)
	// trigger identity register
	r.POST("/sendregister/:ip", handleTriggerRegisterService)
	// trigger identity unregister from other identity services
	r.POST("/un-register", handleTriggerUnRegisterFromServices)
	// election algorithm endpoint
	r.POST(election.RouteElection, handleElection)

	// start api server
	err := r.Run(":" + port)
	if err != nil {
		logrus.Fatalf("[api.StartAPI] Error running server with error %s", err)
	}
}

/*
POST handle REGISTER SERVICE
 */
func HandleRegisterService(c *gin.Context) {
	var serviceRegisterInfo RegisterInfo
	err := c.BindJSON(&serviceRegisterInfo)
	if err != nil {
		logrus.Fatalf("[api.StartAPI] Error marshal serviceRegisterInfo with error %s", err)
	}
	serviceRegisterResponse := ReceiveServiceRegister(serviceRegisterInfo)
	// return all registered users to new identity
	c.JSON(200, serviceRegisterResponse)
}

/*
POST handle TRIGGER REGISTER TO SERVICE
*/
// swagger:operation POST /materials materials addMaterial
// Add a new material
// ---
// consumes:
// - application/json
// parameters:
// - name: material
//   in: body
//   description: New material
//   required: true
//   schema:
//     "$ref": "#/definitions/MaterialCreateDTO"
// responses:
//  '201':
//    description: Successfully created material
func handleTriggerRegisterService(c *gin.Context) {
	// send post request to other endpoint to trigger connection cycle
	ip, _ := c.Params.Get("ip")
	msg := RegisterToService(ip)
	// response check only if request was success full and has no further impact
	c.String(200, msg)
}

/*
POST handle TRIGGER UNREGISTER FROM SERVICE'S
unregister yourself from other identity services (gracefully shutdown)
*/
func handleTriggerUnRegisterFromServices(c *gin.Context) {
	// TODO unregister identity service
	c.String(403, "this service is not available at the moment")
}

/*
POST handle election algorithm state
election algorithm - get a coordinator
 */
func handleElection(c *gin.Context) {
	// TODO election
	c.String(403, "this service is not available at the moment")
}