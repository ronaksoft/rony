//go:build !windows && !appengine
// +build !windows,!appengine

package tcpGateway

import (
	"encoding/binary"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/allegro/bigcache/v2"
	"github.com/gobwas/ws"
	"github.com/mailru/easygo/netpoll"
	wsutil "github.com/ronaksoft/rony/internal/gateway/tcp/util"
	"github.com/ronaksoft/rony/internal/metrics"
	"github.com/ronaksoft/rony/log"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
	"go.uber.org/zap"
)

/*
   Creation Time: 2019 - Feb - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// websocketConn
type websocketConn struct {
	mtx      sync.Mutex
	connID   uint64
	clientIP []byte

	// KV Store
	kvLock tools.SpinLock
	kv     map[string]interface{}

	// Internals
	gateway      *Gateway
	lastActivity int64
	conn         net.Conn
	desc         *netpoll.Desc
	closed       bool
	startTime    int64
}

func newWebsocketConn(g *Gateway, conn net.Conn, clientIP []byte) (*websocketConn, error) {
	desc, err := netpoll.Handle(conn,
		netpoll.EventRead|netpoll.EventHup|netpoll.EventOneShot,
	)
	if err != nil {
		return nil, err
	}

	// Increment total connection counter and connection ID
	totalConns := atomic.AddInt32(&g.connsTotal, 1)
	connID := atomic.AddUint64(&g.connsLastID, 1)
	wsConn := &websocketConn{
		connID:       connID,
		gateway:      g,
		conn:         conn,
		desc:         desc,
		closed:       false,
		kv:           make(map[string]interface{}, 4),
		lastActivity: tools.CPUTicks(),
	}
	wsConn.SetClientIP(clientIP)

	g.connsMtx.Lock()
	g.conns[connID] = wsConn
	g.connsMtx.Unlock()
	if ce := g.cfg.Logger.Check(log.DebugLevel, "websocket connection created"); ce != nil {
		ce.Write(
			zap.Uint64("ConnID", connID),
			zap.String("IP", wsConn.ClientIP()),
			zap.Int32("Total", totalConns),
		)
	}
	g.connGC.monitorConnection(connID)

	return wsConn, nil
}

func (wc *websocketConn) registerDesc() error {
	atomic.StoreInt64(&wc.startTime, tools.CPUTicks())
	err := wc.gateway.poller.Start(wc.desc, wc.startEvent)
	if err != nil {
		if err != netpoll.ErrRegistered {
			wc.release(1)

			return err
		}
		_ = wc.gateway.poller.Stop(wc.desc)
		err = wc.gateway.poller.Start(wc.desc, wc.startEvent)
		if err != nil {
			wc.release(1)

			return err
		}
	}

	return nil
}

func (wc *websocketConn) release(_ int) {
	// delete the reference from the gateway's conns
	g := wc.gateway
	g.connsMtx.Lock()
	_, ok := g.conns[wc.connID]
	if !ok {
		g.connsMtx.Unlock()

		return
	}
	delete(g.conns, wc.connID)
	g.connsMtx.Unlock()

	// Decrease the total connection counter
	totalConns := atomic.AddInt32(&g.connsTotal, -1)

	if ce := g.cfg.Logger.Check(log.DebugLevel, "websocket connection removed"); ce != nil {
		ce.Write(
			zap.Uint64("ConnID", wc.connID),
			zap.Int32("Total", totalConns),
		)
	}

	wc.mtx.Lock()
	if wc.desc != nil {
		_ = wc.gateway.poller.Stop(wc.desc)
		err := wc.desc.Close()
		if err != nil {
			if ce := g.cfg.Logger.Check(log.DebugLevel, "got error on closing desc"); ce != nil {
				ce.Write(zap.Error(err))
			}
		}
	}
	_ = wc.conn.Close()

	if !wc.closed {
		g.delegate.OnClose(wc)
		wc.closed = true
		wc.conn = nil
	}

	wc.mtx.Unlock()
}

func (wc *websocketConn) startEvent(event netpoll.Event) {
	if atomic.LoadInt32(&wc.gateway.stop) == 1 {
		return
	}

	if event&netpoll.EventReadHup != 0 {
		wc.release(2)

		return
	}

	if event&netpoll.EventRead != 0 {
		atomic.StoreInt64(&wc.lastActivity, tools.CPUTicks())
		wc.gateway.waitGroupReaders.Add(1)

		err := goPoolNB.Submit(
			func() {
				waitGroup := pools.AcquireWaitGroup()
				err := wc.gateway.websocketReadPump(wc, waitGroup)
				if err != nil {
					wc.release(3)
				} else {
					_ = wc.gateway.poller.Resume(wc.desc)
				}
				waitGroup.Wait()
				pools.ReleaseWaitGroup(waitGroup)
				wc.gateway.waitGroupReaders.Done()
			},
		)
		if err != nil {
			wc.gateway.cfg.Logger.Warn("got error on start event go-routine pool", zap.Error(err))
		}
	}
}

func (wc *websocketConn) read(ms []wsutil.Message) ([]wsutil.Message, error) {
	var err error
	wc.mtx.Lock()
	if wc.conn != nil {
		_ = wc.conn.SetReadDeadline(time.Now().Add(defaultReadTimout))
		ms, err = wsutil.ReadMessage(wc.conn, ws.StateServerSide, ms)
	} else {
		err = ErrConnectionClosed
	}
	wc.mtx.Unlock()

	return ms, err
}

func (wc *websocketConn) write(opCode ws.OpCode, payload []byte) (err error) {
	wc.mtx.Lock()
	if wc.conn != nil {
		_ = wc.conn.SetWriteDeadline(time.Now().Add(defaultWriteTimeout))
		err = wsutil.WriteMessage(wc.conn, ws.StateServerSide, opCode, payload)
	} else {
		err = ErrWriteToClosedConn
	}
	wc.mtx.Unlock()

	return
}

func (wc *websocketConn) Get(key string) interface{} {
	wc.kvLock.Lock()
	v := wc.kv[key]
	wc.kvLock.Unlock()

	return v
}

func (wc *websocketConn) Set(key string, val interface{}) {
	wc.kvLock.Lock()
	wc.kv[key] = val
	wc.kvLock.Unlock()
}

func (wc *websocketConn) Walk(f func(k string, v interface{}) bool) {
	wc.kvLock.Lock()
	defer wc.kvLock.Unlock()

	for k, v := range wc.kv {
		if !f(k, v) {
			return
		}
	}
}

func (wc *websocketConn) ConnID() uint64 {
	return atomic.LoadUint64(&wc.connID)
}

func (wc *websocketConn) ClientIP() string {
	return string(wc.clientIP)
}

func (wc *websocketConn) SetClientIP(ip []byte) {
	wc.clientIP = append(wc.clientIP[:0], ip...)
}

// WriteBinary
// Make sure you don't use payload after calling this function, because its underlying
// array will be put back into the pool to be reused.
func (wc *websocketConn) WriteBinary(streamID int64, payload []byte) error {
	if wc == nil || wc.closed {
		return ErrWriteToClosedConn
	}
	wc.gateway.waitGroupWriters.Add(1)

	opCode := ws.OpBinary
	if wc.gateway.cfg.TextDataFrame {
		opCode = ws.OpText
	}
	wr := acquireWriteRequest(wc, opCode)
	wr.CopyPayload(payload)
	err := wc.gateway.websocketWritePump(wr)
	if err != nil {
		wc.release(4)
	}
	releaseWriteRequest(wr)
	metrics.IncCounter(metrics.CntGatewayOutgoingWebsocketMessage)

	return nil
}

func (wc *websocketConn) Disconnect() {
	wc.release(5)
}

func (wc *websocketConn) Persistent() bool {
	return true
}

// writeRequest
type writeRequest struct {
	wc      *websocketConn
	opCode  ws.OpCode
	payload []byte
}

func (wr *writeRequest) CopyPayload(p []byte) {
	wr.payload = append(wr.payload[:0], p...)
}

// websocketConnGC the garbage collector of the stalled websocket connections
type websocketConnGC struct {
	bg     *bigcache.BigCache
	gw     *Gateway
	inChan chan uint64
}

func newWebsocketConnGC(gw *Gateway) *websocketConnGC {
	gc := &websocketConnGC{
		gw:     gw,
		inChan: make(chan uint64, 1000),
	}
	bgConf := bigcache.DefaultConfig(time.Duration(gw.maxIdleTime))
	bgConf.CleanWindow = time.Second
	bgConf.Verbose = false
	bgConf.OnRemoveWithReason = gc.onRemove
	bgConf.Shards = 128
	bgConf.MaxEntrySize = 8
	bgConf.MaxEntriesInWindow = 100000
	gc.bg, _ = bigcache.NewBigCache(bgConf)

	// background job for receiving connIDs
	go func() {
		b := make([]byte, 8)
		for connID := range gc.inChan {
			binary.BigEndian.PutUint64(b, connID)
			_ = gc.bg.Set(tools.ByteToStr(b), b)
		}
	}()

	return gc
}

func (gc *websocketConnGC) onRemove(key string, entry []byte, reason bigcache.RemoveReason) {
	switch reason {
	case bigcache.Expired:
		connID := binary.BigEndian.Uint64(entry)
		if wsConn := gc.gw.getConnection(connID); wsConn != nil {
			if tools.CPUTicks()-atomic.LoadInt64(&wsConn.lastActivity) > gc.gw.maxIdleTime {
				wsConn.release(6)
			} else {
				gc.monitorConnection(connID)
			}
		}
	}
}

func (gc *websocketConnGC) monitorConnection(connID uint64) {
	gc.inChan <- connID
}
