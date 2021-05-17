package localdb

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
	defaultBatchWorkers    = 16
	defaultBatchSize       = 256
)

type Config struct {
	DirPath             string
	ConflictRetries     int
	ConflictMaxInterval time.Duration
	BatchWorkers        int
	BatchSize           int
}

func DefaultConfig(dataPath string) Config {
	return Config{
		DirPath:             dataPath,
		ConflictRetries:     100,
		ConflictMaxInterval: time.Millisecond,
		BatchSize:           defaultBatchSize,
		BatchWorkers:        defaultBatchWorkers,
	}
}
