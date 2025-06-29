// Copyright (c) 2025 A Bit of Help, Inc.

// Package errors provides error types and interfaces for the domain layer
package errors

import (
	"errors"
	"fmt"
)

// Error is the base error interface
type Error interface {
	error
	Code() string
	Message() string
	Cause() error
}

// ValidationError represents a validation error
type ValidationError interface {
	Error
	Field() string
}

// DomainError represents a domain-specific error
type DomainError interface {
	Error
}

// DatabaseError represents a database-related error
type DatabaseError interface {
	Error
	Operation() string
	Table() string
}

// NotFoundError represents a resource not found error
type NotFoundError interface {
	Error
	ResourceType() string
	ResourceID() string
}

// baseError implements the Error interface
type baseError struct {
	code    string
	message string
	cause   error
}

// Code returns the error code
func (e *baseError) Code() string {
	return e.code
}

// Message returns the error message
func (e *baseError) Message() string {
	return e.message
}

// Cause returns the underlying cause of the error
func (e *baseError) Cause() error {
	return e.cause
}

// Error returns the error message
func (e *baseError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("[%s] %s: %s", e.code, e.message, e.cause.Error())
	}
	return fmt.Sprintf("[%s] %s", e.code, e.message)
}

// validationError implements the ValidationError interface
type validationError struct {
	baseError
	field string
}

// Field returns the field that failed validation
func (e *validationError) Field() string {
	return e.field
}

// domainError implements the DomainError interface
type domainError struct {
	baseError
}

// databaseError implements the DatabaseError interface
type databaseError struct {
	baseError
	operation string
	table     string
}

// Operation returns the database operation that failed
func (e *databaseError) Operation() string {
	return e.operation
}

// Table returns the database table that was being accessed
func (e *databaseError) Table() string {
	return e.table
}

// notFoundError implements the NotFoundError interface
type notFoundError struct {
	baseError
	resourceType string
	resourceID   string
}

// ResourceType returns the type of resource that was not found
func (e *notFoundError) ResourceType() string {
	return e.resourceType
}

// ResourceID returns the ID of the resource that was not found
func (e *notFoundError) ResourceID() string {
	return e.resourceID
}

// NewDomainError creates a new domain error with the given code, message, and cause
func NewDomainError(code string, message string, cause error) error {
	return &domainError{
		baseError: baseError{
			code:    code,
			message: message,
			cause:   cause,
		},
	}
}

// NewValidationError creates a new validation error with the given message, field, and cause
func NewValidationError(message string, field string, cause error) error {
	return &validationError{
		baseError: baseError{
			code:    "VALIDATION_ERROR",
			message: message,
			cause:   cause,
		},
		field: field,
	}
}

// NewDatabaseError creates a new database error with the given message, operation, table, and cause
func NewDatabaseError(message string, operation string, table string, cause error) error {
	return &databaseError{
		baseError: baseError{
			code:    "DATABASE_ERROR",
			message: message,
			cause:   cause,
		},
		operation: operation,
		table:     table,
	}
}

// NewNotFoundError creates a new not found error with the given resource type, resource ID, and cause
func NewNotFoundError(resourceType string, resourceID string, cause error) error {
	return &notFoundError{
		baseError: baseError{
			code:    "NOT_FOUND_ERROR",
			message: fmt.Sprintf("%s with ID %s not found", resourceType, resourceID),
			cause:   cause,
		},
		resourceType: resourceType,
		resourceID:   resourceID,
	}
}

// IsValidationError checks if the given error is a validation error
func IsValidationError(err error) bool {
	_, ok := err.(ValidationError)
	return ok
}

// IsDomainError checks if the given error is a domain error
func IsDomainError(err error) bool {
	_, ok := err.(DomainError)
	return ok
}

// IsDatabaseError checks if the given error is a database error
func IsDatabaseError(err error) bool {
	_, ok := err.(DatabaseError)
	return ok
}

// IsNotFoundError checks if the given error is a not found error
func IsNotFoundError(err error) bool {
	_, ok := err.(NotFoundError)
	return ok
}

// GetErrorCode extracts the error code from an error if it implements the Code() method
func GetErrorCode(err error) string {
	if err == nil {
		return ""
	}

	// Try to cast to our Error interface
	if e, ok := err.(Error); ok {
		return e.Code()
	}

	// Try to extract from wrapped error
	if cause := GetErrorCause(err); cause != nil && cause != err {
		return GetErrorCode(cause)
	}

	return ""
}

// GetErrorMessage extracts the error message from an error if it implements the Message() method
func GetErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	// Try to cast to our Error interface
	if e, ok := err.(Error); ok {
		return e.Message()
	}

	// Use the standard error message
	return err.Error()
}

// GetErrorCause extracts the cause from an error if it implements the Cause() method
func GetErrorCause(err error) error {
	if err == nil {
		return nil
	}

	// Try to cast to our Error interface
	if e, ok := err.(Error); ok {
		return e.Cause()
	}

	// Try to unwrap using standard errors package
	return errors.Unwrap(err)
}

// FormatError formats an error with its code, message, and cause
func FormatError(err error) string {
	if err == nil {
		return ""
	}

	code := GetErrorCode(err)
	message := GetErrorMessage(err)
	cause := GetErrorCause(err)

	if code != "" {
		if cause != nil {
			return fmt.Sprintf("[%s] %s: %s", code, message, cause.Error())
		}
		return fmt.Sprintf("[%s] %s", code, message)
	}

	return message
}

// WrapError wraps an error with a message
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// Domain-specific error types

// FamilyTooManyParentsError represents an error when a family has too many parents
type FamilyTooManyParentsError struct {
	baseError
}

// NewFamilyTooManyParentsError creates a new FamilyTooManyParentsError
func NewFamilyTooManyParentsError(message string, cause error) error {
	return &FamilyTooManyParentsError{
		baseError: baseError{
			code:    "FAMILY_TOO_MANY_PARENTS",
			message: message,
			cause:   cause,
		},
	}
}

// FamilyParentExistsError represents an error when a parent already exists in a family
type FamilyParentExistsError struct {
	baseError
}

// NewFamilyParentExistsError creates a new FamilyParentExistsError
func NewFamilyParentExistsError(message string, cause error) error {
	return &FamilyParentExistsError{
		baseError: baseError{
			code:    "FAMILY_PARENT_EXISTS",
			message: message,
			cause:   cause,
		},
	}
}

// FamilyParentDuplicateError represents an error when a duplicate parent is added to a family
type FamilyParentDuplicateError struct {
	baseError
}

// NewFamilyParentDuplicateError creates a new FamilyParentDuplicateError
func NewFamilyParentDuplicateError(message string, cause error) error {
	return &FamilyParentDuplicateError{
		baseError: baseError{
			code:    "FAMILY_PARENT_DUPLICATE",
			message: message,
			cause:   cause,
		},
	}
}

// FamilyChildExistsError represents an error when a child already exists in a family
type FamilyChildExistsError struct {
	baseError
}

// NewFamilyChildExistsError creates a new FamilyChildExistsError
func NewFamilyChildExistsError(message string, cause error) error {
	return &FamilyChildExistsError{
		baseError: baseError{
			code:    "FAMILY_CHILD_EXISTS",
			message: message,
			cause:   cause,
		},
	}
}

// FamilyCannotRemoveLastParentError represents an error when trying to remove the last parent from a family
type FamilyCannotRemoveLastParentError struct {
	baseError
}

// NewFamilyCannotRemoveLastParentError creates a new FamilyCannotRemoveLastParentError
func NewFamilyCannotRemoveLastParentError(message string, cause error) error {
	return &FamilyCannotRemoveLastParentError{
		baseError: baseError{
			code:    "FAMILY_CANNOT_REMOVE_LAST_PARENT",
			message: message,
			cause:   cause,
		},
	}
}

// FamilyNotMarriedError represents an error when trying to divorce a family that is not married
type FamilyNotMarriedError struct {
	baseError
}

// NewFamilyNotMarriedError creates a new FamilyNotMarriedError
func NewFamilyNotMarriedError(message string, cause error) error {
	return &FamilyNotMarriedError{
		baseError: baseError{
			code:    "FAMILY_NOT_MARRIED",
			message: message,
			cause:   cause,
		},
	}
}

// FamilyDivorceRequiresTwoParentsError represents an error when trying to divorce a family that doesn't have two parents
type FamilyDivorceRequiresTwoParentsError struct {
	baseError
}

// NewFamilyDivorceRequiresTwoParentsError creates a new FamilyDivorceRequiresTwoParentsError
func NewFamilyDivorceRequiresTwoParentsError(message string, cause error) error {
	return &FamilyDivorceRequiresTwoParentsError{
		baseError: baseError{
			code:    "FAMILY_DIVORCE_REQUIRES_TWO_PARENTS",
			message: message,
			cause:   cause,
		},
	}
}

// FamilyCreateFailedError represents an error when creating a family fails
type FamilyCreateFailedError struct {
	baseError
}

// NewFamilyCreateFailedError creates a new FamilyCreateFailedError
func NewFamilyCreateFailedError(message string, cause error) error {
	return &FamilyCreateFailedError{
		baseError: baseError{
			code:    "FAMILY_CREATE_FAILED",
			message: message,
			cause:   cause,
		},
	}
}