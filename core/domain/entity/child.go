// Copyright (c) 2025 A Bit of Help, Inc.

package entity

import (
	"fmt"
	"regexp"
	"time"

	domainerrors "github.com/abitofhelp/family-service/core/domain/errors"
	"github.com/abitofhelp/family-service/infrastructure/adapters/errorswrapper"
	"github.com/abitofhelp/family-service/infrastructure/adapters/identificationwrapper"
	"github.com/abitofhelp/family-service/infrastructure/adapters/validationwrapper"
)

// Child represents a child entity in the family domain
type Child struct {
	id        identificationwrapper.ID
	firstName identificationwrapper.Name
	lastName  identificationwrapper.Name
	birthDate identificationwrapper.DateOfBirth
	deathDate *identificationwrapper.DateOfDeath
}

// NewChild creates a new Child entity with validation
func NewChild(id, firstName, lastName string, birthDate time.Time, deathDate *time.Time) (*Child, error) {
	// Create ID value object
	idVO, err := identificationwrapper.NewIDFromString(id)
	if err != nil {
		return nil, errorswrapper.NewValidationError("invalid ID: "+err.Error(), "ID", err)
	}

	// Create FirstName value object
	firstNameVO, err := identificationwrapper.NewName(firstName)
	if err != nil {
		return nil, errorswrapper.NewValidationError("invalid FirstName: "+err.Error(), "FirstName", err)
	}

	// Create LastName value object
	lastNameVO, err := identificationwrapper.NewName(lastName)
	if err != nil {
		return nil, errorswrapper.NewValidationError("invalid LastName: "+err.Error(), "LastName", err)
	}

	// Create BirthDate value object
	year, month, day := birthDate.Date()
	birthDateVO, err := identificationwrapper.NewDateOfBirth(year, int(month), day)
	if err != nil {
		return nil, errorswrapper.NewValidationError("invalid BirthDate: "+err.Error(), "BirthDate", err)
	}

	// Create DeathDate value object if provided
	var deathDateVO *identificationwrapper.DateOfDeath
	if deathDate != nil {
		year, month, day := deathDate.Date()
		dod, err := identificationwrapper.NewDateOfDeath(year, int(month), day)
		if err != nil {
			return nil, errorswrapper.NewValidationError("invalid DeathDate: "+err.Error(), "DeathDate", err)
		}
		deathDateVO = &dod
	}

	c := &Child{
		id:        idVO,
		firstName: firstNameVO,
		lastName:  lastNameVO,
		birthDate: birthDateVO,
		deathDate: deathDateVO,
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return c, nil
}

// Validate ensures the Child entity is valid
func (c *Child) Validate() error {
	result := validationwrapper.NewValidationResult()

	// Value objects already have their own validation, but we can add additional validation here
	// For example, we can validate that the death date is after the birth date
	if c.deathDate != nil && !c.deathDate.Date().After(c.birthDate.Date()) {
		result.AddError("death date must be after birth date", "DeathDate")
	}

	// Validate minimum length for names (value objects only validate that they're not empty)
	if len(c.firstName.String()) < 2 {
		result.AddError("must be at least 2 characters long", "FirstName")
	}

	if len(c.lastName.String()) < 2 {
		result.AddError("must be at least 2 characters long", "LastName")
	}

	// Enhanced validation: Validate name doesn't contain special characters
	nameRegex := "^[a-zA-Z\\s-]+$"
	if matched, _ := regexp.MatchString(nameRegex, c.firstName.String()); !matched {
		result.AddError("must contain only letters, spaces, and hyphens", "FirstName")
	}

	if matched, _ := regexp.MatchString(nameRegex, c.lastName.String()); !matched {
		result.AddError("must contain only letters, spaces, and hyphens", "LastName")
	}

	// Enhanced validation: Validate birth date is not in the future
	if c.birthDate.Date().After(time.Now()) {
		result.AddError("birth date cannot be in the future", "BirthDate")
	}

	// Enhanced validation: Validate death date is not in the future
	if c.deathDate != nil && c.deathDate.Date().After(time.Now()) {
		result.AddError("death date cannot be in the future", "DeathDate")
	}

	// Enhanced validation: Validate maximum age (e.g., 150 years)
	maxAge := 150
	now := time.Now()
	birthDate := c.birthDate.Date()
	age := now.Year() - birthDate.Year()

	// Adjust age if birthday hasn't occurred yet this year
	if now.Month() < birthDate.Month() || (now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}

	if age > maxAge {
		result.AddError(fmt.Sprintf("age cannot exceed %d years", maxAge), "BirthDate")
	}

	return result.Error()
}

// ID returns the child's ID
func (c *Child) ID() string {
	return c.id.String()
}

// FirstName returns the child's first name
func (c *Child) FirstName() string {
	return c.firstName.String()
}

// LastName returns the child's last name
func (c *Child) LastName() string {
	return c.lastName.String()
}

// BirthDate returns the child's birth date
func (c *Child) BirthDate() time.Time {
	return c.birthDate.Date()
}

// DeathDate returns the child's death date
func (c *Child) DeathDate() *time.Time {
	if c.deathDate == nil {
		return nil
	}
	// Return a copy to prevent modification
	date := c.deathDate.Date()
	return &date
}

// FullName returns the child's full name
func (c *Child) FullName() string {
	return c.firstName.String() + " " + c.lastName.String()
}

// IsDeceased returns true if the child is deceased
func (c *Child) IsDeceased() bool {
	return c.deathDate != nil
}

// MarkDeceased marks the child as deceased with the given death date
func (c *Child) MarkDeceased(deathDate time.Time) error {
	if c.deathDate != nil {
		return domainerrors.NewChildAlreadyDeceasedError("child is already marked as deceased", nil)
	}

	// Validate death date is after birth date
	if !deathDate.After(c.birthDate.Date()) {
		return errorswrapper.NewValidationError("death date must be after birth date", "DeathDate", nil)
	}

	// Create DateOfDeath value object
	year, month, day := deathDate.Date()
	dod, err := identificationwrapper.NewDateOfDeath(year, int(month), day)
	if err != nil {
		return errorswrapper.NewValidationError("invalid death date: "+err.Error(), "DeathDate", err)
	}

	c.deathDate = &dod
	return nil
}

// Equals checks if two children are the same based on ID
func (c *Child) Equals(other *Child) bool {
	if other == nil {
		return false
	}
	return c.id.Equals(other.id)
}

// ToDTO converts the Child entity to a data transfer object for external use
func (c *Child) ToDTO() ChildDTO {
	var deathDate *time.Time
	if c.deathDate != nil {
		date := c.deathDate.Date()
		deathDate = &date
	}

	dto := ChildDTO{
		ID:        c.id.String(),
		FirstName: c.firstName.String(),
		LastName:  c.lastName.String(),
		BirthDate: c.birthDate.Date(),
		DeathDate: deathDate,
	}
	return dto
}

// ChildDTO is a data transfer object for the Child entity
type ChildDTO struct {
	ID        string
	FirstName string
	LastName  string
	BirthDate time.Time
	DeathDate *time.Time
}

// ChildFromDTO creates a Child entity from a data transfer object
func ChildFromDTO(dto ChildDTO) (*Child, error) {
	return NewChild(dto.ID, dto.FirstName, dto.LastName, dto.BirthDate, dto.DeathDate)
}