package tcpGateway

import (
	"bytes"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/mailru/easygo/netpoll"
)

/*
   Creation Time: 2020 - Aug - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	wrapConnPool sync.Pool
)

// filer describes an object that has ability to return os.File.
type filer interface {
	// File returns a copy of object's file descriptor.
	File() (*os.File, error)
}

type wrapConn struct {
	c   net.Conn
	r   io.Reader
	buf *bytes.Buffer
}

func (wc *wrapConn) Write(b []byte) (n int, err error) {
	return wc.c.Write(b)
}

func (wc *wrapConn) Close() error {
	err := wc.c.Close()
	releaseWrapConn(wc)

	return err
}

func (wc *wrapConn) LocalAddr() net.Addr {
	return wc.c.LocalAddr()
}

func (wc *wrapConn) RemoteAddr() net.Addr {
	return wc.c.RemoteAddr()
}

func (wc *wrapConn) SetDeadline(t time.Time) error {
	return wc.c.SetDeadline(t)
}

func (wc *wrapConn) SetReadDeadline(t time.Time) error {
	return wc.c.SetReadDeadline(t)
}

func (wc *wrapConn) SetWriteDeadline(t time.Time) error {
	return wc.c.SetWriteDeadline(t)
}

func newWrapConn(c net.Conn) *wrapConn {
	wc := &wrapConn{
		c:   c,
		buf: bytes.NewBuffer(make([]byte, 0, 128)),
	}
	wc.r = io.TeeReader(wc.c, wc.buf)

	return wc
}

func acquireWrapConn(c net.Conn) *wrapConn {
	wc, ok := wrapConnPool.Get().(*wrapConn)
	if !ok {
		return newWrapConn(c)
	}
	wc.c = c
	wc.r = io.TeeReader(wc.c, wc.buf)

	return wc
}

func releaseWrapConn(wc *wrapConn) {
	wc.buf.Reset()
	wrapConnPool.Put(wc)
}

func (wc *wrapConn) UnsafeConn() net.Conn {
	return wc.c
}

func (wc *wrapConn) Read(p []byte) (int, error) {
	return wc.r.Read(p)
}

func (wc *wrapConn) File() (*os.File, error) {
	x, ok := wc.c.(filer)
	if !ok {
		return nil, netpoll.ErrNotFiler
	}

	return x.File()
}

func (wc *wrapConn) ReadyForUpgrade() {
	wc.r = io.MultiReader(wc.buf, wc.c)
}
