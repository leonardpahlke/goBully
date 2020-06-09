# for more information about go docker apk images: https://www.callicoder.com/docker-golang-image-container-example/
# build image to construct binary executable
FROM golang:latest as builder

LABEL maintainer="Leonard Pahlke <leonardpahlke@icloud.com>"

# execute further commands inside /app directory
WORKDIR /app

# copy go.mod file to download dependecies
COPY go.mod go.sum ./

# download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# copy everything in the root directory into the workingdirectory directory
COPY . .

# run go build to compile the binary executable
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# ---------------------------
# final image is smaller in size
FROM alpine:latest

# update application and add git
RUN apk --no-cache add ca-certificates

# set working directory
WORKDIR /root/

# copy build executables
COPY --from=builder /app/main .

# set default enviorment variables
ENV USERID exampleUser
ENV ENDPOINT 127.0.0.1:8080

# expose default port
EXPOSE 8080

# execute binary
CMD ["./main"]