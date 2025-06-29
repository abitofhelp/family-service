# GraphQL Resolvers

## Overview

The GraphQL Resolvers package implements the resolvers for the Family Service GraphQL API. In GraphQL, resolvers are functions that resolve (fetch or compute) the data for fields in a GraphQL query. This package contains the resolver implementations that translate GraphQL requests into calls to application services and convert the results back into the format expected by GraphQL. This approach keeps the GraphQL-specific code isolated in this package, allowing the application services to remain focused on business logic without any knowledge of GraphQL.

## Features

- **Query Operations**: Resolvers for retrieving family data
- **Mutation Operations**: Resolvers for creating, updating, and deleting family data
- **Type Resolvers**: Resolvers for specific GraphQL types like Family, Parent, and Child
- **Authentication**: Authentication and authorization for GraphQL operations
- **Error Handling**: Proper error handling and translation to GraphQL errors
- **Context Propagation**: Context propagation for request-scoped data
- **Dependency Injection**: Clean dependency management through constructor injection

## Installation

```bash
go get github.com/abitofhelp/family-service/interface/adapters/graphql/resolver
```

## Quick Start

See the [Quick Start example](../../../../EXAMPLES/graphql/basic_usage/README.md) for a complete, runnable example of how to use the GraphQL resolvers.

## Configuration

The GraphQL Resolvers package can be configured with the following options:

- **Authentication**: Configure authentication and authorization rules for GraphQL operations
- **Validation**: Configure validation rules for GraphQL inputs
- **Error Handling**: Configure how errors are handled and translated to GraphQL errors
- **Logging**: Configure logging for GraphQL operations
- **Tracing**: Configure distributed tracing for GraphQL operations

Example configuration:

```go
package main

import (
    "github.com/99designs/gqlgen/graphql/handler"
    "github.com/99designs/gqlgen/graphql/handler/extension"
    "github.com/abitofhelp/family-service/interface/adapters/graphql/generated"
    "github.com/abitofhelp/family-service/interface/adapters/graphql/resolver"
    "github.com/abitofhelp/family-service/interface/adapters/graphql/directives"
    "github.com/hashicorp/golang-lru"
)

func main() {
    // Configure the resolver with dependencies
    resolver := resolver.NewResolver(
        familyService,  // Application service for family operations
        familyMapper,   // Mapper for converting between GraphQL and domain models
    )

    // Configure the GraphQL server with the resolver
    srv := handler.NewDefaultServer(
        generated.NewExecutableSchema(
            generated.Config{
                Resolvers: resolver,
                Directives: generated.DirectiveRoot{
                    IsAuthorized: directives.IsAuthorized,
                },
            },
        ),
    )

    // Configure middleware for the GraphQL server
    srv.Use(extension.Introspection{})
    srv.Use(extension.AutomaticPersistedQuery{
        Cache: lru.New(100),
    })
}
```

## API Documentation

### Core Types

#### Resolver

The Resolver struct serves as a dependency injection container for GraphQL resolvers.

```
// Resolver serves as a dependency injection container for GraphQL resolvers.
//
// This struct holds references to the application services and other dependencies
// needed by the resolvers. It follows the Dependency Injection pattern, where
// dependencies are provided from the outside rather than created internally.
//
// The Resolver acts as a facade between the GraphQL layer and the application layer,
// delegating client requests to the appropriate application services.
type Resolver struct {
    familyService ports.FamilyApplicationService // Application service for family operations
    mapper       dto.FamilyMapper               // Mapper for converting between GraphQL and domain models
}
```

#### Query Resolver

The queryResolver struct implements the QueryResolver interface for handling GraphQL queries.

```
// queryResolver implements the QueryResolver interface for handling GraphQL queries.
type queryResolver struct {
    *Resolver // Embeds the main Resolver for access to dependencies
}
```

#### Mutation Resolver

The mutationResolver struct implements the MutationResolver interface for handling GraphQL mutations.

```
// mutationResolver implements the MutationResolver interface for handling GraphQL mutations.
type mutationResolver struct {
    *Resolver // Embeds the main Resolver for access to dependencies
}
```

### Key Methods

#### NewResolver

Creates a new resolver with the given dependencies.

```
// NewResolver creates a new resolver with the given dependencies.
//
// This function creates a new Resolver instance with the provided dependencies.
// It follows the Dependency Injection pattern, requiring all dependencies
// to be provided rather than creating them internally.
func NewResolver(familyService ports.FamilyApplicationService, mapper dto.FamilyMapper) *Resolver
```

#### Query

Returns the query resolver implementation.

```
// Query returns the query resolver implementation.
//
// This method returns a resolver for GraphQL query operations.
// In GraphQL, queries are used to retrieve data (similar to GET in REST).
func (r *Resolver) Query() generated.QueryResolver
```

#### Mutation

Returns the mutation resolver implementation.

```
// Mutation returns the mutation resolver implementation.
//
// This method returns a resolver for GraphQL mutation operations.
// In GraphQL, mutations are used to modify data (similar to POST, PUT, DELETE in REST).
func (r *Resolver) Mutation() generated.MutationResolver
```

## Examples

Example of using the GraphQL API:

```graphql
# Query all families
query GetAllFamilies {
  getAllFamilies {
    id
    status
    parents {
      id
      firstName
      lastName
    }
    children {
      id
      firstName
      lastName
    }
  }
}

# Create a new family
mutation CreateFamily {
  createFamily(input: {
    id: "fam-123"
    status: SINGLE
    parents: [{
      id: "par-123"
      firstName: "John"
      lastName: "Doe"
      birthDate: "1980-01-01T00:00:00Z"
    }]
    children: []
  }) {
    id
    status
    parents {
      id
      firstName
      lastName
    }
  }
}
```

## Best Practices

1. **Separation of Concerns**: Keep GraphQL-specific code isolated in resolvers
2. **Dependency Injection**: Use constructor injection for dependencies
3. **Error Handling**: Properly translate domain errors to GraphQL errors
4. **Context Propagation**: Propagate context throughout the request lifecycle
5. **Type Safety**: Use strong typing for all resolver parameters and return values

## Troubleshooting

### Common Issues

#### Authentication Errors

If you encounter authentication errors, ensure that the authentication directives are properly configured and that the user has the necessary permissions.

#### Performance Issues

If you encounter performance issues, consider using DataLoader to batch and cache database queries, and ensure that resolvers are efficiently implemented.

## Related Components

- [Application Services](../../../../core/application/services/README.md) - The application services used by these resolvers
- [GraphQL DTOs](../dto/README.md) - The DTOs used for mapping between GraphQL and domain models
- [GraphQL Generated Code](../generated/README.md) - The generated GraphQL code used by these resolvers

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../../LICENSE) file for details.
