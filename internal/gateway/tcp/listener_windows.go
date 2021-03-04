// +build windows

package tcpGateway

import (
	"net"
)

/*
   Creation Time: 2021 - Mar - 04
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type wrapListener struct {
	l net.Listener
}

func (w *wrapListener) Accept() (net.Conn, error) {
	c, err := w.l.Accept()
	if err != nil {
		return nil, err
	}
	return acquireWrapConn(c), nil
}

func (w *wrapListener) Close() error {
	return w.l.Close()
}

func (w *wrapListener) Addr() net.Addr {
	return w.l.Addr()
}

func newWrapListener(listenOn string) (wl *wrapListener, err error) {
	wl = &wrapListener{}
	wl.l, err = net.Listen("tcp4", listenOn)
	return
}
