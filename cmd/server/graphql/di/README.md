# Dependency Injection Container

## Overview

The Dependency Injection (DI) package provides a container for managing dependencies in the GraphQL server. It follows the principles of Clean Architecture and Hexagonal Architecture, ensuring proper separation of concerns and dependency inversion. The container initializes and manages all the dependencies required by the GraphQL server, including repositories, domain services, application services, and infrastructure components.

## Architecture

The DI package sits in the infrastructure layer of the application and is responsible for wiring together components from different layers. It follows these principles:

- **Dependency Inversion**: High-level modules don't depend on low-level modules, but both depend on abstractions
- **Clean Architecture**: The container respects the dependency rules of Clean Architecture
- **Hexagonal Architecture**: The container connects the application core to external adapters
- **Generic Programming**: Uses Go's generics for type-safe dependency management

The package is organized into:

- **Container**: The base container that manages common dependencies
- **FamilyContainer**: A specialized container for the family domain that extends the base container
- **Initializers**: Functions for initializing specific dependencies like repositories

## Implementation Details

The DI package implements the following design patterns:

1. **Dependency Injection Pattern**: Provides dependencies to components rather than having components create their own
2. **Factory Pattern**: Creates and configures dependencies
3. **Singleton Pattern**: Ensures single instances of dependencies
4. **Strategy Pattern**: Supports different database strategies (MongoDB, PostgreSQL, SQLite)
5. **Generic Programming**: Uses Go's generics for type-safe dependency management

Key implementation details:

- **Database Abstraction**: Supports multiple database types (MongoDB, PostgreSQL, SQLite)
- **Lifecycle Management**: Manages the lifecycle of dependencies, including initialization and cleanup
- **Configuration Integration**: Uses application configuration to configure dependencies
- **Error Handling**: Provides clear error messages for initialization failures
- **Resource Cleanup**: Ensures proper cleanup of resources when the container is closed

## Features

- **Multi-Database Support**: Supports MongoDB, PostgreSQL, and SQLite
- **Generic Container**: Provides a generic implementation that works with any repository
- **Dependency Lifecycle Management**: Manages the lifecycle of all dependencies
- **Configuration-Driven**: Uses application configuration to configure dependencies
- **Resource Cleanup**: Ensures proper cleanup of resources
- **Type Safety**: Leverages Go's generics for type-safe dependency management

## API Documentation

### Core Types

#### Container

The Container is the base dependency injection container that manages common dependencies.

```
// Container is a dependency injection container for the GraphQL server
type Container struct {
    *basedi.Container
    familyRepo          domainports.FamilyRepository
    familyDomainService *domainservices.FamilyDomainService
    familyAppService    appports.FamilyApplicationService
    familyMapper        dto.FamilyMapper
    authService         *auth.Auth
    dbType              string
    cache               *cache.Cache
}
```

#### FamilyContainer

The FamilyContainer is a specialized container for the family domain that extends the base container.

```
// FamilyContainer is a dependency injection container for the family domain
type FamilyContainer[T domainports.FamilyRepository] struct {
    *Container
    familyRepo          T
    familyDomainService *domainservices.FamilyDomainService
    familyAppService    ports.FamilyApplicationService
    cache               *cache.Cache
}
```

### Key Methods

#### NewContainer

Creates a new dependency injection container.

```
// NewContainer creates a new dependency injection container for the GraphQL server
func NewContainer(ctx context.Context, logger *zap.Logger, cfg *config.Config) (*Container, error)
```

#### NewFamilyContainer

Creates a new family dependency injection container.

```
// NewFamilyContainer creates a new family dependency injection container
func NewFamilyContainer[T domainports.FamilyRepository](
    ctx context.Context,
    logger *zap.Logger,
    cfg *config.Config,
    initRepo FamilyRepositoryInitializer[T],
    connectionString string,
) (*FamilyContainer[T], error)
```

#### Close

Closes all resources managed by the container.

```
// Close closes all resources
func (c *Container) Close() error
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [GraphQL Server Example](../../../EXAMPLES/graphql_server/README.md) - Shows how to use the DI container in a GraphQL server

Example of using the Container:

```
// Create a new container
container, err := di.NewContainer(ctx, logger, cfg)
if err != nil {
    // Handle error
}
defer container.Close()

// Get dependencies from the container
familyRepo := container.GetFamilyRepository()
familyService := container.GetFamilyApplicationService()
authService := container.GetAuthService()
```

## Configuration

The DI container is configured using the application configuration. The following configuration options are available:

- **Database Configuration**: Type, connection string, and other database-specific options
- **Cache Configuration**: TTL, size, and other cache-specific options
- **Authentication Configuration**: JWT secret key, token duration, and other auth-specific options

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

auth:
  jwt:
    secret_key: your-secret-key
    issuer: family-service
    token_duration: 24h
```

## Testing

The DI package is tested through:

1. **Unit Tests**: Each container method has unit tests
2. **Integration Tests**: Tests that verify the container works correctly with real dependencies
3. **Mock Tests**: Tests that use mock dependencies to isolate the container

Key testing approaches:

- **Mock Dependencies**: Tests use mock dependencies to isolate the container
- **Configuration Testing**: Tests verify that the container is correctly configured
- **Resource Cleanup**: Tests verify that resources are properly cleaned up
- **Error Handling**: Tests verify that errors are properly handled and propagated

## Design Notes

1. **Dependency Inversion**: The container follows the Dependency Inversion Principle by ensuring that high-level modules don't depend on low-level modules
2. **Clean Architecture**: The container respects the dependency rules of Clean Architecture
3. **Hexagonal Architecture**: The container connects the application core to external adapters
4. **Generic Programming**: The container uses Go's generics for type-safe dependency management
5. **Resource Lifecycle**: The container manages the lifecycle of all dependencies
6. **Error Handling**: The container provides clear error messages for initialization failures

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Dependency Injection](https://en.wikipedia.org/wiki/Dependency_injection)
- [Go Generics](https://go.dev/doc/tutorial/generics)
- [GraphQL Server](../../../cmd/server/graphql/README.md) - Uses this container to wire together the GraphQL server
- [Family Domain Services](../../../core/domain/services/README.md) - Domain services managed by this container
- [Family Application Services](../../../core/application/services/README.md) - Application services managed by this container
- [Repository Implementations](../../../infrastructure/adapters/repository/README.md) - Repository implementations managed by this container