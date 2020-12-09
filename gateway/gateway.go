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
)

// Conn defines the Connection interface
type Conn interface {
	ConnID() uint64
	ClientIP() string
	Push(m *rony.MessageEnvelope)
	Pop() *rony.MessageEnvelope
	SendBinary(streamID int64, data []byte) error
	Persistent() bool
	Get(key string) interface{}
	Set(key string, val interface{})
}

// Gateway defines the gateway interface where clients could connect
// and communicate with the edge server
type Gateway interface {
	Start()
	Run()
	Shutdown()
	GetConn(connID uint64) Conn
	Addr() []string
}

type ConnectHandler func(c Conn)
type MessageHandler func(c Conn, streamID int64, data []byte, kvs ...KeyValue)
type CloseHandler func(c Conn)
type KeyValue struct {
	Key   string
	Value string
}
