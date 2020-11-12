package tcp

import (
	"github.com/gobwas/ws"
	"github.com/mailru/easygo/netpoll"
	wsutil "github.com/ronaksoft/rony/gateway/tcp/util"
	"github.com/ronaksoft/rony/tools"
	"github.com/valyala/fasthttp"
	"net"
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

var httpConnPool sync.Pool

func acquireHttpConn(gw *Gateway, req *fasthttp.RequestCtx) *httpConn {
	c, ok := httpConnPool.Get().(*httpConn)
	if !ok {
		return &httpConn{
			gateway: gw,
			req:     req,
			buf:     tools.NewLinkedList(),
		}
	}
	c.gateway = gw
	c.req = req
	return c
}

func releaseHttpConn(c *httpConn) {
	c.clientIP = c.clientIP[:0]
	c.clientType = c.clientType[:0]
	for k := range c.kv {
		delete(c.kv, k)
	}
	c.buf.Reset()
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

var websocketConnPool sync.Pool

func acquireWebsocketConn(gw *Gateway, connID uint64, conn net.Conn, desc *netpoll.Desc) *websocketConn {
	c, ok := websocketConnPool.Get().(*websocketConn)
	if !ok {
		return &websocketConn{
			connID:       connID,
			buf:          tools.NewLinkedList(),
			gateway:      gw,
			conn:         conn,
			desc:         desc,
			closed:       false,
			lastActivity: tools.TimeUnix(),
		}
	}
	c.gateway = gw
	c.connID = connID
	c.desc = desc
	c.conn = conn
	c.lastActivity = tools.TimeUnix()
	return c
}

func releaseWebsocketConn(wc *websocketConn) {
	wc.clientIP = wc.clientIP[:0]
	wc.buf.Reset()
	wc.conn = nil
	for k := range wc.kv {
		delete(wc.kv, k)
	}
	wc.closed = false
	websocketConnPool.Put(wc)
}

var websocketMessagePool sync.Pool

func acquireWebsocketMessage() []wsutil.Message {
	x, ok := websocketMessagePool.Get().([]wsutil.Message)
	if !ok {
		return make([]wsutil.Message, 0, 8)
	}
	return x
}

func releaseWebsocketMessage(x []wsutil.Message) {
	websocketMessagePool.Put(x)
}
