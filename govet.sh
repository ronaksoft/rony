#!/usr/bin/env bash

go vet ./... || exit
go fmt ./... || exit