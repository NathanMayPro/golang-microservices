# API test

The goal of this file is to test the creation of a new API.

net_http.go will test the creation with native package net/http.

gin.go will test the creation with Gin framework.

## Explaination
go.mod is the file that contains the dependencies of the project.
go.sum is the file that contains the checksum of the dependencies.
The checksum is used to verify that the content of the dependencies is the same as the one that was downloaded.


## DOCKERFILE
The dockerfile is used to create a docker image of the project.
The docker image is used to run the project in a container.
The docker image is created with the command:
```bash
docker build -t api-test .
```
The docker image is run with the command:
```bash
docker run -p 8080:8080 api-test
```
The docker image is stopped with the command:
```bash
docker stop <container_id>
```
The docker image is removed with the command:
```bash
docker rmi api-test
```


The content of the dockerfile is:
```dockerfile
FROM golang:1.16.3-alpine3.13

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main .

CMD ["/app/main"]
```