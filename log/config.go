package log

import "go.uber.org/zap/zapcore"

type Option func(cfg *config)

type config struct {
	level           Level
	TimeEncoder     TimeEncoder
	LevelEncoder    LevelEncoder
	DurationEncoder DurationEncoder
	CallerEncoder   CallerEncode
	sentryDSN       string
	sentryLevel     Level
	release         string
	environment     string
	skipCaller      int
	syslogTag       string
}

var defaultConfig = config{
	level:           InfoLevel,
	sentryDSN:       "",
	sentryLevel:     WarnLevel,
	release:         "",
	environment:     "",
	skipCaller:      1,
	TimeEncoder:     timeEncoder,
	LevelEncoder:    zapcore.CapitalLevelEncoder,
	DurationEncoder: zapcore.StringDurationEncoder,
	CallerEncoder:   zapcore.ShortCallerEncoder,
}

func WithLevel(lvl Level) Option {
	return func(cfg *config) {
		cfg.level = lvl
	}
}

func WithRelease(version string) Option {
	return func(cfg *config) {
		cfg.release = version
	}
}

func WithEnvironment(env string) Option {
	return func(cfg *config) {
		cfg.environment = env
	}
}

func WithSkipCaller(skip int) Option {
	return func(cfg *config) {
		cfg.skipCaller = skip
	}
}

func WithSentry(dsn string, lvl Level) Option {
	return func(cfg *config) {
		cfg.sentryDSN = dsn
		cfg.sentryLevel = lvl
	}
}

func WithSyslog(tag string) Option {
	return func(cfg *config) {
		cfg.syslogTag = tag
	}
}

func WithMongoDB(dsn string) Option {
	return func(cfg *config) {

	}
}
