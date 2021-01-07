package main

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/pools"
	"google.golang.org/protobuf/proto"
)

type dispatcher struct{}

func (d dispatcher) OnOpen(conn rony.Conn, kvs ...*rony.KeyValue) {

}

func (d dispatcher) OnClose(conn rony.Conn) {

}

func (d dispatcher) OnMessage(ctx *edge.DispatchCtx, envelope *rony.MessageEnvelope) {
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

func (d dispatcher) Interceptor(ctx *edge.DispatchCtx, data []byte) (err error) {
	return ctx.UnmarshalEnvelope(data)
}

func (d dispatcher) Done(ctx *edge.DispatchCtx) {}
