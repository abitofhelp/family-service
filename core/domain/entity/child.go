package entity

import (
	"time"

	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/validation"
)

// Child represents a child entity in the family domain
type Child struct {
	id        string
	firstName string
	lastName  string
	birthDate time.Time
	deathDate *time.Time
}

// NewChild creates a new Child entity with validation
func NewChild(id, firstName, lastName string, birthDate time.Time, deathDate *time.Time) (*Child, error) {
	c := &Child{
		id:        id,
		firstName: firstName,
		lastName:  lastName,
		birthDate: birthDate,
		deathDate: deathDate,
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return c, nil
}

// Validate ensures the Child entity is valid
func (c *Child) Validate() error {
	result := validation.NewValidationResult()

	validation.ValidateID(c.id, "ID", result)
	validation.Required(c.firstName, "FirstName", result)
	validation.Required(c.lastName, "LastName", result)
	validation.MinLength(c.firstName, 2, "FirstName", result)
	validation.MinLength(c.lastName, 2, "LastName", result)
	validation.PastDate(c.birthDate, "BirthDate", result)

	if c.deathDate != nil {
		validation.PastDate(*c.deathDate, "DeathDate", result)
		validation.ValidDateRange(c.birthDate, *c.deathDate, "BirthDate", "DeathDate", result)
	}

	return result.Error()
}

// ID returns the child's ID
func (c *Child) ID() string {
	return c.id
}

// FirstName returns the child's first name
func (c *Child) FirstName() string {
	return c.firstName
}

// LastName returns the child's last name
func (c *Child) LastName() string {
	return c.lastName
}

// BirthDate returns the child's birth date
func (c *Child) BirthDate() time.Time {
	return c.birthDate
}

// DeathDate returns the child's death date
func (c *Child) DeathDate() *time.Time {
	if c.deathDate == nil {
		return nil
	}
	// Return a copy to prevent modification
	copy := *c.deathDate
	return &copy
}

// FullName returns the child's full name
func (c *Child) FullName() string {
	return c.firstName + " " + c.lastName
}

// IsDeceased returns true if the child is deceased
func (c *Child) IsDeceased() bool {
	return c.deathDate != nil
}

// MarkDeceased marks the child as deceased with the given death date
func (c *Child) MarkDeceased(deathDate time.Time) error {
	if c.deathDate != nil {
		return errors.NewDomainError(nil, "child is already marked as deceased", "ALREADY_DECEASED")
	}

	if !deathDate.After(c.birthDate) {
		return errors.NewValidationError("death date must be after birth date")
	}

	if deathDate.After(time.Now()) {
		return errors.NewValidationError("death date cannot be in the future")
	}

	c.deathDate = &deathDate
	return nil
}

// Equals checks if two children are the same based on ID
func (c *Child) Equals(other *Child) bool {
	if other == nil {
		return false
	}
	return c.id == other.id
}

// ToDTO converts the Child entity to a data transfer object for external use
func (c *Child) ToDTO() ChildDTO {
	dto := ChildDTO{
		ID:        c.id,
		FirstName: c.firstName,
		LastName:  c.lastName,
		BirthDate: c.birthDate,
		DeathDate: c.deathDate,
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
