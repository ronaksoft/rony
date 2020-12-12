package edgec

import (
	"context"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/ronaksoft/rony"
	log "github.com/ronaksoft/rony/internal/logger"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"net"
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
type Config struct {
	HostPort     string
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
}
type Websocket struct {
	hostPort       string
	secure         bool
	idleTimeout    time.Duration
	dialTimeout    time.Duration
	requestTimeout time.Duration
	contextTimeout time.Duration
	requestRetry   int
	dialer         ws.Dialer
	conn           net.Conn
	stop           bool
	connected      bool
	connectMtx     sync.Mutex
	header         map[string]string

	// parameters related to handling request/responses
	h          MessageHandler
	pendingMtx sync.RWMutex
	pending    map[uint64]chan *rony.MessageEnvelope
	nextReqID  uint64
}

func NewWebsocket(config Config) *Websocket {
	if config.DialTimeout == 0 {
		config.DialTimeout = dialTimeout
	}
	if config.IdleTimeout == 0 {
		config.IdleTimeout = idleTimeout
	}
	if config.RequestMaxRetry == 0 {
		config.RequestMaxRetry = requestRetry
	}
	if config.RequestTimeout == 0 {
		config.RequestTimeout = requestTimeout
	}
	if config.ContextTimeout == 0 {
		config.ContextTimeout = config.RequestTimeout * time.Duration(config.RequestMaxRetry)
	}
	c := Websocket{
		nextReqID:      tools.RandomUint64(0),
		idleTimeout:    config.IdleTimeout,
		dialTimeout:    config.DialTimeout,
		requestTimeout: config.RequestTimeout,
		contextTimeout: config.ContextTimeout,
		requestRetry:   config.RequestMaxRetry,
		hostPort:       config.HostPort,
		h:              config.Handler,
		pending:        make(map[uint64]chan *rony.MessageEnvelope, 100),
		header:         config.Header,
	}

	// start the connection
	if config.ForceConnect {
		c.connect()
	} else {
		go c.connect()
	}
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

	sb := strings.Builder{}
	for k, v := range c.header {
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(v)
		sb.WriteRune('\n')
	}
	c.dialer.Header = ws.HandshakeHeaderString(sb.String())
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

func (c *Websocket) waitUntilConnect(ctx context.Context) error {
	for !c.stop {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		if c.connected {
			break
		}
		time.Sleep(time.Millisecond)
	}
	return nil
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
		_ = x.Unmarshal(e.Message)
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
		if c.h != nil {
			c.h(e.Clone())
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
		c.h(e.Clone())
	}
}

func (c *Websocket) Send(req, res *rony.MessageEnvelope) (err error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), c.contextTimeout)
	err = c.SendWithDetails(ctx, req, res, true, c.requestRetry, c.requestTimeout)
	cancelFunc()
	return
}

func (c *Websocket) SendWithContext(ctx context.Context, req, res *rony.MessageEnvelope) (err error) {
	err = c.SendWithDetails(ctx, req, res, true, c.requestRetry, c.requestTimeout)
	return
}

func (c *Websocket) SendWithDetails(ctx context.Context, req, res *rony.MessageEnvelope, waitToConnect bool, retry int, timeout time.Duration) (err error) {
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
	// Check if context is canceled on each loop
	select {
	case <-ctx.Done():
		return ctx.Err()
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
			_ = proto.Unmarshal(e.Message, x)
			c.connectMtx.Lock()
			if len(x.LeaderHostPort) > 0 && c.hostPort != x.LeaderHostPort[0] {
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
