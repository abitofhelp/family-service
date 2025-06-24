// Copyright (c) 2025 A Bit of Help, Inc.

package di

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// setupTest sets up common test dependencies
func setupTest(t *testing.T) (context.Context, *zap.Logger) {
	ctx := context.Background()
	logger := zaptest.NewLogger(t)
	return ctx, logger
}

// TestInitMongoRepository tests the InitMongoRepository function
func TestInitMongoRepository(t *testing.T) {
	// Setup
	ctx, logger := setupTest(t)

	// Create a mock MongoDB collection
	collection := &mongo.Collection{}

	// Test successful initialization
	successInitializer := func(ctx context.Context, uri string, dbName string, collectionName string, logger *zap.Logger) (interface{}, error) {
		// Verify input parameters
		assert.Equal(t, "mongodb://localhost:27017", uri)
		assert.Equal(t, "family_service", dbName)
		assert.Equal(t, "families", collectionName)
		assert.NotNil(t, logger)

		// Return the mock collection
		return collection, nil
	}

	repo, err := InitMongoRepository(ctx, "mongodb://localhost:27017", logger, successInitializer)

	// Verify results
	require.NoError(t, err)
	require.NotNil(t, repo)
	assert.Equal(t, collection, repo.Collection)

	// Test error case: Initializer returns an error
	errorInitializer := func(ctx context.Context, uri string, dbName string, collectionName string, logger *zap.Logger) (interface{}, error) {
		return nil, errors.New("connection error")
	}

	repo, err = InitMongoRepository(ctx, "mongodb://localhost:27017", logger, errorInitializer)

	// Verify error handling
	require.Error(t, err)
	assert.Nil(t, repo)
	assert.Contains(t, err.Error(), "connection error")

	// Test error case: Initializer returns wrong type
	wrongTypeInitializer := func(ctx context.Context, uri string, dbName string, collectionName string, logger *zap.Logger) (interface{}, error) {
		return "not a collection", nil
	}

	repo, err = InitMongoRepository(ctx, "mongodb://localhost:27017", logger, wrongTypeInitializer)

	// Verify error handling
	require.Error(t, err)
	assert.Nil(t, repo)
	assert.Contains(t, err.Error(), "failed to cast MongoDB collection")
}

// TestInitPostgresRepository tests the InitPostgresRepository function
func TestInitPostgresRepository(t *testing.T) {
	// Setup
	ctx, logger := setupTest(t)

	// Create a mock PostgreSQL pool
	pool := &pgxpool.Pool{}

	// Test successful initialization
	successInitializer := func(ctx context.Context, dsn string, logger *zap.Logger) (interface{}, error) {
		// Verify input parameters
		assert.Equal(t, "postgres://user:pass@localhost:5432/familydb", dsn)
		assert.NotNil(t, logger)

		// Return the mock pool
		return pool, nil
	}

	repo, err := InitPostgresRepository(ctx, "postgres://user:pass@localhost:5432/familydb", logger, successInitializer)

	// Verify results
	require.NoError(t, err)
	require.NotNil(t, repo)
	assert.Equal(t, pool, repo.DB)

	// Test error case: Initializer returns an error
	errorInitializer := func(ctx context.Context, dsn string, logger *zap.Logger) (interface{}, error) {
		return nil, errors.New("connection error")
	}

	repo, err = InitPostgresRepository(ctx, "postgres://user:pass@localhost:5432/familydb", logger, errorInitializer)

	// Verify error handling
	require.Error(t, err)
	assert.Nil(t, repo)
	assert.Contains(t, err.Error(), "connection error")

	// Test error case: Initializer returns wrong type
	wrongTypeInitializer := func(ctx context.Context, dsn string, logger *zap.Logger) (interface{}, error) {
		return "not a pool", nil
	}

	repo, err = InitPostgresRepository(ctx, "postgres://user:pass@localhost:5432/familydb", logger, wrongTypeInitializer)

	// Verify error handling
	require.Error(t, err)
	assert.Nil(t, repo)
	assert.Contains(t, err.Error(), "failed to cast PostgreSQL pool")
}

// TestInitSQLiteRepository tests the InitSQLiteRepository function
func TestInitSQLiteRepository(t *testing.T) {
	// Setup
	ctx, logger := setupTest(t)

	// Create a mock SQLite database
	db := &sql.DB{}

	// Test successful initialization
	successInitializer := func(ctx context.Context, uri string, logger *zap.Logger) (interface{}, error) {
		// Verify input parameters
		assert.Equal(t, "file:data/dev/sqlite/family_service.db?cache=shared&mode=rwc", uri)
		assert.NotNil(t, logger)

		// Return the mock database
		return db, nil
	}

	repo, err := InitSQLiteRepository(ctx, "file:data/dev/sqlite/family_service.db?cache=shared&mode=rwc", logger, successInitializer)

	// Verify results
	require.NoError(t, err)
	require.NotNil(t, repo)
	assert.Equal(t, db, repo.DB)

	// Test error case: Initializer returns an error
	errorInitializer := func(ctx context.Context, uri string, logger *zap.Logger) (interface{}, error) {
		return nil, errors.New("connection error")
	}

	repo, err = InitSQLiteRepository(ctx, "file:data/dev/sqlite/family_service.db?cache=shared&mode=rwc", logger, errorInitializer)

	// Verify error handling
	require.Error(t, err)
	assert.Nil(t, repo)
	assert.Contains(t, err.Error(), "connection error")

	// Test error case: Initializer returns wrong type
	wrongTypeInitializer := func(ctx context.Context, uri string, logger *zap.Logger) (interface{}, error) {
		return "not a database", nil
	}

	repo, err = InitSQLiteRepository(ctx, "file:data/dev/sqlite/family_service.db?cache=shared&mode=rwc", logger, wrongTypeInitializer)

	// Verify error handling
	require.Error(t, err)
	assert.Nil(t, repo)
	assert.Contains(t, err.Error(), "failed to cast SQLite database")
}
