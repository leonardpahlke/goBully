package goBully

import (
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	logrus.Infof("[main] Starting Container \n get environment variables")
	userID := os.Getenv("USERID")
	endpoint := os.Getenv("ENDPOINT")
	callbackEndpoint := endpoint + "election/callback" // TODO check this later
	YourUserInformation = UserInformation{
		UserID:           userID,
		CallbackEndpoint: callbackEndpoint,
		Endpoint:         endpoint,
	}
	logrus.Infof("[main] Service Information set")
	// TODO start api
}