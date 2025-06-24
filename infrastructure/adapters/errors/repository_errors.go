// Copyright (c) 2025 A Bit of Help, Inc.

// Package errors provides common error handling functions for repository adapters.
package errors

import (
	"github.com/abitofhelp/servicelib/errors"
)

// Common error codes used across repository adapters
const (
	// Database-specific error codes
	DatabaseErrorCode = "DATABASE_ERROR"
	MongoErrorCode    = "MONGO_ERROR"
	PostgresErrorCode = "POSTGRES_ERROR"
	SQLiteErrorCode   = "SQLITE_ERROR"

	// Data processing error codes
	JSONErrorCode       = "JSON_ERROR"
	DataFormatErrorCode = "DATA_FORMAT_ERROR"
	ConversionErrorCode = "CONVERSION_ERROR"
)

// NewRepositoryError is a helper function that wraps errors.NewDatabaseError
// to provide a consistent error handling approach across all repository adapters.
// It maps error codes to appropriate operations and tables.
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

	return errors.NewDatabaseError(message, operation, table, err)
}

// NewValidationError creates a new validation error with the given message, field, and cause.
func NewValidationError(message string, field string, cause error) error {
	return errors.NewValidationError(message, field, cause)
}

// NewNotFoundError creates a new not found error with the given resource type, resource ID, and cause.
func NewNotFoundError(resourceType string, resourceID string, cause error) error {
	return errors.NewNotFoundError(resourceType, resourceID, cause)
}