# Domain Ports

## Overview

The Domain Ports package defines interfaces that act as ports in the Hexagonal Architecture pattern. These interfaces are defined in the domain layer but implemented in the infrastructure layer, providing a clean separation between the domain logic and the infrastructure concerns. This package follows the Dependency Inversion Principle, allowing the domain layer to define the interfaces it needs without depending on specific implementations.

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
    repository.Repository[*entity.Family]

    // FindByParentID finds families that contain a specific parent
    FindByParentID(ctx context.Context, parentID string) ([]*entity.Family, error)

    // FindByChildID finds the family that contains a specific child
    FindByChildID(ctx context.Context, childID string) (*entity.Family, error)
}
```

### Mock Implementations

The package includes mock implementations of the interfaces for testing purposes. These mocks are generated using GoMock and can be found in the `mock` subdirectory.

## Best Practices

1. **Interface Segregation**: Keep interfaces focused on specific responsibilities
2. **Dependency Inversion**: Define interfaces in the domain layer, implement them in the infrastructure layer
3. **Consistent Naming**: Follow consistent naming conventions for interface methods
4. **Context Propagation**: Always include context.Context as the first parameter
5. **Error Handling**: Return meaningful errors that can be handled by the caller

## Related Components

- [Domain Entities](../entity/README.md) - The entities managed by these repositories
- [Domain Services](../services/README.md) - Services that use these repositories
- [Infrastructure Adapters](../../../infrastructure/adapters/repository/README.md) - Implementations of these interfaces
- [Mock Implementations](./mock/README.md) - Mock implementations for testing