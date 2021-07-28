---
title: "Getting Started"
date: 2021-07-26T14:33:20+04:30
draft: false
weight: 50
menu: "main"
tags: ["wiki", "install", "guide"]
categories: ["wiki"]
weight: 20
contentCopyright: MIT
mathjax: true
autoCollapseToc: true

---

# Prepare
To install Rony we need to have protoc compiler installed in our system. With a quick search we can
find good resources. Some useful references provided here:
* https://grpc.io/docs/protoc-installation/ 
* http://google.github.io/proto-lens/installing-protoc.html


We also need to install protoc golang plugin from [here](https://github.com/golang/protobuf)

# Install Rony
To install Rony just run: `GO111MODULE=on go get -u github.com/ronaksoft/rony/...` or in Go 1.16 and above use: `go install github.com/ronaksoft/rony/...`

Rony has its own protoc plugin which is required to generate codes based on your `proto` files. Make sure `/cmd/protoc-gen-gorony` is installed
on your path. **Usually `go install` takes care of this.**


# First Project
To crate your first project, we need to use `rony` executable which is installed in previous section. We could 
find out about the flags by calling `rony -h`. To create our new project skeleton run :

```shell
mkdir sample-project
cd ./sample-project
rony create-project --project.name github.com/ronaksoft/sample --project.name sample-project
```

There are some examples in `./example` directory. 
