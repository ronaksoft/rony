// +build !windows,!appengine

package tcp

import (
	"encoding/binary"
	"github.com/allegro/bigcache/v2"
	"github.com/gobwas/ws"
	"github.com/mailru/easygo/netpoll"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/tools"
	"sync"
	"time"

	"net"
	"sync/atomic"
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
	sync.Mutex
	connID   uint64
	clientIP []byte

	// KV Store
	mtx sync.RWMutex
	kv  map[string]interface{}

	// Internals
	buf          *tools.LinkedList
	gateway      *Gateway
	lastActivity int64
	conn         net.Conn
	desc         *netpoll.Desc
	closed       bool
}

func (wc *websocketConn) Get(key string) interface{} {
	wc.mtx.RLock()
	v := wc.kv[key]
	wc.mtx.RUnlock()
	return v
}

func (wc *websocketConn) Set(key string, val interface{}) {
	wc.mtx.Lock()
	wc.kv[key] = val
	wc.mtx.Unlock()
}

func (wc *websocketConn) ConnID() uint64 {
	return atomic.LoadUint64(&wc.connID)
}

func (wc *websocketConn) ClientIP() string {
	return net.IP(wc.clientIP).String()
}

func (wc *websocketConn) SetClientIP(ip []byte) {
	wc.clientIP = append(wc.clientIP[:0], ip...)
}

func (wc *websocketConn) Push(m *rony.MessageEnvelope) {
	wc.buf.Append(m)
}

func (wc *websocketConn) Pop() *rony.MessageEnvelope {
	v := wc.buf.PickHeadData()
	if v != nil {
		return v.(*rony.MessageEnvelope)
	}
	return nil
}

func (wc *websocketConn) startEvent(event netpoll.Event) {
	if atomic.LoadInt32(&wc.gateway.stop) == 1 {
		return
	}
	if event&netpoll.EventRead != 0 {
		wc.lastActivity = tools.TimeUnix()
		wc.gateway.waitGroupReaders.Add(1)

		ms := acquireWebsocketMessage()
		wc.gateway.readPump(wc, ms)
		releaseWebsocketMessage(ms)
	}
}

// SendBinary
// Make sure you don't use payload after calling this function, because its underlying
// array will be put back into the pool to be reused.
func (wc *websocketConn) SendBinary(streamID int64, payload []byte) error {
	if wc.closed {
		return ErrWriteToClosedConn
	}
	wc.gateway.waitGroupWriters.Add(1)

	wr := acquireWriteRequest(wc, ws.OpBinary)
	wr.CopyPayload(payload)
	wc.gateway.writePump(wr)
	releaseWriteRequest(wr)
	return nil
}

func (wc *websocketConn) Disconnect() {
	wc.gateway.removeConnection(wc.connID)
}

func (wc *websocketConn) Persistent() bool {
	return true
}

type writeRequest struct {
	wc      *websocketConn
	opCode  ws.OpCode
	payload []byte
}

func (wr *writeRequest) CopyPayload(p []byte) {
	wr.payload = append(wr.payload[:0], p...)
}

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
	bgConf := bigcache.DefaultConfig(time.Duration(gw.maxIdleTime) * time.Second)
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
			if tools.TimeUnix()-wsConn.lastActivity > gc.gw.maxIdleTime {
				gc.gw.removeConnection(connID)
			} else {
				gc.monitorConnection(connID)
			}
		}
	}
}

func (gc *websocketConnGC) monitorConnection(connID uint64) {
	gc.inChan <- connID
}
