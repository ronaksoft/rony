package edge

import (
	"github.com/ronaksoft/rony/internal/cluster"
	"github.com/ronaksoft/rony/internal/gateway"
	"github.com/ronaksoft/rony/internal/gateway/tcp/router"
	"github.com/ronaksoft/rony/internal/tunnel"
)

/*
   Creation Time: 2021 - Feb - 27
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type (
	Gateway         = gateway.Gateway
	GatewayProtocol = gateway.Protocol
	GatewayRouter   = router.Router
	Tunnel          = tunnel.Tunnel
	Cluster         = cluster.Cluster
	ClusterMode     = cluster.Mode
)

// Cluster Modes
const (
	// SingleReplica if set then each replica set is only one node. i.e. raft is OFF.
	SingleReplica ClusterMode = "singleReplica"
	// MultiReplica if set then each replica set is a raft cluster
	MultiReplica ClusterMode = "multiReplica"
)

// Gateway Protocols
const (
	Undefined GatewayProtocol = 0
	Dummy     GatewayProtocol = 1 << iota
	Http
	Websocket
	Quic
	Grpc
	// Mixed
	TCP GatewayProtocol = 0x0003 // Http & Websocket
)

// NewHttpRouter create a router which let you wrapper around your RPCs.
func NewHttpRouter() *router.Router {
	return router.New()
}