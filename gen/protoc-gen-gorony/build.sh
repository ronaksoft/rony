#!/usr/bin/env bash

protoc  -I=./testdata  --gorony_out=./testdata ./testdata/*.proto