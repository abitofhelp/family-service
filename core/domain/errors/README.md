# Domain Errors

## Overview

The Domain Errors package provides domain-specific error codes and error handling utilities for the Family Service application. It defines a set of error codes and functions for creating domain errors with specific error codes, making it easier to identify and handle domain-specific error conditions.

## Architecture

The Domain Errors package is part of the core domain layer in the Clean Architecture and Hexagonal Architecture patterns. It sits at the center of the application and has no dependencies on other layers. The architecture follows these principles:

- **Domain-Driven Design (DDD)**: Errors are designed based on the ubiquitous language of the domain
- **Clean Architecture**: Domain errors are independent of infrastructure concerns
- **Hexagonal Architecture**: Domain errors are used by both the domain and application layers

The package is organized into:

- **Error Codes**: Constants that define specific error conditions
- **Error Creation Functions**: Functions that create domain errors with specific error codes
- **Error Categories**: Logical groupings of related errors (family structure, family status, etc.)

## Implementation Details

The Domain Errors package implements the following design patterns:

1. **Factory Pattern**: Functions create domain errors with specific error codes
2. **Decorator Pattern**: Errors wrap underlying errors to maintain error context
3. **Categorization Pattern**: Errors are organized into logical categories

Key implementation details:

- **Error Codes**: Constants that define specific error conditions
- **Error Messages**: Clear and descriptive error messages
- **Error Wrapping**: Support for wrapping underlying errors
- **Error Categories**: Logical groupings of related errors
- **Error Creation Functions**: Functions that create domain errors with specific error codes

The package uses the `errorswrapper` package from the infrastructure layer to create domain errors. This wrapper provides a consistent error handling approach throughout the application while maintaining the dependency inversion principle.

## Features

- **Domain-Specific Error Codes**: Defines error codes for various domain-specific error conditions
- **Error Creation Functions**: Provides functions for creating domain errors with specific error codes
- **Error Categorization**: Organizes errors into categories (family structure, family status, parent-related, child-related)
- **Error Wrapping**: Supports wrapping underlying errors to maintain error context
- **Clean Error Messages**: Ensures error messages are clear and descriptive

## Examples

There may be additional examples in the /EXAMPLES directory.

Example of creating and handling a domain error:

```
// Create a domain error
err := errors.NewFamilyTooManyParentsError("family cannot have more than two parents", nil)

// Handle the error
if errors.IsFamilyTooManyParentsError(err) {
    // Handle the specific error
    fmt.Println("Too many parents error:", err)
} else if errors.IsFamilyError(err) {
    // Handle any family-related error
    fmt.Println("Family error:", err)
} else {
    // Handle other errors
    fmt.Println("Other error:", err)
}
```

## Configuration

The Domain Errors package doesn't require any specific configuration as it contains pure domain logic. However, it does have some configurable aspects:

- **Error Codes**: The error codes are defined as constants and can be modified if needed
- **Error Messages**: The error messages are provided when creating errors and can be customized
- **Error Wrapping**: The package supports wrapping underlying errors, which can be enabled or disabled

## Testing

The Domain Errors package is tested through:

1. **Unit Tests**: Each error creation function has unit tests
2. **Integration Tests**: Tests that verify error handling across layers
3. **Error Type Tests**: Tests that verify error type checking functions

Key testing approaches:

- **Error Creation Testing**: Tests that verify error creation functions create the correct error types
- **Error Type Testing**: Tests that verify error type checking functions correctly identify error types
- **Error Message Testing**: Tests that verify error messages are correct
- **Error Wrapping Testing**: Tests that verify error wrapping works correctly

Example of a test case:

```
func TestNewFamilyTooManyParentsError(t *testing.T) {
    // Create a domain error
    message := "family cannot have more than two parents"
    cause := fmt.Errorf("some underlying error")
    err := errors.NewFamilyTooManyParentsError(message, cause)

    // Verify the error
    assert.Error(t, err)
    assert.True(t, errors.IsFamilyTooManyParentsError(err))
    assert.True(t, errors.IsFamilyError(err))
    assert.Contains(t, err.Error(), message)
    assert.Contains(t, err.Error(), cause.Error())
}
```

## Design Notes

1. **Error Codes**: Error codes are defined as constants to ensure consistency
2. **Error Categories**: Errors are organized into logical categories to make them easier to understand and handle
3. **Error Wrapping**: Errors wrap underlying errors to maintain error context
4. **Error Messages**: Error messages are clear and descriptive to make debugging easier
5. **Error Type Checking**: Functions are provided to check error types, making error handling more robust
6. **Dependency Inversion**: The package uses the `errorswrapper` package from the infrastructure layer to create errors, maintaining the dependency inversion principle

## API Documentation

### Error Codes

The package defines several domain-specific error codes:

```
// Family structure errors
const (
    FamilyTooManyParentsCode     = "FAMILY_TOO_MANY_PARENTS"
    FamilyParentExistsCode       = "FAMILY_PARENT_EXISTS"
    FamilyParentDuplicateCode    = "FAMILY_PARENT_DUPLICATE"
    FamilyChildExistsCode        = "FAMILY_CHILD_EXISTS"
    FamilyCannotRemoveLastParent = "FAMILY_CANNOT_REMOVE_LAST_PARENT"
)

// Family status errors
const (
    FamilyNotMarriedCode         = "FAMILY_NOT_MARRIED"
    FamilyDivorceRequiresTwoCode = "FAMILY_DIVORCE_REQUIRES_TWO_PARENTS"
    FamilyCreateFailedCode       = "FAMILY_CREATE_FAILED"
    FamilyStatusUpdateFailedCode = "FAMILY_STATUS_UPDATE_FAILED"
)

// Parent-related errors
const (
    ParentAlreadyDeceasedCode = "PARENT_ALREADY_DECEASED"
)

// Child-related errors
const (
    ChildAlreadyDeceasedCode = "CHILD_ALREADY_DECEASED"
)
```

### Key Functions

#### NewFamilyTooManyParentsError

Creates a new domain error for when a family has too many parents.

```
// NewFamilyTooManyParentsError creates a new domain error for when a family has too many parents
func NewFamilyTooManyParentsError(message string, cause error) error
```

#### NewFamilyNotMarriedError

Creates a new domain error for when a family is not in a married state.

```
// NewFamilyNotMarriedError creates a new domain error for when a family is not in a married state
func NewFamilyNotMarriedError(message string, cause error) error
```

#### NewFamilyDivorceRequiresTwoParentsError

Creates a new domain error for when a divorce is attempted without two parents.

```
// NewFamilyDivorceRequiresTwoParentsError creates a new domain error for when a divorce is attempted without two parents
func NewFamilyDivorceRequiresTwoParentsError(message string, cause error) error
```

## Best Practices

1. **Descriptive Error Codes**: Use descriptive error codes that clearly indicate the error condition
2. **Consistent Naming**: Follow a consistent naming convention for error codes and error creation functions
3. **Error Wrapping**: Always wrap underlying errors to maintain error context
4. **Clear Error Messages**: Provide clear and descriptive error messages
5. **Error Categorization**: Organize errors into logical categories

## References

- [Domain-Driven Design](https://domainlanguage.com/ddd/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Error Handling in Go](https://blog.golang.org/error-handling-and-go)
- [Domain Entities](../entity/README.md) - Uses these errors to indicate domain rule violations
- [Domain Services](../services/README.md) - Uses these errors to indicate domain operation failures
- [Application Services](../../application/services/README.md) - Handles these errors at the application layer
- [Error Wrapper](../../../infrastructure/adapters/errorswrapper/README.md) - Provides the underlying error wrapping functionality
