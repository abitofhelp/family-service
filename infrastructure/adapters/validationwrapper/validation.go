// Copyright (c) 2025 A Bit of Help, Inc.

// Package validationwrapper provides validation utilities for the domain layer
// without direct dependencies on external libraries.
package validationwrapper

import (
	"fmt"

	"github.com/abitofhelp/family-service/infrastructure/adapters/errorswrapper"
)

// ValidationRule defines a validation rule that can be applied to an entity
type ValidationRule interface {
	// Validate applies the validation rule to the given entity
	Validate(entity interface{}) error
	// Name returns the name of the rule
	Name() string
}

// simpleRule implements a simple validation rule
type simpleRule struct {
	name         string
	validateFunc func(entity interface{}) error
}

// NewValidationRule creates a new validation rule
func NewValidationRule(name string, validateFunc func(entity interface{}) error) ValidationRule {
	return &simpleRule{
		name:         name,
		validateFunc: validateFunc,
	}
}

// Validate applies the validation rule to the given entity
func (r *simpleRule) Validate(entity interface{}) error {
	return r.validateFunc(entity)
}

// Name returns the name of the rule
func (r *simpleRule) Name() string {
	return r.name
}

// ValidationPipeline defines a pipeline of validation rules
type ValidationPipeline interface {
	// AddRule adds a validation rule to the pipeline
	AddRule(rule ValidationRule) ValidationPipeline
	// Validate applies all validation rules to the given entity
	Validate(entity interface{}) error
}

// pipeline implements a validation pipeline
type pipeline struct {
	rules []ValidationRule
}

// NewValidationPipeline creates a new validation pipeline
func NewValidationPipeline() ValidationPipeline {
	return &pipeline{
		rules: make([]ValidationRule, 0),
	}
}

// AddRule adds a validation rule to the pipeline
func (p *pipeline) AddRule(rule ValidationRule) ValidationPipeline {
	p.rules = append(p.rules, rule)
	return p
}

// Validate applies all validation rules to the given entity
func (p *pipeline) Validate(entity interface{}) error {
	for _, rule := range p.rules {
		if err := rule.Validate(entity); err != nil {
			return err
		}
	}
	return nil
}

// CompositeRule defines a composite validation rule that combines multiple rules
type CompositeRule interface {
	ValidationRule
	// AddRule adds a validation rule to the composite rule
	AddRule(rule ValidationRule) CompositeRule
}

// compositeRule implements a composite validation rule
type compositeRule struct {
	name  string
	rules []ValidationRule
}

// NewCompositeRule creates a new composite validation rule
func NewCompositeRule(name string) CompositeRule {
	return &compositeRule{
		name:  name,
		rules: make([]ValidationRule, 0),
	}
}

// AddRule adds a validation rule to the composite rule
func (r *compositeRule) AddRule(rule ValidationRule) CompositeRule {
	r.rules = append(r.rules, rule)
	return r
}

// Validate applies all validation rules to the given entity
func (r *compositeRule) Validate(entity interface{}) error {
	for _, rule := range r.rules {
		if err := rule.Validate(entity); err != nil {
			return err
		}
	}
	return nil
}

// Name returns the name of the rule
func (r *compositeRule) Name() string {
	return r.name
}

// ValidationResult represents the result of a validation
type ValidationResult interface {
	// AddError adds an error to the validation result
	AddError(message string, field string)
	// Error returns an error if there are validation errors, or nil if there are none
	Error() error
}

// validationResult implements a validation result
type validationResult struct {
	errors []validationError
}

// validationError represents a validation error
type validationError struct {
	message string
	field   string
}

// NewValidationResult creates a new validation result
func NewValidationResult() ValidationResult {
	return &validationResult{
		errors: make([]validationError, 0),
	}
}

// AddError adds an error to the validation result
func (r *validationResult) AddError(message string, field string) {
	r.errors = append(r.errors, validationError{
		message: message,
		field:   field,
	})
}

// Error returns an error if there are validation errors, or nil if there are none
func (r *validationResult) Error() error {
	if len(r.errors) == 0 {
		return nil
	}

	// Create a validation error with the first error message
	firstError := r.errors[0]
	return errorswrapper.NewValidationError(firstError.message, firstError.field, nil)
}

// Helper functions for common validations

// ValidateNotNil validates that a value is not nil
func ValidateNotNil(value interface{}, fieldName string) error {
	if value == nil {
		return errorswrapper.NewValidationError("value cannot be nil", fieldName, nil)
	}
	return nil
}

// ValidateNotEmpty validates that a string is not empty
func ValidateNotEmpty(value string, fieldName string) error {
	if value == "" {
		return errorswrapper.NewValidationError("value cannot be empty", fieldName, nil)
	}
	return nil
}

// ValidateMinLength validates that a string has a minimum length
func ValidateMinLength(value string, minLength int, fieldName string) error {
	if len(value) < minLength {
		return errorswrapper.NewValidationError(fmt.Sprintf("value must have a minimum length of %d", minLength), fieldName, nil)
	}
	return nil
}

// ValidateMaxLength validates that a string has a maximum length
func ValidateMaxLength(value string, maxLength int, fieldName string) error {
	if len(value) > maxLength {
		return errorswrapper.NewValidationError(fmt.Sprintf("value must have a maximum length of %d", maxLength), fieldName, nil)
	}
	return nil
}
