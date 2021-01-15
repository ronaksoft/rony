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
	// log.Info("OnMessage", zap.String("ServerID", ctx.ServerID()), zap.String("Kind", ctx.Kind().String()))
	switch ctx.Kind() {
	case edge.GatewayMessage, edge.TunnelMessage:
		mo := proto.MarshalOptions{UseCachedSize: true}

		buf := pools.Buffer.GetCap(mo.Size(envelope))
		b, _ := mo.MarshalAppend(*buf.Bytes(), envelope)
		buf.SetBytes(&b)
		err := ctx.Conn().SendBinary(ctx.StreamID(), *buf.Bytes())
		if err != nil {
			fmt.Println("Error On SendBinary", err)
		}
		pools.Buffer.Put(buf)
	}
}

func (d dispatcher) Interceptor(ctx *edge.DispatchCtx, data []byte) (err error) {
	return ctx.UnmarshalEnvelope(data)
}

func (d dispatcher) Done(ctx *edge.DispatchCtx) {}
