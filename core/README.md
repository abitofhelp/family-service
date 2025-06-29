# Core

## Overview

This directory contains the core business logic of the Family Service, organized according to Domain-Driven Design (DDD) principles. It is divided into two main subdirectories: `domain` and `application`.

## Features

- **Domain Layer**: Contains the business entities, value objects, domain services, and domain interfaces
- **Application Layer**: Contains application services that orchestrate domain operations
- **Clean Architecture**: Follows clean architecture principles with dependencies pointing inward
- **Domain-Driven Design**: Implements DDD concepts like entities, value objects, and aggregates
- **Hexagonal Architecture**: Uses ports and adapters pattern for flexible infrastructure integration

## Installation

This code is part of the Family Service and does not require separate installation.

## Quick Start

The core business logic is used by the service's interface layer. To understand how it works, explore the domain and application layers:

```
core/
├── application/  # Application services and use cases
│   ├── ports/    # Application ports (interfaces)
│   └── services/ # Application services
└── domain/       # Domain model and business rules
    ├── entity/   # Domain entities
    ├── errors/   # Domain-specific errors
    ├── metrics/  # Domain metrics
    ├── ports/    # Domain ports (interfaces)
    ├── services/ # Domain services
    └── validation/ # Domain validation rules
```

## API Documentation

### Core Types

#### Domain Layer

The domain layer contains the business entities and logic of the Family Service.

#### Application Layer

The application layer orchestrates domain operations and provides use cases for the interface layer.

## Examples

There may be additional examples in the /EXAMPLES directory.

## Best Practices

1. **Dependency Rule**: Dependencies always point inward (domain ← application ← interface ← infrastructure)
2. **Domain Purity**: Keep the domain layer pure and free from infrastructure concerns
3. **Use Cases**: Implement use cases as application services
4. **Domain Events**: Use domain events for cross-aggregate communication
5. **Validation**: Validate inputs at the application layer and business rules at the domain layer

## Troubleshooting

### Common Issues

#### Circular Dependencies

If you encounter circular dependencies, review your architecture to ensure dependencies point inward.

#### Business Logic Leakage

If business logic leaks into the interface or infrastructure layers, refactor to move it to the domain layer.

## Related Components

- [Interface Layer](../interface/README.md) - Adapters that expose the core functionality
- [Infrastructure Layer](../infrastructure/README.md) - External services and technical concerns
- [Command Layer](../cmd/README.md) - Application entry points

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.