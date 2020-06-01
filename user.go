package goBully

// store all active users
var YourUserInformation UserInformation
var Users []UserInformation

// adds a user to the user pool
func AddUser(userInformation UserInformation) {
	Users = append(Users, userInformation)
}

// deletes a user from the user pool
func DeleteUser(userInformation UserInformation) {
	for i, user := range Users {
		if user.UserID == userInformation.UserID {
			// delete user from the list
			Users[i] = Users[len(Users)-1]
			Users = Users[:len(Users)-1]
			break
		}
	}
}
