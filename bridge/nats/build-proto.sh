#!/usr/bin/env bash

currentWorkingDir=`pwd`
rm -f ./*.pb.go
protoc  -I=${GOPATH}/src/git.ronaksoftware.com -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf -I=${GOPATH}/src -I=. --gogofaster_out=./ ./*.proto
go fmt




