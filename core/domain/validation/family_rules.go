// Copyright (c) 2025 A Bit of Help, Inc.

package validation

import (
	"context"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/servicelib/errors"
)

// ParentAgeRule validates that parents meet minimum age requirements
type ParentAgeRule struct {
	minimumAge int
}

// NewParentAgeRule creates a new ParentAgeRule
func NewParentAgeRule(minimumAge int) *ParentAgeRule {
	return &ParentAgeRule{
		minimumAge: minimumAge,
	}
}

// Validate checks that all parents in the family meet the minimum age requirement
func (r *ParentAgeRule) Validate(ctx context.Context, e interface{}) error {
	family, ok := e.(*entity.Family)
	if !ok {
		return errors.NewValidationError("entity is not a Family", "entity", nil)
	}

	now := time.Now()
	for _, parent := range family.Parents() {
		birthDate := parent.BirthDate()
		age := now.Year() - birthDate.Year()

		// Adjust age if birthday hasn't occurred yet this year
		if now.Month() < birthDate.Month() || (now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
			age--
		}

		if age < r.minimumAge {
			return errors.NewValidationError("parent does not meet minimum age requirement", "Parent", nil)
		}
	}

	return nil
}

// ChildBirthDateRule validates that children's birth dates are after parents' birth dates
type ChildBirthDateRule struct{}

// NewChildBirthDateRule creates a new ChildBirthDateRule
func NewChildBirthDateRule() *ChildBirthDateRule {
	return &ChildBirthDateRule{}
}

// Validate checks that all children in the family were born after their parents
func (r *ChildBirthDateRule) Validate(ctx context.Context, e interface{}) error {
	family, ok := e.(*entity.Family)
	if !ok {
		return errors.NewValidationError("entity is not a Family", "entity", nil)
	}

	parents := family.Parents()
	children := family.Children()

	for _, child := range children {
		childBirthDate := child.BirthDate()

		for _, parent := range parents {
			parentBirthDate := parent.BirthDate()

			if !childBirthDate.After(parentBirthDate) {
				return errors.NewValidationError("child's birth date before parent's birth date", "Child", nil)
			}
		}
	}

	return nil
}

// FamilyStatusRule validates that the family status is consistent with the number of parents
type FamilyStatusRule struct{}

// NewFamilyStatusRule creates a new FamilyStatusRule
func NewFamilyStatusRule() *FamilyStatusRule {
	return &FamilyStatusRule{}
}

// Validate checks that the family status is consistent with the number of parents
func (r *FamilyStatusRule) Validate(ctx context.Context, e interface{}) error {
	family, ok := e.(*entity.Family)
	if !ok {
		return errors.NewValidationError("entity is not a Family", "entity", nil)
	}

	status := family.Status()
	parentCount := family.CountParents()

	if status == entity.Status("MARRIED") && parentCount != 2 {
		return errors.NewValidationError("married family must have exactly two parents", "Status", nil)
	}

	if status == entity.Status("SINGLE") && parentCount > 1 {
		return errors.NewValidationError("single family cannot have more than one parent", "Status", nil)
	}

	return nil
}

// CreateFamilyValidationPipeline creates a validation pipeline for family entities
func CreateFamilyValidationPipeline() *Pipeline {
	return NewPipeline(
		NewParentAgeRule(18),
		NewChildBirthDateRule(),
		NewFamilyStatusRule(),
	)
}
