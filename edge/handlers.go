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

type (
	Handler = func(ctx *RequestCtx, in *rony.MessageEnvelope)
)

// HandlerOption is a structure holds all the information required for a handler.
type HandlerOption struct {
	constructors map[int64]struct{}
	handlers     []Handler
	tunnel       bool
	gateway      bool

	// internal use
	builtin bool
}

func NewHandlerOptions() *HandlerOption {
	return &HandlerOption{
		constructors: map[int64]struct{}{},
		handlers:     nil,
		gateway:      true,
		tunnel:       true,
	}
}

func (ho *HandlerOption) setBuiltin() *HandlerOption {
	ho.builtin = true
	return ho
}

func (ho *HandlerOption) SetConstructor(constructors ...int64) *HandlerOption {
	for _, c := range constructors {
		ho.constructors[c] = struct{}{}
	}
	return ho
}

// SetHandler replaces the handlers for this constructor with h
func (ho *HandlerOption) SetHandler(h ...Handler) *HandlerOption {
	ho.handlers = append(ho.handlers[:0], h...)
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
