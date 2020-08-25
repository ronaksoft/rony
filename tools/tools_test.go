package tools

import (
	"fmt"
	"runtime"
	"testing"
	"time"
)

/*
   Creation Time: 2019 - Nov - 25
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
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

func TestSanitizePhone(t *testing.T) {
	phones := map[string]string{
		"989121228718":  "989121228718",
		"+989121228718": "989121228718",
		"9121228718":    "989121228718",
	}

	for ph, cph := range phones {
		sph := SanitizePhone(ph, "IR")
		if sph != cph {
			t.Fatal()
		}
	}
}
