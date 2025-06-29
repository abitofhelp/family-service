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

// Parent represents a parent entity in the family domain.
//
// In Domain-Driven Design (DDD), an "entity" is an object with a distinct identity
// that runs through time and different states. The Parent entity represents a person
// who is a parent in one or more families.
//
// The Parent struct uses Value Objects (like ID, Name, DateOfBirth) for its properties
// to encapsulate validation and business rules. This is a common pattern in DDD that
// helps maintain the integrity of the domain model.
//
// All fields are private to enforce that changes must go through methods that can
// validate business rules, maintaining data integrity.
type Parent struct {
	id        identificationwrapper.ID            // Unique identifier for the parent
	firstName identificationwrapper.Name          // First name of the parent
	lastName  identificationwrapper.Name          // Last name of the parent
	birthDate identificationwrapper.DateOfBirth   // Birth date of the parent
	deathDate *identificationwrapper.DateOfDeath  // Death date of the parent (nil if alive)
}

// NewParent creates a new Parent entity with validation.
//
// This is a factory function that follows the Domain-Driven Design pattern
// for creating valid entities. It ensures that all business rules are
// satisfied before a Parent can be created.
//
// The function converts primitive types (strings, time.Time) into Value Objects
// that encapsulate validation and business rules. This is a key aspect of DDD
// that helps maintain the integrity of the domain model.
//
// Parameters:
//   - id: Unique identifier for the parent
//   - firstName: First name of the parent
//   - lastName: Last name of the parent
//   - birthDate: Birth date of the parent
//   - deathDate: Death date of the parent (nil if alive)
//
// Returns:
//   - A pointer to the new Parent if valid
//   - An error if validation fails
//
// Example usage:
//
//	// Create a living parent
//	parent, err := NewParent("p1", "John", "Doe", time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC), nil)
//	if err != nil {
//	    // Handle error
//	}
//
//	// Create a deceased parent
//	deathDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
//	parent, err := NewParent("p2", "Jane", "Doe", time.Date(1985, 1, 1, 0, 0, 0, 0, time.UTC), &deathDate)
func NewParent(id, firstName, lastName string, birthDate time.Time, deathDate *time.Time) (*Parent, error) {
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

	p := &Parent{
		id:        idVO,
		firstName: firstNameVO,
		lastName:  lastNameVO,
		birthDate: birthDateVO,
		deathDate: deathDateVO,
	}

	// Validate the entire entity to ensure it meets all business rules
	if err := p.Validate(); err != nil {
		return nil, err
	}

	return p, nil
}

// Validate ensures the Parent entity is valid according to business rules.
//
// This method checks all invariants (business rules that must always be true)
// for a Parent. It validates:
//   - Names meet minimum length requirements
//   - Names contain only valid characters
//   - Birth date is not in the future
//   - Death date (if present) is after birth date and not in the future
//   - Parent meets minimum age requirement (18 years)
//   - Parent's age is within reasonable limits
//
// This is a crucial part of Domain-Driven Design as it ensures the entity
// always remains in a valid state.
func (p *Parent) Validate() error {
	result := validationwrapper.NewValidationResult()

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

	// Enhanced validation: Validate name doesn't contain special characters
	nameRegex := "^[a-zA-Z\\s-]+$"
	if matched, _ := regexp.MatchString(nameRegex, p.firstName.String()); !matched {
		result.AddError("must contain only letters, spaces, and hyphens", "FirstName")
	}

	if matched, _ := regexp.MatchString(nameRegex, p.lastName.String()); !matched {
		result.AddError("must contain only letters, spaces, and hyphens", "LastName")
	}

	// Enhanced validation: Validate birth date is not in the future
	if p.birthDate.Date().After(time.Now()) {
		result.AddError("birth date cannot be in the future", "BirthDate")
	}

	// Enhanced validation: Validate death date is not in the future
	if p.deathDate != nil && p.deathDate.Date().After(time.Now()) {
		result.AddError("death date cannot be in the future", "DeathDate")
	}

	// Enhanced validation: Validate minimum age for a parent (18 years)
	minAge := 18
	now := time.Now()
	birthDate := p.birthDate.Date()
	age := now.Year() - birthDate.Year()

	// Adjust age if birthday hasn't occurred yet this year
	if now.Month() < birthDate.Month() || (now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		age--
	}

	if age < minAge {
		result.AddError(fmt.Sprintf("parent does not meet minimum age requirement (%d years)", minAge), "BirthDate")
	}

	// Enhanced validation: Validate maximum age (e.g., 150 years)
	maxAge := 150
	if age > maxAge {
		result.AddError(fmt.Sprintf("age cannot exceed %d years", maxAge), "BirthDate")
	}

	return result.Error()
}

// ID returns the parent's unique identifier.
//
// This is a getter method that provides read-only access to the private id field.
// In Domain-Driven Design, we use getters to control access to entity properties
// while keeping the internal state encapsulated.
func (p *Parent) ID() string {
	return string(p.id)
}

// FirstName returns the parent's first name.
//
// This getter provides read-only access to the firstName value object.
// It returns the string representation of the Name value object.
func (p *Parent) FirstName() string {
	return p.firstName.String()
}

// LastName returns the parent's last name.
//
// This getter provides read-only access to the lastName value object.
// It returns the string representation of the Name value object.
func (p *Parent) LastName() string {
	return p.lastName.String()
}

// BirthDate returns the parent's birth date.
//
// This getter provides read-only access to the birthDate value object.
// It returns the time.Time representation of the DateOfBirth value object.
func (p *Parent) BirthDate() time.Time {
	return p.birthDate.Date()
}

// DeathDate returns the parent's death date, or nil if the parent is alive.
//
// This getter provides read-only access to the deathDate value object.
// It returns a copy of the time.Time value to prevent modification of the
// internal state, which is an important pattern in Domain-Driven Design.
//
// Returns:
//   - A pointer to a copy of the death date if the parent is deceased
//   - nil if the parent is alive
func (p *Parent) DeathDate() *time.Time {
	if p.deathDate == nil {
		return nil
	}
	// Return a copy to prevent modification
	date := p.deathDate.Date()
	return &date
}

// FullName returns the parent's full name (first name + last name).
//
// This is a derived property that combines the first and last names.
// In Domain-Driven Design, derived properties are calculated from other
// properties and don't need to be stored separately.
func (p *Parent) FullName() string {
	return p.firstName.String() + " " + p.lastName.String()
}

// IsDeceased returns true if the parent is deceased (has a death date).
//
// This is a convenience method that checks if the deathDate field is set.
// It's an example of how domain entities can provide methods that express
// domain concepts directly, making the code more readable and expressive.
func (p *Parent) IsDeceased() bool {
	return p.deathDate != nil
}

// MarkDeceased marks the parent as deceased with the given death date.
//
// This method handles the important life event of a parent's death by:
// 1. Checking if the parent is already marked as deceased
// 2. Validating that the death date is after the birth date
// 3. Creating a DateOfDeath value object
// 4. Updating the parent's state
//
// This is an example of how domain logic encapsulates real-world events and
// ensures that all related state changes happen consistently.
//
// Parameters:
//   - deathDate: The date when the parent died
//
// Returns:
//   - nil if the parent was successfully marked as deceased
//   - ParentAlreadyDeceasedError if the parent is already marked as deceased
//   - ValidationError if the death date is invalid or before the birth date
func (p *Parent) MarkDeceased(deathDate time.Time) error {
	if p.deathDate != nil {
		return domainerrors.NewParentAlreadyDeceasedError("parent is already marked as deceased", nil)
	}

	// Validate death date is after birth date
	if !deathDate.After(p.birthDate.Date()) {
		return errorswrapper.NewValidationError("death date must be after birth date", "DeathDate", nil)
	}

	// Create DateOfDeath value object
	year, month, day := deathDate.Date()
	dod, err := identificationwrapper.NewDateOfDeath(year, int(month), day)
	if err != nil {
		return errorswrapper.NewValidationError("invalid death date: "+err.Error(), "DeathDate", err)
	}

	p.deathDate = &dod
	return nil
}

// Equals checks if two parents are the same based on ID.
//
// In Domain-Driven Design, entities are distinguished by their identity, not their
// attributes. This method implements the equality comparison based on the ID,
// which is the identity of the Parent entity.
//
// This is important for operations like checking if a parent already exists in a
// collection, or finding a specific parent in a list.
//
// Parameters:
//   - other: Another Parent entity to compare with
//
// Returns:
//   - true if the parents have the same ID
//   - false if the other parent is nil or has a different ID
func (p *Parent) Equals(other *Parent) bool {
	if other == nil {
		return false
	}
	return p.id.Equals(other.id)
}

// ToDTO converts the Parent entity to a data transfer object for external use.
//
// This method creates a DTO (Data Transfer Object) that can be safely passed
// across layer boundaries in the application. DTOs are important in Domain-Driven
// Design as they:
// 1. Decouple the domain model from external representations
// 2. Allow for different representations of the same domain object
// 3. Prevent accidental modifications to the domain model
//
// The DTO contains only simple data types and structures, making it suitable
// for serialization (e.g., to JSON for API responses).
//
// Returns:
//   - A ParentDTO containing all the relevant data from this Parent
func (p *Parent) ToDTO() ParentDTO {
	var deathDate *time.Time
	if p.deathDate != nil {
		date := p.deathDate.Date()
		deathDate = &date
	}

	dto := ParentDTO{
		ID:        string(p.id),
		FirstName: p.firstName.String(),
		LastName:  p.lastName.String(),
		BirthDate: p.birthDate.Date(),
		DeathDate: deathDate,
	}
	return dto
}

// ParentDTO is a data transfer object for the Parent entity.
//
// This struct represents the Parent entity in a format suitable for
// transferring data between layers of the application. It contains
// only simple data types and structures that can be easily serialized.
//
// In Domain-Driven Design, DTOs help maintain the separation between
// the domain model and external interfaces, preventing domain logic
// from leaking into other layers.
type ParentDTO struct {
	ID        string     // Unique identifier for the parent
	FirstName string     // First name of the parent
	LastName  string     // Last name of the parent
	BirthDate time.Time  // Birth date of the parent
	DeathDate *time.Time // Death date of the parent (nil if alive)
}

// ParentFromDTO creates a Parent entity from a data transfer object.
//
// This function is the counterpart to ToDTO and is used to reconstruct
// a valid Parent domain entity from a DTO. It delegates to the NewParent
// factory function to ensure that all business rules are satisfied.
//
// This is typically used when receiving data from external sources
// (like API requests) and needing to work with it in the domain layer.
//
// Parameters:
//   - dto: The ParentDTO containing the data to convert
//
// Returns:
//   - A pointer to the new Parent if valid
//   - An error if validation fails
func ParentFromDTO(dto ParentDTO) (*Parent, error) {
	return NewParent(dto.ID, dto.FirstName, dto.LastName, dto.BirthDate, dto.DeathDate)
}
