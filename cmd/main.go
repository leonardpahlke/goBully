package main

import (
	"gobully/internal/identity"
	"gobully/internal/service"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Infof("[main] Starting Container \n get environment variables")
	// environment variables
	userID := os.Getenv("USERID")
	endpoint := os.Getenv("ENDPOINT")
	port := strings.SplitAfter(endpoint, ":")[1]

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
