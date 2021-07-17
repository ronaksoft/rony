#!/bin/bash

# Generate the auto-generated codes
go generate ./...

# Check for errors and warnings
go vet ./...

# Format your code
go fmt ./...