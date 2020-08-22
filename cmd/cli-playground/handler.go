package main

import (
	"git.ronaksoft.com/ronak/rony/edge"
	"git.ronaksoft.com/ronak/rony/internal/testEnv/pb"
	"time"
)

/*
   Creation Time: 2020 - Feb - 24
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type SampleServer struct {
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
	res.Timestamp = time.Now().UnixNano()
	res.Delay = res.Timestamp - req.Timestamp
}

func (h *SampleServer) Ask(ctx *edge.RequestCtx, req *pb.AskRequest, res *pb.AskResponse) {
	res.Responder = req.ServerID
	res.Coordinator = req.ServerID
}
