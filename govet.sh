#!/usr/bin/env bash

# Generate codes
go generate ./... || exit

# Make sure the code guide lines are met
go vet ./... || exit

# Format the code
go fmt ./... || exit

go install ./cmd/protoc-gen-gorony
go install ./cmd/protoc-gen-goexport
go install ./cmd/rony

mkdir -p "$GOPATH"/src/github.com/ronaksoft/rony
cp ./*.proto "$GOPATH"/src/github.com/ronaksoft/rony


#golangci-lint run