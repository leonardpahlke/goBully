FROM golang:1.12.0-alpine3.9

# set enviorment variables (docker-compose overwrites these variables)
ENV USERID exampleUser
ENV ENDPOINT localhost:8080

# create /app directory within the dockerimage which will hold the application source files
RUN mkdir /app

# execute further commands inside /app directory
WORKDIR /app

# Copy go mod and sum files (dependencies)
COPY go.mod go.sum ./

RUN apk add --update --no-cache ca-certificates git

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download


# We copy everything in the root directory into the /app directory
COPY . .

RUN ls -la

# run go build to compile the binary executable
RUN go build -o main . # TODO

# start the application
CMD ["/app/cmd/main"]