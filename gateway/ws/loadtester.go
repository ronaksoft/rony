// +build ignore

package main

import (
	"context"
	"crypto/rand"
	"git.ronaksoftware.com/ronak/rony/internal/logger"
	websocketGateway "git.ronaksoftware.com/ronak/toolbox/gateway/ws"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
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
   Copyright Ronak Software Group 2018
*/

var (
	cntReceivedPacket int32
	cntWritePacket    int32
	cntServerReceived int32
	_Log              log.Logger
)

func main() {
	n := 1
	m := 100
	runServer()
	time.Sleep(time.Second * 3)
	startTime := time.Now()
	runClients(n, m)
	log.Info("Clients Finished",
		zap.Duration("Elapsed", time.Now().Sub(startTime)),
	)

	for i := 0; i < 30; i++ {
		log.Info("Stats",
			zap.Int32("Received", cntReceivedPacket),
			zap.Int32("Sent", cntWritePacket),
			zap.Int32("Server Received", cntServerReceived),
		)
		time.Sleep(time.Second)
	}

}

func runClients(n, m int) {
	waitGroup := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		waitGroup.Add(1)
		go func(i int) {
			runClient(m)
			waitGroup.Done()
		}(i)
		time.Sleep(10 * time.Millisecond)
	}
	waitGroup.Wait()
	log.Debug("Running Clients Finished")
}

func runClient(m int) {
	var c net.Conn
	var err error
	for {
		c, _, _, err = ws.Dialer{}.Dial(context.Background(), "ws://localhost:8080")
		if err == nil {
			break
		}
		log.Warn("error on connect: ", zap.Error(err))
		time.Sleep(time.Millisecond * 10)
	}
	go func() {
		time.Sleep(time.Second)
		for {
			_, _, err = wsutil.ReadServerData(c)
			if err != nil {
				if nerr, ok := err.(net.Error); ok && !nerr.Temporary() {
					log.Error(err.Error())
					return
				} else {
					log.Warn(err.Error())
				}
				return
			}
			atomic.AddInt32(&cntReceivedPacket, 1)
		}
	}()
	sentData := make([]byte, 1<<20)
	rand.Read(sentData)
	for i := 0; i < m; i++ {
		err := wsutil.WriteClientBinary(c, sentData)
		if err != nil {
			log.Error(err.Error())
			c.Close()
			return
		} else {
			atomic.AddInt32(&cntWritePacket, 1)
		}
	}
}

func runServer() {
	websocketGateway.SetLogger(_Log)
	gw, err := websocketGateway.New(websocketGateway.Config{
		CloseHandler: func(c *websocketGateway.Conn) {},
		MessageHandler: func(c *websocketGateway.Conn, date []byte) {
			atomic.AddInt32(&cntServerReceived, 1)
			time.Sleep(time.Millisecond * 1000)
			c.SendBinary(tools.StrToByte("Hi"))
		},
		ConnectHandler: func(connID uint64) {},
		FlushFunc: func(c *websocketGateway.Conn) [][]byte {
			// return [][]byte{{1}, {2}}
			return nil
		},
		NewConnectionWorkers: 10,
		MaxConcurrency:       10000,
		ListenAddress:        ":8080",
		PingerPeriod:         time.Minute,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
	go gw.Run()

}
