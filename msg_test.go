package rony_test

import (
	"sync"
	"testing"

	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/tools"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/protobuf/proto"
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
			Constructor: tools.RandomUint64(0),
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
			Constructor: tools.RandomUint64(0),
			Auth:        tools.StrToByte(tools.RandomID(10)),
		}

		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				dst := src.Clone()
				c.So(proto.Equal(dst, src), ShouldBeTrue)
				wg.Done()
			}()
		}
		wg.Wait()
	})
}

func TestMessageEnvelope_Get(t *testing.T) {
	Convey("Get (MessageEnvelope)", t, func(c C) {
		s := &rony.MessageEnvelope{
			Constructor: tools.RandomUint64(0),
			RequestID:   tools.RandomUint64(0),
			Message:     tools.S2B(tools.RandomID(1024)),
			Auth:        nil,
			Header: []*rony.KeyValue{
				{Key: "Key1", Value: "V1"},
				{Key: "Key2", Value: "V2"},
			},
		}
		sb, _ := s.Marshal()
		wg := &sync.WaitGroup{}
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				x := rony.PoolMessageEnvelope.Get()
				_ = x.Unmarshal(sb)
				c.So(proto.Equal(x, s), ShouldBeTrue)
				rony.PoolMessageEnvelope.Put(x)
				wg.Done()
			}()
		}
		wg.Wait()
	})
}
func BenchmarkMessageEnvelope_GetWithPool(b *testing.B) {
	s := &rony.MessageEnvelope{
		Constructor: tools.RandomUint64(0),
		RequestID:   tools.RandomUint64(0),
		Message:     tools.S2B(tools.RandomID(1024)),
		Auth:        nil,
		Header: []*rony.KeyValue{
			{Key: "Key1", Value: "V1"},
			{Key: "Key2", Value: "V2"},
		},
	}
	sb, _ := s.Marshal()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			x := rony.PoolMessageEnvelope.Get()
			_ = x.Unmarshal(sb)
			rony.PoolMessageEnvelope.Put(x)
		}
	})
}

func BenchmarkMessageEnvelope_Get(b *testing.B) {
	s := &rony.MessageEnvelope{
		Constructor: tools.RandomUint64(0),
		RequestID:   tools.RandomUint64(0),
		Message:     tools.S2B(tools.RandomID(1024)),
		Auth:        nil,
		Header: []*rony.KeyValue{
			{Key: "Key1", Value: "V1"},
			{Key: "Key2", Value: "V2"},
		},
	}
	sb, _ := s.Marshal()
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			x := &rony.MessageEnvelope{}
			_ = x.Unmarshal(sb)
		}
	})
}
