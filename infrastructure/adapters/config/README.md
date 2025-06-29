# Infrastructure Adapters - Configuration

## Overview

The Configuration adapter provides implementations for configuration-related ports defined in the core domain and application layers. This adapter connects the application to configuration sources and frameworks, following the Ports and Adapters (Hexagonal) architecture pattern. By isolating configuration implementations in adapter classes, the core business logic remains independent of specific configuration technologies, making the system more maintainable, testable, and flexible.

## Features

- Loading configuration from various sources (files, environment variables, etc.)
- Configuration validation
- Dynamic configuration updates
- Configuration change notifications
- Hierarchical configuration support
- Environment-specific configuration
- Secure configuration handling (for secrets)

## Installation

```bash
go get github.com/abitofhelp/family-service/infrastructure/adapters/config
```

## Configuration

The configuration adapter itself can be configured according to specific requirements. Here's an example of setting up the configuration adapter:

```
// Pseudocode example - not actual Go code
// This demonstrates how to set up and use a configuration adapter

// 1. Import necessary packages
import config, logging

// 2. Create a logger
logger = logging.NewLogger()

// 3. Configure the configuration adapter
configOptions = {
    configPath: "./config",
    environment: "development",
    defaultConfigFile: "config.yaml",
    envPrefix: "APP_",
    watchForChanges: true,
    reloadInterval: 30 seconds
}

// 4. Create the configuration adapter
configAdapter = config.NewConfigAdapter(configOptions, logger)

// 5. Use the configuration adapter
dbConfig = configAdapter.GetDatabaseConfig()
serverPort = configAdapter.GetInt("server.port", 8080)
apiKeys = configAdapter.GetStringMap("security.apiKeys")
```

## API Documentation

### Core Concepts

The configuration adapter follows these core concepts:

1. **Adapter Pattern**: Implements configuration ports defined in the core domain or application layer
2. **Dependency Injection**: Receives dependencies through constructor injection
3. **Configuration Sources**: Supports multiple configuration sources with priority order
4. **Logging**: Uses a consistent logging approach
5. **Error Handling**: Handles configuration errors gracefully

### Key Adapter Functions

```
// Pseudocode example - not actual Go code
// This demonstrates a configuration adapter implementation

// Configuration adapter structure
type ConfigAdapter {
    options       // Configuration options
    logger        // Logger for logging operations
    contextLogger // Context-aware logger
    provider      // Configuration provider
}

// Constructor for the configuration adapter
function NewConfigAdapter(options, logger) {
    provider = createConfigProvider(options)
    return new ConfigAdapter {
        options: options,
        logger: logger,
        contextLogger: new ContextLogger(logger),
        provider: provider
    }
}

// Method to get a string value
function ConfigAdapter.GetString(key, defaultValue) {
    // Implementation would include:
    // 1. Logging the operation
    // 2. Retrieving the value from the configuration provider
    // 3. Handling errors gracefully
    // 4. Returning the value or default
}

// Method to get a typed configuration section
function ConfigAdapter.GetDatabaseConfig() {
    // Implementation would include:
    // 1. Logging the operation
    // 2. Retrieving the database configuration section
    // 3. Validating the configuration
    // 4. Mapping to a typed structure
    // 5. Returning the typed configuration
}
```

## Best Practices

1. **Separation of Concerns**: Keep configuration logic separate from domain logic
2. **Interface Segregation**: Define focused configuration interfaces in the domain layer
3. **Dependency Injection**: Use constructor injection for adapter dependencies
4. **Error Handling**: Handle configuration errors gracefully with sensible defaults
5. **Consistent Logging**: Use a consistent logging approach
6. **Configuration Validation**: Validate configuration at startup
7. **Testing**: Write unit and integration tests for configuration adapters

## Troubleshooting

### Common Issues

#### Configuration Loading Failures

If you encounter issues with configuration loading, check the following:
- Configuration files exist in the expected locations
- File permissions allow reading the configuration files
- Environment variables are set correctly
- Configuration format is valid (JSON, YAML, etc.)

#### Configuration Type Mismatches

If you encounter type mismatch issues with configuration values, consider the following:
- Validate configuration values against expected types
- Provide clear error messages for type mismatches
- Use default values when types don't match
- Consider using a schema for configuration validation

## Related Components

- [Domain Layer](../../core/domain/README.md) - The domain layer that defines the configuration ports
- [Application Layer](../../core/application/README.md) - The application layer that uses configuration
- [Interface Adapters](../../interface/adapters/README.md) - The interface adapters that use configuration

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.