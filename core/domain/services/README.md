# Domain Services

## Overview

The Domain Services package contains services that coordinate operations on domain entities and aggregates. These services encapsulate complex business logic that doesn't naturally fit within a single entity or aggregate. The FamilyDomainService, in particular, orchestrates operations on the Family aggregate, ensuring that business rules are properly applied and that the aggregate remains in a consistent state.

## Features

- **Family Domain Service**: Coordinates operations on the Family aggregate
- **Transaction Management**: Ensures data consistency across complex operations
- **Business Rule Enforcement**: Applies domain-specific business rules
- **Distributed Tracing**: Integrates with OpenTelemetry for distributed tracing
- **Comprehensive Logging**: Provides detailed logging of all operations
- **Metrics Collection**: Collects metrics for monitoring application behavior
- **Error Handling**: Properly handles and propagates domain-specific errors

## API Documentation

### Core Types

#### FamilyDomainService

The FamilyDomainService coordinates operations on the Family aggregate, ensuring that business rules are properly applied.

```
// FamilyDomainService is a domain service that coordinates operations on the Family aggregate
type FamilyDomainService struct {
    repo   ports.FamilyRepository
    logger *loggingwrapper.ContextLogger
    tracer trace.Tracer
}
```

### Key Methods

#### CreateFamily

Creates a new family with validation.

```
// CreateFamily creates a new family
func (s *FamilyDomainService) CreateFamily(ctx context.Context, dto entity.FamilyDTO) (*entity.FamilyDTO, error)
```

#### AddParent

Adds a parent to a family, updating the family status if necessary.

```
// AddParent adds a parent to a family
func (s *FamilyDomainService) AddParent(ctx context.Context, familyID string, parentDTO entity.ParentDTO) (*entity.FamilyDTO, error)
```

#### Divorce

Handles the divorce process, creating a new family for the non-custodial parent.

```
// Divorce handles the divorce process
func (s *FamilyDomainService) Divorce(ctx context.Context, familyID string, custodialParentID string) (*entity.FamilyDTO, error)
```

## Best Practices

1. **Single Responsibility**: Each domain service should focus on a specific domain concept
2. **Statelessness**: Domain services should be stateless, with all state maintained in the entities
3. **Dependency Injection**: Use dependency injection to provide repositories and other dependencies
4. **Comprehensive Logging**: Log all operations with appropriate context for observability
5. **Proper Error Handling**: Handle errors appropriately and provide meaningful error messages

## Related Components

- [Domain Entities](../entity/README.md) - The entities managed by these services
- [Domain Ports](../ports/README.md) - The interfaces used by these services for data access
- [Domain Errors](../errors/README.md) - Domain-specific error types used by these services
- [Application Services](../../application/services/README.md) - Higher-level services that use these domain services