# Infrastructure Adapters - Error Wrapper

## Overview

The Error Wrapper adapter provides implementations for error handling-related ports defined in the core domain and application layers. This adapter connects the application to error handling frameworks and libraries, following the Ports and Adapters (Hexagonal) architecture pattern. By isolating error handling implementations in adapter classes, the core business logic remains independent of specific error handling technologies, making the system more maintainable, testable, and flexible.

## Features

- Structured error handling
- Error categorization and classification
- Error wrapping and unwrapping
- Error translation between layers
- Error context enrichment
- Stack trace management
- Error logging integration
- Error reporting and monitoring

## Installation

```bash
go get github.com/abitofhelp/family-service/infrastructure/adapters/errorswrapper
```

## Configuration

The error wrapper can be configured according to specific requirements. Here's an example of configuring the error wrapper:

```
// Pseudocode example - not actual Go code
// This demonstrates how to configure and use an error wrapper

// 1. Import necessary packages
import errors, config, logging

// 2. Create a logger
logger = logging.NewLogger()

// 3. Configure the error wrapper
errorConfig = {
    includeStackTrace: true,
    maxStackDepth: 20,
    errorReporting: {
        enabled: true,
        sampleRate: 0.1,
        endpoint: "https://errors.example.com"
    },
    sensitiveFields: ["password", "token", "secret"]
}

// 4. Create the error wrapper
errorWrapper = errors.NewErrorWrapper(errorConfig, logger)

// 5. Use the error wrapper
err = someOperation()
if err != nil {
    // Wrap a low-level error with domain context
    domainErr = errorWrapper.Wrap(err, "failed to process family data")
    
    // Categorize the error
    if errorWrapper.IsNotFound(err) {
        // Handle not found case
    }
    
    // Log the error with context
    errorWrapper.LogError(context, domainErr)
}
```

## API Documentation

### Core Concepts

The error wrapper follows these core concepts:

1. **Adapter Pattern**: Implements error handling ports defined in the core domain or application layer
2. **Dependency Injection**: Receives dependencies through constructor injection
3. **Configuration**: Configured through a central configuration system
4. **Logging**: Uses a consistent logging approach
5. **Error Classification**: Categorizes errors into meaningful types

### Key Adapter Functions

```
// Pseudocode example - not actual Go code
// This demonstrates an error wrapper implementation

// Error wrapper structure
type ErrorWrapper {
    config        // Error wrapper configuration
    logger        // Logger for logging operations
    contextLogger // Context-aware logger
}

// Constructor for the error wrapper
function NewErrorWrapper(config, logger) {
    return new ErrorWrapper {
        config: config,
        logger: logger,
        contextLogger: new ContextLogger(logger)
    }
}

// Method to wrap an error with additional context
function ErrorWrapper.Wrap(err, message) {
    // Implementation would include:
    // 1. Creating a new error with the original as cause
    // 2. Adding the message
    // 3. Capturing stack trace if configured
    // 4. Preserving error type information
    // 5. Returning the wrapped error
}

// Method to check if an error is of a specific type
function ErrorWrapper.IsNotFound(err) {
    // Implementation would include:
    // 1. Unwrapping the error if needed
    // 2. Checking if it's a not found error
    // 3. Returning the result
}
```

## Best Practices

1. **Separation of Concerns**: Keep error handling logic separate from domain logic
2. **Interface Segregation**: Define focused error handling interfaces in the domain layer
3. **Dependency Injection**: Use constructor injection for adapter dependencies
4. **Error Classification**: Categorize errors into meaningful types
5. **Consistent Logging**: Use a consistent logging approach for errors
6. **Context Enrichment**: Add relevant context to errors
7. **Security**: Avoid including sensitive information in errors
8. **Testing**: Write unit and integration tests for error handling

## Troubleshooting

### Common Issues

#### Error Information Loss

If you encounter issues with error information being lost, consider the following:
- Ensure errors are properly wrapped rather than replaced
- Preserve stack traces when wrapping errors
- Use error types or codes to maintain categorization
- Add sufficient context when wrapping errors

#### Error Handling Performance

If error handling is impacting application performance, consider the following:
- Optimize stack trace capture
- Use sampling for error reporting
- Implement asynchronous error reporting
- Optimize error serialization
- Cache error templates or messages

## Related Components

- [Domain Layer](../../core/domain/README.md) - The domain layer that defines the error handling ports
- [Application Layer](../../core/application/README.md) - The application layer that uses error handling
- [Errors Package](../errors/README.md) - The errors package that defines error types

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.