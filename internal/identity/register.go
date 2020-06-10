package identity

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"goBully/pkg"
)

// static routes (api discovery would set these vars in a prod application)
const RegisterRoute = "/register"
const SendRegisterRoute = "/sendregister"

const UnRegisterRoute = "/unregister"
const SendUnRegisterRoute = "/sendunregister"

// REGISTER
/*
ReceiveServiceRegister - get user credentials from a new user and send them to the other connected users if new user send data directly to you
 */
func ReceiveServiceRegister(serviceRegisterInfo RegisterInfoDTO) RegisterResponseDTO {
	logrus.Infof("[identity.ReceiveServiceRegister] register information received from %s", serviceRegisterInfo.DistributingUserId)
	// check if sending id is also new id (-> do we have to notify other services?)
	distributingUserIsNewUser := serviceRegisterInfo.DistributingUserId == serviceRegisterInfo.NewUserId
	// create newUser information
	newUser := InformationUserDTO{ // info given (exec input)
		UserId:   serviceRegisterInfo.NewUserId,
		Endpoint: serviceRegisterInfo.Endpoint,
	}
	// add new id to users
	AddUser(newUser)
	// send other users the newUser information
	if distributingUserIsNewUser {
		// create payload to send others newUser info
		var myRegisterInfo = RegisterInfoDTO{
			DistributingUserId: YourUserInformation.UserId,
			NewUserId:          newUser.UserId,
			Endpoint:           newUser.Endpoint,
		}
		payload, err := json.Marshal(myRegisterInfo)
		if err != nil {
			logrus.Fatalf("[identity.ReceiveServiceRegister] Error marshal newUser with error %s", err)
		}
		var sendRegistrationTo = "["
		for _, user := range Users {
			// don't send register messages to yourself and the new user
			if (user.UserId != YourUserInformation.UserId) && (user.UserId != newUser.UserId) && (user.UserId != serviceRegisterInfo.DistributingUserId) {
				// IDEA we could wait if the other api answers and kick him out if he doesn't TODO do that?
				sendRegistrationTo = sendRegistrationTo + user.UserId + ", "
				res, err := pkg.RequestPOST(user.Endpoint +RegisterRoute, string(payload))
				if err != nil {
					logrus.Fatalf("[identity.ReceiveServiceRegister] Error sending post request with error %s", err)
				}
				var registerResponse RegisterResponseDTO
				err = json.Unmarshal(res, &registerResponse)
				if err != nil {
					logrus.Fatalf("[identity.ReceiveServiceRegister] Error Unmarshal post response with error %s", err)
				}
				logrus.Infof("[identity.ReceiveServiceRegister] register information send to service: %s", user.Endpoint)
			}
		}
		sendRegistrationTo = sendRegistrationTo + "]"
		return RegisterResponseDTO{
			Message:     "send register information to the other clients: " + sendRegistrationTo,
			UserIdInfos: Users,
		}
	}
	return RegisterResponseDTO{
		Message: YourUserInformation.UserId + " here, I have added new id " + newUser.UserId + " to my id pool",
		UserIdInfos: Users,
	}
}

/*
RegisterToService - send a registration message containing id details to an another endpoint
 */
func RegisterToService(userEndpoint string) string {
	endpointToRegisterTo := "http://" + userEndpoint
	// send YourUserInformation details as a payload to the api to get your identification
	payload, err := json.Marshal(RegisterInfoDTO{
		DistributingUserId: YourUserInformation.UserId,
		NewUserId:          YourUserInformation.UserId,
		Endpoint:           YourUserInformation.Endpoint,
	})
	if err != nil {
		logrus.Fatalf("[identity.RegisterToService] Error marshal newUser with error %s", err)
	}
	logrus.Info("[api.RegisterToService] prepare POST to register to endpoint: " + endpointToRegisterTo)
	res, err := pkg.RequestPOST(endpointToRegisterTo +RegisterRoute, string(payload))
	if err != nil {
		logrus.Fatalf("[identity.RegisterToService] Error sending POST request with error %s", err)
	}
	var registerResponse RegisterResponseDTO
	err = json.Unmarshal(res, &registerResponse)
	if err != nil {
		logrus.Fatalf("[identity.RegisterToService] Error Unmarshal registerResponse with error %s", err)
	}
	// set Users with all UserIdInfos (yours included)
	Users = registerResponse.UserIdInfos
	return "ok"
}

// UNREGISTER
/*
UnregisterUserFromYourUserList - unregister (without election algorithm)
 */
func UnregisterUserFromYourUserList(userInformation InformationUserDTO) bool {
	logrus.Info("[identity.UnregisterUserFromYourUserList] user: " + userInformation.UserId)
	return DeleteUser(userInformation)
}

/*
SendUnregisterUserFromYourUserList - unregister from all other user services
*/
func SendUnregisterUserFromYourUserList() bool {
	logrus.Info("[identity.SendUnregisterUserFromYourUserList] sending POST messages")
	payload, err := json.Marshal(YourUserInformation)
	if err != nil {
		logrus.Fatalf("[identity.SendUnregisterUserFromYourUserList] Error Unmarshal YourUserInformation with error %s", err)
	}
	for _, user := range Users {
		_, err = pkg.RequestPOST(user.Endpoint +UnRegisterRoute, string(payload))
		if err != nil {
			logrus.Fatalf("[identity.SendUnregisterUserFromYourUserList] Error RequestPOST with error %s", err)
		}
	}
	logrus.Info("[identity.SendUnregisterUserFromYourUserList] POST messages send")
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
	UserIdInfos []InformationUserDTO `json:"user_id_infos"`
}
