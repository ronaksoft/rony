package edgec

import (
	"context"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/ronaksoft/rony"
	wsutil "github.com/ronaksoft/rony/internal/gateway/tcp/util"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"net"
	"strings"
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
	serverID   string
	ws         *Websocket
	stop       bool
	conn       net.Conn
	dialer     ws.Dialer
	connected  bool
	hostPorts  []string
	secure     bool
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

			return nil, ErrNoConnection
		},
		OnStatusError: nil,
		OnHeader:      nil,
		TLSClient:     nil,
		TLSConfig:     nil,
		WrapConn:      nil,
	}
}

func (c *wsConn) connect() {
	urlPrefix := "ws://"
	if c.secure {
		urlPrefix = "wss://"
	}
ConnectLoop:
	log.Debug("Connect", zap.Strings("H", c.hostPorts))
	c.createDialer(c.ws.cfg.DialTimeout)

	sb := strings.Builder{}
	if hf := c.ws.cfg.HeaderFunc; hf != nil {
		for k, v := range hf() {
			sb.WriteString(k)
			sb.WriteString(": ")
			sb.WriteString(v)
			sb.WriteRune('\n')
		}
	}

	c.dialer.Header = ws.HandshakeHeaderString(sb.String())
	conn, _, _, err := c.dialer.Dial(context.Background(), fmt.Sprintf("%s%s", urlPrefix, c.hostPorts[0]))
	if err != nil {
		log.Debug("Dial failed", zap.Error(err), zap.Strings("Host", c.hostPorts))
		time.Sleep(time.Duration(tools.RandomInt64(2000))*time.Millisecond + time.Second)

		goto ConnectLoop
	}
	c.conn = conn
	c.connected = true

	go c.receiver()
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
				c.connected = false
				c.connect()
			}

			break
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
			defaultHandler(e)
		}
		return
	}

	c.ws.pendingMtx.Lock()
	ch := c.ws.pending[e.GetRequestID()]
	delete(c.ws.pending, e.GetRequestID())
	c.ws.pendingMtx.Unlock()

	if ch != nil {
		ch <- e.Clone()
	} else {
		defaultHandler(e)
	}
}

func (c *wsConn) close() error {
	// by setting the stop flag, we are making sure no reconnection will happen
	c.stop = true

	if c.conn == nil {
		return nil
	}

	_ = wsutil.WriteMessage(c.conn, ws.StateClientSide, ws.OpClose, nil)

	// by setting the read deadline we make the receiver() routine stops
	return c.conn.SetReadDeadline(time.Now())
}

func (c *wsConn) send(req *rony.MessageEnvelope) error {
	if !c.connected {
		c.connect()
	}
	mo := proto.MarshalOptions{UseCachedSize: true}
	buf := pools.Buffer.GetCap(mo.Size(req))
	defer pools.Buffer.Put(buf)

	b, err := mo.MarshalAppend(*buf.Bytes(), req)
	if err != nil {
		return err
	}

	return wsutil.WriteMessage(c.conn, ws.StateClientSide, ws.OpBinary, b)
}
