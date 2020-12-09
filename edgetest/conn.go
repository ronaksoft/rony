package edgetest

import (
	"github.com/ronaksoft/rony"
	dummyGateway "github.com/ronaksoft/rony/gateway/dummy"
	"github.com/ronaksoft/rony/tools"
	"google.golang.org/protobuf/proto"
	"sync"
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

type conn struct {
	mtx         sync.Mutex
	id          uint64
	req         []byte
	expect      map[int64]CheckFunc
	gw          *dummyGateway.Gateway
	err         error
	receiveChan chan struct{}
}

func newConn(gw *dummyGateway.Gateway) *conn {
	c := &conn{
		id:          tools.RandomUint64(0),
		expect:      make(map[int64]CheckFunc),
		receiveChan: make(chan struct{}, 10),
		gw:          gw,
	}
	return c
}

// Request set the request you wish to send to the server
func (c *conn) Request(constructor int64, p proto.Message, kv ...*rony.KeyValue) *conn {
	data, _ := proto.Marshal(p)
	e := &rony.MessageEnvelope{
		Constructor: constructor,
		RequestID:   tools.RandomUint64(0),
		Message:     data,
		Auth:        nil,
		Header:      kv,
	}
	c.req, c.err = proto.Marshal(e)
	return c
}

// Expect let you set what you expect to receive. If cf is set, then you can do more checks on the response and return error
// if the response was not fully acceptable
func (c *conn) Expect(constructor int64, cf CheckFunc) *conn {
	c.expect[constructor] = cf
	return c
}

func (c *conn) check(e *rony.MessageEnvelope) {
	c.receiveChan <- struct{}{}
	c.mtx.Lock()
	f := c.expect[e.Constructor]
	c.mtx.Unlock()
	if f == nil {
		return
	}
	c.err = f(e.Message, e.Auth, e.Header...)
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
func (c *conn) Run(timeout time.Duration) error {
	// Open Connection
	c.gw.OpenConn(c.id, func(connID uint64, streamID int64, data []byte) {
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
	})

	// Send the Request
	c.gw.SendToConn(c.id, 0, c.req)

	// Wait for Response(s)
	for {
		select {
		case <-c.receiveChan:
		case <-time.After(timeout):
			return ErrTimeout
		}
		if c.expectCount() == 0 {
			break
		}
	}

	// Check if all the expectations have been passed
	if c.expectCount() > 0 {
		c.err = ErrExpectationFailed
	}

	return c.err
}
