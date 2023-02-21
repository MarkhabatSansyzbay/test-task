## Overview
This project consists of two services:

- Service 1 is a web-server with an endpoint (/generate-salt) for generating salt.
- Service 2 is a web-server that stores users' data. It has 2 endpoints (/create-user, /get-user/{email}). At endpoint "/create-user" it checks email for validity (regex) and uniqueness (should not be repeated in the database). Then it calls service 1 to get the salt, hashes that salt with password and saves the user's data. At endpoint "/get-user/{email}" it gets a user by its email from the database if user exists. Otherwise it returns 404 status.

## Usage

To run service1
```
go run cmd/main.go
```


To run service2
```
go run server/main.go 
```
and 
```
go run client/main.go
```
