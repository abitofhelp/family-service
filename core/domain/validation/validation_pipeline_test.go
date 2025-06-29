// Copyright (c) 2025 A Bit of Help, Inc.

package validation

import (
	"context"
	"testing"

	"github.com/abitofhelp/family-service/infrastructure/adapters/errorswrapper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRule is a mock implementation of the Rule interface for testing
type MockRule struct {
	mock.Mock
}

func (m *MockRule) Validate(ctx context.Context, entity interface{}) error {
	args := m.Called(ctx, entity)
	return args.Error(0)
}

// TestNewPipeline tests creating a new validation pipeline
func TestNewPipeline(t *testing.T) {
	// Create mock rules
	rule1 := new(MockRule)
	rule2 := new(MockRule)

	// Create a new pipeline with the mock rules
	pipeline := NewPipeline(rule1, rule2)

	// Assert that the pipeline was created with the correct rules
	assert.Equal(t, 2, len(pipeline.rules), "Pipeline should have 2 rules")
	assert.Equal(t, rule1, pipeline.rules[0], "First rule should be rule1")
	assert.Equal(t, rule2, pipeline.rules[1], "Second rule should be rule2")
}

// TestAddRule tests adding a rule to a validation pipeline
func TestAddRule(t *testing.T) {
	// Create a new pipeline with no rules
	pipeline := NewPipeline()
	assert.Equal(t, 0, len(pipeline.rules), "Pipeline should have 0 rules initially")

	// Create a mock rule
	rule := new(MockRule)

	// Add the rule to the pipeline
	pipeline.AddRule(rule)

	// Assert that the rule was added
	assert.Equal(t, 1, len(pipeline.rules), "Pipeline should have 1 rule after adding")
	assert.Equal(t, rule, pipeline.rules[0], "The added rule should be in the pipeline")
}

// TestValidate_Success tests validating an entity with a pipeline (success case)
func TestValidate_Success(t *testing.T) {
	// Create mock rules that will succeed
	rule1 := new(MockRule)
	rule1.On("Validate", mock.Anything, mock.Anything).Return(nil)

	rule2 := new(MockRule)
	rule2.On("Validate", mock.Anything, mock.Anything).Return(nil)

	// Create a pipeline with the mock rules
	pipeline := NewPipeline(rule1, rule2)

	// Validate an entity
	entity := "test entity"
	err := pipeline.Validate(context.Background(), entity)

	// Assert that validation succeeded
	assert.Nil(t, err, "Validation should succeed")

	// Verify that both rules were called
	rule1.AssertCalled(t, "Validate", mock.Anything, entity)
	rule2.AssertCalled(t, "Validate", mock.Anything, entity)
}

// TestValidate_Failure tests validating an entity with a pipeline (failure case)
func TestValidate_Failure(t *testing.T) {
	// Create a mock rule that will succeed
	rule1 := new(MockRule)
	rule1.On("Validate", mock.Anything, mock.Anything).Return(nil)

	// Create a mock rule that will fail
	rule2 := new(MockRule)
	validationErr := errorswrapper.NewValidationError("validation failed", "field", nil)
	rule2.On("Validate", mock.Anything, mock.Anything).Return(validationErr)

	// Create a pipeline with the mock rules
	pipeline := NewPipeline(rule1, rule2)

	// Validate an entity
	entity := "test entity"
	err := pipeline.Validate(context.Background(), entity)

	// Assert that validation failed
	assert.Error(t, err, "Validation should fail")

	// Verify that both rules were called
	rule1.AssertCalled(t, "Validate", mock.Anything, entity)
	rule2.AssertCalled(t, "Validate", mock.Anything, entity)
}

// TestValidate_MultipleFailures tests validating an entity with a pipeline (multiple failures)
func TestValidate_MultipleFailures(t *testing.T) {
	// Create mock rules that will fail
	rule1 := new(MockRule)
	validationErr1 := errorswrapper.NewValidationError("validation failed 1", "field1", nil)
	rule1.On("Validate", mock.Anything, mock.Anything).Return(validationErr1)

	rule2 := new(MockRule)
	validationErr2 := errorswrapper.NewValidationError("validation failed 2", "field2", nil)
	rule2.On("Validate", mock.Anything, mock.Anything).Return(validationErr2)

	// Create a pipeline with the mock rules
	pipeline := NewPipeline(rule1, rule2)

	// Validate an entity
	entity := "test entity"
	err := pipeline.Validate(context.Background(), entity)

	// Assert that validation failed with multiple errors
	assert.Error(t, err, "Validation should fail")

	// Verify that both rules were called
	rule1.AssertCalled(t, "Validate", mock.Anything, entity)
	rule2.AssertCalled(t, "Validate", mock.Anything, entity)
}

// TestNewContextCompositeRule tests creating a new composite rule
func TestNewContextCompositeRule(t *testing.T) {
	// Create mock rules
	rule1 := new(MockRule)
	rule2 := new(MockRule)

	// Create a new composite rule with the mock rules
	compositeRule := NewContextCompositeRule("test composite", rule1, rule2)

	// Assert that the composite rule was created with the correct name and rules
	assert.Equal(t, "test composite", compositeRule.name, "Composite rule should have the correct name")
	assert.Equal(t, 2, len(compositeRule.rules), "Composite rule should have 2 rules")
	assert.Equal(t, rule1, compositeRule.rules[0], "First rule should be rule1")
	assert.Equal(t, rule2, compositeRule.rules[1], "Second rule should be rule2")
}

// TestContextCompositeRule_Validate_Success tests validating an entity with a composite rule (success case)
func TestContextCompositeRule_Validate_Success(t *testing.T) {
	// Create mock rules that will succeed
	rule1 := new(MockRule)
	rule1.On("Validate", mock.Anything, mock.Anything).Return(nil)

	rule2 := new(MockRule)
	rule2.On("Validate", mock.Anything, mock.Anything).Return(nil)

	// Create a composite rule with the mock rules
	compositeRule := NewContextCompositeRule("test composite", rule1, rule2)

	// Validate an entity
	entity := "test entity"
	err := compositeRule.Validate(context.Background(), entity)

	// Assert that validation succeeded
	assert.Nil(t, err, "Validation should succeed")

	// Verify that both rules were called
	rule1.AssertCalled(t, "Validate", mock.Anything, entity)
	rule2.AssertCalled(t, "Validate", mock.Anything, entity)
}

// TestContextCompositeRule_Validate_Failure tests validating an entity with a composite rule (failure case)
func TestContextCompositeRule_Validate_Failure(t *testing.T) {
	// Create a mock rule that will succeed
	rule1 := new(MockRule)
	rule1.On("Validate", mock.Anything, mock.Anything).Return(nil)

	// Create a mock rule that will fail
	rule2 := new(MockRule)
	validationErr := errorswrapper.NewValidationError("validation failed", "field", nil)
	rule2.On("Validate", mock.Anything, mock.Anything).Return(validationErr)

	// Create a composite rule with the mock rules
	compositeRule := NewContextCompositeRule("test composite", rule1, rule2)

	// Validate an entity
	entity := "test entity"
	err := compositeRule.Validate(context.Background(), entity)

	// Assert that validation failed
	assert.Error(t, err, "Validation should fail")
	assert.Contains(t, err.Error(), "composite rule 'test composite' failed", "Error message should contain the composite rule name")

	// Verify that both rules were called
	rule1.AssertCalled(t, "Validate", mock.Anything, entity)
	rule2.AssertCalled(t, "Validate", mock.Anything, entity)
}

// TestContextCompositeRule_Validate_NonValidationError tests handling non-ValidationError errors
func TestContextCompositeRule_Validate_NonValidationError(t *testing.T) {
	// Create a mock rule that will return a non-ValidationError
	rule := new(MockRule)
	rule.On("Validate", mock.Anything, mock.Anything).Return(assert.AnError)

	// Create a composite rule with the mock rule
	compositeRule := NewContextCompositeRule("test composite", rule)

	// Validate an entity
	entity := "test entity"
	err := compositeRule.Validate(context.Background(), entity)

	// Assert that validation failed
	assert.Error(t, err, "Validation should fail")
	assert.Contains(t, err.Error(), "composite rule 'test composite' failed", "Error message should contain the composite rule name")

	// Verify that the rule was called
	rule.AssertCalled(t, "Validate", mock.Anything, entity)
}
