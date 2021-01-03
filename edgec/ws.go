package edgec

import (
	"context"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/tools"
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

type RouterFunc func(m *rony.MessageEnvelope) (replicaSet uint64, forceLeader bool)

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
	Router         RouterFunc
}

// Websocket client which could handle multiple connections
type Websocket struct {
	cfg            WebsocketConfig
	pool           *connPool
	sessionReplica uint64
	nextReqID      uint64
}

func NewWebsocket(config WebsocketConfig) (*Websocket, error) {
	c := Websocket{
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
		config.ContextTimeout = config.RequestTimeout * time.Duration(config.RequestMaxRetry)
	}
	if c.cfg.Router == nil {
		c.cfg.Router = c.defaultRouter
	}

	err := c.initConn()
	if err != nil {
		return nil, err
	}
	// start the connection
	// if config.ForceConnect {
	// 	c.connect()
	// } else {
	// 	go c.connect()
	// }
	return &c, nil
}

func (c *Websocket) GetRequestID() uint64 {
	return atomic.AddUint64(&c.nextReqID, 1)
}

func (c *Websocket) initConn() error {
	wsConn := newConn("", 0, c.cfg.SeedHostPort)
	ctx, cf := context.WithTimeout(context.Background(), c.cfg.ContextTimeout)
	defer cf()
	req := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(req)
	res := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(res)
	req.Fill(c.GetRequestID(), rony.C_GetNodes, &rony.GetNodes{})
	sessionReplica, err := wsConn.send(ctx, req, res, true, requestRetry, requestTimeout)
	if err != nil {
		return err
	}
	c.sessionReplica = sessionReplica
	switch res.Constructor {
	case rony.C_NodeInfoMany:
		x := &rony.NodeInfoMany{}
		_ = x.Unmarshal(res.Message)
		for _, n := range x.Nodes {
			wsc := newConn(n.ServerID, n.ReplicaSet, n.HostPorts...)
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

func (c *Websocket) isLeaderOnly(req *rony.MessageEnvelope) bool {
	return true
}

func (c *Websocket) defaultRouter(req *rony.MessageEnvelope) (replicaSet uint64, leaderOnly bool) {
	return c.sessionReplica, c.isLeaderOnly(req)
}

func (c *Websocket) Send(req, res *rony.MessageEnvelope) (err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.cfg.ContextTimeout)
	err = c.SendWithDetails(ctx, req, res, true, c.cfg.RequestMaxRetry, c.cfg.RequestTimeout)
	cancelFunc()
	return
}

func (c *Websocket) SendWithContext(ctx context.Context, req, res *rony.MessageEnvelope) (err error) {
	err = c.SendWithDetails(ctx, req, res, true, c.cfg.RequestMaxRetry, c.cfg.RequestTimeout)
	return
}

func (c *Websocket) SendWithDetails(ctx context.Context, req, res *rony.MessageEnvelope, waitToConnect bool, retry int, timeout time.Duration) (err error) {
	rs, leader := c.cfg.Router(req)
	wsc := c.pool.getConn(rs, leader)
SendLoop:
	rs, err = wsc.send(ctx, req, res, waitToConnect, retry, timeout)
	switch err {
	case nil:
		return nil
	case ErrReplicaMaster:
		wsc = c.pool.getConn(rs, true)
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
