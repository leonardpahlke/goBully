package main

import (
	"gobully/internal/identity"
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
	// your identity service endpoint
	endpoint := "http://" + host + ":" + port

	// set identity information
	identity.YourUserInformation = identity.InformationUserDTO{
		UserId:   userID,
		Endpoint: endpoint,
	}
	// add yourself to identity list
	identity.AddUser(identity.YourUserInformation)

	// start api
	logrus.Infof("[main] Service Information set, starting api")
	service.StartAPI(port)
}
