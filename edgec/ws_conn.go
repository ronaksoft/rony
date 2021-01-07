package edgec

import (
	"context"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/ronaksoft/rony"
	wsutil "github.com/ronaksoft/rony/gateway/tcp/util"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"net"
	"strings"
	"sync"
	"time"
)

/*
   Creation Time: 2021 - Jan - 04
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type wsConn struct {
	replicaSet uint64
	id         string
	stop       bool
	ws         *Websocket
	conn       net.Conn
	dialer     ws.Dialer
	connected  bool
	mtx        sync.Mutex
	hostPorts  []string
	secure     bool
	pendingMtx sync.RWMutex
	pending    map[uint64]chan *rony.MessageEnvelope
}

func (c *wsConn) createDialer(timeout time.Duration) {
	c.dialer = ws.Dialer{
		ReadBufferSize:  32 * 1024, // 32kB
		WriteBufferSize: 32 * 1024, // 32kB
		Timeout:         timeout,
		NetDial: func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			ips, err := net.LookupIP(host)
			if err != nil {
				return nil, err
			}
			log.Debug("DNS LookIP", zap.String("Addr", addr), zap.Any("IPs", ips))
			d := net.Dialer{Timeout: timeout}
			for _, ip := range ips {
				if ip.To4() != nil {
					conn, err = d.DialContext(ctx, "tcp4", net.JoinHostPort(ip.String(), port))
					if err != nil {
						continue
					}
					return
				}
			}
			return nil, fmt.Errorf("no connection")
		},
		OnStatusError: nil,
		OnHeader:      nil,
		TLSClient:     nil,
		TLSConfig:     nil,
		WrapConn:      nil,
	}
}

func (c *wsConn) connect() {
	if c.isConnected() {
		return
	}
	urlPrefix := "ws://"
	if c.secure {
		urlPrefix = "wss://"
	}
ConnectLoop:
	log.Debug("Connect", zap.Strings("H", c.hostPorts))
	c.createDialer(c.ws.cfg.DialTimeout)

	sb := strings.Builder{}
	for k, v := range c.ws.cfg.Header {
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(v)
		sb.WriteRune('\n')
	}
	c.dialer.Header = ws.HandshakeHeaderString(sb.String())
	conn, _, _, err := c.dialer.Dial(context.Background(), fmt.Sprintf("%s%s", urlPrefix, c.hostPorts[0]))
	if err != nil {
		log.Debug("Dial failed", zap.Error(err), zap.Strings("Host", c.hostPorts))
		time.Sleep(time.Duration(tools.RandomInt64(2000))*time.Millisecond + time.Second)
		goto ConnectLoop
	}
	c.conn = conn
	c.mtx.Lock()
	c.connected = true
	c.mtx.Unlock()

	go c.receiver()
	return
}

func (c *wsConn) isConnected() bool {
	c.mtx.Lock()
	b := c.connected
	c.mtx.Unlock()
	return b
}

func (c *wsConn) reconnect() {
	c.mtx.Lock()
	c.connected = false
	c.mtx.Unlock()
	_ = c.conn.SetReadDeadline(time.Now())
}

func (c *wsConn) waitUntilConnect(ctx context.Context) error {
	step := time.Duration(10)
	for !c.stop {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if c.isConnected() {
			break
		}
		time.Sleep(time.Millisecond * step)
		if step < 1000 {
			step += 10
		}
	}
	return nil
}

func (c *wsConn) receiver() {
	var (
		ms []wsutil.Message
	)
	// Receive Loop
	for {
		ms = ms[:0]
		_ = c.conn.SetReadDeadline(time.Now().Add(c.ws.cfg.IdleTimeout))
		ms, err := wsutil.ReadMessage(c.conn, ws.StateClientSide, ms)
		if err != nil {
			_ = c.conn.Close()
			if !c.stop {
				c.mtx.Lock()
				c.connected = false
				c.mtx.Unlock()
				c.connect()
			}
			break
		}
		for idx := range ms {
			switch ms[idx].OpCode {
			case ws.OpBinary, ws.OpText:
				e := rony.PoolMessageEnvelope.Get()
				uo := proto.UnmarshalOptions{}
				_ = uo.Unmarshal(ms[idx].Payload, e)
				c.extractor(e)
				rony.PoolMessageEnvelope.Put(e)
			default:
			}
		}
	}
}

func (c *wsConn) extractor(e *rony.MessageEnvelope) {
	switch e.GetConstructor() {
	case rony.C_MessageContainer:
		x := rony.PoolMessageContainer.Get()
		_ = x.Unmarshal(e.Message)
		for idx := range x.Envelopes {
			c.handler(x.Envelopes[idx])
		}
		rony.PoolMessageContainer.Put(x)
	default:
		c.handler(e)
	}
}

func (c *wsConn) handler(e *rony.MessageEnvelope) {
	defaultHandler := c.ws.cfg.Handler
	if e.GetRequestID() == 0 {
		if defaultHandler != nil {
			defaultHandler(e.Clone())
		}
		return
	}

	c.pendingMtx.Lock()
	h := c.pending[e.GetRequestID()]
	delete(c.pending, e.GetRequestID())
	c.pendingMtx.Unlock()
	if h != nil {
		h <- e.Clone()
	} else {
		defaultHandler(e.Clone())
	}
}

func (c *wsConn) close() error {
	// by setting the stop flag, we are making sure no reconnection will happen
	c.stop = true

	// by setting the read deadline we make the receiver() routine stops
	return c.conn.SetReadDeadline(time.Now())
}

func (c *wsConn) send(ctx context.Context, req, res *rony.MessageEnvelope, waitToConnect bool, retry int, timeout time.Duration) (replicaSet uint64, err error) {
	replicaSet = c.replicaSet
	mo := proto.MarshalOptions{UseCachedSize: true}
	b := pools.Bytes.GetCap(mo.Size(req))
	defer pools.Bytes.Put(b)

	b, err = mo.MarshalAppend(b, req)
	if err != nil {
		return
	}

	t := pools.AcquireTimer(timeout)
	defer pools.ReleaseTimer(t)

SendLoop:
	// Check if context is canceled on each loop
	select {
	case <-ctx.Done():
		err = ctx.Err()
		return
	default:
	}

	// If we exceeds the maximum retry then we return
	if retry--; retry < 0 {
		err = rony.ErrRetriesExceeded(err)
		return
	}

	// If it is required to wait until the connection is established before try sending the
	// request over the wire
	if waitToConnect {
		err = c.waitUntilConnect(ctx)
		if err != nil {
			return
		}
	}

	resChan := make(chan *rony.MessageEnvelope, 1)
	c.pendingMtx.Lock()
	c.pending[req.GetRequestID()] = resChan
	c.pendingMtx.Unlock()
	c.mtx.Lock()
	err = wsutil.WriteMessage(c.conn, ws.StateClientSide, ws.OpBinary, b)
	c.mtx.Unlock()
	if err != nil {
		c.pendingMtx.Lock()
		delete(c.pending, req.GetRequestID())
		c.pendingMtx.Unlock()
		goto SendLoop
	}
	pools.ResetTimer(t, timeout)

	select {
	case e := <-resChan:
		switch e.GetConstructor() {
		case rony.C_Redirect:
			x := &rony.Redirect{}
			err = proto.Unmarshal(e.Message, x)
			if err != nil {
				log.Warn("Error On Unmarshal Redirect", zap.Error(err))
				goto SendLoop
			}
			replicaSet, err = c.redirect(x)
			return
		}
		e.DeepCopy(res)
	case <-t.C:
		c.pendingMtx.Lock()
		delete(c.pending, req.GetRequestID())
		c.pendingMtx.Unlock()
		err = ErrTimeout
		goto SendLoop
	}
	return
}

func (c *wsConn) redirect(x *rony.Redirect) (replicaSet uint64, err error) {
	replicaSet = c.replicaSet
	if ce := log.Check(log.InfoLevel, "Redirect"); ce != nil {
		ce.Write(
			zap.Any("Leader", x.Leader),
			zap.Any("Followers", x.Followers),
			zap.Any("Wait", x.WaitInSec),
		)
	}

	c.ws.pool.addConn(
		x.Leader.ServerID, x.Leader.ReplicaSet, true,
		c.ws.newConn(x.Leader.ServerID, x.Leader.ReplicaSet, x.Leader.HostPorts...),
	)
	replicaSet = x.Leader.ReplicaSet
	for _, n := range x.Followers {
		c.ws.pool.addConn(
			n.ServerID, n.ReplicaSet, false,
			c.ws.newConn(n.ServerID, n.ReplicaSet, n.HostPorts...),
		)
	}

	switch x.Reason {
	case rony.RedirectReason_ReplicaMaster:
		err = ErrReplicaMaster
	case rony.RedirectReason_ReplicaSetSession:
		c.ws.sessionReplica = replicaSet
		err = ErrReplicaSetSession
	case rony.RedirectReason_ReplicaSetRequest:
		replicaSet = x.Leader.ReplicaSet
		err = ErrReplicaSetRequest
	default:
		err = ErrUnknownResponse
	}

	return
}

type connPool struct {
	mtx       sync.RWMutex
	pool      map[uint64]map[string]*wsConn
	leaderIDs map[uint64]string
}

func newConnPool() *connPool {
	cp := &connPool{
		pool:      make(map[uint64]map[string]*wsConn, 16),
		leaderIDs: make(map[uint64]string, 16),
	}
	return cp
}

func (cp *connPool) addConn(serverID string, replicaSet uint64, leader bool, c *wsConn) {
	log.Debug("Pool connection added",
		zap.String("ServerID", serverID),
		zap.Uint64("RS", replicaSet),
		zap.Bool("Leader", leader),
	)
	cp.mtx.Lock()
	defer cp.mtx.Unlock()

	if cp.pool[replicaSet] == nil {
		cp.pool[replicaSet] = make(map[string]*wsConn, 16)
	}
	cp.pool[replicaSet][serverID] = c
	if leader || replicaSet == 0 {
		cp.leaderIDs[replicaSet] = serverID
	}
}

func (cp *connPool) removeConn(serverID string, replicaSet uint64) {
	cp.mtx.Lock()
	defer cp.mtx.Unlock()
}

func (cp *connPool) getConn(replicaSet uint64, onlyLeader bool) *wsConn {
	cp.mtx.RLock()
	defer cp.mtx.RUnlock()

	if onlyLeader {
		leaderID := cp.leaderIDs[replicaSet]
		if leaderID == "" {
			return nil
		}
		m := cp.pool[replicaSet]
		if m != nil {
			c := m[leaderID]
			go c.connect()
			return c
		}
	} else {
		m := cp.pool[replicaSet]
		if m != nil {
			for _, c := range m {
				go c.connect()
				return c
			}
		}
	}
	return nil
}
