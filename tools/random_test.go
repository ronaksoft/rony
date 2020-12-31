package tools

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

/*
   Creation Time: 2021 - Jan - 01
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func BenchmarkRandomInt64(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = RandomInt64(0)
		}
	})
}

func BenchmarkRandomID(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = RandomID(10)
		}
	})
}

func TestRandomID(t *testing.T) {
	x := RandomID(10)
	fmt.Println(x)
	for i := 0; i < 1000; i++ {
		RandomID(10)
	}

	time.Sleep(time.Second)
	runtime.GC()
	fmt.Println(x)
}
