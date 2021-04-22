package edge

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/gateway/tcp"
)

/*
   Creation Time: 2020 - Nov - 09
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// OnGatewayMessage is exposed only for test packages
func (edge *Server) OnGatewayMessage(conn rony.Conn, streamID int64, data []byte, bypass bool) {
	edge.onGatewayMessage(conn, streamID, data, bypass)
}

// GatewayConns is exposed only for test packages
func (edge *Server) GatewayConns() int {
	g, _ := edge.gateway.(*tcpGateway.Gateway)
	if g == nil {
		return 0
	}
	return g.TotalConnections()
}
