// Copyright (c) 2025 A Bit of Help, Inc.

package validation

import (
	"context"
	"fmt"

	"github.com/abitofhelp/family-service/infrastructure/adapters/errorswrapper"
	"github.com/abitofhelp/family-service/infrastructure/adapters/validationwrapper"
)

// Rule represents a validation rule that can be applied to a domain entity
type Rule interface {
	// Validate applies the rule to the entity and returns an error if validation fails
	Validate(ctx context.Context, entity interface{}) error
}

// Pipeline represents a validation pipeline that can apply multiple rules to an entity
type Pipeline struct {
	rules []Rule
}

// NewPipeline creates a new validation pipeline
func NewPipeline(rules ...Rule) *Pipeline {
	return &Pipeline{
		rules: rules,
	}
}

// AddRule adds a rule to the pipeline
func (p *Pipeline) AddRule(rule Rule) {
	p.rules = append(p.rules, rule)
}

// Validate applies all rules in the pipeline to the entity
func (p *Pipeline) Validate(ctx context.Context, entity interface{}) error {
	result := validationwrapper.NewValidationResult()

	for _, rule := range p.rules {
		if err := rule.Validate(ctx, entity); err != nil {
			if errorswrapper.IsValidationError(err) {
				// Add the error message to the result
				result.AddError(err.Error(), "")
			} else {
				result.AddError(err.Error(), "")
			}
		}
	}

	return result.Error()
}

// ContextCompositeRule is a rule that combines multiple rules
type ContextCompositeRule struct {
	rules []Rule
	name  string
}

// NewContextCompositeRule creates a new composite rule
func NewContextCompositeRule(name string, rules ...Rule) *ContextCompositeRule {
	return &ContextCompositeRule{
		rules: rules,
		name:  name,
	}
}

// Validate applies all rules in the composite rule to the entity
func (r *ContextCompositeRule) Validate(ctx context.Context, entity interface{}) error {
	result := validationwrapper.NewValidationResult()

	for _, rule := range r.rules {
		if err := rule.Validate(ctx, entity); err != nil {
			if errorswrapper.IsValidationError(err) {
				// Add the error message to the result
				result.AddError(err.Error(), "")
			} else {
				result.AddError(err.Error(), "")
			}
		}
	}

	if err := result.Error(); err != nil {
		return errorswrapper.NewValidationError(fmt.Sprintf("composite rule '%s' failed", r.name), r.name, err)
	}

	return nil
}
