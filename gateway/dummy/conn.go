package dummy

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/tools"
	"sync"
)

/*
   Creation Time: 2020 - Oct - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Conn struct {
	id         uint64
	clientIP   string
	persistent bool
	mtx        sync.Mutex
	buf        *tools.LinkedList
	kv         map[string]interface{}
	onMessage  func(connID uint64, streamID int64, data []byte)
}

func (c *Conn) Get(key string) interface{} {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	return c.kv[key]
}

func (c *Conn) Set(key string, val interface{}) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.kv[key] = val
}

func (c *Conn) ConnID() uint64 {
	return c.id
}

func (c *Conn) ClientIP() string {
	return c.clientIP
}

func (c *Conn) Push(m *rony.MessageEnvelope) {
	c.buf.Append(m)
}

func (c *Conn) Pop() *rony.MessageEnvelope {
	v := c.buf.PickHeadData()
	if v != nil {
		return v.(*rony.MessageEnvelope)
	}
	return nil
}

func (c *Conn) SendBinary(streamID int64, data []byte) error {
	c.onMessage(c.id, streamID, data)
	return nil
}

func (c *Conn) Persistent() bool {
	return c.persistent
}

func (c *Conn) SetPersistent(b bool) {
	c.persistent = b
}
