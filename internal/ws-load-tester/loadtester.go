package main

import (
	"context"
	"fmt"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/internal/logger"
	"github.com/ronaksoft/rony/internal/testEnv/pb"
	"github.com/ronaksoft/rony/tools"
	"go.uber.org/zap"
	"net"
	"os"
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
	latency           int64
)

func main() {
	fmt.Println("RUN...")
	log.Init(log.DefaultConfig)
	log.SetLevel(log.InfoLevel)

	var n, m, port int64
	switch len(os.Args) {
	case 4:
		n = tools.StrToInt64(os.Args[1])
		m = tools.StrToInt64(os.Args[2])
		port = tools.StrToInt64(os.Args[3])
	default:
		fmt.Println("needs 3 args, n m port")
		return
	}

	go func() {
		for {
			log.Info("Stats",
				zap.Int32("Received", cntReceivedPacket),
				zap.Int32("Sent", cntWritePacket),
				zap.Int32("Connected", cntConnected),
			)
			time.Sleep(time.Second)
		}
	}()

	runClients(int(n), int(m), int(port))
	log.Info("Final Stats",
		zap.Int32("Received", cntReceivedPacket),
		zap.Int32("Sent", cntWritePacket),
		zap.Int32("Connected", cntConnected),
		zap.Duration("Avg Latency", time.Duration(latency/int64(cntReceivedPacket))),
	)
}

func runClients(n, m, port int) {
	waitGroup := sync.WaitGroup{}
	for i := 0; i < n; i++ {
		waitGroup.Add(1)
		go func(i int) {
			runClient(&waitGroup, m, port)
			waitGroup.Done()
		}(i)
		time.Sleep(time.Millisecond)
	}
	waitGroup.Wait()
	log.Debug("Running Clients Finished")
}

func runClient(wg *sync.WaitGroup, m int, port int) {
	var (
		c   net.Conn
		err error
	)

	for {
		c, _, _, err = ws.Dialer{}.Dial(context.Background(), fmt.Sprintf("ws://localhost:%d", port))
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
			startTime := tools.CPUTicks()
			m, err = wsutil.ReadServerMessage(c, m)
			d := time.Duration(tools.CPUTicks() - startTime)
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

	echoRequest := pb.EchoRequest{
		Int:       21313,
		Timestamp: 42342342342,
	}
	req := rony.PoolMessageEnvelope.Get()
	req.Fill(tools.RandomUint64(0), pb.C_Echo, &echoRequest)
	reqBytes, _ := req.Marshal()
	for i := 0; i < m; i++ {
		err := wsutil.WriteClientMessage(c, ws.OpBinary, reqBytes)
		if err != nil {
			log.Error(err.Error())
			_ = c.Close()
			return
		} else {
			atomic.AddInt32(&cntWritePacket, 1)
		}
	}
}
