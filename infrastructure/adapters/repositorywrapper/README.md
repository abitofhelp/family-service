# Infrastructure Adapters - Repository Wrapper

## Overview

The Repository Wrapper adapter provides implementations for repository-related ports defined in the core domain and application layers. This adapter connects the application to repository implementations, following the Ports and Adapters (Hexagonal) architecture pattern. By isolating repository wrapper implementations in adapter classes, the core business logic remains independent of specific repository technologies, making the system more maintainable, testable, and flexible.

## Features

- Repository pattern implementation
- Cross-cutting concerns for repositories (logging, metrics, caching)
- Transaction management
- Repository decorators
- Error handling and translation
- Performance monitoring
- Retry mechanisms
- Circuit breaking for repository operations

## Installation

```bash
go get github.com/abitofhelp/family-service/infrastructure/adapters/repositorywrapper
```

## Configuration

The repository wrapper can be configured according to specific requirements. Here's an example of configuring the repository wrapper:

```
// Pseudocode example - not actual Go code
// This demonstrates how to configure and use a repository wrapper

// 1. Import necessary packages
import repository, config, logging, cache, metrics

// 2. Create a logger
logger = logging.NewLogger()

// 3. Create dependencies
cacheAdapter = cache.NewCacheAdapter(cacheConfig, logger)
metricsAdapter = metrics.NewMetricsAdapter(metricsConfig, logger)

// 4. Configure the repository wrapper
repositoryConfig = {
    cacheEnabled: true,
    cacheTTL: 5 minutes,
    metricsEnabled: true,
    retryEnabled: true,
    maxRetries: 3,
    retryBackoff: 100 milliseconds,
    circuitBreakerEnabled: true,
    circuitBreakerThreshold: 5,
    circuitBreakerTimeout: 30 seconds
}

// 5. Create the base repository
baseRepository = mongo.NewFamilyRepository(mongoAdapter, logger)

// 6. Create the repository wrapper
familyRepository = repository.NewRepositoryWrapper(
    baseRepository,
    repositoryConfig,
    logger,
    cacheAdapter,
    metricsAdapter
)

// 7. Use the repository wrapper
family, err = familyRepository.FindById(context, "family-123")
if err != nil {
    logger.Error("Failed to find family", err)
}
```

## API Documentation

### Core Concepts

The repository wrapper follows these core concepts:

1. **Decorator Pattern**: Wraps repository implementations to add cross-cutting concerns
2. **Dependency Injection**: Receives dependencies through constructor injection
3. **Configuration**: Configured through a central configuration system
4. **Logging**: Uses a consistent logging approach
5. **Error Handling**: Translates repository-specific errors to domain errors

### Key Adapter Functions

```
// Pseudocode example - not actual Go code
// This demonstrates a repository wrapper implementation

// Repository wrapper structure
type RepositoryWrapper {
    repository     // Base repository implementation
    config         // Repository wrapper configuration
    logger         // Logger for logging operations
    contextLogger  // Context-aware logger
    cache          // Cache adapter
    metrics        // Metrics adapter
}

// Constructor for the repository wrapper
function NewRepositoryWrapper(repository, config, logger, cache, metrics) {
    return new RepositoryWrapper {
        repository: repository,
        config: config,
        logger: logger,
        contextLogger: new ContextLogger(logger),
        cache: cache,
        metrics: metrics
    }
}

// Method to find an entity by ID
function RepositoryWrapper.FindById(context, id) {
    // Implementation would include:
    // 1. Logging the operation with context
    // 2. Checking cache if enabled
    // 3. Starting metrics collection
    // 4. Implementing retry logic
    // 5. Checking circuit breaker status
    // 6. Delegating to the base repository
    // 7. Caching the result if enabled
    // 8. Recording metrics
    // 9. Handling errors
    // 10. Returning the result or error
}
```

## Best Practices

1. **Separation of Concerns**: Keep repository wrapper logic separate from domain logic
2. **Interface Segregation**: Define focused repository interfaces in the domain layer
3. **Dependency Injection**: Use constructor injection for adapter dependencies
4. **Error Translation**: Translate repository-specific errors to domain errors
5. **Consistent Logging**: Use a consistent logging approach
6. **Transaction Management**: Implement proper transaction handling
7. **Testing**: Write unit and integration tests for repository wrappers
8. **Performance Monitoring**: Include performance metrics for repository operations

## Troubleshooting

### Common Issues

#### Cache Consistency

If you encounter cache consistency issues, consider the following:
- Implement cache invalidation strategies
- Use appropriate cache TTL values
- Consider cache dependencies for related entities
- Implement cache versioning
- Use write-through or write-behind caching strategies

#### Performance Issues

If you encounter performance issues with repositories, consider the following:
- Optimize the underlying repository implementation
- Adjust cache settings for frequently accessed data
- Review retry and circuit breaker configurations
- Monitor and optimize transaction usage
- Implement batch operations where appropriate

## Related Components

- [Domain Layer](../../core/domain/README.md) - The domain layer that defines the repository interfaces
- [Application Layer](../../core/application/README.md) - The application layer that uses repositories
- [Database Adapters](../mongo/README.md) - The database adapters used by repositories

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.