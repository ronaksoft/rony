package tcpGateway

import (
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"github.com/valyala/fasthttp"
	"net"
)

/*
   Creation Time: 2019 - Nov - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// HttpConn
type HttpConn struct {
	gateway    *Gateway
	req        *fasthttp.RequestCtx
	buf        *tools.LinkedList
	AuthID     int64
	ClientIP   []byte
	ClientType []byte
}

func (c *HttpConn) Push(m *rony.MessageEnvelope) {
	c.buf.Append(m)
}

func (c *HttpConn) Pop() *rony.MessageEnvelope {
	v := c.buf.PickHeadData()
	if v != nil {
		return v.(*rony.MessageEnvelope)
	}
	return nil
}

func (c *HttpConn) GetAuthID() int64 {
	return c.AuthID
}

func (c *HttpConn) SetAuthID(authID int64) {
	c.AuthID = authID
}

func (c *HttpConn) GetConnID() uint64 {
	return c.req.ConnID()
}

func (c *HttpConn) GetClientIP() string {
	return net.IP(c.ClientIP).String()
}

func (c *HttpConn) SendBinary(streamID int64, data []byte) error {
	_, err := c.req.Write(data)
	return err
}

func (c *HttpConn) Persistent() bool {
	return false
}
