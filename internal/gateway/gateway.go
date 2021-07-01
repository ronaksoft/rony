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

// Gateway defines the gateway interface where clients could connect
// and communicate with the edge server
type Gateway interface {
	Start()
	Run()
	Shutdown()
	GetConn(connID uint64) rony.Conn
	Addr() []string
	Protocol() Protocol
}

type RequestCtx interface {

}

type (
	ConnectHandler = func(c rony.Conn, kvs ...*rony.KeyValue)
	MessageHandler = func(c rony.Conn, streamID int64, data []byte, bypass bool)
	CloseHandler   = func(c rony.Conn)
)

