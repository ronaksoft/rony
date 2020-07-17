package edgeClient_test

import (
	"git.ronaksoftware.com/ronak/rony"
	"git.ronaksoftware.com/ronak/rony/edgeClient"
	"git.ronaksoftware.com/ronak/rony/internal/testEnv"
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
	c := edgeClient.NewWebsocket(edgeClient.Config{
		HostPort: "ws://127.0.0.1:8081",
	})
	c.Connect()
	_, err := c.Send(rony.C_MessageContainer, &rony.MessageContainer{})
	if err != nil {
		t.Fatal(err)
	}

}
