package grpcGateway

import (
	"git.ronaksoftware.com/ronak/rony/errors"
	"git.ronaksoftware.com/ronak/rony/gateway"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/msg"
	"go.uber.org/zap"
)

type Conn struct {
	ConnID   uint64
	AuthID   int64
	UserID   int64
	AuthKey  []byte
	ClientIP string

	srv       msg.Edge_ServeRequestServer
	flushFunc FlushFunc
}

func (c *Conn) GetConnID() uint64 {
	return c.ConnID
}

func (c *Conn) GetClientIP() string {
	return c.ClientIP
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

func (c *Conn) IncServerSeq(v int64) int64 {
	panic("implement me")
}

func (c *Conn) Persistent() bool {
	return true
}

func (c *Conn) SendProto(streamID int64, message gateway.ProtoBufferMessage) error {
	if message == nil {
		log.Warn("Received nil MessageEnvelope in GRPC Connection SendProto.",
			zap.Uint64("ConnID", c.ConnID),
		)
		return errors.ErrInvalidData
	}

	err := c.srv.Send(message.(*msg.ProtoMessage))
	if err != nil {
		log.Warn("Failed to send message envelope to stream.",
			zap.Uint64("ConnID", c.ConnID),
			zap.Int64("AuthID", c.AuthID),
			zap.Int64("UserID", c.UserID),
			zap.Error(err),
		)
	}
	return err
}

func (c *Conn) Flush() {
	envelope := c.flushFunc(c)
	_ = c.SendProto(0, envelope)
}
