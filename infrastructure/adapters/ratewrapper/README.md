# Infrastructure Adapters - Rate Wrapper

## Overview

The Rate Wrapper adapter provides implementations for rate limiting and throttling ports defined in the core domain and application layers. This adapter connects the application to rate limiting frameworks and libraries, following the Ports and Adapters (Hexagonal) architecture pattern. By isolating rate limiting implementations in adapter classes, the core business logic remains independent of specific rate limiting technologies, making the system more maintainable, testable, and flexible.

## Features

- Request rate limiting
- Throttling of API calls
- Configurable rate limits
- Token bucket algorithm implementation
- Rate limit headers generation
- Distributed rate limiting support
- Custom rate limiting strategies

## Installation

```bash
go get github.com/abitofhelp/family-service/infrastructure/adapters/ratewrapper
```

## Configuration

The rate wrapper adapter can be configured according to specific requirements. Here's an example of configuring the rate wrapper adapter:

```
// Pseudocode example - not actual Go code
// This demonstrates how to configure and use a rate wrapper adapter

// 1. Import necessary packages
import rate, config, logging, time

// 2. Create a logger
logger = logging.NewLogger()

// 3. Configure the rate limiter
rateConfig = {
    enabled: true,
    requestsPerSecond: 100,
    burstSize: 150,
    strategy: "token-bucket",
    distributedCache: "redis://localhost:6379"
}

// 4. Create the rate wrapper adapter
rateAdapter = rate.NewRateLimiter(rateConfig, logger)

// 5. Use the rate wrapper adapter
allowed, remaining, reset = rateAdapter.Allow(context, "user-123", "api-endpoint")
if !allowed {
    logger.Warn("Rate limit exceeded", {user: "user-123", remaining: 0, resetAt: reset})
    return RateLimitExceededError
}
```

## API Documentation

### Core Concepts

The rate wrapper adapter follows these core concepts:

1. **Adapter Pattern**: Implements rate limiting ports defined in the core domain or application layer
2. **Dependency Injection**: Receives dependencies through constructor injection
3. **Configuration**: Configured through a central configuration system
4. **Logging**: Uses a consistent logging approach
5. **Error Handling**: Translates rate limiting errors to domain errors

### Key Adapter Functions

```
// Pseudocode example - not actual Go code
// This demonstrates a rate wrapper adapter implementation

// Rate wrapper adapter structure
type RateLimiter {
    config        // Rate limiter configuration
    logger        // Logger for logging operations
    contextLogger // Context-aware logger
    store         // Storage for rate limit state
}

// Constructor for the rate wrapper adapter
function NewRateLimiter(config, logger) {
    store = createRateLimitStore(config)
    return new RateLimiter {
        config: config,
        logger: logger,
        contextLogger: new ContextLogger(logger),
        store: store
    }
}

// Method to check if a request is allowed
function RateLimiter.Allow(context, key, action) {
    // Implementation would include:
    // 1. Logging the operation
    // 2. Checking if the request is allowed based on rate limits
    // 3. Updating rate limit counters
    // 4. Returning allow status, remaining requests, and reset time
}
```

## Best Practices

1. **Separation of Concerns**: Keep rate limiting logic separate from domain logic
2. **Interface Segregation**: Define focused rate limiting interfaces in the domain layer
3. **Dependency Injection**: Use constructor injection for adapter dependencies
4. **Error Translation**: Translate rate limiting errors to domain errors
5. **Consistent Logging**: Use a consistent logging approach
6. **Configuration**: Configure rate limiting through a central configuration system
7. **Testing**: Write unit and integration tests for rate limiting adapters

## Troubleshooting

### Common Issues

#### Rate Limit Configuration

If you encounter issues with rate limiting configuration, check the following:
- Rate limit values are appropriate for your application's needs
- Burst size is configured correctly
- Rate limit strategy is appropriate for your use case
- Distributed cache configuration is correct (if using distributed rate limiting)

#### Performance Issues

If rate limiting is impacting application performance, consider the following:
- Use a more efficient rate limiting algorithm
- Implement caching for rate limit state
- Use a distributed cache for rate limit state in a clustered environment
- Optimize the rate limit check process

## Related Components

- [Domain Layer](../../core/domain/README.md) - The domain layer that defines the rate limiting ports
- [Application Layer](../../core/application/README.md) - The application layer that uses rate limiting
- [Interface Adapters](../../interface/adapters/README.md) - The interface adapters that use rate limiting

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.