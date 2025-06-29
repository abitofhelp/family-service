# Logging Wrapper

## Overview

The Logging Wrapper package provides a wrapper around the `github.com/abitofhelp/servicelib/logging` package to ensure that the domain layer doesn't directly depend on external libraries. This follows the principles of Clean Architecture and Hexagonal Architecture (Ports and Adapters), allowing the domain layer to remain isolated from external dependencies.

## Architecture

The Logging Wrapper package follows the Adapter pattern from Hexagonal Architecture, providing a layer of abstraction over the external `servicelib/logging` package. This ensures that the core domain doesn't directly depend on external libraries, maintaining the dependency inversion principle.

The package sits in the infrastructure layer of the application and is used by the domain layer through interfaces defined in the domain layer. The architecture follows these principles:

- **Dependency Inversion**: The domain layer depends on abstractions, not concrete implementations
- **Adapter Pattern**: This package adapts the external library to the domain's needs
- **Context Propagation**: Logging context is propagated through the application using context.Context

## Implementation Details

The Logging Wrapper package implements the following design patterns:

1. **Adapter Pattern**: Adapts the external library to the domain's needs
2. **Facade Pattern**: Provides a simplified interface to the underlying logging library
3. **Decorator Pattern**: Adds context-awareness to the logging functionality

Key implementation details:

- **Logger Wrapper**: The `Logger` type wraps the zap.Logger to provide a consistent interface
- **Context Logger**: The `ContextLogger` type adds context-awareness to logging
- **Context Integration**: Functions for adding and retrieving loggers from context
- **Structured Logging**: Support for structured logging with fields

The package uses the `github.com/abitofhelp/servicelib/logging` package internally but exposes its own API to the domain layer, ensuring that the domain layer doesn't directly depend on the external library.

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Family Service Example](../../../examples/family_service/README.md) - Shows how to use the logging wrapper

Example of using the logging wrapper:

```
// Create a new logger
logger := loggingwrapper.NewLogger(zapLogger)

// Log messages
logger.Debug("Debug message")
logger.Info("Info message")
logger.Warn("Warning message")
logger.Error("Error message")

// Create a logger with additional fields
logger = logger.With("key", "value")

// Create a context logger
contextLogger := loggingwrapper.NewContextLogger(zapLogger)

// Log messages with context
contextLogger.Debug(ctx, "Debug message")
contextLogger.Info(ctx, "Info message")
contextLogger.Warn(ctx, "Warning message")
contextLogger.Error(ctx, "Error message")

// Get a logger from context
logger = loggingwrapper.FromContext(ctx)

// Add a logger to context
ctx = loggingwrapper.WithContext(ctx, logger)
```

## Configuration

The Logging Wrapper package doesn't require any specific configuration itself, but it relies on the configuration of the underlying zap.Logger. The zap.Logger can be configured with various options:

- **Log Level**: Debug, Info, Warn, Error, Fatal
- **Output**: Console, File, or custom io.Writer
- **Format**: JSON or Console
- **Sampling**: Enable or disable sampling
- **Caller Info**: Include or exclude caller information
- **Stack Traces**: Include or exclude stack traces for errors

These configurations are typically provided when creating the zap.Logger that is passed to the `NewLogger` function.

## Testing

The Logging Wrapper package is tested through:

1. **Unit Tests**: Each function and method has unit tests
2. **Integration Tests**: Tests that verify the wrapper works correctly with the underlying library
3. **Context Tests**: Tests that verify context integration works correctly

Key testing approaches:

- **Mock Loggers**: Tests use mock loggers to verify logging behavior
- **Context Integration**: Tests verify that loggers can be added to and retrieved from context
- **Log Level Testing**: Tests verify that log level filtering works correctly
- **Structured Logging**: Tests verify that structured logging with fields works correctly

Example of a test case:

```
// Test that the logger logs messages at the correct level
func TestLogger_Info(t *testing.T) {
    // Create a buffer to capture log output
    var buf bytes.Buffer

    // Create a zap logger that writes to the buffer
    zapLogger := zap.New(
        zapcore.NewCore(
            zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
            zapcore.AddSync(&buf),
            zapcore.InfoLevel,
        ),
    )

    // Create a logger wrapper
    logger := loggingwrapper.NewLogger(zapLogger)

    // Log a message
    logger.Info("test message")

    // Verify the message was logged
    assert.Contains(t, buf.String(), "test message")
    assert.Contains(t, buf.String(), "\"level\":\"info\"")
}
```

## Design Notes

1. **Simplified API**: The wrapper provides a simplified API compared to the underlying library
2. **Context Integration**: The wrapper adds context integration to the logging functionality
3. **Structured Logging**: The wrapper supports structured logging with fields
4. **Log Levels**: The wrapper supports all standard log levels (Debug, Info, Warn, Error, Fatal)
5. **Dependency Inversion**: The package follows the Dependency Inversion Principle by ensuring that the domain layer depends on abstractions rather than concrete implementations
6. **Performance**: The wrapper is designed to add minimal overhead to logging operations

## API Documentation

### Logger

The `Logger` type is a wrapper around the zap logger:

```
type Logger struct {
    logger *zap.Logger
}
```

### ContextLogger

The `ContextLogger` type is a wrapper around the servicelib context logger:

```
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

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Zap Logger](https://github.com/uber-go/zap)
- [Go Context Package](https://golang.org/pkg/context/)
- [Domain Services](../../../core/domain/services/README.md) - Uses this logger for domain operations
- [Application Services](../../../core/application/services/README.md) - Uses this logger for application operations
- [GraphQL Resolvers](../../../interface/adapters/graphql/resolver/README.md) - Uses this logger for GraphQL operations
