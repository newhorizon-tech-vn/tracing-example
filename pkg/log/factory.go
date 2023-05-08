package log

import (
	"context"
	"os"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Factory is the default logging wrapper that can create
// logger instances either for a given Context or context-less.
type Factory struct {
	logger *zap.Logger
}

type ConsoleConfiguration struct {
}

type FileConfiguration struct {
	Filename   string // "./chat.log"
	MaxSize    int    // megabytes
	MaxAge     int    // days
	MaxBackups int    // files
}

// If Console isnot nil, write log to std. If File isnot nil, write log to file.
type Configuration struct {
	JSONFormat      bool
	LogLevel        string
	StacktraceLevel string
	File            *FileConfiguration
	Console         *ConsoleConfiguration
}

// NewFactory creates a new Factory.
func NewFactory(config *Configuration) Factory {
	cores := []zapcore.Core{}
	level := getZapLevel(config.LogLevel) // log level

	// write to stdout
	if config.Console != nil {
		core := zapcore.NewCore(getEncoder(config.JSONFormat), zapcore.Lock(os.Stdout), level)
		cores = append(cores, core)
	}

	// write to file
	if config.File != nil {
		writer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.File.Filename,   // log file name
			MaxSize:    config.File.MaxSize,    // max size per file (in MB). Default 100MB
			MaxBackups: config.File.MaxBackups, // count backup file. Default keep all files
			MaxAge:     config.File.MaxAge,     // Max time file old log. Default keep forever
		})
		core := zapcore.NewCore(getEncoder(config.JSONFormat), writer, level)
		cores = append(cores, core)
	}

	combinedCore := zapcore.NewTee(cores...)

	stacktraceLevel := getZapLevel(config.StacktraceLevel) // trace level
	logger := zap.New(combinedCore,
		zap.AddCallerSkip(1),
		zap.AddCaller(),
		zap.AddStacktrace(stacktraceLevel),
	)

	return Factory{
		logger: logger,
	}
}

// Bg creates a context-unaware logger.
func (b Factory) Bg() Logger {
	return logger{logger: b.logger}
}

// For returns a context-aware Logger. If the context
// contains an OpenTracing span, all logging calls are also
// echo-ed into the span.
func (b Factory) For(ctx context.Context) Logger {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return b.Bg()
	}

	logger := logger{logger: b.logger}
	logger.spanFields = []zapcore.Field{
		zap.String("trace_id", span.SpanContext().TraceID().String()),
		zap.String("span_id", span.SpanContext().SpanID().String()),
	}

	return logger
}

// With creates a child logger, and optionally adds some context fields to that logger.
func (b Factory) With(fields ...zapcore.Field) Factory {
	return Factory{logger: b.logger.With(fields...)}
}

// Logger ...
func (b Factory) Logger() *zap.Logger {
	return b.logger
}

func getZapLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case DebugLevel:
		return zapcore.DebugLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	case PanicLevel:
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

func getEncoder(isJSON bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if isJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}
