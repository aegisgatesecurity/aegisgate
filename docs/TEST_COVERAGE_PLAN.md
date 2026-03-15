# Test Coverage Plan - Priority 3

## Current Coverage by Package

### High Coverage (>60%) - Good ✅
| Package | Coverage |
|---------|----------|
| pkg/immutable-config/readonly | 98.6% |
| pkg/immutable-config/rollback | 96.7% |
| pkg/immutable-config/logging | 94.1% |
| pkg/immutable-config/integrity | 90.0% |
| pkg/immutable-config/watcher | 88.5% |
| pkg/hash_chain | 86.6% |
| pkg/immutable-config/wal | 87.2% |
| pkg/immutable-config/snapshot | 78.3% |
| pkg/secrets | 77.7% |
| pkg/crypto/fips | 76.3% |
| pkg/certificate | 75.9% |
| pkg/sandbox | 71.7% |
| pkg/adapters | 69.4% |
| pkg/plugin | 67.4% |
| pkg/metrics | 65.1% |
| pkg/i18n | 59.3% |
| pkg/immutable-config | 58.2% |
| pkg/scanner | 58.4% |
| pkg/websocket | 55.2% |
| pkg/ml | 53.9% |
| pkg/immutable-config/filesystem | 52.1% |
| pkg/webhook | 49.0% |
| pkg/core | 46.8% |
| pkg/dashboard | 38.8% |
| pkg/config | 37.6% |
| pkg/api | 36.6% |
| pkg/tenant | 36.7% |
| pkg/threatintel | 35.9% |
| pkg/signature_verification | 35.8% |
| pkg/trustdomain | 32.4% |

### Medium Coverage (20-35%) - Needs Improvement
| Package | Coverage |
|---------|----------|
| pkg/auth | 79.7% |
| pkg/proxy | 27.9% |
| pkg/compliance | 22.1% |
| pkg/sso | 21.5% |
| pkg/siem | 24.2% |
| pkg/tls | 14.8% |
| pkg/pkiattest | 11.2% |

### Low Coverage (<15%) - Needs Tests
| Package | Coverage |
|---------|----------|
| pkg/security | 9.6% |
| tests/integration | 8.8% |

## Priority Improvements

### 1. pkg/security - Target 40%
- Add tests for: panic_recovery.go, recovery.go
- Current: 9.6% → Target: 40%

### 2. pkg/compliance - Target 40%  
- Add tests for: atlas.go, compliance.go, framework_mapping.go
- Current: 22.1% → Target: 40%

### 3. pkg/proxy - Target 40%
- Add tests for: proxy.go, http3.go, mitm.go
- Current: 27.9% → Target: 40%

### 4. pkg/tls - Target 30%
- Add tests for: ca.go, manager.go
- Current: 14.8% → Target: 30%

## Implementation

Run coverage with:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Benchmarks Already Exist ✅
- pkg/proxy/proxy_benchmark_test.go
- pkg/scanner/scanner_benchmark_test.go  
- pkg/security/middleware_bench_test.go

Run benchmarks:
```bash
go test -bench=. -benchmem ./pkg/proxy/...
go test -bench=. -benchmem ./pkg/scanner/...
```
