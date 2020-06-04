package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"gobully/internal"
)

func main() {
	logrus.Infof("[main] Starting Container \n get environment variables")
	// environment variables
	userID := os.Getenv("USERID")
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	endpoint := "http://" + host + ":" + port

	// set user information
	internal.YourUserInformation = internal.UserInformation{
		UserID:           userID,
		Endpoint:         endpoint,
	}

	// start api
	logrus.Infof("[main] Service Information set, starting api")
	internal.StartAPI(endpoint, port)
}
