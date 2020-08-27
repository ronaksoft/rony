package tools

import (
	"sync/atomic"
	"time"
)

/*
   Creation Time: 2020 - Apr - 09
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	timeInSec int64
)

func init() {
	timeInSec = time.Now().Unix()
	go func() {
		for {
			time.Sleep(time.Second)
			atomic.AddInt64(&timeInSec, time.Now().Unix()-atomic.LoadInt64(&timeInSec))
		}
	}()
}

func TimeUnix() int64 {
	return atomic.LoadInt64(&timeInSec)
}
