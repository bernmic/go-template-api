# go-template-api

This is a template for a simple HTTP REST api server written in Go. It uses Gin Gonic and delivers simple CRUD functionality for a users table.

config.go handles the cli parameters for configuration. http.go configures and starts the http listener and has generic functions for handling errors.
user.go contains the REST endpoints and their backend CRUD functions. 

## build

Checkout and build with `go get . && go build .` in the project folder.

## run

Just start the executable built in den "build" step. It accepts on parameter --port. Default is 8080. After it is started REST-API can be called at http://localhost:8080/api/user (or the port which is specified)
