// +build !windows,!appengine

package tcpGateway

import (
	"encoding/binary"
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"github.com/allegro/bigcache/v2"
	"github.com/gobwas/ws"
	"github.com/mailru/easygo/netpoll"
	"time"

	"net"
	"sync"
	"sync/atomic"
)

/*
   Creation Time: 2019 - Feb - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// WebsocketConn
type WebsocketConn struct {
	sync.Mutex
	ConnID   uint64
	AuthID   int64
	ClientIP string

	// Internals
	buf          *tools.LinkedList
	gateway      *Gateway
	lastActivity int64
	conn         net.Conn
	desc         *netpoll.Desc
	closed       bool
}

func (wc *WebsocketConn) GetAuthID() int64 {
	return wc.AuthID
}

func (wc *WebsocketConn) SetAuthID(authID int64) {
	wc.AuthID = authID
}

func (wc *WebsocketConn) GetConnID() uint64 {
	return wc.ConnID
}

func (wc *WebsocketConn) GetClientIP() string {
	return net.IP(wc.ClientIP).String()
}

func (wc *WebsocketConn) Push(m *rony.MessageEnvelope) {
	wc.buf.Append(m)
}

func (wc *WebsocketConn) Pop() *rony.MessageEnvelope {
	v := wc.buf.PickHeadData()
	if v != nil {
		return v.(*rony.MessageEnvelope)
	}
	return nil
}

func (wc *WebsocketConn) startEvent(event netpoll.Event) {
	if atomic.LoadInt32(&wc.gateway.stop) == 1 {
		return
	}
	if event&netpoll.EventRead != 0 {
		wc.lastActivity = tools.TimeUnix()
		wc.gateway.waitGroupReaders.Add(1)
		wc.gateway.readPump(wc)
	}

}

// SendBinary
// Make sure you don't use payload after calling this function, because its underlying
// array will be put back into the pool to be reused.
// You MUST NOT re-use the underlying array of payload, otherwise you might get unexpected results.
func (wc *WebsocketConn) SendBinary(streamID int64, payload []byte) error {
	if wc.closed {
		return ErrWriteToClosedConn
	}
	wc.gateway.waitGroupWriters.Add(1)
	p := pools.Bytes.GetLen(len(payload))
	copy(p, payload)
	wc.gateway.writePump(writeRequest{
		wc:      wc,
		payload: p,
		opCode:  ws.OpBinary,
	})
	return nil
}

func (wc *WebsocketConn) Disconnect() {
	wc.gateway.removeConnection(wc.ConnID)
}

func (wc *WebsocketConn) Persistent() bool {
	return true
}

type writeRequest struct {
	wc      *WebsocketConn
	opCode  ws.OpCode
	payload []byte
}

type websocketConnGC struct {
	bg      *bigcache.BigCache
	gw      *Gateway
	inChan  chan uint64
	timeNow int64
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

	go func() {
		b := make([]byte, 8)
		for connID := range gc.inChan {
			binary.BigEndian.PutUint64(b, connID)
			_ = gc.bg.Set(tools.ByteToStr(b), b)
		}
	}()

	go func() {
		for {
			gc.timeNow = time.Now().Unix()
			time.Sleep(time.Second)
		}
	}()
	return gc
}

func (gc *websocketConnGC) onRemove(key string, entry []byte, reason bigcache.RemoveReason) {
	switch reason {
	case bigcache.Expired:
		connID := binary.BigEndian.Uint64(entry)
		if wsConn := gc.gw.getConnection(connID); wsConn != nil {
			if gc.timeNow-wsConn.lastActivity > gc.gw.maxIdleTime {
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
