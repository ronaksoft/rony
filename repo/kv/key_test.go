package kv

import (
	"github.com/ronaksoft/rony/tools"
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
		for pb.Next() {
			bk := NewBulkKey()
			d := bk.GenKey(tools.FastRand(), "tools.RandomID(10)", 3, 125)
			if len(d) != 38 {
				b.Fatal("invalid size", len(d))
			}
			bk.ReleaseAll()
		}
	})
}
