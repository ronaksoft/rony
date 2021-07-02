package rony

/*
   Creation Time: 2021 - Jul - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// Gateway defines the gateway interface where clients could connect
// and communicate with the edge server
type Gateway interface {
	Start()
	Run()
	Shutdown()
	GetConn(connID uint64) Conn
	Addr() []string
	Protocol() GatewayProtocol
}

type GatewayProtocol int32

const (
	Undefined GatewayProtocol = 0
	Dummy     GatewayProtocol = 1 << iota
	Http
	Websocket
	Quic
	Grpc
	TCP = Http | Websocket // Http & Websocket
)

var protocolNames = map[GatewayProtocol]string{
	Undefined: "Undefined",
	Dummy:     "Dummy",
	Http:      "Http",
	Websocket: "Websocket",
	Quic:      "Quic",
	Grpc:      "Grpc",
	TCP:       "TCP",
}

func (p GatewayProtocol) String() string {
	return protocolNames[p]
}

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
