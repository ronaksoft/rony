package benchs_test

import (
	"flag"
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/edge"
	"github.com/ronaksoft/rony/edgec"
	"github.com/ronaksoft/rony/internal/testEnv"
	"github.com/ronaksoft/rony/internal/testEnv/pb/service"
	"github.com/ronaksoft/rony/pools"
	"github.com/valyala/fasthttp"
	"google.golang.org/protobuf/proto"
	"os"
	"sync/atomic"
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

func BenchmarkEdge(b *testing.B) {
	benches := map[string]func(*testing.B){
		"Http-Client":            benchHttpClient,
		"Http-Single-EdgeClient": benchHttpSingleEdgeClient,
		// "SingleWebsocket": benchSingleWebsocketClient,
		// "MultiWebsocket":  benchMultiWebsocketClient,
	}
	for k, f := range benches {
		b.Run(k, f)
	}
}
func benchHttpClient(b *testing.B) {
	edgeClient := &fasthttp.HostClient{
		Addr: "127.0.0.1:8080",
	}
	echoRequest := service.EchoRequest{
		Int:       100,
		Timestamp: 32809238402,
	}
	reqID := uint64(1)
	cntErr := int64(0)
	cntOK := int64(0)

	b.ResetTimer()
	b.ReportAllocs()
	b.SetParallelism(10)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			rid := atomic.AddUint64(&reqID, 1)
			req := rony.PoolMessageEnvelope.Get()
			res := rony.PoolMessageEnvelope.Get()
			req.Fill(rid, service.C_SampleEcho, &echoRequest)
			err := sendHttp(edgeClient, req, res)
			if err != nil {
				atomic.AddInt64(&cntErr, 1)
			} else {
				atomic.AddInt64(&cntOK, 1)
			}

			rony.PoolMessageEnvelope.Put(req)
			rony.PoolMessageEnvelope.Put(res)
		}
	})
	b.Log("Failed / OK:", cntErr, "/", cntOK)
}
func benchHttpSingleEdgeClient(b *testing.B) {
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
	cntErr := int64(0)
	cntOK := int64(0)

	b.ResetTimer()
	b.ReportAllocs()
	b.SetParallelism(10)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			req := rony.PoolMessageEnvelope.Get()
			res := rony.PoolMessageEnvelope.Get()
			req.Fill(edgeClient.GetRequestID(), service.C_SampleEcho, &echoRequest)
			err = edgeClient.Send(req, res, true)
			if err != nil {
				atomic.AddInt64(&cntErr, 1)
			} else {
				atomic.AddInt64(&cntOK, 1)
			}
			rony.PoolMessageEnvelope.Put(req)
			rony.PoolMessageEnvelope.Put(res)
		}
	})
	b.Log("Failed / OK:", cntErr, "/", cntOK)
	_ = edgeClient.Close()

}
func benchMultiWebsocketClient(b *testing.B) {
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
func benchSingleWebsocketClient(b *testing.B) {
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

func sendHttp(c *fasthttp.HostClient, req, res *rony.MessageEnvelope) error {
	mo := proto.MarshalOptions{UseCachedSize: true}
	b := pools.Buffer.GetCap(mo.Size(req))
	*b.Bytes(), _ = mo.MarshalAppend(*b.Bytes(), req)
	httpReq := fasthttp.AcquireRequest()
	httpRes := fasthttp.AcquireResponse()
	httpReq.Header.SetHost("127.0.0.1:8080")
	httpReq.Header.SetMethod(fasthttp.MethodPost)
	httpReq.SetBody(*b.Bytes())
	err := c.Do(httpReq, httpRes)
	if err != nil {
		return err
	}
	err = res.Unmarshal(httpRes.Body())
	return err
}
