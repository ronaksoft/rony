#!/usr/bin/env bash

go generate ./... || exit
go vet ./... || exit
go fmt ./... || exit