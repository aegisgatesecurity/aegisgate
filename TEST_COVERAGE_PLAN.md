# AegisGate Test Coverage Plan

## Target: 90%+ Code Coverage on All Core Packages

### Current Coverage Analysis

| Package | Current Status | Target | Gap |
|---------|---------------|--------|-----|
| proxy/ | ~70% | 90% | +20% |
| compliance/ | ~50% | 90% | +40% |
| auth/ | ~60% | 90% | +30% |
| tls/ | ~40% | 90% | +50% |
| config/ | ~70% | 90% | +20% |
| sso/ | sso_test.go | 90% | +40% |
| ml/ | detector_test.go | 90% | +30% |
| opsec/ | opsec_test.go | 90% | +30% |
| dashboard/ | dashboard_test.go | 90% | +20% |
| siem/ | siem_test.go | 90% | +20% |
| secrets/ | 2 test files | 90% | +30% |
| security/ | 2 test files | 90% | +30% |
| metrics/ | metrics_test.go | 90% | +20% |
| hash_chain/ | test.go | 90% | +50% |
| signature_verification/ | example_test.go | 90% | +50% |
| trustdomain/ | manager_test.go, validation_test.go | 90% | +20% |
| sandbox/ | sandbox_test.go | 90% | +30% |
| webhook/ | webhook_test.go | 90% | OK |
| websocket/ | websocket_test.go | 90% | OK |

## Coverage Gaps to Address

### Priority 1 - Critical (Low Coverage <50%)

1. **tls/tls_test.go** - Currently 1 line! Needs complete rewrite
2. **hash_chain/** - Only has test.go (10 lines)
3. **signature_verification/** - Only has example_test.go

### Priority 2 - Important (Medium Coverage 50-70%)

4. **compliance/** - Needs more edge case tests
5. **sso/** - Needs more OIDC/SAML scenario tests
6. **ml/** - Needs more detection scenario tests
7. **secrets/** - Provider needs more test scenarios

### Priority 3 - Enhancement (70-90%)

8. **proxy/** - Add more edge cases
9. **auth/** - Add more auth flow tests
10. **security/** - Add more attack vector tests

## Test Coverage Strategy

### 1. Table-Driven Tests
Use Go table-driven tests for maximum coverage with minimal code:
```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        expected string
        wantErr bool
    }{...}
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {...})
    }
}
```

### 2. Property-Based Testing
For cryptographic and complex logic:
```go
func TestProperty(t *testing.T) {
    for i := 0; i < 1000; i++ {
        // Generate random input
        // Test property holds
    }
}
```

### 3. Fuzz Testing
Use go-fuzz for complex parsing:
```go
func FuzzParse(f *testing.F) {
    f.Fuzz(func(t *testing.T, s string) {...})
}
```

### 4. Integration Tests
Test real-world scenarios in tests/integration/

## Execution Plan

### Phase 1: Fix Empty/Near-Empty Test Files
- [ ] tls/tls_test.go - Complete rewrite
- [ ] hash_chain/test.go - Expand to comprehensive
- [ ] signature_verification/ - Add full test suite

### Phase 2: Expand Existing Tests
- [ ] compliance/compliance_test.go - Add 20+ cases
- [ ] sso/sso_test.go - Add OIDC/SAML flows
- [ ] ml/detector_test.go - Add anomaly scenarios
- [ ] opsec/opsec_test.go - Add security scenarios

### Phase 3: Coverage Optimization
- [ ] Run `go test -cover` to identify gaps
- [ ] Add tests for uncovered branches
- [ ] Add tests for error conditions
- [ ] Add tests for edge cases

### Phase 4: Verification
- [ ] Run full test suite with coverage
- [ ] Generate coverage report
- [ ] Document coverage metrics

## Running Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage by package
go test -coverprofile=coverage.out ./pkg/... -covermode=atomic

# HTML report
go tool cover -html=coverage.out -o coverage.html

# Per-package coverage
go test -cover ./pkg/tls/...
go test -cover ./pkg/proxy/...
```
