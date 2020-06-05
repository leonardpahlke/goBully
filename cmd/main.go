package main

import (
	"gobully/internal/api"
	"gobully/internal/service"
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Infof("[main] Starting Container \n get environment variables")
	// environment variables
	userID := os.Getenv("USERID")
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	// your user service endpoint
	endpoint := "http://" + host + ":" + port

	// set user information
	service.YourUserInformation = service.UserInformation{
		UserId:   userID,
		Endpoint: endpoint,
	}
	// add yourself to user list
	service.AddUser(service.YourUserInformation)

	// start api
	logrus.Infof("[main] Service Information set, starting api")
	api.StartAPI(endpoint, port)
}
