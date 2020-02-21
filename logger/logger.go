package log

import (
	"fmt"
	"git.ronaksoftware.com/ronak/rony/config"
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
   Copyright Ronak Software Group 2018
*/

var (
	atomicLevel   zap.AtomicLevel
	DefaultLogger *zapLogger
)

type Level = zapcore.Level
type Field = zapcore.Field
type CheckedEntry = zapcore.CheckedEntry

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

func NewConsoleLogger() *zapLogger {
	l := new(zapLogger)
	l.AtomicLevel = zap.NewAtomicLevel()
	consoleWriteSyncer := zapcore.Lock(os.Stdout)
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "zapLogger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	l.Logger = zap.New(
		zapcore.NewCore(consoleEncoder, consoleWriteSyncer, l.AtomicLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.AddStacktrace(ErrorLevel),
	)

	return l
}

func newFileLogger(filename string) *zapLogger {
	l := new(zapLogger)

	l.AtomicLevel = zap.NewAtomicLevelAt(DebugLevel)
	fileLog, _ := os.Create(filename)

	syncer := zapcore.Lock(fileLog)
	encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	l.Logger = zap.New(
		zapcore.NewCore(encoder, syncer, l.AtomicLevel),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
	)
	return l
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

func NewNop() *zapLogger {
	l := new(zapLogger)
	l.AtomicLevel = zap.NewAtomicLevel()
	l.Logger = zap.NewNop()
	return l
}

func InitLogger(logLevel zapcore.Level, sentryDSN string) Logger {
	atomicLevel = zap.NewAtomicLevelAt(logLevel)
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	cores := append([]zapcore.Core{}[:0],
		zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), atomicLevel),
	)

	logFile := filepath.Join(
		config.GetString(config.LogFilePath),
		fmt.Sprintf("%s-%s.log", config.GetString(config.BundleID), config.GetString(config.InstanceID)),
	)
	if f, err := os.Create(logFile); err == nil {
		cores = append(cores,
			zapcore.NewCore(encoder, zapcore.AddSync(f), zap.NewAtomicLevelAt(WarnLevel)),
		)
	}

	sentryCore, err := NewSentryCore(WarnLevel, nil)
	if err == nil {
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
