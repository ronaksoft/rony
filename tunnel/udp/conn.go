package udpTunnel

import (
	"github.com/panjf2000/gnet"
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
	id uint64
	c  gnet.Conn
	// KV Store
	mtx sync.RWMutex
	kv  map[string]interface{}
}

func newConn(connID uint64, c gnet.Conn) *udpConn {
	uc := &udpConn{
		id: connID,
		c:  c,
	}
	c.SetContext(uc)
	return uc
}

func (u *udpConn) ConnID() uint64 {
	return u.id
}

func (u *udpConn) ClientIP() string {
	if u == nil || u.c == nil {
		return ""
	}
	return u.c.RemoteAddr().String()
}

func (u *udpConn) SendBinary(streamID int64, data []byte) error {
	if u == nil || u.c == nil {
		return nil
	}
	return u.c.SendTo(data)
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
