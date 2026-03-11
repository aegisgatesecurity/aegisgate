# Release Notes v1.0.1

**Release Date:** 2026-03-11

---

## Overview

This release includes bug fixes, workflow improvements, and repository cleanup. The focus was on stability improvements and preparing the community repository for production use.

---

## Changes

### Bug Fixes

#### Context Keys Fix (`pkg/api/versioning.go`)
- **Issue:** Context keys `apiVersionKey` and `versionInfoKey` were undefined
- **Fix:** Properly defined context key types in the versioning package

#### ShouldBlock Logic Fix (`pkg/scanner/patterns.go`)
- **Issue:** Test failure in ShouldBlock function - severity threshold was set too high
- **Fix:** Changed severity threshold from `>= Critical` to `>= High` to properly detect and block threats

#### Dockerfile Configuration Fix (`Dockerfile`)
- **Issue:** Docker build failed - referenced incorrect config file `padlock.yml.example`
- **Fix:** Updated to reference correct `aegisgate.yml.example` configuration file

### Repository Cleanup

- **Removed:** 150+ temporary files including:
  - Check scripts (check_*.py, check_*.ps1)
  - Fix scripts (fix_*.py, fix_*.ps1) 
  - Commit helper scripts (commit*.bat, commit*.py)
  - Test binaries (*.test.exe)
  - Certificate files (*.crt, *.key)
  - Old release notes (v0.x.x)
  - Analysis and memory documentation
  - Configuration backup files
- **Result:** Clean repository with 46 essential files

### CI/CD Improvements

- Fixed GitHub Actions workflow errors
- Improved Docker build reliability
- Streamlined configuration handling

---

## Upgrading from v1.0.0

No special migration steps required. Simply update to the latest version:

```bash
# Rebuild the binary
make build

# Or pull the latest Docker image
docker pull aegisgatesecurity/aegisgate:latest
```

---

## Known Issues

No critical known issues in this release.

---

## What's Next

Expected in v1.0.2:
- Enhanced ML detection algorithms
- Additional compliance framework support
- Performance optimizations
- Kubernetes deployment manifests

---

## Contributors

Thanks to all contributors who helped identify and fix these issues.

---

## Security

For security vulnerabilities, please report to: **security@aegisgatesecurity.io**

See SECURITY.md for full disclosure guidelines and supported versions.

---

## Links

- **Documentation:** https://docs.aegisgatesecurity.io
- **GitHub Issues:** https://github.com/aegisgatesecurity/aegisgate/issues
- **Website:** https://aegisgatesecurity.io

---

**AegisGate** - Protect your AI infrastructure