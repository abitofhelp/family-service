# GraphQL DTO Package

## Overview

The DTO package implements the mapping layer that ensures proper separation between the domain layer and the GraphQL interface layer. This separation is crucial for maintaining clean architecture boundaries and preventing domain model leaks.

## Features

- **Bidirectional Mapping**: Convert between domain DTOs and GraphQL models in both directions
- **Type Safety**: Strong typing throughout to catch potential issues at compile time
- **Error Handling**: Comprehensive error handling for all conversion operations
- **Clean Architecture Compliance**: Ensures domain layer doesn't depend on interface layer
- **Immutability**: Creates new instances rather than modifying existing ones

## Installation

```bash
go get github.com/abitofhelp/family-service/interface/adapters/graphql/dto
```

## Quick Start

See the [Quick Start example](../../../../EXAMPLES/graphql/basic_usage/README.md) for a complete, runnable example of how to use the GraphQL DTO package.

## Configuration

The GraphQL DTO package doesn't require any specific configuration. It provides a set of mappers that can be used as-is. However, you can configure the mapping behavior by:

- **Custom Mappers**: Create custom mappers for new entity types
- **Validation Rules**: Add custom validation rules for the mapping process
- **Error Handling**: Configure how errors are handled during mapping

## API Documentation

### Core Types

Description of the main types provided by the component.

#### FamilyMapper

The `FamilyMapper` struct provides bidirectional mapping between domain DTOs and GraphQL models.

```
// FamilyMapper provides bidirectional mapping between domain DTOs and GraphQL models
type FamilyMapper struct {
    // Fields
}
```

#### ParentMapper

The `ParentMapper` struct provides bidirectional mapping between domain Parent DTOs and GraphQL Parent models.

```
// ParentMapper provides bidirectional mapping between domain Parent DTOs and GraphQL Parent models
type ParentMapper struct {
    // Fields
}
```

### Key Methods

Description of the key methods provided by the component.

#### ToGraphQL

Converts domain DTOs to GraphQL models.

```
// ToGraphQL converts a domain DTO to a GraphQL model
func (m *FamilyMapper) ToGraphQL(dto *domain.FamilyDTO) (*model.Family, error)
```

#### ToDomain

Converts GraphQL inputs to domain DTOs.

```
// ToDomain converts a GraphQL input to a domain DTO
func (m *FamilyMapper) ToDomain(input *model.FamilyInput) (*domain.FamilyDTO, error)
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory.

Example of using the FamilyMapper:

```go
package main

import (
    "fmt"
    "github.com/abitofhelp/family-service/interface/adapters/graphql/dto"
    "github.com/abitofhelp/family-service/core/domain/entity"
    "github.com/abitofhelp/family-service/interface/adapters/graphql/model"
)

func main() {
    // Create a mapper
    mapper := dto.NewFamilyMapper()

    // Convert domain DTO to GraphQL model
    graphqlModel, err := mapper.ToGraphQL(domainDTO)
    if err != nil {
        fmt.Printf("Error converting to GraphQL: %v\n", err)
        return
    }

    // Convert GraphQL input to domain DTO
    domainDTO, err := mapper.ToDomain(graphqlInput)
    if err != nil {
        fmt.Printf("Error converting to domain: %v\n", err)
        return
    }
}
```

## Best Practices

1. **Use Mappers Consistently**: Always use the mapper for conversions between domain and GraphQL types
2. **Handle All Errors**: Always check and handle errors returned by mapper methods
3. **Immutability**: Don't modify the returned objects directly; create new instances if needed
4. **Separation of Concerns**: Keep mapping logic separate from business logic
5. **Type Safety**: Use strong typing to catch potential issues at compile time

## Troubleshooting

### Common Issues

#### Type Conversion Errors

If you encounter type conversion errors, ensure that the source and target types are compatible. The mapper expects specific types and will return errors if the types don't match.

#### Null or Empty Values

The mapper handles null or empty values gracefully. If you encounter issues with null values, check that the mapper is correctly handling these cases.

## Related Components

- [Domain Entities](../../../../core/domain/entity/README.md) - The domain entities that are mapped
- [GraphQL Resolvers](../resolver/README.md) - Uses these mappers to handle GraphQL requests
- [GraphQL Models](../model/README.md) - The GraphQL models that are mapped to/from domain entities

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../../LICENSE) file for details.
