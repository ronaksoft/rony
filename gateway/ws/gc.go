package websocketGateway

import (
	"encoding/binary"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"github.com/allegro/bigcache/v2"
	"time"
)

/*
   Creation Time: 2020 - Jan - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type connGC struct {
	bg     *bigcache.BigCache
	gw     *Gateway
	inChan chan uint64
}

func newGC(gw *Gateway) *connGC {
	gc := &connGC{
		gw:     gw,
		inChan: make(chan uint64, 1000),
	}
	bgConf := bigcache.DefaultConfig(gw.maxIdleTime)
	bgConf.CleanWindow = time.Second
	bgConf.Verbose = false
	bgConf.OnRemoveWithReason = gc.onRemove
	gc.bg, _ = bigcache.NewBigCache(bgConf)

	go func() {
		b := make([]byte, 8)
		for connID := range gc.inChan {
			binary.BigEndian.PutUint64(b, connID)
			_ = gc.bg.Set(tools.ByteToStr(b), b)
		}
	}()
	return gc
}

func (gc *connGC) onRemove(key string, entry []byte, reason bigcache.RemoveReason) {
	switch reason {
	case bigcache.Expired:
		connID := binary.BigEndian.Uint64(entry)
		if wsConn := gc.gw.GetConnection(connID); wsConn != nil {
			if time.Now().Sub(wsConn.lastActivity) > gc.gw.maxIdleTime {
				gc.gw.removeConnection(connID)
			} else {
				gc.monitorConnection(connID)
			}
		}
	}
}

func (gc *connGC) monitorConnection(connID uint64) {
	gc.inChan <- connID
}
