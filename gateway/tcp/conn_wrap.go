package tcp

import (
	"bytes"
	"github.com/mailru/easygo/netpoll"
	"io"
	"net"
	"os"
	"sync"
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
	net.Conn
	io.Reader
	buf *bytes.Buffer
}

func newWrapConn(c net.Conn) *wrapConn {
	wc := &wrapConn{
		Conn: c,
		buf:  bytes.NewBuffer(make([]byte, 0, 128)),
	}
	wc.Reader = io.TeeReader(wc.Conn, wc.buf)
	return wc
}

func acquireWrapConn(c net.Conn) *wrapConn {
	wc, ok := wrapConnPool.Get().(*wrapConn)
	if !ok {
		return newWrapConn(c)
	}
	wc.Conn = c
	wc.Reader = io.TeeReader(wc.Conn, wc.buf)
	return wc
}

func releaseWrapConn(wc *wrapConn) {
	wc.buf.Reset()
	wrapConnPool.Put(wc)
}

func (wc *wrapConn) Read(p []byte) (int, error) {
	return wc.Reader.Read(p)
}

func (wc *wrapConn) File() (*os.File, error) {
	x, ok := wc.Conn.(filer)
	if !ok {
		return nil, netpoll.ErrNotFiler
	}
	return x.File()
}

func (wc *wrapConn) ReadyForUpgrade() {
	wc.Reader = io.MultiReader(wc.buf, wc.Conn)
}
