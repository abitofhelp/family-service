# Application Ports

## Overview

The Application Ports package defines interfaces that act as ports in the Hexagonal Architecture pattern. These interfaces are defined in the application layer but are implemented in the application layer and used by the interface layer. They provide a clear boundary between the application core and the outside world.

## Architecture

The Application Ports package follows the Ports and Adapters (Hexagonal) Architecture pattern. In this pattern:

- **Ports**: Interfaces defined in the application layer that specify how the application interacts with the outside world
- **Adapters**: Implementations of these interfaces that connect the application to external systems

This package contains the port definitions (interfaces) that are implemented by adapters in other layers. The ports are divided into:

- **Primary Ports (Driving)**: Used by the interface layer to drive the application
- **Secondary Ports (Driven)**: Used by the application to interact with infrastructure

The Application Ports sit at the boundary of the application layer, providing a clear separation between the application core and external concerns.

## Implementation Details

The Application Ports package implements the following design patterns:

1. **Interface Segregation Principle**: Each interface is focused on a specific responsibility
2. **Dependency Inversion Principle**: High-level modules (application services) depend on abstractions (interfaces)
3. **Generic Programming**: Uses Go's generics to create reusable interface definitions

Key implementation details:

- **Generic Interfaces**: The `ApplicationService` interface uses type parameters to create a reusable interface for any entity type
- **Interface Composition**: The `FamilyApplicationService` interface embeds the generic interface to inherit its methods
- **Context Propagation**: All methods accept a context.Context parameter for cancellation and value propagation
- **Error Handling**: Methods return errors that can be handled by the caller

## Examples

There may be additional examples in the /EXAMPLES directory.

Example usage of the FamilyApplicationService interface:

```
// In your application code:

// Create a new family
family, err := familyService.CreateFamily(ctx, familyDTO)
if err != nil {
    // Handle error
}

// Get a family by ID
family, err = familyService.GetByID(ctx, "family-123")
if err != nil {
    // Handle error
}

// Get all families
families, err := familyService.GetAll(ctx)
if err != nil {
    // Handle error
}
```

## Configuration

The Application Ports package doesn't require any specific configuration as it only defines interfaces. However, implementations of these interfaces may require configuration for:

- Database connections
- Authentication settings
- Caching parameters
- Logging levels

These configurations are typically provided to the implementations through dependency injection.

## Testing

The Application Ports package is tested through:

1. **Mock Implementations**: Using GoMock to create mock implementations of the interfaces
2. **Integration Tests**: Testing the interfaces with real implementations
3. **Unit Tests**: Testing the interface definitions for correctness

Example of creating a mock implementation:

```
// In your test code:

// Create a mock implementation of the FamilyApplicationService interface
ctrl := gomock.NewController(t)
defer ctrl.Finish()

mockService := mock.NewMockFamilyApplicationService(ctrl)
mockService.EXPECT().GetByID(gomock.Any(), "family-123").Return(expectedFamily, nil)

// Use the mock implementation in tests
family, err := mockService.GetByID(ctx, "family-123")
assert.NoError(t, err)
assert.Equal(t, expectedFamily, family)
```

## Design Notes

1. **Generic Programming**: The use of generics allows for type-safe, reusable interface definitions
2. **Interface Composition**: Embedding interfaces allows for inheritance of methods while maintaining the ability to add specific methods
3. **Context Propagation**: All methods accept a context.Context parameter for cancellation and value propagation
4. **Error Handling**: Methods return errors that can be handled by the caller
5. **Naming Conventions**: Interface names follow the pattern of `EntityNameService` for clarity

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

## References

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Go Generics](https://go.dev/doc/tutorial/generics)
- [Domain-Driven Design](https://domainlanguage.com/ddd/)
- [Application Services](../services/README.md) - Implements these interfaces
- [Domain Entities](../../domain/entity/README.md) - Provides the entity types used by these interfaces
- [GraphQL Resolvers](../../../interface/adapters/graphql/resolver/README.md) - Uses these interfaces to handle GraphQL requests
