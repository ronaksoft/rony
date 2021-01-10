#!/usr/bin/env bash

protoc  -I=./testdata  -I=../.. --go_out=./testdata ./testdata/*.proto
protoc  -I=./testdata  -I=../.. --gorony_out=./testdata ./testdata/*.proto