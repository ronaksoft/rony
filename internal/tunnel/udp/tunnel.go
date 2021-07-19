package udpTunnel

import (
	"context"
	"fmt"
	"github.com/panjf2000/gnet"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/internal/metrics"
	"github.com/ronaksoft/rony/internal/msg"
	"github.com/ronaksoft/rony/internal/tunnel"
	"github.com/ronaksoft/rony/pools"
	"github.com/ronaksoft/rony/tools"
	"go.uber.org/zap"
	"net"
	"sync/atomic"
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

type Config struct {
	ServerID      string
	ListenAddress string
	MaxBodySize   int
	ExternalAddrs []string
}

// Tunnel
type Tunnel struct {
	tunnel.MessageHandler
	cfg      Config
	addrs    []string
	shutdown int32 // atomic shutdown flag
	connID   uint64
}

func New(config Config) (*Tunnel, error) {
	t := &Tunnel{
		cfg: config,
	}

	var hosts []string
	// try to detect the ip address of the listener
	ta, err := net.ResolveUDPAddr("udp", t.cfg.ListenAddress)
	if err != nil {
		return nil, err
	}
	if ta.Port == 0 {
		ta.Port = tools.RandomInt(63000)
		t.cfg.ListenAddress = ta.String()
	}

	if ta.IP.IsUnspecified() {
		addrs, err := net.InterfaceAddrs()

		if err == nil {
			for _, a := range addrs {
				switch x := a.(type) {
				case *net.IPNet:
					if x.IP.To4() == nil || x.IP.IsLoopback() {
						continue
					}
					hosts = append(hosts, x.IP.String())
				case *net.IPAddr:
					if x.IP.To4() == nil || x.IP.IsLoopback() {
						continue
					}
					hosts = append(hosts, x.IP.String())
				case *net.UDPAddr:
					if x.IP.To4() == nil || x.IP.IsLoopback() {
						continue
					}
					hosts = append(hosts, x.IP.String())
				}
			}
		}
	}

	for _, h := range hosts {
		t.addrs = append(t.addrs, fmt.Sprintf("%s:%d", h, ta.Port))
	}

	return t, nil
}

func (t *Tunnel) nextID() uint64 {
	return atomic.AddUint64(&t.connID, 1)
}

func (t *Tunnel) Start() {
	go t.Run()
	time.Sleep(time.Millisecond * 100)
}

func (t *Tunnel) Run() {

	err := gnet.Serve(t, fmt.Sprintf("udp://%s", t.cfg.ListenAddress),
		gnet.WithReusePort(true),
		gnet.WithMulticore(true),
		gnet.WithLockOSThread(true),
		gnet.WithLogLevel(log.WarnLevel),
		gnet.WithLogger(log.DefaultLogger.Logger.Sugar()),
	)

	if err != nil {
		panic(err)
	}
}

func (t *Tunnel) Shutdown() {
	atomic.StoreInt32(&t.shutdown, 1)
	ctx, cf := context.WithTimeout(context.TODO(), time.Second*30)
	defer cf()
	if err := gnet.Stop(ctx, fmt.Sprintf("udp://%s", t.cfg.ListenAddress)); err != nil {
		log.Warn("Error On Stopping Tunnel", zap.Error(err))
	}
}

func (t *Tunnel) Addr() []string {
	if len(t.cfg.ExternalAddrs) > 0 {
		return t.cfg.ExternalAddrs
	}
	return t.addrs
}

func (t *Tunnel) OnInitComplete(server gnet.Server) (action gnet.Action) {
	return gnet.None
}

func (t *Tunnel) OnShutdown(server gnet.Server) {
	log.Info("Tunnel shutdown")
}

func (t *Tunnel) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Info("Tunnel connection opened")
	return nil, gnet.None
}

func (t *Tunnel) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	log.Info("Tunnel connection closed", zap.Error(err))
	return gnet.None
}

func (t *Tunnel) PreWrite() {
}

func (t *Tunnel) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	if atomic.LoadInt32(&t.shutdown) == 1 {
		return nil, gnet.Shutdown
	}

	req := msg.PoolTunnelMessage.Get()
	if err := req.Unmarshal(frame); err != nil {
		log.Warn("Error On Tunnel's data received", zap.Error(err))
		return nil, gnet.Close
	}

	conn := newConn(t.nextID(), c)
	pools.Go(func() {
		metrics.IncCounter(metrics.CntTunnelIncomingMessage)
		t.MessageHandler(conn, req)
		msg.PoolTunnelMessage.Put(req)
	})

	return
}

func (t *Tunnel) Tick() (delay time.Duration, action gnet.Action) {
	return time.Minute, gnet.None
}
