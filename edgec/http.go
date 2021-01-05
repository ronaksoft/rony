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

// HttpConfig holds the configurations for the Http client.
type HttpConfig struct {
	Name         string
	HostPort     string
	Header       map[string]string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Retries      int
}

// Http connects to edge servers with HTTP transport.
type Http struct {
	cfg   HttpConfig
	reqID uint64
	c     *fasthttp.Client
}

func NewHttp(config HttpConfig) *Http {
	h := &Http{
		cfg: config,
		c: &fasthttp.Client{
			Name:                      config.Name,
			MaxIdemponentCallAttempts: 10,
			ReadTimeout:               config.ReadTimeout,
			WriteTimeout:              config.WriteTimeout,
			MaxResponseBodySize:       0,
		},
	}
	return h
}

// Send implements Client interface
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
	httpReq.SetHost(h.cfg.HostPort)
	httpReq.SetBody(b)

	httpRes := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(httpRes)

	retries := h.cfg.Retries
Retries:
	err = h.c.Do(httpReq, httpRes)
	switch err {
	case fasthttp.ErrNoFreeConns:
		retries--
		goto Retries
	}

	if err != nil {
		return
	}
	err = res.Unmarshal(httpRes.Body())
	return
}

// Close implements Client interface
func (h *Http) Close() error {
	h.c.CloseIdleConnections()
	return nil
}

// GetRequestID implements Client interface
func (h *Http) GetRequestID() uint64 {
	return atomic.AddUint64(&h.reqID, 1)
}
