# AegisGate v1.0.10 Release Notes

**Release Date**: March 15, 2026  
**Version**: v1.0.10  
**Type**: Minor Release

---

## Overview

AegisGate v1.0.10 introduces comprehensive developer tooling improvements with the addition of a production-grade pre-commit framework. This release adds multiple automated code quality checks, security scanning, and validation hooks that run before every commit.

---

## Highlights

### 🚀 New Pre-commit Framework

Integrated the popular [pre-commit](https://pre-commit.com/) framework with 15+ automated checks:

| Hook | Purpose |
|------|---------|
| **golangci-lint** | Comprehensive Go linting (18+ linters) |
| **gosec** | Security vulnerability scanning |
| **misspell** | Spell checking |
| **markdownlint** | Markdown validation |
| **hadolint** | Dockerfile linting |
| **gofmt/go vet** | Code formatting & static analysis |
| **go test -short** | Fast test validation |
| **go-mod-tidy** | Dependency management |
| **File validators** | Merge-conflict, large-files, whitespace, YAML/JSON |

### 🔒 Security Improvements

- Replaced exposed private keys with mock test keys
- gosec integration for Go security scanning
- hadolint for Dockerfile security

---

## Breaking Changes

None - backward-compatible release.

---

## Upgrade from v1.0.9

Pre-commit hooks install automatically on `git pull`. First commit may trigger auto-formatting.

---

## Resources

- **Website**: https://aegisgate.io
- **Documentation**: https://docs.aegisgate.io
- **GitHub**: https://github.com/aegisgatesecurity/aegisgate
- **Discord**: https://discord.gg/aegisgate

---

*Happy coding! 🛡️*