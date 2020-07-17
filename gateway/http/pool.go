package httpGateway

import (
	"git.ronaksoftware.com/ronak/rony/tools"
	"github.com/valyala/fasthttp"
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

var connPool sync.Pool

func acquireConn(gw *Gateway, req *fasthttp.RequestCtx) *Conn {
	c, ok := connPool.Get().(*Conn)
	if !ok {
		return &Conn{
			gateway: gw,
			req:     req,
			buf:     tools.NewLinkedList(),
		}
	}
	c.gateway = gw
	c.req = req
	return c
}

func releaseConn(c *Conn) {
	c.ClientIP = c.ClientIP[:0]
	c.ClientType = c.ClientType[:0]
	c.AuthID = 0
	c.buf.Reset()
	connPool.Put(c)
}
