// +build !windows,!appengine

package websocketGateway

import (
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
	ConnID     uint64
	AuthID     int64
	UserID     int64
	AuthKey    []byte
	ClientIP   string
	Counters   ConnCounters
	ServerSeq  int64
	ClientSeq  int64
	ClientType string

	// Internals
	gateway      *Gateway
	lastActivity time.Time
	conn         net.Conn
	desc         *netpoll.Desc
	writeErrors  int32
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

func (wc *Conn) GetAuthKey() []byte {
	return wc.AuthKey
}

func (wc *Conn) SetAuthKey(authKey []byte) {
	wc.AuthKey = authKey
}

func (wc *Conn) IncServerSeq(x int64) int64 {
	return atomic.AddInt64(&wc.ServerSeq, x)
}

func (wc *Conn) GetConnID() uint64 {
	return wc.ConnID
}

func (wc *Conn) GetClientIP() string {
	return wc.ClientIP
}

func (wc *Conn) GetUserID() int64 {
	return wc.UserID
}

func (wc *Conn) SetUserID(userID int64) {
	wc.UserID = userID
}

type ConnCounters struct {
	IncomingCount uint64
	IncomingBytes uint64
	OutgoingCount uint64
	OutgoingBytes uint64
}

func (wc *Conn) startEvent(event netpoll.Event) {
	if atomic.LoadInt32(&wc.gateway.stop) == 1 {
		return
	}
	if event&netpoll.EventRead != 0 {
		wc.lastActivity = time.Now()
		select {
		case wc.gateway.connsInQ <- wc:
			if ce := log.Check(log.DebugLevel, "Websocket StartRead"); ce != nil {
				ce.Write(
					zap.Uint64("ConnID", wc.ConnID),
					zap.Int64("authID", wc.AuthID),
					zap.String("ClientType", wc.ClientType),
					zap.String("IP", wc.ClientIP),
				)
			}
		default:
			log.Info("Websocket Dropped, Busy Server")
			// wc.gateway.removeConnection(wc.ConnID)
		}
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

	select {
	case wc.gateway.connsOutQ <- writeRequest{
		wc:      wc,
		payload: payload,
		opCode:  ws.OpBinary,
	}:
	default:
		return ErrWriteToFullBufferedConn

	}
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

	for idx := 0; idx < len(bytesSlice); idx++ {
		err := wc.SendBinary(0, bytesSlice[idx])
		if err != nil {
			if ce := log.Check(log.DebugLevel, "Error On Write To Websocket Conn"); ce != nil {
				ce.Write(
					zap.Uint64("ConnID", wc.ConnID),
					zap.Int64("authID", wc.AuthID),
					zap.Error(err),
				)
			}
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
