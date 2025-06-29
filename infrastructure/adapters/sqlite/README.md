# Infrastructure Adapters - SQLite

## Overview

The SQLite adapter provides implementations for database-related ports defined in the core domain and application layers, specifically for SQLite database interactions. This adapter connects the application to SQLite database systems, following the Ports and Adapters (Hexagonal) architecture pattern.

> **For Junior Developers**: Think of this adapter as a bridge between your business logic and the SQLite database. It allows your core business code to store and retrieve data without knowing the details of how SQLite works.

By isolating SQLite implementations in adapter classes, the core business logic remains independent of specific database technologies, making the system more maintainable, testable, and flexible.

## Features

- SQLite database connection management
- Query execution and result mapping
- Transaction management
- Migration support
- Data mapping between domain entities and SQLite tables
- Error handling and translation
- Connection pooling
- Performance optimization

## Getting Started

If you're new to this codebase, follow these steps to start using the SQLite adapter:

1. **Understand the purpose**: The SQLite adapter handles all database operations for SQLite
2. **Learn the interfaces**: Look at the domain repository interfaces to understand what operations are available
3. **Database location**: SQLite databases are stored as files, usually in the `./data` directory
4. **Ask questions**: If something isn't clear, ask a more experienced developer

## Installation

```bash
go get github.com/abitofhelp/family-service/infrastructure/adapters/sqlite
```

## Configuration

The SQLite adapter can be configured according to specific requirements. Here's an example of configuring the SQLite adapter:

```
// Pseudocode example - not actual Go code
// This demonstrates how to configure and use a SQLite adapter

// 1. Import necessary packages
import sqlite, config, logging

// 2. Create a logger
// This is needed for the SQLite adapter to log any issues
logger = logging.NewLogger()

// 3. Configure the SQLite connection
// These settings determine how the SQLite database behaves
sqliteConfig = {
    databasePath: "./data/family_service.db",      // Where the database file is stored
    journalMode: "WAL",                            // Write-Ahead Logging for better performance
    busyTimeout: 5000,                             // How long to wait if the database is locked (ms)
    foreignKeys: true,                             // Enable foreign key constraints
    maxOpenConnections: 10,                        // Maximum number of open connections
    maxIdleConnections: 5,                         // Maximum number of idle connections
    connectionMaxLifetime: 1 hour                  // How long a connection can be reused
}

// 4. Create the SQLite adapter
// This is the object you'll use for all database operations
sqliteAdapter = sqlite.NewSQLiteAdapter(sqliteConfig, logger)

// 5. Use the SQLite adapter
// First, connect to the database
err = sqliteAdapter.Connect()
if err != nil {
    logger.Error("Failed to connect to SQLite database", err)
}

// Execute a query with a parameter
// The "?" is a placeholder that will be replaced with "family-123"
results, err = sqliteAdapter.Query(context, "SELECT * FROM families WHERE id = ?", "family-123")
if err != nil {
    logger.Error("Failed to execute query", err)
}
```

## API Documentation

### Core Concepts

> **For Junior Developers**: These concepts are fundamental to understanding how the SQLite adapter works. Take time to understand each one before diving into the code.

The SQLite adapter follows these core concepts:

1. **Adapter Pattern**: Implements database ports defined in the core domain or application layer
   - This means the SQLite adapter implements interfaces defined elsewhere
   - The business logic only knows about these interfaces, not the SQLite implementation details

2. **Dependency Injection**: Receives dependencies through constructor injection
   - Dependencies like loggers are passed in when creating the adapter
   - This makes testing easier and components more loosely coupled

3. **Configuration**: Configured through a central configuration system
   - Settings like connection parameters are defined in configuration
   - This allows changing behavior without changing code

4. **Logging**: Uses a consistent logging approach
   - All database operations are logged for debugging and monitoring
   - Context information is included in logs when available

5. **Error Handling**: Translates SQLite-specific errors to domain errors
   - SQLite errors are converted to domain-specific errors
   - This prevents SQLite-specific error details from leaking into the domain

### Key Adapter Functions

Here are the main functions you'll use when working with the SQLite adapter:

```
// Pseudocode example - not actual Go code
// This demonstrates a SQLite adapter implementation

// SQLite adapter structure
type SQLiteAdapter {
    config        // SQLite configuration
    logger        // Logger for logging operations
    contextLogger // Context-aware logger
    db            // Database connection
}

// Constructor for the SQLite adapter
// This is how you create a new instance of the adapter
function NewSQLiteAdapter(config, logger) {
    return new SQLiteAdapter {
        config: config,
        logger: logger,
        contextLogger: new ContextLogger(logger),
        db: nil  // Database connection is initialized later
    }
}

// Method to connect to the database
// Call this before using any other methods
function SQLiteAdapter.Connect() {
    // Implementation would include:
    // 1. Logging the operation
    // 2. Opening a connection to the SQLite database
    // 3. Configuring connection parameters
    // 4. Setting up connection pooling
    // 5. Handling connection errors
}

// Method to execute a query
// Use this for SELECT statements
function SQLiteAdapter.Query(context, query, args...) {
    // Implementation would include:
    // 1. Logging the operation with context
    // 2. Preparing the SQL statement
    // 3. Executing the query with arguments
    // 4. Mapping results to appropriate structures
    // 5. Handling query errors
    // 6. Returning results or error
}
```

### Common Database Operations

Here are some common operations you might need to perform:

1. **Connecting to the database**:
   ```
   err = sqliteAdapter.Connect()
   if err != nil {
       // Handle connection error
   }
   ```

2. **Executing a query**:
   ```
   results, err = sqliteAdapter.Query(context, "SELECT * FROM users WHERE age > ?", 18)
   if err != nil {
       // Handle query error
   }
   ```

3. **Executing a command (insert, update, delete)**:
   ```
   rowsAffected, err = sqliteAdapter.Execute(context, "UPDATE users SET active = ? WHERE id = ?", true, "user-123")
   if err != nil {
       // Handle execution error
   }
   ```

4. **Using transactions**:
   ```
   tx, err = sqliteAdapter.BeginTransaction(context)
   if err != nil {
       // Handle transaction error
   }

   // Perform operations within the transaction

   err = tx.Commit()
   if err != nil {
       // Handle commit error
   }
   ```

## Best Practices

> **For Junior Developers**: Following these best practices will help you avoid common pitfalls and write more maintainable code.

1. **Separation of Concerns**: Keep database logic separate from domain logic
   - **Why?** Your business logic shouldn't need to know how data is stored or retrieved
   - **Example:** Don't put SQL queries in your domain entities or services

2. **Interface Segregation**: Define focused database interfaces in the domain layer
   - **Why?** Small, specific interfaces are easier to understand and implement
   - **Example:** Have separate repository interfaces for different entity types

3. **Dependency Injection**: Use constructor injection for adapter dependencies
   - **Why?** This makes testing easier and components more loosely coupled
   - **Example:** Pass the logger and configuration to the SQLite adapter constructor

4. **Error Translation**: Translate SQLite-specific errors to domain errors
   - **Why?** Domain code shouldn't need to understand SQLite error codes
   - **Example:** Convert "SQLITE_CONSTRAINT" errors to domain-specific validation errors

5. **Consistent Logging**: Use a consistent logging approach
   - **Why?** Makes it easier to debug issues and monitor performance
   - **Example:** Log all database operations with context information

6. **Transaction Management**: Implement proper transaction handling
   - **Why?** Ensures data consistency and integrity
   - **Example:** Use transactions for operations that modify multiple records

7. **Connection Management**: Properly manage database connections
   - **Why?** Prevents resource leaks and improves performance
   - **Example:** Use connection pooling and close connections when done

8. **Testing**: Write unit and integration tests for SQLite adapters
   - **Why?** Ensures the adapter works correctly and catches regressions
   - **Example:** Use an in-memory SQLite database for testing

## Common Mistakes to Avoid

1. **Not using parameterized queries**
   - **Problem:** Makes your code vulnerable to SQL injection attacks
   - **Solution:** Always use query parameters (?) instead of string concatenation

2. **Forgetting to close resources**
   - **Problem:** Can lead to resource leaks and performance issues
   - **Solution:** Always close statements, result sets, and connections

3. **Not handling concurrent access**
   - **Problem:** SQLite has limitations with concurrent writes
   - **Solution:** Use appropriate locking strategies and journal modes

4. **Using SQLite for high-concurrency workloads**
   - **Problem:** SQLite may not perform well with many concurrent writers
   - **Solution:** Consider other databases for high-concurrency scenarios

## Troubleshooting

### Common Issues

#### Database Connection Issues

If you encounter database connection issues, check the following:

- **Database file path is correct and accessible**
  - **Example:** Verify the path with `os.Stat(databasePath)`
  - **Solution:** Use absolute paths to avoid confusion

- **File permissions allow reading and writing to the database file**
  - **Problem:** The process might not have permission to access the file
  - **Solution:** Check and fix file permissions with `chmod`

- **Database is not locked by another process**
  - **Problem:** SQLite allows only one writer at a time
  - **Solution:** Increase busy timeout or use WAL journal mode

- **Connection string parameters are correct**
  - **Example:** Check for typos in parameter names
  - **Solution:** Use constants for parameter names to avoid typos

#### Performance Issues

If you encounter performance issues with SQLite, consider the following:

- **Use WAL (Write-Ahead Logging) journal mode**
  - **Why?** Allows concurrent reads while writing
  - **Example:** Set `PRAGMA journal_mode=WAL`

- **Implement proper indexing for frequently queried columns**
  - **Why?** Indexes speed up queries dramatically
  - **Example:** `CREATE INDEX idx_users_email ON users(email)`

- **Use prepared statements for repeated queries**
  - **Why?** Prepared statements are compiled once and executed many times
  - **Example:** Prepare a statement once and reuse it with different parameters

- **Optimize query patterns**
  - **Why?** Some query patterns are more efficient than others
  - **Example:** Use `EXISTS` instead of `COUNT(*)` when checking existence

- **Configure appropriate connection pool settings**
  - **Why?** Too few connections limit concurrency, too many waste resources
  - **Example:** Start with `maxOpenConnections = 10` and adjust based on load

## Related Components

> **For Junior Developers**: Understanding how components relate to each other is crucial for working effectively in this codebase.

- [Domain Layer](../../core/domain/README.md) - The domain layer that defines the database ports
  - This is where the interfaces that the SQLite adapter implements are defined
  - Look here to understand what operations are available

- [Application Layer](../../core/application/README.md) - The application layer that uses database operations
  - This layer contains the business logic that uses the SQLite adapter
  - See how database operations are used in business processes

- [Repository Adapters](../repository/README.md) - The repository adapters that use SQLite
  - These adapters implement the repository pattern using SQLite
  - They provide a higher-level abstraction over the SQLite adapter

## Glossary of Terms

- **Adapter Pattern**: A design pattern that allows incompatible interfaces to work together
- **Port**: An interface defined in the domain or application layer
- **Dependency Injection**: A technique where an object receives its dependencies from outside
- **SQLite**: A self-contained, serverless, zero-configuration, transactional SQL database engine
- **WAL**: Write-Ahead Logging, a journal mode that allows concurrent reads while writing
- **Transaction**: A sequence of operations performed as a single logical unit of work
- **Connection Pool**: A cache of database connections maintained for reuse
- **Prepared Statement**: A precompiled SQL statement that can be executed multiple times with different parameters

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.
