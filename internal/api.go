package internal

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/swaggo/gin-swagger/example/basic/docs" // docs is generated by Swag CLI, you have to import it.
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func StartAPI(endpoint string, port string) {
	r := gin.New()

	url := ginSwagger.URL(endpoint + "/swagger/doc.json") // The url pointing to API definition
	// SWAGGER DOCUMENTATION
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// REGISTER SERVICE
	r.POST(RegisterRoute, func(c *gin.Context){
		var serviceRegisterInfo ServiceRegisterInfo
		err := c.BindJSON(&serviceRegisterInfo)
		if err != nil {
			logrus.Fatalf("[api.StartAPI] Error marshal serviceRegisterInfo with error %s", err)
		}
		serviceRegisterResponse := ReceiveServiceRegister(serviceRegisterInfo)
		// return all registered users to new user
		c.JSON(200, serviceRegisterResponse)
	})
	// TRIGGER REGISTER TO SERVICE
	r.POST("/send-register/:ip", func(c *gin.Context){
		// send post request to other endpoint to trigger connection cycle
		ip, _ := c.Params.Get("ip")
		msg := RegisterToService(ip)
		// response check only if request was success full and has no further impact
		c.String(200, msg)
	})

	// TODO TRIGGER UNREGISTER FROM SERVICE'S
	r.POST("/un-register", func(c *gin.Context){
		// unregister yourself from other user services (gracefully shutdown)
	})

	// TODO ELECTION
	r.POST(ElectionRoute, func(c *gin.Context){
		// start election algorithm - get a coordinator
	})

	err := r.Run(":" + port)
	if err != nil {
		logrus.Fatalf("[api.StartAPI] Error running server with error %s", err)
	}
}

// static routes (service discovery would set these vars in a prod application)
const RegisterRoute = "/register"

// types for POST data
type ServiceRegisterInfo struct {
	DistributingUserId string `json:"distributing_user_id"` // user sending new user information (new userId or some other userId)
	NewUserId string `json:"new_user_id"`                   // new user id
	Endpoint  string `json:"endpoint"`                      // new user endpoint
}
type ServiceRegisterResponse struct {
	Message string `json:"message"`
	UserIdInfos []UserInformation `json:"user_id_infos"`
}

func ReceiveServiceRegister(serviceRegisterInfo ServiceRegisterInfo) ServiceRegisterResponse {
	logrus.Info("[api.ReceiveServiceRegister] register information received")
	// check if sending user is also new user (-> do we have to notify other services?)
	distributingUserIsNewUser := serviceRegisterInfo.DistributingUserId == serviceRegisterInfo.NewUserId
	// create newUser information
	newUser := UserInformation{ // info given (exec input)
		UserID:           serviceRegisterInfo.NewUserId,
		Endpoint:         serviceRegisterInfo.Endpoint,
	}
	// add new user to users
	Users = append(Users, newUser)
	// send other users the newUser information
	if distributingUserIsNewUser {
		payload, err := json.Marshal(newUser)
		if err != nil {
			logrus.Fatalf("[api.ReceiveServiceRegister] Error marshal newUser with error %s", err)
		}
		for _, user := range Users {
			if user.UserID != YourUserInformation.UserID {
				// IDEA we could wait if the other service answers and kick him out if he doesn't TODO do that?
				res, err := RequestPOST(user.Endpoint +RegisterRoute, string(payload), "")
				if err != nil {
					logrus.Fatalf("[api.ReceiveServiceRegister] Error sending post request with error %s", err)
				}
				registerResponse := ServiceRegisterResponse{}
				err = json.Unmarshal(res, registerResponse)
				if err != nil {
					logrus.Fatalf("[api.ReceiveServiceRegister] Error Unmarshal post response with error %s", err)
				}
				logrus.Info("[api.ReceiveServiceRegister] received message " + registerResponse.Message)
			}
		}
		logrus.Info("[api.ReceiveServiceRegister] register information send to services")
		return ServiceRegisterResponse{
			Message:     "all registered users that I have noticed. I have send your information to the others..",
			UserIdInfos: Users,
		}
	}
	return ServiceRegisterResponse{
		Message:     YourUserInformation.UserID + " here, I have added new user " + newUser.UserID + " to my user pool",
		// IDEA we could sync users with each time post request via sending Users
		// (we do not do that because of performance concerns)
		// (we could only send our information which should do the job in the end as well - but meh)
		UserIdInfos: []UserInformation{},
	}
}

// RegisterToService - send a registration message containing user details to an another endpoint
func RegisterToService(ip string ) string {
	endpoint := "http://" + ip
	// send YourUserInformation as a payload to the service to get your identification
	payload, err := json.Marshal(ServiceRegisterInfo{
		DistributingUserId: YourUserInformation.UserID,
		NewUserId:          YourUserInformation.UserID,
		Endpoint:           YourUserInformation.Endpoint,
	})
	if err != nil {
		logrus.Fatalf("[api.RegisterToService] Error marshal newUser with error %s", err)
	}
	_, err = RequestPOST(endpoint +RegisterRoute, string(payload), "")
	if err != nil {
		logrus.Fatalf("[api.RegisterToService] Error sending post request with error %s", err)
	}
	logrus.Info("[api.RegisterToService] register information send")
	return "ok"
}