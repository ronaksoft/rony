package tools

import "time"

/*
   Creation Time: 2019 - Jul - 23
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	maxTry         = 10
	slowSleep      = time.Second
	fastSleep      = time.Millisecond * 100
	superFastSleep = time.Microsecond
)

type RetryableFunc func() error

func TrySlow(f RetryableFunc) error {
	return Try(maxTry, slowSleep, f)
}

func TryFast(f RetryableFunc) error {
	return Try(maxTry, fastSleep, f)
}

func TrySuperFast(f RetryableFunc) error {
	return Try(maxTry, superFastSleep, f)
}

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
