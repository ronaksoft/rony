package rony_test

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/tools"
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"testing"
)

/*
   Creation Time: 2020 - Dec - 16
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func TestMessageEnvelope_Clone(t *testing.T) {
	Convey("Clone (MessageEnvelope)", t, func(c C) {
		src := &rony.MessageEnvelope{
			RequestID:   tools.RandomUint64(0),
			Constructor: tools.RandomInt64(0),
			Header: []*rony.KeyValue{
				{
					Key:   "Key1",
					Value: "Value1",
				},
				{
					Key:   "Key2",
					Value: "Value2",
				},
			},
			Auth: tools.StrToByte(tools.RandomID(10)),
		}

		wg := sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				dst := src.Clone()
				c.So(dst, ShouldResemble, src)
				wg.Done()
			}()
		}
		wg.Wait()

		src = &rony.MessageEnvelope{
			RequestID:   tools.RandomUint64(0),
			Constructor: tools.RandomInt64(0),
			Auth:        tools.StrToByte(tools.RandomID(10)),
		}

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				dst := src.Clone()
				c.So(dst, ShouldResemble, src)
				wg.Done()
			}()
		}
		wg.Wait()
	})
}
