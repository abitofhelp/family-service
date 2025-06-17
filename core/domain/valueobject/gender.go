// Copyright (c) 2025 A Bit of Help, Inc.

package valueobject

import (
	"errors"
	"strings"
)

// Gender represents a person's gender value object
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

// NewGender creates a new Gender with validation
func NewGender(gender string) (Gender, error) {
	// Trim whitespace and convert to lowercase
	trimmedGender := strings.ToLower(strings.TrimSpace(gender))

	// Validate gender
	switch trimmedGender {
	case string(GenderMale), string(GenderFemale), string(GenderOther):
		return Gender(trimmedGender), nil
	default:
		return "", errors.New("invalid gender value")
	}
}

// String returns the string representation of the Gender
func (g Gender) String() string {
	return string(g)
}

// Equals checks if two Genders are equal
func (g Gender) Equals(other Gender) bool {
	return g == other
}

// IsEmpty checks if the Gender is empty
func (g Gender) IsEmpty() bool {
	return g == ""
}
