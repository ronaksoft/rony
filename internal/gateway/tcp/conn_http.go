package tcpGateway

import (
	"github.com/valyala/fasthttp"
	"net"
	"sync"
)

/*
   Creation Time: 2019 - Nov - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// httpConn
type httpConn struct {
	gateway    *Gateway
	req        *fasthttp.RequestCtx
	clientIP   []byte
	clientType []byte
	mtx        sync.RWMutex
	kv         map[string]interface{}
}

func (c *httpConn) Get(key string) interface{} {
	c.mtx.RLock()
	v := c.kv[key]
	c.mtx.RUnlock()
	return v
}

func (c *httpConn) Set(key string, val interface{}) {
	c.mtx.Lock()
	c.kv[key] = val
	c.mtx.Unlock()
}

func (c *httpConn) ConnID() uint64 {
	return c.req.ConnID()
}

func (c *httpConn) ClientIP() string {
	return net.IP(c.clientIP).String()
}

func (c *httpConn) SetClientIP(ip []byte) {
	c.clientIP = append(c.clientIP[:0], ip...)
}

func (c *httpConn) SetClientType(ct []byte) {
	c.clientType = append(c.clientType[:0], ct...)
}

func (c *httpConn) SendBinary(streamID int64, data []byte) error {
	_, err := c.req.Write(data)
	return err
}

func (c *httpConn) Persistent() bool {
	return false
}
