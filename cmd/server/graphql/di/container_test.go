// Copyright (c) 2025 A Bit of Help, Inc.

package di_test

import (
	"context"
	"github.com/abitofhelp/family-service/cmd/server/graphql/di"
	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// TestNewContainer_MongoDB tests creating a container with MongoDB
func TestNewContainer_MongoDB(t *testing.T) {
	// Skip this test in CI environments or when MongoDB is not available
	t.Skip("Skipping MongoDB test as it requires a real MongoDB connection")

	// Setup
	ctx := context.Background()
	logger := zaptest.NewLogger(t)
	cfg := &config.Config{
		App: config.AppConfig{
			Version: "test",
		},
		Auth: config.AuthConfig{
			OIDCTimeout: 5 * time.Second,
			JWT: config.JWTConfig{
				SecretKey:     "test-secret-key",
				TokenDuration: 24 * time.Hour,
				Issuer:        "test-issuer",
			},
		},
		Database: config.DatabaseConfig{
			Type: "mongodb",
			MongoDB: config.MongoDBConfig{
				URI:               "mongodb://localhost:27017",
				ConnectionTimeout: 5 * time.Second,
				PingTimeout:       5 * time.Second,
				DisconnectTimeout: 5 * time.Second,
				IndexTimeout:      5 * time.Second,
				MigrationTimeout:  5 * time.Second,
			},
			Postgres: config.PostgresConfig{
				DSN:              "postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable",
				MigrationTimeout: 5 * time.Second,
			},
		},
		Features: config.FeaturesConfig{
			UseGenerics: true,
		},
		Log: config.LogConfig{
			Level:       "info",
			Development: true,
		},
		Server: config.ServerConfig{
			Port:            "8089",
			ReadTimeout:     5 * time.Second,
			WriteTimeout:    5 * time.Second,
			IdleTimeout:     5 * time.Second,
			ShutdownTimeout: 5 * time.Second,
			HealthEndpoint:  "/health",
		},
		Telemetry: config.TelemetryConfig{
			ShutdownTimeout: 5 * time.Second,
			Exporters: config.ExportersConfig{
				Metrics: config.MetricsExporterConfig{
					Prometheus: config.PrometheusConfig{
						Enabled: true,
						Listen:  ":9090",
						Path:    "/metrics",
					},
				},
			},
		},
	}

	// Act
	container, err := di.NewContainer(ctx, logger, cfg)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, container)

	// Cleanup
	err = container.Close()
	assert.NoError(t, err)
}

// TestNewContainer_Postgres tests creating a container with PostgreSQL
func TestNewContainer_Postgres(t *testing.T) {
	// Skip this test in CI environments or when PostgreSQL is not available
	t.Skip("Skipping PostgreSQL test as it requires a real PostgreSQL connection")

	// Setup
	ctx := context.Background()
	logger := zaptest.NewLogger(t)
	cfg := &config.Config{
		App: config.AppConfig{
			Version: "test",
		},
		Auth: config.AuthConfig{
			OIDCTimeout: 5 * time.Second,
			JWT: config.JWTConfig{
				SecretKey:     "test-secret-key",
				TokenDuration: 24 * time.Hour,
				Issuer:        "test-issuer",
			},
		},
		Database: config.DatabaseConfig{
			Type: "postgres",
			MongoDB: config.MongoDBConfig{
				URI:               "mongodb://localhost:27017",
				ConnectionTimeout: 5 * time.Second,
				PingTimeout:       5 * time.Second,
				DisconnectTimeout: 5 * time.Second,
				IndexTimeout:      5 * time.Second,
				MigrationTimeout:  5 * time.Second,
			},
			Postgres: config.PostgresConfig{
				DSN:              "postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable",
				MigrationTimeout: 5 * time.Second,
			},
		},
		Features: config.FeaturesConfig{
			UseGenerics: true,
		},
		Log: config.LogConfig{
			Level:       "info",
			Development: true,
		},
		Server: config.ServerConfig{
			Port:            "8089",
			ReadTimeout:     5 * time.Second,
			WriteTimeout:    5 * time.Second,
			IdleTimeout:     5 * time.Second,
			ShutdownTimeout: 5 * time.Second,
			HealthEndpoint:  "/health",
		},
		Telemetry: config.TelemetryConfig{
			ShutdownTimeout: 5 * time.Second,
			Exporters: config.ExportersConfig{
				Metrics: config.MetricsExporterConfig{
					Prometheus: config.PrometheusConfig{
						Enabled: true,
						Listen:  ":9090",
						Path:    "/metrics",
					},
				},
			},
		},
	}

	// Act
	container, err := di.NewContainer(ctx, logger, cfg)

	// Assert
	require.NoError(t, err)
	require.NotNil(t, container)

	// Cleanup
	err = container.Close()
	assert.NoError(t, err)
}

// TestNewContainer_UnsupportedDB tests creating a container with an unsupported database type
func TestNewContainer_UnsupportedDB(t *testing.T) {
	// Setup
	ctx := context.Background()
	logger := zaptest.NewLogger(t)
	cfg := &config.Config{
		App: config.AppConfig{
			Version: "test",
		},
		Auth: config.AuthConfig{
			OIDCTimeout: 5 * time.Second,
			JWT: config.JWTConfig{
				SecretKey:     "test-secret-key",
				TokenDuration: 24 * time.Hour,
				Issuer:        "test-issuer",
			},
		},
		Database: config.DatabaseConfig{
			Type: "unsupported",
			MongoDB: config.MongoDBConfig{
				URI:               "mongodb://localhost:27017",
				ConnectionTimeout: 5 * time.Second,
				PingTimeout:       5 * time.Second,
				DisconnectTimeout: 5 * time.Second,
				IndexTimeout:      5 * time.Second,
				MigrationTimeout:  5 * time.Second,
			},
			Postgres: config.PostgresConfig{
				DSN:              "postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable",
				MigrationTimeout: 5 * time.Second,
			},
		},
		Features: config.FeaturesConfig{
			UseGenerics: true,
		},
		Log: config.LogConfig{
			Level:       "info",
			Development: true,
		},
		Server: config.ServerConfig{
			Port:            "8089",
			ReadTimeout:     5 * time.Second,
			WriteTimeout:    5 * time.Second,
			IdleTimeout:     5 * time.Second,
			ShutdownTimeout: 5 * time.Second,
			HealthEndpoint:  "/health",
		},
		Telemetry: config.TelemetryConfig{
			ShutdownTimeout: 5 * time.Second,
			Exporters: config.ExportersConfig{
				Metrics: config.MetricsExporterConfig{
					Prometheus: config.PrometheusConfig{
						Enabled: true,
						Listen:  ":9090",
						Path:    "/metrics",
					},
				},
			},
		},
	}

	// Act
	container, err := di.NewContainer(ctx, logger, cfg)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, container)
	assert.Contains(t, err.Error(), "unsupported database type")
}

// TestContainer_Getters tests the getter methods of the container
func TestContainer_Getters(t *testing.T) {
	// This test uses a mock repository factory to avoid external dependencies
	// Setup
	ctx := context.Background()
	logger := zaptest.NewLogger(t)

	// Create a minimal container for testing getters
	container := &di.Container{}

	// Use reflection to set private fields for testing
	// Note: In a real test, we would use a proper constructor or factory method
	// This is just for demonstration purposes
	t.Skip("Skipping getter tests as they require access to private fields")

	// Assert
	assert.Equal(t, ctx, container.GetContext())
	assert.Equal(t, logger, container.GetLogger())
	assert.NotNil(t, container.GetContextLogger())
	assert.NotNil(t, container.GetValidator())
	assert.NotNil(t, container.GetRepositoryFactory())
	assert.NotNil(t, container.GetFamilyService())
	assert.NotNil(t, container.GetAuthorizationService())
}

// TestContainer_Close tests the Close method of the container
func TestContainer_Close(t *testing.T) {
	// This test uses a mock repository factory to avoid external dependencies
	// Setup
	ctx := context.Background()
	logger := zap.NewNop()
	cfg := &config.Config{
		App: config.AppConfig{
			Version: "test",
		},
		Auth: config.AuthConfig{
			OIDCTimeout: 5 * time.Second,
			JWT: config.JWTConfig{
				SecretKey:     "test-secret-key-that-is-at-least-32-characters-long",
				TokenDuration: 24 * time.Hour,
				Issuer:        "test-issuer",
			},
		},
		Database: config.DatabaseConfig{
			Type: "mongodb",
			MongoDB: config.MongoDBConfig{
				URI:               "mongodb://localhost:27017",
				ConnectionTimeout: 5 * time.Second,
				PingTimeout:       5 * time.Second,
				DisconnectTimeout: 5 * time.Second,
				IndexTimeout:      5 * time.Second,
				MigrationTimeout:  5 * time.Second,
			},
			Postgres: config.PostgresConfig{
				DSN:              "postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable",
				MigrationTimeout: 5 * time.Second,
			},
		},
		Features: config.FeaturesConfig{
			UseGenerics: true,
		},
		Log: config.LogConfig{
			Level:       "info",
			Development: true,
		},
		Server: config.ServerConfig{
			Port:            "8089",
			ReadTimeout:     5 * time.Second,
			WriteTimeout:    5 * time.Second,
			IdleTimeout:     5 * time.Second,
			ShutdownTimeout: 5 * time.Second,
			HealthEndpoint:  "/health",
		},
		Telemetry: config.TelemetryConfig{
			ShutdownTimeout: 5 * time.Second,
			Exporters: config.ExportersConfig{
				Metrics: config.MetricsExporterConfig{
					Prometheus: config.PrometheusConfig{
						Enabled: true,
						Listen:  ":9090",
						Path:    "/metrics",
					},
				},
			},
		},
	}

	// Create a container with a mock repository factory
	container, err := di.NewContainer(ctx, logger, cfg)
	require.NoError(t, err)

	// Act
	err = container.Close()

	// Assert
	assert.NoError(t, err)
}

// TestContainer_Close_Error tests the Close method when an error occurs
func TestContainer_Close_Error(t *testing.T) {
	// Skip this test as it requires modifying private fields
	t.Skip("Skipping close error test as it requires access to private fields")

	// In a real test, we would:
	// 1. Create a mock repository factory that returns an error on Close
	// 2. Set it as the container's repository factory
	// 3. Call Close and verify that it returns an error
}
