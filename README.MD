# RONY  **(Fast and Scalable RPC Framework)**

### Features
1. Gossip Protocol
2. Raft
3. ProtocolBuffer friendly rpc
4. Multiple gateways: Websocket, Http, Quic


### Performance
Rony is very fast and with very low overhead. In a non-raft environment it adds < 40us latency, and
in raft-enabled scenario around 1ms latency.

Rony has negligible memory footprint by reusing buffers and pooling.

### Easy Setup
```
func main() {
    server := rony.NewEdgeServer(bundleID, instanceID, &dispatcher{},
    	rony.WithWebsocketGateway(websocketGateway.Config{
    		NewConnectionWorkers: 10,
    		MaxConcurrency:       1000,
    		MaxIdleTime:          0,
    		ListenAddress:        "0.0.0.0:0",
    	}),
    	rony.WithDataPath(filepath.Join("./_hdd", serverID)),
    	rony.WithRaft(port*10, bootstrap),
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


