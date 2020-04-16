// +build !windows,!appengine

package websocketGateway

import (
	"git.ronaksoftware.com/ronak/rony"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"github.com/gobwas/ws"
	"github.com/mailru/easygo/netpoll"
	"time"

	"go.uber.org/zap"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"syscall"
)

/*
   Creation Time: 2019 - Feb - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// Conn
type Conn struct {
	sync.Mutex
	ConnID   uint64
	AuthID   int64
	ClientIP string

	// Internals
	buf          chan *rony.MessageEnvelope
	gateway      *Gateway
	lastActivity int64
	conn         net.Conn
	desc         *netpoll.Desc
	closed       bool
	flushChan    chan bool
	flushing     int32
}

func (wc *Conn) GetAuthID() int64 {
	return wc.AuthID
}

func (wc *Conn) SetAuthID(authID int64) {
	wc.AuthID = authID
}

func (wc *Conn) GetConnID() uint64 {
	return wc.ConnID
}

func (wc *Conn) GetClientIP() string {
	return wc.ClientIP
}

func (wc *Conn) Push(m *rony.MessageEnvelope) {
	wc.buf <- m
}

func (wc *Conn) Pop() *rony.MessageEnvelope {
	return <-wc.buf
}

func (wc *Conn) startEvent(event netpoll.Event) {
	if atomic.LoadInt32(&wc.gateway.stop) == 1 {
		return
	}
	if event&netpoll.EventRead != 0 {
		wc.lastActivity = time.Now().Unix()
		// TODO:: rate limit this
		wc.gateway.waitGroupReaders.Add(1)
		wc.gateway.readPump(wc)
	}

}

// SendBinary
// Make sure you don't use payload after calling this function, because its underlying
// array will be put back into the pool to be reused.
// You MUST NOT re-use the underlying array of payload, otherwise you might get unexpected results.
func (wc *Conn) SendBinary(streamID int64, payload []byte) error {
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

func (wc *Conn) Disconnect() {
	wc.gateway.removeConnection(wc.ConnID)
}

func (wc *Conn) Persistent() bool {
	return true
}

// Flush
func (wc *Conn) Flush() {
	if atomic.CompareAndSwapInt32(&wc.flushing, 0, 1) {
		go wc.flushJob()
	}
	select {
	case wc.flushChan <- true:
	default:
	}
}

func (wc *Conn) flushJob() {
	timer := pools.AcquireTimer(flushDebounceTime)

	keepGoing := true
	cnt := 0
	for keepGoing {
		if cnt++; cnt > flushDebounceMax {
			break
		}
		select {
		case <-wc.flushChan:
			pools.ResetTimer(timer, flushDebounceTime)
		case <-timer.C:
			// Stop the loop
			keepGoing = false
		}
	}

	// Read the flush function
	bytesSlice := wc.gateway.FlushFunc(wc)

	// Reset the flushing flag, let the flusher run again
	atomic.StoreInt32(&wc.flushing, 0)

	err := wc.SendBinary(0, bytesSlice)
	if err != nil {
		if ce := log.Check(log.DebugLevel, "Error On Write To Websocket Conn"); ce != nil {
			ce.Write(
				zap.Uint64("ConnID", wc.ConnID),
				zap.Int64("authID", wc.AuthID),
				zap.Error(err),
			)
		}
	}
	pools.ReleaseTimer(timer)
}

// Check will check the underlying net.Conn if it is stalled or not. This code has been
// copied from mysql-go-driver, and it is not very efficient and only works on unix based
// hosts
func (wc *Conn) Check() error {
	var sysErr error

	sysConn, ok := wc.conn.(syscall.Conn)
	if !ok {
		return nil
	}
	rawConn, err := sysConn.SyscallConn()
	if err != nil {
		return err
	}

	err = rawConn.Read(func(fd uintptr) bool {
		var buf [1]byte
		n, err := syscall.Read(int(fd), buf[:])
		switch {
		case n == 0 && err == nil:
			sysErr = io.EOF
		case n > 0:
			sysErr = ErrUnexpectedSocketRead
		case err == syscall.EAGAIN || err == syscall.EWOULDBLOCK:
			sysErr = nil
		default:
			sysErr = err
		}
		return true
	})
	if err != nil {
		return err
	}

	return sysErr

}

type writeRequest struct {
	wc      *Conn
	opCode  ws.OpCode
	payload []byte
}
