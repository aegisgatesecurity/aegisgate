# Contributing to AegisGate

Thank you for your interest in contributing to AegisGate! We welcome contributions under strict terms to protect our IP.> **⚠️ IMPORTANT: By contributing to this project, you agree to the terms below.**
> 
> All contributions become the exclusive property of AegisGate Security.
> Contributors retain no ownership claims, patents, or intellectual property rights.

---

## Legal Disclaimer

### Intellectual Property Ownership

**By submitting any contribution to AegisGate, you acknowledge and agree that:**

1. **Ownership Transfer**: All contributions, including code, documentation, designs, and ideas, become the sole property of **AegisGate Security**.

2. **No Retained Rights**: You irrevocably transfer all rights, title, and interest in your contributions to the Company.

3. **No Compensation**: Contributions are voluntary and unpaid.

4. **Warranty Representation**: You warrant that your contributions are original to you or properly licensed, and do not infringe third-party rights.

5. **No Patent Claims**: You agree not to assert any patent or trade secret rights in any contribution.

### Disclaimer of Warranties

**CONTRIBUTIONS ARE PROVIDED "AS IS" WITHOUT WARRANTY OF ANY KIND.**

---



---

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Environment](#development-environment)
- [Making Changes](#making-changes)
- [Submitting Changes](#submitting-changes)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)
- [Community Support](#community-support)

---

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md). Please read it before contributing.

**Key Points:**
- Be respectful and inclusive
- Welcome newcomers and help them learn
- Accept constructive criticism professionally
- Focus on what is best for the community

---

## Getting Started

### Prerequisites

| Requirement | Version | Notes |
|-------------|---------|-------|
| Go | 1.21+ | For building from source |
| Docker | Latest | For containerized development |
| Git | Any recent | Version control |
| PostgreSQL | 14+ | For integration testing |

### Clone the Repository

```bash
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate
```

---

## Development Environment

### Local Development with Docker

```bash
# Start local development environment
make dev

# Run tests
make test

# Build binary
make build
```

### Manual Setup

```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o bin/aegisgate ./cmd/aegisgate
```

### Environment Variables

Create a `.env` file for local development:

```bash
# Copy example config
cp config/community.env .env

# Edit your configuration
nano .env
```

---

## Making Changes

### 1. Create a Branch

```bash
# Create a new branch for your feature or bugfix
git checkout -b feature/your-feature-name
# or
git checkout -b fix/bug-description
```

### 2. Make Your Changes

**Keep your changes focused:**
- One feature or fix per branch
- Keep commits atomic and descriptive
- Write tests for new functionality

### 3. Follow Coding Standards

#### Go Standards

```go
// ✓ Good: Clear function names, proper documentation
// GetUser retrieves a user by ID from the database.
func GetUser(ctx context.Context, id string) (*User, error) {
    // implementation
}

// ✗ Bad: Unclear names, no documentation
func get(id string) (*User, error) {
    // implementation
}
```

#### Error Handling

```go
// ✓ Good: Descriptive errors with context
if user == nil {
    return nil, fmt.Errorf("user not found with id %s: %w", id, ErrNotFound)
}

// ✗ Bad: Generic or missing errors
if user == nil {
    return nil, errors.New("error")
}
```

#### Logging

```go
// ✓ Good: Structured logging with appropriate level
logger.Info("request processed",
    "method", r.Method,
    "path", r.URL.Path,
    "duration_ms", duration.Milliseconds(),
)

// ✗ Bad: Unstructured or excessive logging
logger.Info("request processed")
```

---

## Submitting Changes

### 1. Test Your Changes

```bash
# Run all tests
go test ./... -v

# Run specific package tests
go test ./pkg/core/... -v

# Run with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### 2. Commit Your Changes

```bash
# Stage your changes
git add .

# Commit with a descriptive message
git commit -m "feat: add rate limiting middleware

- Implements token bucket algorithm for API rate limiting
- Supports per-tier configuration (Community, Developer, Professional)
- Adds X-RateLimit-* headers to responses

Closes #123"
```

**Commit Message Format:**
```
<type>: <short description>

<long description if needed>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation
- `style`: Code style (formatting)
- `refactor`: Code refactoring
- `test`: Testing
- `chore`: Maintenance

### 3. Push and Create Pull Request

```bash
# Push your branch
git push origin feature/your-feature-name

# Create Pull Request via GitHub UI
# OR use GitHub CLI
gh pr create --fill
```

### 4. Pull Request Checklist

Before submitting, ensure:

- [ ] Tests pass (`go test ./...`)
- [ ] Code follows style guidelines
- [ ] Documentation is updated
- [ ] Commit messages are clear
- [ ] No merge conflicts

---

## Coding Standards

### General Guidelines

1. **Keep it simple** - Prefer simple solutions over clever ones
2. **Be consistent** - Follow existing patterns in the codebase
3. **Write tests** - Aim for meaningful test coverage
4. **Document public APIs** - All exported functions need documentation

### Go-Specific

| Guideline | Example |
|----------|---------|
| Use `context.Context` | `func Foo(ctx context.Context, ...)` |
| Return concrete types | `func Foo() (*MyType, error)` |
| Use named returns | `func Foo() (err error)` for cleanup |
| Group imports | Stdlib, then external, then internal |

### File Organization

```
pkg/
├── component/
│   ├── component.go       # Main implementation
│   ├── component_test.go  # Tests
│   └── doc.go             # Package documentation
```

---

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with race detector
go test -race ./...

# Run specific test
go test -run TestFunctionName ./...
```

### Writing Tests

```go
package core_test

import (
    "testing"

    "github.com/aegisgatesecurity/aegisgate/pkg/core"
)

func TestFeatureX(t *testing.T) {
    // Arrange
    input := "test input"

    // Act
    result, err := core.Process(input)

    // Assert
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result != "expected" {
        t.Errorf("expected %q, got %q", "expected", result)
    }
}
```

### Test Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View in browser
go tool cover -html=coverage.out

# Summary
go tool cover -func=coverage.out
```

---

## Documentation

### Updating Documentation

Documentation lives in the `docs/` directory:

- `README.md` - Main project documentation
- `docs/` - Detailed guides and references
- Code comments - API documentation

### Building Documentation Locally

```bash
# If using a docs generator
make docs
```

### Writing Docs

- Use clear, simple language
- Include code examples
- Keep it up to date with code changes

---

## Community Support

### Getting Help

- **Discord**: [Join our community](https://discord.gg/aegisgate)
- **Forum**: [Community Forum](https://community.aegisgate.example.com)
- **Issues**: [GitHub Issues](https://github.com/aegisgatesecurity/aegisgate/issues)

### Recognitionributors will be recognized

Cont in:
- CONTRIBUTORS.md file
- Release notes
- Social media shoutouts

---

## Thank You!

Your contributions make AegisGate better for everyone. We appreciate your time and effort!
