# Release Notes - v0.28.2

## Overview

**Release Date:** March 4, 2026
**Version:** v0.28.2
**Type:** Feature Release

---

## What's New in v0.28.2

This release introduces the **Resilience Package** for fault-tolerant proxy connections, along with CLI improvements and version synchronization enhancements.

### New Features

#### 1. Resilience Package (pkg/resilience)

A comprehensive resilience framework for building fault-tolerant proxy connections:

- **Circuit Breaker Pattern**
  - Three-state implementation: Closed, Open, Half-Open
  - Configurable failure thresholds
  - Automatic state transitions
  - Health endpoint metrics integration

- **Retry with Exponential Backoff**
  - Configurable retry attempts
  - Exponential backoff with jitter
  - Deadline-aware retries

- **Timeout Executor**
  - Per-operation timeout handling
  - Context-aware cancellation
  - Configurable timeout durations

- **ResilientClient**
  - Combines all resilience patterns
  - Easy-to-use HTTP client wrapper
  - Pluggable configuration via Options struct

- **Proxy Integration**
  - Circuit breaker integrated into proxy
  - Automatic failure prevention
  - Request queuing during Open state

#### 2. CLI Improvements

- **New CLI Flags**
  - --help: Display comprehensive help message
  - --version: Show version information with build details

- **Version Information**
  - Version number display
  - Build date
  - Git commit hash
  - Git branch information

#### 3. Version Synchronization

- **Version Check CI Workflow** (.github/workflows/version-check.yml)
  - Automated VERSION file validation
  - main.go version consistency check
  - Git tag alignment verification
  - Go version consistency across workflows

- **Version Sync Script** (scripts/version-sync.sh)
  - Local version consistency verification
  - Git tag comparison
  - Automated version validation

## Technical Details

### Dependencies Updated

- github.com/prometheus/client_golang v1.20.5
- golang.org/x/net v0.35.0
- golang.org/x/oauth2 v0.35.0

### Go Version

- **Minimum Required:** Go 1.24.0
- **Tested Against:** Go 1.24.0

### Build Improvements

- Go version aligned across all CI workflows (1.24.0)
- Binary cleanup in repository root
- Makefile updated to reflect current version

## Security

- Circuit breaker prevents cascade failures
- Timeout enforcement prevents resource exhaustion
- Improved audit logging

## Breaking Changes

**None** - This release is fully backward compatible with v0.28.1.

## Migration Guide

### Upgrading from v0.28.1

1. Update your installation:
`ash
git pull origin main
go build -o aegisgate ./cmd/aegisgate
`

2. Test the new features:
`ash
./aegisgate --version
./aegisgate --help
`

## Project Statistics

| Metric | Value |
|--------|-------|
| Total Packages | 30+ |
| Go Files | 200+ |
| CI/CD Workflows | 7 |
| Compliance Frameworks | 14+ |

## Known Issues

None reported.

## Support

- **Issues:** Report via GitHub Issues
- **Discussions:** Use GitHub Discussions
- **Security:** See SECURITY.md

---

*Thank you for using AegisGate!*

**Protect Your AI Applications with AegisGate**
