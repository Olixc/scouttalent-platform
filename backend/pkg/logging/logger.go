package logging

import (
	"context"
	"os"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

// NewLogger creates a new structured logger
func NewLogger(serviceName string) (*Logger, error) {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(getLogLevel()),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		InitialFields: map[string]interface{}{
			"service": serviceName,
		},
	}

	logger, err := config.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return nil, err
	}

	return &Logger{logger}, nil
}

// WithContext adds context information to the logger
func (l *Logger) WithContext(ctx context.Context) *zap.Logger {
	fields := []zap.Field{}

	if traceID := trace.SpanContextFromContext(ctx).TraceID(); traceID.IsValid() {
		fields = append(fields, zap.String("trace_id", traceID.String()))
	}

	return l.With(fields...)
}

func getLogLevel() zapcore.Level {
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}