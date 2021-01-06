package udpTunnel

import (
	"fmt"
	"github.com/ronaksoft/rony/tunnel"
	"github.com/tidwall/evio"
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

// Config
type Config struct {
	Concurrency   int
	ListenAddress string
	MaxBodySize   int
	MaxIdleTime   time.Duration
	ExternalAddrs []string
}

// Tunnel
type Tunnel struct {
	tunnel.MessageHandler
	events evio.Events
	// Internals
	listenOn    string
	addrs       []string
	extAddrs    []string
	concurrency int
	maxBodySize int
}

func NewTunnel(config Config) *Tunnel {
	t := &Tunnel{}
	return t
}

func (g *Tunnel) Start() {
	go g.Run()
}

func (g *Tunnel) Run() {
	err := evio.Serve(g.events, fmt.Sprintf("udp://%s?reuseport=true", g.listenOn))
	if err != nil {
		panic(err)
	}
}

func (g *Tunnel) Shutdown() {
	panic("implement me")
}

func (g *Tunnel) Addr() []string {
	panic("implement me")
}
