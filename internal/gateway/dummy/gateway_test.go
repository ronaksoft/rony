package dummyGateway_test

import (
	"github.com/ronaksoft/rony"
	dummyGateway "github.com/ronaksoft/rony/internal/gateway/dummy"
	"github.com/ronaksoft/rony/tools"
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"sync/atomic"
	"testing"
)

/*
   Creation Time: 2021 - Jul - 03
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	_connID uint64
)

func nextConnID() uint64 {
	return atomic.AddUint64(&_connID, 1)
}

func TestDummyGateway(t *testing.T) {
	Convey("DummyGateway", t, func(c C) {
		gw, err := dummyGateway.New(dummyGateway.Config{})
		c.So(err, ShouldBeNil)
		Convey("NonPersistent - Single", nonPersistentSingle(gw, nextConnID()))
		Convey("NonPersistent - Concurrent", nonPersistentConcurrent(gw, nextConnID()))
	})
}
func nonPersistentSingle(gw *dummyGateway.Gateway, testConnID uint64) func(c C) {
	return func(c C) {
		inputData := tools.S2B(tools.RandomID(128))
		wg := sync.WaitGroup{}
		gw.MessageHandler = func(conn rony.Conn, streamID int64, data []byte) {
			err := conn.WriteBinary(streamID, data)
			c.So(err, ShouldBeNil)
		}
		gw.ConnectHandler = func(conn rony.Conn, kvs ...*rony.KeyValue) {
			// c.Println("Connect", conn.ConnID())
			c.So(conn.ConnID(), ShouldEqual, testConnID)
		}
		gw.CloseHandler = func(conn rony.Conn) {
			// c.Println("Close", conn.ConnID())
			c.So(conn.ConnID(), ShouldEqual, testConnID)
		}
		receiver := func(connID uint64, streamID int64, data []byte) {
			c.So(connID, ShouldEqual, testConnID)
			c.So(data, ShouldResemble, inputData)
			wg.Done()
		}

		wg.Add(1)
		gw.OpenConn(testConnID, false, receiver)
		err := gw.RPC(testConnID, 0, inputData)
		c.So(err, ShouldBeNil)
		gw.CloseConn(testConnID)
		wg.Wait()
	}
}
func nonPersistentConcurrent(gw *dummyGateway.Gateway, testConnID uint64) func(c C) {
	return func(c C) {
		inputData := tools.S2B(tools.RandomID(128))
		wg := sync.WaitGroup{}
		gw.MessageHandler = func(conn rony.Conn, streamID int64, data []byte) {
			err := conn.WriteBinary(streamID, data)
			c.So(err, ShouldBeNil)
		}
		gw.ConnectHandler = func(conn rony.Conn, kvs ...*rony.KeyValue) {
			// c.Println("Connect", conn.ConnID())
			c.So(conn.ConnID(), ShouldEqual, testConnID)
		}
		gw.CloseHandler = func(conn rony.Conn) {
			// c.Println("Close", conn.ConnID())
			c.So(conn.ConnID(), ShouldEqual, testConnID)
		}
		receiver := func(connID uint64, streamID int64, data []byte) {
			c.So(connID, ShouldEqual, testConnID)
			c.So(data, ShouldResemble, inputData)
			wg.Done()
		}

		gw.OpenConn(testConnID, false, receiver)
		for i := 0; i < 50; i++ {
			wg.Add(1)
			err := gw.RPC(testConnID, 0, inputData)
			c.So(err, ShouldBeNil)
		}
		gw.CloseConn(testConnID)
		wg.Wait()
	}
}
