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
	GetAuthID() int64
	GetUserID() int64
	GetAuthKey(buf []byte) []byte
	GetConnID() uint64
	GetClientIP() string
	Push(m *rony.MessageEnvelope)
	Pop() *rony.MessageEnvelope
	SendBinary(streamID int64, data []byte) error
	SetAuthID(authID int64)
	SetAuthKey(key []byte)
	SetUserID(userID int64)
	Persistent() bool
}

// Gateway defines the gateway interface where clients could connect
// and communicate with the edge server
type Gateway interface {
	Start()
	Run()
	Shutdown()
	Addr() []string
}

type ConnectHandler func(c Conn)
type MessageHandler func(c Conn, streamID int64, data []byte, kvs ...KeyValue)
type CloseHandler func(c Conn)
type KeyValue struct {
	Key   string
	Value string
}
