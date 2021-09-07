package log

import (
	"fmt"
	"github.com/ronaksoft/rony/di"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"runtime/debug"
	"sync"
)

/*
   Creation Time: 2019 - Mar - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

type Config struct {
	Level           Level
	DirPath         string
	SentryDSN       string
	SentryLevel     Level
	Release         string
	Environment     string
	SkipCaller      int
	TimeEncoder     TimeEncoder
	LevelEncoder    LevelEncoder
	DurationEncoder DurationEncoder
	CallerEncoder   CallerEncode
}

var DefaultConfig = Config{
	Level:           InfoLevel,
	DirPath:         ".",
	SentryDSN:       "",
	SentryLevel:     WarnLevel,
	Release:         "",
	Environment:     "",
	SkipCaller:      2,
	TimeEncoder:     timeEncoder,
	LevelEncoder:    zapcore.CapitalLevelEncoder,
	DurationEncoder: zapcore.StringDurationEncoder,
	CallerEncoder:   zapcore.ShortCallerEncoder,
}

// ronyLogger is a wrapper around zap.Logger and adds a good few features to it.
// It provides layered logs which could be used by separate packages, and could be turned off or on
// separately. Separate layers could also have independent log levels.
// Whenever you change log level it propagates through its children.
type ronyLogger struct {
	prefix     string
	skipCaller int
	z          *zap.Logger
	sz         *zap.SugaredLogger
	lvl        zap.AtomicLevel
}

func New(cfg Config) *ronyLogger {
	l := &ronyLogger{
		lvl:        zap.NewAtomicLevelAt(cfg.Level),
		skipCaller: cfg.SkipCaller,
	}

	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    cfg.LevelEncoder,
		EncodeTime:     cfg.TimeEncoder,
		EncodeDuration: cfg.DurationEncoder,
		EncodeCaller:   cfg.CallerEncoder,
	})

	cores := append([]zapcore.Core{},
		zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), l.lvl),
	)

	sentryCore := NewSentryCore(cfg.SentryDSN, cfg.Release, cfg.Environment, cfg.SentryLevel, nil)
	if sentryCore != nil {
		cores = append(cores, sentryCore)
	}

	l.z = zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddStacktrace(ErrorLevel),
		zap.AddCallerSkip(cfg.SkipCaller),
	)
	l.sz = zap.New(
		l.z.Core(),
		zap.AddCaller(),
		zap.AddStacktrace(ErrorLevel),
		zap.AddCallerSkip(cfg.SkipCaller-1),
	).Sugar()

	return l
}

var (
	once sync.Once
)

// ProvideDI is protected by sync.Once and provides the Logger interface for other packages.
func ProvideDI(cfg Config) {
	once.Do(func() {
		di.MustProvide(func() Config {
			return cfg
		})
		di.MustProvide(New)
	})
}
func (l *ronyLogger) Sugared() *sugaredRonyLogger {
	return &sugaredRonyLogger{
		l: l,
	}
}

func (l *ronyLogger) Check(lvl Level, msg string) *CheckedEntry {
	return l.z.Check(lvl, msg)
}

func (l *ronyLogger) Sync() error {
	return l.z.Sync()
}

func (l *ronyLogger) SetLevel(lvl Level) {
	l.lvl.SetLevel(lvl)
}

func (l *ronyLogger) With(name string) Logger {
	return l.WithSkip(name, l.skipCaller)
}

func (l *ronyLogger) WithSkip(name string, skipCaller int) Logger {
	childLogger := &ronyLogger{
		prefix:     fmt.Sprintf("%s[%s]", l.prefix, name),
		skipCaller: l.skipCaller,
		z: zap.New(
			l.z.Core(),
			zap.AddCaller(),
			zap.AddStacktrace(ErrorLevel),
			zap.AddCallerSkip(skipCaller),
		),
		sz: zap.New(
			l.z.Core(),
			zap.AddCaller(),
			zap.AddStacktrace(ErrorLevel),
			zap.AddCallerSkip(skipCaller-1)).Sugar(),
		lvl: l.lvl,
	}

	return childLogger
}

func (l *ronyLogger) addPrefix(m string) string {
	if l.prefix != "" {
		return fmt.Sprintf("%s %s", l.prefix, m)
	}

	return m
}

func (l *ronyLogger) log(lvl Level, msg string, fields ...Field) {
	if ce := l.z.Check(lvl, l.addPrefix(msg)); ce != nil {
		ce.Write(fields...)
	}
}

func (l *ronyLogger) WarnOnErr(guideTxt string, err error, fields ...Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
		l.Warn(guideTxt, fields...)
	}
}

func (l *ronyLogger) ErrorOnErr(guideTxt string, err error, fields ...Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
		l.Error(guideTxt, fields...)
	}
}

func (l *ronyLogger) Debug(msg string, fields ...Field) {
	l.log(zap.DebugLevel, msg, fields...)
}

func (l *ronyLogger) Warn(msg string, fields ...Field) {
	l.log(zap.WarnLevel, msg, fields...)
}

func (l *ronyLogger) Info(msg string, fields ...Field) {
	l.log(zap.InfoLevel, msg, fields...)
}

func (l *ronyLogger) Error(msg string, fields ...Field) {
	l.log(zap.ErrorLevel, msg, fields...)
}

func (l *ronyLogger) Fatal(msg string, fields ...Field) {
	l.log(zap.FatalLevel, msg, fields...)
}

func (l *ronyLogger) RecoverPanic(funcName string, extraInfo interface{}, compensationFunc func()) {
	if r := recover(); r != nil {
		l.Error("Panic Recovered",
			zap.String("Func", funcName),
			zap.Any("Info", extraInfo),
			zap.Any("Recover", r),
			zap.ByteString("StackTrace", debug.Stack()),
		)
		if compensationFunc != nil {
			go compensationFunc()
		}
	}
}

type sugaredRonyLogger struct {
	l *ronyLogger
}

func (l sugaredRonyLogger) Debugf(template string, args ...interface{}) {
	l.l.sz.Debugf(l.l.addPrefix(template), args...)
}

func (l sugaredRonyLogger) Infof(template string, args ...interface{}) {
	l.l.sz.Infof(l.l.addPrefix(template), args...)
}

func (l sugaredRonyLogger) Printf(template string, args ...interface{}) {
	fmt.Printf(template, args...)
}

func (l sugaredRonyLogger) Warnf(template string, args ...interface{}) {
	l.l.sz.Warnf(l.l.addPrefix(template), args...)
}

func (l sugaredRonyLogger) Errorf(template string, args ...interface{}) {
	l.l.sz.Errorf(l.l.addPrefix(template), args...)
}

func (l sugaredRonyLogger) Fatalf(template string, args ...interface{}) {
	l.l.sz.Fatalf(l.l.addPrefix(template), args...)
}

func (l sugaredRonyLogger) Debug(args ...interface{}) {
	l.l.sz.Debug(args...)
}

func (l sugaredRonyLogger) Info(args ...interface{}) {
	l.l.sz.Info(args...)
}

func (l sugaredRonyLogger) Warn(args ...interface{}) {
	l.l.sz.Warn(args...)
}

func (l sugaredRonyLogger) Error(args ...interface{}) {
	l.l.sz.Error(args...)
}

func (l sugaredRonyLogger) Fatal(args ...interface{}) {
	l.l.sz.Fatal(args...)
}

func (l sugaredRonyLogger) Panic(args ...interface{}) {
	l.l.sz.Panic(args...)
}
