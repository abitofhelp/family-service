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

// InitMongoRepository initializes a MongoDB repository
// This function uses servicelib/di.GenericMongoInitializer internally
func InitMongoRepository(ctx context.Context, uri string, zapLogger *zap.Logger) (*mongo.MongoFamilyRepository, error) {
	// Use the generic initializer from servicelib
	collection, err := servicedi.GenericMongoInitializer(ctx, uri, "family_service", "families", zapLogger)
	if err != nil {
		return nil, err
	}

	// Cast the result to *mongodb.Collection
	mongoCollection, ok := collection.(*mongodb.Collection)
	if !ok {
		return nil, fmt.Errorf("failed to cast MongoDB collection to *mongodb.Collection")
	}

	// Create repository
	return mongo.NewMongoFamilyRepository(mongoCollection), nil
}

// InitPostgresRepository initializes a PostgreSQL repository
// This function uses servicelib/di.GenericPostgresInitializer internally
func InitPostgresRepository(ctx context.Context, dsn string, zapLogger *zap.Logger) (*postgres.PostgresFamilyRepository, error) {
	// Create a context logger
	logger := logging.NewContextLogger(zapLogger)

	// Use the generic initializer from servicelib
	pool, err := servicedi.GenericPostgresInitializer(ctx, dsn, zapLogger)
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

// InitSQLiteRepository initializes a SQLite repository
// This function uses servicelib/di.GenericSQLiteInitializer internally
func InitSQLiteRepository(ctx context.Context, uri string, zapLogger *zap.Logger) (*sqlite.SQLiteFamilyRepository, error) {
	// Create a context logger
	logger := logging.NewContextLogger(zapLogger)

	// Use the generic initializer from servicelib
	sqliteDB, err := servicedi.GenericSQLiteInitializer(ctx, uri, zapLogger)
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
