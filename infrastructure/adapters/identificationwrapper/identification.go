// Copyright (c) 2025 A Bit of Help, Inc.

// Package identificationwrapper provides a wrapper around servicelib/valueobject/identification to ensure
// the domain layer doesn't directly depend on external libraries.
package identificationwrapper

import (
	"time"

	"github.com/abitofhelp/family-service/infrastructure/adapters/errorswrapper"
	"github.com/abitofhelp/servicelib/valueobject/identification"
	"github.com/google/uuid"
)

// ID represents a unique identifier
type ID string

// NewID creates a new ID with a random UUID
func NewID() (ID, error) {
	// Generate a random UUID
	id := uuid.New().String()
	// Use the servicelib function to create an ID
	serviceID, err := identification.NewID(id)
	if err != nil {
		return "", err
	}
	return ID(serviceID), nil
}

// NewIDFromString creates a new ID from a string
func NewIDFromString(id string) (ID, error) {
	if id == "" {
		return "", errorswrapper.NewValidationError("ID cannot be empty", "ID", nil)
	}

	// Validate that the ID is a valid UUID
	_, err := uuid.Parse(id)
	if err != nil {
		return "", errorswrapper.NewValidationError("invalid ID format", "ID", err)
	}

	return ID(id), nil
}

// String returns the string representation of the ID
func (id ID) String() string {
	return string(id)
}

// IsEmpty checks if the ID is empty
func (id ID) IsEmpty() bool {
	return id == ""
}

// Equals checks if the ID equals another ID
func (id ID) Equals(other ID) bool {
	return id == other
}

// MarshalJSON implements the json.Marshaler interface
func (id ID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + id.String() + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (id *ID) UnmarshalJSON(data []byte) error {
	// Remove quotes
	if len(data) >= 2 && data[0] == '"' && data[len(data)-1] == '"' {
		data = data[1 : len(data)-1]
	}

	// Validate that the ID is a valid UUID
	_, err := uuid.Parse(string(data))
	if err != nil {
		return errorswrapper.NewValidationError("invalid ID format", "ID", err)
	}

	*id = ID(data)
	return nil
}

// Name represents a person's name
type Name string

// NewName creates a new Name from a string
func NewName(name string) (Name, error) {
	if name == "" {
		return "", errorswrapper.NewValidationError("name cannot be empty", "Name", nil)
	}
	return Name(name), nil
}

// String returns the string representation of the Name
func (n Name) String() string {
	return string(n)
}

// DateOfBirth represents a person's birth date
type DateOfBirth struct {
	date time.Time
}

// NewDateOfBirth creates a new DateOfBirth
func NewDateOfBirth(year, month, day int) (DateOfBirth, error) {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	// Validate the date
	if date.Year() != year || int(date.Month()) != month || date.Day() != day {
		return DateOfBirth{}, errorswrapper.NewValidationError("invalid date", "DateOfBirth", nil)
	}

	// Validate that the date is not in the future
	if date.After(time.Now()) {
		return DateOfBirth{}, errorswrapper.NewValidationError("birth date cannot be in the future", "DateOfBirth", nil)
	}

	return DateOfBirth{date: date}, nil
}

// Date returns the time.Time representation of the DateOfBirth
func (d DateOfBirth) Date() time.Time {
	return d.date
}

// DateOfDeath represents a person's death date
type DateOfDeath struct {
	date time.Time
}

// NewDateOfDeath creates a new DateOfDeath
func NewDateOfDeath(year, month, day int) (DateOfDeath, error) {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	// Validate the date
	if date.Year() != year || int(date.Month()) != month || date.Day() != day {
		return DateOfDeath{}, errorswrapper.NewValidationError("invalid date", "DateOfDeath", nil)
	}

	// Validate that the date is not in the future
	if date.After(time.Now()) {
		return DateOfDeath{}, errorswrapper.NewValidationError("death date cannot be in the future", "DateOfDeath", nil)
	}

	return DateOfDeath{date: date}, nil
}

// Date returns the time.Time representation of the DateOfDeath
func (d DateOfDeath) Date() time.Time {
	return d.date
}
