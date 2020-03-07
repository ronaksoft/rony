package main

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/cmd/cli-playground/msg"
	"git.ronaksoftware.com/ronak/rony/gateway"
	"github.com/gobwas/pool/pbytes"
)

type dispatcher struct{}

func (d dispatcher) DispatchUpdate(conn gateway.Conn, streamID, authID int64, envelope *rony.UpdateEnvelope) {

}

func (d dispatcher) DispatchMessage(conn gateway.Conn, streamID, authID int64, envelope *rony.MessageEnvelope) {
	proto := &msg.ProtoMessage{}
	proto.AuthID = authID
	proto.Payload, _ = envelope.Marshal()
	protoBytes := pbytes.GetLen(proto.Size())
	_, _ = proto.MarshalTo(protoBytes)
	if conn != nil {
		err := conn.SendBinary(streamID, protoBytes)
		if err != nil {
			fmt.Println("Error On SendBinary", err)
		}
	}

}

func (d dispatcher) DispatchRequest(conn gateway.Conn, streamID int64, data []byte, envelope *rony.MessageEnvelope) (authID int64, err error) {
	proto := &msg.ProtoMessage{}
	err = proto.Unmarshal(data)
	if err != nil {
		return
	}
	err = envelope.Unmarshal(proto.Payload)
	if err != nil {
		return
	}
	authID = proto.AuthID
	return
}
