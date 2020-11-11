#!/usr/bin/env bash

# Embed Assets into the binary
pkger -o ./cmd/rony

# Generate codes
go generate ./... || exit

# Make sure the code guide lines are met
go vet ./... || exit

# Format the code
go fmt ./... || exit


go install ./cmd/protoc-gen-gorony
go install ./cmd/rony