# Error Handling Refactoring Guide

## Overview

This document outlines the changes needed to refactor the error handling in the family-service repository to take full advantage of the error handling capabilities provided by the servicelib package.

## Current Error Handling

The current error handling in the family-service repository uses the following error functions:

1. `errors.NewRepositoryError(err, message, code)` - Used for repository-related errors
2. `errors.NewNotFoundError(resourceType, id)` - Used for resource not found errors
3. `errors.NewValidationError(message)` - Used for validation errors

These functions are from an older version of the servicelib package and should be updated to use the newer error handling functions.

## Recommended Changes

### 1. Replace `errors.NewRepositoryError` with `errors.NewDatabaseError`

The `errors.NewRepositoryError` function should be replaced with the `errors.NewDatabaseError` function, which provides more specific error information for database-related errors.

**Current:**
```go
return errors.NewRepositoryError(err, "failed to create families table", "POSTGRES_ERROR")
```

**Recommended:**
```go
return errors.NewDatabaseError("failed to create families table", "CREATE", "families", err)
```

### 2. Update `errors.NewNotFoundError` to include a cause error

The `errors.NewNotFoundError` function should be updated to include a cause error parameter.

**Current:**
```go
return nil, errors.NewNotFoundError("Family", id)
```

**Recommended:**
```go
return nil, errors.NewNotFoundError("Family", id, err)
```

### 3. Update `errors.NewValidationError` to include field and cause error parameters

The `errors.NewValidationError` function should be updated to include field and cause error parameters.

**Current:**
```go
return nil, errors.NewValidationError("id is required")
```

**Recommended:**
```go
return nil, errors.NewValidationError("id is required", "id", nil)
```

## Implementation Strategy

The error handling refactoring should be done in a systematic way to ensure consistency across the codebase. Here's a recommended strategy:

1. Update the imports to ensure the correct error package is being used
2. Replace all occurrences of `errors.NewRepositoryError` with `errors.NewDatabaseError`
3. Update all occurrences of `errors.NewNotFoundError` to include a cause error parameter
4. Update all occurrences of `errors.NewValidationError` to include field and cause error parameters
5. Test the changes to ensure they work correctly

## Files to Update

The following files need to be updated:

1. `infrastructure/adapters/postgres/repo.go`
2. `infrastructure/adapters/mongo/repo.go`
3. `infrastructure/adapters/sqlite/repo.go`

## Example Replacements

### PostgresFamilyRepository

#### Database Errors

```go
// Before
return errors.NewRepositoryError(err, "failed to create families table", "POSTGRES_ERROR")

// After
return errors.NewDatabaseError("failed to create families table", "CREATE", "families", err)
```

```go
// Before
return nil, errors.NewRepositoryError(err, "failed to get family from PostgreSQL", "POSTGRES_ERROR")

// After
return nil, errors.NewDatabaseError("failed to get family from PostgreSQL", "SELECT", "families", err)
```

```go
// Before
return nil, errors.NewRepositoryError(err, "failed to unmarshal parents data", "JSON_ERROR")

// After
return nil, errors.NewDatabaseError("failed to unmarshal parents data", "UNMARSHAL", "families", err)
```

#### Not Found Errors

```go
// Before
return nil, errors.NewNotFoundError("Family", id)

// After
return nil, errors.NewNotFoundError("Family", id, nil)
```

```go
// Before
return nil, errors.NewNotFoundError("Family with Child", childID)

// After
return nil, errors.NewNotFoundError("Family with Child", childID, nil)
```

#### Validation Errors

```go
// Before
return nil, errors.NewValidationError("id is required")

// After
return nil, errors.NewValidationError("id is required", "id", nil)
```

```go
// Before
return errors.NewValidationError("family cannot be nil")

// After
return errors.NewValidationError("family cannot be nil", "family", nil)
```

## Conclusion

Refactoring the error handling in the family-service repository to take full advantage of the servicelib package's error handling capabilities will improve the consistency and clarity of error handling throughout the codebase. This will make it easier to understand and debug errors, and will provide more specific error information to clients of the API.