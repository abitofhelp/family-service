# Servicelib Refactoring Guide for Family Service

## Overview

This document outlines the changes needed to refactor the family-service repository to take full advantage of the features provided by the servicelib package. The refactoring is divided into several areas:

1. Error Handling
2. Transaction Handling
3. Database Connection Handling
4. Logging
5. Configuration

## 1. Error Handling

The current error handling in the family-service repository uses older error functions from servicelib. These should be updated to use the newer error handling functions.

### 1.1 Replace `errors.NewRepositoryError` with `errors.NewDatabaseError`

The `errors.NewRepositoryError` function should be replaced with the `errors.NewDatabaseError` function, which provides more specific error information for database-related errors.

**Current:**
```go
return errors.NewRepositoryError(err, "failed to create families table", "POSTGRES_ERROR")
```

**Recommended:**
```go
return errors.NewDatabaseError("failed to create families table", "CREATE", "families", err)
```

### 1.2 Update `errors.NewNotFoundError` to include a cause error

The `errors.NewNotFoundError` function should be updated to include a cause error parameter.

**Current:**
```go
return nil, errors.NewNotFoundError("Family", id)
```

**Recommended:**
```go
return nil, errors.NewNotFoundError("Family", id, err)
```

### 1.3 Update `errors.NewValidationError` to include field and cause error parameters

The `errors.NewValidationError` function should be updated to include field and cause error parameters.

**Current:**
```go
return nil, errors.NewValidationError("id is required")
```

**Recommended:**
```go
return nil, errors.NewValidationError("id is required", "id", nil)
```

## 2. Transaction Handling

The current transaction handling in the family-service repository uses manual transaction management. This should be updated to use servicelib's transaction handling functions, which provide retry mechanisms and better error handling.

### 2.1 Replace manual transaction handling with `db.ExecutePostgresTransaction`

**Current:**
```go
tx, err := r.DB.Begin(ctx)
if err != nil {
    r.logger.Error(ctx, "Failed to begin transaction", zap.Error(err), zap.String("family_id", fam.ID()))
    return errors.NewRepositoryError(err, "failed to begin transaction", "POSTGRES_ERROR")
}

var txErr error
defer func() {
    if txErr != nil {
        if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
            // Log rollback error, but don't return it as it would mask the original error
        }
    }
}()

// Execute SQL
_, txErr = tx.Exec(ctx, `...`)

if txErr != nil {
    return errors.NewRepositoryError(txErr, "failed to save family to PostgreSQL", "POSTGRES_ERROR")
}

// Commit transaction
if txErr = tx.Commit(ctx); txErr != nil {
    return errors.NewRepositoryError(txErr, "failed to commit transaction", "POSTGRES_ERROR")
}
```

**Recommended:**
```go
err = db.ExecutePostgresTransaction(ctx, r.DB, func(tx pgx.Tx) error {
    // Execute SQL
    _, err := tx.Exec(ctx, `...`)
    if err != nil {
        r.logger.Error(ctx, "Failed to save family to PostgreSQL", zap.Error(err), zap.String("family_id", fam.ID()))
        return errors.NewDatabaseError("failed to save family to PostgreSQL", "INSERT", "families", err)
    }
    return nil
}, db.RetryConfig{
    MaxRetries:     3,
    InitialBackoff: 100 * time.Millisecond,
    MaxBackoff:     2 * time.Second,
    BackoffFactor:  2.0,
    Logger:         r.logger,
})
```

### 2.2 Replace manual transaction handling with `db.ExecuteSQLTransaction` for SQLite

**Current:**
```go
tx, err := r.DB.Begin()
if err != nil {
    r.logger.Error(ctx, "Failed to begin transaction", zap.Error(err), zap.String("family_id", fam.ID()))
    return errors.NewRepositoryError(err, "failed to begin transaction", "SQLITE_ERROR")
}

var txErr error
defer func() {
    if txErr != nil {
        if rollbackErr := tx.Rollback(); rollbackErr != nil {
            // Log rollback error, but don't return it as it would mask the original error
        }
    }
}()

// Execute SQL
_, txErr = tx.Exec(`...`)

if txErr != nil {
    return errors.NewRepositoryError(txErr, "failed to save family to SQLite", "SQLITE_ERROR")
}

// Commit transaction
if txErr = tx.Commit(); txErr != nil {
    return errors.NewRepositoryError(txErr, "failed to commit transaction", "SQLITE_ERROR")
}
```

**Recommended:**
```go
err = db.ExecuteSQLTransaction(ctx, r.DB, func(tx *sql.Tx) error {
    // Execute SQL
    _, err := tx.Exec(`...`)
    if err != nil {
        r.logger.Error(ctx, "Failed to save family to SQLite", zap.Error(err), zap.String("family_id", fam.ID()))
        return errors.NewDatabaseError("failed to save family to SQLite", "INSERT", "families", err)
    }
    return nil
}, db.RetryConfig{
    MaxRetries:     3,
    InitialBackoff: 100 * time.Millisecond,
    MaxBackoff:     2 * time.Second,
    BackoffFactor:  2.0,
    Logger:         r.logger,
})
```

## 3. Database Connection Handling

The current database connection handling in the family-service repository has been updated to use servicelib's database connection handling functions. This is a good start, but there are still some improvements that can be made.

### 3.1 Use `db.IsTransientError` to check for transient errors

**Current:**
```go
if err != nil {
    // No check for transient errors
    return errors.NewRepositoryError(err, "failed to execute query", "POSTGRES_ERROR")
}
```

**Recommended:**
```go
if err != nil {
    if db.IsTransientError(err) {
        // Retry the operation
        return errors.NewDatabaseError("failed to execute query (transient error)", "SELECT", "families", err)
    }
    return errors.NewDatabaseError("failed to execute query", "SELECT", "families", err)
}
```

### 3.2 Use `db.CheckPostgresHealth` to check database health

**Current:**
```go
// No health check
```

**Recommended:**
```go
if err := db.CheckPostgresHealth(ctx, r.DB); err != nil {
    r.logger.Error(ctx, "Database health check failed", zap.Error(err))
    return nil, err
}
```

## 4. Logging

The current logging in the family-service repository uses servicelib's logging package. This is good, but there are still some improvements that can be made.

### 4.1 Use structured logging consistently

**Current:**
```go
r.logger.Error(ctx, "Failed to create families table in PostgreSQL", zap.Error(err))
```

**Recommended:**
```go
r.logger.Error(ctx, "Failed to create families table in PostgreSQL",
    zap.Error(err),
    zap.String("operation", "CREATE"),
    zap.String("table", "families"))
```

### 4.2 Use appropriate log levels

**Current:**
```go
r.logger.Debug(ctx, "Getting family by ID from PostgreSQL", zap.String("family_id", id))
```

**Recommended:**
```go
r.logger.Debug(ctx, "Getting family by ID from PostgreSQL",
    zap.String("family_id", id),
    zap.String("operation", "SELECT"),
    zap.String("table", "families"))
```

## 5. Configuration

The current configuration in the family-service repository uses servicelib's configuration package. This is good, but there are still some improvements that can be made.

### 5.1 Use `config.Validate` to validate configuration

**Current:**
```go
// No validation
```

**Recommended:**
```go
if err := config.Validate(cfg); err != nil {
    return nil, err
}
```

### 5.2 Use `config.LoadWithOptions` to load configuration with options

**Current:**
```go
cfg, err := config.LoadConfig()
```

**Recommended:**
```go
cfg, err := config.LoadWithOptions(config.Options{
    EnvPrefix:      "FAMILY_SERVICE",
    ConfigFileName: "config",
    ConfigFilePath: "./config",
})
```

## Implementation Strategy

The refactoring should be done in a systematic way to ensure consistency across the codebase. Here's a recommended strategy:

1. Update the imports to ensure the correct packages are being used
2. Replace all occurrences of `errors.NewRepositoryError` with `errors.NewDatabaseError`
3. Update all occurrences of `errors.NewNotFoundError` to include a cause error parameter
4. Update all occurrences of `errors.NewValidationError` to include field and cause error parameters
5. Replace manual transaction handling with servicelib's transaction handling functions
6. Use `db.IsTransientError` to check for transient errors
7. Use `db.CheckPostgresHealth` to check database health
8. Use structured logging consistently
9. Use appropriate log levels
10. Use `config.Validate` to validate configuration
11. Use `config.LoadWithOptions` to load configuration with options
12. Test the changes to ensure they work correctly

## Files to Update

The following files need to be updated:

1. `infrastructure/adapters/postgres/repo.go`
2. `infrastructure/adapters/mongo/repo.go`
3. `infrastructure/adapters/sqlite/repo.go`
4. `infrastructure/adapters/config/config.go`
5. `cmd/server/graphql/main.go`

## Conclusion

Refactoring the family-service repository to take full advantage of servicelib's features will improve the consistency, reliability, and maintainability of the codebase. This will make it easier to understand and debug errors, and will provide more specific error information to clients of the API.