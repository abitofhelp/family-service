// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"errors"
	"github.com/google/uuid"
)

// ID represents a unique identifier value object
type ID string

// NewID creates a new ID with validation
func NewID(id string) (ID, error) {
	if id == "" {
		return "", errors.New("ID cannot be empty")
	}

	// Validate UUID format if not empty
	if _, err := uuid.Parse(id); err != nil {
		return "", errors.New("invalid ID format")
	}

	return ID(id), nil
}

// GenerateID creates a new random ID
func GenerateID() ID {
	return ID(uuid.New().String())
}

// String returns the string representation of the ID
func (id ID) String() string {
	return string(id)
}

// Equals checks if two IDs are equal
func (id ID) Equals(other ID) bool {
	return id == other
}

// IsEmpty checks if the ID is empty
func (id ID) IsEmpty() bool {
	return id == ""
}
