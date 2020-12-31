package tools

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"sync/atomic"
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

func TestNewFlusherPool(t *testing.T) {
	Convey("Flusher", t, func(c C) {
		var out, in int64
		f := NewFlusherPool(10, func(targetID string, entries []FlushEntry) {
			time.Sleep(time.Second)
			atomic.AddInt64(&out, int64(len(entries)))
		})

		wg := sync.WaitGroup{}
		total := 10000
		for i := 0; i < total; i++ {
			wg.Add(1)
			go func() {
				f.Enter(fmt.Sprintf("T%d", RandomInt(3)), RandomInt(10))
				atomic.AddInt64(&in, 1)
				wg.Done()
			}()
		}
		wg.Wait()
		time.Sleep(time.Second * 10)
		for _, q := range f.pool {
			c.So(q.entryChan, ShouldHaveLength, 0)
		}
		c.So(in, ShouldEqual, total)
		c.So(out, ShouldEqual, total)

	})
}
