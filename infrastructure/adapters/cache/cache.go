// Copyright (c) 2025 A Bit of Help, Inc.

// Package cache provides functionality for caching frequently accessed data.
package cache

import (
	"context"
	"sync"
	"time"

	"github.com/abitofhelp/family-service/infrastructure/adapters/config"
	"go.uber.org/zap"
)

// Item represents a cached item with its value and expiration time
type Item struct {
	Value      interface{}
	Expiration int64
}

// Cache is a simple in-memory cache with expiration
type Cache struct {
	items             map[string]Item
	mu                sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	logger            *zap.Logger
	stopCleanup       chan bool
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

	cache := &Cache{
		items:             make(map[string]Item),
		defaultExpiration: cfg.Cache.TTL,
		cleanupInterval:   cfg.Cache.PurgeInterval,
		logger:            logger,
		stopCleanup:       make(chan bool),
	}

	// Start the cleanup goroutine
	go cache.startCleanupTimer()

	logger.Info("Cache initialized successfully")
	return cache, nil
}

// Set adds an item to the cache with the default expiration time
func (c *Cache) Set(key string, value interface{}) {
	if c == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = Item{
		Value:      value,
		Expiration: time.Now().Add(c.defaultExpiration).UnixNano(),
	}
}

// Get retrieves an item from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	if c == nil {
		return nil, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	// Check if the item has expired
	if time.Now().UnixNano() > item.Expiration {
		return nil, false
	}

	return item.Value, true
}

// Delete removes an item from the cache
func (c *Cache) Delete(key string) {
	if c == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.items, key)
}

// startCleanupTimer starts the cleanup timer
func (c *Cache) startCleanupTimer() {
	ticker := time.NewTicker(c.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.stopCleanup:
			return
		}
	}
}

// cleanup removes expired items from the cache
func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UnixNano()
	for k, v := range c.items {
		if now > v.Expiration {
			delete(c.items, k)
		}
	}
}

// Shutdown stops the cleanup timer
func (c *Cache) Shutdown() {
	if c == nil {
		return
	}

	c.logger.Info("Shutting down cache")
	c.stopCleanup <- true
	c.logger.Info("Cache shut down successfully")
}

// WithCache is a middleware that adds caching to a function
func WithCache(cache *Cache, key string, fn func() (interface{}, error)) (interface{}, error) {
	if cache == nil {
		return fn()
	}

	// Try to get from cache
	if value, found := cache.Get(key); found {
		return value, nil
	}

	// If not found, call the function
	value, err := fn()
	if err != nil {
		return nil, err
	}

	// Store the result in cache
	cache.Set(key, value)
	return value, nil
}

// WithContextCache is a middleware that adds caching to a function with context
func WithContextCache(ctx context.Context, cache *Cache, key string, fn func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	if cache == nil {
		return fn(ctx)
	}

	// Try to get from cache
	if value, found := cache.Get(key); found {
		return value, nil
	}

	// If not found, call the function
	value, err := fn(ctx)
	if err != nil {
		return nil, err
	}

	// Store the result in cache
	cache.Set(key, value)
	return value, nil
}
