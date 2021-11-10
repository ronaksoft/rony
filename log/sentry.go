package log

import (
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"
)

/*
   Creation Time: 2019 - Apr - 24
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type sentryCore struct {
	zapcore.LevelEnabler
	hub  *sentry.Hub
	tags map[string]string
}

func NewSentryCore(sentryDSN, release, environment string, level zapcore.Level, tags map[string]string) zapcore.Core {
	if len(sentryDSN) == 0 {
		return nil
	}
	client, err := sentry.NewClient(sentry.ClientOptions{
		Dsn:         sentryDSN,
		Release:     release,
		Environment: environment,
	})
	if err != nil {
		return zapcore.NewNopCore()
	}

	sentryScope := sentry.NewScope()
	sentryHub := sentry.NewHub(client, sentryScope)

	return &sentryCore{
		hub:          sentryHub,
		tags:         tags,
		LevelEnabler: level,
	}
}

func (c *sentryCore) With(fs []zapcore.Field) zapcore.Core {
	return &sentryCore{
		hub:          c.hub,
		tags:         c.tags,
		LevelEnabler: c.LevelEnabler,
	}
}

func (c *sentryCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.LevelEnabler.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}

	return ce
}

func (c *sentryCore) Write(ent zapcore.Entry, fs []zapcore.Field) error {
	m := make(map[string]interface{}, len(fs))
	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fs {
		f.AddTo(enc)
	}
	for k, v := range enc.Fields {
		m[k] = v
	}

	event := sentry.NewEvent()
	event.Message = ent.Message
	event.Timestamp = ent.Time
	event.Level = sentryLevel(ent.Level)
	event.Extra = m
	event.Tags = c.tags
	c.hub.CaptureEvent(event)

	// We may be crashing the program, so should flush any buffered events.
	if ent.Level > zapcore.ErrorLevel {
		c.hub.Flush(time.Second)
	}

	return nil
}

func (c *sentryCore) Sync() error {
	c.hub.Flush(time.Second * 3)

	return nil
}

func sentryLevel(lvl zapcore.Level) sentry.Level {
	switch lvl {
	case zapcore.DebugLevel:
		return sentry.LevelDebug
	case zapcore.InfoLevel:
		return sentry.LevelInfo
	case zapcore.WarnLevel:
		return sentry.LevelWarning
	case zapcore.ErrorLevel:
		return sentry.LevelError
	case zapcore.DPanicLevel:
		return sentry.LevelFatal
	case zapcore.PanicLevel:
		return sentry.LevelFatal
	case zapcore.FatalLevel:
		return sentry.LevelFatal
	default:
		// Unrecognized levels are fatal.
		return sentry.LevelFatal
	}
}
