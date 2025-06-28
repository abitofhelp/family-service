# Family Service Examples

This directory contains example programs that demonstrate various aspects of the Family Service and ServiceLib integration.

## Error Handling Examples

### Basic Error Handling

The `errors` directory contains a simple example that demonstrates how to use the ServiceLib error types:

```bash
cd errors
go run main.go
```

This will output examples of DatabaseError, NotFoundError, and ValidationError.

### Family-Specific Error Handling

The `family_errors` directory contains a more comprehensive example that demonstrates family-specific error handling:

```bash
cd family_errors
go run main.go
```

This will output examples of various error types in the context of the Family Service.

## Authentication Example

The `auth_directive` directory contains an example that demonstrates how to use the GraphQL authentication directive:

```bash
cd auth_directive
go run main.go
```

This will attempt to make GraphQL queries with and without authentication. To test with authentication, set the `AUTH_TOKEN` environment variable to a valid JWT token:

```bash
export AUTH_TOKEN="your.jwt.token"
cd auth_directive
go run main.go
```

## Note

These examples are provided for demonstration purposes and are not part of the main application. They show how to use various features of the Family Service and ServiceLib.