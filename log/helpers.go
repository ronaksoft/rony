package log

import (
	"time"

	"go.uber.org/zap/zapcore"
)

/*
   Creation Time: 2019 - Aug - 29
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("06-01-02T15:04:05"))
}
