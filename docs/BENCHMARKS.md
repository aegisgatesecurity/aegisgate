# AegisGate Performance Benchmarks

## Overview
This document describes the performance benchmarks for AegisGate AI Security Gateway.

## Running Benchmarks

### All benchmarks
```bash
go test -bench=. ./pkg/...
```

### Specific package benchmarks
```bash
# Proxy benchmarks
go test -bench=. ./pkg/proxy/

# Scanner benchmarks
go test -bench=. ./pkg/scanner/

# Compliance benchmarks
go test -bench=. ./pkg/compliance/
```

### With memory profiling
```bash
go test -bench=. -benchmem ./pkg/...
```

### With CPU profiling
```bash
go test -bench=. -cpuprofile=cpu.prof ./pkg/...
go tool pprof cpu.prof
```

## Performance Targets (v1.0)

| Metric | Target | Current |
|--------|--------|---------|
| Request Latency (p50) | < 10ms | TBD |
| Request Latency (p95) | < 50ms | TBD |
| Request Latency (p99) | < 100ms | TBD |
| Throughput | > 10,000 req/s | TBD |
| Memory per Request | < 1KB | TBD |
| Scanner Latency | < 5ms | TBD |
| Compliance Report Gen | < 500ms | TBD |

## Benchmark Categories

### 1. Proxy Benchmarks
- **BenchmarkRequestForwarding**: Basic request forwarding latency
- **BenchmarkScanAndForward**: Request scanning + forwarding
- **BenchmarkSecurityScan**: Security scanning isolation
- **BenchmarkLatency**: End-to-end latency percentiles
- **BenchmarkThroughput**: Requests per second
- **BenchmarkMemoryAllocation**: Memory allocation per request

### 2. Scanner Benchmarks
- **BenchmarkScanRequest**: Single request scanning
- **BenchmarkScanWithPatterns**: Pattern matching performance
- **BenchmarkScanLargePayload**: Large payload handling (1MB)
- **BenchmarkScanMaliciousPayloads**: Known malicious payload detection
- **BenchmarkRegexMatching**: Regex pattern matching

### 3. Compliance Benchmarks
- **BenchmarkGenerateMITREReport**: MITRE ATLAS report generation
- **BenchmarkGenerateNISTReport**: NIST AI RMF report generation
- **BenchmarkGenerateISO42001Report**: ISO 42001 report generation
- **BenchmarkValidateCompliance**: Compliance rule validation
- **BenchmarkCalculateRiskScore**: Risk score calculation
- **BenchmarkExportJSON/PDF**: Report export performance

## Performance Monitoring

### Enabling pprof
```bash
export PPROF_ENABLED=true
./aegisgate
```

### Accessing pprof
- CPU Profile: http://localhost:6060/debug/pprof/profile
- Heap Profile: http://localhost:6060/debug/pprof/heap
- Goroutine Profile: http://localhost:6060/debug/pprof/goroutine
- Full Profile: http://localhost:6060/debug/pprof/

### Memory Analysis
```bash
curl http://localhost:6060/debug/pprof/heap > heap.prof
go tool pprof heap.prof
```

## Continuous Performance Testing

Benchmarks are run automatically in CI for:
- Every pull request
- Every release candidate
- Weekly scheduled runs

Performance regressions > 10% will block the PR.

## Hardware Specifications for Baseline

- CPU: 8 cores @ 2.5GHz
- RAM: 16GB
- Network: 10Gbps
- OS: Ubuntu 22.04 LTS
- Go: 1.23+

Last Updated: 2024-02-23
