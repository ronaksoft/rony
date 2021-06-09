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
		Convey("Without WaitTime", func(c C) {
			var out, in int64
			f := NewFlusherPool(10, 20, func(targetID string, entries []FlushEntry) {
				time.Sleep(time.Millisecond * 100)
				atomic.AddInt64(&out, int64(len(entries)))
			})

			wg := sync.WaitGroup{}
			total := 10000
			for i := 0; i < total; i++ {
				wg.Add(1)
				go func() {
					f.EnterAndWait(fmt.Sprintf("T%d", RandomInt(3)), NewEntry(RandomInt(10)))
					atomic.AddInt64(&in, 1)
					wg.Done()
				}()
			}
			wg.Wait()
			for _, q := range f.pool {
				c.So(q.entryChan, ShouldHaveLength, 0)
			}
			c.So(in, ShouldEqual, total)
			c.So(out, ShouldEqual, total)
		})
		Convey("With WaitTime", func(c C) {
			var out, in int64
			f := NewFlusherPoolWithWaitTime(10, 20, 250*time.Millisecond, func(targetID string, entries []FlushEntry) {
				time.Sleep(time.Millisecond * 100)
				atomic.AddInt64(&out, int64(len(entries)))
			})

			wg := sync.WaitGroup{}
			total := 10000
			for i := 0; i < total; i++ {
				wg.Add(1)
				go func() {
					f.EnterAndWait(fmt.Sprintf("T%d", RandomInt(3)), NewEntry(RandomInt(10)))
					atomic.AddInt64(&in, 1)
					wg.Done()
				}()
			}
			wg.Wait()
			for _, q := range f.pool {
				c.So(q.entryChan, ShouldHaveLength, 0)
			}
			c.So(in, ShouldEqual, total)
			c.So(out, ShouldEqual, total)
		})

	})
}
