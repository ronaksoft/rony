package edgeClient

import (
	"context"
	"fmt"
	"git.ronaksoftware.com/ronak/rony"
	log "git.ronaksoftware.com/ronak/rony/internal/logger"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
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
type Client struct {
	hostPort    string
	idleTimeout time.Duration
	dialTimeout time.Duration
	dialer      ws.Dialer
	conn        net.Conn
	stop        bool
	stopChan    chan struct{}
	h           MessageHandler
	pendingMtx  sync.RWMutex
	pending     map[uint64]chan *rony.MessageEnvelope
}

func SetLogLevel(level log.Level) {
	log.SetLevel(level)
}

func NewWebsocket(hostPort string, dialTimeout time.Duration, h MessageHandler) *Client {
	c := Client{}
	c.stopChan = make(chan struct{}, 1)
	c.dialTimeout = dialTimeout
	c.hostPort = hostPort
	c.h = h
	c.pending = make(map[uint64]chan *rony.MessageEnvelope, 100)
	return &c
}

func (c *Client) createDialer(timeout time.Duration) {
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
			log.Info("DNS LookIP", zap.String("Addr", addr), zap.Any("IPs", ips))
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
func (c *Client) Connect() {
ConnectLoop:
	c.createDialer(c.dialTimeout)
	conn, _, _, err := c.dialer.Dial(context.Background(), c.hostPort)
	if err != nil {
		log.Debug("Dial failed", zap.Error(err), zap.String("Host", c.hostPort))
		time.Sleep(time.Duration(tools.RandomInt64(2000)) * time.Millisecond)
		goto ConnectLoop
	}
	c.conn = conn
	go c.receiver()
	return
}

func (c *Client) receiver() {
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
			c.Connect()
			continue
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

func (c *Client) extractor(e *rony.MessageEnvelope) {
	switch e.Constructor {
	case rony.C_MessageContainer:
		x := rony.PoolMessageContainer.Get()
		_ = x.Unmarshal(e.Message)
		for idx := range x.Envelopes {
			c.pendingMtx.Lock()
			ch := c.pending[x.Envelopes[idx].RequestID]
			delete(c.pending, x.Envelopes[idx].RequestID)
			c.pendingMtx.Unlock()
			if ch != nil {
				ch <- x.Envelopes[idx].Clone()
			} else {
				c.h(x.Envelopes[idx].Clone())
			}
		}
	default:
		c.pendingMtx.Lock()
		ch := c.pending[e.RequestID]
		delete(c.pending, e.RequestID)
		c.pendingMtx.Unlock()
		if ch != nil {
			ch <- e.Clone()
		} else {
			c.h(e.Clone())
		}
	}
}

func (c *Client) Send(constructor int64, m rony.ProtoBufferMessage) (<-chan *rony.MessageEnvelope, error) {
	reqID := tools.RandomUint64()
	e := rony.PoolMessageEnvelope.Get()
	e.Fill(reqID, constructor, m)
	b := pools.Bytes.GetLen(e.Size())
	_, err := e.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}

	resChan := make(chan *rony.MessageEnvelope)
	c.pendingMtx.Lock()
	c.pending[reqID] = resChan
	c.pendingMtx.Unlock()
	err = wsutil.WriteClientMessage(c.conn, ws.OpBinary, b)
	pools.Bytes.Put(b)
	return resChan, err
}
