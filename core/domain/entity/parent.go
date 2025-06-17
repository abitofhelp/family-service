package entity

import (
	"time"

	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/validation"
)

// Parent represents a parent entity in the family domain
type Parent struct {
	id        string
	firstName string
	lastName  string
	birthDate time.Time
	deathDate *time.Time
}

// NewParent creates a new Parent entity with validation
func NewParent(id, firstName, lastName string, birthDate time.Time, deathDate *time.Time) (*Parent, error) {
	p := &Parent{
		id:        id,
		firstName: firstName,
		lastName:  lastName,
		birthDate: birthDate,
		deathDate: deathDate,
	}

	if err := p.Validate(); err != nil {
		return nil, err
	}

	return p, nil
}

// Validate ensures the Parent entity is valid
func (p *Parent) Validate() error {
	result := validation.NewValidationResult()

	validation.ValidateID(p.id, "ID", result)
	validation.Required(p.firstName, "FirstName", result)
	validation.Required(p.lastName, "LastName", result)
	validation.MinLength(p.firstName, 2, "FirstName", result)
	validation.MinLength(p.lastName, 2, "LastName", result)
	validation.PastDate(p.birthDate, "BirthDate", result)

	if p.deathDate != nil {
		validation.PastDate(*p.deathDate, "DeathDate", result)
		validation.ValidDateRange(p.birthDate, *p.deathDate, "BirthDate", "DeathDate", result)
	}

	return result.Error()
}

// ID returns the parent's ID
func (p *Parent) ID() string {
	return p.id
}

// FirstName returns the parent's first name
func (p *Parent) FirstName() string {
	return p.firstName
}

// LastName returns the parent's last name
func (p *Parent) LastName() string {
	return p.lastName
}

// BirthDate returns the parent's birth date
func (p *Parent) BirthDate() time.Time {
	return p.birthDate
}

// DeathDate returns the parent's death date
func (p *Parent) DeathDate() *time.Time {
	if p.deathDate == nil {
		return nil
	}
	// Return a copy to prevent modification
	copy := *p.deathDate
	return &copy
}

// FullName returns the parent's full name
func (p *Parent) FullName() string {
	return p.firstName + " " + p.lastName
}

// IsDeceased returns true if the parent is deceased
func (p *Parent) IsDeceased() bool {
	return p.deathDate != nil
}

// MarkDeceased marks the parent as deceased with the given death date
func (p *Parent) MarkDeceased(deathDate time.Time) error {
	if p.deathDate != nil {
		return errors.NewDomainError(nil, "parent is already marked as deceased", "ALREADY_DECEASED")
	}

	if !deathDate.After(p.birthDate) {
		return errors.NewValidationError("death date must be after birth date")
	}

	if deathDate.After(time.Now()) {
		return errors.NewValidationError("death date cannot be in the future")
	}

	p.deathDate = &deathDate
	return nil
}

// Equals checks if two parents are the same based on ID
func (p *Parent) Equals(other *Parent) bool {
	if other == nil {
		return false
	}
	return p.id == other.id
}

// ToDTO converts the Parent entity to a data transfer object for external use
func (p *Parent) ToDTO() ParentDTO {
	dto := ParentDTO{
		ID:        p.id,
		FirstName: p.firstName,
		LastName:  p.lastName,
		BirthDate: p.birthDate,
		DeathDate: p.deathDate,
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
