package rest

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/gateway"
)

/*
   Creation Time: 2021 - Apr - 22
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type RequestHandler func(ctx *Context) error
type ResponseHandler func(envelope *rony.MessageEnvelope, bodyWriter gateway.BodyWriter, hdrWriter *gateway.HeaderWriter)

type Factory struct {
	onRequest  RequestHandler
	onResponse ResponseHandler
}

func NewFactory(
	onRequest RequestHandler,
	onResponse ResponseHandler,
) *Factory {
	return &Factory{
		onRequest:  onRequest,
		onResponse: onResponse,
	}
}

func (f *Factory) Get() gateway.ProxyHandle {
	return &restHandler{
		requestHandler:  f.onRequest,
		responseHandler: f.onResponse,
	}
}
