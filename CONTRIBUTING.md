# Contributing to vanity-go

First off, thank you for considering contributing to vanity-go! It's people like you that make vanity-go such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check existing issues as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

- **Use a clear and descriptive title** for the issue to identify the problem
- **Describe the exact steps which reproduce the problem** in as many details as possible
- **Provide specific examples to demonstrate the steps**
- **Describe the behavior you observed after following the steps** and point out what exactly is the problem with that behavior
- **Explain which behavior you expected to see instead and why**
- **Include details about your configuration and environment**

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

- **Use a clear and descriptive title** for the issue to identify the suggestion
- **Provide a step-by-step description of the suggested enhancement** in as many details as possible
- **Provide specific examples to demonstrate the steps**
- **Describe the current behavior** and **explain which behavior you expected to see instead** and why
- **Explain why this enhancement would be useful** to most vanity-go users

### Pull Requests

1. Fork the repo and create your branch from `main`
2. If you've added code that should be tested, add tests
3. If you've changed APIs, update the documentation
4. Ensure the test suite passes
5. Make sure your code follows the existing code style
6. Issue that pull request!

## Development Process

### Prerequisites

- Go 1.21 or later
- Make (optional but recommended)
- golangci-lint (for linting)

### Setting Up Your Development Environment

```bash
# Clone your fork
git clone https://github.com/gllm-dev/vanity-go.git
cd vanity-go

# Add upstream remote
git remote add upstream https://github.com/gllm-dev/vanity-go.git

# Install dependencies
go mod download

# Run the application locally
VANITY_DOMAIN=localhost VANITY_REPOSITORY=https://github.com/gllm-dev go run cmd/main.go
```

### Code Style

We follow the standard Go formatting guidelines:

- Run `go fmt` on your code before committing
- Run `go vet` to catch common mistakes
- Run `golangci-lint run` if you have it installed
- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Write clear, idiomatic Go code

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detector
go test -race ./...

# Run specific test
go test -run TestSpecificFunction ./...
```

### Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line
- When only changing documentation, include `[ci skip]` in the commit title

Example:
```
Add support for custom ports

- Add PORT environment variable
- Update documentation with new configuration option
- Add tests for port configuration

Fixes #123
```

### Documentation

- Keep README.md up to date with any changes
- Document new features or changes in behavior
- Use clear and concise language
- Include examples where appropriate

## Project Structure

```
vanity-go/
├── cmd/
│   └── main.go             # Application entry point
├── internal/
│   ├── adapters/
│   │   └── handlers/      # HTTP handlers
│   │       └── rest/      # REST API handlers
│   └── services/          # Business logic services
│       └── gosvc/         # Go vanity import service
├── di/                    # Dependency injection
├── Dockerfile           # Docker build configuration
├── go.mod              # Go module definition
└── go.sum              # Go module checksums
```

## Pull Request Process

1. Update the README.md with details of changes to the interface, if applicable
2. Update the CHANGELOG.md with your changes under the "Unreleased" section
3. The PR will be merged once you have the sign-off of at least one maintainer

## Release Process

Maintainers will handle releases by:

1. Updating CHANGELOG.md to reflect the new version
2. Creating a git tag for the version
3. Building and pushing Docker images
4. Creating a GitHub release

## Questions?

Feel free to open an issue with your question or reach out to the maintainers directly.

Thank you for contributing!