package pools

import (
	"git.ronaksoftware.com/ronak/rony/context"
	"sync"
)

/*
   Creation Time: 2019 - Oct - 14
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var ctxPool = sync.Pool{}

func AcquireContext(authID int64, userID int64, quickReturn, blocking bool) *context.Context {
	var ctx *context.Context
	if v := ctxPool.Get(); v == nil {
		ctx = context.New()
	} else {
		ctx = v.(*context.Context)
		ctx.Clear()
		// Just to make sure channel is empty, or empty it if not
		select {
		case <-ctx.NextChan:
		default:
		}
	}
	ctx.Stop = false
	ctx.QuickReturn = quickReturn
	ctx.AuthID = authID
	ctx.UserID = userID
	ctx.Blocking = blocking
	return ctx
}

func ReleaseContext(ctx *context.Context) {
	ctxPool.Put(ctx)
}
