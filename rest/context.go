package rest

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/gateway"
	"mime/multipart"
)

/*
   Creation Time: 2021 - Apr - 22
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Request struct {
	reqCtx *gateway.RequestCtx
	conn   rony.Conn
}

func (ctx *Request) MultiPart() (*multipart.Form, error) {
	return ctx.reqCtx.MultipartForm()
}

func (ctx *Request) Set(key string, value interface{}) {
	ctx.conn.Set(key, value)
}

func (ctx *Request) Get(key string) interface{} {
	return ctx.conn.Get(key)
}

func (ctx *Request) GetInt64(key string, defaultValue int64) int64 {
	v, ok := ctx.conn.Get(key).(int64)
	if !ok {
		return defaultValue
	}
	return v
}

func (ctx *Request) GetString(key string, defaultValue string) string {
	v, ok := ctx.conn.Get(key).(string)
	if !ok {
		return defaultValue
	}
	return v
}

func (ctx *Request) GetInt32(key string, defaultValue int32) int32 {
	v, ok := ctx.conn.Get(key).(int32)
	if !ok {
		return defaultValue
	}
	return v
}
