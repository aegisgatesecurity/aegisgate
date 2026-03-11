## AegisGate v0.3.0 - CI Pipeline Fixes

### Changes in v0.3.0
- Fixed CI pipeline compilation errors
- Removed conflicting API wrapper files
- Fixed integration tests to use actual implementation APIs
- Fixed WebSocket example tests with expected output
- All workflows passing (Build, Test, Security, SBOM)

### Full Changelog
- fix: remove unused net/http import in websocket example_test.go
- fix: websocket example tests to produce expected output  
- fix: unused imports in integration test and websocket examples
- fix: remove conflicting api wrappers and fix integration tests
- fix: update ml API wrapper to match test expectations

### What's New
This release focuses on build stability and CI/CD pipeline improvements.

### Installation
```bash
# Download release
curl -L https://github.com/aegisgatesecurity/aegisgate/archive/refs/tags/v0.3.0.tar.gz | tar xz
cd aegisgate-0.3.0

# Build
go build -o aegisgate ./cmd/aegisgate

# Run
./aegisgate
```

### Docker
```bash
docker pull aegisgatesecurity/aegisgate:v0.3.0
```

### Checksums
See release assets for checksums.

**Full Changelog**: https://github.com/aegisgatesecurity/aegisgate/commits/v0.3.0