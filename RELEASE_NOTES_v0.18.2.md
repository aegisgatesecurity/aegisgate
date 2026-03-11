# Release Notes v0.18.2

**Release Date:** 2026-02-23

## Bug Fixes

| Issue | File | Fix |
|-------|------|-----|
| CI workflow failure | `.github/workflows/ci.yml` | Set `CGO_ENABLED=0` |
| Test workflow failure | `.github/workflows/test.yml` | Set `CGO_ENABLED=0` |
| Build workflow failure | `.github/workflows/build.yml` | Set `CGO_ENABLED=0` |

## Problem

The GitHub Actions workflows were failing on the Windows self-hosted runner with the following error:

```
cgo: C compiler "gcc" not found: exec: "gcc": executable file not found in %PATH%
```

This was caused by `CGO_ENABLED=1` being set in the ci.yml workflow, requiring a C compiler that was not installed on the Windows runner.

## Solution

Since AegisGate is a pure Go project with no CGO dependencies, setting `CGO_ENABLED=0` in all workflow files resolves the issue:

```yaml
- name: Run tests
  run: go test -v ./...
  env:
    CGO_ENABLED: "0"
```

## Verification

After this fix, all CI workflows should pass successfully:
- Test workflow
- Lint workflow  
- Build workflow
- Security scan

## Upgrade Guide

No code changes required - this is a CI workflow fix only. Existing v0.18.0/v0.18.1 installations are unaffected.

---

**Previous Release:** [v0.18.1](./RELEASE_NOTES_v0.18.1.md)
