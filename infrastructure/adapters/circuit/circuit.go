// Copyright (c) 2025 A Bit of Help, Inc.

// Package circuit provides functionality for circuit breaking on external dependencies.
package circuit

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	"go.uber.org/zap"
)

// CircuitBreaker implements the circuit breaker pattern to protect against
// cascading failures when external dependencies are unavailable.
type CircuitBreaker struct {
	name            string
	config          *config.CircuitConfig
	logger          *zap.Logger
	state           State
	failureCount    int
	lastFailureTime time.Time
	mutex           sync.RWMutex
}

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

	return &CircuitBreaker{
		name:            name,
		config:          cfg,
		logger:          logger,
		state:           Closed,
		failureCount:    0,
		lastFailureTime: time.Time{},
		mutex:           sync.RWMutex{},
	}
}

// Execute executes the given function with circuit breaking
// If the circuit is open, it will return an error immediately
// If the circuit is closed or half-open, it will execute the function
// and update the circuit state based on the result
func (cb *CircuitBreaker) Execute(ctx context.Context, operation string, fn func(ctx context.Context) error) error {
	if cb == nil {
		// If circuit breaker is disabled, just execute the function
		return fn(ctx)
	}

	// Check if the circuit is open
	if !cb.allowRequest() {
		cb.logger.Warn("Circuit is open, rejecting request",
			zap.String("circuit", cb.name),
			zap.String("operation", operation))
		return fmt.Errorf("circuit breaker %s is open", cb.name)
	}

	// Execute the function
	err := fn(ctx)

	// Update the circuit state based on the result
	cb.updateState(err)

	return err
}

// ExecuteWithFallback executes the given function with circuit breaking
// If the circuit is open or the function fails, it will execute the fallback function
func (cb *CircuitBreaker) ExecuteWithFallback(ctx context.Context, operation string, fn func(ctx context.Context) error, fallback func(ctx context.Context, err error) error) error {
	if cb == nil {
		// If circuit breaker is disabled, just execute the function
		err := fn(ctx)
		if err != nil {
			return fallback(ctx, err)
		}
		return nil
	}

	// Check if the circuit is open
	if !cb.allowRequest() {
		cb.logger.Warn("Circuit is open, using fallback",
			zap.String("circuit", cb.name),
			zap.String("operation", operation))
		return fallback(ctx, fmt.Errorf("circuit breaker %s is open", cb.name))
	}

	// Execute the function
	err := fn(ctx)

	// Update the circuit state based on the result
	cb.updateState(err)

	// If the function failed, execute the fallback
	if err != nil {
		return fallback(ctx, err)
	}

	return nil
}

// allowRequest checks if a request should be allowed through the circuit
func (cb *CircuitBreaker) allowRequest() bool {
	cb.mutex.RLock()
	defer cb.mutex.RUnlock()

	switch cb.state {
	case Closed:
		return true
	case Open:
		// Check if the sleep window has elapsed
		if time.Since(cb.lastFailureTime) > cb.config.SleepWindow {
			// Transition to half-open state
			cb.mutex.RUnlock()
			cb.mutex.Lock()
			cb.state = HalfOpen
			cb.mutex.Unlock()
			cb.mutex.RLock()
			return true
		}
		return false
	case HalfOpen:
		// In half-open state, allow a limited number of requests through
		return true
	default:
		return true
	}
}

// updateState updates the state of the circuit based on the result of a request
func (cb *CircuitBreaker) updateState(err error) {
	cb.mutex.Lock()
	defer cb.mutex.Unlock()

	if err != nil {
		// Request failed
		cb.failureCount++
		cb.lastFailureTime = time.Now()

		// Check if the circuit should be opened
		if cb.state == Closed && cb.failureCount >= cb.config.VolumeThreshold && float64(cb.failureCount)/float64(cb.config.VolumeThreshold) >= cb.config.ErrorThreshold {
			cb.logger.Warn("Opening circuit breaker due to error threshold exceeded",
				zap.String("circuit", cb.name),
				zap.Int("failure_count", cb.failureCount),
				zap.Int("volume_threshold", cb.config.VolumeThreshold),
				zap.Float64("error_threshold", cb.config.ErrorThreshold))
			cb.state = Open
		} else if cb.state == HalfOpen {
			// If a request fails in half-open state, go back to open
			cb.logger.Warn("Reopening circuit breaker due to failure in half-open state",
				zap.String("circuit", cb.name))
			cb.state = Open
		}
	} else {
		// Request succeeded
		if cb.state == HalfOpen {
			// If a request succeeds in half-open state, close the circuit
			cb.logger.Info("Closing circuit breaker after successful request in half-open state",
				zap.String("circuit", cb.name))
			cb.state = Closed
			cb.failureCount = 0
		} else if cb.state == Closed {
			// Reset failure count on successful request
			cb.failureCount = 0
		}
	}
}

// GetState returns the current state of the circuit breaker
func (cb *CircuitBreaker) GetState() State {
	if cb == nil {
		return Closed
	}

	cb.mutex.RLock()
	defer cb.mutex.RUnlock()
	return cb.state
}

// Reset resets the circuit breaker to its initial state
func (cb *CircuitBreaker) Reset() {
	if cb == nil {
		return
	}

	cb.mutex.Lock()
	defer cb.mutex.Unlock()
	cb.state = Closed
	cb.failureCount = 0
	cb.lastFailureTime = time.Time{}
}
