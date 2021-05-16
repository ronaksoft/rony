package tools_test

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/tools"
	. "github.com/smartystreets/goconvey/convey"
	"google.golang.org/protobuf/proto"
	"testing"
)

/*
   Creation Time: 2020 - Nov - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Alias int32

const (
	Alias_A1 Alias = iota + 1
	Alias_A2
	Alias_A3
)

func TestNewAllocator(t *testing.T) {
	Convey("Allocator", t, func(c C) {
		alloc := tools.NewAllocator()
		b := alloc.Gen(Alias_A2)
		c.So(b, ShouldHaveLength, 4)
		b = alloc.Gen(Alias_A2, Alias_A1)
		c.So(b, ShouldHaveLength, 8)
		b = alloc.Gen(Alias_A1, "TXT1")
		c.So(b, ShouldHaveLength, 8)
		b = alloc.Gen(Alias_A1, 3232)
		c.So(b, ShouldHaveLength, 12)
		b = alloc.Gen(Alias_A1, 3232, []byte("TXT1"))
		c.So(b, ShouldHaveLength, 16)
		alloc.ReleaseAll()
	})
}
func BenchmarkBulkKey_GenKey(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		bk := tools.NewAllocator()
		for pb.Next() {
			d := bk.Gen(tools.FastRand(), "tools.RandomID(10)", 3, 125)
			if len(d) != 38 {
				b.Fatal("invalid size", len(d))
			}
		}
		bk.ReleaseAll()
	})
}

func BenchmarkAllocator_GenValue(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		m := rony.PoolMessageEnvelope.Get()
		m.RequestID = tools.RandomUint64(0)
		m.Constructor = 3232
		m.Message = append(m.Message, tools.StrToByte("Something here for test ONLY!")...)
		bk := tools.NewAllocator()
		for pb.Next() {
			d := bk.Gen(m)
			if len(d) != proto.Size(m) {
				b.Fatal("invalid size", len(d))
			}
		}
		bk.ReleaseAll()
		rony.PoolMessageEnvelope.Put(m)
	})
}
