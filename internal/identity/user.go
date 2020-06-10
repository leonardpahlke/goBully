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
	if !ContainsUser(Users, userInformation) {
		Users = append(Users, userInformation)
		logrus.Infof("[api.AddUser] user %s added", userInformation.UserId)
	} else {
		logrus.Infof("[api.AddUser] user %s not added - already in user list ", userInformation.UserId)
	}
}

/*
return whether a user is in user list
*/
func ContainsUser(userList []InformationUserDTO, user InformationUserDTO) bool {
	for _, a := range userList {
		if a == user {
			return true
		}
	}
	return false
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
			logrus.Infof("[api.DeleteUser] identity deleted %s", userInformation.UserId)
			return true
		}
	}
	logrus.Warningf("[api.DeleteUser] identity could not be found and deleted %s", userInformation.UserId)
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

// get api user info
// swagger:model
type InformationUserInfoDTO struct {
	// all user linked to the api
	// required: true
	Users   []InformationUserDTO `json:"users"`
	// set coordinator
	Coordinator string `json:"coordinator"`
}
