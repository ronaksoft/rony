package edge

import (
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/gateway"
	httpGateway "git.ronaksoftware.com/ronak/rony/gateway/http"
	tcpGateway "git.ronaksoftware.com/ronak/rony/gateway/tcp"
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
type Option func(edge *Server)

// WithReplicaSet
func WithReplicaSet(replicaSet uint64, bindPort int, bootstrap bool) Option {
	return func(edge *Server) {
		edge.raftFSM = raftFSM{edge: edge}
		edge.replicaSet = replicaSet
		edge.raftEnabled = true
		edge.raftPort = bindPort
		edge.raftBootstrap = bootstrap
	}
}

// WithGossipPort
func WithGossipPort(gossipPort int) Option {
	return func(edge *Server) {
		edge.gossipPort = gossipPort
	}
}

// WithCustomerConstructorName will be used to give a human readable names to messages
func WithCustomConstructorName(h func(constructor int64) (name string)) Option {
	return func(edge *Server) {
		edge.getConstructorName = func(constructor int64) string {
			// Lookup internal messages first
			n := rony.ConstructorNames[constructor]
			if len(n) > 0 {
				return n
			}
			return h(constructor)
		}
	}
}

// WithDataPath set where the internal data for raft and gossip are stored.
func WithDataPath(path string) Option {
	return func(edge *Server) {
		edge.dataPath = path
	}
}

// WithWebsocketGateway set the gateway to websocket.
// Only one gateway could be set and if you set another gateway it panics on runtime.
func WithWebsocketGateway(config websocketGateway.Config) Option {
	return func(edge *Server) {
		if edge.gatewayProtocol != gateway.Undefined {
			panic(rony.ErrGatewayAlreadyInitialized)
		}
		gatewayWebsocket, err := websocketGateway.New(config)
		if err != nil {
			panic(err)
		}
		gatewayWebsocket.MessageHandler = edge.HandleGatewayMessage
		gatewayWebsocket.ConnectHandler = edge.onConnect
		gatewayWebsocket.CloseHandler = edge.onClose
		edge.gatewayProtocol = gateway.Websocket
		edge.gateway = gatewayWebsocket
		return
	}
}

// WithHttpGateway set the gateway to http
// Only one gateway could be set and if you set another gateway it panics on runtime.
func WithHttpGateway(config httpGateway.Config) Option {
	return func(edge *Server) {
		if edge.gatewayProtocol != gateway.Undefined {
			panic(rony.ErrGatewayAlreadyInitialized)
		}
		gatewayHttp, err := httpGateway.New(config)
		if err != nil {
			panic(err)
		}
		gatewayHttp.MessageHandler = edge.HandleGatewayMessage
		edge.gatewayProtocol = gateway.HTTP
		edge.gateway = gatewayHttp
		return
	}
}

// WithTcpGateway set the gateway to tcp which can support http and/or websocket
// Only one gateway could be set and if you set another gateway it panics on runtime.
func WithTcpGateway(config tcpGateway.Config) Option {
	return func(edge *Server) {
		if edge.gatewayProtocol != gateway.Undefined {
			panic(rony.ErrGatewayAlreadyInitialized)
		}
		gatewayTcp, err := tcpGateway.New(config)
		if err != nil {
			panic(err)
		}
		gatewayTcp.MessageHandler = edge.HandleGatewayMessage
		gatewayTcp.ConnectHandler = edge.onConnect
		gatewayTcp.CloseHandler = edge.onClose
		edge.gatewayProtocol = gateway.TCP
		edge.gateway = gatewayTcp
		return
	}
}
