package benchs_test

import (
	"flag"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/edgec"
	"github.com/ronaksoft/rony/internal/testEnv"
	"github.com/ronaksoft/rony/internal/testEnv/pb"
	"os"
	"testing"
	"time"
)

/*
   Creation Time: 2020 - Dec - 28
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	edgeServer *edge.Server
)

func TestMain(m *testing.M) {
	edgeServer = testEnv.InitEdgeServerWithWebsocket("Adam", 8080, 1000)

	pb.RegisterSample(
		&testEnv.Handlers{
			ServerID: edgeServer.GetServerID(),
		}, edgeServer)
	err := edgeServer.StartCluster()
	if err != nil {
		panic(err)
	}

	err = edgeServer.StartGateway()
	if err != nil {
		panic(err)
	}
	flag.Parse()
	code := m.Run()
	edgeServer.Shutdown()
	os.Exit(code)
}

func BenchmarkSingleClient(b *testing.B) {
	edgeClient := edgec.NewWebsocket(edgec.WebsocketConfig{
		SeedHostPort:    "127.0.0.1:8080",
		IdleTimeout:     time.Second,
		DialTimeout:     time.Second,
		ForceConnect:    true,
		Handler:         func(m *rony.MessageEnvelope) {},
		RequestMaxRetry: 10,
		RequestTimeout:  time.Second,
		// ContextTimeout:  time.Second,
	})
	err := edgeClient.Start()
	if err != nil {
		b.Fatal(err)
	}
	echoRequest := pb.EchoRequest{
		Int:       100,
		Timestamp: 32809238402,
	}

	b.ResetTimer()
	b.ReportAllocs()
	// b.SetParallelism(10)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			req := rony.PoolMessageEnvelope.Get()
			res := rony.PoolMessageEnvelope.Get()
			req.Fill(edgeClient.GetRequestID(), pb.C_Echo, &echoRequest)
			_ = edgeClient.Send(req, res)
			rony.PoolMessageEnvelope.Put(req)
			rony.PoolMessageEnvelope.Put(res)
		}
	})
}

func BenchmarkMultiClient(b *testing.B) {
	echoRequest := pb.EchoRequest{
		Int:       100,
		Timestamp: 32809238402,
	}
	b.ResetTimer()
	b.ReportAllocs()
	// b.SetParallelism(10)
	b.RunParallel(func(p *testing.PB) {
		edgeClient := edgec.NewWebsocket(edgec.WebsocketConfig{
			SeedHostPort:    "127.0.0.1:8080",
			IdleTimeout:     time.Second,
			DialTimeout:     time.Second,
			ForceConnect:    true,
			Handler:         func(m *rony.MessageEnvelope) {},
			RequestMaxRetry: 10,
			RequestTimeout:  time.Second,
			// ContextTimeout:  time.Second,
		})
		err := edgeClient.Start()
		if err != nil {
			b.Fatal(err)
		}
		for p.Next() {
			req := rony.PoolMessageEnvelope.Get()
			res := rony.PoolMessageEnvelope.Get()
			req.Fill(edgeClient.GetRequestID(), pb.C_Echo, &echoRequest)
			_ = edgeClient.Send(req, res)
			rony.PoolMessageEnvelope.Put(req)
			rony.PoolMessageEnvelope.Put(res)
		}
	})
}
