// Copyright (c) 2025 A Bit of Help, Inc.

// Package di provides repository initializers for the family-service application.
// This file contains deprecated functions for backward compatibility.
package di

import (
	"context"
	"time"

	"github.com/abitofhelp/family-service/infrastructure/adapters/mongo"
	"github.com/abitofhelp/family-service/infrastructure/adapters/postgres"
	"github.com/abitofhelp/family-service/infrastructure/adapters/sqlite"
	"github.com/abitofhelp/servicelib/db"
	"github.com/abitofhelp/servicelib/di"
	"github.com/abitofhelp/servicelib/logging"
	"go.uber.org/zap"
)

// InitMongoRepository initializes a MongoDB repository
// Deprecated: Use servicelib/di.GenericMongoInitializer instead
func InitMongoRepository(ctx context.Context, uri string, zapLogger *zap.Logger) (*mongo.MongoFamilyRepository, error) {
	// Create a context logger
	logger := logging.NewContextLogger(zapLogger)

	// Initialize MongoDB client using the db package
	client, err := db.InitMongoClient(ctx, uri, di.DefaultTimeout)
	if err != nil {
		return nil, err
	}

	// Log successful connection
	db.LogDatabaseConnection(ctx, logger, "MongoDB")

	// Get the families collection
	collection := client.Database("family_service").Collection("families")

	// Create repository
	return mongo.NewMongoFamilyRepository(collection), nil
}

// InitPostgresRepository initializes a PostgreSQL repository
// Deprecated: Use servicelib/di.GenericPostgresInitializer instead
func InitPostgresRepository(ctx context.Context, dsn string, zapLogger *zap.Logger) (*postgres.PostgresFamilyRepository, error) {
	// Create a context logger
	logger := logging.NewContextLogger(zapLogger)

	// Initialize PostgreSQL connection pool using the db package
	pool, err := db.InitPostgresPool(ctx, dsn, di.DefaultTimeout)
	if err != nil {
		return nil, err
	}

	// Log successful connection
	db.LogDatabaseConnection(ctx, logger, "PostgreSQL")

	// Create repository
	return postgres.NewPostgresFamilyRepository(pool), nil
}

// InitSQLiteRepository initializes a SQLite repository
// Deprecated: Use servicelib/di.GenericSQLiteInitializer instead
func InitSQLiteRepository(ctx context.Context, uri string, zapLogger *zap.Logger) (*sqlite.SQLiteFamilyRepository, error) {
	// Create a context logger
	logger := logging.NewContextLogger(zapLogger)

	// Initialize SQLite database connection using the db package
	sqliteDB, err := db.InitSQLiteDB(ctx, uri, di.DefaultTimeout, time.Hour, 10, 5)
	if err != nil {
		return nil, err
	}

	// Log successful connection
	db.LogDatabaseConnection(ctx, logger, "SQLite")

	// Create repository
	return sqlite.NewSQLiteFamilyRepository(sqliteDB), nil
}