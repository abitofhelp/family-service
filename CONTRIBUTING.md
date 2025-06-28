# Contributing to ServiceLib

Thank you for considering contributing to ServiceLib! This document provides guidelines and instructions for contributing to this project.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct. Please be respectful and considerate of others when contributing.

## How to Contribute

### Reporting Bugs

If you find a bug in the project, please create an issue on GitHub with the following information:

1. A clear, descriptive title
2. A detailed description of the issue
3. Steps to reproduce the bug
4. Expected behavior
5. Actual behavior
6. Any relevant logs or error messages
7. Your environment (Go version, OS, etc.)

### Suggesting Enhancements

If you have an idea for an enhancement, please create an issue on GitHub with the following information:

1. A clear, descriptive title
2. A detailed description of the enhancement
3. Why you believe this enhancement would be valuable
4. Any relevant examples or use cases

### Pull Requests

1. Fork the repository
2. Create a new branch for your changes
3. Make your changes
4. Run tests to ensure your changes don't break existing functionality
5. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.18 or higher
- Git

### Setting Up the Development Environment

1. Clone the repository:
   ```bash
   git clone https://github.com/abitofhelp/servicelib.git
   cd servicelib
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run tests:
   ```bash
   make test
   ```

## Coding Standards

### Go Code Style

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Write meaningful comments and documentation
- Write tests for your code

### Commit Messages

- Use clear, descriptive commit messages
- Start with a short summary line (50 chars or less)
- Optionally, follow with a blank line and a more detailed explanation
- Reference issues and pull requests where appropriate

### Documentation

- Update documentation when changing code
- Follow the existing documentation style
- Use proper grammar and spelling

## Testing

- Write unit tests for all new functionality
- Ensure all tests pass before submitting a pull request
- Aim for high test coverage

## License

By contributing to ServiceLib, you agree that your contributions will be licensed under the project's MIT License.

## Questions

If you have any questions about contributing, please open an issue or contact the project maintainers.

Thank you for your contributions!