package edgec

import (
	"context"
	"github.com/ronaksoft/rony"
	log "github.com/ronaksoft/rony/internal/logger"
	"github.com/ronaksoft/rony/tools"
	"go.uber.org/zap"
	"sync/atomic"
	"time"
)

/*
   Creation Time: 2020 - Jul - 17
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

const (
	requestTimeout = 3 * time.Second
	requestRetry   = 5
	dialTimeout    = 3 * time.Second
	idleTimeout    = time.Minute
)

type MessageHandler func(m *rony.MessageEnvelope)

// WebsocketConfig holds the configs for the Websocket client
type WebsocketConfig struct {
	SeedHostPort string
	IdleTimeout  time.Duration
	DialTimeout  time.Duration
	Handler      MessageHandler
	Header       map[string]string
	Secure       bool
	ForceConnect bool
	// RequestMaxRetry is the maximum number client sends a request if any network layer error occurs
	RequestMaxRetry int
	// RequestTimeout is the timeout for each individual request on each try.
	RequestTimeout time.Duration
	// ContextTimeout is the amount that Send function will wait until times out. This includes all the retries.
	ContextTimeout time.Duration
	Router         Router
}

// Websocket client which could handle multiple connections
type Websocket struct {
	cfg            WebsocketConfig
	pool           *connPool
	sessionReplica uint64
	nextReqID      uint64
}

func NewWebsocket(config WebsocketConfig) *Websocket {
	c := &Websocket{
		nextReqID: tools.RandomUint64(0),
		pool:      newConnPool(),
		cfg:       config,
	}

	// Prepare default config values
	if c.cfg.DialTimeout == 0 {
		c.cfg.DialTimeout = dialTimeout
	}
	if c.cfg.IdleTimeout == 0 {
		c.cfg.IdleTimeout = idleTimeout
	}
	if c.cfg.RequestMaxRetry == 0 {
		c.cfg.RequestMaxRetry = requestRetry
	}
	if c.cfg.RequestTimeout == 0 {
		c.cfg.RequestTimeout = requestTimeout
	}
	if c.cfg.ContextTimeout == 0 {
		c.cfg.ContextTimeout = c.cfg.RequestTimeout * time.Duration(c.cfg.RequestMaxRetry)
	}
	if c.cfg.Router == nil {
		c.cfg.Router = &defaultRouter{
			c: c,
		}
	}

	return c
}

func (c *Websocket) Start() error {
	err := c.initConn()
	if err != nil {
		return err
	}
	return nil
}

func (c *Websocket) GetRequestID() uint64 {
	return atomic.AddUint64(&c.nextReqID, 1)
}

func (c *Websocket) newConn(id string, replicaSet uint64, hostPorts ...string) *wsConn {
	return &wsConn{
		id:         id,
		ws:         c,
		replicaSet: replicaSet,
		hostPorts:  hostPorts,
		pending:    make(map[uint64]chan *rony.MessageEnvelope, 100),
	}
}

func (c *Websocket) initConn() error {
	initConn := c.newConn("", 0, c.cfg.SeedHostPort)
	initConn.connect()
	ctx, cf := context.WithTimeout(context.Background(), c.cfg.ContextTimeout)
	defer cf()
	req := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(req)
	res := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(res)
	req.Fill(c.GetRequestID(), rony.C_GetNodes, &rony.GetNodes{})
	sessionReplica, err := initConn.send(ctx, req, res, true, requestRetry, requestTimeout)
	if err != nil {
		return err
	}
	c.sessionReplica = sessionReplica
	switch res.Constructor {
	case rony.C_NodeInfoMany:
		x := &rony.NodeInfoMany{}
		_ = x.Unmarshal(res.Message)
		found := false
		for _, n := range x.Nodes {
			var wsc *wsConn
			if !found {
				for _, hp := range n.HostPorts {
					if hp == initConn.hostPorts[0] {
						wsc = initConn
						wsc.hostPorts = n.HostPorts
						found = true
					}
				}
			} else {
				wsc = c.newConn(n.ServerID, n.ReplicaSet, n.HostPorts...)
			}

			c.pool.addConn(n.ServerID, n.ReplicaSet, n.Leader, wsc)
			if n.Leader {
				c.sessionReplica = n.ReplicaSet
			}
		}
	default:
		return ErrUnknownResponse

	}

	return nil
}

func (c *Websocket) Send(req, res *rony.MessageEnvelope, leaderOnly bool) (err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.cfg.ContextTimeout)
	err = c.SendWithDetails(ctx, req, res, true, c.cfg.RequestMaxRetry, c.cfg.RequestTimeout, leaderOnly)
	cancelFunc()
	return
}

func (c *Websocket) SendWithContext(ctx context.Context, req, res *rony.MessageEnvelope, leaderOnly bool) (err error) {
	err = c.SendWithDetails(ctx, req, res, true, c.cfg.RequestMaxRetry, c.cfg.RequestTimeout, leaderOnly)
	return
}

func (c *Websocket) SendWithDetails(
	ctx context.Context, req, res *rony.MessageEnvelope,
	waitToConnect bool, retry int, timeout time.Duration,
	leaderOnly bool,
) (err error) {
	rs := c.cfg.Router.GetRoute(req)
	wsc := c.pool.getConn(rs, leaderOnly)
	if ce := log.Check(log.DebugLevel, "SendWithDetails"); ce != nil {
		ce.Write(
			zap.Uint64("ReqID", req.RequestID),
			zap.Uint64("RS", rs),
			zap.Bool("LeaderOnly", leaderOnly),
		)
	}

	if wsc == nil {
		return ErrNoConnection
	}

SendLoop:
	rs, err = wsc.send(ctx, req, res, waitToConnect, retry, timeout)
	switch err {
	case nil:
		return nil
	case ErrReplicaMaster:
		leaderOnly = true
		wsc = c.pool.getConn(rs, leaderOnly)
	case ErrReplicaSetSession, ErrReplicaSetRequest:
		rs = c.sessionReplica
	}
	goto SendLoop
}

func (c *Websocket) Close() error {

	// by setting the read deadline we make the receiver() routine stops
	// TODO:: close all the connections
	return nil
}
