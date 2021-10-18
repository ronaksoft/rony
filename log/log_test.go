package log_test

import (
	"github.com/ronaksoft/rony/di"
	"github.com/ronaksoft/rony/log"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
	"testing"
)

func TestLog(t *testing.T) {
	Convey("Log", t, func(c C) {
		Convey("Different styles", func(c C) {
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
		Convey("Dependency Injection", func(c C) {
			log.ProvideDI(log.DefaultConfig)
			x := NewDI()
			x.Info("DI Hooorraa")

			log.ProvideDI(log.DefaultConfig)
			x2 := NewDI()
			x2.Info("DI Hooorraa2")
		})
	})
}

type DI struct {
	l log.Logger
}

func NewDI() *DI {
	d := &DI{}
	di.MustInvoke(func(l log.Logger) {
		d.l = l
	})

	return d
}

func (di DI) Info(msg string, fields ...log.Field) {
	di.l.Info(msg, fields...)
}
