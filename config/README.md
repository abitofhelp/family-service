# Configuration

## Overview

This directory contains configuration files for the Family Service in different environments and deployment modes.

## Features

- **Development CLI Configuration**: Configuration for running the service in development mode via CLI
- **Development Docker Configuration**: Configuration for running the service in development mode via Docker
- **Environment-Specific Settings**: Settings tailored to different environments
- **Service Configuration**: Database, logging, security, and other service settings
- **Deployment Configuration**: Settings for different deployment scenarios

## Installation

These configuration files are part of the project and do not require separate installation.

## Quick Start

To use these configuration files, specify the appropriate file when starting the service:

```bash
# For CLI development
go run cmd/server/main.go --config=config/config.dev.cli.yaml

# For Docker development (usually set in docker-compose.yml)
# FAMILY_SERVICE_CONFIG=/app/config/config.dev.docker.yaml
```

## Configuration

### Configuration Files

#### config.dev.cli.yaml

This file contains configuration for running the service in development mode via CLI:

```yaml
# Example structure (actual values may differ)
server:
  port: 8080
  timeout: 30s
database:
  type: sqlite
  path: ./data/dev/sqlite/family.db
logging:
  level: debug
```

#### config.dev.docker.yaml

This file contains configuration for running the service in development mode via Docker:

```yaml
# Example structure (actual values may differ)
server:
  port: 8080
  timeout: 30s
database:
  type: postgres
  host: postgres
  port: 5432
  name: family
  user: postgres
  password: postgres
logging:
  level: debug
```

## API Documentation

### Core Types

This directory does not contain code, only configuration files.

## Examples

There may be additional examples in the /EXAMPLES directory.

## Best Practices

1. **Environment Variables**: Use environment variables for sensitive information
2. **Configuration Validation**: Validate configuration at service startup
3. **Default Values**: Provide sensible default values for optional settings
4. **Documentation**: Document all configuration options
5. **Version Control**: Track configuration changes in version control

## Troubleshooting

### Common Issues

#### Configuration Not Found

If the service cannot find the configuration file, ensure the path is correct and the file exists.

#### Invalid Configuration

If the service fails to start due to invalid configuration, check the logs for specific error messages.

## Related Components

- [Infrastructure Config](../infrastructure/adapters/config/README.md) - Configuration loading and validation
- [Server](../infrastructure/server/README.md) - Server that uses these configuration files

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../LICENSE) file for details.