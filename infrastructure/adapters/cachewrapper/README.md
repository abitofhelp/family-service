# Cache Wrapper

## Overview

The Cache Wrapper package provides a wrapper around the `github.com/abitofhelp/servicelib/cache` package to ensure that the domain layer doesn't directly depend on external libraries. This follows the principles of Clean Architecture and Hexagonal Architecture (Ports and Adapters), allowing the domain layer to remain isolated from external dependencies.

## Architecture

The Cache Wrapper package follows the Adapter pattern from Hexagonal Architecture, providing a layer of abstraction over the external `servicelib/cache` package. This ensures that the core domain doesn't directly depend on external libraries, maintaining the dependency inversion principle.

The package sits in the infrastructure layer of the application and is used by the domain layer through interfaces defined in the domain layer. The architecture follows these principles:

- **Dependency Inversion**: The domain layer depends on abstractions, not concrete implementations
- **Adapter Pattern**: This package adapts the external library to the domain's needs
- **Middleware Pattern**: Provides middleware functions for adding caching to application functions

## Implementation Details

The Cache Wrapper package implements the following design patterns:

1. **Adapter Pattern**: Adapts the external library to the domain's needs
2. **Middleware Pattern**: Provides middleware functions for adding caching to application functions
3. **Null Object Pattern**: Handles nil cache gracefully by falling back to the original function
4. **Facade Pattern**: Simplifies the interface to the underlying cache implementation

Key implementation details:

- **Cache Struct**: Wraps the servicelib's cache to provide a consistent interface
- **Configuration Integration**: Uses application configuration to configure the cache
- **Context Propagation**: Supports context-aware caching operations
- **Middleware Functions**: Provides WithCache and WithContextCache for adding caching to functions
- **Graceful Degradation**: Falls back to the original function when cache is disabled or nil
- **Resource Management**: Provides a Shutdown method to stop the cleanup timer

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Family Service Example](../../../examples/family_service/README.md) - Shows how to use the cache wrapper

Example of using the cache wrapper:

```
// Create a new cache
cache, err := cache.NewCache(cfg, logger)
if err != nil {
    // Handle error
}
defer cache.Shutdown()

// Set a value in the cache
cache.Set("key", value)

// Get a value from the cache
value, found := cache.Get("key")
if found {
    // Use the cached value
}

// Delete a value from the cache
cache.Delete("key")

// Use the cache middleware
result, err := cache.WithCache(cache, "key", func() (interface{}, error) {
    // This function will only be called if the key is not in the cache
    return expensiveOperation()
})

// Use the context-aware cache middleware
result, err := cache.WithContextCache(ctx, cache, "key", func(ctx context.Context) (interface{}, error) {
    // This function will only be called if the key is not in the cache
    return expensiveOperationWithContext(ctx)
})
```

## Configuration

The Cache Wrapper package is configured through the application's configuration system. The following configuration options are available:

- **Enabled**: Whether the cache is enabled
- **TTL**: The time-to-live for cache entries
- **MaxSize**: The maximum number of items in the cache
- **PurgeInterval**: The interval at which expired items are purged from the cache

Example configuration:

```yaml
cache:
  enabled: true
  ttl: 5m
  max_size: 1000
  purge_interval: 10m
```

## Testing

The Cache Wrapper package is tested through:

1. **Unit Tests**: Each function and method has unit tests
2. **Integration Tests**: Tests that verify the wrapper works correctly with the underlying cache
3. **Middleware Tests**: Tests that verify the middleware functions work correctly

Key testing approaches:

- **Mock Cache**: Tests use a mock cache to verify caching behavior
- **Expiration Testing**: Tests verify that cache entries expire correctly
- **Middleware Testing**: Tests verify that middleware functions correctly cache function results
- **Error Handling**: Tests verify that errors are properly propagated

Example of a test case:

```
// Create a cache
cfg := &config.Config{
    Cache: config.CacheConfig{
        Enabled: true,
        TTL: 1 * time.Minute,
        MaxSize: 100,
        PurgeInterval: 5 * time.Minute,
    },
}
logger, _ := zap.NewDevelopment()
cache, err := cache.NewCache(cfg, logger)
assert.NoError(t, err)

// Set a value
cache.Set("key", "value")

// Get the value
value, found := cache.Get("key")
assert.True(t, found)
assert.Equal(t, "value", value)

// Get a non-existent value
value, found = cache.Get("nonexistent")
assert.False(t, found)
assert.Nil(t, value)
```

## Design Notes

1. **Dependency Inversion**: The package follows the Dependency Inversion Principle by ensuring that the domain layer depends on abstractions rather than concrete implementations
2. **Graceful Degradation**: The package gracefully handles nil cache by falling back to the original function
3. **Context Propagation**: The package supports context-aware caching operations
4. **Resource Management**: The package provides a Shutdown method to stop the cleanup timer
5. **Middleware Pattern**: The package provides middleware functions for adding caching to application functions
6. **Configuration Integration**: The package uses application configuration to configure the cache

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Adapter Pattern](https://en.wikipedia.org/wiki/Adapter_pattern)
- [Middleware Pattern](https://en.wikipedia.org/wiki/Middleware)
- [Null Object Pattern](https://en.wikipedia.org/wiki/Null_object_pattern)
- [Application Services](../../../core/application/services/README.md) - Uses this cache wrapper for caching
- [Domain Services](../../../core/domain/services/README.md) - Uses this cache wrapper for caching
- [GraphQL Resolvers](../../../interface/adapters/graphql/resolver/README.md) - Uses this cache wrapper for caching