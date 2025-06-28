# Circuit Breaker Guide for Family Service

## Overview

This document describes the circuit breaker pattern used in the Family Service, including its implementation, configuration, and error handling.

## Circuit Breaker Pattern

The circuit breaker pattern is a design pattern used to detect failures and prevent cascading failures in distributed systems. It works by "tripping" when a certain threshold of failures is reached, preventing further requests to the failing service until it has recovered.

## Implementation

The Family Service implements the circuit breaker pattern in the `infrastructure/adapters/circuit` package, which provides a wrapper around the ServiceLib circuit breaker:

```
// CircuitBreaker implements the circuit breaker pattern to protect against
// cascading failures when external dependencies are unavailable.
type CircuitBreaker struct {
    name   string
    cb     *circuit.CircuitBreaker
    logger *zap.Logger
}
```

The circuit breaker can be in one of three states:
- **Closed**: Normal operation, requests are allowed through
- **Open**: Circuit is tripped, requests are immediately rejected
- **HalfOpen**: Testing if the dependency has recovered, allowing a limited number of requests through

## Configuration

The circuit breaker is configured through the application configuration:

```
circuit:
  enabled: true              # Whether the circuit breaker is enabled
  timeout: 1s                # Maximum time allowed for a request
  max_concurrent: 100        # Maximum number of concurrent requests
  error_threshold: 0.5       # Error rate threshold (0.0-1.0) that trips the circuit
  volume_threshold: 10       # Minimum number of requests before error threshold is considered
  sleep_window: 5s           # Time to wait before allowing requests when circuit is open
```

The configuration can be overridden using environment variables:

```
APP_CIRCUIT_ENABLED=true
APP_CIRCUIT_TIMEOUT=1s
APP_CIRCUIT_MAX_CONCURRENT=100
APP_CIRCUIT_ERROR_THRESHOLD=0.5
APP_CIRCUIT_VOLUME_THRESHOLD=10
APP_CIRCUIT_SLEEP_WINDOW=5s
```

## Usage

The circuit breaker is used in the MongoDB and SQLite repository implementations to protect against database failures:

```
// Execute with circuit breaker
_, err := circuit.Execute(ctx, r.circuitBreaker, "GetByID", circuitOpWrapper)
```

Each repository operation is wrapped with the circuit breaker, which will:
1. Check if the circuit is open (if so, immediately return an error)
2. Execute the operation if the circuit is closed or half-open
3. Update circuit state based on the result (success or failure)
4. Return the result or an error

## Error Handling

When the circuit is open, operations return a specific error that is converted to a database error:

```
// Check for errors from circuit breaker
if err == recovery.ErrCircuitBreakerOpen {
    return nil, errors.NewDatabaseError("circuit breaker is open", "query", "families", err)
}
```

This allows the application to distinguish between actual database errors and circuit breaker rejections.

## Fallback Mechanism

The circuit breaker implementation also supports a fallback mechanism, which allows the application to provide alternative behavior when the circuit is open:

```
// Execute with fallback
err := cb.ExecuteWithFallback(ctx, "operation", func(ctx context.Context) error {
    // Primary operation
    return primaryOperation(ctx)
}, func(ctx context.Context, err error) error {
    // Fallback operation
    return fallbackOperation(ctx)
})
```

This is useful for operations that can be gracefully degraded, such as returning cached data instead of querying the database.

## Best Practices

1. **Use Circuit Breakers for External Dependencies**: Apply circuit breakers to all external dependencies, such as databases, APIs, and other services.
2. **Configure Appropriate Thresholds**: Set error thresholds and volume thresholds based on the expected behavior of the dependency.
3. **Implement Fallbacks**: Provide fallback mechanisms for critical operations to maintain service availability.
4. **Monitor Circuit Breaker State**: Log and monitor circuit breaker state changes to detect issues with dependencies.
5. **Test Circuit Breaker Behavior**: Test how your application behaves when the circuit breaker is open, closed, and half-open.

## Examples

For examples of circuit breaker usage, see the repository implementations in the `infrastructure/adapters/mongo/repo.go` and `infrastructure/adapters/sqlite/repo.go` files.