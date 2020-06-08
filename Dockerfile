FROM golang:1.12.0-alpine3.9

ADD . /go/src/goBully

# execute further commands inside /app directory
WORKDIR /go/src/goBully

# Copy go mod and sum files (dependencies)
RUN apk add --update --no-cache ca-certificates git

# copy mod files and sum files to download dependecies
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Check whether all depedencies have been downloaded
RUN go mod verify


# copy everything in the root directory into the workingdirectory directory
COPY . .

RUN go build -o ./build ./cmd
# set enviorment variables (docker-compose overwrites these variables)
ENV USERID exampleUser
ENV ENDPOINT localhost:8080

# run go build to compile the binary executable
EXPOSE 8080

# start the application
CMD ["go", "run", "cmd/main.go"]