package rony

import (
	"github.com/ronaksoft/rony/internal/cluster"
	"github.com/ronaksoft/rony/internal/gateway"
	"github.com/ronaksoft/rony/internal/tunnel"
)

/*
   Creation Time: 2021 - Jan - 07
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type (
	GatewayBodyWriter   = gateway.BodyWriter
	GatewayHeaderWriter = gateway.HeaderWriter
	Gateway             = gateway.Gateway
	Cluster             = cluster.Cluster
	Tunnel              = tunnel.Tunnel
)

// Conn defines the Connection interface
type Conn interface {
	ConnID() uint64
	ClientIP() string
	SendBinary(streamID int64, data []byte) error
	// Persistent returns FALSE if this connection will be closed when edge.DispatchCtx has been done. i.e. HTTP connections
	// It returns TRUE if this connection still alive when edge.DispatchCtx has been done. i.e. WebSocket connections
	Persistent() bool
	Get(key string) interface{}
	Set(key string, val interface{})
}
