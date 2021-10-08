package edgec

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/errors"
	"github.com/ronaksoft/rony/log"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"sync"
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
	Name            string
	SeedHostPort    string
	HeaderFunc      func() map[string]string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ContextTimeout  time.Duration
	RequestMaxRetry int
	Router          Router
	Secure          bool
}

// Http connects to edge servers with HTTP transport.
type Http struct {
	cfg            HttpConfig
	reqID          uint64
	c              *fasthttp.Client
	mtx            sync.RWMutex
	sessionReplica uint64
	hosts          map[uint64]map[string]*httpConn // holds host by replicaSet/hostID
	logger         log.Logger
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
		hosts:  make(map[uint64]map[string]*httpConn, 32),
		logger: log.With("EdgeC(Http)"),
	}

	if h.cfg.Router == nil {
		h.cfg.Router = &httpRouter{
			c: h,
		}
	}
	if h.cfg.RequestMaxRetry == 0 {
		h.cfg.RequestMaxRetry = requestRetry
	}
	if h.cfg.ContextTimeout == 0 {
		h.cfg.ContextTimeout = requestTimeout
	}

	return h
}

func (h *Http) addConn(serverID string, replicaSet uint64, c *httpConn) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	if h.hosts[replicaSet] == nil {
		h.hosts[replicaSet] = make(map[string]*httpConn, 16)
	}
	h.hosts[replicaSet][serverID] = c
}

func (h *Http) getConn(replicaSet uint64) *httpConn {
	h.mtx.RLock()
	defer h.mtx.RUnlock()

	m := h.hosts[replicaSet]
	for _, c := range m {
		return c
	}

	return nil
}

func (h *Http) newConn(id string, replicaSet uint64, hostPorts ...string) *httpConn {
	return &httpConn{
		id:         id,
		h:          h,
		replicaSet: replicaSet,
		hostPorts:  hostPorts,
		secure:     h.cfg.Secure,
	}
}

func (h *Http) Start() error {
	return h.initConn()
}

func (h *Http) initConn() error {
	initConn := h.newConn("", 0, h.cfg.SeedHostPort)
	req := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(req)
	res := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(res)
	req.Fill(h.GetRequestID(), rony.C_GetNodes, &rony.GetNodes{})
	sessionReplica, err := initConn.send(req, res, requestTimeout)
	if err != nil {
		return err
	}
	h.sessionReplica = sessionReplica
	switch res.Constructor {
	case rony.C_Edges:
		x := &rony.Edges{}
		_ = x.Unmarshal(res.Message)
		for _, n := range x.Nodes {
			if ce := h.logger.Check(log.DebugLevel, "NodeInfo"); ce != nil {
				ce.Write(
					zap.String("ServerID", n.ServerID),
					zap.Uint64("RS", n.ReplicaSet),
					zap.Strings("HostPorts", n.HostPorts),
				)
			}
			httpc := h.newConn(n.ServerID, n.ReplicaSet, n.HostPorts...)
			h.addConn(n.ServerID, n.ReplicaSet, httpc)
			h.sessionReplica = n.ReplicaSet
		}
	default:
		return ErrUnknownResponse
	}

	return nil
}

func (h *Http) Send(req *rony.MessageEnvelope, res *rony.MessageEnvelope) error {
	return h.SendWithDetails(req, res, h.cfg.RequestMaxRetry, h.cfg.ContextTimeout)
}

func (h *Http) SendWithDetails(
	req *rony.MessageEnvelope, res *rony.MessageEnvelope,
	retry int, timeout time.Duration,
) (err error) {
	rs := h.cfg.Router.GetRoute(req)
	hc := h.getConn(rs)
	if hc == nil {
		return ErrNoConnection
	}

SendLoop:
	if ce := h.logger.Check(log.DebugLevel, "sending"); ce != nil {
		ce.Write(
			zap.Uint64("ReqID", req.RequestID),
			zap.Uint64("RS", rs),
			zap.Int("Retry", retry),
		)
	}

	rs, err = hc.send(req, res, timeout)
	switch err {
	case nil:
		return nil
	case ErrReplicaSetSession, ErrReplicaSetRequest:
		rs = h.sessionReplica
	}

	// If we exceeds the maximum retry then we return
	if retry--; retry < 0 {
		err = errors.ErrRetriesExceeded(err)
		return
	}

	goto SendLoop
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

type httpRouter struct {
	c *Http
}

func (d *httpRouter) UpdateRoute(req *rony.MessageEnvelope, replicaSet uint64) {
	// TODO:: implement cache maybe
}

func (d *httpRouter) GetRoute(req *rony.MessageEnvelope) (replicaSet uint64) {
	return d.c.sessionReplica
}
