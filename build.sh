#!/usr/bin/env bash

go vet ./...
go fmt ./...
go generate ./...