package edge

import (
	"github.com/ronaksoft/rony/internal/cluster"
	"github.com/ronaksoft/rony/internal/gateway"
	tcpGateway "github.com/ronaksoft/rony/internal/gateway/tcp"
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
	Tunnel          = tunnel.Tunnel
	Cluster         = cluster.Cluster
	ClusterMode     = cluster.Mode
	ProxyHandle     = gateway.ProxyHandle
	HttpProxy       = tcpGateway.HttpProxy
	HttpRequest     = gateway.RequestCtx
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

// HTTP methods were copied from net/http.
const (
	MethodWild    = "*"
	MethodGet     = "GET"     // RFC 7231, 4.3.1
	MethodHead    = "HEAD"    // RFC 7231, 4.3.2
	MethodPost    = "POST"    // RFC 7231, 4.3.3
	MethodPut     = "PUT"     // RFC 7231, 4.3.4
	MethodPatch   = "PATCH"   // RFC 5789
	MethodDelete  = "DELETE"  // RFC 7231, 4.3.5
	MethodConnect = "CONNECT" // RFC 7231, 4.3.6
	MethodOptions = "OPTIONS" // RFC 7231, 4.3.7
	MethodTrace   = "TRACE"   // RFC 7231, 4.3.8
)
