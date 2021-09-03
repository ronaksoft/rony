package log_test

import (
	"github.com/ronaksoft/rony/log"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLog(t *testing.T) {
	Convey("Log", t, func(c C) {
		l := log.New(log.DefaultConfig)
		l.SetLevel(log.DebugLevel)
		l.Info("INFO1")
		l.Sugared().Info("INFO", 1)
		l.Sugared().Infof("INFO %d", 1)

	})
}
