package edgetest

import (
	"github.com/ronaksoft/rony"
	dummyGateway "github.com/ronaksoft/rony/internal/gateway/dummy"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/proto"
	"sync"
	"sync/atomic"
	"time"
)

/*
   Creation Time: 2020 - Dec - 09
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type rpcCtx struct {
	mtx        sync.Mutex
	id         uint64
	reqC       int64
	reqID      uint64
	req        []byte
	expect     map[int64]CheckFunc
	gw         *dummyGateway.Gateway
	err        error
	errH       func(constructor int64, e *rony.Error)
	doneCh     chan struct{}
	kvs        []*rony.KeyValue
	persistent bool
}

func newRPCContext(gw *dummyGateway.Gateway) *rpcCtx {
	c := &rpcCtx{
		id:     atomic.AddUint64(&connID, 1),
		expect: make(map[int64]CheckFunc),
		gw:     gw,
		doneCh: make(chan struct{}, 1),
	}
	return c
}

// Persistent makes the rpcCtx simulate persistent connection e.g. websocket
func (c *rpcCtx) Persistent() *rpcCtx {
	c.persistent = true
	return c
}

// Request set the request you wish to send to the server
func (c *rpcCtx) Request(constructor int64, p proto.Message, kv ...*rony.KeyValue) *rpcCtx {
	data, _ := proto.Marshal(p)
	c.reqID = tools.RandomUint64(0)
	e := &rony.MessageEnvelope{
		Constructor: constructor,
		RequestID:   c.reqID,
		Message:     data,
		Auth:        nil,
		Header:      kv,
	}
	c.reqC = constructor
	c.req, c.err = proto.Marshal(e)
	return c
}

// Expect let you set what you expect to receive. If cf is set, then you can do more checks on the response and return error
// if the response was not fully acceptable
func (c *rpcCtx) Expect(constructor int64, cf CheckFunc) *rpcCtx {
	c.expect[constructor] = cf
	return c
}

func (c *rpcCtx) ExpectConstructor(constructor int64) *rpcCtx {
	return c.Expect(constructor, nil)
}

func (c *rpcCtx) check(e *rony.MessageEnvelope) {
	c.mtx.Lock()
	f, ok := c.expect[e.Constructor]
	c.mtx.Unlock()
	if !ok && e.Constructor == rony.C_Error {
		err := &rony.Error{}
		c.err = err.Unmarshal(e.Message)
		if c.errH != nil {
			c.errH(c.reqC, err)
		}
		return
	}
	if f != nil {
		c.err = f(e.Message, e.Header...)
	}
	c.mtx.Lock()
	delete(c.expect, e.Constructor)
	c.mtx.Unlock()
}

func (c *rpcCtx) expectCount() int {
	c.mtx.Lock()
	n := len(c.expect)
	c.mtx.Unlock()
	return n
}

func (c *rpcCtx) ErrorHandler(f func(constructor int64, e *rony.Error)) *rpcCtx {
	c.errH = f
	return c
}

func (c *rpcCtx) SetRunParameters(kvs ...*rony.KeyValue) *rpcCtx {
	c.kvs = kvs
	return c
}

func (c *rpcCtx) RunShort(kvs ...*rony.KeyValue) error {
	return c.Run(time.Second*10, kvs...)
}

func (c *rpcCtx) RunLong(kvs ...*rony.KeyValue) error {
	return c.Run(time.Minute, kvs...)
}

func (c *rpcCtx) Run(timeout time.Duration, kvs ...*rony.KeyValue) error {
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

func (c *rpcCtx) receiver(connID uint64, streamID int64, data []byte) {
	defer func() {
		c.doneCh <- struct{}{}
	}()
	e := &rony.MessageEnvelope{}
	c.err = e.Unmarshal(data)
	if c.err != nil {
		return
	}
	switch e.Constructor {
	case rony.C_MessageContainer:
		mc := &rony.MessageContainer{}
		c.err = mc.Unmarshal(e.Message)
		if c.err != nil {
			return
		}
		for _, e := range mc.Envelopes {
			c.check(e)
		}
	default:
		c.check(e)
	}
}
