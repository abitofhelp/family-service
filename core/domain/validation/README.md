# Domain Validation

## Overview

The Domain Validation package provides a flexible and extensible framework for validating domain entities. It implements a pipeline-based approach to validation, allowing multiple validation rules to be applied to an entity in a consistent manner. The package includes both generic validation infrastructure and specific rules for validating family entities.

## Features

- **Validation Pipeline**: Applies multiple validation rules to an entity
- **Composite Rules**: Combines multiple rules into a single rule
- **Family-Specific Rules**: Includes rules for validating family entities
- **Extensible Framework**: Easily add new validation rules
- **Context-Aware**: Validation rules have access to context
- **Error Aggregation**: Collects and reports multiple validation errors
- **Clean Separation**: Keeps validation logic separate from entity logic

## API Documentation

### Core Types

#### Rule

The Rule interface defines the contract for validation rules.

```
// Rule represents a validation rule that can be applied to a domain entity
type Rule interface {
    // Validate applies the rule to the entity and returns an error if validation fails
    Validate(ctx context.Context, entity interface{}) error
}
```

#### Pipeline

The Pipeline type applies multiple validation rules to an entity.

```
// Pipeline represents a validation pipeline that can apply multiple rules to an entity
type Pipeline struct {
    rules []Rule
}

// NewPipeline creates a new validation pipeline
func NewPipeline(rules ...Rule) *Pipeline
```

### Family-Specific Rules

#### ParentAgeRule

Validates that parents meet minimum age requirements.

```
// ParentAgeRule validates that parents meet minimum age requirements
type ParentAgeRule struct {
    minimumAge int
}

// NewParentAgeRule creates a new ParentAgeRule
func NewParentAgeRule(minimumAge int) *ParentAgeRule
```

#### ChildBirthDateRule

Validates that children's birth dates are after parents' birth dates.

```
// ChildBirthDateRule validates that children's birth dates are after parents' birth dates
type ChildBirthDateRule struct{}

// NewChildBirthDateRule creates a new ChildBirthDateRule
func NewChildBirthDateRule() *ChildBirthDateRule
```

#### FamilyStatusRule

Validates that the family status is consistent with the number of parents.

```
// FamilyStatusRule validates that the family status is consistent with the number of parents
type FamilyStatusRule struct{}

// NewFamilyStatusRule creates a new FamilyStatusRule
func NewFamilyStatusRule() *FamilyStatusRule
```

## Best Practices

1. **Single Responsibility**: Each validation rule should focus on a single aspect of validation
2. **Composability**: Use composite rules to combine related validation rules
3. **Clear Error Messages**: Provide clear and descriptive error messages
4. **Context Awareness**: Use context to access contextual information during validation
5. **Separation of Concerns**: Keep validation logic separate from entity logic

## Related Components

- [Domain Entities](../entity/README.md) - The entities being validated
- [Domain Services](../services/README.md) - Services that use these validation rules
- [Validation Wrapper](../../../infrastructure/adapters/validationwrapper/README.md) - Infrastructure for validation results
- [Error Wrapper](../../../infrastructure/adapters/errorswrapper/README.md) - Infrastructure for error handling