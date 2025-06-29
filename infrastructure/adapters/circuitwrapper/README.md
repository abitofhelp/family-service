# Infrastructure Adapters - Circuit Wrapper

## Overview

The Circuit Wrapper adapter provides implementations for circuit breaker-related ports defined in the core domain and application layers. This adapter connects the application to circuit breaker frameworks and libraries, following the Ports and Adapters (Hexagonal) architecture pattern. By isolating circuit breaker implementations in adapter classes, the core business logic remains independent of specific circuit breaker technologies, making the system more maintainable, testable, and flexible.

## Features

- Circuit breaker pattern implementation
- Failure detection and handling
- Automatic recovery
- Configurable thresholds and timeouts
- Half-open state management
- Fallback mechanisms
- Circuit state monitoring and reporting
- Integration with various circuit breaker libraries

## Installation

```bash
go get github.com/abitofhelp/family-service/infrastructure/adapters/circuitwrapper
```

## Configuration

The circuit wrapper can be configured according to specific requirements. Here's an example of configuring the circuit wrapper:

```
// Pseudocode example - not actual Go code
// This demonstrates how to configure and use a circuit wrapper

// 1. Import necessary packages
import circuit, config, logging

// 2. Create a logger
logger = logging.NewLogger()

// 3. Configure the circuit breaker
circuitConfig = {
    name: "database-circuit",
    failureThreshold: 5,
    successThreshold: 2,
    timeout: 30 seconds,
    resetTimeout: 60 seconds,
    fallbackEnabled: true,
    monitoringEnabled: true,
    metricsEnabled: true
}

// 4. Create the circuit wrapper
circuitWrapper = circuit.NewCircuitWrapper(circuitConfig, logger)

// 5. Use the circuit wrapper
result, err = circuitWrapper.Execute(context, "database-operation", function() {
    // Operation that might fail
    return databaseOperation()
}, function() {
    // Fallback function
    return getCachedData()
})

// 6. Check circuit state
state = circuitWrapper.GetState("database-circuit")
if state == circuit.StateOpen {
    logger.Warn("Circuit is open, database operations will fail fast")
}
```

## API Documentation

### Core Concepts

The circuit wrapper follows these core concepts:

1. **Circuit Breaker Pattern**: Implements the circuit breaker pattern to prevent cascading failures
2. **Dependency Injection**: Receives dependencies through constructor injection
3. **Configuration**: Configured through a central configuration system
4. **Logging**: Uses a consistent logging approach
5. **Error Handling**: Handles circuit breaker errors gracefully

### Key Adapter Functions

```
// Pseudocode example - not actual Go code
// This demonstrates a circuit wrapper implementation

// Circuit wrapper structure
type CircuitWrapper {
    config        // Circuit wrapper configuration
    logger        // Logger for logging operations
    contextLogger // Context-aware logger
    circuits      // Map of circuit breakers
}

// Constructor for the circuit wrapper
function NewCircuitWrapper(config, logger) {
    return new CircuitWrapper {
        config: config,
        logger: logger,
        contextLogger: new ContextLogger(logger),
        circuits: {}
    }
}

// Method to execute a function with circuit breaker protection
function CircuitWrapper.Execute(context, name, operation, fallback) {
    // Implementation would include:
    // 1. Logging the operation with context
    // 2. Getting or creating the circuit breaker
    // 3. Checking if the circuit is open
    // 4. Executing the operation if allowed
    // 5. Handling success or failure
    // 6. Updating circuit state
    // 7. Executing fallback if needed
    // 8. Returning the result or error
}

// Method to get circuit state
function CircuitWrapper.GetState(name) {
    // Implementation would include:
    // 1. Getting the circuit breaker
    // 2. Returning its state
}
```

## Best Practices

1. **Separation of Concerns**: Keep circuit breaker logic separate from domain logic
2. **Interface Segregation**: Define focused circuit breaker interfaces in the domain layer
3. **Dependency Injection**: Use constructor injection for adapter dependencies
4. **Error Handling**: Handle circuit breaker errors gracefully
5. **Consistent Logging**: Use a consistent logging approach
6. **Monitoring**: Monitor circuit breaker states and transitions
7. **Testing**: Write unit and integration tests for circuit breakers
8. **Fallbacks**: Implement meaningful fallbacks for operations

## Troubleshooting

### Common Issues

#### Circuit Tripping Too Frequently

If your circuit breaker trips too frequently, consider the following:
- Increase the failure threshold
- Adjust the timeout values
- Implement retry mechanisms before failing
- Review the error detection logic
- Consider using a more sophisticated failure detection algorithm

#### Circuit Never Recovering

If your circuit breaker never recovers, consider the following:
- Verify the reset timeout is appropriate
- Ensure the success threshold is achievable
- Check that the half-open state is working correctly
- Implement health checks for dependent services
- Monitor circuit state transitions

## Related Components

- [Domain Layer](../../core/domain/README.md) - The domain layer that defines the circuit breaker ports
- [Application Layer](../../core/application/README.md) - The application layer that uses circuit breakers
- [Repository Wrapper](../repositorywrapper/README.md) - The repository wrapper that might use circuit breakers

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.