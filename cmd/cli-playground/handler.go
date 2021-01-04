package main

import (
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
}

func (h *SampleServer) Echo(ctx *edge.RequestCtx, req *pb.EchoRequest, res *pb.EchoResponse) {
	res.Int = req.Int
	res.Timestamp = tools.NanoTime()
	res.Delay = res.Timestamp - req.Timestamp
	res.ServerID = h.es.GetServerID()
}
