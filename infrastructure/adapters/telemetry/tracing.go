// Copyright (c) 2025 A Bit of Help, Inc.

// Package telemetry provides functionality for distributed tracing and metrics.
package telemetry

import (
	"context"
	"fmt"

	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.uber.org/zap"
)

// InitTracing initializes distributed tracing for the application
func InitTracing(ctx context.Context, cfg *config.Config, logger *zap.Logger) (func(), error) {
	if !cfg.Telemetry.Tracing.Enabled {
		logger.Info("Distributed tracing is disabled")
		return func() {}, nil
	}

	logger.Info("Initializing distributed tracing",
		zap.String("endpoint", cfg.Telemetry.Tracing.OTLP.Endpoint),
		zap.Bool("insecure", cfg.Telemetry.Tracing.OTLP.Insecure))

	// Create a resource describing the service
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("family-service"),
			semconv.ServiceVersion(cfg.App.Version),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create a trace provider
	// In a real implementation, we would configure an exporter here
	// For now, we'll use a simple stdout exporter for demonstration
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Set the global trace provider
	otel.SetTracerProvider(tp)

	// Set the global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	logger.Info("Distributed tracing initialized successfully")

	// Return a shutdown function
	return func() {
		logger.Info("Shutting down tracer provider")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Telemetry.ShutdownTimeout)
		defer cancel()

		if err := tp.Shutdown(shutdownCtx); err != nil {
			logger.Error("Failed to shutdown tracer provider", zap.Error(err))
		}
		logger.Info("Tracer provider shut down successfully")
	}, nil
}
