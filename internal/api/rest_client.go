package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"goBully/internal/election"
	"goBully/internal/identity"
	"goBully/internal/mutex"
	"time"
)

const DefaultSuccessMessage = "successful operation"
const DefaultErrorMessage = "error in operation"
// const DefaultNotAvailableMessage = "operation not available"

func StartAPI(port string) {
	// create api server - gin framework
	r := gin.New()

	// REST_USER
	// new identity register information
	r.GET(identity.RouteUserInfo, adapterUsersInfo)
	// new identity register information
	r.POST(identity.RegisterRoute, adapterRegisterService)
	// trigger identity register
	r.POST(identity.SendRegisterRoute+ "/:userEndpoint" , adapterSendRegisterToService)
	// trigger identity unregister from other identity services
	r.POST(identity.UnRegisterRoute, adapterUnRegisterFromService)
	// trigger identity unregister from other identity services
	r.POST(identity.SendUnRegisterRoute, adapterSendUnRegisterToServices)

	// REST_ELECTION
	// election algorithm endpoint
	r.POST(election.RouteElection, adapterElectionMessage)
	// start election algorithm endpoint
	r.POST(election.StartRouteElection, adapterStartElectionMessage)
	// start test election with static input
	r.POST(election.StartStaticRouteElection, adapterStartStaticElectionMessage)

	// REST_MUTEX
	// mutex message endpoint
	r.POST(mutex.RouteMutexMessage, adapterMutexMessage)
	// mutex state message endpoint
	r.GET(mutex.RouteMutexState, adapterMutexStateMessage)

	// start api server
	err := r.Run(":" + port)
	if err != nil {
		logrus.Fatalf("[api.StartAPI] Error running server with error %s", err)
	}
}

func ConnectToService(connectTo string) {
	time.Sleep(2 * time.Second)
	logrus.Infof("[api.ConnectToService] Connect to service %s", connectTo)
	// set user identities
	msg := identity.RegisterToService(connectTo)

	logrus.Infof("[api.ConnectToService] register response received, message: %s - starting election, ...", msg)
	// start election to find a coordinator
	election.StartElectionAlgorithm(election.DummyElectionInfoDTO())

	logrus.Infof("[api.ConnectToService] Connection to: %s complete, finished election, new coordinator: %s", connectTo, election.CoordinatorUserId)
	logrus.Print("----------------------")
}
