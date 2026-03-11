# Release Notes - v0.18.1

**Release Date**: February 23, 2026

## Overview

AegisGate v0.18.1 is a maintenance release that resolves golangci-lint errors identified in the Phase 2 Enterprise packages after the v0.18.0 release.

## 🐛 Bug Fixes

### Linter Fixes

| Issue | File | Fix |
|-------|------|-----|
| Unused field `raw` | pkg/threatintel/types.go:652 | Removed unused field from Observable struct |
| Unnecessary fmt.Sprintf | pkg/threatintel/stix.go:834 | Replaced with string literal |
| XML tag conflict | pkg/sso/saml.go:691 | Fixed `saml:NameID` to `NameID` |

### Details

#### pkg/threatintel/types.go
- Removed unused `raw map[string]interface{}` field from Observable struct
- Field was intended for future use but caused linting errors

#### pkg/threatintel/stix.go
- Changed `fmt.Sprintf("File hash indicator")` to `"File hash indicator"`
- Simplified string literal per gosimple S1039

#### pkg/sso/saml.go
- Fixed XML tag conflict in LogoutRequest struct
- Changed `xml:"saml:NameID"` to `xml:"NameID"` to avoid SA5008 staticcheck error
- Resolved namespace conflict between struct XMLName field and NameID field

## ✅ Verification

All tests passing:
- `go test ./pkg/threatintel/...` ✅
- `go test ./pkg/sso/...` ✅
- `go test ./pkg/webhook/...` ✅
- `go vet ./...` ✅
- `go build ./...` ✅

## 📦 Packages Affected

| Package | Changes |
|---------|---------|
| pkg/threatintel | types.go, stix.go |
| pkg/sso | saml.go |

## 🔄 Upgrade Guide

### From v0.18.0 to v0.18.1

```bash
# Pull the latest changes
git pull origin main
git checkout v0.18.1

# Update dependencies (no changes)
go mod tidy
go mod download

# Build
go build -o aegisgate ./cmd/aegisgate
```

No configuration changes required for this release.

## 📚 Documentation

- [README.md](README.md) - Comprehensive project documentation
- [CHANGELOG.md](CHANGELOG.md) - Full release history
- [docs/SIEM_INTEGRATION.md](docs/SIEM_INTEGRATION.md) - SIEM setup guide

## 📥 Download

- **Source Code**: [GitHub Releases](https://github.com/aegisgatesecurity/aegisgate/releases/tag/v0.18.1)
- **Docker Image**: `ghcr.io/aegisgatesecurity/aegisgate:v0.18.1`
- **Go Module**: `github.com/aegisgatesecurity/aegisgate v0.18.1`

---

**Previous Release**: [v0.18.0](RELEASE_NOTES_v0.18.0.md) - Phase 2 Enterprise Readiness
