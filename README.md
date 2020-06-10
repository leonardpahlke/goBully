# goBully

Project is under active development and not finished yet

This project implements the bully algorithm with docker containers. 
Several containers are served, each of which is accessible with a rest API. 
For more information, see the code comments and the Swagger documentation. 

## Install

1. **GO** [installation](https://golang.org/doc/install) getting started - *run project binary*  
2. **Docker** [installation](https://docs.docker.com/get-docker/) getting started - *be able to run docker containers* 
3. **Task** [installation](https://taskfile.dev/#/installation) doc - *build tool Taskfile.yml*
4. **Go Swagger** [installation](https://goswagger.io/install.html) doc - *swagger api documentation*

## Build
*execute commands within the project root directory*

**Check commands**
```
task --list
task: Available tasks for this project:
* build:        Build docker container
* run:          Start docker container
* sdown:        Stop docker-compose scenario
* sup:          Start docker-compose scenario
* swagger:      Generate swagger.yml and start local server
* update:       Update project dependencies
```
**Run commands**
```go
// run listed commands 
task <task>
// like
task build
```

**Stop Docker container**
```
docker stop $(docker ps -a -q --filter ancestor=leonardpahlke/gobully:latest --format="{{.ID}}")
```

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
├── Taskfile.yml                // build scripts - powered by Task
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
