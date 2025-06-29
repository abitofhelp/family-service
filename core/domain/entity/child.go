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

// Child represents a child entity in the family domain.
//
// In Domain-Driven Design (DDD), an "entity" is an object with a distinct identity
// that runs through time and different states. The Child entity represents a person
// who is a child in a family.
//
// Unlike the Parent entity, a Child can only belong to one family at a time in our
// domain model. This is an important business rule that reflects how we model
// family relationships.
//
// The Child struct uses Value Objects (like ID, Name, DateOfBirth) for its properties
// to encapsulate validation and business rules. This is a common pattern in DDD that
// helps maintain the integrity of the domain model.
//
// All fields are private to enforce that changes must go through methods that can
// validate business rules, maintaining data integrity.
type Child struct {
	id        identificationwrapper.ID            // Unique identifier for the child
	firstName identificationwrapper.Name          // First name of the child
	lastName  identificationwrapper.Name          // Last name of the child
	birthDate identificationwrapper.DateOfBirth   // Birth date of the child
	deathDate *identificationwrapper.DateOfDeath  // Death date of the child (nil if alive)
}

// NewChild creates a new Child entity with validation.
//
// This is a factory function that follows the Domain-Driven Design pattern
// for creating valid entities. It ensures that all business rules are
// satisfied before a Child can be created.
//
// The function converts primitive types (strings, time.Time) into Value Objects
// that encapsulate validation and business rules. This is a key aspect of DDD
// that helps maintain the integrity of the domain model.
//
// Parameters:
//   - id: Unique identifier for the child
//   - firstName: First name of the child
//   - lastName: Last name of the child
//   - birthDate: Birth date of the child
//   - deathDate: Death date of the child (nil if alive)
//
// Returns:
//   - A pointer to the new Child if valid
//   - An error if validation fails
//
// Example usage:
//
//	// Create a living child
//	child, err := NewChild("c1", "Jane", "Doe", time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC), nil)
//	if err != nil {
//	    // Handle error
//	}
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

// Validate ensures the Child entity is valid according to business rules.
//
// This method checks all invariants (business rules that must always be true)
// for a Child. It validates:
//   - Names meet minimum length requirements
//   - Names contain only valid characters
//   - Birth date is not in the future
//   - Death date (if present) is after birth date and not in the future
//   - Child's age is within reasonable limits
//
// This is a crucial part of Domain-Driven Design as it ensures the entity
// always remains in a valid state.
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

// ID returns the child's unique identifier.
//
// This is a getter method that provides read-only access to the private id field.
// In Domain-Driven Design, we use getters to control access to entity properties
// while keeping the internal state encapsulated.
func (c *Child) ID() string {
	return c.id.String()
}

// FirstName returns the child's first name.
//
// This getter provides read-only access to the firstName value object.
// It returns the string representation of the Name value object.
func (c *Child) FirstName() string {
	return c.firstName.String()
}

// LastName returns the child's last name.
//
// This getter provides read-only access to the lastName value object.
// It returns the string representation of the Name value object.
func (c *Child) LastName() string {
	return c.lastName.String()
}

// BirthDate returns the child's birth date.
//
// This getter provides read-only access to the birthDate value object.
// It returns the time.Time representation of the DateOfBirth value object.
func (c *Child) BirthDate() time.Time {
	return c.birthDate.Date()
}

// DeathDate returns the child's death date, or nil if the child is alive.
//
// This getter provides read-only access to the deathDate value object.
// It returns a copy of the time.Time value to prevent modification of the
// internal state, which is an important pattern in Domain-Driven Design.
//
// Returns:
//   - A pointer to a copy of the death date if the child is deceased
//   - nil if the child is alive
func (c *Child) DeathDate() *time.Time {
	if c.deathDate == nil {
		return nil
	}
	// Return a copy to prevent modification
	date := c.deathDate.Date()
	return &date
}

// FullName returns the child's full name (first name + last name).
//
// This is a derived property that combines the first and last names.
// In Domain-Driven Design, derived properties are calculated from other
// properties and don't need to be stored separately.
func (c *Child) FullName() string {
	return c.firstName.String() + " " + c.lastName.String()
}

// IsDeceased returns true if the child is deceased (has a death date).
//
// This is a convenience method that checks if the deathDate field is set.
// It's an example of how domain entities can provide methods that express
// domain concepts directly, making the code more readable and expressive.
func (c *Child) IsDeceased() bool {
	return c.deathDate != nil
}

// MarkDeceased marks the child as deceased with the given death date.
//
// This method handles the important life event of a child's death by:
// 1. Checking if the child is already marked as deceased
// 2. Validating that the death date is after the birth date
// 3. Creating a DateOfDeath value object
// 4. Updating the child's state
//
// This is an example of how domain logic encapsulates real-world events and
// ensures that all related state changes happen consistently.
//
// Parameters:
//   - deathDate: The date when the child died
//
// Returns:
//   - nil if the child was successfully marked as deceased
//   - ChildAlreadyDeceasedError if the child is already marked as deceased
//   - ValidationError if the death date is invalid or before the birth date
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

// Equals checks if two children are the same based on ID.
//
// In Domain-Driven Design, entities are distinguished by their identity, not their
// attributes. This method implements the equality comparison based on the ID,
// which is the identity of the Child entity.
//
// This is important for operations like checking if a child already exists in a
// collection, or finding a specific child in a list.
//
// Parameters:
//   - other: Another Child entity to compare with
//
// Returns:
//   - true if the children have the same ID
//   - false if the other child is nil or has a different ID
func (c *Child) Equals(other *Child) bool {
	if other == nil {
		return false
	}
	return c.id.Equals(other.id)
}

// ToDTO converts the Child entity to a data transfer object for external use.
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
//   - A ChildDTO containing all the relevant data from this Child
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

// ChildDTO is a data transfer object for the Child entity.
//
// This struct represents the Child entity in a format suitable for
// transferring data between layers of the application. It contains
// only simple data types and structures that can be easily serialized.
//
// In Domain-Driven Design, DTOs help maintain the separation between
// the domain model and external interfaces, preventing domain logic
// from leaking into other layers.
type ChildDTO struct {
	ID        string     // Unique identifier for the child
	FirstName string     // First name of the child
	LastName  string     // Last name of the child
	BirthDate time.Time  // Birth date of the child
	DeathDate *time.Time // Death date of the child (nil if alive)
}

// ChildFromDTO creates a Child entity from a data transfer object.
//
// This function is the counterpart to ToDTO and is used to reconstruct
// a valid Child domain entity from a DTO. It delegates to the NewChild
// factory function to ensure that all business rules are satisfied.
//
// This is typically used when receiving data from external sources
// (like API requests) and needing to work with it in the domain layer.
//
// Parameters:
//   - dto: The ChildDTO containing the data to convert
//
// Returns:
//   - A pointer to the new Child if valid
//   - An error if validation fails
func ChildFromDTO(dto ChildDTO) (*Child, error) {
	return NewChild(dto.ID, dto.FirstName, dto.LastName, dto.BirthDate, dto.DeathDate)
}
