package log

import "go.uber.org/zap"

/*
   Creation Time: 2019 - Aug - 29
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

func PanicOnError(guideText string, err error, args ...interface{}) {
	if err != nil {
		DefaultLogger.Fatal(guideText,
			zap.Error(err),
			zap.Any("Arguments", args),
		)
	}
}

func ErrorOnError(guideText string, err error, args ...interface{}) {
	if err != nil {
		DefaultLogger.Error(guideText,
			zap.Error(err),
			zap.Any("Arguments", args),
		)
	}
}

func WarnOnError(guideText string, err error, args ...interface{}) {
	if err != nil {
		DefaultLogger.Warn(guideText,
			zap.Error(err),
			zap.Any("Arguments", args),
		)
	}
}
