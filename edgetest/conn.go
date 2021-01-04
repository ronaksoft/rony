package edgetest

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/gateway"
	dummyGateway "github.com/ronaksoft/rony/gateway/dummy"
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

type CheckFunc func(b []byte, auth []byte, kv ...*rony.KeyValue) error

var (
	connID uint64 = 1
)

type conn struct {
	mtx    sync.Mutex
	id     uint64
	reqC   int64
	reqID  uint64
	req    []byte
	expect map[int64]CheckFunc
	gw     *dummyGateway.Gateway
	err    error
	errH   func(constructor int64, e *rony.Error)
	doneCh chan struct{}
}

func newConn(gw *dummyGateway.Gateway) *conn {
	c := &conn{
		id:     atomic.AddUint64(&connID, 1),
		expect: make(map[int64]CheckFunc),
		gw:     gw,
		doneCh: make(chan struct{}, 1),
	}
	return c
}

// Request set the request you wish to send to the server
func (c *conn) Request(constructor int64, p proto.Message, kv ...*rony.KeyValue) *conn {
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
func (c *conn) Expect(constructor int64, cf CheckFunc) *conn {
	c.expect[constructor] = cf
	return c
}

func (c *conn) ExpectConstructor(constructor int64) *conn {
	return c.Expect(constructor, nil)
}

func (c *conn) check(e *rony.MessageEnvelope) {
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
		c.err = f(e.Message, e.Auth, e.Header...)
	}
	c.mtx.Lock()
	delete(c.expect, e.Constructor)
	c.mtx.Unlock()
}

func (c *conn) expectCount() int {
	c.mtx.Lock()
	n := len(c.expect)
	c.mtx.Unlock()
	return n
}

func (c *conn) ErrorHandler(f func(constructor int64, e *rony.Error)) *conn {
	c.errH = f
	return c
}

func (c *conn) RunShort(kvs ...gateway.KeyValue) error {
	return c.Run(time.Second*10, kvs...)
}

func (c *conn) RunLong(kvs ...gateway.KeyValue) error {
	return c.Run(time.Minute, kvs...)
}

func (c *conn) Run(timeout time.Duration, kvs ...gateway.KeyValue) error {
	// Open Connection
	c.gw.OpenConn(c.id, func(connID uint64, streamID int64, data []byte) {
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
	}, kvs...)

	// Send the Request
	c.gw.SendToConn(c.id, 0, c.req)

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
