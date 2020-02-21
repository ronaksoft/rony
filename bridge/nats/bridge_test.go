package natsBridge_test

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/testEnv"

	"git.ronaksoftware.com/ronak/rony/config"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	natsBridge "git.ronaksoftware.com/ronak/rony/bridge/nats"
	"git.ronaksoftware.com/ronak/rony/tools"
	"github.com/nats-io/nats.go"

	. "github.com/smartystreets/goconvey/convey"
)

/*
   Creation Time: 2019 - Feb - 27
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

func init() {
	testEnv.Init()
}

func TestBridge(t *testing.T) {
	errHandler := func(e error) {
		t.Log(e.Error())
	}
	var bridgeSender, bridgeReceiver *natsBridge.Bridge
	Convey("Bridge Suite Test", t, func() {

		Convey("Initialize Two Bridge Ends", func(c C) {
			notifyCount := int32(0)
			messageCount := int32(0)
			_ = messageCount

			natsConfig := nats.GetDefaultOptions()
			natsConfig.Url = config.GetString(config.NatsURL)
			natsConfig.AsyncErrorCB = natsErrHandler

			var err error
			bridgeSender, err = natsBridge.NewBridge(natsBridge.Config{
				BundleID:   "Sender",
				InstanceID: "001",
				Options:    natsConfig,
				ErrHandler: errHandler,
				MessageHandler: func(c *natsBridge.Container) bool {
					atomic.AddInt32(&messageCount, int32(len(c.Messages)))
					time.Sleep(125 * time.Millisecond)
					return true
				},
				NotifyHandler: func(connIDs []uint64) {
					atomic.AddInt32(&notifyCount, 1)
				},
			})
			c.So(err, ShouldBeNil)
			c.So(bridgeSender, ShouldNotBeNil)

			bridgeReceiver, err = natsBridge.NewBridge(natsBridge.Config{
				BundleID:   "Receiver",
				InstanceID: "001",
				Options:    natsConfig,
				ErrHandler: errHandler,
				MessageHandler: func(container *natsBridge.Container) bool {
					atomic.AddInt32(&messageCount, int32(len(container.Messages)))
					time.Sleep(125 * time.Millisecond)
					return true
				},
				NotifyHandler: func(connIDs []uint64) {
					atomic.AddInt32(&notifyCount, 1)
				},
			})
			c.So(err, ShouldBeNil)
			c.So(bridgeReceiver, ShouldNotBeNil)
		})

		Convey("Wait 2 Seconds for bridges to connect to NATS", func(c C) {
			time.Sleep(time.Second * 2)
		})

		Convey("Start Sending Messages", func(c C) {
			c.So(bridgeSender, ShouldNotBeNil)
			c.So(bridgeReceiver, ShouldNotBeNil)
			waitGroup := sync.WaitGroup{}
			for i := 0; i < 100; i++ {
				waitGroup.Add(1)
				go func() {
					for j := 0; j < 10; j++ {
						delivered := bridgeSender.SendMessage("Receiver", tools.RandomInt64(0), 10, tools.StrToByte(tools.RandomID(62)))
						c.So(delivered, ShouldBeTrue)
					}
					waitGroup.Done()
				}()
			}
			waitGroup.Wait()
		})

	})
}

func natsErrHandler(nc *nats.Conn, sub *nats.Subscription, natsErr error) {
	fmt.Printf("error: %v\n", natsErr)
	if natsErr == nats.ErrSlowConsumer {
		pendingMsgs, _, err := sub.Pending()
		if err != nil {
			fmt.Printf("couldn't get pending messages: %v", err)
			return
		}
		fmt.Printf("Falling behind with %d pending messages on subject %q.\n",
			pendingMsgs, sub.Subject)
		// Log error, notify operations...
	}
	// check for other errors
}
