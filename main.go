package goBully

import (
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	logrus.Infof("Starting Container \n get environment variables")
	userID := os.Getenv("USERID")
	endpoint := os.Getenv("ENDPOINT")
	callbackEndpoint := "callback"
	YourUserInformation = UserInformation{
		UserID:           userID,
		CallbackEndpoint: callbackEndpoint,
		Endpoint:         endpoint,
	}
	logrus.Infof("Service Information set")
	// TODO start api
}