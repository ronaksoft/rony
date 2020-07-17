package tools

import "time"

/*
   Creation Time: 2020 - Apr - 09
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

var (
	timeInSec int64
)

func init() {
	go func() {
		for {
			timeInSec = time.Now().Unix()
			time.Sleep(time.Second)
		}
	}()
}

func TimeUnix() int64 {
	return timeInSec
}
