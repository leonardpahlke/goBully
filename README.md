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

	- ReceiveMessage()             // get a message from a service (election, coordinator)
	- messageReceivedElection()    // handle incoming election message
	- sendElectionMessage()        // send a election message to another user
      ---------------------
	- messageReceivedCoordinator() // set local coordinator reference with incoming details
	- sendCoordinatorMessages()    // send coordinator messages to other users
	
Election Message received:
```
messageReceivedElection()
1. filter users to send election messages to (UserID > YourID)
2. if |filtered users| <= 0
   	YES: 2.1 you have the highest ID and win - send coordinatorMessages - exit
   	NO : 2.2 transform message and create POST payload
		 2.3 add callback information to local callbackList
         2.4 GO - sendElectionMessage()
            2.4.1 send POST request to client
            2.4.2 if response is OK check client callback
         2.5 wait a few seconds (enough time users can answer request)
         2.6 Sort users who have called back and who are not
         2.7 if |answered users| <= 0
			2.7.1 YES: send coordinatorMessages - exit
			2.7.2 NO : remove all users how didn't answered from userList
         2.8 clear callback list
3. send response back (answer)
```

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