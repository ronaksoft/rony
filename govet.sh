#!/usr/bin/env bash

# Generate codes
go generate ./... || exit

# Make sure the code guide lines are met
go vet ./... || exit

# Format the code
dirs=$(go list -f {{.Dir}} ./...)
for d in $dirs; do goimports -w $d/*.go; done

#go fmt ./... || exit

go install ./cmd/protoc-gen-gorony
go install ./cmd/rony

mkdir -p "$GOPATH"/src/github.com/ronaksoft/rony
cp ./*.proto "$GOPATH"/src/github.com/ronaksoft/rony


golangci-lint run