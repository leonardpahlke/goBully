# ---------------------------
# BUILD IMAGE
FROM golang:1.12-alpine as builder

# SETUP
# add folders
ADD . /go/src/goBully
# execute further commands inside /go/src/goBully directory
WORKDIR /go/src/goBully
# update application and add git
RUN apk --no-cache add ca-certificates git

# DEPENDENCIES
# copy go.mod file to download dependecies
COPY go.mod ./
# download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# check whether all depedencies have been downloaded
RUN go mod verify

# BUILD
# copy everything in the root directory into the workingdirectory directory
COPY . ./
#RUN CGO_ENABLED=0 go build
# run go build to compile the binary executable
RUN go build -o ./build ./cmd

# ---------------------------
# FINAL IMAGE
FROM alpine
# set working directory
WORKDIR /root
# copy build executables
COPY --from=builder /go/src/goBully/goBully .

# set default enviorment variables
ENV USERID exampleUser
ENV ENDPOINT localhost:8080

# expose default port
EXPOSE 8080

# execute
CMD ["./goBully"]
#CMD ["go", "run", "cmd/main.go"]