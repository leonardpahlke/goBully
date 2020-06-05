# goBully

Project is under active development and not finished yet

This project implements the bully algorithm with docker containers. 
Several containers are served, each of which is accessible with a rest API. 
For more information, see the code comments and also the Swagger documentation 

## Build

TODO

## Start Scenario

TODO

## Features

may change

- docker container as users
- container with api for election and discovery
- create multiple containers with docker-compose

**Bully Algorithm Structure** - `internal/election/election.go`

	- receiveMessage()             // get a message from a service (election, answer, CoordinatorUserId)
	- sendElectionMessage()        // send a service an election message and wait for response
	- sendCoordinatorMessages()    // send a service that you are the CoordinatorUserId now
      ---------------------
	- messageReceivedAnswer()      // handle answer messages
	- messageReceivedElection()    // handle election messages
	- messageReceivedCoordinator() // handle coordinator messages

![goBully](assets/goBully.jpg)

## Project folder structure

```
── goBully
├── api
│   └── swagger.yml             // swagger api dcumentation
├── assets
│   └── ...                     // pictures and stuff
├── build
│   ├── Dockerfile              // docker container script
│   └── docker-compose.yml      // dockercompose run szenario
├── cmd
│   └── main.go                 // starting point of the application
├── internal
│   ├── election
│   │   ├── election.go         // election private functions
│   │   └── election_client.go  // election public functions
│   ├── identity
│   │   └── user.go             // user definition
│   └── service
│       ├── register.go         // user register workflow
│       └── rest.go             // api setup - endpoints
├── .gitignore
│   └── request.go              // rest http calls
├── .gitignore
├── go.mod                      // go module information
├── go.sum                      // go module libary imports
└── README.md
```

**swagger api**

[Swagger Doc](api/swagger.yml)

Update swagger yml ([go-swagger](https://github.com/go-swagger/go-swagger) needed)

`swagger generate spec -o ./api/swagger.yml --scan-models`