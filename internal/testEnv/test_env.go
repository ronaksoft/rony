package testEnv

import (
	"github.com/ronaksoft/rony/internal/log"
	"os"
)

/*
   Creation Time: 2019 - Oct - 17
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	_Initialized bool
)

func Init() {
	if !_Initialized {
		_Initialized = true
	} else {
		return
	}

	log.Init(log.Config{
		Level:       log.WarnLevel,
		DirPath:     "",
		Filename:    "",
		SentryDSN:   "",
		SentryLevel: log.WarnLevel,
		Release:     "",
		Environment: "",
	})

	_ = os.RemoveAll("./_hdd")
	_ = os.MkdirAll("./_hdd", os.ModePerm)
}
