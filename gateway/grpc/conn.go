package grpcGateway

import (
	"git.ronaksoftware.com/ronak/rony"
)

type Conn struct {
	ConnID   uint64
	AuthID   int64
	ClientIP string

	srv rony.Rony_PipelineServer
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

func (c *Conn) Persistent() bool {
	return true
}

func (c *Conn) Push(m *rony.MessageEnvelope) {
	panic("implement me")
}

func (c *Conn) Pop() *rony.MessageEnvelope {
	panic("implement me")
}

func (c *Conn) SendBinary(streamID int64, data []byte) error {
	// c.srv.Send()
	return nil
}
