# Error Wrapper

## Overview

This package provides a wrapper around the `github.com/abitofhelp/servicelib/errors` package to ensure that the domain layer doesn't directly depend on external libraries. This follows the principles of Clean Architecture and Hexagonal Architecture (Ports and Adapters).

The purpose of this wrapper is to:

1. Isolate the domain layer from external dependencies
2. Provide a consistent error handling approach throughout the application
3. Make it easier to replace or update the underlying error handling library in the future

## Architecture

The `errorswrapper` package follows the Adapter pattern from Hexagonal Architecture, providing a layer of abstraction over the external `servicelib/errors` package. This ensures that the core domain doesn't directly depend on external libraries, maintaining the dependency inversion principle.

The package sits in the infrastructure layer of the application and is used by the domain layer through interfaces defined in the domain layer.

## Usage

Instead of directly importing `github.com/abitofhelp/servicelib/errors`, import this wrapper:

```go
import "github.com/abitofhelp/family-service/infrastructure/adapters/errorswrapper"
```

Then use the wrapper functions:

```go
// Create a domain error
err := errorswrapper.NewDomainError("ERROR_CODE", "Error message", cause)

// Create a validation error
err := errorswrapper.NewValidationError("Error message", "field", cause)

// Create a database error
err := errorswrapper.NewDatabaseError("Error message", "operation", "table", cause)

// Create a not found error
err := errorswrapper.NewNotFoundError("ResourceType", "resourceID", cause)

// Check error types
if errorswrapper.IsValidationError(err) {
    // Handle validation error
}

// Get error details
code := errorswrapper.GetErrorCode(err)
message := errorswrapper.GetErrorMessage(err)
cause := errorswrapper.GetErrorCause(err)
```

## Implementation Details

The wrapper provides the following error types:

- `Error`: Base error interface
- `ValidationError`: For validation errors
- `DomainError`: For domain-specific errors
- `DatabaseError`: For database-related errors
- `NotFoundError`: For resource not found errors

The wrapper provides the following functions:

- `NewDomainError`: Creates a new domain error
- `NewValidationError`: Creates a new validation error
- `NewDatabaseError`: Creates a new database error
- `NewNotFoundError`: Creates a new not found error
- `IsValidationError`: Checks if an error is a validation error
- `IsDomainError`: Checks if an error is a domain error
- `IsDatabaseError`: Checks if an error is a database error
- `IsNotFoundError`: Checks if an error is a not found error
- `GetErrorCode`: Gets the error code
- `GetErrorMessage`: Gets the error message
- `GetErrorCause`: Gets the error cause
- `FormatError`: Formats an error with its code, message, and cause
- `WrapError`: Wraps an error with a message

## Examples

```go
// Create a domain error
err := errorswrapper.NewDomainError("INVALID_FAMILY", "Family must have at least one parent", nil)

// Check if an error is a domain error
if errorswrapper.IsDomainError(err) {
    // Handle domain error
    code := errorswrapper.GetErrorCode(err)
    message := errorswrapper.GetErrorMessage(err)
    // ...
}

// Create a validation error
validationErr := errorswrapper.NewValidationError("First name cannot be empty", "firstName", nil)

// Create a database error
dbErr := errorswrapper.NewDatabaseError("Failed to insert record", "insert", "families", someError)

// Format an error for logging
formattedErr := errorswrapper.FormatError(err)
logger.Error(formattedErr)
```

## Configuration

The `errorswrapper` package doesn't require any specific configuration. It's a stateless wrapper around the `servicelib/errors` package.

## Testing

The package includes unit tests that verify the correct behavior of all error types and functions. Tests cover:

- Creating different types of errors
- Checking error types
- Getting error details
- Formatting errors
- Wrapping errors

## Design Notes

1. **Dependency Inversion**: The wrapper follows the Dependency Inversion Principle by ensuring that the domain layer depends on abstractions rather than concrete implementations.
2. **Error Types**: The package provides specific error types for different categories of errors, making it easier to handle errors appropriately.
3. **Error Details**: All errors include details such as error code, message, and cause, making it easier to debug issues.
4. **Consistency**: The wrapper ensures consistent error handling throughout the application.

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Go Error Handling Best Practices](https://blog.golang.org/error-handling-and-go)
- [Domain Errors](../../../core/domain/errors/README.md) - Domain-specific error types
- [Validation Wrapper](../validationwrapper/README.md) - Validation utilities that use these errors
