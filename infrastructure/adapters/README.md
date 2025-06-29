# Infrastructure Adapters

## Overview

The Infrastructure Adapters package provides implementations of the ports defined in the domain layer, following the Ports and Adapters (Hexagonal) Architecture pattern. These adapters connect the application core to external systems and libraries, ensuring that the domain layer remains isolated from infrastructure concerns. The package includes adapters for databases, caching, logging, validation, and other infrastructure services.

## Architecture

The Infrastructure Adapters package sits in the infrastructure layer of the application and implements the ports defined in the domain layer. It follows these principles:

- **Hexagonal Architecture (Ports and Adapters)**: Implements the adapter side of the ports and adapters pattern
- **Dependency Inversion**: The domain layer defines what it needs, and the adapters implement these interfaces
- **Clean Architecture**: The adapters depend on the domain layer, not the other way around
- **Separation of Concerns**: Each adapter focuses on a specific infrastructure concern

The package is organized into:

- **Database Adapters**: Implementations for different database types (MongoDB, PostgreSQL, SQLite)
- **Wrapper Adapters**: Wrappers around external libraries to maintain dependency inversion
- **Service Adapters**: Implementations of infrastructure services like caching, logging, and validation
- **Security Adapters**: Implementations of authentication and authorization services

## Implementation Details

The Infrastructure Adapters package implements the following design patterns:

1. **Adapter Pattern**: Adapts external libraries and systems to the interfaces defined in the domain layer
2. **Repository Pattern**: Provides implementations of repository interfaces for different database types
3. **Decorator Pattern**: Adds cross-cutting concerns like logging, caching, and validation to repositories
4. **Factory Pattern**: Creates and configures adapter instances
5. **Strategy Pattern**: Supports different strategies for the same interface (e.g., different database types)

Key implementation details:

- **Database Abstraction**: Supports multiple database types (MongoDB, PostgreSQL, SQLite)
- **External Library Wrappers**: Wraps external libraries to maintain dependency inversion
- **Error Handling**: Translates infrastructure errors to domain errors
- **Configuration Integration**: Uses application configuration to configure adapters
- **Resource Management**: Manages resources like database connections and caches

## Features

- **Multi-Database Support**: Implementations for MongoDB, PostgreSQL, and SQLite
- **External Library Wrappers**: Wrappers for external libraries to maintain dependency inversion
- **Infrastructure Services**: Implementations of caching, logging, validation, and other services
- **Error Translation**: Translates infrastructure errors to domain errors
- **Configuration Integration**: Uses application configuration to configure adapters
- **Resource Management**: Manages resources like database connections and caches

## Components

The Infrastructure Adapters package includes the following components:

- **cachewrapper**: Wrapper around the cache library
- **circuitwrapper**: Wrapper around the circuit breaker library
- **config**: Configuration loading and management
- **datewrapper**: Wrapper around date handling libraries
- **diwrapper**: Wrapper around dependency injection libraries
- **errors**: Infrastructure-specific error types
- **errorswrapper**: Wrapper around error handling libraries
- **identificationwrapper**: Wrapper around ID generation libraries
- **loggingwrapper**: Wrapper around logging libraries
- **mongo**: MongoDB repository implementation
- **postgres**: PostgreSQL repository implementation
- **profiling**: Profiling utilities
- **ratewrapper**: Wrapper around rate limiting libraries
- **repository**: Base repository implementations
- **repositorywrapper**: Wrapper around repository libraries
- **security**: Authentication and authorization services
- **sqlite**: SQLite repository implementation
- **telemetrywrapper**: Wrapper around telemetry libraries
- **validationwrapper**: Wrapper around validation libraries

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Repository Example](../../EXAMPLES/repository/README.md) - Shows how to use the repository adapters
- [Authentication Example](../../EXAMPLES/auth_directive/README.md) - Shows how to use the authentication adapters

Example of using the repository adapters:

```
// Create a repository implementation
repo := postgres.NewFamilyRepository(db)

// Use the repository
family, err := repo.GetByID(ctx, "family-123")
if err != nil {
    // Handle error
}

// Save an entity
err = repo.Save(ctx, family)
if err != nil {
    // Handle error
}
```

## Configuration

The Infrastructure Adapters are configured using the application configuration. The following configuration options are available:

- **Database Configuration**: Type, connection string, and other database-specific options
- **Cache Configuration**: TTL, size, and other cache-specific options
- **Logging Configuration**: Level, format, and other logging-specific options
- **Security Configuration**: JWT secret key, token duration, and other security-specific options

Example configuration:

```yaml
database:
  type: sqlite  # or mongodb, postgres
  sqlite:
    uri: file:data/dev/sqlite/family_service.db?cache=shared&mode=rwc
  mongodb:
    uri: mongodb://localhost:27017
  postgres:
    dsn: postgres://user:pass@localhost:5432/familydb

cache:
  ttl: 5m
  size: 1000

logging:
  level: info
  format: json

auth:
  jwt:
    secret_key: your-secret-key
    issuer: family-service
    token_duration: 24h
```

## Testing

The Infrastructure Adapters package is tested through:

1. **Unit Tests**: Each adapter has unit tests
2. **Integration Tests**: Tests that verify the adapters work correctly with real external systems
3. **Mock Tests**: Tests that use mock external systems to isolate the adapters

Key testing approaches:

- **Mock External Systems**: Tests use mock external systems to isolate the adapters
- **Real External Systems**: Integration tests use real external systems to verify the adapters work correctly
- **Error Handling**: Tests verify that errors are properly handled and translated
- **Resource Management**: Tests verify that resources are properly managed

## Design Notes

1. **Dependency Inversion**: The adapters implement interfaces defined in the domain layer, ensuring that the domain layer doesn't depend on infrastructure concerns
2. **Clean Architecture**: The adapters depend on the domain layer, not the other way around
3. **Hexagonal Architecture**: The adapters connect the application core to external systems and libraries
4. **Separation of Concerns**: Each adapter focuses on a specific infrastructure concern
5. **Error Translation**: The adapters translate infrastructure errors to domain errors
6. **Resource Management**: The adapters manage resources like database connections and caches

## References

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Adapter Pattern](https://en.wikipedia.org/wiki/Adapter_pattern)
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)
- [Domain Ports](../../core/domain/ports/README.md) - The interfaces implemented by these adapters
- [Domain Entities](../../core/domain/entity/README.md) - The entities managed by these adapters
- [Application Services](../../core/application/services/README.md) - Services that use these adapters
- [GraphQL Resolvers](../../interface/adapters/graphql/resolver/README.md) - Resolvers that use these adapters