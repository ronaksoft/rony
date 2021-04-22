package tcpGateway

import (
	"github.com/ronaksoft/rony/internal/gateway"
	"github.com/ronaksoft/rony/internal/metrics"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
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
	ctx        *gateway.RequestCtx
	clientIP   []byte
	clientType []byte
	mtx        tools.SpinLock
	kv         map[string]interface{}
	proxy      gateway.ProxyHandle
}

func (c *httpConn) Get(key string) interface{} {
	c.mtx.Lock()
	v := c.kv[key]
	c.mtx.Unlock()
	return v
}

func (c *httpConn) Set(key string, val interface{}) {
	c.mtx.Lock()
	c.kv[key] = val
	c.mtx.Unlock()
}

func (c *httpConn) ConnID() uint64 {
	return c.ctx.ConnID()
}

func (c *httpConn) ClientIP() string {
	return string(c.clientIP)
}

func (c *httpConn) SetClientIP(ip []byte) {
	c.clientIP = append(c.clientIP[:0], ip...)
}

func (c *httpConn) SetClientType(ct []byte) {
	c.clientType = append(c.clientType[:0], ct...)
}

func (c *httpConn) SendBinary(streamID int64, data []byte) (err error) {
	if c.proxy != nil {
		d, hdr := c.proxy.OnResponse(data)
		for k, v := range hdr {
			c.ctx.Response.Header.Add(k, v)
		}
		_, err = c.ctx.Write(*d.Bytes())
		pools.Buffer.Put(d)
	} else {
		_, err = c.ctx.Write(data)
	}
	metrics.IncCounter(metrics.CntGatewayOutgoingHttpMessage)
	return err
}

func (c *httpConn) Persistent() bool {
	return false
}
