package api

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"goBully/internal/election"
	id "goBully/internal/identity"
	"goBully/pkg"
)

// static routes (api discovery would set these vars in a prod application)
const RegisterRoute = "/register"
const SendRegisterRoute = "/sendregister"
const UnRegisterRoute = "/unregister"
const SendUnRegisterRoute = "/sendunregister"

/*
receiveServiceRegister - get user credentials from a new user and send them to the other connected users if new user send data directly to you
 */
func receiveServiceRegister(serviceRegisterInfo RegisterInfoDTO) RegisterResponseDTO {
	logrus.Info("[api.receiveServiceRegister] register information received")
	// check if sending id is also new id (-> do we have to notify other services?)
	distributingUserIsNewUser := serviceRegisterInfo.DistributingUserId == serviceRegisterInfo.NewUserId
	// create newUser information
	newUser := id.InformationUserDTO{ // info given (exec input)
		UserId:   serviceRegisterInfo.NewUserId,
		Endpoint: serviceRegisterInfo.Endpoint,
	}
	// add new id to users
	id.AddUser(newUser)
	// send other users the newUser information
	if distributingUserIsNewUser {
		payload, err := json.Marshal(newUser)
		if err != nil {
			logrus.Fatalf("[api.receiveServiceRegister] Error marshal newUser with error %s", err)
		}
		for _, user := range id.Users {
			if user.UserId != id.YourUserInformation.UserId {
				// IDEA we could wait if the other api answers and kick him out if he doesn't TODO do that?
				res, err := pkg.RequestPOST(user.Endpoint + RegisterRoute, string(payload))
				if err != nil {
					logrus.Fatalf("[api.receiveServiceRegister] Error sending post request with error %s", err)
				}
				registerResponse := RegisterResponseDTO{}
				err = json.Unmarshal(res, &registerResponse)
				if err != nil {
					logrus.Fatalf("[api.receiveServiceRegister] Error Unmarshal post response with error %s", err)
				}
				logrus.Info("[api.receiveServiceRegister] received message " + registerResponse.Message)
			}
		}
		logrus.Info("[api.receiveServiceRegister] register information send to services")
		return RegisterResponseDTO{
			Message:     "all registered users that I have noticed. I have send your information to the others..",
			UserIdInfos: id.Users,
		}
	}
	return RegisterResponseDTO{
		Message: id.YourUserInformation.UserId + " here, I have added new id " + newUser.UserId + " to my id pool",
		// IDEA we could sync users with each time post request via sending Users
		// (we do not do that because of performance concerns)
		// (we could only send our information which should do the job in the end as well - but meh)
		UserIdInfos: []id.InformationUserDTO{},
	}
}

/*
registerToService - send a registration message containing id details to an another endpoint
 */
func registerToService(ip string, informationElectionDTO election.InformationElectionDTO) string {
	endpoint := "http://" + ip
	// send YourUserInformation as a payload to the api to get your identification
	payload, err := json.Marshal(RegisterInfoDTO{
		DistributingUserId: id.YourUserInformation.UserId,
		NewUserId:          id.YourUserInformation.UserId,
		Endpoint:           id.YourUserInformation.Endpoint,
	})
	if err != nil {
		logrus.Fatalf("[api.registerToService] Error marshal newUser with error %s", err)
	}
	res, err := pkg.RequestPOST(endpoint + RegisterRoute, string(payload))
	if err != nil {
		logrus.Fatalf("[api.registerToService] Error sending post request with error %s", err)
	}
	var registerResponse RegisterResponseDTO
	err = json.Unmarshal(res, &registerResponse)
	if err != nil {
		logrus.Fatalf("[api.registerToService] Error Unmarshal registerResponse with error %s", err)
	}
	// set Users with all UserIdInfos (yours included)
	id.Users = registerResponse.UserIdInfos
	logrus.Info("[api.registerToService] register information send and id info set, starting election, ...")
	election.StartElectionAlgorithm(informationElectionDTO)
	logrus.Info("[api.registerToService] finished election coordinator: " + election.CoordinatorUserId)
	return "ok"
}


/*
unregisterUserFromYourUserList - unregister (without election algorithm)
 */
func unregisterUserFromYourUserList(userInformation id.InformationUserDTO) bool {
	logrus.Info("[api.unregisterUserFromYourUserList] user: " + userInformation.UserId)
	return id.DeleteUser(userInformation)
}

/*
sendUnregisterUserFromYourUserList - unregister from all other user services
*/
func sendUnregisterUserFromYourUserList() bool {
	logrus.Info("[api.sendUnregisterUserFromYourUserList] sending POST messages")
	payload, err := json.Marshal(id.YourUserInformation)
	if err != nil {
		logrus.Fatalf("[api.sendUnregisterUserFromYourUserList] Error Unmarshal YourUserInformation with error %s", err)
	}
	for _, user := range id.Users {
		_, err = pkg.RequestPOST(user.Endpoint + UnRegisterRoute, string(payload))
		if err != nil {
			logrus.Fatalf("[api.sendUnregisterUserFromYourUserList] Error RequestPOST with error %s", err)
		}
	}
	logrus.Info("[api.sendUnregisterUserFromYourUserList] POST messages send")
	return true
}

// object sending id api to register yourself
// swagger:model
type RegisterInfoDTO struct {
	// id sending new id information (new userId or some other userId)
	// required: true
	DistributingUserId string `json:"distributing_user_id"`
	// new userId id, check if Distributing user is also new one to notify others if so
	// required: true
	NewUserId string `json:"new_user_id"`
	// new userId endpoint
	// required: true
	Endpoint  string `json:"endpoint"`
}

// response object after register to id api
// swagger:model
type RegisterResponseDTO struct {
	// dummy message to print response
	// required: true
	Message string                   `json:"message"`
	// all registered users
	// required: true
	UserIdInfos []id.InformationUserDTO `json:"user_id_infos"`
}