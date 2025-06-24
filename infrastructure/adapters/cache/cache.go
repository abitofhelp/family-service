// Copyright (c) 2025 A Bit of Help, Inc.

// Package cache provides functionality for caching frequently accessed data.
// This package is a wrapper around the servicelib cache package.
package cache

import (
	"context"

	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	"github.com/abitofhelp/servicelib/cache"
	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

// Cache is a wrapper around servicelib's cache
type Cache struct {
	cache *cache.Cache[interface{}]
}

// NewCache creates a new cache with the given configuration
func NewCache(cfg *config.Config, logger *zap.Logger) (*Cache, error) {
	if !cfg.Cache.Enabled {
		logger.Info("Cache is disabled")
		return nil, nil
	}

	logger.Info("Initializing cache",
		zap.Duration("ttl", cfg.Cache.TTL),
		zap.Int("max_size", cfg.Cache.MaxSize),
		zap.Duration("purge_interval", cfg.Cache.PurgeInterval))

	// Create servicelib cache configuration
	cacheConfig := cache.DefaultConfig().
		WithEnabled(cfg.Cache.Enabled).
		WithTTL(cfg.Cache.TTL).
		WithMaxSize(cfg.Cache.MaxSize).
		WithPurgeInterval(cfg.Cache.PurgeInterval)

	// Create servicelib cache options
	contextLogger := logging.NewContextLogger(logger)
	options := cache.DefaultOptions().
		WithLogger(contextLogger).
		WithName("family-service-cache")

	// Create servicelib cache
	serviceCache := cache.NewCache[interface{}](cacheConfig, options)

	logger.Info("Cache initialized successfully")
	return &Cache{
		cache: serviceCache,
	}, nil
}

// Set adds an item to the cache with the default expiration time
func (c *Cache) Set(key string, value interface{}) {
	if c == nil || c.cache == nil {
		return
	}

	c.cache.Set(context.Background(), key, value)
}

// Get retrieves an item from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	if c == nil || c.cache == nil {
		return nil, false
	}

	return c.cache.Get(context.Background(), key)
}

// Delete removes an item from the cache
func (c *Cache) Delete(key string) {
	if c == nil || c.cache == nil {
		return
	}

	c.cache.Delete(context.Background(), key)
}

// Shutdown stops the cleanup timer
func (c *Cache) Shutdown() {
	if c == nil || c.cache == nil {
		return
	}

	c.cache.Shutdown()
}

// WithCache is a middleware that adds caching to a function
func WithCache(cache *Cache, key string, fn func() (interface{}, error)) (interface{}, error) {
	if cache == nil || cache.cache == nil {
		return fn()
	}

	// Create a wrapper function that adapts to servicelib's cache.WithCache
	wrapper := func(ctx context.Context) (interface{}, error) {
		return fn()
	}

	return cache.WithContextCache(context.Background(), key, wrapper)
}

// WithContextCache is a middleware that adds caching to a function with context
func WithContextCache(ctx context.Context, cache *Cache, key string, fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	if cache == nil || cache.cache == nil {
		return fn(ctx)
	}

	return cache.WithContextCache(ctx, key, fn)
}

// WithContextCache is a helper method that delegates to servicelib's cache.WithCache
func (c *Cache) WithContextCache(ctx context.Context, key string, fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	return cache.WithCache(ctx, c.cache, key, fn)
}
