// Copyright (c) 2025 A Bit of Help, Inc.

package model

import (
	"github.com/abitofhelp/servicelib/valueobject/identification"
)

// Parent represents a parent in a family
type Parent struct {
	ID        identification.ID `json:"id"`
	FirstName string            `json:"firstName"`
	LastName  string            `json:"lastName"`
	BirthDate string            `json:"birthDate"`
	DeathDate *string           `json:"deathDate,omitempty"`
}

// Child represents a child in a family
type Child struct {
	ID        identification.ID `json:"id"`
	FirstName string            `json:"firstName"`
	LastName  string            `json:"lastName"`
	BirthDate string            `json:"birthDate"`
	DeathDate *string           `json:"deathDate,omitempty"`
}
