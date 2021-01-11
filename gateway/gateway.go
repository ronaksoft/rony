package gateway

import (
	"github.com/ronaksoft/rony"
)

/*
   Creation Time: 2019 - Aug - 31
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Protocol string

const (
	Undefined Protocol = ""
	TCP       Protocol = "tcp"
	Dummy     Protocol = "dummy"
	Websocket Protocol = "ws"
	Http      Protocol = "http"
)

// Gateway defines the gateway interface where clients could connect
// and communicate with the edge server
type Gateway interface {
	Start()
	Run()
	Shutdown()
	GetConn(connID uint64) rony.Conn
	Addr() []string
}

type ConnectHandler func(c rony.Conn, kvs ...*rony.KeyValue)
type MessageHandler func(c rony.Conn, streamID int64, data []byte)
type CloseHandler func(c rony.Conn)
