// Copyright (c) 2025 A Bit of Help, Inc.

// Package repository provides base implementations and utilities for data repositories.
//
// In Hexagonal Architecture (Ports and Adapters), repositories are adapters that
// implement the repository interfaces (ports) defined in the domain layer.
// They provide the actual implementation for persisting and retrieving domain entities.
//
// This package includes:
// 1. A BaseRepository that implements common functionality like resilience patterns
// 2. Configuration utilities for repository settings
// 3. Error handling specific to repository operations
//
// The repositories in this package follow several important patterns:
// - Circuit Breaker: Prevents cascading failures when a database is unavailable
// - Rate Limiting: Protects databases from being overwhelmed with requests
// - Retries with Backoff: Handles transient errors automatically
// - Timeout Management: Ensures operations don't hang indefinitely
//
// These patterns make the repositories more resilient to failures and
// help maintain system stability even when external dependencies have issues.
package repository

import (
	"context"
	"sync"
	"time"

	circuit "github.com/abitofhelp/family-service/infrastructure/adapters/circuitwrapper"
	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	repoerrors "github.com/abitofhelp/family-service/infrastructure/adapters/errors"
	"github.com/abitofhelp/family-service/infrastructure/adapters/loggingwrapper"
	rate "github.com/abitofhelp/family-service/infrastructure/adapters/ratewrapper"
	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/retry"
)

var (
	// globalConfig holds the application-wide configuration instance.
	// This is used to configure repository behavior like retry policies,
	// circuit breaker settings, and rate limiting.
	globalConfig     *config.Config
	globalConfigOnce sync.Once // Ensures the config is set only once
)

// SetGlobalConfig sets the global configuration instance.
//
// This function uses sync.Once to ensure that the configuration is set only once,
// making it safe to call from multiple goroutines. This is important for
// application startup when multiple components might try to initialize
// the configuration simultaneously.
//
// Parameters:
//   - cfg: The configuration instance to use globally
//
// This function is typically called during application startup, before
// any repositories are created.
func SetGlobalConfig(cfg *config.Config) {
	globalConfigOnce.Do(func() {
		globalConfig = cfg
	})
}

// GetRetryConfig returns the retry configuration for repository operations.
//
// This function provides a consistent retry configuration across all repositories.
// It first tries to use the global configuration if available, and falls back to
// sensible defaults if not.
//
// The retry configuration includes:
// - Maximum number of retry attempts
// - Initial backoff duration (how long to wait before the first retry)
// - Maximum backoff duration (upper limit on how long to wait between retries)
//
// Returns:
//   - A retry.Config instance with appropriate settings
//
// The exponential backoff with jitter strategy helps prevent the "thundering herd"
// problem when multiple clients retry at the same time after a failure.
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

// BaseRepository provides common functionality for all repository implementations.
//
// This struct implements resilience patterns that should be used consistently
// across all repository implementations. By centralizing these patterns in a base struct,
// we ensure that all repositories benefit from:
//
// 1. Circuit Breaking: Prevents cascading failures when a database is unavailable
//    by temporarily stopping requests after a threshold of failures is reached.
//
// 2. Rate Limiting: Protects databases from being overwhelmed with requests
//    by limiting the rate at which operations can be performed.
//
// 3. Logging: Provides consistent logging of repository operations and errors
//    for observability and troubleshooting.
//
// 4. Timeout Management: Ensures operations don't hang indefinitely by
//    applying a default timeout to all operations.
//
// Repository implementations should embed this struct and use its ExecuteWithResilience
// method to perform operations with these resilience patterns applied.
type BaseRepository struct {
	Logger         *loggingwrapper.ContextLogger // For logging operations and errors
	CircuitBreaker *circuit.CircuitBreaker       // For preventing cascading failures
	RateLimiter    *rate.RateLimiter             // For protecting against overwhelming the database
	DefaultTimeout time.Duration                 // Default timeout for operations
}

// NewBaseRepository creates a new BaseRepository with the given logger and database name.
//
// This function initializes a BaseRepository with appropriate circuit breaker and
// rate limiter configurations. It uses the global configuration if available,
// or falls back to sensible defaults.
//
// The circuit breaker and rate limiter are configured specifically for the
// given database, allowing different databases to have different settings.
//
// Parameters:
//   - logger: Logger for recording operations and errors
//   - dbName: Name of the database (used for naming the circuit breaker and rate limiter)
//
// Returns:
//   - A new BaseRepository instance
//
// Panics if the logger is nil, as logging is essential for observability.
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

// IsRetryableError determines if an error should be retried based on its type.
//
// This function implements a key part of the retry strategy by distinguishing
// between errors that are likely to be resolved by retrying and those that won't.
//
// Errors that should NOT be retried:
// - Not Found errors: The entity doesn't exist, retrying won't help
// - Validation errors: The input is invalid, retrying won't fix it
//
// Errors that SHOULD be retried:
// - Network errors: Temporary network issues might resolve
// - Timeout errors: The operation might succeed if tried again
// - Transient database errors: Temporary database issues might resolve
//
// This distinction is important for implementing an effective retry strategy
// that balances resilience against wasting resources on futile retries.
//
// Parameters:
//   - err: The error to evaluate
//
// Returns:
//   - true if the error should be retried
//   - false if retrying is unlikely to help
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

// ExecuteWithResilience executes an operation with retry, circuit breaker, and rate limiter.
//
// This method implements a comprehensive resilience strategy by combining multiple
// patterns in the correct order:
//
// 1. Rate Limiting (outermost): Controls the rate of requests to prevent overwhelming
//    the database, regardless of success or failure.
//
// 2. Circuit Breaking (middle): Prevents sending requests to a failing database
//    after a threshold of failures is reached.
//
// 3. Retries with Backoff (innermost): Automatically retries operations that fail
//    with transient errors, using exponential backoff with jitter.
//
// 4. Timeout Management: Ensures the operation doesn't run indefinitely by
//    applying a default timeout.
//
// This layered approach provides robust resilience against various failure modes
// while maintaining good performance characteristics.
//
// Parameters:
//   - ctx: Context for the operation (used for cancellation, tracing, etc.)
//   - operation: The function to execute with resilience
//   - operationName: Name of the operation (used for logging and metrics)
//
// Returns:
//   - An error if the operation ultimately fails after applying all resilience patterns
//   - nil if the operation succeeds
//
// Example usage:
//
//	err := repo.ExecuteWithResilience(ctx, func(ctx context.Context) error {
//	    return repo.db.QueryRow(ctx, "SELECT * FROM families WHERE id = $1", id).Scan(&family.ID, &family.Status)
//	}, "GetFamilyByID")
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
