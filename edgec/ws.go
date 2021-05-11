package edgec

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
	"go.uber.org/zap"
	"strings"
	"sync"
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
	// Handler must not block in function because other incoming messages might get blocked.
	// This handler must returns quickly and pass a deep copy of the MessageEnvelope to other
	// routines.
	Handler    MessageHandler
	HeaderFunc func() map[string]string
	Secure     bool
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
	sessionReplica uint64
	nextReqID      uint64

	// Connection Pool
	connsMtx       sync.RWMutex
	connsByReplica map[uint64]map[string]*wsConn
	connsByID      map[string]*wsConn
	leaderIDs      map[uint64]string

	// FLusher
	flusherPool *tools.FlusherPool

	// Flying Requests
	pendingMtx tools.SpinLock
	pending    map[uint64]chan *rony.MessageEnvelope
}

func NewWebsocket(config WebsocketConfig) *Websocket {
	c := &Websocket{
		nextReqID:      tools.RandomUint64(0),
		cfg:            config,
		connsByReplica: make(map[uint64]map[string]*wsConn, 64),
		connsByID:      make(map[string]*wsConn, 64),
		leaderIDs:      make(map[uint64]string, 64),
		pending:        make(map[uint64]chan *rony.MessageEnvelope, 1024),
	}

	c.flusherPool = tools.NewFlusherPool(1, 100, c.sendFunc)

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
		c.cfg.Router = &wsRouter{
			c: c,
		}
	}
	if c.cfg.Handler == nil {
		c.cfg.Handler = func(m *rony.MessageEnvelope) {}
	}

	return c
}

func (ws *Websocket) GetRequestID() uint64 {
	return atomic.AddUint64(&ws.nextReqID, 1)
}

func (ws *Websocket) Start() error {
	initConn := ws.newConn("", 0, ws.cfg.SeedHostPort)
	ws.addConn("", 0, true, initConn)

	req := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(req)
	res := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(res)
	req.Fill(ws.GetRequestID(), rony.C_GetNodes, &rony.GetNodes{})
	err := ws.Send(req, res, false)
	if err != nil {
		return err
	}
	switch res.Constructor {
	case rony.C_Edges:
		x := &rony.Edges{}
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
			wsc := ws.newConn(n.ServerID, n.ReplicaSet, n.HostPorts...)
			if !found {
				for _, hp := range n.HostPorts {
					if hp == initConn.hostPorts[0] {
						ws.removeConn("", 0)
						found = true
					}
				}
			}

			ws.addConn(n.ServerID, n.ReplicaSet, n.Leader, wsc)
			if n.Leader {
				ws.sessionReplica = n.ReplicaSet
			}
		}

		// If this connection is not our connsByReplica then we just close it.
		if !found {
			_ = initConn.close()
		}
	default:
		return ErrUnknownResponse

	}
	return nil
}

func (ws *Websocket) addConn(serverID string, replicaSet uint64, leader bool, wsc *wsConn) {
	log.Debug("Pool connection added",
		zap.String("ServerID", serverID),
		zap.Uint64("RS", replicaSet),
		zap.Bool("Leader", leader),
	)
	ws.connsMtx.Lock()
	defer ws.connsMtx.Unlock()

	if ws.connsByReplica[replicaSet] == nil {
		ws.connsByReplica[replicaSet] = make(map[string]*wsConn, 16)
	}
	ws.connsByID[serverID] = wsc
	ws.connsByReplica[replicaSet][serverID] = wsc
	if leader || replicaSet == 0 {
		ws.leaderIDs[replicaSet] = serverID
	}
}

func (ws *Websocket) removeConn(serverID string, replicaSet uint64) {
	ws.connsMtx.Lock()
	defer ws.connsMtx.Unlock()

	if ws.connsByReplica[replicaSet] != nil {
		delete(ws.connsByReplica[replicaSet], serverID)
	}
	delete(ws.connsByID, serverID)
}

func (ws *Websocket) getConnByReplica(replicaSet uint64, onlyLeader bool) *wsConn {
	ws.connsMtx.RLock()
	defer ws.connsMtx.RUnlock()

	if onlyLeader {
		leaderID := ws.leaderIDs[replicaSet]
		if leaderID == "" {
			return nil
		}
		m := ws.connsByReplica[replicaSet]
		if m != nil {
			c := m[leaderID]
			return c
		}
	} else {
		m := ws.connsByReplica[replicaSet]
		for _, c := range m {
			return c
		}
	}
	return nil
}

func (ws *Websocket) getConnByID(serverID string) *wsConn {
	ws.connsMtx.RLock()
	defer ws.connsMtx.RUnlock()

	wsc := ws.connsByID[serverID]
	return wsc
}

func (ws *Websocket) newConn(id string, replicaSet uint64, hostPorts ...string) *wsConn {
	wsc := &wsConn{
		serverID:   id,
		ws:         ws,
		replicaSet: replicaSet,
		hostPorts:  hostPorts,
	}
	return wsc
}

func (ws *Websocket) sendFunc(serverID string, entries []tools.FlushEntry) {
	wsc := ws.getConnByID(serverID)
	if wsc == nil {
		// TODO:: for each entry we must return
		return
	}

	// Check if we have active connection
	if !wsc.connected {
		wsc.connect()
	}

	me := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(me)

	switch len(entries) {
	case 0:
		// There is nothing to do, Probably a bug if we are here
		return
	case 1:
		ev := entries[0].Value().(*wsRequest)
		ws.pendingMtx.Lock()
		ws.pending[ev.req.GetRequestID()] = ev.resChan
		ws.pendingMtx.Unlock()
		ev.req.DeepCopy(me)

	default:
		mc := rony.PoolMessageContainer.Get()
		for _, e := range entries {
			ev := e.Value().(*wsRequest)
			ws.pendingMtx.Lock()
			ws.pending[ev.req.GetRequestID()] = ev.resChan
			ws.pendingMtx.Unlock()
			mc.Envelopes = append(mc.Envelopes, ev.req.Clone())
			mc.Length += 1
		}
		me.Fill(0, rony.C_MessageContainer, mc)
		rony.PoolMessageContainer.Put(mc)
	}
	err := wsc.send(me)
	if err != nil {
		log.Warn("EdgeClient (Websocket) got error on sending request", zap.Error(err))
	}

}

func (ws *Websocket) Send(req, res *rony.MessageEnvelope, leaderOnly bool) (err error) {
	err = ws.SendWithDetails(req, res, ws.cfg.RequestMaxRetry, ws.cfg.RequestTimeout, leaderOnly)
	return
}

func (ws *Websocket) SendWithDetails(
	req, res *rony.MessageEnvelope, retry int, timeout time.Duration, leaderOnly bool,
) error {
	rs := ws.cfg.Router.GetRoute(req)

Loop:
	wsc := ws.getConnByReplica(rs, leaderOnly)

	if wsc == nil {
		// TODO:: try to gather information about the target
		return ErrNoConnection
	}

	wsReq := &wsRequest{
		req:     req,
		resChan: make(chan *rony.MessageEnvelope, 1),
	}
	ws.flusherPool.Enter(wsc.serverID, tools.NewEntry(wsReq))

	t := pools.AcquireTimer(timeout)
	defer pools.ReleaseTimer(t)
	select {
	case x := <-wsReq.resChan:
		switch x.GetConstructor() {
		case rony.C_Redirect:
			xx := &rony.Redirect{}
			_ = xx.Unmarshal(x.GetMessage())
			rs, leaderOnly = ws.redirect(xx, leaderOnly)
			if retry--; retry < 0 {
				return rony.ErrRetriesExceeded(fmt.Errorf("redirect"))
			}
			goto Loop
		default:
			x.DeepCopy(res)
		}
	case <-t.C:
		ws.pendingMtx.Lock()
		delete(ws.pending, req.GetRequestID())
		ws.pendingMtx.Unlock()
		return ErrTimeout
	}

	return nil
}

func (ws *Websocket) redirect(x *rony.Redirect, leaderOnly bool) (replicaSet uint64, switchToLeader bool) {
	if ce := log.Check(log.InfoLevel, "EdgeClient (Websocket) received Redirect"); ce != nil {
		ce.Write(
			zap.Any("Leader", x.Leader),
			zap.Any("Followers", x.Followers),
			zap.Any("Wait", x.WaitInSec),
		)
	}

	ws.addConn(
		x.Leader.ServerID, x.Leader.ReplicaSet, true,
		ws.newConn(x.Leader.ServerID, x.Leader.ReplicaSet, x.Leader.HostPorts...),
	)
	for _, n := range x.Followers {
		ws.addConn(
			n.ServerID, n.ReplicaSet, false,
			ws.newConn(n.ServerID, n.ReplicaSet, n.HostPorts...),
		)
	}

	switchToLeader = leaderOnly
	replicaSet = x.Leader.ReplicaSet
	switch x.Reason {
	case rony.RedirectReason_ReplicaMaster:
		switchToLeader = true
	case rony.RedirectReason_ReplicaSetSession:
		ws.sessionReplica = x.Leader.ReplicaSet
	case rony.RedirectReason_ReplicaSetRequest:
	default:
	}

	return
}

func (ws *Websocket) Close() error {
	ws.connsMtx.RLock()
	defer ws.connsMtx.RUnlock()

	for _, conns := range ws.connsByReplica {
		for _, c := range conns {
			_ = c.close()
		}
	}
	return nil
}

func (ws *Websocket) ConnInfo() string {
	sb := strings.Builder{}
	sb.WriteString("\n-----\n")
	ws.connsMtx.Lock()
	for id, wsc := range ws.connsByID {
		sb.WriteString(fmt.Sprintf("%s: [RS=%d] [HostPorts=%v] [Connected: %t]\n", id, wsc.replicaSet, wsc.hostPorts, wsc.connected))
	}
	ws.connsMtx.Unlock()
	sb.WriteString("-----\n")
	return sb.String()
}

func (ws *Websocket) ClusterInfo(replicaSets ...uint64) (*rony.Edges, error) {
	req := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(req)
	res := rony.PoolMessageEnvelope.Get()
	defer rony.PoolMessageEnvelope.Put(res)
	req.Fill(ws.GetRequestID(), rony.C_GetNodes, &rony.GetNodes{ReplicaSet: replicaSets})
	err := ws.Send(req, res, false)
	if err != nil {
		return nil, err
	}
	switch res.GetConstructor() {
	case rony.C_Edges:
		x := &rony.Edges{}
		_ = x.Unmarshal(res.GetMessage())
		return x, nil
	case rony.C_Error:
		x := &rony.Error{}
		_ = x.Unmarshal(res.GetMessage())
		return nil, x
	}
	return nil, ErrUnknownResponse
}

type wsRouter struct {
	c *Websocket
}

func (d *wsRouter) UpdateRoute(req *rony.MessageEnvelope, replicaSet uint64) {
	// TODO:: implement cache maybe
}

func (d *wsRouter) GetRoute(req *rony.MessageEnvelope) (replicaSet uint64) {
	return d.c.sessionReplica
}

type wsRequest struct {
	req     *rony.MessageEnvelope
	resChan chan *rony.MessageEnvelope
}
