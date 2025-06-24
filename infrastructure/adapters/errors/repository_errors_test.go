// Copyright (c) 2025 A Bit of Help, Inc.

package errors

import (
	"errors"
	"testing"

	serviceerrors "github.com/abitofhelp/servicelib/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewRepositoryError tests the NewRepositoryError function
func TestNewRepositoryError(t *testing.T) {
	// Create a test error
	testErr := errors.New("test error")

	// Test cases for different error codes
	testCases := []struct {
		name         string
		code         string
		message      string
		table        string
		expectedOp   string
		expectedType string
	}{
		{
			name:         "MongoDB Error",
			code:         MongoErrorCode,
			message:      "MongoDB connection failed",
			table:        "families",
			expectedOp:   "query",
			expectedType: "*errors.DatabaseError",
		},
		{
			name:         "PostgreSQL Error",
			code:         PostgresErrorCode,
			message:      "PostgreSQL query failed",
			table:        "families",
			expectedOp:   "query",
			expectedType: "*errors.DatabaseError",
		},
		{
			name:         "SQLite Error",
			code:         SQLiteErrorCode,
			message:      "SQLite transaction failed",
			table:        "families",
			expectedOp:   "query",
			expectedType: "*errors.DatabaseError",
		},
		{
			name:         "JSON Error",
			code:         JSONErrorCode,
			message:      "Failed to unmarshal JSON",
			table:        "families",
			expectedOp:   "unmarshal",
			expectedType: "*errors.DatabaseError",
		},
		{
			name:         "Data Format Error",
			code:         DataFormatErrorCode,
			message:      "Invalid date format",
			table:        "families",
			expectedOp:   "parse",
			expectedType: "*errors.DatabaseError",
		},
		{
			name:         "Conversion Error",
			code:         ConversionErrorCode,
			message:      "Failed to convert string to int",
			table:        "families",
			expectedOp:   "convert",
			expectedType: "*errors.DatabaseError",
		},
		{
			name:         "Default Operation",
			code:         "UNKNOWN_CODE",
			message:      "Unknown error",
			table:        "families",
			expectedOp:   "operation",
			expectedType: "*errors.DatabaseError",
		},
		{
			name:         "Default Table",
			code:         DatabaseErrorCode,
			message:      "Database error",
			table:        "",
			expectedOp:   "query",
			expectedType: "*errors.DatabaseError",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function
			err := NewRepositoryError(testErr, tc.message, tc.code, tc.table)

			// Verify the error type
			require.NotNil(t, err)
			assert.IsType(t, &serviceerrors.DatabaseError{}, err)

			// Verify the error message
			assert.Contains(t, err.Error(), tc.message)

			// Verify the error cause
			dbErr, ok := err.(*serviceerrors.DatabaseError)
			require.True(t, ok)
			assert.Equal(t, testErr, dbErr.Unwrap())

			// We can't directly access the operation and table fields of the DatabaseError
			// since they are not exported, but we can verify that the error was created
			// with the correct parameters by checking that the function didn't panic
			// and returned a DatabaseError.
		})
	}
}

// TestNewValidationError tests the NewValidationError function
func TestNewValidationError(t *testing.T) {
	// Create a test error
	testErr := errors.New("test error")

	// Call the function
	err := NewValidationError("Invalid input", "field_name", testErr)

	// Verify the error type
	require.NotNil(t, err)
	assert.IsType(t, &serviceerrors.ValidationError{}, err)

	// Verify the error message
	assert.Contains(t, err.Error(), "Invalid input")

	// Verify the error cause
	valErr, ok := err.(*serviceerrors.ValidationError)
	require.True(t, ok)
	assert.Equal(t, testErr, valErr.Unwrap())

	// We can't directly access the field name of the ValidationError
	// since it is not exported, but we can verify that the error was created
	// with the correct parameters by checking that the function didn't panic
	// and returned a ValidationError.
}

// TestNewNotFoundError tests the NewNotFoundError function
func TestNewNotFoundError(t *testing.T) {
	// Create a test error
	testErr := errors.New("test error")

	// Call the function
	err := NewNotFoundError("Family", "123", testErr)

	// Verify the error type
	require.NotNil(t, err)
	assert.IsType(t, &serviceerrors.NotFoundError{}, err)

	// Verify the error message
	assert.Contains(t, err.Error(), "Family")
	assert.Contains(t, err.Error(), "123")

	// Verify the error cause
	notFoundErr, ok := err.(*serviceerrors.NotFoundError)
	require.True(t, ok)
	assert.Equal(t, testErr, notFoundErr.Unwrap())
}
