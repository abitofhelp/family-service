# Infrastructure Adapters - Errors

## Overview

The Errors adapter provides implementations for error-related ports defined in the core domain and application layers. This adapter defines and manages domain-specific errors and error types, following the Ports and Adapters (Hexagonal) architecture pattern. By isolating error implementations in adapter classes, the core business logic can work with well-defined error types, making the system more maintainable, testable, and flexible.

## Features

- Domain-specific error types
- Error code management
- Error categorization
- Error hierarchies
- Internationalization support for error messages
- Error serialization and deserialization
- HTTP status code mapping
- Error formatting

## Installation

```bash
go get github.com/abitofhelp/family-service/infrastructure/adapters/errors
```

## Configuration

The errors adapter can be configured according to specific requirements. Here's an example of configuring and using the errors adapter:

```
// Pseudocode example - not actual Go code
// This demonstrates how to configure and use the errors adapter

// 1. Import necessary packages
import errors, config, logging

// 2. Create a logger
logger = logging.NewLogger()

// 3. Configure the errors
errorsConfig = {
    defaultLanguage: "en",
    supportedLanguages: ["en", "es", "fr"],
    errorMessagePath: "./errors/messages",
    includeErrorCodes: true
}

// 4. Create the errors factory
errorsFactory = errors.NewErrorsFactory(errorsConfig, logger)

// 5. Use the errors factory to create domain errors
notFoundErr = errorsFactory.NewNotFoundError("family", "123", "Family not found")
validationErr = errorsFactory.NewValidationError("name", "Name is required")
unauthorizedErr = errorsFactory.NewUnauthorizedError("Invalid credentials")

// 6. Check error types
if errors.IsNotFound(err) {
    // Handle not found case
}

// 7. Get localized error message
spanishMessage = errors.GetLocalizedMessage(err, "es")
```

## API Documentation

### Core Concepts

The errors adapter follows these core concepts:

1. **Domain-Specific Errors**: Defines error types that are meaningful in the domain context
2. **Error Categorization**: Categorizes errors into meaningful types (validation, not found, etc.)
3. **Error Factory**: Provides a factory for creating domain-specific errors
4. **Internationalization**: Supports localized error messages
5. **Error Codes**: Associates unique codes with error types for easier identification

### Key Adapter Functions

```
// Pseudocode example - not actual Go code
// This demonstrates an errors adapter implementation

// Errors factory structure
type ErrorsFactory {
    config        // Errors configuration
    logger        // Logger for logging operations
    contextLogger // Context-aware logger
    messages      // Localized error messages
}

// Constructor for the errors factory
function NewErrorsFactory(config, logger) {
    messages = loadErrorMessages(config.errorMessagePath)
    return new ErrorsFactory {
        config: config,
        logger: logger,
        contextLogger: new ContextLogger(logger),
        messages: messages
    }
}

// Method to create a not found error
function ErrorsFactory.NewNotFoundError(entity, id, message) {
    // Implementation would include:
    // 1. Creating a new not found error
    // 2. Setting the entity and ID
    // 3. Setting the error code
    // 4. Setting the message
    // 5. Returning the error
}

// Function to check if an error is a not found error
function IsNotFound(err) {
    // Implementation would include:
    // 1. Checking if the error is of the not found type
    // 2. Returning the result
}
```

## Best Practices

1. **Separation of Concerns**: Keep error definitions separate from domain logic
2. **Interface Segregation**: Define focused error interfaces in the domain layer
3. **Error Categorization**: Categorize errors into meaningful types
4. **Internationalization**: Support localized error messages
5. **Error Codes**: Use unique codes for error types
6. **Consistency**: Maintain consistent error formats across the application
7. **Testing**: Write unit tests for error handling
8. **Documentation**: Document error types and their meanings

## Troubleshooting

### Common Issues

#### Error Type Identification

If you encounter issues with error type identification, consider the following:
- Use type assertions or error wrapping to preserve error types
- Implement Is/As methods for custom error types
- Use error codes for easier identification
- Avoid creating new error instances when wrapping errors

#### Localization Issues

If you encounter issues with error message localization, consider the following:
- Ensure message files exist for all supported languages
- Verify the message keys match between language files
- Implement fallback mechanisms for missing translations
- Use placeholders consistently in message templates

## Related Components

- [Domain Layer](../../core/domain/README.md) - The domain layer that defines error interfaces
- [Application Layer](../../core/application/README.md) - The application layer that uses domain errors
- [Error Wrapper](../errorswrapper/README.md) - The error wrapper that uses these error types

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.