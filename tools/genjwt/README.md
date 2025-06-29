# JWT Token Generator

## Overview

The JWT Token Generator is a utility tool for generating JSON Web Tokens (JWTs) for testing and development purposes. It creates tokens with different roles (admin, editor, viewer) and permissions, which can be used to test authentication and authorization in the Family Service application. The tool also validates the generated tokens to ensure they are correctly formatted and contain the expected claims.

## Architecture

The JWT Token Generator is a standalone command-line tool that uses the servicelib/auth package to generate and validate JWT tokens. It follows these principles:

- **Separation of Concerns**: The tool focuses solely on token generation and validation
- **Configuration-Driven**: Token properties are configured through a configuration object
- **Role-Based Access Control**: Tokens are generated with different roles and permissions

The tool is organized into:

- **Token Generation**: Functions for generating tokens with different roles and permissions
- **Token Validation**: Functions for validating tokens and extracting claims
- **Configuration**: Configuration options for token generation

## Implementation Details

The JWT Token Generator implements the following design patterns:

1. **Factory Pattern**: Creates tokens with specific properties
2. **Builder Pattern**: Uses a configuration object to build tokens
3. **Command Pattern**: Executes specific commands for token generation and validation

Key implementation details:

- **JWT Standard**: Uses the JWT standard for token generation
- **Role-Based Access Control**: Supports different roles (admin, editor, viewer)
- **Scope-Based Permissions**: Supports different scopes (READ, WRITE, DELETE, CREATE)
- **Resource-Based Permissions**: Supports different resources (FAMILY, PARENT, CHILD)
- **Token Validation**: Validates tokens and extracts claims
- **Error Handling**: Provides clear error messages for token generation and validation failures

## Features

- **Multiple Role Support**: Generates tokens for admin, editor, and viewer roles
- **Customizable Permissions**: Supports different scopes and resources
- **Token Validation**: Validates tokens and extracts claims
- **Error Handling**: Provides clear error messages for token generation and validation failures
- **Command-Line Interface**: Simple command-line interface for token generation
- **Integration with servicelib/auth**: Uses the servicelib/auth package for token generation and validation

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Auth Directive Example](../../EXAMPLES/auth_directive/README.md) - Shows how to use JWT tokens for authentication in GraphQL

Example of generating a token with custom roles, scopes, and resources:

```
// Create an auth instance
authInstance, err := auth.New(ctx, config, logger)
if err != nil {
    logger.Fatal("Failed to create auth instance", zap.Error(err))
}

// Generate a custom token
customScopes := []string{"READ", "WRITE"}
customResources := []string{"FAMILY"}
customToken, err := authInstance.GenerateToken(ctx, "custom", []string{"CUSTOM_ROLE"}, customScopes, customResources)
if err != nil {
    logger.Fatal("Failed to generate custom token", zap.Error(err))
}

// Validate the token
claims, err := authInstance.ValidateToken(ctx, customToken)
if err != nil {
    logger.Fatal("Failed to validate custom token", zap.Error(err))
}

// Use the claims
fmt.Printf("Valid Custom Token, Claims: %+v\n", claims)
```

## Usage

To use the JWT Token Generator, simply run the tool from the command line:

```bash
go run tools/genjwt/main.go
```

This will generate three tokens:

1. **Admin Token**: A token with the ADMIN role and all scopes (READ, WRITE, DELETE, CREATE) for all resources (FAMILY, PARENT, CHILD)
2. **Editor Token**: A token with the EDITOR role and all scopes (READ, WRITE, DELETE, CREATE) for all resources (FAMILY, PARENT, CHILD)
3. **Viewer Token**: A token with the VIEWER role and only the READ scope for all resources (FAMILY, PARENT, CHILD)

The tool will also validate each token and display the extracted claims.

## Example Output

```
Admin Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

Editor Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

Viewer Token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

Valid Admin Token, Claims: {Subject:admin Roles:[ADMIN] Scopes:[READ WRITE DELETE CREATE] Resources:[FAMILY PARENT CHILD] ExpiresAt:1625097600 IssuedAt:1625011200 Issuer:family-service}

Valid Editor Token, Claims: {Subject:editor Roles:[EDITOR] Scopes:[READ WRITE DELETE CREATE] Resources:[FAMILY PARENT CHILD] ExpiresAt:1625097600 IssuedAt:1625011200 Issuer:family-service}

Valid Viewer Token, Claims: {Subject:viewer Roles:[VIEWER] Scopes:[READ] Resources:[FAMILY PARENT CHILD] ExpiresAt:1625097600 IssuedAt:1625011200 Issuer:family-service}
```

## Configuration

The JWT Token Generator uses a configuration object to configure token generation. The following configuration options are available:

- **Secret Key**: The secret key used to sign the tokens
- **Token Duration**: The duration for which the tokens are valid
- **Issuer**: The issuer of the tokens

Example configuration:

```
// Example configuration
config := auth.DefaultConfig()
config.JWT.SecretKey = "01234567890123456789012345678901"
```

## Customization

You can customize the token generation by modifying the following parameters in the code:

- **Roles**: The roles assigned to the tokens (e.g., ADMIN, EDITOR, VIEWER)
- **Scopes**: The scopes assigned to the tokens (e.g., READ, WRITE, DELETE, CREATE)
- **Resources**: The resources for which the tokens are valid (e.g., FAMILY, PARENT, CHILD)
- **Subject**: The subject of the tokens (e.g., admin, editor, viewer)

Example customization:

```
// Example customization
customScopes := []string{"READ", "WRITE"}
customResources := []string{"FAMILY"}
customToken, err := authInstance.GenerateToken(ctx, "custom", []string{"CUSTOM_ROLE"}, customScopes, customResources)
```

## Testing

The JWT Token Generator is tested through:

1. **Manual Testing**: Running the tool and verifying the generated tokens
2. **Validation Testing**: Validating the generated tokens to ensure they contain the expected claims
3. **Integration Testing**: Using the generated tokens in the Family Service application to test authentication and authorization

Key testing approaches:

- **Token Generation**: Testing that tokens are generated with the correct roles, scopes, and resources
- **Token Validation**: Testing that tokens can be validated and claims can be extracted
- **Error Handling**: Testing that appropriate error messages are provided for token generation and validation failures
- **Integration**: Testing that the generated tokens work correctly with the Family Service application

## Design Notes

1. **Simplicity**: The tool is designed to be simple and easy to use
2. **Flexibility**: The tool supports generating tokens with different roles, scopes, and resources
3. **Integration**: The tool integrates with the servicelib/auth package for token generation and validation
4. **Security**: The tool uses a secure secret key for token signing
5. **Extensibility**: The tool can be easily extended to support additional token properties

## Related Components

- [Auth Directive Example](../../EXAMPLES/auth_directive/README.md) - Shows how to use JWT tokens for authentication in GraphQL
- [Security Adapter](../../infrastructure/adapters/security/README.md) - Provides authentication and authorization services
- [GraphQL Resolvers](../../interface/adapters/graphql/resolver/README.md) - Uses JWT tokens for authentication and authorization

## References

- [JWT Standard](https://jwt.io/) - The JSON Web Token standard
- [servicelib/auth](https://github.com/abitofhelp/servicelib/auth) - The authentication library used by this tool
- [Role-Based Access Control](https://en.wikipedia.org/wiki/Role-based_access_control) - The access control model used by this tool
