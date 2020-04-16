package httpGateway

import (
	"git.ronaksoftware.com/ronak/rony"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

/*
   Creation Time: 2019 - Nov - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type Conn struct {
	gateway    *Gateway
	req        *fasthttp.RequestCtx
	buf        chan *rony.MessageEnvelope
	AuthID     int64
	AuthKey    []byte
	UserID     int64
	ConnID     uint64
	ClientIP   string
	ClientType string
}

func (c *Conn) Flush() {
	// Read the flush function
	bytesSlice := c.gateway.FlushFunc(c)

	_, err := c.req.Write(bytesSlice)
	if ce := log.Check(log.DebugLevel, "Error On Write To Websocket Conn"); ce != nil {
		ce.Write(
			zap.Uint64("ConnID", c.ConnID),
			zap.Int64("authID", c.AuthID),
			zap.Error(err),
		)
	}
	return
}

func (c *Conn) Push(m *rony.MessageEnvelope) {
	c.buf <- m
}

func (c *Conn) Pop() *rony.MessageEnvelope {
	return <-c.buf
}

func (c *Conn) GetAuthID() int64 {
	return c.AuthID
}

func (c *Conn) SetAuthID(authID int64) {
	c.AuthID = authID
}

func (c *Conn) GetConnID() uint64 {
	return c.ConnID
}

func (c *Conn) GetClientIP() string {
	return c.ClientIP
}

func (c *Conn) SendBinary(streamID int64, data []byte) error {
	_, err := c.req.Write(data)
	return err
}

func (c *Conn) Persistent() bool {
	return false
}
