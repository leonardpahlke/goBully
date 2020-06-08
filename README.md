# goBully

Project is under active development and not finished yet

This project implements the bully algorithm with docker containers. 
Several containers are served, each of which is accessible with a rest API. 
For more information, see the code comments and also the Swagger documentation 

## Build

**Dockerfile**
```go
// Build Container
docker build -t leonardpahlke/gobully:latest .

// Run Container
docker run --rm -itd -p 8080:8080 leonardpahlke/gobully:latest

```

Connect to container http://localhost:8080

```
// Start Swagger Server
// TODO
```

**API [go-swagger](https://github.com/go-swagger/go-swagger) needed**
```bash
// validate swagger yml
swagger validate ./api/swagger.yml

// update swagger yml
swagger generate spec -o ./api/swagger.yml --scan-models

// start swagger server
swagger generate server -A goBully -f ./api/swagger.yml
```

## Start Scenario

**Docker compose**
```go
// Start
docker-compose up // TODO

// Stop
docker-compose down // TODO
```

// TODO scenario image 

## Features

may change
- docker container as user in the network to run the bully algorithm
- bully algorithm scenario with docker-compose simulated 
- detailed swagger documentation [Swagger yml](api/swagger.yml) with [go-swagger](https://github.com/go-swagger/go-swagger)

![goBully](assets/goBully.jpg)

## Project folder structure

```
── goBully
├── api
│   └── swagger.yml             // swagger api dcumentation
├── assets
│   └── ...                     // pictures and stuff
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
├── pkg
│   └── request.go              // rest http calls
├── .gitignore
├── docker-compose.yml          // dockercompose run szenario
├── Dockerfile                  // docker container script
├── go.mod                      // go module information
├── go.sum                      // go module libary imports
└── README.md
```

## Bully Algorithm implementation

`internal/election/election.go`

	- ReceiveMessage()             // get a message from a service (election, coordinator)
	- messageReceivedElection()    // handle incoming election message
	- sendElectionMessage()        // send a election message to another user
      ---------------------
	- messageReceivedCoordinator() // set local coordinator reference with incoming details
	- sendCoordinatorMessages()    // send coordinator messages to other users
	
more details

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