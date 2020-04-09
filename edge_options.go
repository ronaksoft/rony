package rony

import (
	"git.ronaksoftware.com/ronak/rony/gateway"
	httpGateway "git.ronaksoftware.com/ronak/rony/gateway/http"
	websocketGateway "git.ronaksoftware.com/ronak/rony/gateway/ws"
)

/*
   Creation Time: 2020 - Feb - 22
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// Option
type Option func(edge *EdgeServer)

// WithReplicaSet
func WithReplicaSet(replicaSet uint32, bindPort int, bootstrap bool) Option {
	return func(edge *EdgeServer) {
		edge.raftFSM = raftFSM{edge: edge}
		edge.replicaSet = replicaSet
		edge.raftEnabled = true
		edge.raftPort = bindPort
		edge.raftBootstrap = bootstrap
	}
}

// WithShardSet
func WithShardSet(shardSet uint32) Option {
	return func(edge *EdgeServer) {
		edge.shardSet = shardSet
	}
}

// WithGossipPort
func WithGossipPort(gossipPort int) Option {
	return func(edge *EdgeServer) {
		edge.gossipPort = gossipPort
	}
}

// WithCustomerConstructorName will be used to give a human readable names to messages
func WithCustomConstructorName(h func(constructor int64) (name string)) Option {
	return func(edge *EdgeServer) {
		edge.getConstructorName = func(constructor int64) string {
			// Lookup internal messages first
			n := ConstructorNames[constructor]
			if len(n) > 0 {
				return n
			}
			return h(constructor)
		}
	}
}

// WithDataPath set where the internal data for raft and gossip are stored.
func WithDataPath(path string) Option {
	return func(edge *EdgeServer) {
		edge.dataPath = path
	}
}

// WithWebsocketGateway set the gateway to websocket.
// Only one gateway could be set and if you set another gateway it panics on runtime.
func WithWebsocketGateway(config websocketGateway.Config) Option {
	return func(edge *EdgeServer) {
		if edge.gatewayProtocol != gateway.Undefined {
			panic(ErrGatewayAlreadyInitialized)
		}
		gatewayWebsocket, err := websocketGateway.New(config)
		if err != nil {
			panic(err)
		}
		gatewayWebsocket.MessageHandler = edge.onGatewayMessage
		gatewayWebsocket.ConnectHandler = edge.onConnect
		gatewayWebsocket.CloseHandler = edge.onClose
		gatewayWebsocket.FlushFunc = edge.onFlush
		edge.gatewayProtocol = gateway.Websocket
		edge.gateway = gatewayWebsocket
		return
	}
}

// WithHttpGateway set the gateway to http
// Only one gateway could be set and if you set another gateway it panics on runtime.
func WithHttpGateway(config httpGateway.Config) Option {
	return func(edge *EdgeServer) {
		if edge.gatewayProtocol != gateway.Undefined {
			panic(ErrGatewayAlreadyInitialized)
		}
		gatewayHttp := httpGateway.New(config)
		gatewayHttp.MessageHandler = edge.onGatewayMessage
		gatewayHttp.FlushFunc = edge.onFlush

		edge.gatewayProtocol = gateway.HTTP
		edge.gateway = gatewayHttp
		return
	}
}
