# Circuit Wrapper

## Overview

The Circuit Wrapper package provides a wrapper around the `github.com/abitofhelp/servicelib/circuit` package to ensure that the domain layer doesn't directly depend on external libraries. This follows the principles of Clean Architecture and Hexagonal Architecture (Ports and Adapters), allowing the domain layer to remain isolated from external dependencies. The package implements the Circuit Breaker pattern to prevent cascading failures when external dependencies are unavailable.

## Architecture

The Circuit Wrapper package follows the Adapter pattern from Hexagonal Architecture, providing a layer of abstraction over the external `servicelib/circuit` package. This ensures that the core domain doesn't directly depend on external libraries, maintaining the dependency inversion principle.

The package sits in the infrastructure layer of the application and is used by the domain layer through interfaces defined in the domain layer. The architecture follows these principles:

- **Dependency Inversion**: The domain layer depends on abstractions, not concrete implementations
- **Adapter Pattern**: This package adapts the external library to the domain's needs
- **Circuit Breaker Pattern**: Prevents cascading failures when external dependencies are unavailable

## Implementation Details

The Circuit Wrapper package implements the following design patterns:

1. **Adapter Pattern**: Adapts the external library to the domain's needs
2. **Circuit Breaker Pattern**: Prevents cascading failures when external dependencies are unavailable
3. **Null Object Pattern**: Handles nil circuit breaker gracefully by falling back to the original function
4. **State Pattern**: Manages the state of the circuit breaker (Closed, Open, HalfOpen)
5. **Facade Pattern**: Simplifies the interface to the underlying circuit breaker implementation

Key implementation details:

- **State Enum**: Represents the state of the circuit breaker (Closed, Open, HalfOpen)
- **CircuitBreaker Struct**: Implements the circuit breaker pattern
- **Configuration Integration**: Uses application configuration to configure the circuit breaker
- **Context Propagation**: Supports context-aware circuit breaking operations
- **Fallback Support**: Provides ExecuteWithFallback for handling failures with fallback logic
- **Graceful Degradation**: Falls back to the original function when circuit breaker is disabled or nil
- **State Management**: Manages the state of the circuit breaker based on the results of operations

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Family Service Example](../../../examples/family_service/README.md) - Shows how to use the circuit wrapper

Example of using the circuit wrapper:

```
// Create a new circuit breaker
cb := circuit.NewCircuitBreaker("database", cfg.Circuit, logger)

// Execute a function with circuit breaking
err := cb.Execute(ctx, "GetByID", func(ctx context.Context) error {
    // This function will be executed if the circuit is closed or half-open
    return repository.GetByID(ctx, id)
})
if err != nil {
    // Handle error (could be a circuit open error or a repository error)
}

// Execute a function with circuit breaking and fallback
err := cb.ExecuteWithFallback(
    ctx,
    "GetByID",
    func(ctx context.Context) error {
        // This function will be executed if the circuit is closed or half-open
        return repository.GetByID(ctx, id)
    },
    func(ctx context.Context, err error) error {
        // This function will be executed if the circuit is open or the main function fails
        return getFromCache(ctx, id)
    },
)
if err != nil {
    // Handle error (could be a fallback error)
}

// Get the current state of the circuit breaker
state := cb.GetState()
switch state {
case circuit.Closed:
    // Circuit is closed, requests are allowed through
case circuit.Open:
    // Circuit is open, requests are not allowed through
case circuit.HalfOpen:
    // Circuit is half-open, limited requests are allowed through
}

// Reset the circuit breaker
cb.Reset()
```

## Configuration

The Circuit Wrapper package is configured through the application's configuration system. The following configuration options are available:

- **Enabled**: Whether the circuit breaker is enabled
- **Timeout**: The maximum time allowed for a function to execute
- **MaxConcurrent**: The maximum number of concurrent requests allowed
- **ErrorThreshold**: The error threshold percentage (0.0 to 1.0) that triggers the circuit to open
- **VolumeThreshold**: The minimum number of requests required before the error threshold is checked
- **SleepWindow**: The time the circuit stays open before transitioning to half-open

Example configuration:

```yaml
circuit:
  enabled: true
  timeout: 5s
  max_concurrent: 100
  error_threshold: 0.5
  volume_threshold: 20
  sleep_window: 10s
```

## Testing

The Circuit Wrapper package is tested through:

1. **Unit Tests**: Each function and method has unit tests
2. **Integration Tests**: Tests that verify the wrapper works correctly with the underlying circuit breaker
3. **State Transition Tests**: Tests that verify the circuit breaker transitions between states correctly

Key testing approaches:

- **Mock Dependencies**: Tests use mock dependencies to verify circuit breaker behavior
- **State Transition Testing**: Tests verify that the circuit breaker transitions between states correctly
- **Fallback Testing**: Tests verify that fallback functions are called when appropriate
- **Error Handling**: Tests verify that errors are properly propagated and transformed
- **Configuration Testing**: Tests verify that configuration options affect circuit breaker behavior

Example of a test case:

```
// Create a circuit breaker
cfg := &config.CircuitConfig{
    Enabled: true,
    Timeout: 1 * time.Second,
    MaxConcurrent: 10,
    ErrorThreshold: 0.5,
    VolumeThreshold: 5,
    SleepWindow: 5 * time.Second,
}
logger, _ := zap.NewDevelopment()
cb := circuit.NewCircuitBreaker("test", cfg, logger)

// Test successful execution
err := cb.Execute(context.Background(), "test", func(ctx context.Context) error {
    return nil
})
assert.NoError(t, err)

// Test failed execution
err = cb.Execute(context.Background(), "test", func(ctx context.Context) error {
    return fmt.Errorf("test error")
})
assert.Error(t, err)
assert.Equal(t, "test error", err.Error())

// Test circuit open
// Force the circuit to open by making multiple failed requests
for i := 0; i < 10; i++ {
    _ = cb.Execute(context.Background(), "test", func(ctx context.Context) error {
        return fmt.Errorf("test error")
    })
}

// Verify the circuit is open
assert.Equal(t, circuit.Open, cb.GetState())

// Test execution with open circuit
err = cb.Execute(context.Background(), "test", func(ctx context.Context) error {
    return nil
})
assert.Error(t, err)
assert.Contains(t, err.Error(), "circuit breaker test is open")
```

## Design Notes

1. **Dependency Inversion**: The package follows the Dependency Inversion Principle by ensuring that the domain layer depends on abstractions rather than concrete implementations
2. **Graceful Degradation**: The package gracefully handles nil circuit breaker by falling back to the original function
3. **Context Propagation**: The package supports context-aware circuit breaking operations
4. **State Management**: The package manages the state of the circuit breaker based on the results of operations
5. **Fallback Support**: The package provides ExecuteWithFallback for handling failures with fallback logic
6. **Configuration Integration**: The package uses application configuration to configure the circuit breaker

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Circuit Breaker Pattern](https://martinfowler.com/bliki/CircuitBreaker.html)
- [Adapter Pattern](https://en.wikipedia.org/wiki/Adapter_pattern)
- [State Pattern](https://en.wikipedia.org/wiki/State_pattern)
- [Application Services](../../../core/application/services/README.md) - Uses this circuit wrapper for circuit breaking
- [Repository Implementations](../../../infrastructure/adapters/repository/README.md) - Uses this circuit wrapper for database operations
- [GraphQL Resolvers](../../../interface/adapters/graphql/resolver/README.md) - Uses this circuit wrapper for GraphQL operations