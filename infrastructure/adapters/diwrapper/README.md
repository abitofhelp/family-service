# Dependency Injection Wrapper

## Overview

The Dependency Injection Wrapper package provides a clean and flexible way to manage dependencies in the application. It abstracts away the details of the underlying dependency injection framework and provides a consistent interface for registering and resolving dependencies.

## Architecture

This package is part of the infrastructure layer in the Clean Architecture and Hexagonal Architecture patterns. It provides adapters for dependency injection that can be used by the application.

## Implementation Details

The Dependency Injection Wrapper implements the following design patterns:
- Factory Pattern: Creates instances of dependencies
- Adapter Pattern: Adapts external dependency injection libraries to the application's needs
- Service Locator Pattern: Provides a central registry for dependencies

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:
- [Dependency Injection Example](../../../examples/di/README.md) - Shows how to use the dependency injection wrapper

## Configuration

The Dependency Injection Wrapper can be configured with the following options:
- Singleton Registration: Configure which dependencies are registered as singletons
- Transient Registration: Configure which dependencies are registered as transients
- Scoped Registration: Configure which dependencies are registered as scoped

## Testing

The Dependency Injection Wrapper is tested through:
1. Unit Tests: Each dependency injection method has comprehensive unit tests
2. Integration Tests: Tests that verify the dependency injection wrapper works correctly with the application

## Design Notes

1. The Dependency Injection Wrapper uses a container-based approach to dependency management
2. Dependencies are registered at startup and resolved at runtime
3. The wrapper supports both constructor and property injection

## References

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Dependency Injection](https://en.wikipedia.org/wiki/Dependency_injection)