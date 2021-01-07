package udpTunnel

import (
	"fmt"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/tunnel"
	"github.com/tidwall/evio"
	"go.uber.org/zap"
	"sync/atomic"
)

/*
   Creation Time: 2021 - Jan - 04
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// Config
type Config struct {
	Concurrency   int
	ListenAddress string
	MaxBodySize   int
	ExternalAddrs []string
}

// Tunnel
type Tunnel struct {
	tunnel.MessageHandler
	events   evio.Events
	cfg      Config
	addrs    []string
	shutdown int32 // atomic shutdown flag
	connID   uint64
}

func New(config Config) *Tunnel {
	t := &Tunnel{
		cfg: config,
	}
	t.events.NumLoops = t.cfg.Concurrency
	t.events.Data = t.onData
	t.events.Serving = func(server evio.Server) (action evio.Action) {
		log.Info("UDP Tunnel started", zap.Any("Addr", server.Addrs))
		for _, a := range server.Addrs {
			t.addrs = append(t.addrs, a.String())
		}
		return evio.None
	}
	return t
}

func (t *Tunnel) nextID() uint64 {
	return atomic.AddUint64(&t.connID, 1)
}

func (t *Tunnel) onData(c evio.Conn, in []byte) (out []byte, action evio.Action) {
	if atomic.LoadInt32(&t.shutdown) == 1 {
		return nil, evio.Shutdown
	}

	if len(in) > 0 {
		req := rony.PoolTunnelMessage.Get()
		defer rony.PoolTunnelMessage.Put(req)

		err := req.Unmarshal(in)
		if err != nil {
			log.Warn("Error On Tunnel's data received", zap.Error(err))
			return nil, evio.Close
		}

		conn := newConn(t.nextID(), c)
		t.MessageHandler(conn, req)
	}

	uc, _ := c.Context().(*udpConn)
	if uc == nil {
		return nil, evio.None
	}

	bb := uc.Pop()
	if bb == nil {
		return nil, evio.None
	}
	return *bb.Bytes(), evio.None
}

func (t *Tunnel) Start() {
	go t.Run()
}

func (t *Tunnel) Run() {
	err := evio.Serve(t.events, fmt.Sprintf("udp://%s?reuseport=true", t.cfg.ListenAddress))
	if err != nil {
		panic(err)
	}
}

func (t *Tunnel) Shutdown() {
	atomic.StoreInt32(&t.shutdown, 1)
}

func (t *Tunnel) Addr() []string {
	return t.cfg.ExternalAddrs
}
