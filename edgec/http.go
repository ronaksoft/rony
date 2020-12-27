package edgec

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/pools"
	"github.com/valyala/fasthttp"
	"google.golang.org/protobuf/proto"
	"net/http"
	"sync/atomic"
	"time"
)

/*
   Creation Time: 2020 - Dec - 27
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type HttpConfig struct {
	Name         string
	HostPort     string
	Header       map[string]string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type Http struct {
	hostPort string
	reqID    uint64
	c        *fasthttp.Client
}

func NewHttp(config HttpConfig) *Http {
	h := &Http{
		hostPort: config.HostPort,
		c: &fasthttp.Client{
			Name:                      config.Name,
			MaxIdemponentCallAttempts: 10,
			ReadTimeout:               config.ReadTimeout,
			WriteTimeout:              config.WriteTimeout,
			MaxResponseBodySize:       0,
			RetryIf:                   nil,
		},
	}
	return h
}

func (h *Http) Send(req *rony.MessageEnvelope, res *rony.MessageEnvelope) (err error) {
	mo := proto.MarshalOptions{UseCachedSize: true}
	b := pools.Bytes.GetCap(mo.Size(req))
	defer pools.Bytes.Put(b)

	b, err = mo.MarshalAppend(b, req)
	if err != nil {
		return err
	}

	httpReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(httpReq)

	httpReq.Header.SetMethod(http.MethodPost)
	httpReq.SetHost(h.hostPort)
	httpReq.SetBody(b)

	httpRes := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(httpRes)
	err = h.c.Do(httpReq, httpRes)
	if err != nil {
		return
	}

	err = res.Unmarshal(httpRes.Body())
	return
}

func (h *Http) Close() error {
	h.c.CloseIdleConnections()
	return nil
}

func (h *Http) GetRequestID() uint64 {
	return atomic.AddUint64(&h.reqID, 1)
}
