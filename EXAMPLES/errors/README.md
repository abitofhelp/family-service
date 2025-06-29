# Error Handling Example

## Overview

This example demonstrates how to use the error types provided by the ServiceLib errors package. It shows how to create and handle different types of errors that are commonly used in the Family Service application.

## Features

- **DatabaseError**: Shows how to create and use database-related errors
- **NotFoundError**: Shows how to create and use resource not found errors
- **ValidationError**: Shows how to create and use validation errors

## Running the Example

To run this example, navigate to this directory and execute:

```bash
go run main.go
```

## Code Walkthrough

### Creating a DatabaseError

The example shows how to create a database error with a message, operation, table name, and cause:

```
// Test NewDatabaseError
dbErr := errors.NewDatabaseError("failed to create table", "create", "families", fmt.Errorf("some error"))
fmt.Printf("DatabaseError: %v\n", dbErr)
```

### Creating a NotFoundError

The example shows how to create a not found error with a resource type, resource ID, and cause:

```
// Test NewNotFoundError
notFoundErr := errors.NewNotFoundError("Family", "123", fmt.Errorf("some error"))
fmt.Printf("NotFoundError: %v\n", notFoundErr)
```

### Creating a ValidationError

The example shows how to create a validation error with a message, field name, and cause:

```
// Test NewValidationError
validationErr := errors.NewValidationError("invalid input", "field", fmt.Errorf("some error"))
fmt.Printf("ValidationError: %v\n", validationErr)
```

## Expected Output

```
DatabaseError: failed to create table: operation=create, table=families: some error
NotFoundError: Family with ID 123 not found: some error
ValidationError: invalid input: field: some error
```

## Related Examples

- [Auth Directive Example](../auth_directive/README.md) - Shows how to handle authentication errors
- [Family Errors Example](../family_errors/README.md) - Shows how to handle domain-specific errors

## Related Components

- [Error Wrapper Package](../../infrastructure/adapters/errorswrapper/README.md) - The error wrapper used in the application
- [Domain Errors Package](../../core/domain/errors/README.md) - Domain-specific error types
- [Validation Package](../../core/domain/validation/README.md) - Validation utilities that use these errors

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.
