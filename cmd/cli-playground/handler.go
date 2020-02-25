package main

import (
	"git.ronaksoftware.com/ronak/rony/context"
	"git.ronaksoftware.com/ronak/rony/errors"
	"git.ronaksoftware.com/ronak/rony/msg"
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

func EchoHandler(ctx *context.Context, in *msg.MessageEnvelope) {
	// Write your handler code and remove the following line
	req := msg.PoolEchoRequest.Get()
	defer msg.PoolEchoRequest.Put(req)
	res := msg.PoolEchoResponse.Get()
	defer msg.PoolEchoResponse.Put(res)
	err := req.Unmarshal(in.Message)
	if err != nil {
		ctx.PushError(in.RequestID, errors.ErrCodeInvalid, errors.ErrItemRequest)
		return
	}

	res.Bool = req.Bool
	res.Int = req.Int
	res.Timestamp = time.Now().UnixNano()
	res.Delay = res.Timestamp - req.Timestamp

	ctx.PushMessage(ctx.AuthID, in.RequestID, msg.C_EchoResponse, res)
}