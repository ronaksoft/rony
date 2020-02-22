package httpGateway

import (
	"git.ronaksoftware.com/ronak/rony/gateway"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"github.com/gobwas/pool/pbytes"
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
	AuthID     int64
	AuthKey    []byte
	UserID     int64
	ConnID     uint64
	ClientIP   string
	ClientType string
}

func (c *Conn) IncServerSeq(int64) int64 {
	return 1
}

func (c *Conn) Flush() {
	// Read the flush function
	bytesSlice := c.gateway.FlushFunc(c)

	for idx := 0; idx < len(bytesSlice); idx++ {
		_, err := c.req.Write(bytesSlice[idx])
		if ce := log.Check(log.DebugLevel, "Error On Write To Websocket Conn"); ce != nil {
			ce.Write(
				zap.Uint64("ConnID", c.ConnID),
				zap.Int64("AuthID", c.AuthID),
				zap.Error(err),
			)
		}
	}
	return
}

func (c *Conn) GetAuthID() int64 {
	return c.AuthID
}

func (c *Conn) SetAuthID(authID int64) {
	c.AuthID = authID
}

func (c *Conn) GetAuthKey() []byte {
	return c.AuthKey
}

func (c *Conn) SetAuthKey(authKey []byte) {
	c.AuthKey = authKey
}

func (c *Conn) GetUserID() int64 {
	return c.UserID
}

func (c *Conn) SetUserID(userID int64) {
	c.UserID = userID
}

func (c *Conn) GetConnID() uint64 {
	return c.ConnID
}

func (c *Conn) GetClientIP() string {
	return c.ClientIP
}

func (c *Conn) SendProto(streamID int64, pm gateway.ProtoBufferMessage) error {
	protoBytes := pbytes.GetLen(pm.Size())
	_, _ = pm.MarshalTo(protoBytes)
	_, err := c.req.Write(protoBytes)
	pbytes.Put(protoBytes)
	return err
}

func (c *Conn) Persistent() bool {
	return false
}
