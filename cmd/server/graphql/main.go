// Copyright (c) 2025 A Bit of Help, Inc.

// Package main is the entry point for the GraphQL server.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/abitofhelp/family-service/cmd/server/graphql/di"
	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
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

// initBasicLogger creates a basic logger for use during startup
// before the main logger is configured.
func initBasicLogger() *zap.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		// If we can't create a logger, we're in serious trouble
		fmt.Printf("Failed to create basic logger: %v\n", err)
		os.Exit(1)
	}
	return logger
}

// loadConfig loads the application configuration.
// It returns the configuration and any error that occurred.
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

// initLogger initializes the application logger based on configuration.
// It returns the logger and any error that occurred.
func initLogger(cfg *config.Config, basicLogger *zap.Logger) (*zap.Logger, error) {
	basicLogger.Info("Initializing application logger",
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

// initContainer initializes the dependency injection container.
// It returns the container and any error that occurred.
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

// setupRoutes sets up the HTTP routes for the application.
// It returns the HTTP handler with all routes configured.
func setupRoutes(container *di.Container, logger *zap.Logger, cfg *config.Config) http.Handler {
	logger.Info("Setting up HTTP routes")

	// Create a container adapter for health checks
	adapter := health.NewGenericContainerAdapter(container)

	// Create a ServeMux for routing
	mux := http.NewServeMux()

	// Set up metrics endpoint if enabled
	setupMetricsEndpoint(mux, cfg, logger)

	// Set up GraphQL endpoints
	setupGraphQLEndpoints(mux, container)

	// Health check endpoint
	healthEndpoint := cfg.Server.HealthEndpoint
	logger.Info("Setting up health check endpoint", zap.String("endpoint", healthEndpoint))
	configAdapter := pkgconfig.NewGenericConfigAdapter(cfg)
	mux.Handle(healthEndpoint, health.NewHandler(adapter, logger, configAdapter))

	logger.Info("HTTP routes set up successfully")
	return mux
}

// setupMetricsEndpoint sets up the metrics endpoint if enabled in configuration.
func setupMetricsEndpoint(mux *http.ServeMux, cfg *config.Config, logger *zap.Logger) {
	if cfg.Telemetry.Exporters.Metrics.Prometheus.Enabled {
		metricsPath := cfg.Telemetry.Exporters.Metrics.Prometheus.Path
		logger.Info("Setting up Prometheus metrics endpoint",
			zap.String("path", metricsPath),
			zap.String("listen", cfg.Telemetry.Exporters.Metrics.Prometheus.Listen))
		mux.Handle(metricsPath, telemetry.CreatePrometheusHandler())
	} else {
		logger.Info("Prometheus metrics endpoint is disabled")
	}
}

// setupGraphQLEndpoints sets up all GraphQL-related endpoints.
func setupGraphQLEndpoints(mux *http.ServeMux, container *di.Container) {
	// Get the resolver
	resolverInstance := resolver.NewResolver(container.GetFamilyApplicationService(), container.GetContextLogger())

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

// startServer creates and starts the HTTP server.
// It returns the server and any error that occurred.
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

// setupGracefulShutdown sets up graceful shutdown for the server.
// It returns a function that will be called to shut down the server.
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
	handler := setupRoutes(container, logger, cfg)

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
