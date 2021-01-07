package udpTunnel

import (
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
	"github.com/tidwall/evio"
	"sync"
)

/*
   Creation Time: 2021 - Jan - 07
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type udpConn struct {
	id  uint64
	c   evio.Conn
	buf *tools.LinkedList
	// KV Store
	mtx sync.RWMutex
	kv  map[string]interface{}
}

func newConn(connID uint64, c evio.Conn) *udpConn {
	uc := &udpConn{
		id:  connID,
		c:   c,
		buf: tools.NewLinkedList(),
	}
	c.SetContext(uc)
	return uc
}

func (u *udpConn) Push(data []byte) {
	d := pools.BytesBuffer.GetCap(len(data))
	d.Fill(data)
	u.buf.Append(d)
}

func (u *udpConn) Pop() *pools.ByteBuffer {
	v := u.buf.PickHeadData()
	if v != nil {
		return v.(*pools.ByteBuffer)
	}
	return nil
}

func (u *udpConn) ConnID() uint64 {
	return u.id
}

func (u *udpConn) ClientIP() string {
	return u.c.RemoteAddr().String()
}

func (u *udpConn) SendBinary(streamID int64, data []byte) error {
	u.Push(data)
	u.c.Wake()
	return nil
}

func (u *udpConn) Persistent() bool {
	return false
}

func (u *udpConn) Get(key string) interface{} {
	u.mtx.RLock()
	v := u.kv[key]
	u.mtx.RUnlock()
	return v
}

func (u *udpConn) Set(key string, val interface{}) {
	u.mtx.Lock()
	u.kv[key] = val
	u.mtx.Unlock()
}
