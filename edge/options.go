package edge

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/gateway"
	gossipCluster "github.com/ronaksoft/rony/internal/cluster/gossip"
	dummyGateway "github.com/ronaksoft/rony/internal/gateway/dummy"
	tcpGateway "github.com/ronaksoft/rony/internal/gateway/tcp"
	udpTunnel "github.com/ronaksoft/rony/internal/tunnel/udp"
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

type GossipClusterConfig gossipCluster.Config

// WithGossipCluster enables the cluster in gossip mode. This mod is eventually consistent mode but there is
// no need to a central key-value store or any other 3rd party service to run the cluster
func WithGossipCluster(cfg GossipClusterConfig) Option {
	return func(edge *Server) {
		c := gossipCluster.New(gossipCluster.Config(cfg))
		c.ReplicaMessageHandler = edge.onReplicaMessage
		edge.cluster = c
	}
}

type TcpGatewayConfig tcpGateway.Config

// WithTcpGateway set the gateway to tcp which can support http and/or websocket
// Only one gateway could be set and if you set another gateway it panics on runtime.
func WithTcpGateway(config TcpGatewayConfig) Option {
	return func(edge *Server) {
		if edge.gatewayProtocol != gateway.Undefined {
			panic(rony.ErrGatewayAlreadyInitialized)
		}
		gatewayTcp, err := tcpGateway.New(tcpGateway.Config(config))
		if err != nil {
			panic(err)
		}
		gatewayTcp.MessageHandler = edge.onGatewayMessage
		gatewayTcp.ConnectHandler = edge.onGatewayConnect
		gatewayTcp.CloseHandler = edge.onGatewayClose
		edge.gatewayProtocol = gateway.TCP
		edge.gateway = gatewayTcp
		return
	}
}

type DummyGatewayConfig dummyGateway.Config

// WithTestGateway set the gateway to a dummy gateway which is useful for writing tests.
// Only one gateway could be set and if you set another gateway it panics on runtime.
func WithTestGateway(config DummyGatewayConfig) Option {
	return func(edge *Server) {
		if edge.gatewayProtocol != gateway.Undefined {
			panic(rony.ErrGatewayAlreadyInitialized)
		}
		gatewayDummy, err := dummyGateway.New(dummyGateway.Config(config))
		if err != nil {
			panic(err)
		}
		gatewayDummy.MessageHandler = edge.onGatewayMessage
		gatewayDummy.ConnectHandler = edge.onGatewayConnect
		gatewayDummy.CloseHandler = edge.onGatewayClose
		edge.gatewayProtocol = gateway.Dummy
		edge.gateway = gatewayDummy
		return
	}
}

type UdpTunnelConfig udpTunnel.Config

// WithUdpTunnel set the tunnel to a udp based tunnel which provides communication channel between
// edge servers.
func WithUdpTunnel(config UdpTunnelConfig) Option {
	return func(edge *Server) {
		tunnelUDP := udpTunnel.New(udpTunnel.Config(config))
		tunnelUDP.MessageHandler = edge.onTunnelMessage
		edge.tunnel = tunnelUDP
	}
}
