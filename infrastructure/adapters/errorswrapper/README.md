# Error Wrapper

This package provides a wrapper around the `github.com/abitofhelp/servicelib/errors` package to ensure that the domain layer doesn't directly depend on external libraries. This follows the principles of Clean Architecture and Hexagonal Architecture (Ports and Adapters).

## Purpose

The purpose of this wrapper is to:

1. Isolate the domain layer from external dependencies
2. Provide a consistent error handling approach throughout the application
3. Make it easier to replace or update the underlying error handling library in the future

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

## Error Types

The wrapper provides the following error types:

- `Error`: Base error interface
- `ValidationError`: For validation errors
- `DomainError`: For domain-specific errors
- `DatabaseError`: For database-related errors
- `NotFoundError`: For resource not found errors

## Functions

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