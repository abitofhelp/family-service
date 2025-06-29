// Copyright (c) 2025 A Bit of Help, Inc.

package repository

import (
	"context"
	"sync"
	"time"

	"github.com/abitofhelp/family-service/infrastructure/adapters/circuitwrapper"
	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	repoerrors "github.com/abitofhelp/family-service/infrastructure/adapters/errors"
	"github.com/abitofhelp/family-service/infrastructure/adapters/loggingwrapper"
	"github.com/abitofhelp/family-service/infrastructure/adapters/ratewrapper"
	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/retry"
)

var (
	// Global configuration instance
	globalConfig     *config.Config
	globalConfigOnce sync.Once
)

// SetGlobalConfig sets the global configuration instance
func SetGlobalConfig(cfg *config.Config) {
	globalConfigOnce.Do(func() {
		globalConfig = cfg
	})
}

// GetRetryConfig returns the retry configuration
func GetRetryConfig() retry.Config {
	if globalConfig != nil {
		return retry.DefaultConfig().
			WithMaxRetries(globalConfig.Retry.MaxRetries).
			WithInitialBackoff(globalConfig.Retry.InitialBackoff).
			WithMaxBackoff(globalConfig.Retry.MaxBackoff)
	}

	// Fallback to default values if configuration is not available
	return retry.DefaultConfig().
		WithMaxRetries(3).
		WithInitialBackoff(100 * time.Millisecond).
		WithMaxBackoff(1 * time.Second)
}

// BaseRepository provides common functionality for all repository implementations
// It uses the family-service's wrapper implementations for circuit and rate packages
// to ensure consistent usage of these features throughout the codebase.
type BaseRepository struct {
	Logger         *loggingwrapper.ContextLogger
	CircuitBreaker *circuit.CircuitBreaker
	RateLimiter    *rate.RateLimiter
	DefaultTimeout time.Duration
}

// NewBaseRepository creates a new BaseRepository with the given logger
func NewBaseRepository(logger *loggingwrapper.ContextLogger, dbName string) *BaseRepository {
	if logger == nil {
		panic("logger cannot be nil")
	}

	// Get circuit breaker configuration
	var circuitConfig *config.CircuitConfig
	if globalConfig != nil {
		circuitConfig = &globalConfig.Circuit
	} else {
		// Default configuration if global config is not available
		circuitConfig = &config.CircuitConfig{
			Enabled:         true,
			Timeout:         5 * time.Second,
			MaxConcurrent:   100,
			ErrorThreshold:  0.5,
			VolumeThreshold: 20,
			SleepWindow:     10 * time.Second,
		}
	}

	// Get rate limiter configuration
	var rateConfig *config.RateConfig
	if globalConfig != nil {
		rateConfig = &globalConfig.Rate
	} else {
		// Default configuration if global config is not available
		rateConfig = &config.RateConfig{
			Enabled:           true,
			RequestsPerSecond: 100,
			BurstSize:         50,
		}
	}

	// Create circuit breaker using the family-service wrapper
	cb := circuit.NewCircuitBreaker(dbName, circuitConfig, logger.Logger())

	// Create rate limiter using the family-service wrapper
	rl := rate.NewRateLimiter(dbName, rateConfig, logger.Logger())

	return &BaseRepository{
		Logger:         logger,
		CircuitBreaker: cb,
		RateLimiter:    rl,
		DefaultTimeout: 5 * time.Second,
	}
}

// IsRetryableError determines if an error should be retried
func IsRetryableError(err error) bool {
	// Don't retry not found errors
	if _, ok := err.(*errors.NotFoundError); ok {
		return false
	}

	// Don't retry validation errors
	if _, ok := err.(*errors.ValidationError); ok {
		return false
	}

	// Retry network errors, timeouts, and transient database errors
	return retry.IsNetworkError(err) || retry.IsTimeoutError(err) || retry.IsTransientError(err)
}

// ExecuteWithResilience executes an operation with retry, circuit breaker, and rate limiter
// It uses the family-service's wrapper implementations for circuit and rate packages
// to ensure consistent usage of these features throughout the codebase.
func (r *BaseRepository) ExecuteWithResilience(
	ctx context.Context,
	operation func(context.Context) error,
	operationName string,
) error {
	// Create a context with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, r.DefaultTimeout)
	defer cancel()

	var retryErr error

	// Configure retry with backoff
	retryConfig := GetRetryConfig()

	// Wrap the retry operation with circuit breaker
	circuitOperation := func(ctx context.Context) error {
		// Execute with retry
		retryErr = retry.Do(ctx, operation, retryConfig, IsRetryableError)
		return retryErr
	}

	// Wrap the circuit breaker operation with rate limiter
	rateOperation := func(ctx context.Context) error {
		// Execute with circuit breaker
		return r.CircuitBreaker.Execute(ctx, operationName, circuitOperation)
	}

	// Execute with rate limiter
	err := r.RateLimiter.Execute(ctxWithTimeout, operationName, rateOperation)

	// Check for errors from rate limiter or circuit breaker
	if err != nil && retryErr == nil {
		return err
	}

	// Return the retry error if there is one
	return retryErr
}

// HandleRepositoryError handles common repository error patterns
func HandleRepositoryError(
	err error,
	message string,
	errorCode string,
	resourceType string,
	retryErr error,
) error {
	// If it's already a typed error, return it directly
	if _, ok := retryErr.(*errors.NotFoundError); ok {
		return retryErr
	}
	if _, ok := retryErr.(*errors.ValidationError); ok {
		return retryErr
	}
	if _, ok := retryErr.(*errors.DatabaseError); ok {
		return retryErr
	}

	// Otherwise, wrap it in a repository error
	return repoerrors.NewRepositoryError(retryErr, message, errorCode, resourceType)
}
