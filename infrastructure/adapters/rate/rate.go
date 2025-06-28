// Copyright (c) 2025 A Bit of Help, Inc.

// Package rate provides functionality for rate limiting to protect resources.
// This package is a wrapper around the servicelib rate package.
package rate

import (
	"context"

	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/abitofhelp/servicelib/rate"
	"go.uber.org/zap"
)

// RateLimiter implements a token bucket rate limiter to protect resources
// from being overwhelmed by too many requests.
type RateLimiter struct {
	name   string
	rl     *rate.RateLimiter
	logger *zap.Logger
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(name string, cfg *config.RateConfig, logger *zap.Logger) *RateLimiter {
	if !cfg.Enabled {
		logger.Info("Rate limiter is disabled", zap.String("name", name))
		return nil
	}

	logger.Info("Initializing rate limiter",
		zap.String("name", name),
		zap.Int("requests_per_second", cfg.RequestsPerSecond),
		zap.Int("burst_size", cfg.BurstSize))

	// Create servicelib rate configuration
	rateConfig := rate.DefaultConfig().
		WithEnabled(cfg.Enabled).
		WithRequestsPerSecond(cfg.RequestsPerSecond).
		WithBurstSize(cfg.BurstSize)

	// Create servicelib rate options
	contextLogger := logging.NewContextLogger(logger)
	options := rate.DefaultOptions().
		WithLogger(contextLogger).
		WithName(name)

	// Create servicelib rate limiter
	serviceRL := rate.NewRateLimiter(rateConfig, options)

	return &RateLimiter{
		name:   name,
		rl:     serviceRL,
		logger: logger,
	}
}

// Allow checks if a request should be allowed based on the rate limit
// It returns true if the request is allowed, false otherwise
func (rl *RateLimiter) Allow() bool {
	if rl == nil || rl.rl == nil {
		// If rate limiter is disabled, allow all requests
		return true
	}

	return rl.rl.Allow()
}

// Execute executes the given function with rate limiting
// If the rate limit is exceeded, it will return an error immediately
// Otherwise, it will execute the function
func (rl *RateLimiter) Execute(ctx context.Context, operation string, fn func(ctx context.Context) error) error {
	if rl == nil || rl.rl == nil {
		// If rate limiter is disabled, just execute the function
		return fn(ctx)
	}

	// Create a wrapper function that adapts to servicelib's rate.Execute
	wrapper := func(ctx context.Context) (interface{}, error) {
		err := fn(ctx)
		return nil, err
	}

	// Execute the function with rate limiting
	_, err := rate.Execute(ctx, rl.rl, operation, wrapper)

	return err
}

// ExecuteWithWait executes the given function with rate limiting
// If the rate limit is exceeded, it will wait until a token is available
// and then execute the function
func (rl *RateLimiter) ExecuteWithWait(ctx context.Context, operation string, fn func(ctx context.Context) error) error {
	if rl == nil || rl.rl == nil {
		// If rate limiter is disabled, just execute the function
		return fn(ctx)
	}

	// Create a wrapper function that adapts to servicelib's rate.ExecuteWithWait
	wrapper := func(ctx context.Context) (interface{}, error) {
		err := fn(ctx)
		return nil, err
	}

	// Execute the function with rate limiting and waiting
	_, err := rate.ExecuteWithWait(ctx, rl.rl, operation, wrapper)

	return err
}

// Reset resets the rate limiter to its initial state
func (rl *RateLimiter) Reset() {
	if rl == nil || rl.rl == nil {
		return
	}

	rl.rl.Reset()
}

// Execute is a package-level function that executes the given function with rate limiting
// It's a wrapper around the servicelib rate.Execute function
// This function is used by the repository implementations
func Execute(ctx context.Context, rl *RateLimiter, operation string, fn func(ctx context.Context) (bool, error)) (bool, error) {
	if rl == nil || rl.rl == nil {
		// If rate limiter is disabled, just execute the function
		return fn(ctx)
	}

	// Execute the function with rate limiting using the servicelib rate.Execute
	return rate.Execute(ctx, rl.rl, operation, fn)
}
