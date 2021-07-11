package tools_test

import (
	"github.com/ronaksoft/rony/tools"
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"testing"
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
			_ = tools.RandomInt64(0)
		}
	})
}

func BenchmarkRandomID(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = tools.RandomID(10)
		}
	})
}

func TestSecureRandomInt63(t *testing.T) {
	t.Parallel()
	Convey("SecureRandom", t, func(c C) {
		size := 10000
		mtx := tools.SpinLock{}
		wg := sync.WaitGroup{}
		m := make(map[int64]struct{}, size)
		for i := 0; i < size; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				mtx.Lock()
				x := tools.SecureRandomInt63(0)
				_, ok := m[x]
				c.So(ok, ShouldBeFalse)
				m[x] = struct{}{}
				mtx.Unlock()
			}()
		}
		wg.Wait()
	})
}
