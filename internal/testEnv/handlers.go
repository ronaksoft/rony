package testEnv

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/internal/testEnv/pb"
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

type Handlers struct{}

func (h Handlers) Func1(ctx *edge.RequestCtx, req *pb.Req1, res *pb.Res1) {
	panic("implement me")
}

func (h Handlers) Func2(ctx *edge.RequestCtx, req *pb.Req2, res *pb.Res2) {
	panic("implement me")
}

func (h Handlers) Echo(ctx *edge.RequestCtx, req *pb.EchoRequest, res *pb.EchoResponse) {
	res.Int = req.Int
	res.Bool = req.Bool
}

func (h Handlers) Ask(ctx *edge.RequestCtx, req *pb.AskRequest, res *pb.AskResponse) {
	panic("implement me")
}

func EchoSimple(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := pb.EchoRequest{}
	err := req.Unmarshal(in.Message)
	if err != nil {
		ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)
		return
	}

	ctx.PushMessage(pb.C_EchoResponse, &pb.EchoResponse{
		Int:       req.Int,
		Bool:      req.Bool,
		Timestamp: req.Timestamp,
		Delay:     0,
		ServerID:  "ServerID",
	})
}

func EchoHold(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
	req := pb.EchoRequest{}
	err := req.Unmarshal(in.Message)
	if err != nil {
		ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)
		return
	}
	time.Sleep(time.Second * 10)
	ctx.PushMessage(pb.C_EchoResponse, &pb.EchoResponse{
		Int:       req.Int,
		Bool:      req.Bool,
		Timestamp: req.Timestamp,
		Delay:     0,
		ServerID:  "ServerID",
	})
}
