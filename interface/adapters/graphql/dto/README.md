# GraphQL DTO Package

This package provides Data Transfer Object (DTO) mapping functionality between domain DTOs and GraphQL models.

## Overview

The DTO package implements the mapping layer that ensures proper separation between the domain layer and the GraphQL interface layer. This separation is crucial for maintaining clean architecture boundaries and preventing domain model leaks.

## Components

### FamilyMapper

The `FamilyMapper` struct provides bidirectional mapping between domain DTOs and GraphQL models:

- `ToGraphQL`: Converts domain FamilyDTO to GraphQL Family model
- `ToDomain`: Converts GraphQL FamilyInput to domain FamilyDTO

## Usage

```go
mapper := dto.NewFamilyMapper()

// Convert domain DTO to GraphQL model
graphqlModel, err := mapper.ToGraphQL(domainDTO)

// Convert GraphQL input to domain DTO
domainDTO, err := mapper.ToDomain(graphqlInput)
```

## Design Decisions

1. **Separation of Concerns**: The mapper isolates all conversion logic in one place, making it easier to maintain and modify.
2. **Error Handling**: All conversion methods return errors to ensure proper error handling at the interface layer.
3. **Immutability**: The mapper creates new instances rather than modifying existing ones.
4. **Type Safety**: Strong typing is used throughout to catch potential issues at compile time.

## Best Practices

1. Always use the mapper for conversions between domain and GraphQL types
2. Handle all errors returned by the mapper methods
3. Don't modify the returned objects directly; create new instances if needed
4. Keep the mapping logic simple and focused on data conversion only
