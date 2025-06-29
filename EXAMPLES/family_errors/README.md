# Family Service Error Handling Example

## Overview

This example demonstrates how to use both standard servicelib errors and domain-specific errors in the Family Service application. It shows how to create, handle, and check different types of errors that are commonly encountered when working with family data.

## Features

- **Standard Error Types**: Shows how to use database, not found, and validation errors
- **Domain-Specific Errors**: Demonstrates family structure, family status, and person-related errors
- **Error Handling**: Shows how to handle different types of errors
- **Error Type Checking**: Demonstrates how to check the type of an error

## Running the Example

To run this example, navigate to this directory and execute:

```bash
go run main.go
```

## Code Walkthrough

### Standard Error Types

The example shows how to create and use standard error types from the servicelib/errors package:

```
// Test NewDatabaseError for PostgreSQL operations
dbErr := errors.NewDatabaseError("failed to create families table", "create", "families", fmt.Errorf("some error"))
fmt.Printf("DatabaseError: %v\n", dbErr)

// Test NewNotFoundError for a missing family
notFoundErr := errors.NewNotFoundError("Family", "123", nil)
fmt.Printf("NotFoundError: %v\n", notFoundErr)

// Test NewValidationError for a required field
validationErr := errors.NewValidationError("id is required", "id", nil)
fmt.Printf("ValidationError: %v\n", validationErr)
```

### Domain-Specific Error Types

The example demonstrates how to create and use domain-specific error types from the core/domain/errors package:

```
// Family structure errors
tooManyParentsErr := domainerrors.NewFamilyTooManyParentsError("family cannot have more than two parents", nil)
fmt.Printf("FamilyTooManyParentsError: %v\n", tooManyParentsErr)

// Family status errors
notMarriedErr := domainerrors.NewFamilyNotMarriedError("only married families can divorce", nil)
fmt.Printf("FamilyNotMarriedError: %v\n", notMarriedErr)

// Person-related errors
parentDeceasedErr := domainerrors.NewParentAlreadyDeceasedError("parent is already marked as deceased", nil)
fmt.Printf("ParentAlreadyDeceasedError: %v\n", parentDeceasedErr)
```

### Error Handling

The example shows how to handle different types of errors:

```
// Example of handling domain-specific errors
err := domainerrors.NewFamilyNotMarriedError("only married families can divorce", nil)
fmt.Printf("Domain error example: %v\n", err)

// Example of handling standard errors
err = errors.NewValidationError("id is required", "id", nil)
fmt.Printf("Validation error example: %v\n", err)
```

### Error Type Checking

The example demonstrates how to check the type of an error:

```
// Create a domain-specific error
err = domainerrors.NewFamilyTooManyParentsError("family cannot have more than two parents", nil)

// Check if it's a domain error
if _, ok := err.(*errors.DomainError); ok {
    fmt.Println("Error is a domain error")
}

// Create a validation error
err = errors.NewValidationError("id is required", "id", nil)

// Check if it's a validation error
if _, ok := err.(*errors.ValidationError); ok {
    fmt.Println("Error is a validation error")
}
```

## Expected Output

```
=== Standard Error Types ===
DatabaseError: failed to create families table: operation=create, table=families: some error
QueryError: failed to get family from PostgreSQL: operation=query, table=families: some error
NotFoundError: Family with ID 123 not found
ValidationError: id is required: id

=== Domain-Specific Error Types ===
FamilyTooManyParentsError: family cannot have more than two parents
FamilyParentExistsError: parent already exists in family
FamilyChildExistsError: child already exists in family
FamilyNotMarriedError: only married families can divorce
FamilyDivorceRequiresTwoParentsError: divorce requires exactly two parents
ParentAlreadyDeceasedError: parent is already marked as deceased
ChildAlreadyDeceasedError: child is already marked as deceased

=== Error Handling Examples ===
Domain error example: only married families can divorce
Validation error example: id is required: id

=== Error Type Checking ===
Error is a domain error
Error is a validation error
```

## Related Examples

- [Basic Error Handling](../errors/README.md) - Shows basic error handling with servicelib/errors
- [Auth Directive Example](../auth_directive/README.md) - Shows how to handle authentication errors

## Related Components

- [Domain Errors Package](../../core/domain/errors/README.md) - Domain-specific error types
- [Error Wrapper Package](../../infrastructure/adapters/errorswrapper/README.md) - The error wrapper used in the application
- [Validation Package](../../core/domain/validation/README.md) - Validation utilities that use these errors

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.