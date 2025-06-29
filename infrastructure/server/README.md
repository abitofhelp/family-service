# Server Package

## Overview

The Server package provides functionality for HTTP server management in the Family Service application. It extends the standard Go HTTP server with additional features such as logging, context-aware logging, middleware integration, and graceful shutdown capabilities. This package is designed to be used by the GraphQL server and other HTTP-based services in the application.

## Architecture

The Server package sits in the infrastructure layer of the application and is responsible for managing the HTTP server lifecycle. It follows these principles:

- **Clean Architecture**: The server is part of the infrastructure layer, providing services to the interface layer
- **Separation of Concerns**: The server focuses solely on HTTP server management, delegating request handling to the provided handler
- **Dependency Injection**: The server receives its dependencies through constructor injection
- **Context Propagation**: The server propagates context throughout the request lifecycle for cancellation and value propagation

The package is organized into:

- **Server**: The main server type that extends the standard HTTP server
- **Config**: Configuration options for the server
- **Lifecycle Methods**: Methods for starting and shutting down the server

## Implementation Details

The Server package implements the following design patterns:

1. **Decorator Pattern**: Extends the standard HTTP server with additional functionality
2. **Dependency Injection Pattern**: Receives dependencies through constructor injection
3. **Builder Pattern**: Uses a configuration object to build the server
4. **Observer Pattern**: Logs server events for observability

Key implementation details:

- **Extended HTTP Server**: Extends the standard Go HTTP server with additional functionality
- **Middleware Integration**: Integrates with the middleware package for request processing
- **Logging Integration**: Integrates with the logging package for server event logging
- **Context-Aware Logging**: Uses context-aware logging for request-scoped logging
- **Graceful Shutdown**: Implements graceful shutdown to ensure in-flight requests are completed
- **Timeout Management**: Manages various timeouts for request processing

## Examples

Example of using the Server package:

```
// Create server configuration
cfg := server.NewConfig(
    "8080",                // Port
    10*time.Second,        // Read timeout
    10*time.Second,        // Write timeout
    120*time.Second,       // Idle timeout
    30*time.Second,        // Shutdown timeout
)

// Create a new server
srv := server.New(
    cfg,
    handler,               // HTTP handler
    logger,                // Logger
    contextLogger,         // Context-aware logger
)

// Start the server
srv.Start()

// Gracefully shut down the server
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
if err := srv.Shutdown(ctx); err != nil {
    // Handle error
}
```

## Configuration

The Server package is configured using the Config struct. The following configuration options are available:

- **Port**: The port on which the server will listen
- **ReadTimeout**: Maximum duration for reading the entire request
- **WriteTimeout**: Maximum duration before timing out writes of the response
- **IdleTimeout**: Maximum amount of time to wait for the next request
- **ShutdownTimeout**: Maximum time to wait for server shutdown

Example configuration:

```
// Create server configuration
cfg := server.NewConfig(
    "8080",                // Port
    10*time.Second,        // Read timeout
    10*time.Second,        // Write timeout
    120*time.Second,       // Idle timeout
    30*time.Second,        // Shutdown timeout
)
```

## Testing

The Server package is tested through:

1. **Unit Tests**: Each method has unit tests
2. **Integration Tests**: Tests that verify the server works correctly with real HTTP requests
3. **Mock Tests**: Tests that use mock dependencies to isolate the server

Key testing approaches:

- **Mock Dependencies**: Tests use mock dependencies to isolate the server
- **HTTP Testing**: Tests use the httptest package to test HTTP functionality
- **Timeout Testing**: Tests verify that timeouts work correctly
- **Graceful Shutdown Testing**: Tests verify that graceful shutdown works correctly
- **Error Handling**: Tests verify that errors are properly handled and logged

## Design Notes

1. **Extended HTTP Server**: The server extends the standard Go HTTP server to provide additional functionality while maintaining compatibility
2. **Middleware Integration**: The server integrates with the middleware package to provide a consistent approach to middleware
3. **Logging Integration**: The server integrates with the logging package to provide comprehensive logging of server events
4. **Context-Aware Logging**: The server uses context-aware logging to provide request-scoped logging
5. **Graceful Shutdown**: The server implements graceful shutdown to ensure in-flight requests are completed
6. **Non-Blocking Start**: The server starts in a non-blocking manner to allow the application to perform other initialization tasks

## References

- [Go HTTP Server](https://golang.org/pkg/net/http/#Server) - The standard Go HTTP server
- [Context Package](https://golang.org/pkg/context/) - Go's context package for cancellation and value propagation
- [Graceful Shutdown](https://golang.org/pkg/net/http/#Server.Shutdown) - Go's graceful shutdown functionality
- [GraphQL Server](../../cmd/server/graphql/README.md) - Uses this server package to serve GraphQL requests
- [Middleware Package](../../infrastructure/adapters/middleware/README.md) - Provides middleware for request processing
- [Logging Package](../../infrastructure/adapters/loggingwrapper/README.md) - Provides logging functionality