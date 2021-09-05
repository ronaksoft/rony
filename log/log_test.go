package log_test

import (
	"github.com/ronaksoft/rony/log"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
	"testing"
)

func TestLog(t *testing.T) {
	Convey("Log", t, func(c C) {
		log.Info("Info")
		l := log.New(log.DefaultConfig)
		l.SetLevel(log.DebugLevel)
		l.Info("INFO", zap.Int("Index", 1))
		l.Sugared().Info("INFO", 1)
		l.Sugared().Infof("INFO %d", 1)
		ll := l.With("CHILD")
		ll.SetLevel(log.DebugLevel)
		ll.Info("INFO", zap.Int("Index", 1))
		ll.Sugared().Info("INFO", 1)
		ll.Sugared().Infof("INFO %d", 1)
	})
}
