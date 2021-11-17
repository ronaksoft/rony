package log_test

import (
	"os"
	"testing"

	"github.com/ronaksoft/rony/di"
	"github.com/ronaksoft/rony/log"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
)

func TestLog(t *testing.T) {
	Convey("Log", t, func(c C) {
		Convey("Different styles", func(c C) {
			l := log.New()
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
			log.ProvideDI()
			x := NewDI()
			x.Info("DI Hooorraa")

			log.ProvideDI()
			x2 := NewDI()
			x2.Info("DI Hooorraa2")
		})
		Convey("Logger with SysLog", func(c C) {
			l := log.New(log.WithSyslog("rony-logger"))
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
	})
}

func TestLogWithWriter(t *testing.T) {
	Convey("Log With Added Writer", t, func(c C) {
		Convey("File Writer", func(c C) {
			_ = os.MkdirAll("./_hdd", os.ModePerm)
			l := log.New()
			l.Warn("Log1")
			l.Info("Log2")
			f, err := os.OpenFile(
				"./_hdd/testlog.txt",
				os.O_RDWR|os.O_CREATE|os.O_APPEND,
				0666,
			)
			c.So(err, ShouldBeNil)
			l2 := l.WithCore(log.EncoderBuilder().ConsoleEncoder(), f)
			l2.Warn("File Log1")
			l2.Info("File Log2")
			err = f.Close()
			c.So(err, ShouldBeNil)
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
