// Copyright (c) 2025 A Bit of Help, Inc.

// Package errors provides a wrapper around servicelib/errors to ensure
// the domain layer doesn't directly depend on external libraries.
package errors

import (
	"errors"
	"fmt"

	serviceerrors "github.com/abitofhelp/servicelib/errors"
)

// Error types and interfaces to match servicelib error types

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

// Wrapper functions for servicelib/errors

// NewDomainError creates a new domain error with the given code, message, and cause
func NewDomainError(code string, message string, cause error) error {
	// Convert string code to serviceerrors.ErrorCode
	errorCode := serviceerrors.ErrorCode(code)
	return serviceerrors.NewDomainError(errorCode, message, cause)
}

// NewValidationError creates a new validation error with the given message, field, and cause
func NewValidationError(message string, field string, cause error) error {
	return serviceerrors.NewValidationError(message, field, cause)
}

// NewDatabaseError creates a new database error with the given message, operation, table, and cause
func NewDatabaseError(message string, operation string, table string, cause error) error {
	return serviceerrors.NewDatabaseError(message, operation, table, cause)
}

// NewNotFoundError creates a new not found error with the given resource type, resource ID, and cause
func NewNotFoundError(resourceType string, resourceID string, cause error) error {
	return serviceerrors.NewNotFoundError(resourceType, resourceID, cause)
}

// IsValidationError checks if the given error is a validation error
func IsValidationError(err error) bool {
	return serviceerrors.IsValidationError(err)
}

// IsDomainError checks if the given error is a domain error
func IsDomainError(err error) bool {
	return serviceerrors.IsDomainError(err)
}

// IsDatabaseError checks if the given error is a database error
func IsDatabaseError(err error) bool {
	return serviceerrors.IsDatabaseError(err)
}

// IsNotFoundError checks if the given error is a not found error
func IsNotFoundError(err error) bool {
	return serviceerrors.IsNotFoundError(err)
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
