# Infrastructure Adapters - Dependency Injection Wrapper

## Overview

The Dependency Injection Wrapper adapter provides implementations for dependency injection-related ports defined in the core domain and application layers. This adapter connects the application to dependency injection frameworks and libraries, following the Ports and Adapters (Hexagonal) architecture pattern. By isolating dependency injection implementations in adapter classes, the core business logic remains independent of specific DI technologies, making the system more maintainable, testable, and flexible.

## Features

- Dependency registration and resolution
- Lifecycle management of dependencies
- Scoped dependency containers
- Factory pattern support
- Lazy initialization
- Conditional registration
- Configuration-based dependency setup
- Integration with various DI frameworks

## Installation

```bash
go get github.com/abitofhelp/family-service/infrastructure/adapters/diwrapper
```

## Configuration

The dependency injection wrapper can be configured according to specific requirements. Here's an example of configuring the DI wrapper:

```
// Pseudocode example - not actual Go code
// This demonstrates how to configure and use a dependency injection wrapper

// 1. Import necessary packages
import di, config, logging

// 2. Create a logger
logger = logging.NewLogger()

// 3. Create the DI container
container = di.NewContainer(logger)

// 4. Register dependencies
container.Register("logger", logger, {singleton: true})
container.Register("config", configAdapter, {singleton: true})
container.Register("database", databaseAdapter, {singleton: true})
container.Register("familyRepository", function(c) {
    return repository.NewFamilyRepository(
        c.Resolve("database"),
        c.Resolve("logger")
    )
}, {singleton: true})

// 5. Use the container to resolve dependencies
familyRepository = container.Resolve("familyRepository")
```

## API Documentation

### Core Concepts

The dependency injection wrapper follows these core concepts:

1. **Adapter Pattern**: Implements dependency injection ports defined in the core domain or application layer
2. **Dependency Injection**: Facilitates the injection of dependencies into components
3. **Configuration**: Configured through a central configuration system
4. **Logging**: Uses a consistent logging approach
5. **Error Handling**: Handles dependency resolution errors gracefully

### Key Adapter Functions

```
// Pseudocode example - not actual Go code
// This demonstrates a dependency injection wrapper implementation

// DI container structure
type Container {
    logger        // Logger for logging operations
    contextLogger // Context-aware logger
    registrations // Map of registered dependencies
}

// Constructor for the DI container
function NewContainer(logger) {
    return new Container {
        logger: logger,
        contextLogger: new ContextLogger(logger),
        registrations: {}
    }
}

// Method to register a dependency
function Container.Register(name, instance, options) {
    // Implementation would include:
    // 1. Logging the operation
    // 2. Validating the registration
    // 3. Storing the registration with options
    // 4. Handling registration errors
}

// Method to resolve a dependency
function Container.Resolve(name) {
    // Implementation would include:
    // 1. Logging the operation
    // 2. Looking up the registration
    // 3. Creating or returning the instance based on lifecycle
    // 4. Handling resolution errors
    // 5. Returning the resolved instance
}
```

## Best Practices

1. **Separation of Concerns**: Keep dependency injection logic separate from domain logic
2. **Interface Segregation**: Define focused interfaces for components
3. **Constructor Injection**: Use constructor injection for dependencies
4. **Lifecycle Management**: Properly manage the lifecycle of dependencies
5. **Consistent Logging**: Use a consistent logging approach
6. **Configuration**: Configure dependency injection through a central configuration system
7. **Testing**: Write unit and integration tests for dependency injection

## Troubleshooting

### Common Issues

#### Circular Dependencies

If you encounter circular dependency issues, consider the following:
- Refactor components to break circular dependencies
- Use lazy initialization for one of the dependencies
- Introduce an interface to break the cycle
- Use a mediator pattern to decouple components

#### Resolution Failures

If you encounter dependency resolution failures, check the following:
- Ensure all dependencies are registered before resolution
- Verify that dependency names are correct
- Check that factory functions don't have errors
- Ensure that required dependencies for a component are available
- Look for typos in dependency names

## Related Components

- [Domain Layer](../../core/domain/README.md) - The domain layer that defines the ports
- [Application Layer](../../core/application/README.md) - The application layer that uses dependency injection
- [Interface Adapters](../../interface/adapters/README.md) - The interface adapters that use dependency injection

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.