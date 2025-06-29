// Copyright (c) 2025 A Bit of Help, Inc.

// Package errors provides domain-specific error codes and error handling utilities.
package errors

// Domain-specific error codes for family-related errors
const (
	// Family structure errors
	FamilyTooManyParentsCode     = "FAMILY_TOO_MANY_PARENTS"
	FamilyParentExistsCode       = "FAMILY_PARENT_EXISTS"
	FamilyParentDuplicateCode    = "FAMILY_PARENT_DUPLICATE"
	FamilyChildExistsCode        = "FAMILY_CHILD_EXISTS"
	FamilyCannotRemoveLastParent = "FAMILY_CANNOT_REMOVE_LAST_PARENT"

 // Family status errors
	FamilyNotMarriedCode         = "FAMILY_NOT_MARRIED"
	FamilyDivorceRequiresTwoCode = "FAMILY_DIVORCE_REQUIRES_TWO_PARENTS"
	FamilyCreateFailedCode       = "FAMILY_CREATE_FAILED"
	FamilyStatusUpdateFailedCode = "FAMILY_STATUS_UPDATE_FAILED"

	// Parent-related errors
	ParentAlreadyDeceasedCode = "PARENT_ALREADY_DECEASED"

	// Child-related errors
	ChildAlreadyDeceasedCode = "CHILD_ALREADY_DECEASED"
)

// ParentAlreadyDeceasedError represents an error when a parent is already marked as deceased
type ParentAlreadyDeceasedError struct {
	baseError
}

// NewParentAlreadyDeceasedError creates a new ParentAlreadyDeceasedError
func NewParentAlreadyDeceasedError(message string, cause error) error {
	return &ParentAlreadyDeceasedError{
		baseError: baseError{
			code:    ParentAlreadyDeceasedCode,
			message: message,
			cause:   cause,
		},
	}
}

// ChildAlreadyDeceasedError represents an error when a child is already marked as deceased
type ChildAlreadyDeceasedError struct {
	baseError
}

// NewChildAlreadyDeceasedError creates a new ChildAlreadyDeceasedError
func NewChildAlreadyDeceasedError(message string, cause error) error {
	return &ChildAlreadyDeceasedError{
		baseError: baseError{
			code:    ChildAlreadyDeceasedCode,
			message: message,
			cause:   cause,
		},
	}
}

// FamilyStatusUpdateFailedError represents an error when family status update fails
type FamilyStatusUpdateFailedError struct {
	baseError
}

// NewFamilyStatusUpdateFailedError creates a new FamilyStatusUpdateFailedError
func NewFamilyStatusUpdateFailedError(message string, cause error) error {
	return &FamilyStatusUpdateFailedError{
		baseError: baseError{
			code:    FamilyStatusUpdateFailedCode,
			message: message,
			cause:   cause,
		},
	}
}
