// Copyright (c) 2019 The Jaeger Authors.
// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a simplified abstraction of the zap.Logger
type Logger interface {
	With(fields ...zapcore.Field) *zap.Logger
	Debug(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
	Panic(msg string, fields ...zapcore.Field)
}

// logger delegates all calls to the underlying zap.Logger
type logger struct {
	logger     *zap.Logger
	spanFields []zapcore.Field
}

func (l logger) With(fields ...zapcore.Field) *zap.Logger {
	return l.logger.With(l.spanFields...)
}

// Debug logs an debig msg with fields
func (l logger) Debug(msg string, fields ...zapcore.Field) {
	l.logger.With(l.spanFields...).Debug(msg, fields...)
}

// Info logs an info msg with fields
func (l logger) Info(msg string, fields ...zapcore.Field) {
	l.logger.With(l.spanFields...).Info(msg, fields...)
}

// Warn logs an warn msg with fields
func (l logger) Warn(msg string, fields ...zapcore.Field) {
	l.logger.With(l.spanFields...).Warn(msg, fields...)
}

// Error logs an error msg with fields
func (l logger) Error(msg string, fields ...zapcore.Field) {
	l.logger.With(l.spanFields...).Error(msg, fields...)
}

// Fatal logs a fatal error msg with fields
func (l logger) Fatal(msg string, fields ...zapcore.Field) {
	l.logger.With(l.spanFields...).Fatal(msg, fields...)
}

// Panic logs an panic msg with fields
func (l logger) Panic(msg string, fields ...zapcore.Field) {
	l.logger.With(l.spanFields...).Panic(msg, fields...)
}
