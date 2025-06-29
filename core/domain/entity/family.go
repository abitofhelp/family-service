// Copyright (c) 2025 A Bit of Help, Inc.

// Package entity contains the core domain entities for the family service.
// It implements Domain-Driven Design (DDD) principles where entities represent
// the key business objects with identity and lifecycle.
package entity

import (
	"fmt"
	"time"

	domainerrors "github.com/abitofhelp/family-service/core/domain/errors"
	"github.com/abitofhelp/family-service/infrastructure/adapters/errorswrapper"
	"github.com/abitofhelp/family-service/infrastructure/adapters/identificationwrapper"
	"github.com/abitofhelp/family-service/infrastructure/adapters/validationwrapper"
	"github.com/google/uuid"
)

// Status represents the current relationship status of a family.
// In Domain-Driven Design, this is a Value Object that represents an enumeration
// of possible family states.
type Status string

// Family status constants define all possible states a family can be in.
// These statuses help enforce business rules about family composition.
const (
	// Single represents a family with only one parent
	Single Status = "SINGLE"

	// Married represents a family with two parents in a marriage relationship
	Married Status = "MARRIED"

	// Divorced represents a family where the parents were previously married but are now separated
	Divorced Status = "DIVORCED"

	// Widowed represents a family where one parent has died
	Widowed Status = "WIDOWED"

	// Abandoned represents a family where children exist without parents
	Abandoned Status = "ABANDONED"
)

// Family is the root aggregate that represents a family unit in our domain.
// 
// In Domain-Driven Design (DDD), an "aggregate" is a cluster of related objects
// that we treat as a single unit for data changes. The Family aggregate
// encapsulates Parents and Children, ensuring they maintain consistent relationships.
//
// The Family struct uses private fields to enforce that all changes must go through
// methods that can validate business rules, maintaining data integrity.
type Family struct {
	id       identificationwrapper.ID // Unique identifier for the family
	status   Status                   // Current relationship status of the family
	parents  []*Parent                // List of parents in the family (0-2)
	children []*Child                 // List of children in the family
}

// generateID creates a new unique identifier for a family.
// 
// This is an internal helper function that uses UUID v4 (random) to ensure
// each family has a globally unique identifier. We use UUIDs instead of
// sequential IDs to avoid revealing information about the number of families
// in the system and to allow for distributed systems without ID conflicts.
func generateID() string {
	// Generate a UUID v4
	return uuid.New().String()
}

// NewFamily creates a new Family aggregate with validation.
//
// This is a factory function that follows the Domain-Driven Design pattern
// for creating valid aggregates. It ensures that all business rules are
// satisfied before a Family can be created.
//
// If no ID is provided (empty string), a new UUID will be automatically generated.
// The function validates the entire aggregate, including parents and children,
// to ensure it represents a valid family according to our domain rules.
//
// Parameters:
//   - id: Unique identifier for the family (can be empty for auto-generation)
//   - status: Current relationship status of the family
//   - parents: List of parents in the family (must have at least one)
//   - children: List of children in the family (can be empty)
//
// Returns:
//   - A pointer to the new Family if valid
//   - An error if validation fails
//
// Example usage:
//
//	parent, _ := NewParent("p1", "John", "Doe", time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC), nil)
//	family, err := NewFamily("", Single, []*Parent{parent}, []*Child{})
//	if err != nil {
//	    // Handle error
//	}
func NewFamily(id string, status Status, parents []*Parent, children []*Child) (*Family, error) {
	if id == "" {
		id = generateID()
	}

	// Create ID value object
	idVO, err := identificationwrapper.NewIDFromString(id)
	if err != nil {
		return nil, errorswrapper.NewValidationError("invalid ID: "+err.Error(), "ID", err)
	}

	f := &Family{
		id:       idVO,
		status:   status,
		parents:  parents,
		children: children,
	}

	// Validate the entire aggregate to ensure it meets all business rules
	if err := f.Validate(); err != nil {
		return nil, err
	}

	return f, nil
}

// Validate ensures the Family aggregate is valid according to business rules.
//
// This method checks all invariants (business rules that must always be true)
// for a Family. It validates:
//   - The family has a valid ID
//   - The family has a valid status
//   - The family has the correct number of parents based on its status
//   - There are no duplicate parents
//   - All parents are valid
//   - All children are valid
//   - Parent-child relationships make logical sense (e.g., children born after parents)
//
// This is a crucial part of Domain-Driven Design as it ensures the entity
// always remains in a valid state.
func (f *Family) Validate() error {
	result := validationwrapper.NewValidationResult()

	// ID value object has its own validation, so we don't need to validate it here

	// Validate status
	if f.status == "" {
		result.AddError("is required", "Status")
	}

	// A family must have at least one parent
	if len(f.parents) == 0 {
		result.AddError("family must have at least one parent", "Parents")
	}

	// A family cannot have more than two parents
	if len(f.parents) > 2 {
		result.AddError("family cannot have more than two parents", "Parents")
	}

	// Check for duplicate parents
	seen := make(map[string]bool)
	for i, p := range f.parents {
		if p == nil {
			result.AddError(fmt.Sprintf("parent at index %d is nil", i), "Parents")
			continue
		}

		// Validate each parent
		if err := p.Validate(); err != nil {
			result.AddError(fmt.Sprintf("parent at index %d is invalid: %v", i, err), "Parents")
		}

		// Check for duplicates using a composite key
		key := p.FirstName() + p.LastName() + p.BirthDate().Format(time.RFC3339)
		if seen[key] {
			result.AddError("duplicate parent in family", "Parents")
		}
		seen[key] = true

		// Enhanced validation: Validate parent age (minimum 18 years)
		now := time.Now()
		birthDate := p.BirthDate()
		age := now.Year() - birthDate.Year()

		// Adjust age if birthday hasn't occurred yet this year
		if now.Month() < birthDate.Month() || (now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
			age--
		}

		if age < 18 {
			result.AddError(fmt.Sprintf("parent at index %d does not meet minimum age requirement (18 years)", i), "Parents")
		}
	}

	// Validate children
	for i, c := range f.children {
		if c == nil {
			result.AddError(fmt.Sprintf("child at index %d is nil", i), "Children")
			continue
		}

		// Validate each child
		if err := c.Validate(); err != nil {
			result.AddError(fmt.Sprintf("child at index %d is invalid: %v", i, err), "Children")
		}

		// Enhanced validation: Validate child's birth date is after parents' birth dates
		childBirthDate := c.BirthDate()
		for j, p := range f.parents {
			parentBirthDate := p.BirthDate()
			if !childBirthDate.After(parentBirthDate) {
				result.AddError(fmt.Sprintf("child at index %d has birth date before parent at index %d", i, j), "Children")
			}
		}

		// Enhanced validation: Validate child's birth date is not in the future
		if childBirthDate.After(time.Now()) {
			result.AddError(fmt.Sprintf("child at index %d has birth date in the future", i), "Children")
		}
	}

	// Validate status transitions based on parents
	if f.status == Married && len(f.parents) != 2 {
		result.AddError("married family must have exactly two parents", "Status")
	}

	if f.status == Single && len(f.parents) > 1 {
		result.AddError("single family cannot have more than one parent", "Status")
	}

	// Enhanced validation: Validate Divorced status
	if f.status == Divorced && len(f.parents) != 1 {
		result.AddError("divorced family must have exactly one parent", "Status")
	}

	// Enhanced validation: Validate Widowed status
	if f.status == Widowed {
		if len(f.parents) != 1 {
			result.AddError("widowed family must have exactly one parent", "Status")
		} else {
			// Check if the remaining parent is deceased
			if f.parents[0].IsDeceased() {
				result.AddError("widowed family cannot have a deceased parent", "Status")
			}
		}
	}

	// Enhanced validation: Validate Abandoned status
	if f.status == Abandoned && len(f.children) == 0 {
		result.AddError("abandoned family must have at least one child", "Status")
	}

	// Enhanced validation: Validate parent-child age gap
	for i, child := range f.children {
		childBirthDate := child.BirthDate()

		for j, parent := range f.parents {
			parentBirthDate := parent.BirthDate()

			// Calculate the age difference in years
			ageGap := childBirthDate.Year() - parentBirthDate.Year()

			// Adjust for partial years
			if childBirthDate.Month() < parentBirthDate.Month() || 
			   (childBirthDate.Month() == parentBirthDate.Month() && childBirthDate.Day() < parentBirthDate.Day()) {
				ageGap--
			}

			// Minimum 12 years between parent and child
			if ageGap < 12 {
				result.AddError(fmt.Sprintf("child at index %d has too small age gap with parent at index %d", i, j), "Children")
			}
		}
	}

	return result.Error()
}

// ID returns the family's unique identifier.
//
// This is a getter method that provides read-only access to the private id field.
// In Domain-Driven Design, we use getters to control access to entity properties
// while keeping the internal state encapsulated.
func (f *Family) ID() string {
	return string(f.id)
}

// Status returns the family's current relationship status.
//
// This getter provides read-only access to the family's status (Single, Married, etc.).
// The status helps determine what operations are valid for this family and what
// business rules apply.
func (f *Family) Status() Status {
	return f.status
}

// Parents returns a copy of the family's parents.
//
// This method returns a defensive copy of the parents slice rather than the
// original slice. This is an important pattern in Domain-Driven Design that
// prevents external code from modifying the internal state of the aggregate
// without going through the proper methods that enforce business rules.
//
// If you need to modify the parents, use methods like AddParent or RemoveParent
// which ensure all business rules are maintained.
func (f *Family) Parents() []*Parent {
	result := make([]*Parent, len(f.parents))
	copy(result, f.parents)
	return result
}

// Children returns a copy of the family's children.
//
// Like the Parents method, this returns a defensive copy to protect the
// internal state of the aggregate. This ensures that all changes to the
// children collection go through methods that validate business rules.
//
// If you need to modify the children, use methods like AddChild or RemoveChild
// which ensure all business rules are maintained.
func (f *Family) Children() []*Child {
	result := make([]*Child, len(f.children))
	copy(result, f.children)
	return result
}

// CountParents returns the number of parents in the family.
//
// This is a convenience method that returns the current count of parents.
// The family business rules specify that a family can have 0-2 parents
// depending on its status.
func (f *Family) CountParents() int {
	return len(f.parents)
}

// CountChildren returns the number of children in the family.
//
// This is a convenience method that returns the current count of children.
// A family can have any number of children, including zero.
func (f *Family) CountChildren() int {
	return len(f.children)
}

// AddParent adds a parent to the family while enforcing business rules.
//
// This method maintains the integrity of the Family aggregate by:
// 1. Ensuring the parent is not nil
// 2. Checking that the family doesn't exceed the maximum of two parents
// 3. Preventing duplicate parents (either by ID or by name+birthdate)
//
// Note that adding a parent doesn't automatically update the family status.
// This gives the caller control over the family's status, which should be
// set appropriately based on the business context (e.g., marriage, adoption).
//
// Returns:
//   - nil if the parent was successfully added
//   - FamilyTooManyParentsError if the family already has two parents
//   - FamilyParentExistsError if the exact parent already exists in the family
//   - FamilyParentDuplicateError if a parent with the same identity exists
//   - ValidationError if the parent is nil
func (f *Family) AddParent(p *Parent) error {
	if p == nil {
		return errorswrapper.NewValidationError("parent cannot be nil", "Parent", nil)
	}

	if len(f.parents) >= 2 {
		return domainerrors.NewFamilyTooManyParentsError("family cannot have more than two parents", nil)
	}

	// Check for duplicate parent
	for _, existingParent := range f.parents {
		if existingParent.Equals(p) {
			return domainerrors.NewFamilyParentExistsError("parent already exists in family", nil)
		}

		// Check for duplicate based on name and birthdate
		if existingParent.FirstName() == p.FirstName() &&
			existingParent.LastName() == p.LastName() &&
			existingParent.BirthDate().Equal(p.BirthDate()) {
			return domainerrors.NewFamilyParentDuplicateError("parent with same name and birthdate already exists in family", nil)
		}
	}

	f.parents = append(f.parents, p)

	// We don't automatically update the status when adding a parent
	// This allows the caller to control the family status

	return nil
}

// AddChild adds a child to the family while enforcing business rules.
//
// This method maintains the integrity of the Family aggregate by:
// 1. Ensuring the child is not nil
// 2. Preventing duplicate children
//
// In our domain model, a child can only belong to one family at a time.
// If a child needs to move to a different family (e.g., in adoption scenarios),
// it should first be removed from its current family.
//
// Returns:
//   - nil if the child was successfully added
//   - FamilyChildExistsError if the child already exists in the family
//   - ValidationError if the child is nil
func (f *Family) AddChild(c *Child) error {
	if c == nil {
		return errorswrapper.NewValidationError("child cannot be nil", "Child", nil)
	}

	// Check for duplicate child
	for _, existingChild := range f.children {
		if existingChild.Equals(c) {
			return domainerrors.NewFamilyChildExistsError("child already exists in family", nil)
		}
	}

	f.children = append(f.children, c)
	return nil
}

// RemoveChild removes a child from the family by ID.
//
// This method allows a child to be removed from a family, which might happen
// in scenarios like:
// - The child is being adopted by another family
// - The child has reached adulthood and is forming their own family
// - Correcting data entry errors
//
// Returns:
//   - nil if the child was successfully removed
//   - NotFoundError if no child with the given ID exists in the family
func (f *Family) RemoveChild(childID string) error {
	for i, c := range f.children {
		if c.ID() == childID {
			// Remove child at index i
			f.children = append(f.children[:i], f.children[i+1:]...)
			return nil
		}
	}
	return errorswrapper.NewNotFoundError("Child", childID, nil)
}

// RemoveParent removes a parent from the family by ID.
//
// This method maintains the integrity of the Family aggregate by:
// 1. Preventing removal of the last parent (a family must have at least one parent)
// 2. Automatically updating the family status if needed (e.g., from Married to Single)
//
// This might be used in scenarios like:
// - Divorce (though the Divorce method is preferred for this specific case)
// - Death of a parent (though MarkParentDeceased is preferred)
// - Correcting data entry errors
//
// Returns:
//   - nil if the parent was successfully removed
//   - FamilyCannotRemoveLastParentError if attempting to remove the only parent
//   - NotFoundError if no parent with the given ID exists in the family
func (f *Family) RemoveParent(parentID string) error {
	if len(f.parents) <= 1 {
		return domainerrors.NewFamilyCannotRemoveLastParentError("cannot remove the only parent from a family", nil)
	}

	for i, p := range f.parents {
		if p.ID() == parentID {
			// Remove parent at index i
			f.parents = append(f.parents[:i], f.parents[i+1:]...)

			// Update status if needed
			if len(f.parents) == 1 && f.status == Married {
				f.status = Single
			}

			return nil
		}
	}
	return errorswrapper.NewNotFoundError("Parent", parentID, nil)
}

// MarkParentDeceased marks a parent as deceased and updates family status if needed.
//
// This method handles the important life event of a parent's death by:
// 1. Finding the parent by ID
// 2. Marking that parent as deceased with the provided death date
// 3. Updating the family status if appropriate (e.g., from Married to Widowed)
//
// This is an example of how domain logic encapsulates real-world events and
// ensures that all related state changes happen consistently.
//
// Parameters:
//   - parentID: The ID of the parent to mark as deceased
//   - deathDate: The date when the parent died
//
// Returns:
//   - nil if the parent was successfully marked as deceased
//   - NotFoundError if no parent with the given ID exists in the family
//   - ValidationError if the death date is invalid (from the Parent.MarkDeceased method)
func (f *Family) MarkParentDeceased(parentID string, deathDate time.Time) error {
	var foundParent *Parent

	for _, p := range f.parents {
		if p.ID() == parentID {
			foundParent = p
			break
		}
	}

	if foundParent == nil {
		return errorswrapper.NewNotFoundError("Parent", parentID, nil)
	}

	if err := foundParent.MarkDeceased(deathDate); err != nil {
		return err
	}

	// If this was a married family with two parents, and one died, update status
	if f.status == Married && len(f.parents) == 2 {
		f.status = Widowed
	}

	return nil
}

// Divorce handles the divorce process, creating a new family for the non-custodial parent.
//
// This method implements a complex domain operation that models a real-world event.
// When parents divorce:
// 1. The original family keeps the custodial parent and all children
// 2. A new family is created for the non-custodial parent
// 3. Both families' statuses are updated to Divorced
//
// This approach maintains the integrity of family relationships while accurately
// representing the real-world situation after a divorce.
//
// Parameters:
//   - custodialParentID: The ID of the parent who will keep custody of the children
//
// Returns:
//   - A pointer to the new Family created for the non-custodial parent
//   - FamilyNotMarriedError if the family's status isn't Married
//   - FamilyDivorceRequiresTwoParentsError if the family doesn't have exactly two parents
//   - NotFoundError if the custodial parent ID doesn't match any parent in the family
//   - FamilyCreateFailedError if creating the new family fails
//
// Example usage:
//
//	// Assuming we have a married family with two parents and children
//	newFamily, err := family.Divorce("parent-123") // parent-123 is the custodial parent ID
//	if err != nil {
//	    // Handle error
//	}
//	// Now we have two families:
//	// 1. The original family with parent-123 and all children, status = Divorced
//	// 2. newFamily with the other parent, no children, status = Divorced
func (f *Family) Divorce(custodialParentID string) (*Family, error) {
	if f.status != Married {
		return nil, domainerrors.NewFamilyNotMarriedError("only married families can divorce", nil)
	}

	if len(f.parents) != 2 {
		return nil, domainerrors.NewFamilyDivorceRequiresTwoParentsError("divorce requires exactly two parents", nil)
	}

	var custodialParent *Parent
	var remainingParent *Parent

	for _, p := range f.parents {
		if p.ID() == custodialParentID {
			custodialParent = p
		} else {
			remainingParent = p
		}
	}

	if custodialParent == nil {
		return nil, errorswrapper.NewNotFoundError("Parent", custodialParentID, nil)
	}

	// Create a new family for the remaining parent (this will get a new ID)
	// The original family ID will stay with the custodial parent and children
	remainingFamily, err := NewFamily(
		"", // Empty ID will cause a new ID to be generated
		Divorced,
		[]*Parent{remainingParent},
		[]*Child{}, // No children with the remaining parent
	)

	if err != nil {
		return nil, domainerrors.NewFamilyCreateFailedError("failed to create new family for remaining parent", err)
	}

	// Update the original family to keep only the custodial parent
	f.parents = []*Parent{custodialParent}
	f.status = Divorced

	// Return the new family with the remaining parent
	// The original family (with custodial parent and children) is modified in place
	return remainingFamily, nil
}

// ToDTO converts the Family aggregate to a data transfer object for external use.
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
//   - A FamilyDTO containing all the relevant data from this Family
func (f *Family) ToDTO() FamilyDTO {
	parentDTOs := make([]ParentDTO, len(f.parents))
	for i, p := range f.parents {
		parentDTOs[i] = p.ToDTO()
	}

	childDTOs := make([]ChildDTO, len(f.children))
	for i, c := range f.children {
		childDTOs[i] = c.ToDTO()
	}

	return FamilyDTO{
		ID:            string(f.id),
		Status:        string(f.status),
		Parents:       parentDTOs,
		Children:      childDTOs,
		ParentCount:   f.CountParents(),
		ChildrenCount: f.CountChildren(),
	}
}

// FamilyDTO is a data transfer object for the Family aggregate.
//
// This struct represents the Family entity in a format suitable for
// transferring data between layers of the application. It contains
// only simple data types and structures that can be easily serialized.
//
// In Domain-Driven Design, DTOs help maintain the separation between
// the domain model and external interfaces, preventing domain logic
// from leaking into other layers.
type FamilyDTO struct {
	ID            string      // Unique identifier for the family
	Status        string      // Current relationship status as a string
	Parents       []ParentDTO // List of parent DTOs
	Children      []ChildDTO  // List of child DTOs
	ParentCount   int         // Number of parents in the family
	ChildrenCount int         // Number of children in the family
}

// FamilyFromDTO creates a Family aggregate from a data transfer object.
//
// This function is the counterpart to ToDTO and is used to reconstruct
// a valid Family domain entity from a DTO. It:
// 1. Converts all parent DTOs back to Parent entities
// 2. Converts all child DTOs back to Child entities
// 3. Creates a new Family with the converted entities
// 4. Validates the resulting Family to ensure it meets all domain rules
//
// This is typically used when receiving data from external sources
// (like API requests) and needing to work with it in the domain layer.
//
// Parameters:
//   - dto: The FamilyDTO containing the data to convert
//
// Returns:
//   - A pointer to the new Family if valid
//   - An error if validation fails or any conversion fails
func FamilyFromDTO(dto FamilyDTO) (*Family, error) {
	parents := make([]*Parent, len(dto.Parents))
	for i, pDTO := range dto.Parents {
		p, err := ParentFromDTO(pDTO)
		if err != nil {
			return nil, err
		}
		parents[i] = p
	}

	children := make([]*Child, len(dto.Children))
	for i, cDTO := range dto.Children {
		c, err := ChildFromDTO(cDTO)
		if err != nil {
			return nil, err
		}
		children[i] = c
	}

	return NewFamily(dto.ID, Status(dto.Status), parents, children)
}
