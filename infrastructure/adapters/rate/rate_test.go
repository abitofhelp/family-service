// Copyright (c) 2025 A Bit of Help, Inc.

package rate

import (
	"context"
	"testing"
	"time"

	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestNewRateLimiter(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Test cases
	tests := []struct {
		name     string
		config   *config.RateConfig
		expected bool
	}{
		{
			name: "Enabled rate limiter",
			config: &config.RateConfig{
				Enabled:           true,
				RequestsPerSecond: 100,
				BurstSize:         50,
			},
			expected: true,
		},
		{
			name: "Disabled rate limiter",
			config: &config.RateConfig{
				Enabled:           false,
				RequestsPerSecond: 100,
				BurstSize:         50,
			},
			expected: false,
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rl := NewRateLimiter("test", tc.config, logger)

			if tc.expected {
				assert.NotNil(t, rl)
			} else {
				assert.Nil(t, rl)
			}
		})
	}
}

func TestRateLimiter_Allow(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create a rate limiter with a small burst size for testing
	cfg := &config.RateConfig{
		Enabled:           true,
		RequestsPerSecond: 10,
		BurstSize:         5, // Small burst size for testing
	}
	rl := NewRateLimiter("test", cfg, logger)
	require.NotNil(t, rl)

	// First 5 requests should be allowed (burst size)
	for i := 0; i < 5; i++ {
		assert.True(t, rl.Allow(), "Request %d should be allowed", i+1)
	}

	// Next request should be denied (burst size exceeded)
	assert.False(t, rl.Allow(), "Request 6 should be denied")

	// Wait for token refill (at least 100ms for 1 token at 10 RPS)
	time.Sleep(200 * time.Millisecond)

	// Should be allowed again after refill
	assert.True(t, rl.Allow(), "Request after wait should be allowed")
}

func TestRateLimiter_Execute_Success(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create a rate limiter
	cfg := &config.RateConfig{
		Enabled:           true,
		RequestsPerSecond: 10,
		BurstSize:         5,
	}
	rl := NewRateLimiter("test", cfg, logger)
	require.NotNil(t, rl)

	// Test successful execution
	ctx := context.Background()
	callCount := 0

	// Function that should succeed
	fn := func(ctx context.Context) error {
		callCount++
		return nil
	}

	// Execute the function
	err := rl.Execute(ctx, "test-operation", fn)
	assert.NoError(t, err)
	assert.Equal(t, 1, callCount)
}

func TestRateLimiter_Execute_RateLimitExceeded(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create a rate limiter with a small burst size for testing
	cfg := &config.RateConfig{
		Enabled:           true,
		RequestsPerSecond: 10,
		BurstSize:         1, // Small burst size for testing
	}
	rl := NewRateLimiter("test", cfg, logger)
	require.NotNil(t, rl)

	// Test rate limit exceeded
	ctx := context.Background()

	// First request should succeed
	err := rl.Execute(ctx, "test-operation", func(ctx context.Context) error {
		return nil
	})
	assert.NoError(t, err)

	// Second request should fail with rate limit exceeded
	err = rl.Execute(ctx, "test-operation", func(ctx context.Context) error {
		t.Fatal("This function should not be called when rate limit is exceeded")
		return nil
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rate limit exceeded")
}

func TestRateLimiter_ExecuteWithWait(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create a rate limiter with a small burst size for testing
	cfg := &config.RateConfig{
		Enabled:           true,
		RequestsPerSecond: 10,
		BurstSize:         1, // Small burst size for testing
	}
	rl := NewRateLimiter("test", cfg, logger)
	require.NotNil(t, rl)

	// Test execute with wait
	ctx := context.Background()
	callCount := 0

	// First request should succeed immediately
	start := time.Now()
	err := rl.ExecuteWithWait(ctx, "test-operation", func(ctx context.Context) error {
		callCount++
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, callCount)
	assert.Less(t, time.Since(start), 50*time.Millisecond, "First request should not wait")

	// Second request should wait for a token
	start = time.Now()
	err = rl.ExecuteWithWait(ctx, "test-operation", func(ctx context.Context) error {
		callCount++
		return nil
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, callCount)
	assert.GreaterOrEqual(t, time.Since(start), 50*time.Millisecond, "Second request should wait for a token")
}

func TestRateLimiter_ExecuteWithWait_ContextCancellation(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create a rate limiter with a small burst size for testing
	cfg := &config.RateConfig{
		Enabled:           true,
		RequestsPerSecond: 1,
		BurstSize:         1, // Small burst size for testing
	}
	rl := NewRateLimiter("test", cfg, logger)
	require.NotNil(t, rl)

	// Test context cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// First request should succeed
	err := rl.Execute(ctx, "test-operation", func(ctx context.Context) error {
		return nil
	})
	assert.NoError(t, err)

	// Cancel the context
	cancel()

	// Second request should fail with context canceled
	err = rl.ExecuteWithWait(ctx, "test-operation", func(ctx context.Context) error {
		t.Fatal("This function should not be called when context is canceled")
		return nil
	})
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
}

func TestRateLimiter_Reset(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create a rate limiter with a small burst size for testing
	cfg := &config.RateConfig{
		Enabled:           true,
		RequestsPerSecond: 10,
		BurstSize:         2, // Small burst size for testing
	}
	rl := NewRateLimiter("test", cfg, logger)
	require.NotNil(t, rl)

	// Use up all tokens
	assert.True(t, rl.Allow())
	assert.True(t, rl.Allow())
	assert.False(t, rl.Allow())

	// Reset the rate limiter
	rl.Reset()

	// Should have full burst size again
	assert.True(t, rl.Allow())
	assert.True(t, rl.Allow())
	assert.False(t, rl.Allow())
}

func TestRateLimiter_NilSafety(t *testing.T) {
	// Test that nil rate limiter operations don't panic
	var rl *RateLimiter

	assert.NotPanics(t, func() {
		_ = rl.Allow()
	})

	assert.NotPanics(t, func() {
		_ = rl.Execute(context.Background(), "test-operation", func(ctx context.Context) error {
			return nil
		})
	})

	assert.NotPanics(t, func() {
		_ = rl.ExecuteWithWait(context.Background(), "test-operation", func(ctx context.Context) error {
			return nil
		})
	})

	assert.NotPanics(t, func() {
		rl.Reset()
	})
}
