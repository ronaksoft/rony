package edgeClient_test

import (
	"git.ronaksoftware.com/ronak/rony/edgeClient"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv/pb"
	"testing"
)

/*
   Creation Time: 2020 - Jul - 17
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func TestClient_Connect(t *testing.T) {
	testEnv.Init()

	c := pb.NewSampleClient(edgeClient.NewWebsocket(edgeClient.Config{
		HostPort: "ws://127.0.0.1:8081",
	}))
	res, err := c.Func1(&pb.Req1{Item1: 123})
	if err != nil {
		panic(err)
	}
	t.Log(res)
}
