# Domain Ports

## Overview

The Domain Ports package defines interfaces that act as ports in the Hexagonal Architecture pattern. These interfaces are defined in the domain layer but implemented in the infrastructure layer, providing a clean separation between the domain logic and the infrastructure concerns. This package follows the Dependency Inversion Principle, allowing the domain layer to define the interfaces it needs without depending on specific implementations.

## Architecture

The Domain Ports package is part of the core domain layer in the Clean Architecture and Hexagonal Architecture patterns. It sits at the center of the application and has no dependencies on other layers. The architecture follows these principles:

- **Hexagonal Architecture (Ports and Adapters)**: Defines ports (interfaces) that are implemented by adapters in the infrastructure layer
- **Dependency Inversion Principle**: The domain layer defines what it needs from the infrastructure layer
- **Clean Architecture**: The domain layer is independent of infrastructure concerns
- **Domain-Driven Design**: Interfaces are designed based on the ubiquitous language of the domain

The package is organized into:

- **Repository Interfaces**: Interfaces for data persistence operations
- **Service Interfaces**: Interfaces for domain services
- **Generic Interfaces**: Reusable interface definitions using Go's generics
- **Mock Subdirectory**: Contains mock implementations for testing

## Implementation Details

The Domain Ports package implements the following design patterns:

1. **Interface Segregation Principle**: Each interface is focused on a specific responsibility
2. **Dependency Inversion Principle**: High-level modules (domain) depend on abstractions (interfaces)
3. **Generic Programming**: Uses Go's generics to create reusable interface definitions
4. **Repository Pattern**: Defines interfaces for data access operations
5. **Factory Pattern**: Defines interfaces for creating domain objects

Key implementation details:

- **Generic Interfaces**: The `Repository` interface uses type parameters to create a reusable interface for any entity type
- **Context Propagation**: All methods accept a context.Context parameter for cancellation and value propagation
- **Error Handling**: Methods return errors that can be handled by the caller
- **Domain-Specific Methods**: Interfaces include domain-specific methods beyond basic CRUD operations
- **Mock Generation**: Interfaces are designed to work with GoMock for generating mock implementations

## Features

- **Repository Interfaces**: Defines interfaces for data persistence operations
- **Hexagonal Architecture**: Implements the ports side of the ports and adapters pattern
- **Dependency Inversion**: Allows the domain layer to define what it needs from infrastructure
- **Generic Interfaces**: Leverages Go's generics for type-safe interface definitions
- **Clean Separation**: Maintains a clear boundary between domain and infrastructure

## API Documentation

### Core Types

#### FamilyRepository

The FamilyRepository interface defines the contract for family persistence operations. It embeds a generic Repository interface and adds family-specific methods.

```
// FamilyRepository defines the interface for family persistence operations
// This interface represents a port in the Hexagonal Architecture pattern
// It's defined in the domain layer but implemented in the infrastructure layer
type FamilyRepository interface {
    // Embed the generic Repository interface with Family entity
    repositorywrapper.Repository[*entity.Family]

    // FindByParentID finds families that contain a specific parent
    FindByParentID(ctx context.Context, parentID string) ([]*entity.Family, error)

    // FindByChildID finds the family that contains a specific child
    FindByChildID(ctx context.Context, childID string) (*entity.Family, error)
}
```

### Mock Implementations

The package includes mock implementations of the interfaces for testing purposes. These mocks are generated using GoMock and can be found in the `mock` subdirectory.

## Examples

There may be additional examples in the /EXAMPLES directory.

Example of using the FamilyRepository interface:

```
// Create a repository implementation
repo := postgres.NewFamilyRepository(db)

// Get a family by ID
family, err := repo.GetByID(ctx, "family-123")
if err != nil {
    // Handle error
}

// Find families by parent ID
families, err := repo.FindByParentID(ctx, "parent-456")
if err != nil {
    // Handle error
}

// Find family by child ID
family, err = repo.FindByChildID(ctx, "child-789")
if err != nil {
    // Handle error
}

// Save a family
err = repo.Save(ctx, family)
if err != nil {
    // Handle error
}
```

## Configuration

The Domain Ports package doesn't require any specific configuration as it only defines interfaces. However, implementations of these interfaces may require configuration for:

- Database connections
- Connection pooling
- Retry policies
- Timeout settings
- Caching strategies

These configurations are typically provided to the implementations through dependency injection.

## Testing

The Domain Ports package is tested through:

1. **Mock Implementations**: Using GoMock to create mock implementations of the interfaces
2. **Integration Tests**: Testing the interfaces with real implementations
3. **Contract Tests**: Ensuring that all implementations adhere to the interface contract

Key testing approaches:

- **Mock Generation**: Interfaces are designed to work with GoMock for generating mock implementations
- **Behavior Verification**: Tests verify that the correct methods are called with the correct parameters
- **Error Handling**: Tests verify that errors are properly propagated
- **Context Propagation**: Tests verify that context is properly propagated

Example of creating a mock implementation:

```
// Create a mock controller
ctrl := gomock.NewController(t)
defer ctrl.Finish()

// Create a mock repository
mockRepo := mock.NewMockFamilyRepository(ctrl)

// Set up expectations
mockRepo.EXPECT().
    GetByID(gomock.Any(), "family-123").
    Return(&entity.Family{}, nil)

// Use the mock repository
family, err := mockRepo.GetByID(ctx, "family-123")
assert.NoError(t, err)
assert.NotNil(t, family)
```

## Design Notes

1. **Interface Segregation**: Interfaces are focused on specific responsibilities to avoid bloated interfaces
2. **Dependency Inversion**: Interfaces are defined in the domain layer but implemented in the infrastructure layer
3. **Generic Programming**: Go's generics are used to create reusable interface definitions
4. **Context Propagation**: All methods accept a context.Context parameter for cancellation and value propagation
5. **Error Handling**: Methods return errors that can be handled by the caller
6. **Domain-Specific Methods**: Interfaces include domain-specific methods beyond basic CRUD operations
7. **Repository Pattern**: The Repository pattern is used to abstract data access operations

## Best Practices

1. **Interface Segregation**: Keep interfaces focused on specific responsibilities
2. **Dependency Inversion**: Define interfaces in the domain layer, implement them in the infrastructure layer
3. **Consistent Naming**: Follow consistent naming conventions for interface methods
4. **Context Propagation**: Always include context.Context as the first parameter
5. **Error Handling**: Return meaningful errors that can be handled by the caller

## References

- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Dependency Inversion Principle](https://en.wikipedia.org/wiki/Dependency_inversion_principle)
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)
- [Go Generics](https://go.dev/doc/tutorial/generics)
- [Domain Entities](../entity/README.md) - The entities managed by these repositories
- [Domain Services](../services/README.md) - Services that use these repositories
- [Infrastructure Adapters](../../../infrastructure/adapters/repository/README.md) - Implementations of these interfaces
- [Mock Implementations](./mock/README.md) - Mock implementations for testing
