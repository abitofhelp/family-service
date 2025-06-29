# Integration Tests

## Overview

This directory contains integration tests for the Family Service with different database backends (MongoDB, PostgreSQL, and SQLite). These tests verify that the service works correctly with each supported database.

## Features

- **Multi-Database Testing**: Tests with MongoDB, PostgreSQL, and SQLite
- **End-to-End Verification**: Verifies the entire service stack
- **Database Operations**: Tests CRUD operations against real databases
- **Error Handling**: Tests error conditions and edge cases
- **Performance Testing**: Verifies performance with different databases

## Installation

These tests are part of the project and do not require separate installation. However, to run them, you'll need:

1. Docker and Docker Compose (for MongoDB and PostgreSQL tests)
2. Go 1.16 or later

## Quick Start

To run the integration tests:

```bash
# Run all integration tests
cd tests/integration
go test ./...

# Run tests for a specific database
go test -run TestFamilyServiceMongoDB
go test -run TestFamilyServicePostgreSQL
go test -run TestFamilyServiceSQLite
```

## API Documentation

### Core Types

This directory primarily contains test files rather than API types.

## Examples

There may be additional examples in the /EXAMPLES directory.

### Test Files

#### MongoDB Integration Tests

Tests the Family Service with MongoDB as the backend database.

```bash
go test -v -run TestFamilyServiceMongoDB
```

#### PostgreSQL Integration Tests

Tests the Family Service with PostgreSQL as the backend database.

```bash
go test -v -run TestFamilyServicePostgreSQL
```

#### SQLite Integration Tests

Tests the Family Service with SQLite as the backend database.

```bash
go test -v -run TestFamilyServiceSQLite
```

## Best Practices

1. **Database Cleanup**: Always clean up test data after tests
2. **Isolated Environments**: Use separate databases for testing
3. **Realistic Scenarios**: Test realistic usage scenarios
4. **Error Conditions**: Test error conditions and edge cases
5. **Performance Considerations**: Be aware of test performance with different databases

## Troubleshooting

### Common Issues

#### Database Connection Failures

If tests fail to connect to the database, ensure the database is running and accessible.

#### Slow Tests

Integration tests can be slow. Consider running only the tests you need during development.

## Related Components

- [MongoDB Adapter](../../infrastructure/adapters/mongo/README.md) - MongoDB adapter tested by these tests
- [PostgreSQL Adapter](../../infrastructure/adapters/postgres/README.md) - PostgreSQL adapter tested by these tests
- [SQLite Adapter](../../infrastructure/adapters/sqlite/README.md) - SQLite adapter tested by these tests

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../LICENSE) file for details.