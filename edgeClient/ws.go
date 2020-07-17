package edgeClient

import (
	"context"
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/internal/pools"
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"net"
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

type MessageHandler func(m rony.MessageEnvelope)
type Client struct {
	hostPort string
	dialer   ws.Dialer
	conn     net.Conn
	stop     bool
	stopChan chan struct{}
	h        MessageHandler
}

func NewWebsocket(hostPort string, dialTimeout time.Duration, h MessageHandler) *Client {
	c := Client{}
	c.stopChan = make(chan struct{}, 1)
	c.dialer = ws.DefaultDialer
	c.dialer.Timeout = dialTimeout
	c.hostPort = hostPort
	c.h = h
	return &c
}

func (c *Client) Connect() {
	go func() {
		var (
			ms []wsutil.Message
		)
		for {
			// Connect Loop
		ConnectLoop:
			conn, _, _, err := ws.Dial(context.Background(), c.hostPort)
			if err != nil {
				goto ConnectLoop
			}
			c.conn = conn

			// Receive Loop
			for {
				ms = ms[:0]
				ms, err := wsutil.ReadServerMessage(c.conn, ms)
				if err != nil {
					break
				}
				for idx := range ms {
					switch ms[idx].OpCode {
					case ws.OpBinary, ws.OpText:
						e := rony.MessageEnvelope{}
						_ = e.Unmarshal(ms[idx].Payload)
						c.h(e)
					default:
					}

				}
			}
		}

	}()
	return
}

func (c *Client) Send(constructor int64, m rony.ProtoBufferMessage) error {
	e := rony.PoolMessageEnvelope.Get()
	e.Fill(tools.RandomUint64(), constructor, m)
	b := pools.Bytes.GetLen(e.Size())
	_, err := e.MarshalToSizedBuffer(b)
	if err != nil {
		return err
	}
	err = wsutil.WriteClientMessage(c.conn, ws.OpBinary, b)
	pools.Bytes.Put(b)
	return err
}
