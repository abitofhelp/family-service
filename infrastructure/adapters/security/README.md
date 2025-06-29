# Infrastructure Adapters - Security

## Overview

The Security adapter provides implementations for security-related ports defined in the core domain and application layers. This adapter connects the application to security frameworks and libraries, following the Ports and Adapters (Hexagonal) architecture pattern. By isolating security implementations in adapter classes, the core business logic remains independent of specific security technologies, making the system more maintainable, testable, and flexible.

## Features

- Authentication mechanisms (JWT, OAuth, API keys, etc.)
- Authorization and access control
- Encryption and decryption utilities
- Hashing and password management
- Secure token generation and validation
- Security headers management
- CSRF protection
- Rate limiting for security purposes

## Installation

```bash
go get github.com/abitofhelp/family-service/infrastructure/adapters/security
```

## Configuration

The security adapter can be configured according to specific requirements. Here's an example of configuring the security adapter:

```
// Pseudocode example - not actual Go code
// This demonstrates how to configure and use a security adapter

// 1. Import necessary packages
import security, config, logging

// 2. Create a logger
logger = logging.NewLogger()

// 3. Configure the security adapter
securityConfig = {
    jwtSecret: "your-jwt-secret",
    jwtExpirationMinutes: 60,
    passwordHashingStrength: 12,
    enableCSRF: true,
    csrfTokenExpiration: 30 minutes,
    securityHeaders: {
        "Content-Security-Policy": "default-src 'self'",
        "X-XSS-Protection": "1; mode=block",
        "X-Frame-Options": "DENY"
    }
}

// 4. Create the security adapter
securityAdapter = security.NewSecurityAdapter(securityConfig, logger)

// 5. Use the security adapter
token, err = securityAdapter.GenerateJWT(user)
if err != nil {
    logger.Error("Failed to generate JWT", err)
}

isValid, claims = securityAdapter.ValidateJWT(token)
if !isValid {
    logger.Warn("Invalid JWT token")
}
```

## API Documentation

### Core Concepts

The security adapter follows these core concepts:

1. **Adapter Pattern**: Implements security ports defined in the core domain or application layer
2. **Dependency Injection**: Receives dependencies through constructor injection
3. **Configuration**: Configured through a central configuration system
4. **Logging**: Uses a consistent logging approach
5. **Error Handling**: Translates security-specific errors to domain errors

### Key Adapter Functions

```
// Pseudocode example - not actual Go code
// This demonstrates a security adapter implementation

// Security adapter structure
type SecurityAdapter {
    config        // Security configuration
    logger        // Logger for logging operations
    contextLogger // Context-aware logger
}

// Constructor for the security adapter
function NewSecurityAdapter(config, logger) {
    return new SecurityAdapter {
        config: config,
        logger: logger,
        contextLogger: new ContextLogger(logger)
    }
}

// Method to generate a JWT token
function SecurityAdapter.GenerateJWT(user) {
    // Implementation would include:
    // 1. Logging the operation
    // 2. Creating claims based on user information
    // 3. Generating a signed JWT token
    // 4. Handling errors
    // 5. Returning the token or error
}

// Method to validate a JWT token
function SecurityAdapter.ValidateJWT(token) {
    // Implementation would include:
    // 1. Logging the operation
    // 2. Parsing and validating the JWT token
    // 3. Extracting claims from the token
    // 4. Handling validation errors
    // 5. Returning validation status and claims
}
```

## Best Practices

1. **Separation of Concerns**: Keep security logic separate from domain logic
2. **Interface Segregation**: Define focused security interfaces in the domain layer
3. **Dependency Injection**: Use constructor injection for adapter dependencies
4. **Error Translation**: Translate security-specific errors to domain errors
5. **Consistent Logging**: Use a consistent logging approach
6. **Configuration**: Configure security through a central configuration system
7. **Testing**: Write unit and integration tests for security adapters
8. **Secrets Management**: Never hardcode secrets in the code or configuration files

## Troubleshooting

### Common Issues

#### Authentication Failures

If you encounter authentication issues, check the following:
- JWT secret is correctly configured
- Token expiration times are appropriate
- Clock synchronization between systems
- Token validation logic is correct
- User credentials are valid

#### Security Configuration

If you encounter security configuration issues, consider the following:
- Review security headers for correctness
- Ensure CSRF protection is properly configured
- Verify password hashing strength is appropriate
- Check that encryption keys are properly managed
- Ensure secure communication channels are used

## Related Components

- [Domain Layer](../../core/domain/README.md) - The domain layer that defines the security ports
- [Application Layer](../../core/application/README.md) - The application layer that uses security
- [Interface Adapters](../../interface/adapters/README.md) - The interface adapters that use security

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.