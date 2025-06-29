// Copyright (c) 2025 A Bit of Help, Inc.

// Package main is the entry point for the GraphQL server.
//
// This package initializes and runs the GraphQL server for the family service.
// It handles:
// - Configuration loading
// - Dependency injection
// - Logger initialization
// - Telemetry setup (metrics and tracing)
// - HTTP server setup with GraphQL endpoints
// - Graceful shutdown
//
// The server follows a clean startup sequence that ensures all components
// are properly initialized before the server starts accepting requests.
// Error handling during startup is designed to fail fast if critical
// components cannot be initialized.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/abitofhelp/family-service/cmd/server/graphql/di"
	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	infratelemetry "github.com/abitofhelp/family-service/infrastructure/adapters/telemetrywrapper"
	"github.com/abitofhelp/family-service/infrastructure/server"
	"github.com/abitofhelp/family-service/interface/adapters/graphql/generated"
	"github.com/abitofhelp/family-service/interface/adapters/graphql/resolver"
	pkgconfig "github.com/abitofhelp/servicelib/config"
	"github.com/abitofhelp/servicelib/graphql"
	"github.com/abitofhelp/servicelib/health"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/shutdown"
	"github.com/abitofhelp/servicelib/telemetry"
	"go.uber.org/zap"
)

// initBasicLogger creates a basic logger for use during startup before the main logger is configured.
//
// This function creates a simple production-level logger that can be used
// during the early stages of application startup, before the configuration
// is loaded and the main logger is initialized. If the logger creation fails,
// the application will exit immediately, as logging is essential for observability.
//
// Returns:
//   - A configured zap.Logger instance ready for use
//
// Panics:
//   - If logger creation fails (exits the application with status code 1)
func initBasicLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		// If we can't create a logger, we're in serious trouble
		fmt.Printf("Failed to create basic logger: %v\n", err)
		os.Exit(1)
	}
	return logger
}

// loadConfig loads the application configuration from environment variables and files.
//
// This function attempts to load the application configuration using the
// config package's LoadConfig function. It logs the start and result of
// the configuration loading process, providing visibility into this critical
// startup step.
//
// Parameters:
//   - logger: The logger to use for logging the configuration loading process
//
// Returns:
//   - A pointer to the loaded configuration if successful
//   - An error if configuration loading fails
func loadConfig(logger *zap.Logger) (*config.Config, error) {
	logger.Info("Loading application configuration")
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load application configuration", zap.Error(err))
		return nil, err
	}
	logger.Info("Application configuration loaded successfully")
	return cfg, nil
}

// initLogger initializes the application logger based on configuration settings.
//
// This function creates a properly configured logger using the settings from
// the application configuration. It uses the basicLogger to log the initialization
// process, providing visibility into this critical startup step. The resulting
// logger will be used throughout the application for all logging needs.
//
// Parameters:
//   - cfg: The application configuration containing logging settings
//   - basicLogger: A temporary logger to use during initialization
//
// Returns:
//   - A properly configured zap.Logger instance if successful
//   - An error if logger initialization fails
func initLogger(cfg *config.Config, basicLogger *zap.Logger) (*zap.Logger, error) {
	basicLogger.Info("Initializing application logger",
		zap.String("version", cfg.App.Version),
		zap.String("level", cfg.Log.Level),
		zap.Bool("development", cfg.Log.Development))

	logger, err := logging.NewLogger(cfg.Log.Level, cfg.Log.Development)
	if err != nil {
		basicLogger.Error("Failed to initialize logger", zap.Error(err))
		return nil, err
	}

	logger.Info("Application logger initialized successfully")
	return logger, nil
}

// initContainer initializes the dependency injection container with all application services.
//
// This function creates and configures the dependency injection container that
// will hold all application services and their dependencies. The container is
// responsible for managing the lifecycle of these services and providing them
// to the parts of the application that need them.
//
// The container initialization is a critical step in the application startup
// process, as it creates all the services needed by the application, including
// repositories, domain services, and application services.
//
// Parameters:
//   - ctx: The context for the initialization process
//   - logger: The logger to use for logging the initialization process
//   - cfg: The application configuration
//
// Returns:
//   - A configured dependency injection container if successful
//   - An error if container initialization fails
func initContainer(ctx context.Context, logger *zap.Logger, cfg *config.Config) (*di.Container, error) {
	logger.Info("Initializing dependency injection container")

	container, err := di.NewContainer(ctx, logger, cfg)
	if err != nil {
		logger.Error("Failed to initialize dependency injection container", zap.Error(err))
		return nil, err
	}

	logger.Info("Dependency injection container initialized successfully")
	return container, nil
}

// setupRoutes sets up the HTTP routes for the application including GraphQL and health check endpoints.
//
// This function configures all the HTTP routes that the server will handle, including:
// - GraphQL API endpoint
// - GraphQL Playground for interactive API exploration
// - Health check endpoint for monitoring
// - Telemetry endpoints for metrics and tracing
//
// It also applies middleware and sets up telemetry (metrics and tracing) based on
// the application configuration.
//
// Parameters:
//   - ctx: The context for the setup process
//   - container: The dependency injection container with application services
//   - logger: The logger to use for logging the setup process
//   - cfg: The application configuration
//
// Returns:
//   - An HTTP handler with all routes configured
//   - A shutdown function for telemetry components
//   - An error if route setup fails
func setupRoutes(ctx context.Context, container *di.Container, logger *zap.Logger, cfg *config.Config) (http.Handler, func(), error) {
	logger.Info("Setting up HTTP routes")

	// Create a container adapter for health checks
	adapter := health.NewGenericContainerAdapter(container)

	// Create a ServeMux for routing
	mux := http.NewServeMux()

	// Set up telemetry (metrics and tracing)
	telemetryShutdown, err := setupTelemetry(ctx, mux, cfg, logger)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to set up telemetry: %w", err)
	}

	// Set up GraphQL endpoints
	setupGraphQLEndpoints(mux, container)

	// Health check endpoint
	healthEndpoint := cfg.Server.HealthEndpoint
	logger.Info("Setting up health check endpoint", zap.String("endpoint", healthEndpoint))
	configAdapter := pkgconfig.NewGenericConfigAdapter(cfg)
	mux.Handle(healthEndpoint, health.NewHandler(adapter, logger, configAdapter))

	logger.Info("HTTP routes set up successfully")
	return mux, telemetryShutdown, nil
}

// setupTelemetry sets up telemetry (metrics and tracing) based on application configuration.
//
// This function configures the telemetry components of the application, including:
// - Prometheus metrics endpoint for collecting and exposing application metrics
// - Distributed tracing for tracking requests across service boundaries
//
// The telemetry configuration is based on the application configuration, allowing
// for flexible deployment in different environments. The function returns a shutdown
// function that should be called during application shutdown to ensure proper cleanup
// of telemetry resources.
//
// Parameters:
//   - ctx: The context for the setup process
//   - mux: The HTTP ServeMux to register telemetry endpoints on
//   - cfg: The application configuration with telemetry settings
//   - logger: The logger to use for logging the setup process
//
// Returns:
//   - A shutdown function that should be called during application shutdown
//   - An error if telemetry setup fails
func setupTelemetry(ctx context.Context, mux *http.ServeMux, cfg *config.Config, logger *zap.Logger) (func(), error) {
	var shutdownFuncs []func()

	// Set up metrics if enabled
	if cfg.Telemetry.Exporters.Metrics.Prometheus.Enabled {
		metricsPath := cfg.Telemetry.Exporters.Metrics.Prometheus.Path
		logger.Info("Setting up Prometheus metrics endpoint",
			zap.String("path", metricsPath),
			zap.String("listen", cfg.Telemetry.Exporters.Metrics.Prometheus.Listen))
		mux.Handle(metricsPath, telemetry.CreatePrometheusHandler())
	} else {
		logger.Info("Prometheus metrics endpoint is disabled")
	}

	// Set up tracing
	tracingShutdown, err := infratelemetry.InitTracing(ctx, cfg, logger)
	if err != nil {
		logger.Error("Failed to initialize tracing", zap.Error(err))
		return nil, err
	}
	shutdownFuncs = append(shutdownFuncs, tracingShutdown)

	// Return a combined shutdown function
	return func() {
		for _, fn := range shutdownFuncs {
			fn()
		}
	}, nil
}

// setupGraphQLEndpoints sets up all GraphQL-related endpoints and handlers.
//
// This function configures the GraphQL API endpoints and tools, including:
// - The main GraphQL API endpoint for handling queries and mutations
// - GraphQL Playground for interactive API exploration
// - GraphiQL interface for a more feature-rich API exploration experience
// - A landing page at the root URL
//
// It uses the resolver from the dependency injection container to handle
// GraphQL operations and sets up authorization directives for securing
// the API.
//
// Parameters:
//   - mux: The HTTP ServeMux to register GraphQL endpoints on
//   - container: The dependency injection container with application services
func setupGraphQLEndpoints(mux *http.ServeMux, container *di.Container) {
	// Get the resolver
	resolverInstance := resolver.NewResolver(container.GetFamilyApplicationService(), container.GetFamilyMapper())

	// Initialize GraphQL schema
	schema := generated.NewExecutableSchema(generated.Config{
		Resolvers: resolverInstance,
		Directives: generated.DirectiveRoot{
			IsAuthorized: resolverInstance.IsAuthorized,
		},
	})

	// Create GraphQL server with configuration
	gqlServerConfig := graphql.NewDefaultServerConfig()
	gqlServer := graphql.NewServer(schema, container.GetContextLogger(), gqlServerConfig)

	// Serve the landing page at the root
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "interface/adapters/graphql/static/index.html")
	})

	// GraphQL endpoint
	mux.Handle("/graphql", gqlServer)

	// GraphQL Playground
	mux.Handle("/playground", playground.Handler("GraphQL Playground", "/query"))

	// Custom GraphiQL interface
	mux.HandleFunc("/graphiql", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "interface/adapters/graphql/static/graphiql.html")
	})
}

// startServer creates and starts the HTTP server with the configured handler.
//
// This function initializes the HTTP server with the provided handler and
// configuration settings. It configures important server parameters like
// timeouts and port, and then starts the server in a non-blocking way
// (in a separate goroutine).
//
// The server is configured with sensible defaults for production use,
// including appropriate timeout settings to prevent resource exhaustion
// under high load or when clients disconnect unexpectedly.
//
// Parameters:
//   - handler: The HTTP handler that will process all incoming requests
//   - cfg: The application configuration with server settings
//   - logger: The logger for server-level logging
//   - contextLogger: The context-aware logger for request-level logging
//
// Returns:
//   - A running server instance that can be used for shutdown
func startServer(handler http.Handler, cfg *config.Config, logger *zap.Logger, contextLogger *logging.ContextLogger) *server.Server {
	serverConfig := server.NewConfig(
		cfg.Server.Port,
		cfg.Server.ReadTimeout,
		cfg.Server.WriteTimeout,
		cfg.Server.IdleTimeout,
		cfg.Server.ShutdownTimeout,
	)

	srv := server.New(serverConfig, handler, logger, contextLogger)
	srv.Start()

	logger.Info("HTTP server started successfully")
	return srv
}

// setupGracefulShutdown sets up graceful shutdown for the server and all components.
//
// This function creates a shutdown function that will be called when the
// application receives a termination signal (e.g., SIGINT or SIGTERM).
// The shutdown function:
// 1. Cancels the root context to signal all operations to stop
// 2. Creates a separate context with a timeout for server shutdown
// 3. Calls the server's Shutdown method to gracefully close all connections
//
// Graceful shutdown ensures that in-flight requests are allowed to complete
// (up to the shutdown timeout) before the server exits, preventing abrupt
// connection termination that could lead to errors for clients.
//
// Parameters:
//   - rootCtx: The root context for the application
//   - rootCancel: The cancel function for the root context
//   - srv: The HTTP server to shut down
//   - cfg: The application configuration with shutdown timeout settings
//
// Returns:
//   - A function that will perform the graceful shutdown when called
func setupGracefulShutdown(rootCtx context.Context, rootCancel context.CancelFunc, srv *server.Server, cfg *config.Config) func() error {
	return func() error {
		// Cancel the root context to signal all operations to stop
		rootCancel()

		// Create a separate context for server shutdown with a timeout
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
		defer shutdownCancel()

		// Shutdown the server
		return srv.Shutdown(shutdownCtx)
	}
}

// main is the entry point for the GraphQL server application.
//
// This function orchestrates the startup sequence for the application:
// 1. Initialize a basic logger for startup logging
// 2. Create a root context with cancellation for the application
// 3. Load application configuration
// 4. Initialize the main logger based on configuration
// 5. Initialize the dependency injection container with all services
// 6. Set up HTTP routes including GraphQL and health check endpoints
// 7. Apply authentication middleware
// 8. Start the HTTP server
// 9. Set up graceful shutdown
// 10. Wait for termination signal and perform graceful shutdown
//
// The function follows a "fail fast" approach, exiting immediately if any
// critical initialization step fails. This ensures that the application
// doesn't start in a partially initialized state that could lead to
// unpredictable behavior.
func main() {
	// Initialize a basic logger for startup
	basicLogger := initBasicLogger()
	defer basicLogger.Sync()

	// Create a root context with cancellation
	rootCtx, rootCancel := context.WithCancel(context.Background())
	defer rootCancel() // Ensure all resources are cleaned up when main exits

	// Load configuration
	cfg, err := loadConfig(basicLogger)
	if err != nil {
		os.Exit(1)
	}

	// Initialize the main logger
	logger, err := initLogger(cfg, basicLogger)
	if err != nil {
		os.Exit(1)
	}
	defer logger.Sync()

	// From this point on, use the configured logger instead of the basic logger

	// Initialize dependency injection container
	container, err := initContainer(rootCtx, logger, cfg)
	if err != nil {
		os.Exit(1)
	}
	defer func() {
		if err := container.Close(); err != nil {
			logger.Error("Error closing container", zap.Error(err))
		}
	}()

	// Set up HTTP routes
	handler, telemetryShutdown, err := setupRoutes(rootCtx, container, logger, cfg)
	if err != nil {
		logger.Error("Failed to set up HTTP routes", zap.Error(err))
		os.Exit(1)
	}

	// Add telemetry shutdown to the shutdown process
	defer telemetryShutdown()

	// Apply auth middleware to all routes
	// Note: The auth middleware will validate JWT tokens locally.
	// In the future, this should be configured to use a remote authorization server
	// for improved security and centralized management.
	handler = container.GetAuthService().Middleware()(handler)

	// Start the server
	srv := startServer(handler, cfg, logger, container.GetContextLogger())

	// Set up graceful shutdown
	shutdownFunc := setupGracefulShutdown(rootCtx, rootCancel, srv, cfg)

	// Wait for a shutdown signal
	logger.Info("HTTP server is running. Press Ctrl+C to stop")
	if err := shutdown.GracefulShutdown(rootCtx, container.GetContextLogger(), shutdownFunc); err != nil {
		logger.Error("Failed to gracefully shutdown the HTTP server", zap.Error(err))
		os.Exit(1)
	}

	logger.Info("HTTP server shutdown complete")
}
