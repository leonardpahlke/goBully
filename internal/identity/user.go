package identity

import "github.com/sirupsen/logrus"

// your identity information
var YourUserInformation InformationUser
// store all active users (including YourUserInformation)
var Users []InformationUser

/*
adds a identity to your identity pool
 */
func AddUser(userInformation InformationUser) {
	Users = append(Users, userInformation)
	logrus.Info("[service.AddUser] identity added " + userInformation.UserId)
}

/*
deletes a identity from your identity pool
 */
func DeleteUser(userInformation InformationUser) {
	for i, user := range Users {
		if user.UserId == userInformation.UserId {
			// delete identity from the list
			Users[i] = Users[len(Users)-1]
			Users = Users[:len(Users)-1]
			break
		}
	}
	logrus.Info("[service.DeleteUser] identity deleted " + userInformation.UserId)
}

/* STRUCT */
// identity info struct
type InformationUser struct {
	UserId   string `json:"userId"`
	Endpoint string `json:"endpoint"`
}