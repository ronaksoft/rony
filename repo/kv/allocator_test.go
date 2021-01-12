package kv

import (
	"github.com/ronaksoft/rony"
	"github.com/ronaksoft/rony/tools"
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

func BenchmarkBulkKey_GenKey(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		bk := NewAllocator()
		for pb.Next() {
			d := bk.GenKey(tools.FastRand(), "tools.RandomID(10)", 3, 125)
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
		bk := NewAllocator()
		for pb.Next() {
			d := bk.GenValue(m)
			if len(d) != proto.Size(m) {
				b.Fatal("invalid size", len(d))
			}
		}
		bk.ReleaseAll()
		rony.PoolMessageEnvelope.Put(m)
	})
}
