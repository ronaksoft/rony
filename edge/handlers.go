package edge

import (
	"github.com/ronaksoft/rony"
)

/*
   Creation Time: 2021 - Jan - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Handler func(ctx *RequestCtx, in *rony.MessageEnvelope)

type HandlerOptions struct {
	pre  []Handler
	post []Handler
}

func NewHandlerOptions() *HandlerOptions {
	return &HandlerOptions{}
}

func (ho *HandlerOptions) SetPreHandlers(h ...Handler) *HandlerOptions {
	ho.pre = append(ho.pre[:0], h...)
	return ho
}

func (ho *HandlerOptions) SetPostHandlers(h ...Handler) *HandlerOptions {
	ho.post = append(ho.post[:0], h...)
	return ho
}

func (ho *HandlerOptions) ApplyTo(h ...Handler) []Handler {
	out := make([]Handler, 0, len(ho.pre)+len(h)+len(ho.post))
	out = append(out, ho.pre...)
	out = append(out, h...)
	out = append(out, ho.post...)
	return out
}
