package log

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"sync"

	"github.com/ronaksoft/rony/di"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*
   Creation Time: 2019 - Mar - 02
   Created by:  (ehsan)
   Maintainers:
      1.  Ehsan N. Moosa (E2)
   Auditor: Ehsan N. Moosa (E2)
   Copyright Ronak Software Group 2020
*/

// ronyLogger is a wrapper around zap.Logger and adds a good few features to it.
// It provides layered logs which could be used by separate packages, and could be turned off or on
// separately. Separate layers could also have independent log levels.
// Whenever you change log level it propagates through its children.
type ronyLogger struct {
	prefix     string
	skipCaller int
	encoder    zapcore.Encoder
	z          *zap.Logger
	sz         *zap.SugaredLogger
	lvl        zap.AtomicLevel
}

func New(opts ...Option) *ronyLogger {
	cfg := defaultConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	l := &ronyLogger{
		lvl:        zap.NewAtomicLevelAt(cfg.level),
		skipCaller: cfg.skipCaller,
	}

	l.encoder = EncoderBuilder().
		WithTimeKey("ts").
		WithLevelKey("level").
		WithNameKey("name").
		WithCallerKey("caller").
		WithMessageKey("msg").
		ConsoleEncoder()

	cores := append([]zapcore.Core{},
		zapcore.NewCore(l.encoder, zapcore.Lock(os.Stdout), l.lvl),
	)

	if cfg.syslogTag != "" {
		syslogCore, err := NewSyslogCore(l.lvl, l.encoder, cfg.syslogTag)
		if err != nil {
			fmt.Println("got error on enabling syslog:", err)
		} else {
			cores = append(cores, syslogCore)
		}
	}

	if cfg.sentryDSN != "" {
		sentryCore := NewSentryCore(cfg.sentryDSN, cfg.release, cfg.environment, cfg.sentryLevel, nil)
		if sentryCore != nil {
			cores = append(cores, sentryCore)
		}
	}

	l.z = zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddStacktrace(ErrorLevel),
		zap.AddCallerSkip(cfg.skipCaller),
	)

	l.sz = l.z.Sugar()

	return l
}

func newNOP() *ronyLogger {
	l := &ronyLogger{}
	l.z = zap.NewNop()
	l.sz = l.z.Sugar()

	return l
}

var (
	once sync.Once
)

// ProvideDI is protected by sync.Once and provides the Logger interface for other packages.
func ProvideDI(opts ...Option) {
	once.Do(func() {
		di.MustProvide(
			func() Logger {
				return New(opts...)
			},
		)
	})
}

func (l *ronyLogger) Sugared() *sugaredRonyLogger {
	return &sugaredRonyLogger{
		l: l,
	}
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
	return l.with(l.z.Core(), name, skipCaller)
}

func (l *ronyLogger) WithCore(enc Encoder, w io.Writer) Logger {
	core := zapcore.NewTee(
		l.z.Core(),
		zapcore.NewCore(enc, zapcore.AddSync(w), l.lvl),
	)

	return l.with(core, "", l.skipCaller)
}

func (l *ronyLogger) with(core zapcore.Core, name string, skip int) Logger {
	prefix := l.prefix
	if name != "" {
		prefix = fmt.Sprintf("%s[%s]", l.prefix, name)
	}
	childLogger := &ronyLogger{
		prefix:     prefix,
		skipCaller: l.skipCaller,
		encoder:    l.encoder.Clone(),
		z: zap.New(
			core,
			zap.AddCaller(),
			zap.AddStacktrace(ErrorLevel),
			zap.AddCallerSkip(skip),
		),
		sz: zap.New(
			core,
			zap.AddCaller(),
			zap.AddStacktrace(ErrorLevel),
			zap.AddCallerSkip(skip)).Sugar(),
		lvl: l.lvl,
	}

	return childLogger
}

func (l *ronyLogger) addPrefix(in string) (out string) {
	if l.prefix != "" {
		sb := &strings.Builder{}
		sb.WriteString(l.prefix)
		sb.WriteRune(' ')
		sb.WriteString(in)
		out = sb.String()

		return out
	}

	return in
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

func (l *ronyLogger) checkLevel(lvl Level) bool {
	// Check the level first to reduce the cost of disabled log calls.
	// Since Panic and higher may exit, we skip the optimization for those levels.
	if lvl < zapcore.DPanicLevel && !l.z.Core().Enabled(lvl) {
		return false
	}

	return true
}

func (l *ronyLogger) Check(lvl Level, msg string) *CheckedEntry {
	if !l.checkLevel(lvl) {
		return nil
	}

	return l.z.Check(lvl, l.addPrefix(msg))
}

func (l *ronyLogger) Debug(msg string, fields ...Field) {
	if !l.checkLevel(DebugLevel) {
		return
	}
	if ce := l.z.Check(DebugLevel, l.addPrefix(msg)); ce != nil {
		ce.Write(fields...)
	}
}

func (l *ronyLogger) Info(msg string, fields ...Field) {
	if !l.checkLevel(InfoLevel) {
		return
	}
	if ce := l.z.Check(InfoLevel, l.addPrefix(msg)); ce != nil {
		ce.Write(fields...)
	}
}

func (l *ronyLogger) Warn(msg string, fields ...Field) {
	if !l.checkLevel(WarnLevel) {
		return
	}
	if ce := l.z.Check(WarnLevel, l.addPrefix(msg)); ce != nil {
		ce.Write(fields...)
	}
}

func (l *ronyLogger) Error(msg string, fields ...Field) {
	if !l.checkLevel(ErrorLevel) {
		return
	}
	if ce := l.z.Check(ErrorLevel, l.addPrefix(msg)); ce != nil {
		ce.Write(fields...)
	}
}

func (l *ronyLogger) Fatal(msg string, fields ...Field) {
	l.z.Fatal(l.addPrefix(msg), fields...)
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
