package gateway

import (
	"github.com/gogo/protobuf/proto"
)

/*
   Creation Time: 2019 - Aug - 31
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type Protocol string

const (
	Undefined Protocol = ""
	Websocket Protocol = "websocket"
	QUIC      Protocol = "quic"
	HTTP      Protocol = "http"
	GRPC      Protocol = "grpc"
)

type Client string

const (
	User Client = "user"
	Bot  Client = "bot"
)

type Conn interface {
	GetConnID() uint64
	GetClientIP() string
	SendProto(streamID int64, pm ProtoBufferMessage) error
	GetAuthID() int64
	SetAuthID(int64)
	GetAuthKey() []byte
	SetAuthKey([]byte)
	GetUserID() int64
	SetUserID(int64)
	IncServerSeq(int64) int64
	Flush()
	Persistent() bool
}

type ConnectHandler func(connID uint64)
type MessageHandler func(c Conn, streamID int64, date []byte)
type CloseHandler func(c Conn)
type FailedWriteHandler func(c Conn, data []byte, err error)
type SuccessWriteHandler func(c Conn)
type FlushFunc func(c Conn) [][]byte

type ProtoBufferMessage interface {
	proto.Marshaler
	proto.Sizer
	proto.Unmarshaler
	MarshalTo([]byte) (int, error)
}
