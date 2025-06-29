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

## Installation

```bash
go get github.com/abitofhelp/family-service/core/domain/validation
```

## Quick Start

See the [Quick Start example](../../../EXAMPLES/validation/basic_usage/README.md) for a complete, runnable example of how to use the validation framework.

## Configuration

The Domain Validation package can be configured with the following options:

- **Rule Parameters**: Configure parameters for specific rules (e.g., minimum parent age)
- **Pipeline Composition**: Configure which rules are included in the validation pipeline
- **Error Handling**: Configure how validation errors are handled and reported
- **Context Values**: Configure values that are passed through context to validation rules

## Examples

There may be additional examples in the /EXAMPLES directory.

```go
package main

import (
    "github.com/abitofhelp/family-service/core/domain/validation"
)

func main() {
    // Configure the minimum parent age
    minimumParentAge := 18

    // Create a rule with the configured parameter
    parentAgeRule := validation.NewParentAgeRule(minimumParentAge)

    // Configure the validation pipeline
    pipeline := validation.NewPipeline(
        parentAgeRule,
        validation.NewChildBirthDateRule(),
        validation.NewFamilyStatusRule(),
    )
}
```

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

### Key Methods

#### Validate

Validates an entity using the validation pipeline.

```
// Validate applies all rules in the pipeline to the entity
func (p *Pipeline) Validate(ctx context.Context, entity interface{}) error
```

#### NewCompositeRule

Creates a composite rule that combines multiple rules.

```
// NewCompositeRule creates a new composite rule that combines multiple rules
func NewCompositeRule(rules ...Rule) *CompositeRule
```

## Examples

For complete, runnable examples, see the following directories in the EXAMPLES directory:

- [Validation Example](../../../EXAMPLES/validation/README.md) - Shows how to use the validation framework

Example of using the validation pipeline:

```go
package main

import (
    "context"
    "fmt"
    "github.com/abitofhelp/family-service/core/domain/validation"
    "github.com/abitofhelp/family-service/core/domain/entity"
)

func main() {
    // Create validation rules
    parentAgeRule := validation.NewParentAgeRule(18)
    childBirthDateRule := validation.NewChildBirthDateRule()
    familyStatusRule := validation.NewFamilyStatusRule()

    // Create a validation pipeline
    pipeline := validation.NewPipeline(
        parentAgeRule,
        childBirthDateRule,
        familyStatusRule,
    )

    // Create a family to validate
    family := &entity.Family{
        // Family details
    }

    // Validate the family
    ctx := context.Background()
    err := pipeline.Validate(ctx, family)
    if err != nil {
        fmt.Printf("Validation error: %v\n", err)
        return
    }

    fmt.Println("Family is valid!")
}
```

## Best Practices

1. **Single Responsibility**: Each validation rule should focus on a single aspect of validation
2. **Composability**: Use composite rules to combine related validation rules
3. **Clear Error Messages**: Provide clear and descriptive error messages
4. **Context Awareness**: Use context to access contextual information during validation
5. **Separation of Concerns**: Keep validation logic separate from entity logic

## Troubleshooting

### Common Issues

#### Validation Errors

If you encounter validation errors, check that your entities meet all the validation criteria. The error message should provide details about which validation rule failed.

#### Context Cancellation

If you're using a context with a timeout or cancellation, ensure that your validation rules handle context cancellation properly.

## Related Components

- [Domain Entities](../entity/README.md) - The entities being validated
- [Domain Services](../services/README.md) - Services that use these validation rules
- [Validation Wrapper](../../../infrastructure/adapters/validationwrapper/README.md) - Infrastructure for validation results
- [Error Wrapper](../../../infrastructure/adapters/errorswrapper/README.md) - Infrastructure for error handling

## Contributing

Contributions to this component are welcome! Please see the [Contributing Guide](../../../CONTRIBUTING.md) for more information.

## License

This project is licensed under the MIT License - see the [LICENSE](../../../LICENSE) file for details.
