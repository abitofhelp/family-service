// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"errors"
	"strings"
)

// Name represents a person's name value object
type Name string

// NewName creates a new Name with validation
func NewName(name string) (Name, error) {
	// Trim whitespace
	trimmedName := strings.TrimSpace(name)

	// Validate name is not empty
	if trimmedName == "" {
		return "", errors.New("name cannot be empty")
	}

	// Validate name length
	if len(trimmedName) > 100 {
		return "", errors.New("name is too long")
	}

	return Name(trimmedName), nil
}

// String returns the string representation of the Name
func (n Name) String() string {
	return string(n)
}

// Equals checks if two Names are equal
func (n Name) Equals(other Name) bool {
	return strings.EqualFold(string(n), string(other))
}

// IsEmpty checks if the Name is empty
func (n Name) IsEmpty() bool {
	return n == ""
}
