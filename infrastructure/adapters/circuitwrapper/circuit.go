// Copyright (c) 2025 A Bit of Help, Inc.

// Package circuit provides functionality for circuit breaking on external dependencies.
// This package is a wrapper around the servicelib circuit package.
package circuit

import (
	"context"
	"fmt"

	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	"github.com/abitofhelp/servicelib/circuit"
	"github.com/abitofhelp/servicelib/errors/recovery"
	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

// State represents the state of the circuit breaker
type State int

const (
	// Closed means the circuit is closed and requests are allowed through
	Closed State = iota
	// Open means the circuit is open and requests are not allowed through
	Open
	// HalfOpen means the circuit is allowing a limited number of requests through to test if the dependency is healthy
	HalfOpen
)

// CircuitBreaker implements the circuit breaker pattern to protect against
// cascading failures when external dependencies are unavailable.
type CircuitBreaker struct {
	name   string
	cb     *circuit.CircuitBreaker
	logger *zap.Logger
}

// NewCircuitBreaker creates a new circuit breaker
func NewCircuitBreaker(name string, cfg *config.CircuitConfig, logger *zap.Logger) *CircuitBreaker {
	if !cfg.Enabled {
		logger.Info("Circuit breaker is disabled", zap.String("name", name))
		return nil
	}

	logger.Info("Initializing circuit breaker",
		zap.String("name", name),
		zap.Duration("timeout", cfg.Timeout),
		zap.Int("max_concurrent", cfg.MaxConcurrent),
		zap.Float64("error_threshold", cfg.ErrorThreshold),
		zap.Int("volume_threshold", cfg.VolumeThreshold),
		zap.Duration("sleep_window", cfg.SleepWindow))

	// Create servicelib circuit configuration
	circuitConfig := circuit.DefaultConfig().
		WithEnabled(cfg.Enabled).
		WithTimeout(cfg.Timeout).
		WithMaxConcurrent(cfg.MaxConcurrent).
		WithErrorThreshold(cfg.ErrorThreshold).
		WithVolumeThreshold(cfg.VolumeThreshold).
		WithSleepWindow(cfg.SleepWindow)

	// Create servicelib circuit options
	contextLogger := logging.NewContextLogger(logger)
	options := circuit.DefaultOptions().
		WithLogger(contextLogger).
		WithName(name)

	// Create servicelib circuit breaker
	serviceCB := circuit.NewCircuitBreaker(circuitConfig, options)

	return &CircuitBreaker{
		name:   name,
		cb:     serviceCB,
		logger: logger,
	}
}

// Execute executes the given function with circuit breaking
// If the circuit is open, it will return an error immediately
// If the circuit is closed or half-open, it will execute the function
// and update the circuit state based on the result
func (cb *CircuitBreaker) Execute(ctx context.Context, operation string, fn func(ctx context.Context) error) error {
	if cb == nil || cb.cb == nil {
		// If circuit breaker is disabled, just execute the function
		return fn(ctx)
	}

	// Create a wrapper function that adapts to servicelib's circuit.Execute
	wrapper := func(ctx context.Context) (interface{}, error) {
		err := fn(ctx)
		return nil, err
	}

	// Execute the function with circuit breaking
	_, err := circuit.Execute(ctx, cb.cb, operation, wrapper)

	// Convert servicelib's circuit breaker error to our format if needed
	if err == recovery.ErrCircuitBreakerOpen {
		return fmt.Errorf("circuit breaker %s is open", cb.name)
	}

	return err
}

// ExecuteWithFallback executes the given function with circuit breaking
// If the circuit is open or the function fails, it will execute the fallback function
func (cb *CircuitBreaker) ExecuteWithFallback(ctx context.Context, operation string, fn func(ctx context.Context) error, fallback func(ctx context.Context, err error) error) error {
	if cb == nil || cb.cb == nil {
		// If circuit breaker is disabled, just execute the function
		err := fn(ctx)
		if err != nil {
			return fallback(ctx, err)
		}
		return nil
	}

	// Create wrapper functions that adapt to servicelib's circuit.ExecuteWithFallback
	wrapper := func(ctx context.Context) (interface{}, error) {
		err := fn(ctx)
		return nil, err
	}

	fallbackWrapper := func(ctx context.Context, err error) (interface{}, error) {
		// Convert servicelib's circuit breaker error to our format if needed
		if err == recovery.ErrCircuitBreakerOpen {
			err = fmt.Errorf("circuit breaker %s is open", cb.name)
		}

		fallbackErr := fallback(ctx, err)
		return nil, fallbackErr
	}

	// Execute the function with circuit breaking and fallback
	_, err := circuit.ExecuteWithFallback(ctx, cb.cb, operation, wrapper, fallbackWrapper)
	return err
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() State {
	if cb == nil || cb.cb == nil {
		return Closed
	}

	// Convert servicelib's circuit state to our state
	switch cb.cb.GetState() {
	case circuit.Closed:
		return Closed
	case circuit.Open:
		return Open
	case circuit.HalfOpen:
		return HalfOpen
	default:
		return Closed
	}
}

// Reset resets the circuit breaker to its initial state
func (cb *CircuitBreaker) Reset() {
	if cb == nil || cb.cb == nil {
		return
	}

	cb.cb.Reset()
}

// Execute is a package-level function that executes the given function with circuit breaking
// It's a wrapper around the servicelib circuit.Execute function
// This function is used by the repository implementations
func Execute(ctx context.Context, cb *CircuitBreaker, operation string, fn func(ctx context.Context) (bool, error)) (bool, error) {
	if cb == nil || cb.cb == nil {
		// If circuit breaker is disabled, just execute the function
		return fn(ctx)
	}

	// Execute the function with circuit breaking using the servicelib circuit.Execute
	return circuit.Execute(ctx, cb.cb, operation, fn)
}
