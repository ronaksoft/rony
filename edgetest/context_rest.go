package edgetest

import (
	"github.com/ronaksoft/rony"
	dummyGateway "github.com/ronaksoft/rony/internal/gateway/dummy"
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

type restCtx struct {
	id     uint64
	gw     *dummyGateway.Gateway
	err    error
	expect CheckFunc
	doneCh chan struct{}
	kvs    []*rony.KeyValue
	method string
	path   string
	body   []byte
}

func newRESTContext(gw *dummyGateway.Gateway) *restCtx {
	c := &restCtx{
		id:     atomic.AddUint64(&connID, 1),
		gw:     gw,
		doneCh: make(chan struct{}, 1),
	}

	return c
}

// Request set the request you wish to send to the server
func (c *restCtx) Request(method, path string, body []byte, kvs ...*rony.KeyValue) *restCtx {
	c.method = method
	c.path = path
	c.body = body
	c.kvs = kvs

	return c
}

// Expect let you set what you expect to receive. If cf is set, then you can do more checks on the response and return error
// if the response was not fully acceptable
func (c *restCtx) Expect(cf CheckFunc) *restCtx {
	c.expect = cf

	return c
}

func (c *restCtx) RunShort(kvs ...*rony.KeyValue) error {
	return c.Run(time.Second*10, kvs...)
}

func (c *restCtx) RunLong(kvs ...*rony.KeyValue) error {
	return c.Run(time.Minute, kvs...)
}

func (c *restCtx) Run(timeout time.Duration, kvs ...*rony.KeyValue) error {
	// We return error early if we have encountered error before Run
	if c.err != nil {
		return c.err
	}

	go func() {
		respBody, respHdr := c.gw.REST(c.id, c.method, c.path, c.body, kvs...)
		hdr := make([]*rony.KeyValue, 0, len(respHdr))
		for k, v := range respHdr {
			hdr = append(hdr, &rony.KeyValue{Key: k, Value: v})
		}
		c.err = c.expect(respBody, hdr...)
		c.doneCh <- struct{}{}
	}()

	// Wait for Response(s)
	select {
	case <-c.doneCh:
	case <-time.After(timeout):
		c.err = ErrTimeout
	}

	return c.err
}
