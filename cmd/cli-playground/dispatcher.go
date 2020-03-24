package main

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/cmd/cli-playground/msg"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
)

type dispatcher struct{}

func (d dispatcher) DispatchUpdate(ctx *rony.DispatchCtx, authID int64, envelope *rony.UpdateEnvelope) {

}

func (d dispatcher) DispatchMessage(ctx *rony.DispatchCtx, authID int64, envelope *rony.MessageEnvelope) {
	proto := &msg.ProtoMessage{}
	proto.AuthID = authID
	proto.Payload, _ = envelope.Marshal()
	protoBytes := pools.Bytes.GetLen(proto.Size())
	_, _ = proto.MarshalTo(protoBytes)
	if ctx.Conn() != nil {
		err := ctx.Conn().SendBinary(ctx.StreamID(), protoBytes)
		if err != nil {
			fmt.Println("Error On SendBinary", err)
		}
	}

}

func (d dispatcher) DispatchRequest(ctx *rony.DispatchCtx,  data []byte) (err error) {
	proto := &msg.ProtoMessage{}
	err = proto.Unmarshal(data)
	if err != nil {
		return
	}
	err = ctx.UnmarshalEnvelope(proto.Payload)
	if err != nil {
		return
	}
	ctx.SetAuthID(proto.AuthID)
	return
}
