package log

import "go.uber.org/zap/zapcore"

type encoderBuilder struct {
	cfg zapcore.EncoderConfig
}

func EncoderBuilder() *encoderBuilder {
	return &encoderBuilder{
		cfg: zapcore.EncoderConfig{
			TimeKey:        "",
			LevelKey:       "",
			NameKey:        "",
			CallerKey:      "",
			MessageKey:     "",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeTime:     timeEncoder,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}
}

func (eb *encoderBuilder) WithTimeKey(key string) *encoderBuilder {
	eb.cfg.TimeKey = key

	return eb
}

func (eb *encoderBuilder) WithLevelKey(key string) *encoderBuilder {
	eb.cfg.LevelKey = key

	return eb
}

func (eb *encoderBuilder) WithNameKey(key string) *encoderBuilder {
	eb.cfg.NameKey = key

	return eb
}

func (eb *encoderBuilder) WithMessageKey(key string) *encoderBuilder {
	eb.cfg.MessageKey = key

	return eb
}

func (eb *encoderBuilder) WithCallerKey(key string) *encoderBuilder {
	eb.cfg.CallerKey = key

	return eb
}

func (eb *encoderBuilder) ConsoleEncoder() Encoder {
	return zapcore.NewConsoleEncoder(eb.cfg)
}

func (eb *encoderBuilder) JsonEncoder() Encoder {
	return zapcore.NewJSONEncoder(eb.cfg)
}
