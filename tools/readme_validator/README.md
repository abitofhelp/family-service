# README Validator

## Overview

The README Validator is a tool that ensures all README.md files in the project follow a consistent structure and format. It validates that each README.md file contains the required sections based on its location in the project and checks for consistent formatting of example paths.

## Architecture

The README Validator follows a simple command-line tool architecture. It recursively searches for all README.md files in the project and validates each one against a set of rules. The tool is designed to be run as part of the CI/CD pipeline to ensure consistency across all documentation.

The architecture consists of:

- **Main Function**: Entry point that coordinates the validation process
- **File Discovery**: Finds all README.md files in the project
- **Validation Logic**: Checks each README.md file against the required structure
- **Error Reporting**: Reports validation errors in a clear and actionable format

## Implementation Details

The README Validator implements the following design patterns:

1. **Visitor Pattern**: Visits each README.md file in the project to validate it
2. **Strategy Pattern**: Uses different validation strategies based on the file's location
3. **Command Pattern**: Implements a command-line interface for running the validation

Key implementation details:

- **Regular Expressions**: Uses regular expressions to identify section headings and example paths
- **File System Traversal**: Uses filepath.Walk to recursively find all README.md files
- **Error Aggregation**: Collects all validation errors for comprehensive reporting
- **Exit Codes**: Returns non-zero exit code if validation fails, enabling integration with CI/CD pipelines

## Examples

Example usage of the README Validator:

```
# Validate all README.md files in the current directory and subdirectories
go run tools/readme_validator/main.go

# Validate all README.md files in a specific directory
go run tools/readme_validator/main.go path/to/directory
```

## Configuration

The README Validator has two sets of required sections:

1. **Component README Files**:
   - Overview
   - Architecture
   - Implementation Details
   - Examples
   - Configuration
   - Testing
   - Design Notes
   - References

2. **Example README Files**:
   - Overview
   - Features
   - Running the Example
   - Code Walkthrough
   - Expected Output
   - Related Examples
   - Related Components
   - License

The validator automatically determines which set of required sections to use based on the file's location:
- README.md files in the EXAMPLES directory use the Example README template
- All other README.md files use the Component README template

## Testing

The README Validator is tested through:

1. **Manual Testing**: Running the validator on the project to ensure it correctly identifies validation issues
2. **Edge Case Testing**: Testing with various README.md files to ensure all validation rules are correctly applied

Key testing approaches:

- **Valid Files**: Testing with valid README.md files to ensure they pass validation
- **Invalid Files**: Testing with invalid README.md files to ensure they fail validation with the correct error messages
- **Special Cases**: Testing with special cases like the main README.md and template files that are excluded from validation

## Design Notes

1. **Simplicity**: The validator is designed to be simple and easy to understand
2. **Consistency**: The validator ensures consistency across all README.md files in the project
3. **Automation**: The validator can be integrated into CI/CD pipelines for automated validation
4. **Flexibility**: The validator supports different templates for different types of README.md files
5. **Extensibility**: The validator can be easily extended to support additional validation rules

## References

- [Markdown Guide](https://www.markdownguide.org/) - Guide to Markdown syntax
- [Component README Template](../../COMPONENT_README_TEMPLATE.md) - Template for component README.md files
- [Example README Template](../../EXAMPLES/EXAMPLE_README_TEMPLATE.md) - Template for example README.md files
- [Go filepath Package](https://golang.org/pkg/path/filepath/) - Used for file system traversal
- [Go regexp Package](https://golang.org/pkg/regexp/) - Used for regular expression matching