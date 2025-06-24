// Copyright (c) 2025 A Bit of Help, Inc.

// Package errors provides domain-specific error codes and error handling utilities.
package errors

import (
	"github.com/abitofhelp/servicelib/errors"
)

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

// NewFamilyTooManyParentsError creates a new domain error for when a family has too many parents
func NewFamilyTooManyParentsError(message string, cause error) error {
	return errors.NewDomainError(FamilyTooManyParentsCode, message, cause)
}

// NewFamilyParentExistsError creates a new domain error for when a parent already exists in a family
func NewFamilyParentExistsError(message string, cause error) error {
	return errors.NewDomainError(FamilyParentExistsCode, message, cause)
}

// NewFamilyParentDuplicateError creates a new domain error for when a parent with the same name and birthdate already exists
func NewFamilyParentDuplicateError(message string, cause error) error {
	return errors.NewDomainError(FamilyParentDuplicateCode, message, cause)
}

// NewFamilyChildExistsError creates a new domain error for when a child already exists in a family
func NewFamilyChildExistsError(message string, cause error) error {
	return errors.NewDomainError(FamilyChildExistsCode, message, cause)
}

// NewFamilyCannotRemoveLastParentError creates a new domain error for when attempting to remove the only parent
func NewFamilyCannotRemoveLastParentError(message string, cause error) error {
	return errors.NewDomainError(FamilyCannotRemoveLastParent, message, cause)
}

// NewFamilyNotMarriedError creates a new domain error for when a family is not in a married state
func NewFamilyNotMarriedError(message string, cause error) error {
	return errors.NewDomainError(FamilyNotMarriedCode, message, cause)
}

// NewFamilyDivorceRequiresTwoParentsError creates a new domain error for when a divorce is attempted without two parents
func NewFamilyDivorceRequiresTwoParentsError(message string, cause error) error {
	return errors.NewDomainError(FamilyDivorceRequiresTwoCode, message, cause)
}

// NewFamilyCreateFailedError creates a new domain error for when family creation fails
func NewFamilyCreateFailedError(message string, cause error) error {
	return errors.NewDomainError(FamilyCreateFailedCode, message, cause)
}

// NewParentAlreadyDeceasedError creates a new domain error for when a parent is already marked as deceased
func NewParentAlreadyDeceasedError(message string, cause error) error {
	return errors.NewDomainError(ParentAlreadyDeceasedCode, message, cause)
}

// NewChildAlreadyDeceasedError creates a new domain error for when a child is already marked as deceased
func NewChildAlreadyDeceasedError(message string, cause error) error {
	return errors.NewDomainError(ChildAlreadyDeceasedCode, message, cause)
}

// NewFamilyStatusUpdateFailedError creates a new domain error for when family status update fails
func NewFamilyStatusUpdateFailedError(message string, cause error) error {
	return errors.NewDomainError(FamilyStatusUpdateFailedCode, message, cause)
}
