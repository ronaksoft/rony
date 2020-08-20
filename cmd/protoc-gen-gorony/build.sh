#!/usr/bin/env bash

protoc  -I=./testdata  --go_out=./testdata ./testdata/*.proto
protoc  -I=./testdata  --gorony_out=./testdata ./testdata/*.proto