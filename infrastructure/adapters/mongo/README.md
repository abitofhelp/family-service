# Infrastructure Adapters - MongoDB

## Overview

The MongoDB adapter provides implementations for database-related ports defined in the core domain and application layers, specifically for MongoDB database interactions. This adapter connects the application to MongoDB database systems, following the Ports and Adapters (Hexagonal) architecture pattern. By isolating MongoDB implementations in adapter classes, the core business logic remains independent of specific database technologies, making the system more maintainable, testable, and flexible.

## Features

- MongoDB connection management
- Document querying and manipulation
- Transaction support
- Data mapping between domain entities and MongoDB documents
- Error handling and translation
- Connection pooling
- Performance optimization
- Support for MongoDB-specific features (aggregation, geospatial queries, etc.)
- Index management

## Installation

```bash
go get github.com/abitofhelp/family-service/infrastructure/adapters/mongo
```

## Configuration

The MongoDB adapter can be configured according to specific requirements. Here's an example of configuring the MongoDB adapter:

```
// Pseudocode example - not actual Go code
// This demonstrates how to configure and use a MongoDB adapter

// 1. Import necessary packages
import mongo, config, logging

// 2. Create a logger
logger = logging.NewLogger()

// 3. Configure the MongoDB connection
mongoConfig = {
    connectionString: "mongodb://localhost:27017",
    database: "family_service",
    username: "user",
    password: "password",
    authSource: "admin",
    maxPoolSize: 100,
    minPoolSize: 10,
    maxConnIdleTime: 30 seconds,
    connectTimeout: 10 seconds,
    serverSelectionTimeout: 5 seconds,
    retryWrites: true,
    writeConcern: "majority"
}

// 4. Create the MongoDB adapter
mongoAdapter = mongo.NewMongoAdapter(mongoConfig, logger)

// 5. Use the MongoDB adapter
err = mongoAdapter.Connect()
if err != nil {
    logger.Error("Failed to connect to MongoDB", err)
}

// Find a document
family, err = mongoAdapter.FindOne(context, "families", {id: "family-123"})
if err != nil {
    logger.Error("Failed to find family", err)
}
```

## API Documentation

### Core Concepts

The MongoDB adapter follows these core concepts:

1. **Adapter Pattern**: Implements database ports defined in the core domain or application layer
2. **Dependency Injection**: Receives dependencies through constructor injection
3. **Configuration**: Configured through a central configuration system
4. **Logging**: Uses a consistent logging approach
5. **Error Handling**: Translates MongoDB-specific errors to domain errors

### Key Adapter Functions

```
// Pseudocode example - not actual Go code
// This demonstrates a MongoDB adapter implementation

// MongoDB adapter structure
type MongoAdapter {
    config        // MongoDB configuration
    logger        // Logger for logging operations
    contextLogger // Context-aware logger
    client        // MongoDB client
    database      // MongoDB database
}

// Constructor for the MongoDB adapter
function NewMongoAdapter(config, logger) {
    return new MongoAdapter {
        config: config,
        logger: logger,
        contextLogger: new ContextLogger(logger),
        client: nil,
        database: nil
    }
}

// Method to connect to the database
function MongoAdapter.Connect() {
    // Implementation would include:
    // 1. Logging the operation
    // 2. Creating a MongoDB client
    // 3. Connecting to the MongoDB server
    // 4. Getting a database reference
    // 5. Handling connection errors
}

// Method to find a single document
function MongoAdapter.FindOne(context, collection, filter) {
    // Implementation would include:
    // 1. Logging the operation with context
    // 2. Getting the collection
    // 3. Executing the find operation
    // 4. Mapping the result to a domain entity
    // 5. Handling query errors
    // 6. Returning the result or error
}
```

## Best Practices

1. **Separation of Concerns**: Keep database logic separate from domain logic
2. **Interface Segregation**: Define focused database interfaces in the domain layer
3. **Dependency Injection**: Use constructor injection for adapter dependencies
4. **Error Translation**: Translate MongoDB-specific errors to domain errors
5. **Consistent Logging**: Use a consistent logging approach
6. **Transaction Management**: Implement proper transaction handling where supported
7. **Connection Management**: Properly manage database connections
8. **Testing**: Write unit and integration tests for MongoDB adapters
9. **Indexing**: Create appropriate indexes for query performance

## Troubleshooting

### Common Issues

#### Database Connection Issues

If you encounter database connection issues, check the following:
- MongoDB server is running and accessible
- Connection string parameters are correct
- Network connectivity between the application and the database
- Proper authentication credentials are provided
- Firewall rules allow the connection

#### Performance Issues

If you encounter performance issues with MongoDB, consider the following:
- Create appropriate indexes for frequently queried fields
- Use projection to limit returned fields
- Optimize query patterns
- Configure appropriate connection pool settings
- Use query analysis tools to identify slow queries
- Consider database-specific optimizations (like read preferences)

## Related Components

- [Domain Layer](../../core/domain/README.md) - The domain layer that defines the database ports
- [Application Layer](../../core/application/README.md) - The application layer that uses database operations
- [Repository Adapters](../repository/README.md) - The repository adapters that use MongoDB

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.