package service

import "github.com/sirupsen/logrus"

// your user information
var YourUserInformation UserInformation
// store all active users (including YourUserInformation)
var Users []UserInformation

/*
adds a user to your user pool
 */
func AddUser(userInformation UserInformation) {
	Users = append(Users, userInformation)
	logrus.Info("[service.AddUser] user added " + userInformation.UserId)
}

/*
deletes a user from your user pool
 */
func DeleteUser(userInformation UserInformation) {
	for i, user := range Users {
		if user.UserId == userInformation.UserId {
			// delete user from the list
			Users[i] = Users[len(Users)-1]
			Users = Users[:len(Users)-1]
			break
		}
	}
	logrus.Info("[service.DeleteUser] user deleted " + userInformation.UserId)
}

/* STRUCT */
// user info struct
type UserInformation struct {
	UserId   string `json:"userId"`
	Endpoint string `json:"endpoint"`
}