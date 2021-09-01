package log

import (
	"go.uber.org/zap"
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
	atomicLevel   zap.AtomicLevel
	DefaultLogger *ronyLogger
)

type (
	Level        = zapcore.Level
	Field        = zapcore.Field
	FieldType    = zapcore.FieldType
	CheckedEntry = zapcore.CheckedEntry
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

func Init(config Config) Logger {
	DefaultLogger = New(config)
	return DefaultLogger
}

func Debug(msg string, fields ...Field) {
	DefaultLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	DefaultLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	DefaultLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	DefaultLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	DefaultLogger.Fatal(msg, fields...)
}

func Check(level Level, msg string) *CheckedEntry {
	return DefaultLogger.Check(level, msg)
}

func SetLevel(level Level) {
	atomicLevel.SetLevel(level)
}
