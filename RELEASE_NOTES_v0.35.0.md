# AegisGate v0.35.0 Release Notes

**Release Date**: 2026-03-06  
**Version**: 0.35.0  
**Go Version**: 1.24.0

---

## 📋 Overview

v0.35.0 adds comprehensive HTTP/3 E2E test coverage and repository cleanup.

---

## ✅ Changes

### Repository Cleanup

| Change | Description |
|--------|-------------|
| .gitignore Update | Added development scripts to ignore list |
| Deleted Files | Removed 22 temporary/phase files |

### HTTP/3 E2E Test Coverage

| Test Category | Tests Added |
|---------------|-------------|
| Configuration Tests | 3 |
| TLS Configuration Tests | 2 |
| Server Lifecycle Tests | 3 |
| Request Handling Tests | 3 |
| Request Validation Tests | 5 |
| Metrics Tests | 2 |
| Enable/Disable Tests | 2 |
| Backend URL Tests | 2 |
| End-to-End Tests | 3 |
| Edge Case Tests | 3 |
| Integration Tests | 3 |
| **Total** | **35** |

### Files Changed

- `.gitignore` - Updated
- `PROJECT_MEMORY.md` - Updated
- `VERSION` - Updated to 0.35.0
- `cmd/aegisgate/main.go` - Version bump
- `pkg/proxy/http3_integration_test.go` - **NEW** (1014 lines)

---

## 🧪 Testing

All tests pass:

| Package | Tests | Status |
|---------|-------|--------|
| pkg/hash_chain | 21 | ✅ PASS |
| pkg/proxy | 55+ | ✅ PASS |
| pkg/scanner | 28 | ✅ PASS |

---

## 🚀 Next Steps (v0.36.0)

1. **Penetration Testing Preparation** - Begin internal security audit
2. **HTTP/3 E2E Tests** - Expand test coverage for HTTP/3
3. **Plugin Architecture Design** - Design extensibility interface

---

## 📦 Downloads

- **Binary**: `go install github.com/aegisgatesecurity/aegisgate/cmd/aegisgate@latest`
- **Docker**: `docker pull aegisgatesecurity/aegisgate:v0.35.0`
- **Source**: https://github.com/aegisgatesecurity/aegisgate/tree/v0.35.0

---

## 🔗 Links

- **Repository**: https://github.com/aegisgatesecurity/aegisgate
- **Releases**: https://github.com/aegisgatesecurity/aegisgate/releases
- **Issues**: https://github.com/aegisgatesecurity/aegisgate/issues

---

**🔒 AegisGate - Secure Your AI Infrastructure**
