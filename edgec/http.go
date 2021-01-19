package edgec

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/log"
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
	Name           string
	SeedHostPort   string
	Header         map[string]string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	ContextTimeout time.Duration
	Retries        int
	Router         Router
	Secure         bool
}

// Http connects to edge servers with HTTP transport.
type Http struct {
	cfg            HttpConfig
	reqID          uint64
	c              *fasthttp.Client
	mtx            sync.RWMutex
	sessionReplica uint64
	hosts          map[uint64]map[string]*httpConn // holds host by replicaSet/hostID
	leaders        map[uint64]string
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
	if h.cfg.Router == nil {
		h.cfg.Router = &httpRouter{
			c: h,
		}
	}

	return h
}

func (h *Http) addConn(serverID string, replicaSet uint64, leader bool, c *httpConn) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	if h.hosts[replicaSet] == nil {
		h.hosts[replicaSet] = make(map[string]*httpConn, 16)
	}
	h.hosts[replicaSet][serverID] = c
	if leader || replicaSet == 0 {
		h.leaders[replicaSet] = serverID
	}
}

func (h *Http) getConn(replicaSet uint64, onlyLeader bool) *httpConn {
	h.mtx.RLock()
	defer h.mtx.RUnlock()

	if onlyLeader {
		leaderID := h.leaders[replicaSet]
		if leaderID == "" {
			return nil
		}
		m := h.hosts[replicaSet]
		if m != nil {
			c := m[leaderID]
			return c
		}
	} else {
		m := h.hosts[replicaSet]
		if m != nil {
			for _, c := range m {
				return c
			}
		}
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
	err := h.initConn()
	if err != nil {
		return err
	}
	return nil
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
	case rony.C_NodeInfoMany:
		x := &rony.NodeInfoMany{}
		_ = x.Unmarshal(res.Message)
		found := false
		for _, n := range x.Nodes {
			if ce := log.Check(log.DebugLevel, "NodeInfo"); ce != nil {
				ce.Write(
					zap.String("ServerID", n.ServerID),
					zap.Uint64("RS", n.ReplicaSet),
					zap.Bool("Leader", n.Leader),
					zap.Strings("HostPorts", n.HostPorts),
				)
			}
			httpc := h.newConn(n.ServerID, n.ReplicaSet, n.HostPorts...)
			if !found {
				for _, hp := range n.HostPorts {
					if hp == initConn.hostPorts[0] {
						httpc = initConn
						httpc.hostPorts = n.HostPorts
						found = true
					}
				}
			}

			h.addConn(n.ServerID, n.ReplicaSet, n.Leader, httpc)
			if n.Leader {
				h.sessionReplica = n.ReplicaSet
			}
		}
	default:
		fmt.Println(res)
		return ErrUnknownResponse

	}

	return nil
}

func (h *Http) Send(req *rony.MessageEnvelope, res *rony.MessageEnvelope, leaderOnly bool) error {
	return h.SendWithDetails(req, res, h.cfg.ContextTimeout, leaderOnly)
}

// Send implements Client interface
func (h *Http) SendWithDetails(req *rony.MessageEnvelope, res *rony.MessageEnvelope, timeout time.Duration, leaderOnly bool) (err error) {
	rs := h.cfg.Router.GetRoute(req)
	hc := h.getConn(rs, leaderOnly)
	if ce := log.Check(log.DebugLevel, "Send"); ce != nil {
		ce.Write(
			zap.Uint64("ReqID", req.RequestID),
			zap.Uint64("RS", rs),
			zap.Bool("LeaderOnly", leaderOnly),
		)
	}

	if hc == nil {
		return ErrNoConnection
	}

SendLoop:
	rs, err = hc.send(req, res, timeout)
	switch err {
	case nil:
		return nil
	case ErrReplicaMaster:
		leaderOnly = true
		hc = h.getConn(rs, leaderOnly)
	case ErrReplicaSetSession, ErrReplicaSetRequest:
		rs = h.sessionReplica
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
