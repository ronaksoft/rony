package pools

import (
	"sync"
	"time"
)

/*
   Creation Time: 2019 - Aug - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var timerPool sync.Pool

func AcquireTimer(timeout time.Duration) *time.Timer {
	tv := timerPool.Get()
	if tv == nil {
		return time.NewTimer(timeout)
	}

	t := tv.(*time.Timer)
	if t.Reset(timeout) {
		panic("BUG: Active timer trapped into acquireTimer()")
	}

	return t
}

func ReleaseTimer(t *time.Timer) {
	if !t.Stop() {
		// Collect possibly added time from the channel
		// if timer has been stopped and nobody collected its' value.
		select {
		case <-t.C:
		default:
		}
	}
	timerPool.Put(t)
}

func ResetTimer(t *time.Timer, period time.Duration) {
	if !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
	t.Reset(period)
}
