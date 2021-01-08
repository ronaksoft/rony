package main

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/internal/testEnv/pb"
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

func (h *SampleServer) EchoLeaderOnly(ctx *edge.RequestCtx, req *pb.EchoRequest, res *pb.EchoResponse) {
	res.Int = req.Int
	res.Timestamp = tools.NanoTime()
	res.Delay = res.Timestamp - req.Timestamp
	res.ServerID = h.es.GetServerID()
}

func (h *SampleServer) EchoTunnel(ctx *edge.RequestCtx, req *pb.EchoRequest, res *pb.EchoResponse) {
	res.Int = req.Int
	res.Timestamp = tools.NanoTime()
	res.Delay = res.Timestamp - req.Timestamp
	res.ServerID = h.es.GetServerID()

	switch ctx.Kind() {
	case edge.GatewayMessage:
		out := rony.PoolMessageEnvelope.Get()
		defer rony.PoolMessageEnvelope.Put(out)
		in := rony.PoolMessageEnvelope.Get()
		defer rony.PoolMessageEnvelope.Put(in)
		out.Fill(ctx.ReqID(), pb.C_EchoTunnel, req)
		err := ctx.ExecuteRemote(req.ReplicaSet, true, out, in)
		if err != nil {
			ctx.PushError(rony.ErrCodeInternal, err.Error())
			return
		}
		switch in.Constructor {
		case pb.C_EchoResponse:
			_ = res.Unmarshal(in.Message)
		default:
			ctx.PushError(rony.ErrCodeInternal, "invalid constructor in response")
		}
	default:
		return

	}

}

func (h *SampleServer) Echo(ctx *edge.RequestCtx, req *pb.EchoRequest, res *pb.EchoResponse) {
	res.Int = req.Int
	res.Timestamp = tools.NanoTime()
	res.Delay = res.Timestamp - req.Timestamp
	res.ServerID = h.es.GetServerID()
}

func (h *SampleServer) EchoDelay(ctx *edge.RequestCtx, req *pb.EchoRequest, res *pb.EchoResponse) {
	res.Int = req.Int
	res.Timestamp = tools.NanoTime()
	res.Delay = res.Timestamp - req.Timestamp
	res.ServerID = h.es.GetServerID()
	time.Sleep(time.Second * 2)
}
