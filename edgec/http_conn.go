package edgec

import (
	"net/http"
	"time"

	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/log"
	"github.com/ronaksoft/rony/pools"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

/*
   Creation Time: 2021 - Jan - 05
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type httpConn struct {
	h          *Http
	replicaSet uint64
	id         string
	hostPorts  []string
	secure     bool
}

func (c *httpConn) send(req, res *rony.MessageEnvelope, timeout time.Duration) (replicaSet uint64, err error) {
	replicaSet = c.replicaSet

	mo := proto.MarshalOptions{UseCachedSize: true}
	buf := pools.Buffer.GetCap(mo.Size(req))
	defer pools.Buffer.Put(buf)

	b, err := mo.MarshalAppend(*buf.Bytes(), req)
	if err != nil {
		return
	}

	httpReq := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(httpReq)

	httpReq.Header.SetMethod(http.MethodPost)
	if c.secure {
		httpReq.URI().SetScheme("https")
	}
	httpReq.SetHost(c.hostPorts[0])
	httpReq.SetBody(b)
	if hf := c.h.cfg.HeaderFunc; hf != nil {
		for k, v := range hf() {
			httpReq.Header.Set(k, v)
		}
	}

	httpRes := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(httpRes)

SendLoop:
	err = c.h.c.DoTimeout(httpReq, httpRes, timeout)
	switch err {
	case fasthttp.ErrNoFreeConns:
		goto SendLoop
	}
	if err != nil {
		return
	}
	err = res.Unmarshal(httpRes.Body())
	if err != nil {
		return
	}
	switch res.GetConstructor() {
	case rony.C_Redirect:
		x := &rony.Redirect{}
		err = proto.Unmarshal(res.Message, x)
		if err != nil {
			return
		}
		replicaSet, err = c.redirect(x)

		return
	}

	return
}

func (c *httpConn) redirect(x *rony.Redirect) (replicaSet uint64, err error) {
	if ce := c.h.logger.Check(log.InfoLevel, "Redirect"); ce != nil {
		ce.Write(
			zap.Any("Edges", x.Edges),
			zap.Any("Wait", x.WaitInSec),
		)
	}

	replicaSet = x.Edges[0].ReplicaSet
	for _, n := range x.Edges {
		c.h.addConn(n.ServerID, n.ReplicaSet, c.h.newConn(n.ServerID, n.ReplicaSet, n.HostPorts...))
	}

	switch x.Reason {
	case rony.RedirectReason_ReplicaSetSession:
		c.h.sessionReplica = replicaSet
		err = ErrReplicaSetSession
	case rony.RedirectReason_ReplicaSetRequest:
		replicaSet = x.Edges[0].ReplicaSet
		err = ErrReplicaSetRequest
	default:
		err = ErrUnknownResponse
	}

	return
}
