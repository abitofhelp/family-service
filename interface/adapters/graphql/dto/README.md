# GraphQL DTO Package

## Overview

The DTO package implements the mapping layer that ensures proper separation between the domain layer and the GraphQL interface layer. This separation is crucial for maintaining clean architecture boundaries and preventing domain model leaks.

## Architecture

The GraphQL DTO package follows the Adapter pattern from Hexagonal Architecture, providing a layer of abstraction between the domain layer and the GraphQL interface layer. This ensures that the domain layer doesn't directly depend on GraphQL-specific types, maintaining the dependency inversion principle.

The package sits in the interface layer of the application and is used by the GraphQL resolvers to convert between domain DTOs and GraphQL models. The architecture follows these principles:

- **Dependency Inversion**: The domain layer depends on abstractions, not concrete implementations
- **Adapter Pattern**: This package adapts the domain types to the GraphQL interface
- **Data Transfer Object Pattern**: Uses DTOs to transfer data between layers

## Implementation Details

The GraphQL DTO package implements the following design patterns:

1. **Adapter Pattern**: Adapts the domain types to the GraphQL interface
2. **Data Transfer Object Pattern**: Uses DTOs to transfer data between layers
3. **Factory Pattern**: Factory methods create new instances of GraphQL models and domain DTOs

Key implementation details:

- **Bidirectional Mapping**: The mapper provides methods for converting in both directions
- **Error Handling**: All conversion methods return errors to ensure proper error handling
- **Immutability**: The mapper creates new instances rather than modifying existing ones
- **Type Safety**: Strong typing is used throughout to catch potential issues at compile time

## Components

### FamilyMapper

The `FamilyMapper` struct provides bidirectional mapping between domain DTOs and GraphQL models:

- `ToGraphQL`: Converts domain FamilyDTO to GraphQL Family model
- `ToDomain`: Converts GraphQL FamilyInput to domain FamilyDTO

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [GraphQL API Example](../../../examples/graphql_api/README.md) - Shows how to use the GraphQL DTO package with the GraphQL API

Example of using the FamilyMapper:

```
// Create a mapper
mapper := NewFamilyMapper()

// Convert domain DTO to GraphQL model
graphqlModel, err := mapper.ToGraphQL(domainDTO)
if err != nil {
    // Handle error
}

// Convert GraphQL input to domain DTO
domainDTO, err := mapper.ToDomain(graphqlInput)
if err != nil {
    // Handle error
}
```

## Configuration

The GraphQL DTO package doesn't require any specific configuration. It provides a set of mappers that can be used as-is. However, you can configure the mapping behavior by:

- **Custom Mappers**: Create custom mappers for new entity types
- **Validation Rules**: Add custom validation rules for the mapping process
- **Error Handling**: Configure how errors are handled during mapping

## Testing

The GraphQL DTO package is tested through:

1. **Unit Tests**: Each mapper has comprehensive unit tests
2. **Integration Tests**: Tests that verify the mappers work correctly with the GraphQL resolvers
3. **Edge Case Tests**: Tests that verify the mappers handle edge cases correctly

Key testing approaches:

- **Bidirectional Testing**: Tests that verify mapping works correctly in both directions
- **Error Case Testing**: Tests that verify proper error handling for invalid inputs
- **Null/Empty Testing**: Tests that verify proper handling of null or empty values
- **Type Conversion Testing**: Tests that verify proper type conversion between domain and GraphQL types

Example of a test case:

```
// Example test for the FamilyMapper
func TestFamilyMapper_ToGraphQL(t *testing.T) {
    // Create a domain DTO
    domainDTO := createTestFamilyDTO()

    // Create a mapper
    mapper := NewFamilyMapper()

    // Convert to GraphQL model
    graphqlModel, err := mapper.ToGraphQL(domainDTO)

    // Verify the conversion
    assert.NoError(t, err)
    assert.Equal(t, domainDTO.ID, graphqlModel.ID)
    assert.Equal(t, string(domainDTO.Status), string(graphqlModel.Status))
    assert.Len(t, graphqlModel.Parents, 1)
    assert.Equal(t, domainDTO.Parents[0].ID, graphqlModel.Parents[0].ID)
}
```

## Design Notes

1. **Separation of Concerns**: The mapper isolates all conversion logic in one place, making it easier to maintain and modify.
2. **Error Handling**: All conversion methods return errors to ensure proper error handling at the interface layer.
3. **Immutability**: The mapper creates new instances rather than modifying existing ones.
4. **Type Safety**: Strong typing is used throughout to catch potential issues at compile time.
5. **Clean Architecture Compliance**: The package follows Clean Architecture principles by ensuring that the domain layer doesn't depend on the interface layer.
6. **Testability**: The package is designed to be easily testable, with clear interfaces and dependencies.

## Best Practices

1. Always use the mapper for conversions between domain and GraphQL types
2. Handle all errors returned by the mapper methods
3. Don't modify the returned objects directly; create new instances if needed
4. Keep the mapping logic simple and focused on data conversion only

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Data Transfer Object Pattern](https://martinfowler.com/eaaCatalog/dataTransferObject.html)
- [GraphQL](https://graphql.org/)
- [Domain Entities](../../../../core/domain/entity/README.md) - The domain entities that are mapped
- [GraphQL Resolvers](../resolver/README.md) - Uses these mappers to handle GraphQL requests
- [GraphQL Models](../model/README.md) - The GraphQL models that are mapped to/from domain entities
