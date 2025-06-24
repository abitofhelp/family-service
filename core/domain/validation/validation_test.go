// Copyright (c) 2025 A Bit of Help, Inc.

package validation

import (
	"context"
	"testing"
	"time"

	"github.com/abitofhelp/family-service/core/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a valid parent
func createValidParent(birthDate time.Time, deathDate *time.Time) *entity.Parent {
	parent, _ := entity.NewParent(uuid.New().String(), "Test", "Parent", birthDate, deathDate)
	return parent
}

// Helper function to create a valid child
func createValidChild(birthDate time.Time) *entity.Child {
	child, _ := entity.NewChild(uuid.New().String(), "Test", "Child", birthDate, nil)
	return child
}

// Helper function to create a valid married family
func createValidMarriedFamily() *entity.Family {
	parent1 := createValidParent(time.Now().AddDate(-30, 0, 0), nil)
	parent2 := createValidParent(time.Now().AddDate(-28, 0, 0), nil)
	child := createValidChild(time.Now().AddDate(-5, 0, 0))
	family, _ := entity.NewFamily(uuid.New().String(), entity.Married, []*entity.Parent{parent1, parent2}, []*entity.Child{child})
	return family
}

// Helper function to create a valid single family
func createValidSingleFamily() *entity.Family {
	parent := createValidParent(time.Now().AddDate(-30, 0, 0), nil)
	child := createValidChild(time.Now().AddDate(-5, 0, 0))
	family, _ := entity.NewFamily(uuid.New().String(), entity.Single, []*entity.Parent{parent}, []*entity.Child{child})
	return family
}

// Helper function to create a valid divorced family
func createValidDivorcedFamily() *entity.Family {
	parent := createValidParent(time.Now().AddDate(-30, 0, 0), nil)
	child := createValidChild(time.Now().AddDate(-5, 0, 0))
	family, _ := entity.NewFamily(uuid.New().String(), entity.Divorced, []*entity.Parent{parent}, []*entity.Child{child})
	return family
}

// Helper function to create a valid widowed family
func createValidWidowedFamily() *entity.Family {
	parent := createValidParent(time.Now().AddDate(-30, 0, 0), nil)
	child := createValidChild(time.Now().AddDate(-5, 0, 0))
	family, _ := entity.NewFamily(uuid.New().String(), entity.Widowed, []*entity.Parent{parent}, []*entity.Child{child})
	return family
}

// Helper function to create a valid abandoned family
func createValidAbandonedFamily() *entity.Family {
	parent := createValidParent(time.Now().AddDate(-30, 0, 0), nil)
	child := createValidChild(time.Now().AddDate(-5, 0, 0))
	family, _ := entity.NewFamily(uuid.New().String(), entity.Abandoned, []*entity.Parent{parent}, []*entity.Child{child})
	return family
}

// TestComplexRuleValidator_ValidateFamily tests the ComplexRuleValidator
func TestComplexRuleValidator_ValidateFamily(t *testing.T) {
	validator := NewComplexRuleValidator()

	t.Run("Valid married family", func(t *testing.T) {
		family := createValidMarriedFamily()
		err := validator.ValidateFamily(context.Background(), family)
		assert.Nil(t, err, "Validation should succeed for a valid married family")
	})

	t.Run("Valid single family", func(t *testing.T) {
		family := createValidSingleFamily()
		err := validator.ValidateFamily(context.Background(), family)
		assert.Nil(t, err, "Validation should succeed for a valid single family")
	})

	t.Run("Valid divorced family", func(t *testing.T) {
		family := createValidDivorcedFamily()
		err := validator.ValidateFamily(context.Background(), family)
		assert.Nil(t, err, "Validation should succeed for a valid divorced family")
	})

	t.Run("Valid widowed family", func(t *testing.T) {
		family := createValidWidowedFamily()
		err := validator.ValidateFamily(context.Background(), family)
		assert.Nil(t, err, "Validation should succeed for a valid widowed family")
	})

	t.Run("Valid abandoned family", func(t *testing.T) {
		family := createValidAbandonedFamily()
		err := validator.ValidateFamily(context.Background(), family)
		assert.Nil(t, err, "Validation should succeed for a valid abandoned family")
	})
}

// TestFamilyConsistencyRule tests the FamilyConsistencyRule
func TestFamilyConsistencyRule(t *testing.T) {
	rule := NewFamilyConsistencyRule()

	t.Run("Invalid entity type", func(t *testing.T) {
		err := rule.Validate(context.Background(), "not a family")
		assert.NotNil(t, err, "Should fail for non-Family entity")
		assert.Contains(t, err.Error(), "entity is not a Family")
	})

	t.Run("Valid married family", func(t *testing.T) {
		family := createValidMarriedFamily()
		err := rule.Validate(context.Background(), family)
		assert.Nil(t, err, "Should pass for valid married family")
	})

	t.Run("Valid single family", func(t *testing.T) {
		family := createValidSingleFamily()
		err := rule.Validate(context.Background(), family)
		assert.Nil(t, err, "Should pass for valid single family")
	})

	t.Run("Valid divorced family", func(t *testing.T) {
		family := createValidDivorcedFamily()
		err := rule.Validate(context.Background(), family)
		assert.Nil(t, err, "Should pass for valid divorced family")
	})

	t.Run("Valid widowed family", func(t *testing.T) {
		family := createValidWidowedFamily()
		err := rule.Validate(context.Background(), family)
		assert.Nil(t, err, "Should pass for valid widowed family")
	})

	// Negative test cases
	t.Run("Married family with one parent", func(t *testing.T) {
		// Create a valid single family first
		family := createValidSingleFamily()

		// Modify it to have a married status (which is invalid with one parent)
		family = &entity.Family{}
		family, _ = entity.NewFamily(uuid.New().String(), entity.Married, []*entity.Parent{createValidParent(time.Now().AddDate(-30, 0, 0), nil)}, []*entity.Child{createValidChild(time.Now().AddDate(-5, 0, 0))})

		// If NewFamily returns nil due to validation, skip this test
		if family == nil {
			t.Skip("Could not create test family due to validation")
		}

		err := rule.Validate(context.Background(), family)
		assert.NotNil(t, err, "Should fail for married family with one parent")
		assert.Contains(t, err.Error(), "married family must have exactly two parents")
	})

	t.Run("Single family with two parents", func(t *testing.T) {
		// Create a valid married family first
		family := createValidMarriedFamily()

		// Modify it to have a single status (which is invalid with two parents)
		family = &entity.Family{}
		family, _ = entity.NewFamily(uuid.New().String(), entity.Single, []*entity.Parent{createValidParent(time.Now().AddDate(-30, 0, 0), nil), createValidParent(time.Now().AddDate(-28, 0, 0), nil)}, []*entity.Child{createValidChild(time.Now().AddDate(-5, 0, 0))})

		// If NewFamily returns nil due to validation, skip this test
		if family == nil {
			t.Skip("Could not create test family due to validation")
		}

		err := rule.Validate(context.Background(), family)
		assert.NotNil(t, err, "Should fail for single family with two parents")
		assert.Contains(t, err.Error(), "single family cannot have more than one parent")
	})
}

// TestParentChildAgeGapRule tests the ParentChildAgeGapRule
func TestParentChildAgeGapRule(t *testing.T) {
	rule := NewParentChildAgeGapRule()

	t.Run("Invalid entity type", func(t *testing.T) {
		err := rule.Validate(context.Background(), "not a family")
		assert.NotNil(t, err, "Should fail for non-Family entity")
		assert.Contains(t, err.Error(), "entity is not a Family")
	})

	t.Run("Valid age gap", func(t *testing.T) {
		family := createValidSingleFamily()
		err := rule.Validate(context.Background(), family)
		assert.Nil(t, err, "Should pass for valid age gap")
	})

	t.Run("Edge case - exactly minimum age gap", func(t *testing.T) {
		// Create a parent born exactly 12 years before the child
		parentBirthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		childBirthDate := time.Date(2012, 1, 1, 0, 0, 0, 0, time.UTC)

		parent := createValidParent(parentBirthDate, nil)
		child := createValidChild(childBirthDate)

		family, _ := entity.NewFamily(uuid.New().String(), entity.Single, []*entity.Parent{parent}, []*entity.Child{child})

		err := rule.Validate(context.Background(), family)
		assert.Nil(t, err, "Should pass for exactly minimum age gap")
	})

	// Negative test cases
	t.Run("Insufficient age gap", func(t *testing.T) {
		// Try to create a family with insufficient age gap
		parentBirthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		childBirthDate := time.Date(2011, 1, 1, 0, 0, 0, 0, time.UTC) // 11 years gap

		parent := createValidParent(parentBirthDate, nil)
		child := createValidChild(childBirthDate)

		family, _ := entity.NewFamily(uuid.New().String(), entity.Single, []*entity.Parent{parent}, []*entity.Child{child})

		// If NewFamily returns nil due to validation, skip this test
		if family == nil {
			t.Skip("Could not create test family due to validation")
		}

		err := rule.Validate(context.Background(), family)
		assert.NotNil(t, err, "Should fail for insufficient age gap")
		assert.Contains(t, err.Error(), "too small age gap between parent and child")
	})

	t.Run("Edge case - just below minimum age gap", func(t *testing.T) {
		// Try to create a family with age gap just below minimum
		parentBirthDate := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
		childBirthDate := time.Date(2011, 12, 31, 0, 0, 0, 0, time.UTC) // 11 years and 364 days

		parent := createValidParent(parentBirthDate, nil)
		child := createValidChild(childBirthDate)

		family, _ := entity.NewFamily(uuid.New().String(), entity.Single, []*entity.Parent{parent}, []*entity.Child{child})

		// If NewFamily returns nil due to validation, skip this test
		if family == nil {
			t.Skip("Could not create test family due to validation")
		}

		err := rule.Validate(context.Background(), family)
		assert.NotNil(t, err, "Should fail for age gap just below minimum")
		assert.Contains(t, err.Error(), "too small age gap between parent and child")
	})
}

// TestFamilyStatusConsistencyRule tests the FamilyStatusConsistencyRule
func TestFamilyStatusConsistencyRule(t *testing.T) {
	rule := NewFamilyStatusConsistencyRule()

	t.Run("Invalid entity type", func(t *testing.T) {
		err := rule.Validate(context.Background(), "not a family")
		assert.NotNil(t, err, "Should fail for non-Family entity")
		assert.Contains(t, err.Error(), "entity is not a Family")
	})

	t.Run("Valid widowed family", func(t *testing.T) {
		family := createValidWidowedFamily()
		err := rule.Validate(context.Background(), family)
		assert.Nil(t, err, "Should pass for valid widowed family")
	})

	t.Run("Valid abandoned family", func(t *testing.T) {
		family := createValidAbandonedFamily()
		err := rule.Validate(context.Background(), family)
		assert.Nil(t, err, "Should pass for valid abandoned family")
	})

	// Negative test cases
	t.Run("Widowed family with deceased parent", func(t *testing.T) {
		// Try to create a widowed family with a deceased parent
		deathDate := time.Now().AddDate(-1, 0, 0)
		parent := createValidParent(time.Now().AddDate(-30, 0, 0), &deathDate)
		child := createValidChild(time.Now().AddDate(-5, 0, 0))

		family, _ := entity.NewFamily(uuid.New().String(), entity.Widowed, []*entity.Parent{parent}, []*entity.Child{child})

		// If NewFamily returns nil due to validation, skip this test
		if family == nil {
			t.Skip("Could not create test family due to validation")
		}

		err := rule.Validate(context.Background(), family)
		assert.NotNil(t, err, "Should fail for widowed family with deceased parent")
		assert.Contains(t, err.Error(), "widowed family cannot have a deceased parent")
	})

	t.Run("Abandoned family without children", func(t *testing.T) {
		// Try to create an abandoned family without children
		parent := createValidParent(time.Now().AddDate(-30, 0, 0), nil)

		family, _ := entity.NewFamily(uuid.New().String(), entity.Abandoned, []*entity.Parent{parent}, []*entity.Child{})

		// If NewFamily returns nil due to validation, skip this test
		if family == nil {
			t.Skip("Could not create test family due to validation")
		}

		err := rule.Validate(context.Background(), family)
		assert.NotNil(t, err, "Should fail for abandoned family without children")
		assert.Contains(t, err.Error(), "abandoned family must have at least one child")
	})
}
