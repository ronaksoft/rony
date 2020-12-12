package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/ronaksoft/rony/gateway"
	tcpGateway "github.com/ronaksoft/rony/gateway/tcp"
	"github.com/ronaksoft/rony/internal/logger"
	"go.uber.org/zap"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

/*
   Creation Time: 2019 - Jun - 05
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	cntConnected      int32
	cntReceivedPacket int32
	cntWritePacket    int32
	cntServerReceived int32
	latency           int64
	ps                = 4 * 1024
)

func main() {
	fmt.Println("RUN...")
	log.Init(log.DefaultConfig)
	log.SetLevel(log.InfoLevel)
	n := 5000
	m := 10
	log.Info("Server Starting ...")
	gw := runServer()
	log.Info("Server Started.")
	time.Sleep(time.Second * 1)

	go func() {
		for {
			log.Info("Stats",
				zap.Int32("Received", cntReceivedPacket),
				zap.Int32("Sent", cntWritePacket),
				zap.Int32("Server Received", cntServerReceived),
				zap.Int32("Connected", cntConnected),
			)
			time.Sleep(time.Second)
		}
	}()

	runClients(n, m)
	log.Info("Final Stats",
		zap.Int32("Received", cntReceivedPacket),
		zap.Int32("Sent", cntWritePacket),
		zap.Int32("Server Received", cntServerReceived),
		zap.Int32("Connected", cntConnected),
		zap.Duration("Avg Latency", time.Duration(latency/int64(cntReceivedPacket))),
	)
	gw.Shutdown()
}

func runClients(n, m int) {
	waitGroup := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		waitGroup.Add(1)
		go func(i int) {
			runClient(&waitGroup, m)
			waitGroup.Done()
		}(i)
		time.Sleep(time.Millisecond)
	}
	waitGroup.Wait()
	log.Debug("Running Clients Finished")
}

func runClient(wg *sync.WaitGroup, m int) {
	var (
		c   net.Conn
		err error
	)

	for {
		c, _, _, err = ws.Dialer{}.Dial(context.Background(), "ws://localhost:8080")
		if err == nil {
			break
		}
		log.Warn("error on connect: ", zap.Error(err))
		time.Sleep(time.Millisecond)
	}
	atomic.AddInt32(&cntConnected, 1)
	wg.Add(1)
	go func() {
		time.Sleep(time.Second)
		defer wg.Done()
		var m []wsutil.Message
		for {
			m = m[:0]
			_ = c.SetReadDeadline(time.Now().Add(time.Second * 3))
			startTime := time.Now()
			m, err = wsutil.ReadServerMessage(c, m)
			d := time.Now().Sub(startTime)
			if err != nil {
				err = c.Close()
				if err != nil {
					log.Warn("Error on Close", zap.Error(err))
				}
				return
			}
			atomic.AddInt64(&latency, int64(d))
			atomic.AddInt32(&cntReceivedPacket, int32(len(m)))

		}
	}()
	sentData := make([]byte, ps)
	rand.Read(sentData)
	for i := 0; i < m; i++ {
		err := wsutil.WriteClientMessage(c, ws.OpBinary, sentData)
		if err != nil {
			log.Error(err.Error())
			_ = c.Close()
			return
		} else {
			atomic.AddInt32(&cntWritePacket, 1)
		}
	}
}

func runServer() gateway.Gateway {
	gw, err := tcpGateway.New(tcpGateway.Config{
		Concurrency:   1000,
		ListenAddress: ":8080",
	})

	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	gw.MessageHandler = func(c gateway.Conn, streamID int64, data []byte) {
		atomic.AddInt32(&cntServerReceived, 1)
		err := c.SendBinary(0, []byte("HI"))
		if err != nil {
			log.Warn("Error On SendBinary", zap.Error(err))
		}
	}
	go gw.Run()
	return gw
}
