//go:build ignore
// +build ignore

// Copyright (c) 2025 A Bit of Help, Inc.

package main

import (
	"context"
	"fmt"
	domainerrors "github.com/abitofhelp/family-service/core/domain/errors"
	"github.com/abitofhelp/servicelib/errors"
)

func main() {
	// Create a context (unused in this test)
	_ = context.Background()

	// Part 1: Standard servicelib errors
	fmt.Println("=== Standard Error Types ===")

	// Test NewDatabaseError for PostgreSQL operations
	dbErr := errors.NewDatabaseError("failed to create families table", "create", "families", fmt.Errorf("some error"))
	fmt.Printf("DatabaseError: %v\n", dbErr)

	// Test NewDatabaseError for a query operation
	queryErr := errors.NewDatabaseError("failed to get family from PostgreSQL", "query", "families", fmt.Errorf("some error"))
	fmt.Printf("QueryError: %v\n", queryErr)

	// Test NewNotFoundError for a missing family
	notFoundErr := errors.NewNotFoundError("Family", "123", nil)
	fmt.Printf("NotFoundError: %v\n", notFoundErr)

	// Test NewValidationError for a required field
	validationErr := errors.NewValidationError("id is required", "id", nil)
	fmt.Printf("ValidationError: %v\n", validationErr)

	// Part 2: Domain-specific error codes
	fmt.Println("\n=== Domain-Specific Error Types ===")

	// Family structure errors
	tooManyParentsErr := domainerrors.NewFamilyTooManyParentsError("family cannot have more than two parents", nil)
	fmt.Printf("FamilyTooManyParentsError: %v\n", tooManyParentsErr)

	parentExistsErr := domainerrors.NewFamilyParentExistsError("parent already exists in family", nil)
	fmt.Printf("FamilyParentExistsError: %v\n", parentExistsErr)

	childExistsErr := domainerrors.NewFamilyChildExistsError("child already exists in family", nil)
	fmt.Printf("FamilyChildExistsError: %v\n", childExistsErr)

	// Family status errors
	notMarriedErr := domainerrors.NewFamilyNotMarriedError("only married families can divorce", nil)
	fmt.Printf("FamilyNotMarriedError: %v\n", notMarriedErr)

	divorceRequiresTwoErr := domainerrors.NewFamilyDivorceRequiresTwoParentsError("divorce requires exactly two parents", nil)
	fmt.Printf("FamilyDivorceRequiresTwoParentsError: %v\n", divorceRequiresTwoErr)

	// Person-related errors
	parentDeceasedErr := domainerrors.NewParentAlreadyDeceasedError("parent is already marked as deceased", nil)
	fmt.Printf("ParentAlreadyDeceasedError: %v\n", parentDeceasedErr)

	childDeceasedErr := domainerrors.NewChildAlreadyDeceasedError("child is already marked as deceased", nil)
	fmt.Printf("ChildAlreadyDeceasedError: %v\n", childDeceasedErr)

	// Demonstrating error handling
	fmt.Println("\n=== Error Handling Examples ===")

	// Example of handling domain-specific errors
	err := domainerrors.NewFamilyNotMarriedError("only married families can divorce", nil)
	fmt.Printf("Domain error example: %v\n", err)

	// Example of handling standard errors
	err = errors.NewValidationError("id is required", "id", nil)
	fmt.Printf("Validation error example: %v\n", err)

	// Demonstrating error type checking
	fmt.Println("\n=== Error Type Checking ===")

	// Create a domain-specific error
	err = domainerrors.NewFamilyTooManyParentsError("family cannot have more than two parents", nil)

	// Check if it's a domain error
	if _, ok := err.(*errors.DomainError); ok {
		fmt.Println("Error is a domain error")
	}

	// Create a validation error
	err = errors.NewValidationError("id is required", "id", nil)

	// Check if it's a validation error
	if _, ok := err.(*errors.ValidationError); ok {
		fmt.Println("Error is a validation error")
	}
}
