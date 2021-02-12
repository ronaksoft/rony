package benchs_test

import (
	"flag"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/edgec"
	"github.com/ronaksoft/rony/internal/testEnv"
	"github.com/ronaksoft/rony/internal/testEnv/pb/service"
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
	edgeServer = testEnv.InitEdgeServer("Adam", 8080, 100000)
	rony.SetLogLevel(0)
	service.RegisterSample(
		&testEnv.Handlers{
			ServerID: edgeServer.GetServerID(),
		},
		edgeServer,
	)

	edgeServer.Start()

	flag.Parse()
	code := m.Run()
	edgeServer.Shutdown()
	os.Exit(code)
}

func BenchmarkSingleHttpClient(b *testing.B) {
	edgeClient := edgec.NewHttp(edgec.HttpConfig{
		Name:            "Benchmark",
		SeedHostPort:    "127.0.0.1:8080",
		ReadTimeout:     time.Second,
		WriteTimeout:    time.Second,
		ContextTimeout:  time.Second,
		RequestMaxRetry: 10,
	})
	err := edgeClient.Start()
	if err != nil {
		b.Fatal(err)
	}
	echoRequest := service.EchoRequest{
		Int:       100,
		Timestamp: 32809238402,
	}

	b.ResetTimer()
	b.ReportAllocs()
	b.SetParallelism(10)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			req := rony.PoolMessageEnvelope.Get()
			res := rony.PoolMessageEnvelope.Get()
			req.Fill(edgeClient.GetRequestID(), service.C_SampleEcho, &echoRequest)
			_ = edgeClient.Send(req, res, true)
			rony.PoolMessageEnvelope.Put(req)
			rony.PoolMessageEnvelope.Put(res)
		}
	})
	_ = edgeClient.Close()
}

func BenchmarkSingleWebsocketClient(b *testing.B) {
	edgeClient := edgec.NewWebsocket(edgec.WebsocketConfig{
		SeedHostPort:    "127.0.0.1:8080",
		IdleTimeout:     time.Second,
		DialTimeout:     time.Second,
		Handler:         func(m *rony.MessageEnvelope) {},
		RequestMaxRetry: 10,
		RequestTimeout:  time.Second,
		// ContextTimeout:  time.Second,
	})
	err := edgeClient.Start()
	if err != nil {
		b.Fatal(err)
	}
	echoRequest := service.EchoRequest{
		Int:       100,
		Timestamp: 32809238402,
	}

	b.ResetTimer()
	b.ReportAllocs()
	b.SetParallelism(10)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			req := rony.PoolMessageEnvelope.Get()
			res := rony.PoolMessageEnvelope.Get()
			req.Fill(edgeClient.GetRequestID(), service.C_SampleEcho, &echoRequest)
			_ = edgeClient.Send(req, res, true)
			rony.PoolMessageEnvelope.Put(req)
			rony.PoolMessageEnvelope.Put(res)
		}
	})
	_ = edgeClient.Close()
}

func BenchmarkMultiHttpClient(b *testing.B) {
	echoRequest := service.EchoRequest{
		Int:       100,
		Timestamp: 32809238402,
	}
	b.ResetTimer()
	b.ReportAllocs()
	// b.SetParallelism(5)
	b.RunParallel(func(p *testing.PB) {
		edgeClient := edgec.NewHttp(edgec.HttpConfig{
			Name:            "Benchmark",
			SeedHostPort:    "127.0.0.1:8080",
			ReadTimeout:     time.Second,
			WriteTimeout:    time.Second,
			ContextTimeout:  time.Second,
			RequestMaxRetry: 10,
		})
		err := edgeClient.Start()
		if err != nil {
			b.Fatal(err)
		}
		for p.Next() {
			req := rony.PoolMessageEnvelope.Get()
			res := rony.PoolMessageEnvelope.Get()
			req.Fill(edgeClient.GetRequestID(), service.C_SampleEcho, &echoRequest)
			_ = edgeClient.Send(req, res, true)
			rony.PoolMessageEnvelope.Put(req)
			rony.PoolMessageEnvelope.Put(res)
		}
		_ = edgeClient.Close()
	})
}

func BenchmarkMultiWebsocketClient(b *testing.B) {
	echoRequest := service.EchoRequest{
		Int:       100,
		Timestamp: 32809238402,
	}
	b.ResetTimer()
	b.ReportAllocs()
	// b.SetParallelism(5)
	b.RunParallel(func(p *testing.PB) {
		edgeClient := edgec.NewWebsocket(edgec.WebsocketConfig{
			SeedHostPort:    "127.0.0.1:8080",
			IdleTimeout:     time.Second,
			DialTimeout:     time.Second,
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
			req.Fill(edgeClient.GetRequestID(), service.C_SampleEcho, &echoRequest)
			_ = edgeClient.Send(req, res, true)
			rony.PoolMessageEnvelope.Put(req)
			rony.PoolMessageEnvelope.Put(res)
		}
		_ = edgeClient.Close()
	})

}
