# Error Handling Guide for Family Service

## Overview

This document describes the error handling approach used in the Family Service, including the standard error types and domain-specific error codes.

## Error Types

The Family Service uses a combination of standard error types from the `servicelib/errors` package and domain-specific error codes defined in the `core/domain/errors` package.

### Standard Error Types

The following standard error types are used throughout the application:

1. **ValidationError**: Used for validation failures, such as invalid input data.
2. **NotFoundError**: Used when a requested resource cannot be found.
3. **DatabaseError**: Used for database-related errors.
4. **DomainError**: Used for domain-specific business rule violations.

### Domain-Specific Error Codes

To provide more granular error categorization, the Family Service defines domain-specific error codes that are used with the `DomainError` type:

#### Family Structure Errors

- `FAMILY_TOO_MANY_PARENTS`: A family cannot have more than two parents.
- `FAMILY_PARENT_EXISTS`: A parent already exists in the family.
- `FAMILY_PARENT_DUPLICATE`: A parent with the same name and birthdate already exists in the family.
- `FAMILY_CHILD_EXISTS`: A child already exists in the family.
- `FAMILY_CANNOT_REMOVE_LAST_PARENT`: Cannot remove the only parent from a family.

#### Family Status Errors

- `FAMILY_NOT_MARRIED`: Only married families can divorce.
- `FAMILY_DIVORCE_REQUIRES_TWO_PARENTS`: Divorce requires exactly two parents.
- `FAMILY_CREATE_FAILED`: Failed to create a new family.
- `FAMILY_STATUS_UPDATE_FAILED`: Failed to update family status.

#### Person-Related Errors

- `PARENT_ALREADY_DECEASED`: A parent is already marked as deceased.
- `CHILD_ALREADY_DECEASED`: A child is already marked as deceased.

## Error Handling

### Creating Errors

The Family Service provides helper functions for creating domain-specific errors:

```go
// Example: Creating a domain-specific error
err := domainerrors.NewFamilyTooManyParentsError("family cannot have more than two parents", nil)

// Example: Creating a standard error
err := errors.NewValidationError("id is required", "id", nil)
```

### Handling Errors

When handling errors, you can use type assertions to determine the error type:

```go
// Check if it's a domain error
if domainErr, ok := err.(*errors.DomainError); ok {
    // Handle domain error
    fmt.Printf("Domain error: %v\n", domainErr)
}

// Check if it's a validation error
if validationErr, ok := err.(*errors.ValidationError); ok {
    // Handle validation error
    fmt.Printf("Validation error: %v\n", validationErr)
}
```

## Best Practices

1. **Use Specific Error Types**: Always use the most specific error type available for the situation.
2. **Include Context**: Include relevant context in error messages to help with debugging.
3. **Propagate Errors**: Propagate errors up the call stack, adding context as needed.
4. **Log Errors**: Log errors at appropriate levels (debug, info, warn, error) based on severity.
5. **Return Appropriate Status Codes**: Map errors to appropriate HTTP or GraphQL status codes in API responses.

## Examples

For examples of error handling, see the `examples/family_errors/main.go` file, which demonstrates creating and handling different error types.