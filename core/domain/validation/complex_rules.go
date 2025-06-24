// Copyright (c) 2025 A Bit of Help, Inc.

package validation

import (
	"context"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/abitofhelp/servicelib/errors"
)

// ComplexRuleValidator is a service that validates complex business rules
// that span multiple entities or require context-specific validation
type ComplexRuleValidator struct {
	pipeline *Pipeline
}

// NewComplexRuleValidator creates a new ComplexRuleValidator
func NewComplexRuleValidator() *ComplexRuleValidator {
	// Create a pipeline with complex business rules
	pipeline := NewPipeline(
		NewFamilyConsistencyRule(),
		NewParentChildAgeGapRule(),
		NewFamilyStatusConsistencyRule(),
	)

	return &ComplexRuleValidator{
		pipeline: pipeline,
	}
}

// ValidateFamily validates a family against complex business rules
func (v *ComplexRuleValidator) ValidateFamily(ctx context.Context, family *entity.Family) error {
	return v.pipeline.Validate(ctx, family)
}

// FamilyConsistencyRule validates that a family's structure is consistent
type FamilyConsistencyRule struct{}

// NewFamilyConsistencyRule creates a new FamilyConsistencyRule
func NewFamilyConsistencyRule() *FamilyConsistencyRule {
	return &FamilyConsistencyRule{}
}

// Validate checks that a family's structure is consistent
func (r *FamilyConsistencyRule) Validate(ctx context.Context, e interface{}) error {
	family, ok := e.(*entity.Family)
	if !ok {
		return errors.NewValidationError("entity is not a Family", "entity", nil)
	}

 // Check that married families have exactly two parents
	if family.Status() == entity.Status("MARRIED") && family.CountParents() != 2 {
		return errors.NewValidationError("married family must have exactly two parents", "Family", nil)
	}

	// Check that single families have exactly one parent
	if family.Status() == entity.Status("SINGLE") && family.CountParents() != 1 {
		return errors.NewValidationError("single family must have exactly one parent", "Family", nil)
	}

	// Check that divorced families have exactly one parent
	if family.Status() == entity.Status("DIVORCED") && family.CountParents() != 1 {
		return errors.NewValidationError("divorced family must have exactly one parent", "Family", nil)
	}

	// Check that widowed families have exactly one parent
	if family.Status() == entity.Status("WIDOWED") && family.CountParents() != 1 {
		return errors.NewValidationError("widowed family must have exactly one parent", "Family", nil)
	}

	return nil
}

// ParentChildAgeGapRule validates that there is a reasonable age gap between parents and children
type ParentChildAgeGapRule struct {
	minimumAgeGap int
}

// NewParentChildAgeGapRule creates a new ParentChildAgeGapRule
func NewParentChildAgeGapRule() *ParentChildAgeGapRule {
	return &ParentChildAgeGapRule{
		minimumAgeGap: 12, // Minimum 12 years between parent and child
	}
}

// Validate checks that there is a reasonable age gap between parents and children
func (r *ParentChildAgeGapRule) Validate(ctx context.Context, e interface{}) error {
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

			// Calculate the age difference in years
			ageGap := childBirthDate.Year() - parentBirthDate.Year()

			// Adjust for partial years
			if childBirthDate.Month() < parentBirthDate.Month() || 
			   (childBirthDate.Month() == parentBirthDate.Month() && childBirthDate.Day() < parentBirthDate.Day()) {
				ageGap--
			}

			if ageGap < r.minimumAgeGap {
				return errors.NewValidationError("too small age gap between parent and child", "Family", nil)
			}
		}
	}

	return nil
}

// FamilyStatusConsistencyRule validates that a family's status is consistent with its members
type FamilyStatusConsistencyRule struct{}

// NewFamilyStatusConsistencyRule creates a new FamilyStatusConsistencyRule
func NewFamilyStatusConsistencyRule() *FamilyStatusConsistencyRule {
	return &FamilyStatusConsistencyRule{}
}

// Validate checks that a family's status is consistent with its members
func (r *FamilyStatusConsistencyRule) Validate(ctx context.Context, e interface{}) error {
	family, ok := e.(*entity.Family)
	if !ok {
		return errors.NewValidationError("entity is not a Family", "entity", nil)
	}

	// Check that widowed families don't have deceased parents
	if family.Status() == entity.Status("WIDOWED") {
		for _, parent := range family.Parents() {
			if parent.IsDeceased() {
				return errors.NewValidationError("widowed family cannot have a deceased parent", "Family", nil)
			}
		}
	}

	// Check that abandoned families have at least one child
	if family.Status() == entity.Status("ABANDONED") && family.CountChildren() == 0 {
		return errors.NewValidationError("abandoned family must have at least one child", "Family", nil)
	}

	return nil
}
