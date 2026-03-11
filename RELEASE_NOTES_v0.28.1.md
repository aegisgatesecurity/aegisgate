# AegisGate v0.28.1 Release Notes

## Release Date: March 4, 2026

## Overview

This release introduces comprehensive observability features with Prometheus metrics integration, enhanced health endpoints, and improved production-readiness for Kubernetes deployments.

## Whats New

### Prometheus Metrics Integration

The platform now includes a full Prometheus metrics export system at /metrics:

- aegisgate_requests_total: Total requests processed
- aegisgate_blocked_requests_total: Blocked requests by reason
- aegisgate_violations_total: Security violations by severity
- aegisgate_request_duration_seconds: Request duration
- aegisgate_goroutines: Number of goroutines
- aegisgate_memory_alloc_bytes: Memory allocated
- aegisgate_component_health: Component health status

### Enhanced Health Endpoints

New Kubernetes-ready health endpoints:

| Endpoint | Purpose |
|----------|---------|
| /health | Enhanced health with system metrics |
| /health/components | Detailed component status |

## Files Changed

| File | Status |
|------|--------|
| pkg/dashboard/observability.go | NEW |
| pkg/dashboard/dashboard.go | MODIFIED |
| go.mod | MODIFIED |
| README.md | MODIFIED |

## Dependencies

- github.com/prometheus/client_golang v1.20.5

## Breaking Changes

None - backward compatible.

## Upgrade Guide

1. git pull origin main
2. go mod tidy
3. go build -o aegisgate ./cmd/aegisgate

## Next Release (v0.29.0)

- Advanced ML-based threat detection
- Enhanced SIEM integrations
- Kubernetes operator

---

Full Changelog: https://github.com/aegisgatesecurity/aegisgate/compare/v0.28.0...v0.28.1
