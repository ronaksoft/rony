package edgetest

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/goccy/go-json"
	"github.com/ronaksoft/rony/registry"

	"github.com/ronaksoft/rony"
	dummyGateway "github.com/ronaksoft/rony/internal/gateway/dummy"
	"github.com/ronaksoft/rony/tools"
)

/*
   Creation Time: 2020 - Dec - 09
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type jrpcCtx struct {
	mtx        sync.Mutex
	id         uint64
	reqC       uint64
	reqID      uint64
	req        []byte
	expect     map[uint64]CheckFunc
	gw         *dummyGateway.Gateway
	err        error
	errH       func(constructor uint64, e *rony.Error)
	doneCh     chan struct{}
	kvs        []*rony.KeyValue
	persistent bool
}

func newJSONRPCContext(gw *dummyGateway.Gateway) *jrpcCtx {
	c := &jrpcCtx{
		id:     atomic.AddUint64(&connID, 1),
		expect: make(map[uint64]CheckFunc),
		gw:     gw,
		doneCh: make(chan struct{}, 1),
	}

	return c
}

// Persistent makes the jrpcCtx simulate persistent connection e.g. websocket
func (c *jrpcCtx) Persistent() *jrpcCtx {
	c.persistent = true

	return c
}

// Request set the request you wish to send to the server
func (c *jrpcCtx) Request(constructor uint64, p rony.IMessage, kvs ...*rony.KeyValue) *jrpcCtx {
	c.reqID = tools.RandomUint64(0)
	c.reqC = constructor

	data, _ := p.MarshalJSON()
	e := &rony.MessageEnvelopeJSON{
		Constructor: registry.C(constructor),
		RequestID:   c.reqID,
		Message:     data,
		Header:      map[string]string{},
	}
	for _, kv := range kvs {
		e.Header[kv.Key] = kv.Value
	}
	c.req, c.err = json.Marshal(e)

	return c
}

// Expect let you set what you expect to receive. If cf is set, then you can do more checks
// on the response and return error if the response was not fully acceptable
func (c *jrpcCtx) Expect(constructor uint64, cf CheckFunc) *jrpcCtx {
	c.expect[constructor] = cf

	return c
}

func (c *jrpcCtx) ExpectConstructor(constructor uint64) *jrpcCtx {
	return c.Expect(constructor, nil)
}

func (c *jrpcCtx) check(e *rony.MessageEnvelopeJSON) {
	c.mtx.Lock()
	f, ok := c.expect[registry.N(e.Constructor)]
	c.mtx.Unlock()
	if !ok && registry.N(e.Constructor) == rony.C_Error {
		err := &rony.Error{}
		c.err = err.Unmarshal(e.Message)
		if c.errH != nil {
			c.errH(c.reqC, err)
		}

		return
	}
	if f != nil {
		var kvs = make([]*rony.KeyValue, 0, len(e.Header))
		for k, v := range e.Header {
			kvs = append(kvs, &rony.KeyValue{Key: k, Value: v})
		}
		c.err = f(e.Message, kvs...)
	}
	c.mtx.Lock()
	delete(c.expect, registry.N(e.Constructor))
	c.mtx.Unlock()
}

func (c *jrpcCtx) expectCount() int {
	c.mtx.Lock()
	n := len(c.expect)
	c.mtx.Unlock()

	return n
}

func (c *jrpcCtx) ErrorHandler(f func(constructor uint64, e *rony.Error)) *jrpcCtx {
	c.errH = f

	return c
}

func (c *jrpcCtx) SetRunParameters(kvs ...*rony.KeyValue) *jrpcCtx {
	c.kvs = kvs

	return c
}

func (c *jrpcCtx) RunShort(kvs ...*rony.KeyValue) error {
	return c.Run(time.Second*10, kvs...)
}

func (c *jrpcCtx) RunLong(kvs ...*rony.KeyValue) error {
	return c.Run(time.Minute, kvs...)
}

func (c *jrpcCtx) Run(timeout time.Duration, kvs ...*rony.KeyValue) error {
	// We return error early if we have encountered error before Run
	if c.err != nil {
		return c.err
	}

	c.SetRunParameters(kvs...)

	// Open Connection
	c.gw.OpenConn(c.id, c.persistent, c.receiver, c.kvs...)

	// Send the Request
	err := c.gw.RPC(c.id, 0, c.req)
	if err != nil {
		return err
	}

	// Wait for Response(s)
Loop:
	for {
		select {
		case <-c.doneCh:
			// Check if all the expectations have been passed
			if c.expectCount() == 0 {
				break Loop
			}
		case <-time.After(timeout):
			break Loop
		}
	}

	if c.expectCount() > 0 {
		c.err = ErrExpectationFailed
	}

	return c.err
}

func (c *jrpcCtx) receiver(connID uint64, streamID int64, data []byte) {
	defer func() {
		c.doneCh <- struct{}{}
	}()
	e := &rony.MessageEnvelopeJSON{}
	c.err = json.Unmarshal(data, e)
	if c.err != nil {
		return
	}
	c.check(e)
}
