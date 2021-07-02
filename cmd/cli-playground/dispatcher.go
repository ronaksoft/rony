package main

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/pools"
	"go.uber.org/zap"
)

type dispatcher struct{}

func (d dispatcher) OnOpen(conn rony.Conn, kvs ...*rony.KeyValue) {

}

func (d dispatcher) OnClose(conn rony.Conn) {

}

func (d dispatcher) OnMessage(ctx *edge.DispatchCtx, envelope *rony.MessageEnvelope) {
	log.Info("OnMessage", zap.String("ServerID", ctx.ServerID()), zap.String("Kind", ctx.Kind().String()))

}

func (d dispatcher) Interceptor(ctx *edge.DispatchCtx, data []byte) (err error) {
	return ctx.UnmarshalEnvelope(data)
}

func (d dispatcher) Done(ctx *edge.DispatchCtx) {
	log.Info("Done", zap.String("ServerID", ctx.ServerID()), zap.String("Kind", ctx.Kind().String()))
	for {
		envelope := ctx.BufferPop()
		if envelope == nil {
			break
		}
		switch ctx.Kind() {
		case edge.GatewayMessage, edge.TunnelMessage:
			buf := pools.Buffer.FromProto(envelope)
			err := ctx.Conn().WriteBinary(ctx.StreamID(), *buf.Bytes())
			if err != nil {
				fmt.Println("Error On WriteBinary", err)
			}
			pools.Buffer.Put(buf)
		}
	}

}
