package kv

import (
	"time"
)

/*
   Creation Time: 2020 - Nov - 10
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

const (
	defaultConflictRetries = 100
	defaultMaxInterval     = time.Millisecond
)

type Config struct {
	DirPath             string
	ConflictRetries     int
	ConflictMaxInterval time.Duration
}

var (
	DefaultConfig = Config{
		DirPath:             "./_hdd",
		ConflictRetries:     100,
		ConflictMaxInterval: time.Millisecond,
	}
)
