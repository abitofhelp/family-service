// Copyright (c) 2025 A Bit of Help, Inc.

package repository

import (
	"context"
	"sync"
	"time"

	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	repoerrors "github.com/abitofhelp/family-service/infrastructure/adapters/errors"
	"github.com/abitofhelp/servicelib/circuit"
	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/rate"
	"github.com/abitofhelp/servicelib/retry"
	"go.uber.org/zap"
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
type BaseRepository struct {
	Logger         *logging.ContextLogger
	CircuitBreaker *circuit.CircuitBreaker
	RateLimiter    *rate.RateLimiter
	DefaultTimeout time.Duration
}

// NewBaseRepository creates a new BaseRepository with the given logger
func NewBaseRepository(logger *logging.ContextLogger, dbName string) *BaseRepository {
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

	// Create a new zap logger for the circuit breaker and rate limiter
	zapLogger, _ := zap.NewProduction()
	contextLogger := logging.NewContextLogger(zapLogger)

	// Create circuit breaker config
	circuitBreakerConfig := circuit.DefaultConfig().
		WithEnabled(circuitConfig.Enabled).
		WithTimeout(circuitConfig.Timeout).
		WithMaxConcurrent(circuitConfig.MaxConcurrent).
		WithErrorThreshold(circuitConfig.ErrorThreshold).
		WithVolumeThreshold(circuitConfig.VolumeThreshold).
		WithSleepWindow(circuitConfig.SleepWindow)

	// Create circuit breaker options
	circuitBreakerOptions := circuit.DefaultOptions().
		WithName(dbName).
		WithLogger(contextLogger)

	// Create circuit breaker
	cb := circuit.NewCircuitBreaker(circuitBreakerConfig, circuitBreakerOptions)

	// Create rate limiter config
	rateLimiterConfig := rate.DefaultConfig().
		WithEnabled(rateConfig.Enabled).
		WithRequestsPerSecond(rateConfig.RequestsPerSecond).
		WithBurstSize(rateConfig.BurstSize)

	// Create rate limiter options
	rateLimiterOptions := rate.DefaultOptions().
		WithName(dbName).
		WithLogger(contextLogger)

	// Create rate limiter
	rl := rate.NewRateLimiter(rateLimiterConfig, rateLimiterOptions)

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
		// We need to wrap the circuitOperation to match the generic function signature
		circuitOpWrapper := func(ctx context.Context) (interface{}, error) {
			err := circuitOperation(ctx)
			return nil, err
		}
		_, err := circuit.Execute(ctx, r.CircuitBreaker, operationName, circuitOpWrapper)
		return err
	}

	// Execute with rate limiter
	// We need to wrap the rateOperation to match the generic function signature
	rateOpWrapper := func(ctx context.Context) (interface{}, error) {
		err := rateOperation(ctx)
		return nil, err
	}
	_, err := rate.Execute(ctxWithTimeout, r.RateLimiter, operationName, rateOpWrapper)

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