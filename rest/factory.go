package rest

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/gateway"
	"github.com/ronaksoft/rony/pools"
)

/*
   Creation Time: 2021 - Apr - 22
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Factory struct {
	onRequest  func(ctx *Context) error
	onResponse func(envelope *rony.MessageEnvelope) (*pools.ByteBuffer, map[string]string)
}

func NewFactory(
	onRequest func(ctx *Context) error,
	onResponse func(envelope *rony.MessageEnvelope) (*pools.ByteBuffer, map[string]string),
) *Factory {
	return &Factory{
		onRequest:  onRequest,
		onResponse: onResponse,
	}
}

func (f *Factory) Get() gateway.ProxyHandle {
	return &Handle{
		onRequest:  f.onRequest,
		onResponse: f.onResponse,
	}
}

func (f *Factory) Release(h gateway.ProxyHandle) {

}
