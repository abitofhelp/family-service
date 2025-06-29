# Infrastructure Adapters - Profiling

## Overview

The Profiling adapter provides implementations for profiling-related ports defined in the core domain and application layers. This adapter connects the application to profiling frameworks and libraries, following the Ports and Adapters (Hexagonal) architecture pattern. By isolating profiling implementations in adapter classes, the core business logic remains independent of specific profiling technologies, making the system more maintainable, testable, and flexible.

## Features

- CPU profiling
- Memory profiling
- Goroutine profiling
- Block profiling
- Mutex profiling
- Trace collection
- Profile visualization
- Performance metrics collection
- Profiling endpoint exposure

## Installation

```bash
go get github.com/abitofhelp/family-service/infrastructure/adapters/profiling
```

## Configuration

The profiling adapter can be configured according to specific requirements. Here's an example of configuring the profiling adapter:

```
// Pseudocode example - not actual Go code
// This demonstrates how to configure and use a profiling adapter

// 1. Import necessary packages
import profiling, config, logging

// 2. Create a logger
logger = logging.NewLogger()

// 3. Configure the profiler
profilingConfig = {
    enabled: true,
    cpuProfileEnabled: true,
    memProfileEnabled: true,
    blockProfileEnabled: true,
    mutexProfileEnabled: true,
    traceEnabled: true,
    profilePath: "/tmp/profiles",
    httpEndpoint: "/debug/pprof",
    sampleRate: 100,
    profileDuration: 30 seconds
}

// 4. Create the profiling adapter
profilingAdapter = profiling.NewProfilingAdapter(profilingConfig, logger)

// 5. Use the profiling adapter
profilingAdapter.Start()

// 6. Capture a CPU profile
profilingAdapter.StartCPUProfile("api-request")
// ... code to profile ...
profilingAdapter.StopCPUProfile()

// 7. Capture a memory profile
profilingAdapter.CaptureMemoryProfile("after-processing")

// 8. Stop profiling when done
defer profilingAdapter.Stop()
```

## API Documentation

### Core Concepts

The profiling adapter follows these core concepts:

1. **Adapter Pattern**: Implements profiling ports defined in the core domain or application layer
2. **Dependency Injection**: Receives dependencies through constructor injection
3. **Configuration**: Configured through a central configuration system
4. **Logging**: Uses a consistent logging approach
5. **Error Handling**: Handles profiling errors gracefully without affecting core functionality

### Key Adapter Functions

```
// Pseudocode example - not actual Go code
// This demonstrates a profiling adapter implementation

// Profiling adapter structure
type ProfilingAdapter {
    config        // Profiling configuration
    logger        // Logger for logging operations
    contextLogger // Context-aware logger
    active        // Flag indicating if profiling is active
}

// Constructor for the profiling adapter
function NewProfilingAdapter(config, logger) {
    return new ProfilingAdapter {
        config: config,
        logger: logger,
        contextLogger: new ContextLogger(logger),
        active: false
    }
}

// Method to start profiling
function ProfilingAdapter.Start() {
    // Implementation would include:
    // 1. Logging the operation
    // 2. Setting up profiling based on configuration
    // 3. Starting HTTP server for pprof if configured
    // 4. Handling startup errors
    // 5. Setting active flag
}

// Method to capture a CPU profile
function ProfilingAdapter.StartCPUProfile(name) {
    // Implementation would include:
    // 1. Logging the operation
    // 2. Creating a profile file
    // 3. Starting CPU profiling
    // 4. Handling profiling errors
}
```

## Best Practices

1. **Separation of Concerns**: Keep profiling logic separate from domain logic
2. **Interface Segregation**: Define focused profiling interfaces in the domain layer
3. **Dependency Injection**: Use constructor injection for adapter dependencies
4. **Error Handling**: Handle profiling errors gracefully without affecting core functionality
5. **Consistent Logging**: Use a consistent logging approach
6. **Configuration**: Configure profiling through a central configuration system
7. **Resource Management**: Be mindful of the performance impact of profiling in production

## Troubleshooting

### Common Issues

#### Performance Impact

If profiling is impacting application performance, consider the following:
- Disable profiling in production or use sampling
- Limit the types of profiling enabled
- Use shorter profiling durations
- Profile specific operations rather than continuous profiling
- Adjust sampling rates for block and mutex profiling

#### File System Issues

If you encounter issues with profile file generation, check the following:
- Ensure the profile directory exists and is writable
- Check disk space availability
- Verify file permissions
- Use absolute paths for profile locations
- Implement cleanup of old profile files

## Related Components

- [Domain Layer](../../core/domain/README.md) - The domain layer that defines the profiling ports
- [Application Layer](../../core/application/README.md) - The application layer that uses profiling
- [Interface Adapters](../../interface/adapters/README.md) - The interface adapters that use profiling

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.