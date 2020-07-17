package main

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/edge"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv/pb"
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

func GenAskHandler(serverID string) edge.Handler {
	return func(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
		req := pb.PoolAskRequest.Get()
		defer pb.PoolAskRequest.Put(req)
		res := pb.PoolAskResponse.Get()
		defer pb.PoolAskResponse.Put(res)
		err := req.Unmarshal(in.Message)
		if err != nil {
			ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)
			return
		}

		if req.ServerID != serverID {
			ctx.PushClusterMessage(req.ServerID, ctx.AuthID(), in.RequestID, in.Constructor, req)
		} else {
			res.Responder = serverID
			ctx.PushMessage(pb.C_AskResponse, res)
		}
	}
}

func GenEchoHandler(serverID string) edge.Handler {
	return func(ctx *edge.RequestCtx, in *rony.MessageEnvelope) {
		req := pb.PoolEchoRequest.Get()
		defer pb.PoolEchoRequest.Put(req)
		res := pb.PoolEchoResponse.Get()
		defer pb.PoolEchoResponse.Put(res)
		err := req.Unmarshal(in.Message)
		if err != nil {
			ctx.PushError(rony.ErrCodeInvalid, rony.ErrItemRequest)
			return
		}

		if req.Bool {
			fmt.Println(serverID)
		}

		res.Bool = req.Bool
		res.Int = req.Int
		res.ServerID = serverID
		res.Timestamp = time.Now().UnixNano()
		res.Delay = res.Timestamp - req.Timestamp

		ctx.PushMessage(pb.C_EchoResponse, res)
	}

}
