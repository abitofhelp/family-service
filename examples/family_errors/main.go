// Copyright (c) 2025 A Bit of Help, Inc.

package main

import (
	"context"
	"fmt"
	"github.com/abitofhelp/servicelib/errors"
	"github.com/jackc/pgx/v5"
)

func main() {
	// Create a context (unused in this test)
	_ = context.Background()

	// Test NewDatabaseError for PostgreSQL operations
	dbErr := errors.NewDatabaseError("failed to create families table", "create", "families", fmt.Errorf("some error"))
	fmt.Printf("DatabaseError: %v\n", dbErr)

	// Test NewDatabaseError for a query operation
	queryErr := errors.NewDatabaseError("failed to get family from PostgreSQL", "query", "families", fmt.Errorf("some error"))
	fmt.Printf("QueryError: %v\n", queryErr)

	// Test NewDatabaseError for a JSON unmarshaling error
	jsonErr := errors.NewDatabaseError("failed to unmarshal parents data", "unmarshal", "families", fmt.Errorf("some error"))
	fmt.Printf("JSONError: %v\n", jsonErr)

	// Test NewNotFoundError for a missing family
	notFoundErr := errors.NewNotFoundError("Family", "123", nil)
	fmt.Printf("NotFoundError: %v\n", notFoundErr)

	// Test NewNotFoundError for a missing family with a cause
	notFoundWithCauseErr := errors.NewNotFoundError("Family", "123", pgx.ErrNoRows)
	fmt.Printf("NotFoundWithCauseError: %v\n", notFoundWithCauseErr)

	// Test NewValidationError for a required field
	validationErr := errors.NewValidationError("id is required", "id", nil)
	fmt.Printf("ValidationError: %v\n", validationErr)

	// Test NewValidationError for a family validation
	familyValidationErr := errors.NewValidationError("family cannot be nil", "family", nil)
	fmt.Printf("FamilyValidationError: %v\n", familyValidationErr)
}
