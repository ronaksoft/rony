package edgeClient

import (
	"context"
	"fmt"
	"git.ronaksoftware.com/ronak/rony"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/pools"
	"git.ronaksoftware.com/ronak/rony/tools"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"go.uber.org/zap"
	"net"
	"sync"
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
	h           MessageHandler
	pendingMtx  sync.RWMutex
	pending     map[uint64]chan *rony.MessageEnvelope
}

func SetLogLevel(level log.Level) {
	log.SetLevel(level)
}

func NewWebsocket(config Config) Client {
	c := Websocket{}
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

	c.connect()
	return &c
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
ConnectLoop:
	c.createDialer(c.dialTimeout)
	urlPrefix := "ws://"
	if c.secure {
		urlPrefix = "wss://"
	}
	conn, _, _, err := c.dialer.Dial(context.Background(), fmt.Sprintf("%s%s", urlPrefix, c.hostPort))
	if err != nil {
		log.Debug("Dial failed", zap.Error(err), zap.String("Host", c.hostPort))
		time.Sleep(time.Duration(tools.RandomInt64(2000))*time.Millisecond + time.Second)
		goto ConnectLoop
	}
	c.conn = conn
	go c.receiver()
	return
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
				_ = e.Unmarshal(ms[idx].Payload)
				c.extractor(e)
				rony.PoolMessageEnvelope.Put(e)
			default:
			}
		}
	}
}

func (c *Websocket) extractor(e *rony.MessageEnvelope) {

	switch e.Constructor {
	case rony.C_MessageContainer:
		x := rony.PoolMessageContainer.Get()
		_ = x.Unmarshal(e.Message)
		for idx := range x.Envelopes {
			c.handler(x.Envelopes[idx])
		}
	default:
		c.handler(e)
	}
}
func (c *Websocket) handler(e *rony.MessageEnvelope) {
	c.pendingMtx.Lock()
	h := c.pending[e.RequestID]
	delete(c.pending, e.RequestID)
	c.pendingMtx.Unlock()

	if h != nil {
		h <- e.Clone()
	} else {
		c.h(e.Clone())
	}
}

func (c *Websocket) Send(req, res *rony.MessageEnvelope) error {
	b := pools.Bytes.GetLen(req.Size())
	_, err := req.MarshalToSizedBuffer(b)
	if err != nil {
		return err
	}

	resChan := make(chan *rony.MessageEnvelope, 1)
	c.pendingMtx.Lock()
	c.pending[req.RequestID] = resChan
	c.pendingMtx.Unlock()
	err = wsutil.WriteClientMessage(c.conn, ws.OpBinary, b)
	pools.Bytes.Put(b)
	if err != nil {
		c.pendingMtx.Lock()
		delete(c.pending, req.RequestID)
		c.pendingMtx.Unlock()
		return err
	}

	select {
	case e := <-resChan:
		switch e.Constructor {
		case rony.C_Redirect:
			x := &rony.Redirect{}
			_ = x.Unmarshal(e.Message)
			c.hostPort = x.LeaderHostPort[0]
			_ = c.conn.SetReadDeadline(time.Now())
			return ErrLeaderRedirect
		}
		e.CopyTo(res)
	case <-time.After(time.Second * 5):
		c.pendingMtx.Lock()
		delete(c.pending, req.RequestID)
		c.pendingMtx.Unlock()
		return ErrTimeout
	}
	return err
}

func (c *Websocket) Close() error {
	c.stop = true
	return c.conn.Close()
}
