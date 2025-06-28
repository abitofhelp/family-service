// Copyright (c) 2025 A Bit of Help, Inc.

// Package loggingwrapper provides a wrapper around servicelib/logging to ensure
// the domain layer doesn't directly depend on external libraries.
package loggingwrapper

import (
	"context"

	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

// contextLoggerKey is the key used to store and retrieve the logger from the context
type contextLoggerKeyType struct{}

var contextLoggerKey = contextLoggerKeyType{}

// Logger is a wrapper around the servicelib logger
type Logger struct {
	logger *zap.Logger
}

// NewLogger creates a new logger
func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

// With returns a logger with the given fields
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		logger: l.logger.With(fields...),
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Info logs an info message
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

// Error logs an error message
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

// ContextLogger is a wrapper around the servicelib context logger
type ContextLogger struct {
	logger *logging.ContextLogger
	zapLogger *zap.Logger
}

// NewContextLogger creates a new context logger
func NewContextLogger(logger *zap.Logger) *ContextLogger {
	return &ContextLogger{
		logger: logging.NewContextLogger(logger),
		zapLogger: logger,
	}
}

// Logger returns the underlying zap logger
func (l *ContextLogger) Logger() *zap.Logger {
	return l.zapLogger
}

// FromContext returns a logger from the context
func (l *ContextLogger) FromContext(ctx context.Context) *Logger {
	// Try to get logger from context using a context key
	loggerValue := ctx.Value(contextLoggerKey)
	if loggerValue == nil {
		// Return a no-op logger if none is found in the context
		return &Logger{
			logger: zap.NewNop(),
		}
	}

	// Try to cast to *zap.Logger
	if zapLogger, ok := loggerValue.(*zap.Logger); ok {
		return &Logger{
			logger: zapLogger,
		}
	}

	// Try to cast to *Logger
	if logger, ok := loggerValue.(*Logger); ok {
		return logger
	}

	// Return a no-op logger if the value is not a logger
	return &Logger{
		logger: zap.NewNop(),
	}
}

// WithContext returns a new context with the logger attached
func (l *ContextLogger) WithContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, contextLoggerKey, logger.logger)
}

// Debug logs a debug message with context
func (l *ContextLogger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Debug(ctx, msg, fields...)
}

// Info logs an info message with context
func (l *ContextLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Info(ctx, msg, fields...)
}

// Warn logs a warning message with context
func (l *ContextLogger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Warn(ctx, msg, fields...)
}

// Error logs an error message with context
func (l *ContextLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Error(ctx, msg, fields...)
}

// Fatal logs a fatal message with context and exits
func (l *ContextLogger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.Fatal(ctx, msg, fields...)
}
