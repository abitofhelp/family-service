# Repository Wrapper

## Overview

The Repository Wrapper package provides a wrapper around the servicelib/repository package to ensure the domain layer doesn't directly depend on external libraries. This follows the Dependency Inversion Principle, where high-level modules should not depend on low-level modules, but both should depend on abstractions.

## Architecture

The Repository Wrapper package follows the Adapter pattern from Hexagonal Architecture, providing a layer of abstraction over the external `servicelib/repository` package. This ensures that the core domain doesn't directly depend on external libraries, maintaining the dependency inversion principle.

The package sits in the infrastructure layer of the application and is used by the domain layer through interfaces defined in the domain layer. The architecture follows these principles:

- **Dependency Inversion**: The domain layer depends on abstractions, not concrete implementations
- **Adapter Pattern**: This package adapts the external library to the domain's needs
- **Repository Pattern**: Provides a consistent interface for data access operations
- **Generic Programming**: Uses Go's generics for type-safe repository operations

## Implementation Details

The Repository Wrapper package implements the following design patterns:

1. **Adapter Pattern**: Adapts the external library to the domain's needs
2. **Repository Pattern**: Provides a consistent interface for data access operations
3. **Generic Programming**: Uses Go's generics for type-safe repository operations
4. **Facade Pattern**: Simplifies the interface to the underlying repository implementation

Key implementation details:

- **Generic Interface**: The `Repository` interface uses type parameters to create a reusable interface for any entity type
- **Delegation**: The `RepositoryWrapper` delegates to the underlying repository implementation
- **Context Propagation**: All methods accept a context.Context parameter for cancellation and value propagation
- **Error Handling**: Methods return errors that can be handled by the caller
- **Type Safety**: Leverages Go's generics for type-safe repository operations

## Features

- **Dependency Inversion**: Decouples the domain layer from external libraries
- **Generic Interface**: Provides a generic Repository interface that can be used with any entity type
- **Adapter Pattern**: Implements the adapter pattern to wrap the servicelib/repository functionality
- **Type Safety**: Leverages Go's generics for type-safe repository operations

## API Documentation

### Core Types

#### Repository

The Repository interface defines the contract for basic repository operations. It's a generic interface that can be used with any entity type.

```
// Repository is a generic interface for repository operations
// It wraps the servicelib/repository.Repository interface
type Repository[T any] interface {
    // GetByID retrieves an entity by its ID
    GetByID(ctx context.Context, id string) (T, error)

    // Save persists an entity
    Save(ctx context.Context, entity T) error

    // GetAll retrieves all entities
    GetAll(ctx context.Context) ([]T, error)
}
```

#### RepositoryWrapper

The RepositoryWrapper struct implements the Repository interface by delegating to the servicelib/repository.Repository interface.

```
// RepositoryWrapper is a wrapper around servicelib/repository.Repository
type RepositoryWrapper[T any] struct {
    repo repository.Repository[T]
}
```

### Key Methods

#### NewRepositoryWrapper

Creates a new RepositoryWrapper instance.

```
// NewRepositoryWrapper creates a new RepositoryWrapper
func NewRepositoryWrapper[T any](repo repository.Repository[T]) *RepositoryWrapper[T]
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Repository Example](../../../examples/repository/README.md) - Shows how to use the repository wrapper

Example of using the repository wrapper:

```
// Create a repository implementation
repo := postgres.NewFamilyRepository(db)

// Wrap it with the repository wrapper
wrappedRepo := repositorywrapper.NewRepositoryWrapper(repo)

// Use the wrapped repository
family, err := wrappedRepo.GetByID(ctx, "family-123")
if err != nil {
    // Handle error
}

// Save an entity
err = wrappedRepo.Save(ctx, family)
if err != nil {
    // Handle error
}

// Get all entities
families, err := wrappedRepo.GetAll(ctx)
if err != nil {
    // Handle error
}
```

## Configuration

The Repository Wrapper package doesn't require any specific configuration. It's a stateless wrapper around the underlying repository implementation. However, the underlying repository implementation may require configuration for:

- Database connections
- Connection pooling
- Retry policies
- Timeout settings
- Caching strategies

These configurations are typically provided to the underlying repository implementation through dependency injection.

## Testing

The Repository Wrapper package is tested through:

1. **Unit Tests**: Each method has unit tests that verify it correctly delegates to the underlying repository
2. **Integration Tests**: Tests that verify the wrapper works correctly with real repository implementations
3. **Mock Tests**: Tests that use mock repositories to verify the wrapper's behavior

Key testing approaches:

- **Mock Repositories**: Tests use mock repositories to verify that the wrapper correctly delegates to the underlying repository
- **Error Propagation**: Tests verify that errors from the underlying repository are properly propagated
- **Context Propagation**: Tests verify that context is properly propagated to the underlying repository

Example of a test case:

```
func TestRepositoryWrapper_GetByID(t *testing.T) {
    // Create a mock repository
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()
    mockRepo := mock.NewMockRepository[*entity.Family](ctrl)

    // Set up expectations
    expectedFamily := &entity.Family{}
    mockRepo.EXPECT().
        GetByID(gomock.Any(), "family-123").
        Return(expectedFamily, nil)

    // Create the wrapper
    wrapper := repositorywrapper.NewRepositoryWrapper(mockRepo)

    // Call the method
    family, err := wrapper.GetByID(context.Background(), "family-123")

    // Verify the result
    assert.NoError(t, err)
    assert.Equal(t, expectedFamily, family)
}
```

## Best Practices

1. **Dependency Inversion**: Use this wrapper to avoid direct dependencies on external libraries in the domain layer
2. **Interface Segregation**: Keep interfaces focused on specific responsibilities
3. **Consistent Naming**: Follow consistent naming conventions for interface methods
4. **Context Propagation**: Always include context.Context as the first parameter
5. **Error Handling**: Return meaningful errors that can be handled by the caller

## Design Notes

1. **Adapter Pattern**: The Repository Wrapper implements the Adapter pattern to provide a layer of abstraction over the external repository library
2. **Delegation**: The wrapper delegates all operations to the underlying repository implementation
3. **Type Safety**: The wrapper leverages Go's generics to provide type-safe repository operations
4. **Minimal Interface**: The wrapper exposes only the methods needed by the domain layer
5. **Context Propagation**: All methods accept a context.Context parameter for cancellation and value propagation
6. **Error Handling**: All methods return errors that can be handled by the caller
7. **Dependency Inversion**: The wrapper follows the Dependency Inversion Principle by ensuring that the domain layer depends on abstractions rather than concrete implementations

## Related Components

- [Domain Ports](../../../core/domain/ports/README.md) - The domain layer interfaces that use this wrapper
- [Infrastructure Adapters](../../adapters/README.md) - Other infrastructure adapters

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Adapter Pattern](https://en.wikipedia.org/wiki/Adapter_pattern)
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)
- [Go Generics](https://go.dev/doc/tutorial/generics)
- [Domain Entities](../../../core/domain/entity/README.md) - The entities managed by repositories
- [Domain Ports](../../../core/domain/ports/README.md) - The interfaces that define repository contracts
- [Application Services](../../../core/application/services/README.md) - Services that use repositories
