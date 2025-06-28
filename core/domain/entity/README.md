# Domain Entities

## Overview

The Domain Entities package contains the core business objects of the Family Service application. These entities encapsulate the state and behavior of the domain model, implementing business rules and invariants. The package follows Domain-Driven Design principles, with Family as the root aggregate that maintains consistency boundaries.

## Features

- **Rich Domain Model**: Implements a comprehensive model of family relationships
- **Invariant Enforcement**: Ensures business rules are consistently applied
- **Value Objects**: Uses value objects for identifiers and other immutable concepts
- **Aggregate Roots**: Implements Family as an aggregate root to maintain consistency
- **Domain Events**: Supports domain events for important state changes
- **Data Transfer Objects**: Provides DTOs for transferring data between layers

## API Documentation

### Core Types

#### Family

The Family aggregate represents a family unit with parents and children. It enforces business rules such as requiring at least one parent, preventing duplicate parents, and ensuring proper family status based on composition.

```
// Family is the root aggregate that represents a family unit
type Family struct {
    id       identificationwrapper.ID
    status   Status
    parents  []*Parent
    children []*Child
}
```

#### Parent

The Parent entity represents a parent in a family. It includes personal information and can be marked as deceased.

```
// Parent represents a parent in a family
type Parent struct {
    id        identificationwrapper.ID
    firstName string
    lastName  string
    birthDate time.Time
    deathDate *time.Time
}
```

#### Child

The Child entity represents a child in a family. It includes personal information similar to a parent.

```
// Child represents a child in a family
type Child struct {
    id        identificationwrapper.ID
    firstName string
    lastName  string
    birthDate time.Time
}
```

### Key Methods

#### NewFamily

Creates a new Family aggregate with validation.

```
// NewFamily creates a new Family aggregate with validation
func NewFamily(id string, status Status, parents []*Parent, children []*Child) (*Family, error)
```

#### Divorce

Handles the divorce process, creating a new family for the non-custodial parent.

```
// Divorce handles the divorce process, creating a new family for the remaining parent
func (f *Family) Divorce(custodialParentID string) (*Family, error)
```

#### MarkParentDeceased

Marks a parent as deceased and updates family status if needed.

```
// MarkParentDeceased marks a parent as deceased and updates family status if needed
func (f *Family) MarkParentDeceased(parentID string, deathDate time.Time) error
```

## Best Practices

1. **Encapsulation**: Keep entity state private and provide methods for manipulation
2. **Validation**: Validate entities at creation and when state changes
3. **Immutability**: Use immutable value objects where appropriate
4. **Rich Behavior**: Implement domain logic in entity methods rather than external services
5. **Consistency Boundaries**: Use aggregates to maintain consistency boundaries

## Related Components

- [Domain Services](../services/README.md) - Services that operate on these entities
- [Domain Validation](../validation/README.md) - Validation rules for these entities
- [Domain Errors](../errors/README.md) - Domain-specific error types
- [Application Services](../../application/services/README.md) - Uses these entities to implement use cases