package udpTunnel

import (
	"encoding/binary"
	"fmt"
	"github.com/panjf2000/gnet"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/log"
	"github.com/ronaksoft/rony/tunnel"
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
	cfg       Config
	addrs     []string
	shutdown  int32 // atomic shutdown flag
	connID    uint64
	encConfig gnet.EncoderConfig
	decConfig gnet.DecoderConfig
}

func New(config Config) (*Tunnel, error) {
	t := &Tunnel{
		cfg: config,
	}
	t.encConfig = gnet.EncoderConfig{
		ByteOrder:                       binary.BigEndian,
		LengthFieldLength:               4,
		LengthAdjustment:                0,
		LengthIncludesLengthFieldLength: false,
	}
	t.decConfig = gnet.DecoderConfig{
		ByteOrder:           binary.BigEndian,
		LengthFieldOffset:   0,
		LengthFieldLength:   4,
		LengthAdjustment:    0,
		InitialBytesToStrip: 0,
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
		gnet.WithNumEventLoop(t.cfg.Concurrency),
		gnet.WithCodec(gnet.NewLengthFieldBasedFrameCodec(t.encConfig, t.decConfig)),
	)

	if err != nil {
		panic(err)
	}
}

func (t *Tunnel) Shutdown() {
	atomic.StoreInt32(&t.shutdown, 1)
}

func (t *Tunnel) Addr() []string {
	if len(t.cfg.ExternalAddrs) > 0 {
		return t.cfg.ExternalAddrs
	}
	return t.addrs
}

func (t *Tunnel) OnInitComplete(server gnet.Server) (action gnet.Action) {
	var hosts []string
	// try to detect the ip address of the listener
	ta, err := net.ResolveUDPAddr("udp", t.cfg.ListenAddress)
	if err != nil {
		return gnet.None
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

	ta, err = net.ResolveUDPAddr("udp", server.Addr.String())
	if err != nil {
		panic(err)
	}
	for _, h := range hosts {
		t.addrs = append(t.addrs, fmt.Sprintf("%s:%d", h, ta.Port))
	}

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
	log.Info("Tunnel connection closed")
	return gnet.None
}

func (t *Tunnel) PreWrite() {
}

func (t *Tunnel) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	if atomic.LoadInt32(&t.shutdown) == 1 {
		return nil, gnet.Shutdown
	}
	// log.Info("Tunnel On Data", zap.String("Remote", c.RemoteAddr().String()))
	if len(frame) > 0 {
		req := rony.PoolTunnelMessage.Get()
		defer rony.PoolTunnelMessage.Put(req)

		err := req.Unmarshal(frame)
		if err != nil {
			log.Warn("Error On Tunnel's data received", zap.Error(err))
			return nil, gnet.Close
		}

		conn := newConn(t.nextID(), c)
		t.MessageHandler(conn, req)
	}

	uc, _ := c.Context().(*udpConn)
	if uc == nil {
		return nil, gnet.None
	}

	bb := uc.Pop()
	if bb == nil {
		return nil, gnet.None
	}
	return *bb.Bytes(), gnet.None
}

func (t *Tunnel) Tick() (delay time.Duration, action gnet.Action) {
	return 0, gnet.None
}
