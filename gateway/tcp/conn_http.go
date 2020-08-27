package tcpGateway

import (
	"git.ronaksoft.com/ronak/rony"
	"git.ronaksoft.com/ronak/rony/tools"
	"github.com/valyala/fasthttp"
	"net"
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
	buf        *tools.LinkedList
	authID     int64
	clientIP   []byte
	clientType []byte
}

func (c *httpConn) Push(m *rony.MessageEnvelope) {
	c.buf.Append(m)
}

func (c *httpConn) Pop() *rony.MessageEnvelope {
	v := c.buf.PickHeadData()
	if v != nil {
		return v.(*rony.MessageEnvelope)
	}
	return nil
}

func (c *httpConn) GetAuthID() int64 {
	return c.authID
}

func (c *httpConn) SetAuthID(authID int64) {
	c.authID = authID
}

func (c *httpConn) GetConnID() uint64 {
	return c.req.ConnID()
}

func (c *httpConn) GetClientIP() string {
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
