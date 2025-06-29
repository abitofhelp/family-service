# Infrastructure Adapters

## Overview

The Infrastructure Adapters package provides implementations of the ports defined in the core domain and application layers. These adapters connect the application to external systems, frameworks, and libraries, following the Ports and Adapters (Hexagonal) architecture pattern. By isolating these implementations in adapter classes, the core business logic remains independent of infrastructure concerns, making the system more maintainable, testable, and flexible.

## Features

- **Cache Adapters**: Wrappers for caching mechanisms to improve performance
- **Circuit Breaker Adapters**: Implementations for circuit breaking patterns to handle failures gracefully
- **Configuration Adapters**: Adapters for loading and accessing application configuration
- **Date/Time Adapters**: Wrappers for date and time operations
- **Dependency Injection Adapters**: Adapters for managing dependencies and their lifecycle
- **Error Handling Adapters**: Adapters for consistent error handling across the application
- **Identification Adapters**: Adapters for generating and validating identifiers
- **Logging Adapters**: Wrappers for logging frameworks
- **Database Adapters**: Implementations for various database systems (MongoDB, PostgreSQL, SQLite)
- **Profiling Adapters**: Adapters for application profiling and performance monitoring
- **Rate Limiting Adapters**: Adapters for rate limiting to protect resources
- **Repository Adapters**: Implementations of the repository interfaces defined in the domain layer
- **Security Adapters**: Adapters for authentication, authorization, and other security concerns
- **Telemetry Adapters**: Adapters for collecting and reporting telemetry data
- **Validation Adapters**: Adapters for input validation

## Installation

```bash
go get github.com/abitofhelp/family-service/infrastructure/adapters
```

## Configuration

Each adapter can be configured according to its specific requirements. Here's an example of configuring the cache adapter:

```
// Pseudocode example - not actual Go code
// This demonstrates how to configure and use a cache adapter

// 1. Import necessary packages
import cache, config, logging, time

// 2. Create a logger
logger = logging.NewLogger()

// 3. Configure the cache
cacheConfig = {
    enabled: true,
    ttl: 5 minutes,
    maxSize: 1000,
    purgeInterval: 1 minute
}

// 4. Create the cache adapter
cacheAdapter = cache.NewCache(cacheConfig, logger)

// 5. Use the cache adapter
cacheAdapter.Set("key", "value")
value, found = cacheAdapter.Get("key")
if found {
    logger.Info("Found value in cache", value)
}
```

## API Documentation

### Core Concepts

The infrastructure adapters follow these core concepts:

1. **Adapter Pattern**: Each adapter implements a port (interface) defined in the core domain or application layer
2. **Dependency Injection**: Adapters receive their dependencies through constructor injection
3. **Configuration**: Adapters are configured through a central configuration system
4. **Logging**: Adapters use a consistent logging approach
5. **Error Handling**: Adapters translate infrastructure-specific errors to domain errors

### Key Adapter Categories

#### Database Adapters

Database adapters provide implementations for repository interfaces defined in the domain layer. They handle the persistence of domain entities in various database systems.

```
// Pseudocode example - not actual Go code
// This demonstrates a repository adapter implementation

// Repository adapter structure
type FamilyRepository {
    database      // Database connection
    logger        // Logger for logging operations
    contextLogger // Context-aware logger
}

// Constructor for the repository adapter
function NewFamilyRepository(database, logger) {
    return new FamilyRepository {
        database: database,
        logger: logger,
        contextLogger: new ContextLogger(logger)
    }
}

// Method to find a family by ID
function FamilyRepository.FindById(context, id) {
    // Implementation would include:
    // 1. Logging the operation
    // 2. Querying the database
    // 3. Mapping database results to domain entities
    // 4. Handling errors and returning appropriate domain errors
    // 5. Returning the family entity or error
}
```

#### Wrapper Adapters

Wrapper adapters provide enhanced functionality around existing libraries or frameworks, adding features like logging, error handling, and metrics.

```
// Pseudocode example - not actual Go code
// This demonstrates a wrapper adapter implementation

// Wrapper adapter structure
type LoggingWrapper {
    logger // Logger instance
}

// Constructor for the wrapper adapter
function NewLoggingWrapper(logger) {
    return new LoggingWrapper {
        logger: logger
    }
}

// Method to create a context-aware logger
function LoggingWrapper.WithContext(context) {
    // Create a new context logger with the provided context
    return new ContextLogger(this.logger).WithContext(context)
}
```

## Best Practices

1. **Separation of Concerns**: Keep adapter implementations separate from domain logic
2. **Interface Segregation**: Define focused interfaces in the domain layer that adapters implement
3. **Dependency Injection**: Use constructor injection for adapter dependencies
4. **Error Translation**: Translate infrastructure-specific errors to domain errors
5. **Consistent Logging**: Use a consistent logging approach across all adapters
6. **Configuration**: Configure adapters through a central configuration system
7. **Testing**: Write unit and integration tests for adapters

## Troubleshooting

### Common Issues

#### Database Connection Issues

If you encounter database connection issues, check the following:
- Database connection string is correct
- Database server is running
- Network connectivity between the application and the database
- Proper authentication credentials are provided

#### Performance Issues

If you encounter performance issues, consider the following:
- Use caching for frequently accessed data
- Implement proper database indexing
- Use connection pooling for database connections
- Implement circuit breakers for external service calls

## Related Components

- [Domain Layer](../../core/domain/README.md) - The domain layer that defines the ports implemented by these adapters
- [Application Layer](../../core/application/README.md) - The application layer that uses these adapters
- [Interface Adapters](../../interface/adapters/README.md) - The interface adapters that use these infrastructure adapters

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.