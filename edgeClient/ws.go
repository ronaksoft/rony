package edgeClient

import (
	"context"
	"fmt"
	"git.ronaksoft.com/ronak/rony"
	log "git.ronaksoft.com/ronak/rony/internal/logger"
	"git.ronaksoft.com/ronak/rony/pools"
	"git.ronaksoft.com/ronak/rony/tools"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"net"
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
)

type MessageHandler func(m *rony.MessageEnvelope)
type Config struct {
	HostPort    string
	IdleTime    time.Duration
	DialTimeout time.Duration
	Handler     MessageHandler
	Secure      bool
}
type Websocket struct {
	hostPort    string
	secure      bool
	idleTimeout time.Duration
	dialTimeout time.Duration
	dialer      ws.Dialer
	conn        net.Conn
	stop        bool
	connected   bool
	connectMtx  sync.Mutex

	// parameters related to handling request/responses
	h          MessageHandler
	pendingMtx sync.RWMutex
	pending    map[uint64]chan *rony.MessageEnvelope
	nextReqID  uint64
}

func SetLogLevel(level log.Level) {
	log.SetLevel(level)
}

func NewWebsocket(config Config) Client {
	c := Websocket{
		nextReqID: tools.RandomUint64(),
	}
	c.hostPort = config.HostPort
	c.h = config.Handler
	if config.DialTimeout == 0 {
		config.DialTimeout = time.Second * 3
	}
	if config.IdleTime == 0 {
		config.IdleTime = time.Minute
	}
	c.idleTimeout = config.IdleTime
	c.dialTimeout = config.DialTimeout
	c.pending = make(map[uint64]chan *rony.MessageEnvelope, 100)

	// start the connection
	c.connect()
	return &c
}

func (c *Websocket) GetRequestID() uint64 {
	return atomic.AddUint64(&c.nextReqID, 1)
}

func (c *Websocket) createDialer(timeout time.Duration) {
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

func (c *Websocket) connect() {
	urlPrefix := "ws://"
	if c.secure {
		urlPrefix = "wss://"
	}
ConnectLoop:
	log.Debug("Connect", zap.String("H", c.hostPort))
	c.createDialer(c.dialTimeout)
	conn, _, _, err := c.dialer.Dial(context.Background(), fmt.Sprintf("%s%s", urlPrefix, c.hostPort))
	if err != nil {
		log.Debug("Dial failed", zap.Error(err), zap.String("Host", c.hostPort))
		time.Sleep(time.Duration(tools.RandomInt64(2000))*time.Millisecond + time.Second)
		goto ConnectLoop
	}
	c.conn = conn
	c.connected = true
	go c.receiver()
	return
}

func (c *Websocket) reconnect() {
	c.connected = false
	_ = c.conn.SetReadDeadline(time.Now())
}

func (c *Websocket) waitUntilConnect() {
	for !c.stop {
		if c.connected {
			break
		}
		time.Sleep(time.Millisecond)
	}
}

func (c *Websocket) receiver() {
	var (
		ms []wsutil.Message
	)
	// Receive Loop
	for {
		ms = ms[:0]
		_ = c.conn.SetReadDeadline(time.Now().Add(c.idleTimeout))
		ms, err := wsutil.ReadServerMessage(c.conn, ms)
		if err != nil {
			_ = c.conn.Close()
			if !c.stop {
				c.connect()
			}
			return
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

func (c *Websocket) extractor(e *rony.MessageEnvelope) {
	switch e.GetConstructor() {
	case rony.C_MessageContainer:
		x := rony.PoolMessageContainer.Get()
		uo := proto.UnmarshalOptions{
			Merge: true,
		}
		_ = uo.Unmarshal(e.Message, x)
		for idx := range x.Envelopes {
			c.handler(x.Envelopes[idx])
		}
		rony.PoolMessageContainer.Put(x)
	default:
		c.handler(e)
	}
}

func (c *Websocket) handler(e *rony.MessageEnvelope) {
	if e.GetRequestID() == 0 {
		c.h(e.Clone())
		return
	}

	c.pendingMtx.Lock()
	h := c.pending[e.GetRequestID()]
	delete(c.pending, e.GetRequestID())
	c.pendingMtx.Unlock()
	if h != nil {
		h <- e.Clone()
	} else {
		c.h(e.Clone())
	}
}

func (c *Websocket) Send(req, res *rony.MessageEnvelope) (err error) {
	return c.SendWithRetry(req, res, requestRetry, requestTimeout)
}

func (c *Websocket) SendWithRetry(req, res *rony.MessageEnvelope, retry int, timeout time.Duration) (err error) {
	mo := proto.MarshalOptions{UseCachedSize: true}
	b := pools.Bytes.GetCap(mo.Size(req))
	defer pools.Bytes.Put(b)

	b, err = mo.MarshalAppend(b, req)
	if err != nil {
		return err
	}

	t := pools.AcquireTimer(timeout)
	defer pools.ReleaseTimer(t)

SendLoop:
	if retry--; retry < 0 {
		return
	}
	c.waitUntilConnect()
	resChan := make(chan *rony.MessageEnvelope, 1)
	c.pendingMtx.Lock()
	c.pending[req.GetRequestID()] = resChan
	c.pendingMtx.Unlock()
	err = wsutil.WriteClientMessage(c.conn, ws.OpBinary, b)
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
			uo := proto.UnmarshalOptions{
				Merge: true,
			}
			_ = uo.Unmarshal(e.Message, x)
			c.connectMtx.Lock()
			if c.hostPort != x.LeaderHostPort[0] {
				c.hostPort = x.LeaderHostPort[0]
				c.reconnect()
			}
			c.connectMtx.Unlock()
			err = ErrLeaderRedirect
			goto SendLoop
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

func (c *Websocket) Close() error {
	// by setting the stop flag, we are making sure no reconnection will happen
	c.stop = true

	// by setting the read deadline we make the receiver() routine stops
	return c.conn.SetReadDeadline(time.Now())
}
