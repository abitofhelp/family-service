# Cache Wrapper

## Overview

The Cache Wrapper package provides a wrapper around the `github.com/abitofhelp/servicelib/cache` package to ensure that the domain layer doesn't directly depend on external libraries. This follows the principles of Clean Architecture and Hexagonal Architecture (Ports and Adapters), allowing the domain layer to remain isolated from external dependencies.

> **For Junior Developers**: Think of this adapter as a middleman between your business logic and the caching library. It allows your core business code to use caching without knowing the details of how the specific cache library works.

## Getting Started

If you're new to this codebase, follow these steps to start using the Cache Wrapper:

1. **Understand the purpose**: The Cache Wrapper helps you store and retrieve data quickly without having to recalculate or refetch it
2. **Learn the interfaces**: Look at the domain ports to understand what caching operations are available
3. **Ask questions**: If something isn't clear, ask a more experienced developer

## Architecture

The Cache Wrapper package follows the Adapter pattern from Hexagonal Architecture, providing a layer of abstraction over the external `servicelib/cache` package. This ensures that the core domain doesn't directly depend on external libraries, maintaining the dependency inversion principle.

The package sits in the infrastructure layer of the application and is used by the domain layer through interfaces defined in the domain layer. The architecture follows these principles:

- **Dependency Inversion**: The domain layer depends on abstractions, not concrete implementations
- **Adapter Pattern**: This package adapts the external library to the domain's needs
- **Middleware Pattern**: Provides middleware functions for adding caching to application functions

## API Documentation

### Core Concepts

> **For Junior Developers**: These concepts are fundamental to understanding how the Cache Wrapper works. Take time to understand each one before diving into the code.

The Cache Wrapper follows these core concepts:

1. **Adapter Pattern**: Implements caching ports defined in the core domain or application layer
   - This means the Cache Wrapper implements interfaces defined elsewhere
   - The business logic only knows about these interfaces, not the implementation details

2. **Dependency Injection**: Receives dependencies through constructor injection
   - Dependencies like loggers are passed in when creating the wrapper
   - This makes testing easier and components more loosely coupled

3. **Configuration**: Configured through a central configuration system
   - Settings like TTL and cache size are defined in configuration
   - This allows changing behavior without changing code

4. **Middleware Pattern**: Provides functions that can be wrapped around other functions
   - This allows adding caching to any function with minimal code changes
   - The middleware handles checking the cache and storing results automatically

5. **Null Object Pattern**: Handles nil cache gracefully
   - If the cache is disabled or nil, the original function is called directly
   - This prevents errors when the cache is not available

### Key Adapter Functions

Here are the main functions you'll use when working with the Cache Wrapper:

```
// Pseudocode example - not actual Go code
// This demonstrates a Cache Wrapper implementation

// Cache wrapper structure
type CacheWrapper {
    config        // Cache configuration
    logger        // Logger for logging operations
    contextLogger // Context-aware logger
    cache         // The underlying cache implementation
}

// Constructor for the cache wrapper
// This is how you create a new instance of the wrapper
function NewCacheWrapper(config, logger) {
    // Implementation would include:
    // 1. Validating the configuration
    // 2. Creating the underlying cache
    // 3. Setting up the cleanup timer
    // 4. Returning the wrapper
    return new CacheWrapper {
        config: config,
        logger: logger,
        contextLogger: new ContextLogger(logger),
        cache: createCache(config)
    }
}

// Method to get a value from the cache
// Use this to retrieve cached values
function CacheWrapper.Get(key) {
    // Implementation would include:
    // 1. Logging the operation
    // 2. Getting the value from the cache
    // 3. Returning the value and whether it was found
}

// Method to set a value in the cache
// Use this to store values in the cache
function CacheWrapper.Set(key, value) {
    // Implementation would include:
    // 1. Logging the operation
    // 2. Setting the value in the cache with the configured TTL
}

// Method to add caching to a function
// This is the middleware pattern in action
function CacheWrapper.WithCache(cache, key, fn) {
    // Implementation would include:
    // 1. Checking if the cache is enabled
    // 2. If disabled, just call the original function
    // 3. If enabled, try to get the value from the cache
    // 4. If found, return it
    // 5. If not found, call the original function
    // 6. Store the result in the cache
    // 7. Return the result
}
```

### Common Cache Operations

Here are some common operations you might need to perform:

1. **Creating a new cache**:
   ```
   cache, err := cache.NewCache(cfg, logger)
   if err != nil {
       // Handle error
   }
   defer cache.Shutdown()  // Don't forget to shut down the cache when done
   ```

2. **Setting a value in the cache**:
   ```
   cache.Set("user-123", userObject)  // Store a user object with key "user-123"
   ```

3. **Getting a value from the cache**:
   ```
   value, found := cache.Get("user-123")
   if found {
       user := value.(*User)  // Type assertion to convert to your type
       // Use the cached user
   } else {
       // Value not in cache, need to fetch it
   }
   ```

4. **Deleting a value from the cache**:
   ```
   cache.Delete("user-123")  // Remove the user from cache
   ```

5. **Using the cache middleware for a function**:
   ```
   // This will automatically cache the result of getUserById
   user, err := cache.WithCache(cache, "user-123", func() (interface{}, error) {
       return getUserById("123")  // Only called if not in cache
   })
   ```

6. **Using the context-aware cache middleware**:
   ```
   // This passes the context to the function and handles caching
   user, err := cache.WithContextCache(ctx, cache, "user-123", func(ctx context.Context) (interface{}, error) {
       return getUserByIdWithContext(ctx, "123")
   })
   ```

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

## Best Practices

> **For Junior Developers**: Following these best practices will help you avoid common pitfalls and write more maintainable code.

1. **Choose Appropriate Cache Keys**: Use meaningful and unique cache keys
   - **Why?** Poorly chosen keys can lead to cache collisions or difficulty debugging
   - **Example:** Use `"user-123"` instead of just `"123"` to provide context

2. **Set Appropriate TTL Values**: Choose time-to-live values based on data volatility
   - **Why?** Too short TTLs reduce cache effectiveness; too long TTLs can serve stale data
   - **Example:** User preferences might have a longer TTL than real-time stock prices

3. **Handle Cache Misses Gracefully**: Always have a fallback for cache misses
   - **Why?** Your application should work even if the cache is empty or disabled
   - **Example:** Use the middleware pattern to automatically handle cache misses

4. **Use Type Assertions Carefully**: Be careful when converting cached values back to their original types
   - **Why?** Type assertions can panic if the cached value is not of the expected type
   - **Example:** Use type assertions with a check: `user, ok := value.(*User)`

5. **Don't Cache Everything**: Be selective about what you cache
   - **Why?** Caching everything can waste memory and might not improve performance
   - **Example:** Cache expensive database queries but not simple in-memory operations

6. **Implement Cache Invalidation**: Have a strategy for invalidating cache entries when data changes
   - **Why?** Without invalidation, users might see outdated data
   - **Example:** Delete cache entries when the underlying data is updated

7. **Monitor Cache Performance**: Track cache hit rates and memory usage
   - **Why?** This helps you optimize your caching strategy
   - **Example:** Log cache hits and misses to calculate hit rate

8. **Shutdown the Cache Properly**: Always call Shutdown() when you're done with the cache
   - **Why?** This ensures resources are released properly
   - **Example:** Use `defer cache.Shutdown()` after creating the cache

## Common Mistakes to Avoid

1. **Caching mutable objects without copying**
   - **Problem:** If you cache a reference to a mutable object, changes to the object will affect the cached value
   - **Solution:** Cache immutable objects or deep copies of mutable objects

2. **Using non-serializable objects as cache keys**
   - **Problem:** Complex objects might not work well as cache keys
   - **Solution:** Use strings or simple types for cache keys

3. **Not handling cache errors**
   - **Problem:** Cache operations can fail, especially in distributed systems
   - **Solution:** Always check for and handle errors from cache operations

4. **Caching sensitive data**
   - **Problem:** Sensitive data in cache might be accessible to unauthorized users
   - **Solution:** Don't cache sensitive data, or encrypt it before caching

5. **Forgetting to call Shutdown()**
   - **Problem:** Not shutting down the cache can lead to resource leaks
   - **Solution:** Always use `defer cache.Shutdown()` after creating the cache

## Design Notes

1. **Dependency Inversion**: The package follows the Dependency Inversion Principle by ensuring that the domain layer depends on abstractions rather than concrete implementations
2. **Graceful Degradation**: The package gracefully handles nil cache by falling back to the original function
3. **Context Propagation**: The package supports context-aware caching operations
4. **Resource Management**: The package provides a Shutdown method to stop the cleanup timer
5. **Middleware Pattern**: The package provides middleware functions for adding caching to application functions
6. **Configuration Integration**: The package uses application configuration to configure the cache

## Troubleshooting

### Common Issues

#### Cache Not Working

If your cache doesn't seem to be working, check the following:

- **Is the cache enabled in configuration?**
  - **Problem:** The cache might be disabled in your configuration
  - **Solution:** Check your configuration and ensure `cache.enabled` is set to `true`

- **Are you using the correct cache key?**
  - **Problem:** Cache misses might occur if keys are inconsistent
  - **Solution:** Use consistent key naming conventions and verify the keys being used

- **Is the TTL too short?**
  - **Problem:** Items might expire before you try to access them again
  - **Solution:** Increase the TTL in your configuration or when setting items

- **Is the cache being properly initialized?**
  - **Problem:** The cache might not be properly created or might be nil
  - **Solution:** Check for errors when creating the cache and ensure it's properly initialized

#### Memory Issues

If you're experiencing memory issues with your cache:

- **Is your cache size too large?**
  - **Problem:** Setting a very large max size can consume too much memory
  - **Solution:** Reduce the `maxSize` in your configuration

- **Are you caching very large objects?**
  - **Problem:** Caching large objects can quickly fill up the cache
  - **Solution:** Consider caching only the most important parts of large objects

- **Is your purge interval appropriate?**
  - **Problem:** If expired items aren't purged frequently enough, they can waste memory
  - **Solution:** Adjust the `purgeInterval` in your configuration

#### Type Assertion Errors

If you're getting type assertion panics:

- **Are you using the correct types?**
  - **Problem:** The cached value might not be of the expected type
  - **Solution:** Use type assertions with a check: `value, ok := cachedValue.(*YourType)`

- **Has the cached value been modified?**
  - **Problem:** If the cached value is a reference type and has been modified, it might not match the expected type
  - **Solution:** Cache immutable objects or deep copies of mutable objects

## Related Components

> **For Junior Developers**: Understanding how components relate to each other is crucial for working effectively in this codebase.

- [Domain Layer](../../core/domain/README.md) - The domain layer that defines the caching ports
  - This is where the interfaces that the Cache Wrapper implements are defined
  - Look here to understand what operations are available

- [Application Layer](../../core/application/README.md) - The application layer that uses caching
  - This layer contains the business logic that uses the Cache Wrapper
  - See how caching is used in business processes

- [Application Services](../../../core/application/services/README.md) - Uses this cache wrapper for caching
  - These services implement business logic and use caching to improve performance
  - Look here to see examples of how caching is used in services

- [Domain Services](../../../core/domain/services/README.md) - Uses this cache wrapper for caching
  - These services implement domain logic and use caching to improve performance
  - See how caching is used in domain operations

- [GraphQL Resolvers](../../../interface/adapters/graphql/resolver/README.md) - Uses this cache wrapper for caching
  - These resolvers handle GraphQL queries and use caching to improve response times
  - Look here to see how caching is used in API responses

## Glossary of Terms

- **Adapter Pattern**: A design pattern that allows incompatible interfaces to work together
- **Cache**: A component that stores data so future requests for that data can be served faster
- **TTL (Time-To-Live)**: The duration for which a cached item is considered valid
- **Cache Key**: A unique identifier used to store and retrieve items in the cache
- **Cache Hit**: When a requested item is found in the cache
- **Cache Miss**: When a requested item is not found in the cache
- **Middleware**: Software that acts as a bridge between different components
- **Dependency Injection**: A technique where an object receives its dependencies from outside
- **Graceful Degradation**: The ability of a system to continue functioning when parts of it fail

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Adapter Pattern](https://en.wikipedia.org/wiki/Adapter_pattern)
- [Middleware Pattern](https://en.wikipedia.org/wiki/Middleware)
- [Null Object Pattern](https://en.wikipedia.org/wiki/Null_object_pattern)
