// Copyright (c) 2025 A Bit of Help, Inc.

package entity

import (
	"time"

	domainerrors "github.com/abitofhelp/family-service/core/domain/errors"
	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/validation"
	"github.com/abitofhelp/servicelib/valueobject/identification"
)

// Parent represents a parent entity in the family domain
type Parent struct {
	id        identification.ID
	firstName identification.Name
	lastName  identification.Name
	birthDate identification.DateOfBirth
	deathDate *identification.DateOfDeath
}

// NewParent creates a new Parent entity with validation
func NewParent(id, firstName, lastName string, birthDate time.Time, deathDate *time.Time) (*Parent, error) {
	// Create ID value object
	idVO, err := identification.NewID(id)
	if err != nil {
		return nil, errors.NewValidationError("invalid ID: "+err.Error(), "ID", err)
	}

	// Create FirstName value object
	firstNameVO, err := identification.NewName(firstName)
	if err != nil {
		return nil, errors.NewValidationError("invalid FirstName: "+err.Error(), "FirstName", err)
	}

	// Create LastName value object
	lastNameVO, err := identification.NewName(lastName)
	if err != nil {
		return nil, errors.NewValidationError("invalid LastName: "+err.Error(), "LastName", err)
	}

	// Create BirthDate value object
	year, month, day := birthDate.Date()
	birthDateVO, err := identification.NewDateOfBirth(year, int(month), day)
	if err != nil {
		return nil, errors.NewValidationError("invalid BirthDate: "+err.Error(), "BirthDate", err)
	}

	// Create DeathDate value object if provided
	var deathDateVO *identification.DateOfDeath
	if deathDate != nil {
		year, month, day := deathDate.Date()
		dod, err := identification.NewDateOfDeath(year, int(month), day)
		if err != nil {
			return nil, errors.NewValidationError("invalid DeathDate: "+err.Error(), "DeathDate", err)
		}
		deathDateVO = &dod
	}

	p := &Parent{
		id:        idVO,
		firstName: firstNameVO,
		lastName:  lastNameVO,
		birthDate: birthDateVO,
		deathDate: deathDateVO,
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	return p, nil
}

// Validate ensures the Parent entity is valid
func (p *Parent) Validate() error {
	result := validation.NewValidationResult()

	// Value objects already have their own validation, but we can add additional validation here
	// For example, we can validate that the death date is after the birth date
	if p.deathDate != nil && !p.deathDate.Date().After(p.birthDate.Date()) {
		result.AddError("death date must be after birth date", "DeathDate")
	}

	// Validate minimum length for names (value objects only validate that they're not empty)
	if len(p.firstName.String()) < 2 {
		result.AddError("must be at least 2 characters long", "FirstName")
	}

	if len(p.lastName.String()) < 2 {
		result.AddError("must be at least 2 characters long", "LastName")
	}

	return result.Error()
}

// ID returns the parent's ID
func (p *Parent) ID() string {
	return p.id.String()
}

// FirstName returns the parent's first name
func (p *Parent) FirstName() string {
	return p.firstName.String()
}

// LastName returns the parent's last name
func (p *Parent) LastName() string {
	return p.lastName.String()
}

// BirthDate returns the parent's birth date
func (p *Parent) BirthDate() time.Time {
	return p.birthDate.Date()
}

// DeathDate returns the parent's death date
func (p *Parent) DeathDate() *time.Time {
	if p.deathDate == nil {
		return nil
	}
	// Return a copy to prevent modification
	date := p.deathDate.Date()
	return &date
}

// FullName returns the parent's full name
func (p *Parent) FullName() string {
	return p.firstName.String() + " " + p.lastName.String()
}

// IsDeceased returns true if the parent is deceased
func (p *Parent) IsDeceased() bool {
	return p.deathDate != nil
}

// MarkDeceased marks the parent as deceased with the given death date
func (p *Parent) MarkDeceased(deathDate time.Time) error {
	if p.deathDate != nil {
		return domainerrors.NewParentAlreadyDeceasedError("parent is already marked as deceased", nil)
	}

	// Validate death date is after birth date
	if !deathDate.After(p.birthDate.Date()) {
		return errors.NewValidationError("death date must be after birth date", "DeathDate", nil)
	}

	// Create DateOfDeath value object
	year, month, day := deathDate.Date()
	dod, err := identification.NewDateOfDeath(year, int(month), day)
	if err != nil {
		return errors.NewValidationError("invalid death date: "+err.Error(), "DeathDate", err)
	}

	p.deathDate = &dod
	return nil
}

// Equals checks if two parents are the same based on ID
func (p *Parent) Equals(other *Parent) bool {
	if other == nil {
		return false
	}
	return p.id.Equals(other.id)
}

// ToDTO converts the Parent entity to a data transfer object for external use
func (p *Parent) ToDTO() ParentDTO {
	var deathDate *time.Time
	if p.deathDate != nil {
		date := p.deathDate.Date()
		deathDate = &date
	}

	dto := ParentDTO{
		ID:        p.id.String(),
		FirstName: p.firstName.String(),
		LastName:  p.lastName.String(),
		BirthDate: p.birthDate.Date(),
		DeathDate: deathDate,
	}
	return dto
}

// ParentDTO is a data transfer object for the Parent entity
type ParentDTO struct {
	ID        string
	FirstName string
	LastName  string
	BirthDate time.Time
	DeathDate *time.Time
}

// ParentFromDTO creates a Parent entity from a data transfer object
func ParentFromDTO(dto ParentDTO) (*Parent, error) {
	return NewParent(dto.ID, dto.FirstName, dto.LastName, dto.BirthDate, dto.DeathDate)
}
