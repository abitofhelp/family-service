# Domain Errors

## Overview

The Domain Errors package provides domain-specific error codes and error handling utilities for the Family Service application. It defines a set of error codes and functions for creating domain errors with specific error codes, making it easier to identify and handle domain-specific error conditions.

## Features

- **Domain-Specific Error Codes**: Defines error codes for various domain-specific error conditions
- **Error Creation Functions**: Provides functions for creating domain errors with specific error codes
- **Error Categorization**: Organizes errors into categories (family structure, family status, parent-related, child-related)
- **Error Wrapping**: Supports wrapping underlying errors to maintain error context
- **Clean Error Messages**: Ensures error messages are clear and descriptive

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

## Related Components

- [Domain Entities](../entity/README.md) - Uses these errors to indicate domain rule violations
- [Domain Services](../services/README.md) - Uses these errors to indicate domain operation failures
- [Application Services](../../application/services/README.md) - Handles these errors at the application layer
- [Error Wrapper](../../../infrastructure/adapters/errorswrapper/README.md) - Provides the underlying error wrapping functionality