# Domain Entities

## Overview

The Domain Entities package contains the core business objects of the Family Service application. These entities encapsulate the state and behavior of the domain model, implementing business rules and invariants. The package follows Domain-Driven Design principles, with Family as the root aggregate that maintains consistency boundaries.

## Architecture

The Domain Entities package is part of the core domain layer in the Clean Architecture and Hexagonal Architecture patterns. It sits at the center of the application and has no dependencies on other layers. The architecture follows these principles:

- **Domain-Driven Design (DDD)**: Entities are designed based on the ubiquitous language of the domain
- **Aggregate Pattern**: Family is an aggregate root that maintains consistency boundaries
- **Value Objects**: Immutable objects that represent concepts with no identity
- **Entity Pattern**: Objects with identity and lifecycle
- **Rich Domain Model**: Business logic is encapsulated in the entities themselves

The package is organized into:

- **Entities**: Family, Parent, Child
- **Value Objects**: Status, ID
- **Data Transfer Objects (DTOs)**: For transferring data between layers
- **Factories**: For creating valid entities

## Implementation Details

The Domain Entities package implements the following design patterns:

1. **Aggregate Root Pattern**: Family is an aggregate root that maintains consistency boundaries
2. **Factory Pattern**: Factory methods create valid entities with proper validation
3. **Value Object Pattern**: Immutable objects for concepts like Status
4. **Builder Pattern**: Used for constructing complex entities
5. **Visitor Pattern**: For operations that need to traverse the entity hierarchy

Key implementation details:

- **Private Fields**: All entity fields are private to enforce encapsulation
- **Constructor Validation**: Entities validate their state during construction
- **Immutable Value Objects**: Value objects are immutable to prevent unexpected state changes
- **Domain Events**: Important state changes trigger domain events
- **Defensive Programming**: Methods validate parameters and state before making changes
- **Error Types**: Domain-specific error types for clear error handling

## Features

- **Rich Domain Model**: Implements a comprehensive model of family relationships
- **Invariant Enforcement**: Ensures business rules are consistently applied
- **Value Objects**: Uses value objects for identifiers and other immutable concepts
- **Aggregate Roots**: Implements Family as an aggregate root to maintain consistency
- **Domain Events**: Supports domain events for important state changes
- **Data Transfer Objects**: Provides DTOs for transferring data between layers

## Examples

There may be additional examples in the /EXAMPLES directory.

Example of creating a new family:

```
// Create parents
parent1, err := NewParent("par-123", "John", "Doe", birthDate, nil)
if err != nil {
    // Handle error
}

parent2, err := NewParent("par-456", "Jane", "Doe", birthDate, nil)
if err != nil {
    // Handle error
}

// Create a child
child, err := NewChild("chi-123", "Jimmy", "Doe", birthDate)
if err != nil {
    // Handle error
}

// Create a family
family, err := NewFamily("fam-123", StatusMarried, []*Parent{parent1, parent2}, []*Child{child})
if err != nil {
    // Handle error
}
```

## Configuration

The Domain Entities package doesn't require any specific configuration as it contains pure domain logic. However, it does have some configurable aspects:

- **ID Generation**: The package uses the identificationwrapper package for ID generation, which can be configured
- **Validation Rules**: Validation rules are defined in the validation package and can be configured
- **Date Handling**: Date handling uses the datewrapper package, which can be configured for different date formats

## Testing

The Domain Entities package is tested through:

1. **Unit Tests**: Each entity and value object has comprehensive unit tests
2. **Property-Based Testing**: Tests with randomized inputs to find edge cases
3. **Scenario-Based Testing**: Tests that simulate real-world scenarios

Key testing approaches:

- **Invariant Testing**: Tests that verify business rules are enforced
- **Boundary Testing**: Tests at the boundaries of valid input
- **Error Case Testing**: Tests that verify proper error handling
- **Lifecycle Testing**: Tests that verify entity lifecycle (creation, modification, deletion)

Example of a test case:

```
func TestFamily_Divorce(t *testing.T) {
    // Setup test data
    family := createTestFamily(t)

    // Execute the operation
    newFamily, err := family.Divorce("par-123")

    // Verify results
    assert.NoError(t, err)
    assert.Equal(t, StatusDivorced, family.Status())
    assert.Equal(t, StatusSingle, newFamily.Status())
    assert.Len(t, family.Parents(), 1)
    assert.Len(t, newFamily.Parents(), 1)
}
```

## Design Notes

1. **Encapsulation**: All entity state is private and can only be modified through methods
2. **Validation**: Entities validate their state during construction and when state changes
3. **Immutability**: Value objects are immutable to prevent unexpected state changes
4. **Rich Behavior**: Domain logic is implemented in entity methods rather than external services
5. **Consistency Boundaries**: Aggregates maintain consistency boundaries
6. **Identity**: Entities have identity that persists across state changes
7. **Value Objects**: Concepts with no identity are implemented as value objects
8. **Domain Events**: Important state changes trigger domain events for cross-aggregate consistency

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

## References

- [Domain-Driven Design](https://domainlanguage.com/ddd/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Value Objects](https://martinfowler.com/bliki/ValueObject.html)
- [Aggregate Pattern](https://martinfowler.com/bliki/DDD_Aggregate.html)
- [Domain Services](../services/README.md) - Services that operate on these entities
- [Domain Validation](../validation/README.md) - Validation rules for these entities
- [Domain Errors](../errors/README.md) - Domain-specific error types
- [Application Services](../../application/services/README.md) - Uses these entities to implement use cases
