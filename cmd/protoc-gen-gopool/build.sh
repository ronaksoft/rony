#!/usr/bin/env bash

protoc  -I=./testdata  --gopool_out=./testdata ./testdata/*.proto