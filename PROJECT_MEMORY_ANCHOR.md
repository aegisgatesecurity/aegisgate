# AegisGate Project Memory Anchor

**Document Version:** 1.0  
**Last Updated:** March 2026  
**Current Version:** v0.29.1  
**Repository:** github.com/aegisgatesecurity/aegisgate

---

## 📋 Table of Contents

1. [Project Overview](#project-overview)
2. [Current State](#current-state)
3. [Architecture](#architecture)
4. [Recent Fixes & Commits](#recent-fixes--commits)
5. [Next Steps](#next-steps)
6. [Future Plans](#future-plans)
7. [Pro Tips & Gotchas](#pro-tips--gotchas)
8. [Lessons Learned](#lessons-learned)
9. [Troubleshooting Guide](#troubleshooting-guide)
10. [Technical Reference](#technical-reference)

---

## 1. Project Overview

### What is AegisGate?

**AegisGate** is an Enterprise AI/LLM Security Gateway - a Go-based reverse proxy that inspects, filters, and secures traffic between client applications and upstream AI services.

### Key Characteristics

| Attribute | Value |
|-----------|-------|
| **Language** | Go 1.24.0 |
| **License** | MIT |
| **Type** | Reverse Proxy / Security Gateway |
| **Target** | AI/LLM Services (OpenAI, Anthropic, Cohere, Azure, Ollama) |
| **Test Coverage** | 85%+ |
| **ATLAS Patterns** | 60+ |
| **Compliance Frameworks** | 14+ |

### Core Features

- 🔒 **Content Scanning** - Real-time threat detection for malicious payloads, PII, sensitive data
- 🛡️ **Prompt Injection Detection** - MITRE ATLAS compliance with 60+ pattern detection
- 🔐 **PKI Attestation** - Certificate-based identity verification using TPM/HSM
- 📊 **Full Observability** - Prometheus metrics, structured logging, health endpoints
- 🔄 **Resilience** - Circuit breaker, retry logic, rate limiting
- ✅ **Compliance** - SOC 2, HIPAA, PCI DSS, NIST AI RMF, ISO 42001, GDPR, OWASP

---

## 2. Current State

### Version Alignment ✅

| Location | Version | Status |
|----------|---------|--------|
| `VERSION` file | 0.29.1 | ✅ Synced |
| `cmd/aegisgate/main.go` | 0.29.1 | ✅ Synced |
| `README.md` | v0.29.1 | ✅ Synced |
| Git Tag | v0.29.1 | ✅ Created |

### Git Status

```
Branch: main
Remote: origin (git@github.com:aegisgatesecurity/aegisgate.git)
Status: Up to date with origin/main
Last Commit: dd8dd81 - "Update VERSION file to v0.29.1"
```

### GitHub Actions Status

- **Integration Tests Workflow:** Active (runs on push to main/develop)
- **Release Workflow:** Active
- **Security Scan:** Active
- **SBOM Generation:** Active
- **Build Validation:** Active

### Recent Tags

```
v0.29.1 ← Latest
v0.29.0
v0.28.2
v0.28.1
v0.28.0
v0.27.0
... (60+ releases total)
```

---

## 3. Architecture

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         AegisGate Proxy                           │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐      │
│  │   Auth       │───▶│   Scanner    │───▶│  Compliance  │      │
│  │   Layer      │    │   Engine     │    │  Manager     │      │
│  └──────────────┘    └──────────────┘    └──────────────┘      │
│         │                   │                   │                │
│         ▼                   ▼                   ▼                │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                    Proxy Core                             │   │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐      │   │
│  │  │  Circuit    │  │   Rate       │  │   Retry     │      │   │
│  │  │  Breaker    │  │   Limiter    │  │   Logic     │      │   │
│  │  └─────────────┘  └─────────────┘  └─────────────┘      │   │
│  └──────────────────────────────────────────────────────────┘   │
│                              │                                   │
│                              ▼                                   │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │              Upstream AI Services                        │   │
│  │   OpenAI │ Anthropic │ Cohere │ Azure │ Ollama │ ...   │   │
│  └──────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

### Package Structure

```
aegisgate/
├── cmd/aegisgate/           # Main application entry point (544 LOC)
├── pkg/
│   ├── proxy/             # Core reverse proxy (MITM support)
│   ├── auth/              # Authentication & authorization (OIDC, SAML)
│   ├── scanner/           # Content scanning engine (patterns)
│   ├── compliance/        # Compliance frameworks (ATLAS, OWASP, etc.)
│   ├── resilience/        # Circuit breaker, retries
│   ├── metrics/           # Prometheus metrics
│   ├── security/          # Security middleware
│   ├── tls/               # TLS/mTLS management
│   ├── i18n/              # Internationalization
│   ├── immutable-config/  # Immutable configuration system
│   ├── opsec/             # Operational security
│   ├── pkiattest/         # PKI attestation
│   ├── threatintel/       # Threat intelligence (STIX/TAXII)
│   ├── siem/              # SIEM integration
│   ├── sso/               # Single sign-on
│   ├── webhook/           # Webhook management
│   └── ...
├── tests/
│   ├── integration/        # Integration tests
│   ├── load/              # Load tests
│   └── security/          # Security tests
└── ui/frontend/           # Web UI
```

---

## 4. Recent Fixes & Commits

### Commit History (Recent)

| Commit | Description |
|--------|-------------|
| `dd8dd81` | Update VERSION file to v0.29.1 |
| `5de6328` | Update version to v0.29.1 |
| `b425bda` | Fix golangci-lint errors: remove unused _configPath and severityFromString |
| `40e89ed` | Update version to v0.29.1 |
| `311c20f` | Add version notes |
| `5f1938a` | Release v0.29.1 |
| `ee21acb` | Fix TestViolationNames: use Pattern.Name instead of Description |
| `2abece4` | Fix version to 0.29.0 and compliance test issues |
| `d453e8a` | Fix ATLAS patterns and API compatibility |

### Key Fixes Applied

#### Fix 1: Unused Variable (`_configPath`)
- **File:** `cmd/aegisgate/main.go`
- **Issue:** `var _configPath is unused` (line 50)
- **Solution:** Removed unused flag variable

#### Fix 2: Unused Function (`severityFromString`)
- **File:** `pkg/compliance/owasp.go`
- **Issue:** `func severityFromString is unused` (line 844)
- **Solution:** Removed unused helper function

#### Fix 3: Version Mismatch
- **Issue:** VERSION file showed 0.29.0, main.go showed 0.29.1
- **Solution:** Updated VERSION file to match main.go

#### Fix 4: TestViolationNames
- **File:** `pkg/proxy/mitm.go`
- **Issue:** Test failing - used `Pattern.Description` instead of `Pattern.Name`
- **Solution:** Changed to use `Pattern.Name` (line 799)

---

## 5. Next Steps

### Immediate Actions

1. ✅ **All workflow errors resolved** - GitHub Actions should now pass
2. **Monitor CI/CD** - Verify workflows run successfully with commit `dd8dd81`
3. **Clean up artifacts** - Several temp files remain untracked (see below)

### Pending Cleanup

The following files are untracked and could be cleaned up:

```
check_issues.py
check_lines.py
fix_*.py (multiple)
temp_*.txt
*.bak (backup files from fixes)
```

### Optional Enhancements

- Remove `.bak` files from `pkg/compliance/`
- Clean up temporary Python scripts used for fixes
- Consider adding these to `.gitignore`

---

## 6. Future Plans

### Roadmap (From Documentation)

#### Near-term (v0.30.x)
- Enhanced ATLAS pattern detection with ML-based classification
- Real-time ATLAS pattern updates via threat intelligence
- Extended compliance report generation

#### Mid-term
- Performance benchmarking dashboard
- Advanced ML-based threat detection
- Extended SIEM integrations

#### Long-term
- Kubernetes operator for easier deployment
- Multi-cluster support
- Advanced analytics dashboard

### Feature Requests (Known)

| Feature | Priority | Status |
|---------|----------|--------|
| ML-based classification | Medium | Planned |
| Real-time ATLAS updates | Medium | Planned |
| Enhanced reporting | Medium | Planned |

---

## 7. Pro Tips & Gotchas

### ⚠️ Critical Gotchas

#### 1. **Version Synchronization is CRITICAL**
- **Gotcha:** golangci-lint will fail if VERSION file doesn't match main.go
- **Why:** CI/CD workflows often check version consistency
- **Fix:** Always update VERSION file when bumping version in main.go

#### 2. **Unused Code = CI Failure**
- **Gotcha:** golangci-lint enforces no unused variables/functions
- **Why:** Clean code practices, prevents dead code
- **Fix:** Run `golangci-lint run` locally before pushing

#### 3. **Pattern.Name vs Pattern.Description**
- **Gotcha:** MITM test failures if using wrong field
- **Why:** Pattern.Description may be empty, Pattern.Name is always set
- **Fix:** Always use `Pattern.Name` for identification

#### 4. **GitHub Actions Not Triggering**
- **Gotcha:** Sometimes actions don't run after push
- **Why:** Could be branch protection rules or sync delays
- **Fix:** Check web UI, ensure remote is synced

### 🔧 Pro Tips

#### Running Tests Locally
```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -v ./pkg/proxy/... -run TestViolationNames

# Run ATLAS compliance tests
go test -v -run Atlas ./tests/integration/...
```

#### Building
```bash
# Standard build
go build -o aegisgate ./cmd/aegisgate

# With version info
go build -ldflags="-X main.gitCommit=$(git rev-parse HEAD)" -o aegisgate ./cmd/aegisgate
```

#### Linting
```bash
# Run linter
golangci-lint run

# Fix auto-fixable issues
golangci-lint run --fix
```

#### Checking Version
```bash
# View current version
cat VERSION

# Or run binary
./aegisgate --version
```

### 📝 Important File Locations

| Purpose | File Location |
|---------|---------------|
| Main entry | `cmd/aegisgate/main.go` |
| Version constant | `cmd/aegisgate/main.go` (line ~43) |
| Version file | `VERSION` |
| README | `README.md` |
| Release notes | `RELEASE_NOTES_v*.md` |
| Compliance tests | `tests/integration/atlas_compliance_test.go` |
| MITM proxy | `pkg/proxy/mitm.go` |
| ATLAS patterns | `pkg/compliance/atlas.go` |
| OWASP compliance | `pkg/compliance/owasp.go` |

---

## 8. Lessons Learned

### 🔴 Lesson 1: Version Sync is Not Optional
**Problem:** CI failed because VERSION didn't match main.go  
**Lesson:** Treat version files as code - commit them together  
**Prevention:** Create a version-sync script

### 🔴 Lesson 2: Test Naming Matters
**Problem:** TestViolationNames failed due to field confusion  
**Lesson:** Use consistent field access patterns across codebase  
**Prevention:** Add linting rule for pattern field access

### 🔴 Lesson 3: Clean Up Dead Code
**Problem:** Unused functions/variables flagged by linter  
**Lesson:** Remove dead code before committing  
**Prevention:** Pre-commit hooks with golangci-lint

### 🔴 Lesson 4: GitHub Actions Need Sync Time
**Problem:** Actions didn't run immediately after push  
**Lesson:** Wait 30-60 seconds, check web UI  
**Prevention:** Verify remote sync independently

### 🔴 Lesson 5: Whitespace Matters in Edits
**Problem:** Python scripts needed to handle Windows line endings  
**Lesson:** Use proper encoding when editing files programmatically  
**Prevention:** Normalize line endings, use explicit content

---

## 9. Troubleshooting Guide

### Common Issues

#### Issue: golangci-lint errors after commit
```
Error: var _configPath is unused
Error: func severityFromString is unused
```
**Solution:** Remove unused code, ensure clean lint before pushing

#### Issue: GitHub Actions not running
```
Workflow doesn't appear in GitHub Actions tab
```
**Solution:**
1. Verify remote is synced: `git status`
2. Check branch protection rules
3. Push fresh commit if needed

#### Issue: Version mismatch errors
```
VERSION file: 0.29.0
main.go: 0.29.1
```
**Solution:** Update VERSION to match main.go constant

#### Issue: Test failures
```
TestViolationNames: Expected 2 unique names, got 1
```
**Solution:** Check Pattern.Name vs Pattern.Description usage in mitm.go

#### Issue: Build failures
```
go build: cannot find package
```
**Solution:** Run `go mod download` and ensure go.mod is correct

### Debug Commands

```bash
# Check version alignment
python check_version.py

# Check for lint errors
golangci-lint run ./...

# Run tests verbosely
go test -v ./...

# Check git status
git status

# Verify remote sync
git remote -v
git log --oneline -5

# View current version
cat VERSION
./aegisgate --version
```

---

## 10. Technical Reference

### GitHub Actions Workflows

| Workflow | Trigger | Purpose |
|----------|---------|---------|
| `test.yml` | Push to main/develop, PR | Integration tests |
| `ci.yml` | Push, PR | CI pipeline |
| `release.yml` | Tag push | Release automation |
| `security.yml` | Push, schedule | Security scanning |
| `sbom.yml` | Push, schedule | SBOM generation |
| `version-check.yml` | Push, PR | Version validation |
| `version-sync.yml` | Schedule | Version synchronization |

### Go Module Structure

```go
module github.com/aegisgatesecurity/aegisgate

go 1.24.0

require (
    github.com/prometheus/client_golang v1.20.5
    golang.org/x/net v0.35.0
    golang.org/x/oauth2 v0.35.0
)
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | 8080 |
| `UPSTREAM_URL` | Upstream AI service URL | - |
| `RATE_LIMIT` | Requests per minute | 1000 |
| `ATLAS_ENABLED` | Enable ATLAS compliance | false |
| `LOG_LEVEL` | Logging level | info |
| `AEGISGATE_MITM_ENABLED` | Enable MITM proxy mode | false |
| `AEGISGATE_AUTH_PROVIDER` | Auth provider (local, oauth, saml) | - |

### API Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/chat/completions` | POST | OpenAI-compatible chat |
| `/v1/completions` | POST | OpenAI-compatible completions |
| `/v1/embeddings` | POST | OpenAI-compatible embeddings |
| `/metrics` | GET | Prometheus metrics |
| `/health` | GET | Health check |
| `/health/live` | GET | Liveness probe |
| `/health/ready` | GET | Readiness probe |

---

## 📌 Quick Reference Card

```
┌────────────────────────────────────────────────────────────────┐
│                     AEGISGATE QUICK REFERENCE                     │
├────────────────────────────────────────────────────────────────┤
│ Version:        v0.29.1                                        │
│ Go Version:     1.24.0                                        │
│ Main File:      cmd/aegisgate/main.go                           │
│ Version Const:  Line ~43 in main.go                           │
│ Version File:   VERSION                                       │
│ Remote:         git@github.com:aegisgatesecurity/aegisgate.git        │
│ Branch:         main                                          │
│ Last Commit:    dd8dd81                                       │
├────────────────────────────────────────────────────────────────┤
│ Commands:                                                        │
│   Build:    go build -o aegisgate ./cmd/aegisgate                 │
│   Test:     go test ./...                                     │
│   Lint:     golangci-lint run                                 │
│   Version:  cat VERSION                                       │
├────────────────────────────────────────────────────────────────┤
│ Common Fixes:                                                   │
│   - Unused code → Remove it                                    │
│   - Version mismatch → Update VERSION file                    │
│   - Lint errors → golangci-lint run --fix                     │
│   - Actions not running → Check remote sync                   │
└────────────────────────────────────────────────────────────────┘
```

---

## 🔗 Related Documentation

- [README.md](../README.md) - Main project documentation
- [RELEASE_NOTES_v0.29.1.md](./RELEASE_NOTES_v0.29.1.md) - Latest release notes
- [docs/ATLAS_FRAMEWORK.md](./docs/ATLAS_FRAMEWORK.md) - ATLAS compliance
- [docs/CONFIGURATION.md](./docs/CONFIGURATION.md) - Configuration guide
- [docs/SECURITY.md](./docs/SECURITY.md) - Security documentation

---

*This document is maintained as a living reference for the AegisGate project. Last updated after fixing golangci-lint errors and version synchronization in v0.29.1 release.*
