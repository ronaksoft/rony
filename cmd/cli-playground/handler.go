package main

import (
	"fmt"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/internal/testEnv/pb"
	"github.com/ronaksoft/rony/tools"
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

func (h *SampleServer) Func1(ctx *edge.RequestCtx, req *pb.Req1, res *pb.Res1) {
	res.Item1 = req.Item1
}

func (h *SampleServer) Func2(ctx *edge.RequestCtx, req *pb.Req2, res *pb.Res2) {
	res.Item1 = req.Item1
}

func (h *SampleServer) Echo(ctx *edge.RequestCtx, req *pb.EchoRequest, res *pb.EchoResponse) {
	res.Bool = req.Bool
	res.Int = req.Int
	res.Timestamp = tools.NanoTime()
	res.Delay = res.Timestamp - req.Timestamp
}

func (h *SampleServer) Ask(ctx *edge.RequestCtx, req *pb.AskRequest, res *pb.AskResponse) {
	res.Responder = h.es.GetServerID()
	res.Coordinator = ctx.ServerID()

	if ctx.Kind() == edge.ClusterMessage {
		fmt.Printf("%s :: Cluster ASK: %#v\n", h.es.GetServerID(), res)
	} else {
		ctx.PushClusterMessage(req.ServerID, ctx.ReqID(), pb.C_Ask, req)
	}
}
