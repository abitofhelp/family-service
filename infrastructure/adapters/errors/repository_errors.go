// Copyright (c) 2025 A Bit of Help, Inc.

// Package errors provides common error handling functions for repository adapters.
//
// This package centralizes error handling for repository operations, ensuring
// consistent error types, messages, and codes across different storage implementations.
// It provides error constants and helper functions that make it easier for junior
// developers to create appropriate error types without having to understand all
// the details of the error handling system.
package errors

import (
	serviceerrors "github.com/abitofhelp/servicelib/errors"
)

// Database and data processing error codes used across repository adapters.
// These constants provide standardized error codes that can be used to
// identify specific types of errors in logs, metrics, and client responses.
const (
	// DatabaseErrorCode represents a generic database operation error.
	// Use this when a more specific error code is not applicable.
	DatabaseErrorCode = "DATABASE_ERROR"

	// MongoErrorCode represents a MongoDB-specific error.
	// Use this for errors that are specific to MongoDB operations.
	MongoErrorCode = "MONGO_ERROR"

	// PostgresErrorCode represents a PostgreSQL-specific error.
	// Use this for errors that are specific to PostgreSQL operations.
	PostgresErrorCode = "POSTGRES_ERROR"

	// SQLiteErrorCode represents a SQLite-specific error.
	// Use this for errors that are specific to SQLite operations.
	SQLiteErrorCode = "SQLITE_ERROR"

	// JSONErrorCode represents an error in JSON processing.
	// Use this for errors related to JSON serialization or deserialization.
	JSONErrorCode = "JSON_ERROR"

	// DataFormatErrorCode represents an error in data format.
	// Use this for errors related to invalid or unexpected data formats.
	DataFormatErrorCode = "DATA_FORMAT_ERROR"

	// ConversionErrorCode represents an error in data type conversion.
	// Use this for errors that occur when converting between data types.
	ConversionErrorCode = "CONVERSION_ERROR"
)

// NewRepositoryError creates a new database error with appropriate operation and table information.
//
// This helper function wraps errors.NewDatabaseError to provide a consistent error
// handling approach across all repository adapters. It automatically maps error codes
// to appropriate operations based on the provided code, simplifying error creation
// for repository implementations.
//
// Parameters:
//   - err: The original error that caused the repository operation to fail
//   - message: A human-readable description of what went wrong
//   - code: One of the error code constants defined in this package
//   - table: The name of the database table or collection being accessed
//
// Returns:
//   - A properly formatted database error that can be returned to callers
//
// Example usage:
//
//	if err := db.Collection("families").FindOne(ctx, filter).Decode(&family); err != nil {
//	    return nil, NewRepositoryError(err, "failed to find family", MongoErrorCode, "families")
//	}
func NewRepositoryError(err error, message string, code string, table string) error {
	// Default values
	operation := "operation"

	// If table is not specified, use "families" as default
	if table == "" {
		table = "families"
	}

	// Map common codes to operations
	switch code {
	case MongoErrorCode, PostgresErrorCode, SQLiteErrorCode, DatabaseErrorCode:
		operation = "query"
	case JSONErrorCode:
		operation = "unmarshal"
	case DataFormatErrorCode:
		operation = "parse"
	case ConversionErrorCode:
		operation = "convert"
	}

	return serviceerrors.NewDatabaseError(message, operation, table, err)
}
