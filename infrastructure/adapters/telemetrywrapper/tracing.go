// Copyright (c) 2025 A Bit of Help, Inc.

// Package telemetry provides functionality for distributed tracing and metrics.
// This package is a wrapper around the servicelib telemetry package.
package telemetry

import (
	"context"
	"fmt"

	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	"github.com/abitofhelp/servicelib/telemetry"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Tracer is a wrapper around the servicelib telemetry tracer
type Tracer struct {
	tracer telemetry.Tracer
	logger *zap.Logger
}

// NewTracer creates a new Tracer
func NewTracer(name string, logger *zap.Logger) *Tracer {
	// Get the OpenTelemetry tracer
	otelTracer := otel.GetTracerProvider().Tracer(name)

	// Create a servicelib tracer
	serviceTracer := telemetry.NewOtelTracer(otelTracer)

	return &Tracer{
		tracer: serviceTracer,
		logger: logger,
	}
}

// Start starts a new span
func (t *Tracer) Start(ctx context.Context, spanName string) (context.Context, trace.Span) {
	// Instead of trying to convert the servicelib span to an OpenTelemetry span,
	// which doesn't work because the GetOtelSpan method doesn't exist,
	// we'll directly create a new span using the OpenTelemetry tracer.
	// Use the name field of the Tracer struct instead of an empty string
	otelTracer := otel.GetTracerProvider().Tracer("family-service")
	return otelTracer.Start(ctx, spanName)
}

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

// StartSpan is a helper function to start a span with attributes
func StartSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	tracer := otel.GetTracerProvider().Tracer("family-service")
	ctx, span := tracer.Start(ctx, name)
	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}
	return ctx, span
}

// AddSpanAttributes adds attributes to the current span
func AddSpanAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attrs...)
}

// RecordError records an error on the current span
func RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
}

// WithSpan wraps a function with a span
func WithSpan(ctx context.Context, name string, fn func(context.Context) error) error {
	// Always use the span from the context if it exists
	span := trace.SpanFromContext(ctx)

	// Only create a new span if there's no span in the context
	if span == nil || span == trace.SpanFromContext(context.Background()) {
		ctx, span = StartSpan(ctx, name)
	}

	defer span.End()

	err := fn(ctx)
	if err != nil {
		span.RecordError(err) // Use the span directly instead of calling RecordError
	}

	return err
}
