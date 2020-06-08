package identity

import "github.com/sirupsen/logrus"

// your identity information
var YourUserInformation InformationUserDTO
// store all active users (including YourUserInformation)
var Users []InformationUserDTO

/*
adds a identity to your identity pool
 */
func AddUser(userInformation InformationUserDTO) {
	Users = append(Users, userInformation)
	logrus.Info("[service.AddUser] identity added " + userInformation.UserId)
}

/*
deletes a identity from your identity pool
 */
func DeleteUser(userInformation InformationUserDTO) bool {
	for i, user := range Users {
		if user.UserId == userInformation.UserId {
			// delete identity from the list
			Users[i] = Users[len(Users)-1]
			Users = Users[:len(Users)-1]
			logrus.Info("[service.DeleteUser] identity deleted " + userInformation.UserId)
			return true
		}
	}
	logrus.Warning("[service.DeleteUser] identity could not be found and deleted " + userInformation.UserId)
	return false
}

// identity info struct
// swagger:model
type InformationUserDTO struct {
	// user identification which should be unique
	// required: true
	UserId   string `json:"userId"`
	// user endpoint to send http request
	// required: true
	Endpoint string `json:"endpoint"`
}

// get service user info
// swagger:model
type InformationUserInfoDTO struct {
	// all user linked to the service
	// required: true
	Users   []InformationUserDTO `json:"users"`
	// set coordinator
	Coordinator string `json:"coordinator"`
}