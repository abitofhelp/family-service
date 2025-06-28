# Application Services

## Overview

The Application Services package implements the application layer in the Clean Architecture pattern. It contains services that orchestrate the execution of use cases by coordinating domain services, repositories, and other infrastructure components. These services act as a bridge between the interface layer and the domain layer, ensuring that business rules are properly applied while handling cross-cutting concerns like logging, caching, and error handling.

## Features

- **Family Application Service**: Implements operations for managing families, parents, and children
- **Caching Integration**: Uses caching to improve performance for frequently accessed data
- **Comprehensive Logging**: Detailed logging of all operations for observability
- **Error Handling**: Proper error handling and propagation
- **Transaction Management**: Ensures data consistency across operations
- **Clean Architecture Compliance**: Follows the principles of Clean Architecture

## API Documentation

### Core Types

#### FamilyApplicationService

The FamilyApplicationService implements the application service for family-related use cases. It provides methods for creating and managing families, adding and removing parents and children, handling divorces, and finding families by parent or child.

```
// FamilyApplicationService implements the application service for family-related use cases
type FamilyApplicationService struct {
    BaseApplicationService[*entity.Family, *entity.FamilyDTO]
    familyService *domainservices.FamilyDomainService
    familyRepo    domainports.FamilyRepository
    logger        *logging.ContextLogger
    cache         *cache.Cache
}
```

### Key Methods

#### Create

Creates a new family.

```
// Create creates a new family
func (s *FamilyApplicationService) Create(ctx context.Context, dto *entity.FamilyDTO) (*entity.FamilyDTO, error)
```

#### AddParent

Adds a parent to a family.

```
// AddParent adds a parent to a family
func (s *FamilyApplicationService) AddParent(ctx context.Context, familyID string, parentDTO entity.ParentDTO) (*entity.FamilyDTO, error)
```

#### Divorce

Handles the divorce process, creating a new family for the custodial parent and children.

```
// Divorce handles the divorce process
func (s *FamilyApplicationService) Divorce(ctx context.Context, familyID string, custodialParentID string) (*entity.FamilyDTO, error)
```

## Best Practices

1. **Separation of Concerns**: Application services should focus on orchestrating use cases, delegating domain logic to domain services
2. **Comprehensive Logging**: Log all operations with appropriate context for observability
3. **Proper Error Handling**: Handle errors appropriately and provide meaningful error messages
4. **Caching Strategy**: Use caching for frequently accessed data to improve performance
5. **Transaction Management**: Ensure data consistency across operations

## Related Components

- [Application Ports](../ports/README.md) - Defines the interfaces implemented by these services
- [Domain Services](../../domain/services/README.md) - Provides the domain logic used by these services
- [Domain Entities](../../domain/entity/README.md) - Defines the entity types used by these services
- [Repositories](../../domain/ports/README.md) - Provides data access for these services