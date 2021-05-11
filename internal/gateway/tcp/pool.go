package tcpGateway

import (
	"github.com/gobwas/ws"
	"github.com/panjf2000/ants/v2"
	"github.com/ronaksoft/rony/internal/gateway"
	wsutil "github.com/ronaksoft/rony/internal/gateway/tcp/util"
	"sync"
)

/*
   Creation Time: 2020 - Apr - 17
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var goPoolB, goPoolNB *ants.Pool

var httpConnPool sync.Pool

func acquireHttpConn(gw *Gateway, req *gateway.RequestCtx) *httpConn {
	c, ok := httpConnPool.Get().(*httpConn)
	if !ok {
		return &httpConn{
			gateway: gw,
			ctx:     req,
			kv:      make(map[string]interface{}, 4),
		}
	}
	c.gateway = gw
	c.ctx = req
	return c
}

func releaseHttpConn(c *httpConn) {
	c.clientIP = c.clientIP[:0]
	c.clientType = c.clientType[:0]
	c.proxy = nil
	for k := range c.kv {
		delete(c.kv, k)
	}
	httpConnPool.Put(c)
}

var writeRequestPool sync.Pool

func acquireWriteRequest(wc *websocketConn, opCode ws.OpCode) *writeRequest {
	wr, ok := writeRequestPool.Get().(*writeRequest)
	if !ok {
		return &writeRequest{
			wc:     wc,
			opCode: opCode,
		}
	}
	wr.wc = wc
	wr.opCode = opCode
	return wr
}

func releaseWriteRequest(wr *writeRequest) {
	wr.wc = nil
	wr.payload = wr.payload[:0]
	writeRequestPool.Put(wr)
}

var websocketMessagePool sync.Pool

func acquireWebsocketMessage() *[]wsutil.Message {
	x, ok := websocketMessagePool.Get().(*[]wsutil.Message)
	if !ok {
		arr := make([]wsutil.Message, 0, 8)
		return &arr
	}
	return x
}

func releaseWebsocketMessage(x *[]wsutil.Message) {
	websocketMessagePool.Put(x)
}
