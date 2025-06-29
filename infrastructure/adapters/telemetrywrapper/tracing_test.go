// Copyright (c) 2025 A Bit of Help, Inc.

package telemetry

import (
	"context"
	"errors"
	"testing"

	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap/zaptest"
)

// setupTestTracer sets up a tracer provider for testing
func setupTestTracer(t *testing.T) func() {
	// Create a resource describing the test
	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName("test-service"),
			semconv.ServiceVersion("test-version"),
		),
	)
	require.NoError(t, err)

	// Create a trace provider with a sampler that always samples
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

	// Return a cleanup function
	return func() {
		_ = tp.Shutdown(context.Background())
	}
}

func TestNewTracer(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create a tracer
	tracer := NewTracer("test", logger)
	require.NotNil(t, tracer)
	assert.NotNil(t, tracer.tracer)
	assert.NotNil(t, tracer.logger)
}

func TestTracer_Start(t *testing.T) {
	// Set up a tracer provider for testing
	cleanup := setupTestTracer(t)
	defer cleanup()

	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create a tracer
	tracer := NewTracer("test", logger)
	require.NotNil(t, tracer)

	// Start a span
	ctx := context.Background()
	spanName := "test-span"
	ctx, span := tracer.Start(ctx, spanName)
	require.NotNil(t, span)

	// Verify the span is valid
	assert.NotEqual(t, trace.SpanID{}, span.SpanContext().SpanID())

	// End the span
	span.End()
}

func TestInitTracing(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Test cases
	tests := []struct {
		name     string
		config   *config.TelemetryConfig
		expected bool
	}{
		{
			name: "Enabled tracing",
			config: &config.TelemetryConfig{
				Tracing: config.TracingConfig{
					Enabled: true,
					OTLP: config.OTLPConfig{
						Endpoint: "localhost:4317",
						Insecure: true,
					},
				},
				ShutdownTimeout: 5,
			},
			expected: true,
		},
		{
			name: "Disabled tracing",
			config: &config.TelemetryConfig{
				Tracing: config.TracingConfig{
					Enabled: false,
				},
				ShutdownTimeout: 5,
			},
			expected: false,
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			shutdown, err := InitTracing(ctx, &config.Config{
				Telemetry: *tc.config,
				App: config.AppConfig{
					Version: "1.1.0",
				},
			}, logger)
			require.NoError(t, err)
			require.NotNil(t, shutdown)

			// Call the shutdown function
			shutdown()
		})
	}
}

func TestStartSpan(t *testing.T) {
	// Set up a tracer provider for testing
	cleanup := setupTestTracer(t)
	defer cleanup()

	// Start a span
	ctx := context.Background()
	spanName := "test-span"
	attrs := []attribute.KeyValue{
		attribute.String("key1", "value1"),
		attribute.Int("key2", 42),
	}
	ctx, span := StartSpan(ctx, spanName, attrs...)
	require.NotNil(t, span)

	// Verify the span is valid
	assert.NotEqual(t, trace.SpanID{}, span.SpanContext().SpanID())

	// End the span
	span.End()
}

func TestAddSpanAttributes(t *testing.T) {
	// Create a span
	ctx := context.Background()
	spanName := "test-span"
	ctx, span := StartSpan(ctx, spanName)
	require.NotNil(t, span)

	// Add attributes
	attrs := []attribute.KeyValue{
		attribute.String("key1", "value1"),
		attribute.Int("key2", 42),
	}
	AddSpanAttributes(ctx, attrs...)

	// End the span
	span.End()
}

func TestRecordError(t *testing.T) {
	// Create a span
	ctx := context.Background()
	spanName := "test-span"
	ctx, span := StartSpan(ctx, spanName)
	require.NotNil(t, span)

	// Record an error
	err := errors.New("test error")
	RecordError(ctx, err)

	// End the span
	span.End()
}

func TestWithSpan(t *testing.T) {
	// Test successful execution
	ctx := context.Background()
	spanName := "test-span"
	callCount := 0

	// Function that should succeed
	fn := func(ctx context.Context) error {
		callCount++
		return nil
	}

	// Execute with span
	err := WithSpan(ctx, spanName, fn)
	assert.NoError(t, err)
	assert.Equal(t, 1, callCount)

	// Test error handling
	testErr := errors.New("test error")
	err = WithSpan(ctx, spanName, func(ctx context.Context) error {
		return testErr
	})
	assert.Error(t, err)
	assert.Equal(t, testErr, err)
}

// MockSpan is a mock implementation of trace.Span for testing
type MockSpan struct {
	trace.Span
	ended       bool
	attributes  []attribute.KeyValue
	recordedErr error
}

func (s *MockSpan) End(options ...trace.SpanEndOption) {
	s.ended = true
}

func (s *MockSpan) SetAttributes(attrs ...attribute.KeyValue) {
	s.attributes = append(s.attributes, attrs...)
}

func (s *MockSpan) RecordError(err error, opts ...trace.EventOption) {
	s.recordedErr = err
}

func TestWithSpan_MockSpan(t *testing.T) {
	// Set up a tracer provider for testing
	cleanup := setupTestTracer(t)
	defer cleanup()

	// Create a mock span
	mockSpan := &MockSpan{Span: noop.Span{}}
	ctx := trace.ContextWithSpan(context.Background(), mockSpan)

	// Test successful execution
	err := WithSpan(ctx, "test-span", func(ctx context.Context) error {
		return nil
	})
	assert.NoError(t, err)
	assert.True(t, mockSpan.ended)

	// Test error handling
	mockSpan = &MockSpan{Span: noop.Span{}}
	ctx = trace.ContextWithSpan(context.Background(), mockSpan)
	testErr := errors.New("test error")
	err = WithSpan(ctx, "test-span", func(ctx context.Context) error {
		return testErr
	})
	assert.Error(t, err)
	assert.Equal(t, testErr, err)
	assert.Equal(t, testErr, mockSpan.recordedErr)
	assert.True(t, mockSpan.ended)
}
