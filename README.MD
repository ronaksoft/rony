# RONY  **(Fast and Scalable RPC Framework)**


---
### Features
1. Gossip Protocol
2. Raft
3. ProtocolBuffer friendly rpc
4. Multiple gateways: Websocket, Http, Quic

---
### Performance
Rony is very fast and with very low overhead. In a non-raft environment it adds < 25us latency, and
in raft-enabled scenario around 1ms latency.

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
### Easy Setup
```
func main() {
    server := rony.NewEdgeServer(serverID, &dispatcher{},
    	rony.WithWebsocketGateway(websocketGateway.Config{
    		NewConnectionWorkers: 10,
    		MaxConcurrency:       1000,
    		MaxIdleTime:          0,
    		ListenAddress:        "0.0.0.0:0",
    	}),
    	rony.WithDataPath(filepath.Join("./_hdd", serverID)),
    	rony.WithReplicaSet(100, port*10, bootstrap),
    	rony.WithGossipPort(port),
    )
 
    server.AddHandler(msg.C_EchoRequest, EchoHandler)
    
    server.Run()
}

func EchoHandler(ctx *context.Context, in *msg.MessageEnvelope) {
    req := msg.EchoRequest {}
    res := msg.EchoResponse {}
    _ = req.Unmarshal(in.Message)

    res.Bool = req.Bool
    	res.Int = req.Int
    	res.Timestamp = time.Now().UnixNano()
    	res.Delay = res.Timestamp - req.Timestamp
    
    ctx.PushMessage(ctx.AuthID, in.RequestID, msg.C_EchoResponse, res)
}
```
**This code does not run, please check example directory for working examples**

---
### Summary

### Contribution
We need to make a clear and understandable documentation for Rony, so any help would be appreciated. We also appreciate benchmarking Rony 
against other platforms for common scenarios and will be cited.



