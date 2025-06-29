// Copyright (c) 2025 A Bit of Help, Inc.

// Package datewrapper provides a wrapper around the servicelib/date package
// to avoid violating the dependency inversion principle.
package datewrapper

import (
	"time"

	"github.com/abitofhelp/servicelib/date"
)

const (
	// StandardDateFormat is the standard date format used throughout the application.
	// This format follows the RFC3339 standard (2006-01-02T15:04:05Z07:00),
	// which is both human-readable and machine-parsable, and includes
	// date, time, and timezone information.
	StandardDateFormat = date.StandardDateFormat
)

// ParseDate parses a date string in the standard format (RFC3339).
// It converts a string representation of a date and time into a time.Time object.
// If the string cannot be parsed according to the StandardDateFormat,
// it returns a validation error with details about the failure.
//
// Parameters:
//   - dateStr: The date string to parse (e.g., "2023-04-15T14:30:00Z")
//
// Returns:
//   - time.Time: The parsed time if successful
//   - error: A validation error if parsing fails, nil otherwise
func ParseDate(dateStr string) (time.Time, error) {
	return date.ParseDate(dateStr)
}

// ParseOptionalDate parses an optional date string in the standard format (RFC3339).
// This function is designed to handle nullable date strings (represented as pointers).
// If the input pointer is nil, the function returns nil without attempting to parse.
// Otherwise, it behaves similarly to ParseDate, converting the string to a time.Time.
//
// This is particularly useful when working with optional date fields in APIs or databases.
//
// Parameters:
//   - dateStr: Pointer to the date string to parse (e.g., "2023-04-15T14:30:00Z"), can be nil
//
// Returns:
//   - *time.Time: Pointer to the parsed time if successful, nil if input is nil
//   - error: A validation error if parsing fails, nil otherwise
func ParseOptionalDate(dateStr *string) (*time.Time, error) {
	return date.ParseOptionalDate(dateStr)
}

// FormatDate formats a time.Time as a string in the standard format (RFC3339).
// It converts a time.Time object into a standardized string representation
// that can be used for display, storage, or transmission.
//
// This function ensures consistent date formatting throughout the application.
//
// Parameters:
//   - date: The time.Time value to format
//
// Returns:
//   - string: The formatted date string (e.g., "2023-04-15T14:30:00Z")
func FormatDate(d time.Time) string {
	return date.FormatDate(d)
}

// FormatOptionalDate formats an optional time.Time as a string in the standard format (RFC3339).
// This function is designed to handle nullable time values (represented as pointers).
// If the input pointer is nil, the function returns nil without attempting to format.
// Otherwise, it behaves similarly to FormatDate, converting the time.Time to a string.
//
// This is particularly useful when working with optional date fields in APIs or databases.
//
// Parameters:
//   - date: Pointer to the time.Time value to format, can be nil
//
// Returns:
//   - *string: Pointer to the formatted date string (e.g., "2023-04-15T14:30:00Z") if input is not nil,
//     nil if input is nil
func FormatOptionalDate(d *time.Time) *string {
	return date.FormatOptionalDate(d)
}