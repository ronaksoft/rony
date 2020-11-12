package edge

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/gateway"
	dummyGateway "github.com/ronaksoft/rony/gateway/dummy"
	tcpGateway "github.com/ronaksoft/rony/gateway/tcp"
)

/*
   Creation Time: 2020 - Feb - 22
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// Option
type Option func(edge *Server)

// WithReplicaSet
func WithReplicaSet(replicaSet uint64, bindPort int, bootstrap bool) Option {
	if replicaSet == 0 {
		panic("replica-set could not be set zero")
	}
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
		gatewayTcp.MessageHandler = edge.onGatewayMessage
		gatewayTcp.ConnectHandler = edge.onConnect
		gatewayTcp.CloseHandler = edge.onClose
		edge.gatewayProtocol = gateway.TCP
		edge.gateway = gatewayTcp
		return
	}
}

// WithTestGateway set the gateway to a dummy gateway which is useful for writing tests.
// Only one gateway could be set and if you set another gateway it panics on runtime.
func WithTestGateway(config dummyGateway.Config) Option {
	return func(edge *Server) {
		if edge.gatewayProtocol != gateway.Undefined {
			panic(rony.ErrGatewayAlreadyInitialized)
		}
		gatewayDummy, err := dummyGateway.New(config)
		if err != nil {
			panic(err)
		}
		gatewayDummy.MessageHandler = edge.onGatewayMessage
		gatewayDummy.ConnectHandler = edge.onConnect
		gatewayDummy.CloseHandler = edge.onClose
		edge.gatewayProtocol = gateway.Dummy
		edge.gateway = gatewayDummy
		return
	}
}
