// Copyright (c) 2025 A Bit of Help, Inc.

package entity

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/abitofhelp/servicelib/errors"
	"github.com/abitofhelp/servicelib/validation"
)

// Status represents the current status of a family
type Status string

// Family status constants
const (
	Single    Status = "SINGLE"
	Married   Status = "MARRIED"
	Divorced  Status = "DIVORCED"
	Widowed   Status = "WIDOWED"
	Abandoned Status = "ABANDONED"
)

// Family is the root aggregate that represents a family unit
type Family struct {
	id       string //FIXME: We'll keep this as a string for now since it's using a custom generation function
	status   Status
	parents  []*Parent
	children []*Child
}

// generateID generates a simple unique ID for a family
func generateID() string {
	// Simple ID generation using timestamp and random number
	return "fam-" + strconv.FormatInt(time.Now().UnixNano(), 10) + "-" + strconv.Itoa(rand.Intn(1000))
}

// NewFamily creates a new Family aggregate with validation
func NewFamily(id string, status Status, parents []*Parent, children []*Child) (*Family, error) {
	if id == "" {
		id = generateID()
	}

	f := &Family{
		id:       id,
		status:   status,
		parents:  parents,
		children: children,
	}

	if err := f.Validate(); err != nil {
		return nil, err
	}

	return f, nil
}

// Validate ensures the Family aggregate is valid
func (f *Family) Validate() error {
	result := validation.NewValidationResult()

	validation.ValidateID(f.id, "ID", result)

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
		key := p.FirstName() + p.LastName() + p.BirthDate().Format("2006-01-02")
		if seen[key] {
			result.AddError("duplicate parent in family", "Parents")
		}
		seen[key] = true
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
	}

	// Validate status transitions based on parents
	if f.status == Married && len(f.parents) != 2 {
		result.AddError("married family must have exactly two parents", "Status")
	}

	if f.status == Single && len(f.parents) > 1 {
		result.AddError("single family cannot have more than one parent", "Status")
	}

	return result.Error()
}

// ID returns the family's ID
func (f *Family) ID() string {
	return f.id
}

// Status returns the family's status
func (f *Family) Status() Status {
	return f.status
}

// Parents returns a copy of the family's parents
func (f *Family) Parents() []*Parent {
	result := make([]*Parent, len(f.parents))
	copy(result, f.parents)
	return result
}

// Children returns a copy of the family's children
func (f *Family) Children() []*Child {
	result := make([]*Child, len(f.children))
	copy(result, f.children)
	return result
}

// CountParents returns the number of parents in the family
func (f *Family) CountParents() int {
	return len(f.parents)
}

// CountChildren returns the number of children in the family
func (f *Family) CountChildren() int {
	return len(f.children)
}

// AddParent adds a parent to the family
func (f *Family) AddParent(p *Parent) error {
	if p == nil {
		return errors.NewValidationError("parent cannot be nil", "Parent", nil)
	}

	if len(f.parents) >= 2 {
		return errors.NewDomainError(errors.BusinessRuleViolationCode, "family cannot have more than two parents", nil)
	}

	// Check for duplicate parent
	for _, existingParent := range f.parents {
		if existingParent.Equals(p) {
			return errors.NewDomainError(errors.BusinessRuleViolationCode, "parent already exists in family", nil)
		}

		// Check for duplicate based on name and birthdate
		if existingParent.FirstName() == p.FirstName() &&
			existingParent.LastName() == p.LastName() &&
			existingParent.BirthDate().Equal(p.BirthDate()) {
			return errors.NewDomainError(errors.BusinessRuleViolationCode, "parent with same name and birthdate already exists in family", nil)
		}
	}

	f.parents = append(f.parents, p)

	// We don't automatically update the status when adding a parent
	// This allows the caller to control the family status

	return nil
}

// AddChild adds a child to the family
func (f *Family) AddChild(c *Child) error {
	if c == nil {
		return errors.NewValidationError("child cannot be nil", "Child", nil)
	}

	// Check for duplicate child
	for _, existingChild := range f.children {
		if existingChild.Equals(c) {
			return errors.NewDomainError(errors.BusinessRuleViolationCode, "child already exists in family", nil)
		}
	}

	f.children = append(f.children, c)
	return nil
}

// RemoveChild removes a child from the family
func (f *Family) RemoveChild(childID string) error {
	for i, c := range f.children {
		if c.ID() == childID {
			// Remove child at index i
			f.children = append(f.children[:i], f.children[i+1:]...)
			return nil
		}
	}
	return errors.NewNotFoundError("Child", childID, nil)
}

// RemoveParent removes a parent from the family
func (f *Family) RemoveParent(parentID string) error {
	if len(f.parents) <= 1 {
		return errors.NewDomainError(errors.BusinessRuleViolationCode, "cannot remove the only parent from a family", nil)
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
	return errors.NewNotFoundError("Parent", parentID, nil)
}

// MarkParentDeceased marks a parent as deceased and updates family status if needed
func (f *Family) MarkParentDeceased(parentID string, deathDate time.Time) error {
	var foundParent *Parent

	for _, p := range f.parents {
		if p.ID() == parentID {
			foundParent = p
			break
		}
	}

	if foundParent == nil {
		return errors.NewNotFoundError("Parent", parentID, nil)
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

// Divorce handles the divorce process, creating a new family for the remaining parent
// The original family keeps the custodial parent and children
func (f *Family) Divorce(custodialParentID string) (*Family, error) {
	if f.status != Married {
		return nil, errors.NewDomainError(errors.BusinessRuleViolationCode, "only married families can divorce", nil)
	}

	if len(f.parents) != 2 {
		return nil, errors.NewDomainError(errors.BusinessRuleViolationCode, "divorce requires exactly two parents", nil)
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
		return nil, errors.NewNotFoundError("Parent", custodialParentID, nil)
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
		return nil, errors.NewDomainError(errors.BusinessRuleViolationCode, "failed to create new family for remaining parent", err)
	}

	// Update the original family to keep only the custodial parent
	f.parents = []*Parent{custodialParent}
	f.status = Divorced

	// Return the new family with the remaining parent
	// The original family (with custodial parent and children) is modified in place
	return remainingFamily, nil
}

// ToDTO converts the Family aggregate to a data transfer object for external use
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
		ID:            f.id,
		Status:        string(f.status),
		Parents:       parentDTOs,
		Children:      childDTOs,
		ParentCount:   f.CountParents(),
		ChildrenCount: f.CountChildren(),
	}
}

// FamilyDTO is a data transfer object for the Family aggregate
type FamilyDTO struct {
	ID            string
	Status        string
	Parents       []ParentDTO
	Children      []ChildDTO
	ParentCount   int
	ChildrenCount int
}

// FamilyFromDTO creates a Family aggregate from a data transfer object
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
