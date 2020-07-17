package main

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/edge"
	"git.ronaksoftware.com/ronak/rony/gateway"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
)

type dispatcher struct{}

func (d dispatcher) OnUpdate(ctx *edge.DispatchCtx, authID int64, envelope *rony.UpdateEnvelope) {

}

func (d dispatcher) OnMessage(ctx *edge.DispatchCtx, authID int64, envelope *rony.MessageEnvelope) {
	if ctx.Conn() != nil {
		protoBytes := pools.Bytes.GetLen(envelope.Size())
		_, _ = envelope.MarshalToSizedBuffer(protoBytes)
		err := ctx.Conn().SendBinary(ctx.StreamID(), protoBytes)
		if err != nil {
			fmt.Println("Error On SendBinary", err)
		}
		pools.Bytes.Put(protoBytes)
	}

}

func (d dispatcher) Prepare(ctx *edge.DispatchCtx, data []byte, kvs ...gateway.KeyValue) (err error) {
	return ctx.UnmarshalEnvelope(data)
}

func (d dispatcher) Done(ctx *edge.DispatchCtx) {}
