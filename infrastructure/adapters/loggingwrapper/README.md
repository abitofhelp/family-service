# Logging Wrapper

This package provides a wrapper around the `github.com/abitofhelp/servicelib/logging` package to ensure that the domain layer doesn't directly depend on external libraries. This follows the principles of Clean Architecture and Hexagonal Architecture (Ports and Adapters).

## Purpose

The purpose of this wrapper is to:

1. Isolate the domain layer from external dependencies
2. Provide a consistent logging approach throughout the application
3. Make it easier to replace or update the underlying logging library in the future

## Usage

Instead of directly importing `github.com/abitofhelp/servicelib/logging`, import this wrapper:

```go
import "github.com/abitofhelp/family-service/infrastructure/adapters/loggingwrapper"
```

Then use the wrapper types and functions:

```go
// Create a new logger
logger := loggingwrapper.NewLogger(zapLogger)

// Log messages
logger.Debug("Debug message", zap.String("key", "value"))
logger.Info("Info message", zap.String("key", "value"))
logger.Warn("Warning message", zap.String("key", "value"))
logger.Error("Error message", zap.String("key", "value"))
logger.Fatal("Fatal message", zap.String("key", "value"))

// Create a logger with additional fields
logger = logger.With(zap.String("key", "value"))

// Create a context logger
contextLogger := loggingwrapper.NewContextLogger(zapLogger)

// Log messages with context
contextLogger.Debug(ctx, "Debug message", zap.String("key", "value"))
contextLogger.Info(ctx, "Info message", zap.String("key", "value"))
contextLogger.Warn(ctx, "Warning message", zap.String("key", "value"))
contextLogger.Error(ctx, "Error message", zap.String("key", "value"))
contextLogger.Fatal(ctx, "Fatal message", zap.String("key", "value"))

// Get a logger from context
logger = loggingwrapper.FromContext(ctx)

// Add a logger to context
ctx = loggingwrapper.WithContext(ctx, logger)
```

## Components

### Logger

The `Logger` type is a wrapper around the zap logger:

```go
type Logger struct {
    logger *zap.Logger
}
```

### ContextLogger

The `ContextLogger` type is a wrapper around the servicelib context logger:

```go
type ContextLogger struct {
    logger *logging.ContextLogger
}
```

### Functions

The package provides the following functions:

- `NewLogger(logger *zap.Logger)`: Creates a new logger
- `(l *Logger) With(fields ...zap.Field)`: Returns a logger with the given fields
- `(l *Logger) Debug(msg string, fields ...zap.Field)`: Logs a debug message
- `(l *Logger) Info(msg string, fields ...zap.Field)`: Logs an info message
- `(l *Logger) Warn(msg string, fields ...zap.Field)`: Logs a warning message
- `(l *Logger) Error(msg string, fields ...zap.Field)`: Logs an error message
- `(l *Logger) Fatal(msg string, fields ...zap.Field)`: Logs a fatal message and exits
- `NewContextLogger(logger *zap.Logger)`: Creates a new context logger
- `FromContext(ctx context.Context)`: Returns a logger from the context
- `WithContext(ctx context.Context, logger *Logger)`: Returns a new context with the logger attached
- `(l *ContextLogger) Debug(ctx context.Context, msg string, fields ...zap.Field)`: Logs a debug message with context
- `(l *ContextLogger) Info(ctx context.Context, msg string, fields ...zap.Field)`: Logs an info message with context
- `(l *ContextLogger) Warn(ctx context.Context, msg string, fields ...zap.Field)`: Logs a warning message with context
- `(l *ContextLogger) Error(ctx context.Context, msg string, fields ...zap.Field)`: Logs an error message with context
- `(l *ContextLogger) Fatal(ctx context.Context, msg string, fields ...zap.Field)`: Logs a fatal message with context and exits