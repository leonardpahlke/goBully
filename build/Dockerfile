FROM golang:1.12.0-alpine3.9
# create /app directory within the dockerimage which will hold the application source files
RUN mkdir /app

# We copy everything in the root directory into our /app directory
ADD . /app

# execute further commands inside /app directory
WORKDIR /app

# run go build to compile the binary executable
RUN go build -o main .

# set enviorment variables (docker-compose overwrites these variables)
ENV USERID exampleUser
ENV ENDPOINT localhost:8080

# start the application
CMD ["/app/cmd/main"]