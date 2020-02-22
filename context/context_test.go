package context

import (
	"git.ronaksoftware.com/ronak/rony/internal/tools"
	"hash/crc32"
	"testing"
)

/*
   Creation Time: 2019 - Oct - 14
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	randomKeys []string
)

func init() {
	for i := 0; i < 1000; i++ {
		randomKeys = append(randomKeys, tools.RandomID(32))
	}
}

func BenchmarkCrc(b *testing.B) {
	benchmarks := map[string]func(*testing.B){
		"Crc32":  benchCrc32,
		"String": benchString,
	}
	for name, bench := range benchmarks {
		b.Run(name, bench)
	}
}
func benchCrc32(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m := make(map[uint32]interface{})
			m[crc32.ChecksumIEEE(tools.StrToByte(tools.RandomID(24)))] = 10
		}
	})
}

func benchString(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		m := make(map[string]interface{})
		for pb.Next() {
			m[tools.RandomID(24)] = 10
		}
	})

}
