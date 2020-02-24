package context

import "sync"

/*
   Creation Time: 2020 - Feb - 22
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var ctxPool = sync.Pool{}

func Acquire(connID uint64, authID int64, quickReturn, blocking bool) *Context {
	var ctx *Context
	if v := ctxPool.Get(); v == nil {
		ctx = New()
	} else {
		ctx = v.(*Context)
		ctx.Clear()
		// Just to make sure channel is empty, or empty it if not
		select {
		case <-ctx.NextChan:
		default:
		}
	}
	ctx.Stop = false
	ctx.QuickReturn = quickReturn
	ctx.ConnID = connID
	ctx.AuthID = authID
	ctx.Blocking = blocking
	return ctx
}

func Release(ctx *Context) {
	ctxPool.Put(ctx)
}
