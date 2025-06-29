//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

// Package main demonstrates how to use the error types from the servicelib/errors package.
//
// This example shows how to:
// - Create and use DatabaseError for database operation failures
// - Create and use NotFoundError for resource not found situations
// - Create and use ValidationError for input validation failures
//
// Each error type provides structured error information that can be used
// for logging, error handling, and generating appropriate responses to clients.
package main

import (
	"fmt"
	"github.com/abitofhelp/servicelib/errors"
)

func main() {
	// Test NewDatabaseError
	dbErr := errors.NewDatabaseError("failed to create table", "create", "families", fmt.Errorf("some error"))
	fmt.Printf("DatabaseError: %v\n", dbErr)

	// Test NewNotFoundError
	notFoundErr := errors.NewNotFoundError("Family", "123", fmt.Errorf("some error"))
	fmt.Printf("NotFoundError: %v\n", notFoundErr)

	// Test NewValidationError
	validationErr := errors.NewValidationError("invalid input", "field", fmt.Errorf("some error"))
	fmt.Printf("ValidationError: %v\n", validationErr)
}
