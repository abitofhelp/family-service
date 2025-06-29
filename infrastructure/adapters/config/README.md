# Configuration Adapter

## Overview

The Configuration Adapter package provides functionality for loading and accessing application configuration from various sources.

## Features

- **Environment Variables**: Load configuration from environment variables
- **Configuration Files**: Load configuration from YAML, JSON, and TOML files
- **Hierarchical Configuration**: Support for nested configuration values
- **Default Values**: Provide default values for configuration options
- **Validation**: Validate configuration values at startup

## Installation

```bash
go get github.com/abitofhelp/family-service/infrastructure/adapters/config
```

## Quick Start

See the [Quick Start example](../../../EXAMPLES/config/basic_usage/README.md) for a complete, runnable example of how to use the configuration adapter.

## Configuration

The Configuration Adapter can be configured with the following options:
- Environment Variables: Configure which environment variables are used
- Configuration Files: Configure which configuration files are loaded

## API Documentation

### Core Types

Description of the main types provided by the component.

#### Config

The main configuration interface that provides access to configuration values.

```
// Config interface for accessing configuration values
type Config interface {
    GetString(key string) string
    GetInt(key string) int
    GetBool(key string) bool
    GetFloat(key string) float64
    GetDuration(key string) time.Duration
}
```

#### ConfigOptions

Options for configuring the configuration adapter.

```
// ConfigOptions for configuring the configuration adapter
type ConfigOptions struct {
    ConfigFile string
    EnvPrefix  string
}
```

### Key Methods

Description of the key methods provided by the component.

#### New

Creates a new configuration instance.

```
// New creates a new configuration instance
func New(configFile string) (Config, error)
```

#### NewWithOptions

Creates a new configuration instance with custom options.

```
// NewWithOptions creates a new configuration instance with custom options
func NewWithOptions(options ConfigOptions) (Config, error)
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory.

### Configuration Example

```go
package main

import (
    "fmt"
    "github.com/abitofhelp/family-service/infrastructure/adapters/config"
)

func main() {
    // Initialize the configuration
    cfg, err := config.New("config.yaml")
    if err != nil {
        fmt.Printf("Error initializing configuration: %v\n", err)
        return
    }

    // Access configuration values
    dbHost := cfg.GetString("database.host")
    dbPort := cfg.GetInt("database.port")

    fmt.Printf("Database connection: %s:%d\n", dbHost, dbPort)

    // You can also use environment variables that override file settings
    // export APP_DATABASE_HOST=localhost
    // export APP_DATABASE_PORT=5432
}
```

## Best Practices

1. **Use Environment Variables for Secrets**: Never store secrets in configuration files
2. **Validate Configuration at Startup**: Fail fast if required configuration is missing
3. **Use Hierarchical Keys**: Organize configuration using hierarchical keys
4. **Provide Default Values**: Always provide sensible default values
5. **Use Strong Typing**: Use the typed getter methods instead of generic ones

## Troubleshooting

### Common Issues

#### Missing Configuration File

If the configuration file is missing, the adapter will return an error. Make sure the file exists and is readable.

#### Invalid Configuration Values

If the configuration values are invalid, the adapter will return an error. Make sure the values are of the correct type.

## Related Components

- [Logging Wrapper](../loggingwrapper/README.md) - Used for logging configuration errors
- [Error Wrapper](../errorswrapper/README.md) - Used for handling configuration errors

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details.
