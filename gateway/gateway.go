package gateway

import (
	"git.ronaksoft.com/ronak/rony"
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
	Websocket Protocol = "websocket"
	HTTP      Protocol = "http"
	TCP       Protocol = "tcp"
)

type Conn interface {
	GetAuthID() int64
	GetConnID() uint64
	GetClientIP() string
	Push(m *rony.MessageEnvelope)
	Pop() *rony.MessageEnvelope
	SendBinary(streamID int64, data []byte) error
	SetAuthID(int64)
	Persistent() bool
}

type Gateway interface {
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
