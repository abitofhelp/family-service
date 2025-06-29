# PostgreSQL Adapter

## Overview

The PostgreSQL Adapter package provides functionality for interacting with PostgreSQL databases. It implements the repository interfaces defined in the domain layer, allowing the application to store and retrieve data from PostgreSQL.

## Architecture

This package is part of the infrastructure layer in the Clean Architecture and Hexagonal Architecture patterns. It provides adapters for PostgreSQL that can be used by the application.

## Implementation Details

The PostgreSQL Adapter implements the following design patterns:
- Repository Pattern: Provides a clean interface for data access
- Adapter Pattern: Adapts the PostgreSQL driver to the application's needs
- Factory Pattern: Creates instances of PostgreSQL repositories

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:
- [Repository Example](../../../examples/repository/README.md) - Shows how to use the repository adapters

## Configuration

The PostgreSQL Adapter can be configured with the following options:
- Connection String: Configure the PostgreSQL connection string
- Database Name: Configure the PostgreSQL database name
- Table Names: Configure the PostgreSQL table names
- Connection Pool Size: Configure the PostgreSQL connection pool size
- Timeout Settings: Configure the PostgreSQL timeout settings

## Testing

The PostgreSQL Adapter is tested through:
1. Unit Tests: Each repository method has comprehensive unit tests
2. Integration Tests: Tests that verify the PostgreSQL adapter works correctly with a real PostgreSQL instance

## Design Notes

1. The PostgreSQL Adapter uses the official PostgreSQL Go driver
2. Repositories are implemented as adapters for the domain repository interfaces
3. The adapter handles PostgreSQL-specific concerns like connection management and query building
4. Error handling is consistent with the application's error handling strategy

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [PostgreSQL Go Driver](https://github.com/lib/pq)
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)