package main

import (
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/errors"
	"github.com/ronaksoft/rony/internal/testEnv/pb/service"
	"github.com/ronaksoft/rony/store"
	"github.com/ronaksoft/rony/tools"
	"time"
)

/*
   Creation Time: 2020 - Feb - 24
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type SampleServer struct {
	es *edge.Server
}

func (h *SampleServer) Set(ctx *edge.RequestCtx, req *service.SetRequest, res *service.SetResponse) {
	err := ctx.Store().Update(func(txn store.Txn) error {
		alloc := tools.NewAllocator()
		defer alloc.ReleaseAll()
		return txn.Set(alloc, req.Value, req.Key)
	})
	if err != nil {
		ctx.PushError(errors.ErrInternal(err.Error(), nil))
		return
	}
	res.OK = true
	res.PushToContext(ctx)
}

func (h *SampleServer) Get(ctx *edge.RequestCtx, req *service.GetRequest, res *service.GetResponse) {
	err := ctx.Store().View(func(txn store.Txn) error {
		alloc := tools.NewAllocator()
		defer alloc.ReleaseAll()
		v, err := txn.Get(alloc, req.Key)
		if err != nil {
			return err
		}
		res.Value = append(res.Value, v...)
		return nil
	})
	if err != nil {
		ctx.PushError(errors.ErrInternal(err.Error(), nil))
		return
	}
	res.PushToContext(ctx)
}

func (h *SampleServer) EchoInternal(ctx *edge.RequestCtx, req *service.EchoRequest, res *service.EchoResponse) {
	panic("implement me")
}

func (h *SampleServer) EchoLeaderOnly(ctx *edge.RequestCtx, req *service.EchoRequest, res *service.EchoResponse) {
	res.Int = req.Int
	res.Timestamp = tools.NanoTime()
	res.Delay = res.Timestamp - req.Timestamp
	res.ServerID = h.es.GetServerID()
}

func (h *SampleServer) EchoTunnel(ctx *edge.RequestCtx, req *service.EchoRequest, res *service.EchoResponse) {
	res.Int = req.Int
	res.Timestamp = tools.NanoTime()
	res.Delay = res.Timestamp - req.Timestamp
	res.ServerID = h.es.GetServerID()

	switch ctx.Kind() {
	case edge.GatewayMessage:
		err := service.ExecuteRemoteSampleEchoTunnel(ctx, req.ReplicaSet, req, res)
		if err != nil {
			ctx.PushError(errors.ErrInternal(err.Error(), nil))
			return
		}
	default:
		return

	}

}

func (h *SampleServer) Echo(ctx *edge.RequestCtx, req *service.EchoRequest, res *service.EchoResponse) {
	res.Int = req.Int
	res.Timestamp = tools.NanoTime()
	res.Delay = res.Timestamp - req.Timestamp
	res.ServerID = h.es.GetServerID()
}

func (h *SampleServer) EchoDelay(ctx *edge.RequestCtx, req *service.EchoRequest, res *service.EchoResponse) {
	res.Int = req.Int
	res.Timestamp = tools.NanoTime()
	res.Delay = res.Timestamp - req.Timestamp
	res.ServerID = h.es.GetServerID()
	time.Sleep(time.Second * 2)
}
