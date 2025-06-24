// Copyright (c) 2025 A Bit of Help, Inc.

package cache

import (
	"context"
	"testing"
	"time"

	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestNewCache(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Test cases
	tests := []struct {
		name     string
		config   *config.CacheConfig
		expected bool
	}{
		{
			name: "Enabled cache",
			config: &config.CacheConfig{
				Enabled:       true,
				TTL:           5 * time.Minute,
				MaxSize:       1000,
				PurgeInterval: 10 * time.Minute,
			},
			expected: true,
		},
		{
			name: "Disabled cache",
			config: &config.CacheConfig{
				Enabled:       false,
				TTL:           5 * time.Minute,
				MaxSize:       1000,
				PurgeInterval: 10 * time.Minute,
			},
			expected: false,
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cache, err := NewCache(&config.Config{Cache: *tc.config}, logger)
			require.NoError(t, err)

			if tc.expected {
				assert.NotNil(t, cache)
			} else {
				assert.Nil(t, cache)
			}
		})
	}
}

func TestCache_SetGet(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create a cache
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled:       true,
			TTL:           5 * time.Minute,
			MaxSize:       1000,
			PurgeInterval: 10 * time.Minute,
		},
	}
	cache, err := NewCache(cfg, logger)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Test Set and Get
	key := "test-key"
	value := "test-value"

	// Set the value
	cache.Set(key, value)

	// Get the value
	result, found := cache.Get(key)
	assert.True(t, found)
	assert.Equal(t, value, result)

	// Get a non-existent key
	result, found = cache.Get("non-existent-key")
	assert.False(t, found)
	assert.Nil(t, result)
}

func TestCache_Delete(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create a cache
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled:       true,
			TTL:           5 * time.Minute,
			MaxSize:       1000,
			PurgeInterval: 10 * time.Minute,
		},
	}
	cache, err := NewCache(cfg, logger)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Set a value
	key := "test-key"
	value := "test-value"
	cache.Set(key, value)

	// Verify it exists
	result, found := cache.Get(key)
	assert.True(t, found)
	assert.Equal(t, value, result)

	// Delete the key
	cache.Delete(key)

	// Verify it's gone
	result, found = cache.Get(key)
	assert.False(t, found)
	assert.Nil(t, result)
}

func TestCache_Expiration(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create a cache with a short TTL for testing
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled:       true,
			TTL:           100 * time.Millisecond, // Short TTL for testing
			MaxSize:       1000,
			PurgeInterval: 10 * time.Minute,
		},
	}
	cache, err := NewCache(cfg, logger)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Set a value
	key := "test-key"
	value := "test-value"
	cache.Set(key, value)

	// Verify it exists
	result, found := cache.Get(key)
	assert.True(t, found)
	assert.Equal(t, value, result)

	// Wait for expiration
	time.Sleep(200 * time.Millisecond)

	// Verify it's expired
	result, found = cache.Get(key)
	assert.False(t, found)
	assert.Nil(t, result)
}

func TestWithCache(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create a cache
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled:       true,
			TTL:           5 * time.Minute,
			MaxSize:       1000,
			PurgeInterval: 10 * time.Minute,
		},
	}
	cache, err := NewCache(cfg, logger)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Test WithCache
	key := "test-key"
	value := "test-value"
	callCount := 0

	// Function that should only be called once
	fn := func() (interface{}, error) {
		callCount++
		return value, nil
	}

	// First call should execute the function
	result, err := WithCache(cache, key, fn)
	require.NoError(t, err)
	assert.Equal(t, value, result)
	assert.Equal(t, 1, callCount)

	// Second call should use the cached value
	result, err = WithCache(cache, key, fn)
	require.NoError(t, err)
	assert.Equal(t, value, result)
	assert.Equal(t, 1, callCount) // Still 1, function not called again
}

func TestWithContextCache(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create a cache
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled:       true,
			TTL:           5 * time.Minute,
			MaxSize:       1000,
			PurgeInterval: 10 * time.Minute,
		},
	}
	cache, err := NewCache(cfg, logger)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Test WithContextCache
	ctx := context.Background()
	key := "test-key"
	value := "test-value"
	callCount := 0

	// Function that should only be called once
	fn := func(ctx context.Context) (interface{}, error) {
		callCount++
		return value, nil
	}

	// First call should execute the function
	result, err := WithContextCache(ctx, cache, key, fn)
	require.NoError(t, err)
	assert.Equal(t, value, result)
	assert.Equal(t, 1, callCount)

	// Second call should use the cached value
	result, err = WithContextCache(ctx, cache, key, fn)
	require.NoError(t, err)
	assert.Equal(t, value, result)
	assert.Equal(t, 1, callCount) // Still 1, function not called again
}

func TestCache_Shutdown(t *testing.T) {
	// Create a test logger
	logger := zaptest.NewLogger(t)

	// Create a cache
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled:       true,
			TTL:           5 * time.Minute,
			MaxSize:       1000,
			PurgeInterval: 10 * time.Minute,
		},
	}
	cache, err := NewCache(cfg, logger)
	require.NoError(t, err)
	require.NotNil(t, cache)

	// Shutdown should not panic
	assert.NotPanics(t, func() {
		cache.Shutdown()
	})
}

func TestCache_NilSafety(t *testing.T) {
	// Test that nil cache operations don't panic
	var cache *Cache

	assert.NotPanics(t, func() {
		cache.Set("key", "value")
	})

	assert.NotPanics(t, func() {
		_, _ = cache.Get("key")
	})

	assert.NotPanics(t, func() {
		cache.Delete("key")
	})

	assert.NotPanics(t, func() {
		cache.Shutdown()
	})

	assert.NotPanics(t, func() {
		_, _ = WithCache(cache, "key", func() (interface{}, error) {
			return "value", nil
		})
	})

	assert.NotPanics(t, func() {
		_, _ = WithContextCache(context.Background(), cache, "key", func(ctx context.Context) (interface{}, error) {
			return "value", nil
		})
	})
}
