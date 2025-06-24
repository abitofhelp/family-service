// Copyright (c) 2025 A Bit of Help, Inc.

// Package di provides repository initializers for the family-service application.
// This file contains functions for initializing repositories using servicelib.
package di

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/abitofhelp/family-service/infrastructure/adapters/mongo"
	"github.com/abitofhelp/family-service/infrastructure/adapters/postgres"
	"github.com/abitofhelp/family-service/infrastructure/adapters/sqlite"
	servicedi "github.com/abitofhelp/servicelib/di"
	"github.com/abitofhelp/servicelib/logging"
	"github.com/jackc/pgx/v5/pgxpool"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

// MongoInitializerFunc defines the function signature for initializing a MongoDB collection
type MongoInitializerFunc func(ctx context.Context, uri string, dbName string, collectionName string, logger *zap.Logger) (interface{}, error)

// DefaultMongoInitializer is the default implementation using servicelib
var DefaultMongoInitializer MongoInitializerFunc = servicedi.GenericMongoInitializer

// InitMongoRepository initializes a MongoDB repository
// This function uses servicelib/di.GenericMongoInitializer internally by default,
// but allows for dependency injection of a different initializer for testing
func InitMongoRepository(ctx context.Context, uri string, zapLogger *zap.Logger, initializer ...MongoInitializerFunc) (*mongo.MongoFamilyRepository, error) {
	// Create a context logger
	logger := logging.NewContextLogger(zapLogger)

	// Use the provided initializer or the default one
	initFunc := DefaultMongoInitializer
	if len(initializer) > 0 && initializer[0] != nil {
		initFunc = initializer[0]
	}

	// Use the initializer
	collection, err := initFunc(ctx, uri, "family_service", "families", zapLogger)
	if err != nil {
		return nil, err
	}

	// Cast the result to *mongodb.Collection
	mongoCollection, ok := collection.(*mongodb.Collection)
	if !ok {
		return nil, fmt.Errorf("failed to cast MongoDB collection to *mongodb.Collection")
	}

	// Create repository with skipIndexCreation=true for tests to avoid nil pointer dereference
	return mongo.NewMongoFamilyRepository(mongoCollection, logger, true), nil
}

// PostgresInitializerFunc defines the function signature for initializing a PostgreSQL pool
type PostgresInitializerFunc func(ctx context.Context, dsn string, logger *zap.Logger) (interface{}, error)

// DefaultPostgresInitializer is the default implementation using servicelib
var DefaultPostgresInitializer PostgresInitializerFunc = servicedi.GenericPostgresInitializer

// InitPostgresRepository initializes a PostgreSQL repository
// This function uses servicelib/di.GenericPostgresInitializer internally by default,
// but allows for dependency injection of a different initializer for testing
func InitPostgresRepository(ctx context.Context, dsn string, zapLogger *zap.Logger, initializer ...PostgresInitializerFunc) (*postgres.PostgresFamilyRepository, error) {
	// Create a context logger
	logger := logging.NewContextLogger(zapLogger)

	// Use the provided initializer or the default one
	initFunc := DefaultPostgresInitializer
	if len(initializer) > 0 && initializer[0] != nil {
		initFunc = initializer[0]
	}

	// Use the initializer
	pool, err := initFunc(ctx, dsn, zapLogger)
	if err != nil {
		return nil, err
	}

	// Cast the result to *pgxpool.Pool
	pgxPool, ok := pool.(*pgxpool.Pool)
	if !ok {
		return nil, fmt.Errorf("failed to cast PostgreSQL pool to *pgxpool.Pool")
	}

	// Create repository
	return postgres.NewPostgresFamilyRepository(pgxPool, logger), nil
}

// SQLiteInitializerFunc defines the function signature for initializing a SQLite database
type SQLiteInitializerFunc func(ctx context.Context, uri string, logger *zap.Logger) (interface{}, error)

// DefaultSQLiteInitializer is the default implementation using servicelib
var DefaultSQLiteInitializer SQLiteInitializerFunc = servicedi.GenericSQLiteInitializer

// InitSQLiteRepository initializes a SQLite repository
// This function uses servicelib/di.GenericSQLiteInitializer internally by default,
// but allows for dependency injection of a different initializer for testing
func InitSQLiteRepository(ctx context.Context, uri string, zapLogger *zap.Logger, initializer ...SQLiteInitializerFunc) (*sqlite.SQLiteFamilyRepository, error) {
	// Create a context logger
	logger := logging.NewContextLogger(zapLogger)

	// Use the provided initializer or the default one
	initFunc := DefaultSQLiteInitializer
	if len(initializer) > 0 && initializer[0] != nil {
		initFunc = initializer[0]
	}

	// Use the initializer
	sqliteDB, err := initFunc(ctx, uri, zapLogger)
	if err != nil {
		return nil, err
	}

	// Cast the result to *sql.DB
	sqlDB, ok := sqliteDB.(*sql.DB)
	if !ok {
		return nil, fmt.Errorf("failed to cast SQLite database to *sql.DB")
	}

	// Create repository
	return sqlite.NewSQLiteFamilyRepository(sqlDB, logger), nil
}
