package dummyGateway

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/errors"
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
	if config.Exposer != nil {
		config.Exposer(g)
	}
	return g, nil
}

// OpenConn opens a persistent connection to the gateway. For short lived use RPC or REST methods.
func (g *Gateway) OpenConn(
	connID uint64, persistent bool,
	onReceiveMessage func(connID uint64, streamID int64, data []byte),
	hdr ...*rony.KeyValue,
) {
	dConn := g.openConn(connID, persistent)
	dConn.onMessage = func(connID uint64, streamID int64, data []byte, _ map[string]string) {
		onReceiveMessage(connID, streamID, data)
	}

	g.ConnectHandler(dConn, hdr...)
}

func (g *Gateway) openConn(connID uint64, persistent bool) *Conn {
	dConn := &Conn{
		id:         connID,
		kv:         make(map[string]interface{}),
		persistent: persistent,
	}
	g.connsMtx.Lock()
	g.conns[connID] = dConn
	g.connsMtx.Unlock()

	atomic.AddInt32(&g.connsTotal, 1)
	return dConn
}

// CloseConn closed the connection
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

// RPC emulates sending an RPC command to the connection. It opens a non-persistent connection if connID
// does not exists
func (g *Gateway) RPC(connID uint64, streamID int64, data []byte) error {
	g.connsMtx.RLock()
	conn := g.conns[connID]
	g.connsMtx.RUnlock()
	if conn == nil {
		return errors.ErrConnectionNotExists
	}

	go g.MessageHandler(conn, streamID, data)
	return nil
}

// REST emulates a Http REST request.
func (g *Gateway) REST(connID uint64, method, path string, body []byte, kvs ...*rony.KeyValue) (respBody []byte, respHdr map[string]string) {
	conn := g.openConn(connID, false)
	conn.onMessage = func(connID uint64, streamID int64, data []byte, hdr map[string]string) {
		respHdr = hdr
		respBody = data
	}
	conn.method = method
	conn.path = path
	conn.body = body
	for _, kv := range kvs {
		conn.kv[kv.Key] = kv.Key
	}
	defer g.CloseConn(connID)

	g.MessageHandler(conn, 0, nil)
	return
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

func (g *Gateway) Protocol() rony.GatewayProtocol {
	return rony.Dummy
}
