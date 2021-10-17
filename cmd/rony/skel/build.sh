#!/bin/bash

# Generate the auto-generated codes
rony gen-proto rpc

# Check for errors and warnings
go vet ./...

# Format your code
go fmt ./...