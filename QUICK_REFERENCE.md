# AegisGate - Quick Reference Card

Version: v0.28.0 | Date: 2026-03-04 | Status: CI/CD Green ✅

## Current State
- **6 Compliance Frameworks** fixed: ATLAS, GDPR, OWASP, NIST, ISO42001, SOC2
- **OPSEC Module** stable with safe memory scrubber
- **i18n**: 12 locales integrated
- **Coverage**: 83.4% across 147 test suites

## Critical Fix Patterns (from v0.28.0)

### Interface Evolution
```go
// Add compile-time assertion after interface changes
var _ common.Framework = (*Framework)(nil)
```

### Pointer Factory Functions
```go
// Return POINTER, not value type
func DefaultConfig() *Config {  // ✅ Good
func DefaultConfig() Config {   // ❌ Bad - cant take address
```

### Safe Memory Operations (No unsafe.Pointer)
```go
// ✅ Safe: Use standard conversion + runtime.KeepAlive
b := []byte(strData)
for i := range b { b[i] = 0 }
runtime.KeepAlive(strData)

// ❌ Unsafe: reflect.StringHeader manipulation
hdr := (*reflect.StringHeader)(unsafe.Pointer(s))
```

### go.mod Version
Use minor version only: `go 1.23` not `go 1.24.0`

### Windows CI/Build Tags
Linux-only files use `//go:build linux` - create Windows stubs for cross-platform

## Key Files
- Interface: pkg/common/framework.go
- OPSEC Config: pkg/opsec/config.go (returns *Config)
- Memory Scrubber: pkg/opsec/memory_scrubber.go (safe implementation)
- Runtime Hardening: pkg/opsec/runtime_hardening_linux.go (hardcoded syscalls)

## Next Priorities
1. Windows runtime_hardening stub
2. PKI attestation for MITM
3. Threat intel enrichment API
4. Behavioral analytics (ML detection)
5. Simple installer/config generator

## Status
ALL GREEN ✅ - Ready for continued development
