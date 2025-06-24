// Copyright (c) 2025 A Bit of Help, Inc.

// Package rate provides functionality for rate limiting to protect resources.
package rate

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	"github.com/abitofhelp/servicelib/errors"
	"go.uber.org/zap"
)

// RateLimiter implements a token bucket rate limiter to protect resources
// from being overwhelmed by too many requests.
type RateLimiter struct {
	name           string
	config         *config.RateConfig
	logger         *zap.Logger
	tokens         int
	lastRefillTime time.Time
	mutex          sync.Mutex
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

	return &RateLimiter{
		name:           name,
		config:         cfg,
		logger:         logger,
		tokens:         cfg.BurstSize,
		lastRefillTime: time.Now(),
		mutex:          sync.Mutex{},
	}
}

// Allow checks if a request should be allowed based on the rate limit
// It returns true if the request is allowed, false otherwise
func (rl *RateLimiter) Allow() bool {
	if rl == nil {
		// If rate limiter is disabled, allow all requests
		return true
	}

	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	// Refill tokens based on time elapsed since last refill
	rl.refillTokens()

	// Check if we have tokens available
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

// Execute executes the given function with rate limiting
// If the rate limit is exceeded, it will return an error immediately
// Otherwise, it will execute the function
func (rl *RateLimiter) Execute(ctx context.Context, operation string, fn func(ctx context.Context) error) error {
	if rl == nil {
		// If rate limiter is disabled, just execute the function
		return fn(ctx)
	}

	// Check if the request is allowed
	if !rl.Allow() {
		rl.logger.Warn("Rate limit exceeded, rejecting request",
			zap.String("rate_limiter", rl.name),
			zap.String("operation", operation))
		return errors.NewApplicationError(errors.InternalErrorCode, fmt.Sprintf("rate limit exceeded for %s", rl.name), nil)
	}

	// Execute the function
	return fn(ctx)
}

// ExecuteWithWait executes the given function with rate limiting
// If the rate limit is exceeded, it will wait until a token is available
// and then execute the function
func (rl *RateLimiter) ExecuteWithWait(ctx context.Context, operation string, fn func(ctx context.Context) error) error {
	if rl == nil {
		// If rate limiter is disabled, just execute the function
		return fn(ctx)
	}

	// Wait until a token is available or context is canceled
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if rl.Allow() {
				// Execute the function
				return fn(ctx)
			}
			// Wait a bit before trying again
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// refillTokens refills tokens based on time elapsed since last refill
func (rl *RateLimiter) refillTokens() {
	now := time.Now()
	elapsed := now.Sub(rl.lastRefillTime)
	tokensToAdd := int(elapsed.Seconds() * float64(rl.config.RequestsPerSecond))

	if tokensToAdd > 0 {
		rl.tokens = min(rl.tokens+tokensToAdd, rl.config.BurstSize)
		rl.lastRefillTime = now
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Reset resets the rate limiter to its initial state
func (rl *RateLimiter) Reset() {
	if rl == nil {
		return
	}

	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	rl.tokens = rl.config.BurstSize
	rl.lastRefillTime = time.Now()
}
