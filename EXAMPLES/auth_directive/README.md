# Auth Directive Example

## Overview

This example demonstrates how to use authentication with the Family Service GraphQL API. It shows how to make GraphQL requests with and without an authorization token.

## Features

- **Authentication Testing**: Tests GraphQL queries with and without authentication
- **JWT Token Usage**: Shows how to include a JWT token in the Authorization header
- **Error Handling**: Demonstrates how to handle authentication errors
- **Environment Variables**: Uses environment variables for token configuration

## Running the Example

To run this example, navigate to this directory and execute:

```bash
# Set your JWT token as an environment variable
export AUTH_TOKEN="your.jwt.token.here"

# Run the example
go run main.go
```

## Code Walkthrough

### GraphQL Request Setup

The example creates a GraphQL request with a query to fetch all families:

```go
// Create the request
req := GraphQLRequest{
    Query:         query,
    OperationName: operationName,
}

// Convert request to JSON
reqBody, err := json.Marshal(req)
```

### Authentication Header

The example shows how to add the JWT token to the Authorization header:

```go
// Set headers
httpReq.Header.Set("Content-Type", "application/json")
if token != "" {
    httpReq.Header.Set("Authorization", "Bearer "+token)
}
```

### Testing Without Authentication

The example first tests a query without any authentication:

```go
// Test with no authorization
fmt.Println("Testing with no authorization:")
testQuery(url, "", "GetAllFamilies", `
    query GetAllFamilies {
        getAllFamilies {
            id
            status
        }
    }
`)
```

## Expected Output

When run without a valid token:

```
Testing with no authorization:
Status: 200 OK
Errors:
  - Access denied: User not authenticated

No AUTH_TOKEN environment variable found. Set it to a valid JWT token for testing with authorization.
```

When run with a valid token:

```
Testing with no authorization:
Status: 200 OK
Errors:
  - Access denied: User not authenticated

Testing with authorization:
Status: 200 OK
Data received successfully
```

## Related Examples

- [Errors Example](../errors/README.md) - Shows how to handle GraphQL errors
- [Family Errors Example](../family_errors/README.md) - Shows how to handle domain-specific errors

## Related Components

- [GraphQL Resolver](../../interface/adapters/graphql/resolver/README.md) - The GraphQL resolver implementation
- [Auth Middleware](../../infrastructure/adapters/security/README.md) - The authentication middleware

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.