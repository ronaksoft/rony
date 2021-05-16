package testEnv

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/internal/testEnv/pb/service"
	"github.com/ronaksoft/rony/store"
	"time"
)

/*
   Creation Time: 2020 - Jul - 17
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Handlers struct {
	ServerID string
}

func (h *Handlers) Set(ctx *edge.RequestCtx, req *service.SetRequest, res *service.SetResponse) {
	err := ctx.Store().Update(func(txn store.Txn) error {
		return nil
	})
	if err != nil {
		ctx.PushError(rony.ErrCodeInternal, err.Error())
		return
	}
	res.OK = true
	res.PushToContext(ctx)
}

func (h *Handlers) Get(ctx *edge.RequestCtx, req *service.GetRequest, res *service.GetResponse) {
	err := ctx.Store().View(func(txn store.Txn) error {
		return nil
	})
	if err != nil {
		ctx.PushError(rony.ErrCodeInternal, err.Error())
		return
	}
	res.PushToContext(ctx)
}

func (h *Handlers) EchoInternal(ctx *edge.RequestCtx, req *service.EchoRequest, res *service.EchoResponse) {
	panic("implement me")
}

func (h *Handlers) Echo(ctx *edge.RequestCtx, req *service.EchoRequest, res *service.EchoResponse) {
	res.ServerID = h.ServerID
	res.Timestamp = req.Timestamp
	res.Int = req.Int
	res.Responder = h.ServerID
}

func (h *Handlers) EchoDelay(ctx *edge.RequestCtx, req *service.EchoRequest, res *service.EchoResponse) {
	res.ServerID = h.ServerID
	res.Timestamp = req.Timestamp
	res.Int = req.Int
	res.Responder = h.ServerID
	time.Sleep(time.Second * 1)
}

func (h *Handlers) EchoTunnel(ctx *edge.RequestCtx, req *service.EchoRequest, res *service.EchoResponse) {
	res.ServerID = h.ServerID
	res.Timestamp = req.Timestamp
	res.Int = req.Int
	res.Responder = h.ServerID
}
