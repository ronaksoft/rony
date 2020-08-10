package tcpGateway

import (
	"git.ronaksoftware.com/ronak/rony/internal/tools"
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

func acquireHttpConn(gw *Gateway, req *fasthttp.RequestCtx) *HttpConn {
	c, ok := connPool.Get().(*HttpConn)
	if !ok {
		return &HttpConn{
			gateway: gw,
			req:     req,
			buf:     tools.NewLinkedList(),
		}
	}
	c.gateway = gw
	c.req = req
	return c
}

func releaseHttpConn(c *HttpConn) {
	c.ClientIP = c.ClientIP[:0]
	c.ClientType = c.ClientType[:0]
	c.AuthID = 0
	c.buf.Reset()
	connPool.Put(c)
}
