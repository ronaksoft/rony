package main

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/gateway"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/msg"
	"github.com/gobwas/pool/pbytes"
)

type dispatcher struct{}

func (d dispatcher) DispatchUpdate(conn gateway.Conn, streamID, authID int64, envelope *msg.UpdateEnvelope) {

}

func (d dispatcher) DispatchMessage(conn gateway.Conn, streamID, authID int64, envelope *msg.MessageEnvelope) {
	proto := pools.AcquireProtoMessage()
	proto.AuthID = authID
	proto.Payload, _ = envelope.Marshal()
	protoBytes := pbytes.GetLen(proto.Size())
	_, _ = proto.MarshalTo(protoBytes)
	err := conn.SendBinary(0, protoBytes)
	if err != nil {
		fmt.Println("Error On SendBinary", err)
	}
	pools.ReleaseProtoMessage(proto)
}

func (d dispatcher) DispatchRequest(conn gateway.Conn, streamID int64, data []byte, envelope *msg.MessageEnvelope) (authID int64, err error) {
	proto := pools.AcquireProtoMessage()
	err = proto.Unmarshal(data)
	if err != nil {
		return
	}
	err = envelope.Unmarshal(proto.Payload)
	if err != nil {
		return
	}
	authID = proto.AuthID
	pools.ReleaseProtoMessage(proto)
	return
}
