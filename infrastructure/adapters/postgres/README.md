# Infrastructure Adapters - PostgreSQL

## Overview

The PostgreSQL adapter provides implementations for database-related ports defined in the core domain and application layers, specifically for PostgreSQL database interactions. This adapter connects the application to PostgreSQL database systems, following the Ports and Adapters (Hexagonal) architecture pattern. By isolating PostgreSQL implementations in adapter classes, the core business logic remains independent of specific database technologies, making the system more maintainable, testable, and flexible.

## Features

- PostgreSQL database connection management
- Query execution and result mapping
- Transaction management
- Migration support
- Data mapping between domain entities and PostgreSQL tables
- Error handling and translation
- Connection pooling
- Performance optimization
- Support for PostgreSQL-specific features (JSON, arrays, etc.)

## Installation

```bash
go get github.com/abitofhelp/family-service/infrastructure/adapters/postgres
```

## Configuration

The PostgreSQL adapter can be configured according to specific requirements. Here's an example of configuring the PostgreSQL adapter:

```
// Pseudocode example - not actual Go code
// This demonstrates how to configure and use a PostgreSQL adapter

// 1. Import necessary packages
import postgres, config, logging

// 2. Create a logger
logger = logging.NewLogger()

// 3. Configure the PostgreSQL connection
postgresConfig = {
    host: "localhost",
    port: 5432,
    username: "postgres",
    password: "password",
    database: "family_service",
    sslMode: "disable",
    maxOpenConnections: 25,
    maxIdleConnections: 10,
    connectionMaxLifetime: 1 hour,
    statementTimeout: 30 seconds
}

// 4. Create the PostgreSQL adapter
postgresAdapter = postgres.NewPostgresAdapter(postgresConfig, logger)

// 5. Use the PostgreSQL adapter
err = postgresAdapter.Connect()
if err != nil {
    logger.Error("Failed to connect to PostgreSQL database", err)
}

// Execute a query
results, err = postgresAdapter.Query(context, "SELECT * FROM families WHERE id = $1", "family-123")
if err != nil {
    logger.Error("Failed to execute query", err)
}
```

## API Documentation

### Core Concepts

The PostgreSQL adapter follows these core concepts:

1. **Adapter Pattern**: Implements database ports defined in the core domain or application layer
2. **Dependency Injection**: Receives dependencies through constructor injection
3. **Configuration**: Configured through a central configuration system
4. **Logging**: Uses a consistent logging approach
5. **Error Handling**: Translates PostgreSQL-specific errors to domain errors

### Key Adapter Functions

```
// Pseudocode example - not actual Go code
// This demonstrates a PostgreSQL adapter implementation

// PostgreSQL adapter structure
type PostgresAdapter {
    config        // PostgreSQL configuration
    logger        // Logger for logging operations
    contextLogger // Context-aware logger
    db            // Database connection
}

// Constructor for the PostgreSQL adapter
function NewPostgresAdapter(config, logger) {
    return new PostgresAdapter {
        config: config,
        logger: logger,
        contextLogger: new ContextLogger(logger),
        db: nil
    }
}

// Method to connect to the database
function PostgresAdapter.Connect() {
    // Implementation would include:
    // 1. Logging the operation
    // 2. Building the connection string
    // 3. Opening a connection to the PostgreSQL database
    // 4. Configuring connection parameters
    // 5. Setting up connection pooling
    // 6. Handling connection errors
}

// Method to execute a query
function PostgresAdapter.Query(context, query, args...) {
    // Implementation would include:
    // 1. Logging the operation with context
    // 2. Preparing the SQL statement
    // 3. Executing the query with arguments
    // 4. Mapping results to appropriate structures
    // 5. Handling query errors
    // 6. Returning results or error
}
```

## Best Practices

1. **Separation of Concerns**: Keep database logic separate from domain logic
2. **Interface Segregation**: Define focused database interfaces in the domain layer
3. **Dependency Injection**: Use constructor injection for adapter dependencies
4. **Error Translation**: Translate PostgreSQL-specific errors to domain errors
5. **Consistent Logging**: Use a consistent logging approach
6. **Transaction Management**: Implement proper transaction handling
7. **Connection Management**: Properly manage database connections
8. **Testing**: Write unit and integration tests for PostgreSQL adapters
9. **Migrations**: Use a migration strategy for database schema changes

## Troubleshooting

### Common Issues

#### Database Connection Issues

If you encounter database connection issues, check the following:
- Database server is running and accessible
- Connection string parameters are correct
- Network connectivity between the application and the database
- Proper authentication credentials are provided
- Firewall rules allow the connection

#### Performance Issues

If you encounter performance issues with PostgreSQL, consider the following:
- Implement proper indexing for frequently queried columns
- Use prepared statements for repeated queries
- Optimize query patterns
- Configure appropriate connection pool settings
- Use query analysis tools to identify slow queries
- Consider database-specific optimizations (like partitioning)

## Related Components

- [Domain Layer](../../core/domain/README.md) - The domain layer that defines the database ports
- [Application Layer](../../core/application/README.md) - The application layer that uses database operations
- [Repository Adapters](../repository/README.md) - The repository adapters that use PostgreSQL

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.