package service

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gobully/internal/api"
	"gobully/internal/election"
)

// static routes (service discovery would set these vars in a prod application)
const RegisterRoute = "/register"

/*
TODO ReceiveServiceRegister - description
 */
func ReceiveServiceRegister(serviceRegisterInfo RegisterInfo) RegisterResponse {
	logrus.Info("[service.ReceiveServiceRegister] register information received")
	// check if sending user is also new user (-> do we have to notify other services?)
	distributingUserIsNewUser := serviceRegisterInfo.DistributingUserId == serviceRegisterInfo.NewUserId
	// create newUser information
	newUser := UserInformation{ // info given (exec input)
		UserId:   serviceRegisterInfo.NewUserId,
		Endpoint: serviceRegisterInfo.Endpoint,
	}
	// add new user to users
	Users = append(Users, newUser)
	// send other users the newUser information
	if distributingUserIsNewUser {
		payload, err := json.Marshal(newUser)
		if err != nil {
			logrus.Fatalf("[service.ReceiveServiceRegister] Error marshal newUser with error %s", err)
		}
		for _, user := range Users {
			if user.UserId != YourUserInformation.UserId {
				// IDEA we could wait if the other service answers and kick him out if he doesn't TODO do that?
				res, err := api.RequestPOST(user.Endpoint +RegisterRoute, string(payload), "")
				if err != nil {
					logrus.Fatalf("[service.ReceiveServiceRegister] Error sending post request with error %s", err)
				}
				registerResponse := RegisterResponse{}
				err = json.Unmarshal(res, registerResponse)
				if err != nil {
					logrus.Fatalf("[service.ReceiveServiceRegister] Error Unmarshal post response with error %s", err)
				}
				logrus.Info("[service.ReceiveServiceRegister] received message " + registerResponse.Message)
			}
		}
		logrus.Info("[service.ReceiveServiceRegister] register information send to services")
		return RegisterResponse{
			Message:     "all registered users that I have noticed. I have send your information to the others..",
			UserIdInfos: Users,
		}
	}
	return RegisterResponse{
		Message: YourUserInformation.UserId + " here, I have added new user " + newUser.UserId + " to my user pool",
		// IDEA we could sync users with each time post request via sending Users
		// (we do not do that because of performance concerns)
		// (we could only send our information which should do the job in the end as well - but meh)
		UserIdInfos: []UserInformation{},
	}
}

/*
RegisterToService - send a registration message containing user details to an another endpoint
 */
func RegisterToService(ip string ) string {
	endpoint := "http://" + ip
	// send YourUserInformation as a payload to the service to get your identification
	payload, err := json.Marshal(RegisterInfo{
		DistributingUserId: YourUserInformation.UserId,
		NewUserId:          YourUserInformation.UserId,
		Endpoint:           YourUserInformation.Endpoint,
	})
	if err != nil {
		logrus.Fatalf("[service.RegisterToService] Error marshal newUser with error %s", err)
	}
	res, err := api.RequestPOST(endpoint + RegisterRoute, string(payload), "")
	if err != nil {
		logrus.Fatalf("[service.RegisterToService] Error sending post request with error %s", err)
	}
	var registerResponse RegisterResponse
	err = json.Unmarshal(res, &registerResponse)
	if err != nil {
		logrus.Fatalf("[service.RegisterToService] Error Unmarshal registerResponse with error %s", err)
	}
	// set Users with all UserIdInfos (yours included)
	Users = registerResponse.UserIdInfos
	logrus.Info("[service.RegisterToService] register information send and user info set, starting election, ...")
	election.StartElectionAlgorithm()
	logrus.Info("[service.RegisterToService] finished election coordinator: " + election.CoordinatorUserId)
	return "ok"
}

// TODO seed register (without election algorithm)

/* STRUCT */
// object sending user service to register yourself
type RegisterInfo struct {
	DistributingUserId string `json:"distributing_user_id"` // user sending new user information (new userId or some other userId)
	NewUserId string `json:"new_user_id"`                   // new user id
	Endpoint  string `json:"endpoint"`                      // new user endpoint
}
// response object after register to user service
type RegisterResponse struct {
	Message string                `json:"message"` 			// dummy message to print response
	UserIdInfos []UserInformation `json:"user_id_infos"` 	// all registered users
}