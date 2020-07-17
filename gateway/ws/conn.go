// +build !windows,!appengine

package websocketGateway

import (
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/pools"
	"git.ronaksoftware.com/ronak/rony/tools"
	"github.com/gobwas/ws"
	"github.com/mailru/easygo/netpoll"
	"time"

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
	buf          *tools.LinkedList
	gateway      *Gateway
	lastActivity int64
	conn         net.Conn
	desc         *netpoll.Desc
	closed       bool
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
	wc.buf.Append(m)
}

func (wc *Conn) Pop() *rony.MessageEnvelope {
	v := wc.buf.PickHeadData()
	if v != nil {
		return v.(*rony.MessageEnvelope)
	}
	return nil
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
