# Infrastructure Adapters - Logging Wrapper

## Overview

The Logging Wrapper adapter provides implementations for logging-related ports defined in the core domain and application layers. This adapter connects the application to logging frameworks and libraries, following the Ports and Adapters (Hexagonal) architecture pattern. By isolating logging implementations in adapter classes, the core business logic remains independent of specific logging technologies, making the system more maintainable, testable, and flexible.

## Features

- Structured logging support
- Log level management
- Context-aware logging
- Log formatting options
- Multiple output destinations
- Performance optimizations
- Log correlation (request ID, trace ID)
- Integration with various logging frameworks

## Installation

```bash
go get github.com/abitofhelp/family-service/infrastructure/adapters/loggingwrapper
```

## Configuration

The logging wrapper can be configured according to specific requirements. Here's an example of configuring the logging wrapper:

```
// Pseudocode example - not actual Go code
// This demonstrates how to configure and use a logging wrapper

// 1. Import necessary packages
import logging, config

// 2. Configure the logger
logConfig = {
    level: "info",
    format: "json",
    outputs: ["stdout", "/var/log/family-service.log"],
    includeTimestamp: true,
    includeCallerInfo: true,
    includeHostname: true,
    samplingRate: 1.0
}

// 3. Create the logger
logger = logging.NewLogger(logConfig)

// 4. Use the logger
logger.Info("Application started", {version: "1.0.0"})

// 5. Create a context-aware logger
contextLogger = logger.WithContext(context)
contextLogger.Debug("Processing request", {requestId: "req-123"})

// 6. Log errors
err = someOperation()
if err != nil {
    contextLogger.Error("Operation failed", {error: err})
}
```

## API Documentation

### Core Concepts

The logging wrapper follows these core concepts:

1. **Adapter Pattern**: Implements logging ports defined in the core domain or application layer
2. **Dependency Injection**: Receives dependencies through constructor injection
3. **Configuration**: Configured through a central configuration system
4. **Structured Logging**: Uses structured logging for better searchability and analysis
5. **Context Awareness**: Supports context-aware logging for request tracing

### Key Adapter Functions

```
// Pseudocode example - not actual Go code
// This demonstrates a logging wrapper implementation

// Logger structure
type Logger {
    config // Logger configuration
    writer // Log writer
}

// Constructor for the logger
function NewLogger(config) {
    writer = createLogWriter(config)
    return new Logger {
        config: config,
        writer: writer
    }
}

// Method to log at info level
function Logger.Info(message, fields) {
    // Implementation would include:
    // 1. Checking if info level is enabled
    // 2. Formatting the log entry
    // 3. Adding standard fields
    // 4. Writing to configured outputs
}

// Method to create a context-aware logger
function Logger.WithContext(context) {
    // Implementation would include:
    // 1. Creating a new context logger
    // 2. Extracting context information (request ID, trace ID)
    // 3. Returning the context logger
}
```

## Best Practices

1. **Separation of Concerns**: Keep logging logic separate from domain logic
2. **Interface Segregation**: Define focused logging interfaces in the domain layer
3. **Dependency Injection**: Use constructor injection for adapter dependencies
4. **Structured Logging**: Use structured logging with consistent field names
5. **Context Awareness**: Include context information in logs for request tracing
6. **Log Levels**: Use appropriate log levels for different types of information
7. **Performance**: Consider performance implications of logging in high-throughput systems

## Troubleshooting

### Common Issues

#### Log Configuration

If you encounter issues with log configuration, check the following:
- Log level is set appropriately
- Output destinations are valid and writable
- Log format is supported
- Configuration is properly loaded

#### Performance Impact

If logging is impacting application performance, consider the following:
- Reduce log verbosity in production
- Use sampling for high-volume log events
- Ensure debug logs are guarded by level checks
- Use asynchronous logging where appropriate
- Optimize log serialization

## Related Components

- [Domain Layer](../../core/domain/README.md) - The domain layer that defines the logging ports
- [Application Layer](../../core/application/README.md) - The application layer that uses logging
- [Interface Adapters](../../interface/adapters/README.md) - The interface adapters that use logging

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.