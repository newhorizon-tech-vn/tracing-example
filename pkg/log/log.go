package log

import (
	"context"

	"go.uber.org/zap/zapcore"
)

var (
	instance Factory
)

func InitLogger(consoleLevel, stacktraceLevel string) {
	if consoleLevel == "" {
		consoleLevel = DebugLevel
	}

	if stacktraceLevel == "" {
		stacktraceLevel = PanicLevel
	}

	instance = NewFactory(Configuration{
		EnableConsole:     true,
		ConsoleJSONFormat: true,
		ConsoleLevel:      consoleLevel,
		StacktraceLevel:   stacktraceLevel,
	})
}

// Inst ...
func Inst() Factory {
	return instance
}

// Bg creates a context-unaware logger.
func Bg() Logger {
	return instance.Bg()
}

// For returns a context-aware Logger. If the context
// contains an OpenTracing span, all logging calls are also
// echo-ed into the span.
func For(ctx context.Context) Logger {
	return instance.For(ctx)
}

// Debug logs an debig msg with fields
func Debug(msg string, fields ...zapcore.Field) {
	instance.Bg().Debug(msg, fields...)
}

// Info logs an info msg with fields
func Info(msg string, fields ...zapcore.Field) {
	instance.Bg().Info(msg, fields...)
}

// Warn logs an warn msg with fields
func Warn(msg string, fields ...zapcore.Field) {
	instance.Bg().Warn(msg, fields...)
}

// Error logs an error msg with fields
func Error(msg string, fields ...zapcore.Field) {
	instance.Bg().Error(msg, fields...)
}

// Fatal logs a fatal error msg with fields
func Fatal(msg string, fields ...zapcore.Field) {
	instance.Bg().Fatal(msg, fields...)
}

// Panic logs an panic msg with fields
func Panic(msg string, fields ...zapcore.Field) {
	instance.Bg().Panic(msg, fields...)
}
