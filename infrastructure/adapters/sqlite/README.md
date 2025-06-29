# SQLite Adapter

## Overview

The SQLite Adapter package provides functionality for interacting with SQLite databases. It implements the repository interfaces defined in the domain layer, allowing the application to store and retrieve data from SQLite.

## Architecture

This package is part of the infrastructure layer in the Clean Architecture and Hexagonal Architecture patterns. It provides adapters for SQLite that can be used by the application.

## Implementation Details

The SQLite Adapter implements the following design patterns:
- Repository Pattern: Provides a clean interface for data access
- Adapter Pattern: Adapts the SQLite driver to the application's needs
- Factory Pattern: Creates instances of SQLite repositories

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:
- [Repository Example](../../../examples/repository/README.md) - Shows how to use the repository adapters

## Configuration

The SQLite Adapter can be configured with the following options:
- Database Path: Configure the SQLite database file path
- Table Names: Configure the SQLite table names
- Connection Pool Size: Configure the SQLite connection pool size
- Timeout Settings: Configure the SQLite timeout settings
- Journal Mode: Configure the SQLite journal mode

## Testing

The SQLite Adapter is tested through:
1. Unit Tests: Each repository method has comprehensive unit tests
2. Integration Tests: Tests that verify the SQLite adapter works correctly with a real SQLite database

## Design Notes

1. The SQLite Adapter uses the official SQLite Go driver
2. Repositories are implemented as adapters for the domain repository interfaces
3. The adapter handles SQLite-specific concerns like connection management and query building
4. Error handling is consistent with the application's error handling strategy
5. The adapter supports both in-memory and file-based SQLite databases

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [SQLite Go Driver](https://github.com/mattn/go-sqlite3)
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)