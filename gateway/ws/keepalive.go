package websocketGateway

import (
	"fmt"
	"net"
	"os"
	"syscall"
	"time"
)

/*
   Creation Time: 2019 - Nov - 21
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

// Enable TCP keepalive in non-blocking mode with given settings for
// the connection, which must be a *tcp.TCPConn.
func SetKeepAlive(c net.Conn, idleTime time.Duration, count int, interval time.Duration) (err error) {

	conn, ok := c.(*net.TCPConn)
	if !ok {
		return fmt.Errorf("bad connection type: %T", c)
	}

	if err := conn.SetKeepAlive(true); err != nil {
		return err
	}

	var f *os.File

	if f, err = conn.File(); err != nil {
		return err
	}
	defer f.Close()

	fd := int(f.Fd())

	if err = setIdle(fd, secs(idleTime)); err != nil {
		return err
	}

	if err = setCount(fd, count); err != nil {
		return err
	}

	if err = setInterval(fd, secs(interval)); err != nil {
		return err
	}

	if err = setNonblock(fd); err != nil {
		return err
	}

	return nil
}

func secs(d time.Duration) int {
	d += (time.Second - time.Nanosecond)
	return int(d.Seconds())
}

func setNonblock(fd int) error {
	return os.NewSyscallError("setsockopt", syscall.SetNonblock(fd, true))

}
