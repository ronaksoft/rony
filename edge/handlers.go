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

type Handler = func(ctx *RequestCtx, in *rony.MessageEnvelope)

// HandlerOption is a structure holds all the information required for a handler.
type HandlerOption struct {
	constructor      int64
	handlers         []Handler
	inconsistentRead bool
	tunnel           bool
	gateway          bool
	// internal user
	builtin bool
}

func NewHandlerOptions(constructor int64, h ...Handler) *HandlerOption {
	return &HandlerOption{
		constructor:      constructor,
		handlers:         h,
		gateway:          true,
		tunnel:           true,
		inconsistentRead: false,
	}
}

func (ho *HandlerOption) setBuiltin() *HandlerOption {
	ho.builtin = true
	return ho
}

// GatewayOnly makes this method only available through gateway messages.
func (ho *HandlerOption) GatewayOnly() *HandlerOption {
	ho.tunnel = false
	ho.gateway = true
	return ho
}

// TunnelOnly makes this method only available through tunnel messages.
func (ho *HandlerOption) TunnelOnly() *HandlerOption {
	ho.tunnel = true
	ho.gateway = false
	return ho
}

// InconsistentRead makes this method (constructor) available on edges in follower state
func (ho *HandlerOption) InconsistentRead() *HandlerOption {
	ho.inconsistentRead = true
	return ho
}

// Set replaces the handlers for this constructor with h
func (ho *HandlerOption) Set(h ...Handler) *HandlerOption {
	ho.handlers = append(ho.handlers[:0], h...)
	return ho
}

// Prepend adds the h handlers before already set handlers
func (ho *HandlerOption) Prepend(h ...Handler) *HandlerOption {
	nh := make([]Handler, 0, len(ho.handlers)+len(h))
	nh = append(nh, h...)
	nh = append(nh, ho.handlers...)
	ho.handlers = nh
	return ho
}

// Append adds the h handlers after already set handlers
func (ho *HandlerOption) Append(h ...Handler) *HandlerOption {
	ho.handlers = append(ho.handlers, h...)
	return ho
}
