# goBully

docker container as users

container with api for election and discovery

create multiple containers with docker-compose

**election.go**

	- receiveMessage()             // get a message from a service (election, answer, CoordinatorUserId)
	- sendElectionMessage()        // send a service an election message and wait for response
	- sendCoordinatorMessages()    // send a service that you are the CoordinatorUserId now
      ---------------------
	- messageReceivedAnswer()      // handle answer message
	- messageReceivedElection()    // handle election message
	- messageReceivedCoordinator() // handle coordinator message

![goBully](assets/goBully.jpg)

## Folder Structure

```
── goBully
├── api
│   └── ...                     // swagger files..
├── assets
│   └── ...
├── build
│   ├── Dockerfile
│   └── docker-compose.yml      // run szenario
├── cmd
│   └── main.go                 // starting point of the application
├── internal
│   ├── api
│   │   ├── request.go          // rest http calls
│   │   └── rest.go             // api setup - endpoints
│   ├── election
│   │   ├── election.go         // election private functions
│   │   └── election_client.go  // election public functions
│   └── service
│       ├── register.go         // user register workflow
│       └── user.go             // user definition
├── .gitignore
├── go.mod                      // project namespace
├── go.sum                      // project imports
└── README.md
```
