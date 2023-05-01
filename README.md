### This repo is dedicated for "Pact and Go: Contract Testing of HTTP-based applications" article:
https://medium.com/@dees3g/build-simple-go-rest-api-in-seconds-7b67bf414064

## Pre-requisites:

Install pact https://github.com/pact-foundation/pact-go#installation

## Steps to perform Contract Testing:

#1 Generate/update a contract:
1. cd /client
2. go test

#2 Run provider test to validate it against the contract:
1. cd /server
2. go test
