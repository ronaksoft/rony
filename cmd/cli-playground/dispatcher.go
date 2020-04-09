package main

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/edge"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv/pb"
)

type dispatcher struct{}

func (d dispatcher) DispatchUpdate(ctx *edge.DispatchCtx, authID int64, envelope *rony.UpdateEnvelope) {

}

func (d dispatcher) DispatchMessage(ctx *edge.DispatchCtx, authID int64, envelope *rony.MessageEnvelope) {
	proto := &pb.ProtoMessage{}
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

func (d dispatcher) DispatchRequest(ctx *edge.DispatchCtx, data []byte) (err error) {
	proto := &pb.ProtoMessage{}
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
