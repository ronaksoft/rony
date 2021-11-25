package log

import (
	"io"

	"go.uber.org/zap/zapcore"
)

/*
   Creation Time: 2021 - Sep - 01
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	DefaultLogger *ronyLogger
	NopLogger     *ronyLogger
)

type (
	Level           = zapcore.Level
	Field           = zapcore.Field
	FieldType       = zapcore.FieldType
	CheckedEntry    = zapcore.CheckedEntry
	DurationEncoder = zapcore.DurationEncoder
	CallerEncode    = zapcore.CallerEncoder
	LevelEncoder    = zapcore.LevelEncoder
	TimeEncoder     = zapcore.TimeEncoder
	Encoder         = zapcore.Encoder
	Core            = zapcore.Core
)

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Check(Level, string) *CheckedEntry
	Sugared() *sugaredRonyLogger
	SetLevel(level Level)
	With(name string) Logger
	WithCore(enc Encoder, w io.Writer) Logger
}

type SugaredLogger interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Printf(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
}

func init() {
	DefaultLogger = New(
		WithSkipCaller(1),
	)

	NopLogger = newNOP()
}

func SetLevel(level Level) {
	DefaultLogger.SetLevel(level)
}

func With(name string) Logger {
	return DefaultLogger.With(name)
}
