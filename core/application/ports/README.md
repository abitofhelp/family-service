# Application Ports

## Overview

The Application Ports package defines interfaces that act as ports in the Hexagonal Architecture pattern. These interfaces are defined in the application layer but are implemented in the application layer and used by the interface layer. They provide a clear boundary between the application core and the outside world.

## Features

- **Generic Application Service Interface**: Provides a reusable interface for common CRUD operations
- **Family Application Service Interface**: Extends the generic interface with family-specific operations
- **Clean Architecture Compliance**: Follows the principles of Clean Architecture by defining clear boundaries
- **Type Safety**: Leverages Go's generics for type-safe interface definitions
- **Extensibility**: Easily extendable for new entity types and operations

## API Documentation

### Core Types

#### ApplicationService

A generic interface for application services that provides common CRUD operations.

```
// ApplicationService is a generic interface for application services
type ApplicationService[T any, D any] interface {
    // Create creates a new entity
    Create(ctx context.Context, dto D) (D, error)

    // GetByID retrieves an entity by ID
    GetByID(ctx context.Context, id string) (D, error)

    // GetAll retrieves all entities
    GetAll(ctx context.Context) ([]D, error)
}
```

#### FamilyApplicationService

An interface that extends the generic ApplicationService for family-specific operations.

```
// FamilyApplicationService defines the interface for family application services
type FamilyApplicationService interface {
    // Embed the generic ApplicationService interface with Family entity and DTO
    ApplicationService[*entity.Family, *entity.FamilyDTO]

    // Embed the servicelib ApplicationService interface
    di.ApplicationService

    // CreateFamily creates a new family (alias for Create)
    CreateFamily(ctx context.Context, dto entity.FamilyDTO) (*entity.FamilyDTO, error)

    // Additional family-specific methods...
}
```

## Best Practices

1. **Interface Segregation**: Keep interfaces focused on specific responsibilities
2. **Dependency Inversion**: Use these interfaces to invert dependencies between layers
3. **Consistent Naming**: Follow consistent naming conventions for interface methods
4. **Context Propagation**: Always include context.Context as the first parameter
5. **Error Handling**: Return meaningful errors that can be handled by the caller

## Related Components

- [Application Services](../services/README.md) - Implements these interfaces
- [Domain Entities](../../domain/entity/README.md) - Provides the entity types used by these interfaces
- [GraphQL Resolvers](../../../interface/adapters/graphql/resolver/README.md) - Uses these interfaces to handle GraphQL requests
