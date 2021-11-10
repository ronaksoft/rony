package log

import (
	"log/syslog"

	"github.com/pkg/errors"
	"go.uber.org/zap/zapcore"
)

type syslogCore struct {
	zapcore.LevelEnabler
	encoder zapcore.Encoder
	w       *syslog.Writer
}

func NewSyslogCore(levelEnabler zapcore.LevelEnabler, encoder zapcore.Encoder, tag string) (*syslogCore, error) {
	c := &syslogCore{
		LevelEnabler: levelEnabler,
		encoder:      encoder,
	}
	var err error
	c.w, err = syslog.New(syslog.LOG_DAEMON, tag)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *syslogCore) clone() *syslogCore {
	return &syslogCore{
		LevelEnabler: s.LevelEnabler,
		encoder:      s.encoder.Clone(),
		w:            s.w,
	}
}

func (s *syslogCore) With(fields []zapcore.Field) zapcore.Core {
	clone := s.clone()
	for _, field := range fields {
		field.AddTo(clone.encoder)
	}

	return clone
}

func (s *syslogCore) Check(entry zapcore.Entry, checkedEntry *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if s.Enabled(entry.Level) {
		return checkedEntry.AddCore(entry, s)
	}

	return checkedEntry
}

func (s *syslogCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	// Generate the message.
	buffer, err := s.encoder.EncodeEntry(entry, fields)
	if err != nil {
		return errors.Wrap(err, "failed to encode log entry")
	}
	defer buffer.Free()

	message := buffer.String()

	// Write the message.
	switch entry.Level {
	case zapcore.DebugLevel:
		return s.w.Debug(message)

	case zapcore.InfoLevel:
		return s.w.Info(message)

	case zapcore.WarnLevel:
		return s.w.Warning(message)

	case zapcore.ErrorLevel:
		return s.w.Err(message)

	case zapcore.DPanicLevel:
		return s.w.Crit(message)

	case zapcore.PanicLevel:
		return s.w.Crit(message)

	case zapcore.FatalLevel:
		return s.w.Crit(message)

	default:
		return errors.Errorf("unknown log level: %v", entry.Level)
	}
}

func (s *syslogCore) Sync() error {
	return nil
}
