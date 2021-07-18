package rpc

import (
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/errors"
	"go.uber.org/zap"
)

/*
   Creation Time: 2021 - Jul - 04
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

//go:generate protoc -I=. -I=../../.. --go_out=paths=source_relative:. sample.proto
//go:generate protoc -I=. -I=../../.. --gorony_out=paths=source_relative:. sample.proto
func init() {}

// Sample implements auto-generated service.ISample interface
type Sample struct{}

func (s *Sample) InfoWithClientRedirect(ctx *edge.RequestCtx, req *InfoRequest, res *InfoResponse) {
	ctx.Log().Warn("Received (Client)", zap.Uint64("ReqRS", req.GetReplicaSet()), zap.Uint64("ServerRS", ctx.ReplicaSet()))
	if req.GetReplicaSet() != ctx.ReplicaSet() {
		ctx.PushRedirectRequest(req.GetReplicaSet())
		return
	}
	res.ServerID = ctx.ServerID()
	res.RandomText = req.GetRandomText()
}

func (s *Sample) InfoWithServerRedirect(ctx *edge.RequestCtx, req *InfoRequest, res *InfoResponse) {
	ctx.Log().Warn("Received", zap.Uint64("ReplicaSet", req.GetReplicaSet()))
	if req.GetReplicaSet() != ctx.ReplicaSet() {
		err := TunnelRequestSampleInfoWithServerRedirect(ctx, req.GetReplicaSet(), req, res)
		if err != nil {
			ctx.Log().Warn("Got Error", zap.Error(err))
			ctx.PushError(errors.ErrInternalServer)
			return
		}
		return
	}
	res.ServerID = ctx.ServerID()
	res.RandomText = req.GetRandomText()
}
