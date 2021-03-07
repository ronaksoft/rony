package gateway

import (
	"github.com/ronaksoft/rony"
	"github.com/valyala/fasthttp"
)

/*
   Creation Time: 2019 - Aug - 31
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Protocol int32

const (
	Undefined Protocol = 0
	Dummy     Protocol = 1 << iota
	Http
	Websocket
	Quic
	Grpc
	// Mixed
	TCP Protocol = 0x0003 // Http & Websocket
)

var protocolNames = map[Protocol]string{
	Undefined: "Undefined",
	Dummy:     "Dummy",
	Http:      "Http",
	Websocket: "Websocket",
	Quic:      "Quic",
	Grpc:      "Grpc",
	TCP:       "TCP",
}

func (p Protocol) String() string {
	return protocolNames[p]
}

// Gateway defines the gateway interface where clients could connect
// and communicate with the edge server
type Gateway interface {
	Start()
	Run()
	Shutdown()
	GetConn(connID uint64) rony.Conn
	Addr() []string
}

type (
	RequestCtx     = fasthttp.RequestCtx
	ConnectHandler = func(c rony.Conn, kvs ...*rony.KeyValue)
	MessageHandler = func(c rony.Conn, streamID int64, data []byte)
	CloseHandler   = func(c rony.Conn)
)

type ProxyHandle interface {
	OnRequest(conn rony.Conn, ctx *RequestCtx) []byte
	OnResponse(data []byte) ([]byte, map[string]string)
}
