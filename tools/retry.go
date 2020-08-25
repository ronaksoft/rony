package tools

import "time"

/*
   Creation Time: 2019 - Jul - 23
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2018
*/

type RetryableFunc func() error

func Try(attempts int, waitTime time.Duration, f RetryableFunc) (err error) {
	for attempts > 0 {
		err = f()
		if err == nil {
			break
		}
		attempts--
		time.Sleep(waitTime)
	}
	return
}
