package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
)

/*
   Creation Time: 2019 - Mar - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

var (
	atomicLevel   zap.AtomicLevel
	DefaultLogger *zapLogger
)

type (
	Level        = zapcore.Level
	Field        = zapcore.Field
	FieldType    = zapcore.FieldType
	CheckedEntry = zapcore.CheckedEntry
)

type Logger interface {
	Log(level Level, msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Check(Level, string) *CheckedEntry
	Sync() error
	SetLevel(level Level)
}

type Config struct {
	Level       Level
	DirPath     string
	Filename    string
	SentryDSN   string
	SentryLevel Level
	Release     string
	Environment string
}

var DefaultConfig = Config{
	Level:       InfoLevel,
	DirPath:     ".",
	Filename:    "",
	SentryDSN:   "",
	SentryLevel: WarnLevel,
	Release:     "",
	Environment: "",
}

func newZapLogger(core zapcore.Core, skip int) *zapLogger {
	l := new(zapLogger)
	l.Logger = zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(ErrorLevel),
		zap.AddCallerSkip(skip),
	)

	return l
}

func Init(config Config) Logger {
	atomicLevel = zap.NewAtomicLevelAt(config.Level)
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	cores := append([]zapcore.Core{}[:0],
		zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), atomicLevel),
	)

	logFile := filepath.Join(config.DirPath, config.Filename)
	if f, err := os.Create(logFile); err == nil {
		cores = append(cores,
			zapcore.NewCore(encoder, zapcore.AddSync(f), zap.NewAtomicLevelAt(WarnLevel)),
		)
	}

	sentryCore := NewSentryCore(config.SentryDSN, config.Release, config.Environment, config.SentryLevel, nil)
	if sentryCore != nil {
		cores = append(cores, sentryCore)
	}

	DefaultLogger = newZapLogger(zapcore.NewTee(cores...), 2)

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

type zapLogger struct {
	*zap.Logger
	zap.AtomicLevel
}

func (l *zapLogger) Log(level Level, msg string, fields ...Field) {
	l.Check(level, msg).Write(fields...)
}

func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.Logger.Debug(msg, fields...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.Logger.Info(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.Logger.Warn(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.Logger.Error(msg, fields...)
}

func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.Logger.Fatal(msg, fields...)
}

func (l *zapLogger) Check(level Level, msg string) *CheckedEntry {
	return l.Logger.Check(level, msg)
}

func (l *zapLogger) SetLevel(level Level) {
	l.AtomicLevel.SetLevel(level)
}
