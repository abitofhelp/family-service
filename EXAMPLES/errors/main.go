//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

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
