// Copyright (c) 2025 A Bit of Help, Inc.

package circuit

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestNewCircuitBreaker(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Test cases
	tests := []struct {
		name     string
		config   *config.CircuitConfig
		expected bool
	}{
		{
			name: "Enabled circuit breaker",
			config: &config.CircuitConfig{
				Enabled:         true,
				Timeout:         5 * time.Second,
				MaxConcurrent:   100,
				ErrorThreshold:  0.5,
				VolumeThreshold: 20,
				SleepWindow:     10 * time.Second,
			},
			expected: true,
		},
		{
			name: "Disabled circuit breaker",
			config: &config.CircuitConfig{
				Enabled:         false,
				Timeout:         5 * time.Second,
				MaxConcurrent:   100,
				ErrorThreshold:  0.5,
				VolumeThreshold: 20,
				SleepWindow:     10 * time.Second,
			},
			expected: false,
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cb := NewCircuitBreaker("test", tc.config, logger)

			if tc.expected {
				assert.NotNil(t, cb)
				assert.Equal(t, Closed, cb.GetState())
			} else {
				assert.Nil(t, cb)
			}
		})
	}
}

func TestCircuitBreaker_Execute_Success(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create a circuit breaker
	cfg := &config.CircuitConfig{
		Enabled:         true,
		Timeout:         5 * time.Second,
		MaxConcurrent:   100,
		ErrorThreshold:  0.5,
		VolumeThreshold: 20,
		SleepWindow:     10 * time.Second,
	}
	cb := NewCircuitBreaker("test", cfg, logger)
	require.NotNil(t, cb)

	// Test successful execution
	ctx := context.Background()
	callCount := 0

	// Function that should succeed
	fn := func(ctx context.Context) error {
		callCount++
		return nil
	}

	// Execute the function
	err := cb.Execute(ctx, "test-operation", fn)
	assert.NoError(t, err)
	assert.Equal(t, 1, callCount)
	assert.Equal(t, Closed, cb.GetState())
}

func TestCircuitBreaker_Execute_Failure(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create a circuit breaker with a low volume threshold for testing
	cfg := &config.CircuitConfig{
		Enabled:         true,
		Timeout:         5 * time.Second,
		MaxConcurrent:   100,
		ErrorThreshold:  0.5,
		VolumeThreshold: 2, // Low threshold for testing
		SleepWindow:     1 * time.Second,
	}
	cb := NewCircuitBreaker("test", cfg, logger)
	require.NotNil(t, cb)

	// Test failing execution
	ctx := context.Background()
	testErr := errors.New("test error")

	// Function that should fail
	fn := func(ctx context.Context) error {
		return testErr
	}

	// Execute the function multiple times to trip the circuit
	for i := 0; i < 2; i++ {
		err := cb.Execute(ctx, "test-operation", fn)
		assert.Error(t, err)
		assert.Equal(t, testErr, err)
	}

	// The third execution should trip the circuit
	err := cb.Execute(ctx, "test-operation", fn)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit breaker test is open")

	// Circuit should be open now
	assert.Equal(t, Open, cb.GetState())

	// Try to execute again, should fail with circuit open error
	err = cb.Execute(ctx, "test-operation", func(ctx context.Context) error {
		t.Fatal("This function should not be called when circuit is open")
		return nil
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "circuit breaker test is open")
}

func TestCircuitBreaker_ExecuteWithFallback(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create a circuit breaker
	cfg := &config.CircuitConfig{
		Enabled:         true,
		Timeout:         5 * time.Second,
		MaxConcurrent:   100,
		ErrorThreshold:  0.5,
		VolumeThreshold: 2, // Low threshold for testing
		SleepWindow:     1 * time.Second,
	}
	cb := NewCircuitBreaker("test", cfg, logger)
	require.NotNil(t, cb)

	// Test with fallback
	ctx := context.Background()
	testErr := errors.New("test error")
	fallbackCalled := false

	// Function that should fail
	fn := func(ctx context.Context) error {
		return testErr
	}

	// Fallback function
	fallback := func(ctx context.Context, err error) error {
		fallbackCalled = true
		assert.Equal(t, testErr, err)
		return nil
	}

	// Execute with fallback
	err := cb.ExecuteWithFallback(ctx, "test-operation", fn, fallback)
	assert.NoError(t, err)
	assert.True(t, fallbackCalled)
}

func TestCircuitBreaker_HalfOpen(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create a circuit breaker with a short sleep window for testing
	cfg := &config.CircuitConfig{
		Enabled:         true,
		Timeout:         5 * time.Second,
		MaxConcurrent:   100,
		ErrorThreshold:  0.5,
		VolumeThreshold: 2, // Low threshold for testing
		SleepWindow:     100 * time.Millisecond, // Short sleep window for testing
	}
	cb := NewCircuitBreaker("test", cfg, logger)
	require.NotNil(t, cb)

	// Trip the circuit
	ctx := context.Background()
	testErr := errors.New("test error")

	// Function that should fail
	fn := func(ctx context.Context) error {
		return testErr
	}

	// Execute the function multiple times to trip the circuit
	for i := 0; i < 3; i++ {
		_ = cb.Execute(ctx, "test-operation", fn)
	}

	// Circuit should be open now
	assert.Equal(t, Open, cb.GetState())

	// Wait for sleep window to elapse
	time.Sleep(200 * time.Millisecond)

	// Next request should put the circuit in half-open state
	successCalled := false
	err := cb.Execute(ctx, "test-operation", func(ctx context.Context) error {
		successCalled = true
		return nil
	})
	assert.NoError(t, err)
	assert.True(t, successCalled)

	// Circuit should be closed now after successful request in half-open state
	assert.Equal(t, Closed, cb.GetState())
}

func TestCircuitBreaker_Reset(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create a circuit breaker
	cfg := &config.CircuitConfig{
		Enabled:         true,
		Timeout:         5 * time.Second,
		MaxConcurrent:   100,
		ErrorThreshold:  0.5,
		VolumeThreshold: 2, // Low threshold for testing
		SleepWindow:     10 * time.Second,
	}
	cb := NewCircuitBreaker("test", cfg, logger)
	require.NotNil(t, cb)

	// Trip the circuit
	ctx := context.Background()
	testErr := errors.New("test error")

	// Function that should fail
	fn := func(ctx context.Context) error {
		return testErr
	}

	// Execute the function multiple times to trip the circuit
	for i := 0; i < 3; i++ {
		_ = cb.Execute(ctx, "test-operation", fn)
	}

	// Circuit should be open now
	assert.Equal(t, Open, cb.GetState())

	// Reset the circuit
	cb.Reset()

	// Circuit should be closed now
	assert.Equal(t, Closed, cb.GetState())
}

func TestCircuitBreaker_NilSafety(t *testing.T) {
	// Test that nil circuit breaker operations don't panic
	var cb *CircuitBreaker

	assert.NotPanics(t, func() {
		_ = cb.Execute(context.Background(), "test-operation", func(ctx context.Context) error {
			return nil
		})
	})

	assert.NotPanics(t, func() {
		_ = cb.ExecuteWithFallback(context.Background(), "test-operation", func(ctx context.Context) error {
			return errors.New("test error")
		}, func(ctx context.Context, err error) error {
			return nil
		})
	})

	assert.NotPanics(t, func() {
		_ = cb.GetState()
	})

	assert.NotPanics(t, func() {
		cb.Reset()
	})
}
