package edge

import (
	"github.com/ronaksoft/rony"
	gossipCluster "github.com/ronaksoft/rony/cluster/gossip"
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

// WithGossipCluster enables the cluster in gossip mode. This mod is eventually consistent mode but there is
// no need to a central key-value store or any other 3rd party service to run the cluster
func WithGossipCluster(cfg gossipCluster.Config) Option {
	return func(edge *Server) {
		c := gossipCluster.New(cfg)
		c.ReplicaMessageHandler = edge.onReplicaMessage
		edge.cluster = c
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
