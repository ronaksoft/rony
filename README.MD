# Rony (Fast and Scalable RPC Framework)

![](https://img.shields.io/github/stars/ronaksoft/rony?color=808080)
![](https://img.shields.io/github/license/ronaksoft/rony?color=808080)

![GitHub release (latest SemVer including pre-releases)](https://img.shields.io/github/v/release/ronaksoft/rony?include_prereleases&logoColor=808080)
![](https://github.com/ronaksoft/rony/actions/workflows/go.yml/badge.svg?color=379c9c&style=flat-square)
![](https://github.com/ronaksoft/rony/actions/workflows/devskim-analysis.yml/badge.svg?color=379c9c&style=flat-square)
![](https://github.com/ronaksoft/rony/actions/workflows/codacy-analysis.yml/badge.svg?color=379c9c&style=flat-square)

## About
**Rony** lets you create a clustered aware service easily. Basically when your write a service with
**Rony** framework, you develop two interfaces one is for the client side which clients of your service connect and communicate. We call it 'Gateway'. The other interface
communicates with other instances in the cluster which we call it 'Tunnel'.
Basically when you are building your service with **Rony** you spin up an 'Edge' server. Each edge server has three components which you can enable them if you need to. '
Gateway', 'Tunnel' and 'Cluster'. You cannot enable 'Tunnel' without enabling 'Cluster'.
Each **Edge** server could be member of a replica set if cluster enabled. 
Each **Edge** server could have its 'Tunnel' component enabled. This is how the edge servers could communicate with each other. For example, imagine you have a request from a
client connecting to an edge with replica set to '1' but to return the response we need to execute a command on replica set '2'. This could be easily done by using 
the '
Tunnel' features in **Rony**.

When we want to write a service to support sharding and replication, developers need to write a lot of code to manage their communication and how to write the code to handle these
interactions between services and clients. **Rony** makes it easy.   
However, to develop your service using **Rony** framework, there are more works to initialize your code comparing to other RESTful frameworks. To help you set up the skeleton even
faster there is an executable to makes the life easier for you.

### Overview Diagram
![](docs/rony_overview.png?raw=true)
### Installation

Use go get to install the Go client-and-server generator:

```
GO111MODULE=on go get -u github.com/ronaksoft/rony/...
```

You will also need:

* protoc, the protobuf compiler. You need version 3+.
* github.com/golang/protobuf/protoc-gen-go, the Go protobuf generator plugin. Get this with `go get`

### Getting Started

After we have installed two executable files `protoc-gen-rony` and `rony` you are ready to create your project.

```shell script
1. mkdir sample-project
2. cd ./sample-project
3. rony create-project --project.name github.com/ronaksoft/sample --project.name sample-project
```

After running the 'create-project' command, you should go to the `service` directory and open `service.proto`
file and write the appropriate protobuf file. We edit the protobuf then we run the following command 
from the root of our project.

```shell script
rony gen-proto
```

The above command generates the boiler plate codes for our service.

---

### Performance

Rony is very fast and with very low overhead. In a non-raft environment it adds < 25us latency, and in raft-enabled scenario around 1ms latency.

Rony has negligible memory footprint by reusing buffers and pooling.

```commandline
BenchmarkEdgeServerMessageSerial
BenchmarkEdgeServerMessageSerial-16                       901370              1195 ns/op              87 B/op          2 allocs/op
BenchmarkEdgeServerMessageParallel
BenchmarkEdgeServerMessageParallel-16                    4521645               272 ns/op              83 B/op          2 allocs/op
BenchmarkEdgeServerWithRaftMessageSerial
BenchmarkEdgeServerWithRaftMessageSerial-16                 9541            132065 ns/op            5034 B/op        116 allocs/op
BenchmarkEdgeServerWithRaftMessageParallel
BenchmarkEdgeServerWithRaftMessageParallel-16             124658              8438 ns/op            4462 B/op         51 allocs/op

BenchmarkServerWithWebsocket-16            46514             25138 ns/op             691 B/op         19 allocs/op
```

---

### Easy Setup for advanced scenarios

```go
package main

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"os"
)

func main() {
	serverID := "UniqueServerID"
	edgeServer := edge.NewServer(serverID,
		edge.WithTcpGateway(edge.TcpGatewayConfig{
			Concurrency:   1000,
			MaxIdleTime:   0,
			ListenAddress: "0.0.0.0:80",
			ExternalAddrs: []string{"127.0.0.1:80"}, // ExternalAddr could be used when the server is behind proxy or nats
		}),
		edge.WithUdpTunnel(edge.UdpTunnelConfig{
			ListenAddress: "0.0.0.0:8080",
			MaxBodySize:   0,
			ExternalAddrs: nil,
		}),
		edge.WithGossipCluster(edge.GossipClusterConfig{
			Bootstrap:  true,
			ReplicaSet: 1,
			GossipPort: 7091,
		}),
	)

	service.RegisterSampleService(&service.SampleService{}, edgeServer)
	edgeServer.Start()

	edgeServer.ShutdownWithSignal(os.Kill)
}
```

**This code does not run, please check example directory for working examples**

### Redirect Handling

### Code Generators

You must install the following protoc plugins which generate the appropriate codes
> 1. protoc-gen-gorony

```
go get -u github.com/ronaksoft/rony/cmd/protoc-gen-gorony

This is a protoc pluging which could generate codes based on the protobuf. In addition to
protobuffer syntax there are directive commands available which let the user generate more
customized code. 

```

---

## Shoulders

Rony is made of big and popular packages. Without these great libraries building Rony was not possible.

* [Hashicorp MemberList](https://github.com/hashicorp/memberlist)
* [SPF13 Cobra](https://github.com/spf13/cobra)
* [SPF13 Viper](https://github.com/spf13/viper)
* [Valyala FastHttp](https://github.com/valyala/fasthttp)
* [Uber Zap](https://go.uber.org/zap)
* [Gobwas Websocket](https://github.com/gobwas/ws)
* [Allegro BigCache](https://github.com/allegro/bigcache)
* [GoBuffalo Genny](https://github.com/gobuffalo/genny)
* [GoBuffalo Plush](https://github.com/gobuffalo/plush)

### Contribution

We need to make a clear and understandable documentation for Rony, so any help would be appreciated. We also appreciate benchmarking Rony against other platforms for common
scenarios and will be cited.

### TODOs

- [x] Auto-generate REST handlers
- [ ] Auto-generate test units
- [ ] Implement scenario tests using `edgetest` framework
- [ ] Write Examples and Tutorials
