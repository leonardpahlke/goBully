package service

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"gobully/internal/election"
	id "gobully/internal/identity"
	"gobully/pkg"
)

// static routes (service discovery would set these vars in a prod application)
const RegisterRoute = "/register"

/*
TODO ReceiveServiceRegister - description
 */
func ReceiveServiceRegister(serviceRegisterInfo RegisterInfo) RegisterResponse {
	logrus.Info("[service.ReceiveServiceRegister] register information received")
	// check if sending id is also new id (-> do we have to notify other services?)
	distributingUserIsNewUser := serviceRegisterInfo.DistributingUserId == serviceRegisterInfo.NewUserId
	// create newUser information
	newUser := id.InformationUser{ // info given (exec input)
		UserId:   serviceRegisterInfo.NewUserId,
		Endpoint: serviceRegisterInfo.Endpoint,
	}
	// add new id to users
	id.Users = append(id.Users, newUser)
	// send other users the newUser information
	if distributingUserIsNewUser {
		payload, err := json.Marshal(newUser)
		if err != nil {
			logrus.Fatalf("[service.ReceiveServiceRegister] Error marshal newUser with error %s", err)
		}
		for _, user := range id.Users {
			if user.UserId != id.YourUserInformation.UserId {
				// IDEA we could wait if the other service answers and kick him out if he doesn't TODO do that?
				res, err := pkg.RequestPOST(user.Endpoint +RegisterRoute, string(payload), "")
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
			UserIdInfos: id.Users,
		}
	}
	return RegisterResponse{
		Message: id.YourUserInformation.UserId + " here, I have added new id " + newUser.UserId + " to my id pool",
		// IDEA we could sync users with each time post request via sending Users
		// (we do not do that because of performance concerns)
		// (we could only send our information which should do the job in the end as well - but meh)
		UserIdInfos: []id.InformationUser{},
	}
}

/*
RegisterToService - send a registration message containing id details to an another endpoint
 */
func RegisterToService(ip string ) string {
	endpoint := "http://" + ip
	// send YourUserInformation as a payload to the service to get your identification
	payload, err := json.Marshal(RegisterInfo{
		DistributingUserId: id.YourUserInformation.UserId,
		NewUserId:          id.YourUserInformation.UserId,
		Endpoint:           id.YourUserInformation.Endpoint,
	})
	if err != nil {
		logrus.Fatalf("[service.RegisterToService] Error marshal newUser with error %s", err)
	}
	res, err := pkg.RequestPOST(endpoint + RegisterRoute, string(payload), "")
	if err != nil {
		logrus.Fatalf("[service.RegisterToService] Error sending post request with error %s", err)
	}
	var registerResponse RegisterResponse
	err = json.Unmarshal(res, &registerResponse)
	if err != nil {
		logrus.Fatalf("[service.RegisterToService] Error Unmarshal registerResponse with error %s", err)
	}
	// set Users with all UserIdInfos (yours included)
	id.Users = registerResponse.UserIdInfos
	logrus.Info("[service.RegisterToService] register information send and id info set, starting election, ...")
	election.StartElectionAlgorithm()
	logrus.Info("[service.RegisterToService] finished election coordinator: " + election.CoordinatorUserId)
	return "ok"
}

// TODO seed register (without election algorithm)

/* STRUCT */
// object sending id service to register yourself
type RegisterInfo struct {
	DistributingUserId string `json:"distributing_user_id"` // id sending new id information (new userId or some other userId)
	NewUserId string `json:"new_user_id"`                   // new id id
	Endpoint  string `json:"endpoint"`                      // new id endpoint
}
// response object after register to id service
type RegisterResponse struct {
	Message string                   `json:"message"`       // dummy message to print response
	UserIdInfos []id.InformationUser `json:"user_id_infos"` // all registered users
}