### Contract Testing HTTP-based service example


## Pre-requisites:

Install pact https://github.com/pact-foundation/pact-go#installation

## Steps to perform Contract Testing:

#1 Generate/update a contract:
1. cd contract/pkg/client
2. go test

#2 Run a local server (provider):
1. cd contract/cmd/server
2. go run main.go

#3 Run a server/provider test:
1. cd contract/pkg/server
2. go test
