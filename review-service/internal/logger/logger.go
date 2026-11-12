package logger

import (
	"context"
	"fmt"
	"io"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ctxLoggerKey struct{}

func New(level string, w io.Writer, options ...zap.Option) (*zap.SugaredLogger, error) {
	zapLevel, err := parseLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.AddSync(w),
		zapLevel,
	)

	l := zap.New(
		core,
		append(options, zap.AddCaller())...,
	)

	return l.Sugar(), nil
}

func WithLogger(ctx context.Context, l *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, l)
}

func FromContext(ctx context.Context) *zap.SugaredLogger {
	l, ok := ctx.Value(ctxLoggerKey{}).(*zap.SugaredLogger)
	if ok && l != nil {
		return l
	}

	return zap.NewNop().Sugar()
}

func parseLevel(level string) (zapcore.Level, error) {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel, nil

	case "info":
		return zapcore.InfoLevel, nil

	case "warn":
		return zapcore.WarnLevel, nil

	case "error":
		return zapcore.ErrorLevel, nil

	default:
		return zapcore.InfoLevel,
			fmt.Errorf("invalid log level: %s", level)
	}
}
