package edge

import (
	"github.com/ronaksoft/rony/gateway"
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
func (edge *Server) OnGatewayMessage(conn gateway.Conn, streamID int64, data []byte, kvs ...gateway.KeyValue) {
	edge.onGatewayMessage(conn, streamID, data, kvs...)
}
