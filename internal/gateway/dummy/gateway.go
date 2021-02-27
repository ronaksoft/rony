package dummyGateway

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/gateway"
	"sync"
	"sync/atomic"
)

/*
   Creation Time: 2020 - Oct - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Config struct {
	Exposer func(gw *Gateway)
}

type Gateway struct {
	gateway.ConnectHandler
	gateway.MessageHandler
	gateway.CloseHandler

	conns      map[uint64]*Conn
	connsMtx   sync.RWMutex
	connsTotal int32
}

func New(config Config) (*Gateway, error) {
	g := &Gateway{
		conns: make(map[uint64]*Conn, 8192),
	}

	// Call the exposer make caller have access to this gateway object
	config.Exposer(g)
	return g, nil
}

func (g *Gateway) OpenConn(connID uint64, onReceiveMessage func(connID uint64, streamID int64, data []byte), kvs ...*rony.KeyValue) {
	dConn := &Conn{
		id:        connID,
		kv:        make(map[string]interface{}),
		onMessage: onReceiveMessage,
	}
	g.connsMtx.Lock()
	g.conns[connID] = dConn
	g.connsMtx.Unlock()
	g.ConnectHandler(dConn, kvs...)
	atomic.AddInt32(&g.connsTotal, 1)
	return
}

func (g *Gateway) CloseConn(connID uint64) {
	g.connsMtx.Lock()
	c := g.conns[connID]
	delete(g.conns, connID)
	g.connsMtx.Unlock()
	if c != nil {
		g.CloseHandler(c)
		atomic.AddInt32(&g.connsTotal, -1)
	}
}

func (g *Gateway) SendToConn(connID uint64, streamID int64, data []byte) {
	g.connsMtx.RLock()
	conn := g.conns[connID]
	g.connsMtx.RUnlock()
	if conn == nil {
		return
	}

	g.MessageHandler(conn, streamID, data)
}

func (g *Gateway) Start() {
	// Do nothing
}

func (g *Gateway) Run() {
	// Do nothing
}

func (g *Gateway) Shutdown() {
	// Do nothing
}

func (g *Gateway) GetConn(connID uint64) rony.Conn {
	g.connsMtx.RLock()
	conn := g.conns[connID]
	g.connsMtx.RUnlock()
	if conn == nil {
		return nil
	}
	return conn
}

func (g *Gateway) Addr() []string {
	return []string{"TEST"}
}
