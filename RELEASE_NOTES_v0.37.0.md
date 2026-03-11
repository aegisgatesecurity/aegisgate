# Release Notes - v0.37.0

**Release Date**: 2026-03-07  
**Type**: Patch / Maintenance Release  
**Status**: ✅ Stable

---

## 📝 Overview

v0.37.0 is a maintenance release focused on fixing build errors, test failures, and improving code quality across multiple packages. This release ensures the project builds cleanly and all tests pass.

---

## 🐛 Bug Fixes

### GraphQL Package (`pkg/graphql`)

- **Type Aliases**: Added type aliases for external package types (`sso.Provider`, `siem.Stats`, `auth.PasswordPolicy`)
- **Resolver Types**: Added `Pagination`, `UserFilter`, `ViolationFilter` and other resolver types
- **Executor Fixes**: Rewrote resolver call logic with proper error handling
- **Handler Simplification**: Simplified authentication middleware
- **Server Fixes**: Added executor field to server initialization
- **Import Fixes**: Fixed import issues in `types.go`
- **Subscription Fixes**: Fixed unused variable warnings

### gRPC Package (`pkg/grpc`)

- **TLS Types**: Added TLS types to `generated.go`
- **Compatibility Layers**: Created gRPC compatibility files:
  - `pkg/compliance/grpc_compat.go`
  - `pkg/proxy/grpc_compat.go`
  - `pkg/tls/grpc_compat.go`

### API Versioning (`pkg/api`)

- **Handler Registration**: Fixed `RegisterHandler` and `AddRoute` to not double-add "v" prefix
- **Version Negotiation**: Fixed version normalization (e.g., "1.0" → "v1") for:
  - Query parameters
  - Accept-Version headers
  - Content-Type negotiation
  - URL path parsing
- **Handler Lookup**: Added fallback logic for major-only version lookup
- **Deprecation**: Fixed `RegisterDeprecated()` to handle major-only version lookup
- **Unsupported Versions**: Fixed `RegisterUnsupported()` to handle major-only version lookup
- **Test Expectations**: Updated test expectations to match normalized version format
- **Syntax Errors**: Fixed orphan `})` and incomplete function definitions

### Configuration (`pkg/config`)

- **Lock Coping Warnings**: Fixed 3 go vet warnings by copying fields individually instead of copying entire struct:
  - `Config.Get()` method
  - `Config.Update()` method
  - `Manager.Get()` method

### Main Application (`cmd/aegisgate`)

- **API Mismatch**: Completely rewrote main.go to fix multiple API mismatches
- **Auth Manager**: Fixed `auth.NewManager` initialization
- **Proxy**: Fixed `proxy.New` initialization
- **GraphQL**: Fixed `graphql.NewResolver` and `graphql.NewHandler` initialization

### Test Files

- **config_test.go**: Simplified test file
- **sdk/go/aegisgate.go**: Fixed `fmt.Errorf` format string
- **api/versioning_test.go**: Fixed syntax errors in test functions

---

## ✅ Quality Improvements

### Code Quality

- **Go Vet**: All warnings resolved
- **Build**: Clean build with no errors
- **Tests**: All 65+ test packages passing

### Documentation

- **README**: Comprehensive documentation with architecture overview
- **Module Descriptions**: Added detailed module descriptions

---

## 🔄 Upgrading from v0.36.0

### Breaking Changes

None. This is a backward-compatible patch release.

### Migration Steps

1. Pull the latest changes:
   ```bash
   git pull origin main
   ```

2. Rebuild the application:
   ```bash
   go build -o aegisgate.exe ./cmd/aegisgate
   ```

3. Run tests to verify:
   ```bash
   go test ./...
   ```

---

## 📦 Included Changes

### Commits

- `4a1fe13` fix: Use valid URL in malformed URL test
- `aef5f1a` docs: Update README and release notes for v0.36.0
- `4acf9de` fix: Remove unused startPeriodicTasks function
- `f724149` fix: Resolve deadlock and unused function warnings
- `ffc6232` feat: Implement Plugin Architecture for v0.36.0
- *(plus additional commits from previous sessions)*

---

## 🧪 Testing

### Test Coverage

| Category | Status |
|----------|--------|
| Unit Tests | ✅ Passing |
| Integration Tests | ✅ Passing |
| Build Tests | ✅ Passing |
| Vet/Lint | ✅ Passing |

### Verified Packages

- `pkg/api` - All tests passing
- `pkg/auth` - All tests passing
- `pkg/config` - All tests passing
- `pkg/graphql` - All tests passing
- `pkg/grpc` - All tests passing
- `pkg/proxy` - All tests passing
- *(60+ packages total)*

---

## � Dependencies

### Go Version

- **Minimum**: Go 1.24
- **Recommended**: Go 1.24+

### Key Dependencies

- golang.org/x/net (HTTP/3, networking)
- golang.org/x/crypto (cryptography)
- gRPC ecosystem packages

---

## 🏆 Acknowledgments

Thanks to all contributors who helped identify and fix these issues.

---

## 📋 Previous Releases

- [v0.36.0](./RELEASE_NOTES_v0.36.0.md) - Plugin Architecture
- [v0.35.0](./RELEASE_NOTES_v0.35.0.md) - HTTP/3 Support
- [v0.34.0](./RELEASE_NOTES_v0.34.0.md) - SIEM IPv6 Fixes
- [v0.33.0](./RELEASE_NOTES_v0.33.0.md) - Hash Chain Integrity
- [v0.32.0](./RELEASE_NOTES_v0.32.0.md) - Development Updates

---

## ❓ Support

- **Issues**: Report bugs via [GitHub Issues](https://github.com/aegisgatesecurity/aegisgate/issues)
- **Discussions**: Ask questions via [GitHub Discussions](https://github.com/aegisgatesecurity/aegisgate/discussions)

---

<p align="center">
  <strong>AegisGate v0.37.0</strong><br>
  Secure • Comply • Protect
</p>
