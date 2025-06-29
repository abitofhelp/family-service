# MongoDB Adapter

## Overview

The MongoDB Adapter package provides functionality for interacting with MongoDB databases. It implements the repository interfaces defined in the domain layer, allowing the application to store and retrieve data from MongoDB.

## Architecture

This package is part of the infrastructure layer in the Clean Architecture and Hexagonal Architecture patterns. It provides adapters for MongoDB that can be used by the application.

## Implementation Details

The MongoDB Adapter implements the following design patterns:
- Repository Pattern: Provides a clean interface for data access
- Adapter Pattern: Adapts the MongoDB driver to the application's needs
- Factory Pattern: Creates instances of MongoDB repositories

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:
- [Repository Example](../../../examples/repository/README.md) - Shows how to use the repository adapters

## Configuration

The MongoDB Adapter can be configured with the following options:
- Connection String: Configure the MongoDB connection string
- Database Name: Configure the MongoDB database name
- Collection Names: Configure the MongoDB collection names
- Connection Pool Size: Configure the MongoDB connection pool size
- Timeout Settings: Configure the MongoDB timeout settings

## Testing

The MongoDB Adapter is tested through:
1. Unit Tests: Each repository method has comprehensive unit tests
2. Integration Tests: Tests that verify the MongoDB adapter works correctly with a real MongoDB instance

## Design Notes

1. The MongoDB Adapter uses the official MongoDB Go driver
2. Repositories are implemented as adapters for the domain repository interfaces
3. The adapter handles MongoDB-specific concerns like connection management and query building
4. Error handling is consistent with the application's error handling strategy

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [MongoDB Go Driver](https://github.com/mongodb/mongo-go-driver)
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)