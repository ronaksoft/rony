package main

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/edge"
	"git.ronaksoftware.com/ronak/rony/gateway"
	"git.ronaksoftware.com/ronak/rony/pools"
	"google.golang.org/protobuf/proto"
)

type dispatcher struct{}

func (d dispatcher) OnOpen(conn gateway.Conn) {

}

func (d dispatcher) OnClose(conn gateway.Conn) {

}

func (d dispatcher) OnMessage(ctx *edge.DispatchCtx, authID int64, envelope *rony.MessageEnvelope) {
	if ctx.Conn() != nil {
		mo := proto.MarshalOptions{
			UseCachedSize: true,
		}
		protoBytes := pools.Bytes.GetCap(mo.Size(envelope))
		protoBytes, _ = mo.MarshalAppend(protoBytes, envelope)
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
