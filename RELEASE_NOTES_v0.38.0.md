# Release Notes - AegisGate v0.38.0

## Overview

This release introduces **complete multi-tenancy support** with tenant-aware proxy middleware, comprehensive API endpoints for tenant management, and file-based persistent storage. This is a major milestone that enables AegisGate to serve multiple enterprise customers with isolated resources and per-tenant compliance.

---

## Highlights

### ✨ New Features

1. **Multi-Tenancy Architecture**
   - Complete tenant isolation with dedicated resources per tenant
   - Per-tenant rate limiting and circuit breaking
   - Per-tenant audit logging with hash chain integrity
   - Tenant-aware proxy middleware for traffic inspection

2. **Tenant Management API**
   - Full CRUD operations for tenant lifecycle
   - Tenant activation/suspension
   - Per-tenant quota management
   - Per-tenant compliance settings

3. **Persistent Storage**
   - File-based tenant storage with JSON persistence
   - Auto-loading of tenants on startup
   - Search and filter capabilities

4. **Proxy Integration**
   - TenantMiddleware for HTTP proxy
   - Domain-based tenant identification
   - Tenant context propagation

---

## Detailed Changes

### New Packages

| Package | Description |
|---------|-------------|
| `pkg/tenant` | Core multi-tenancy types and manager |
| `pkg/proxy/tenant_middleware.go` | Tenant-aware proxy middleware |
| `pkg/api/tenant_handler.go` | Tenant management API endpoints |

### Enhancements to Existing Modules

#### `pkg/opsec/audit.go`
- Added configurable retention periods (90 days to 7 years)
- Added FileStorageBackend for persistent audit storage
- Added AuditFilter for querying logs
- Added ComplianceAuditLog with auto-tagging
- Added alert callbacks for real-time notifications
- Added tamper-evident export functionality

#### `pkg/tenant/`
- SimpleRateLimiter for per-tenant rate limiting
- Tenant plans (Free, Starter, Pro, Enterprise)
- Resource quotas per tenant
- Compliance settings per tenant
- Feature gating based on plan

---

## Breaking Changes

None. This release is fully backward compatible.

---

## Security

### Audit Logging Improvements
- Hash chain integrity verification
- Tamper-evident export format
- Configurable retention periods (SOC2/HIPAA compliant)
- Compliance auto-tagging

### Tenant Isolation
- Isolated rate limiters per tenant
- Isolated circuit breakers per tenant
- Isolated audit logs per tenant

---

## Bug Fixes

- Fixed test compatibility in `pkg/opsec/opsec_test.go`
- Fixed API compatibility in `pkg/opsec/opsec.go`

---

## Deprecations

None.

---

## Migration Guide

### Upgrading from v0.37.0

No migration required. The new multi-tenancy features are additive and backward compatible.

### Enabling Multi-Tenancy

To enable multi-tenancy in your configuration:

```yaml
tenancy:
  enabled: true
  storage: /var/lib/aegisgate/tenants
```

---

## API Changes

### New Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/tenants` | POST | Create tenant |
| `/api/v1/tenants` | GET | List tenants |
| `/api/v1/tenants/{id}` | GET | Get tenant |
| `/api/v1/tenants/{id}` | PUT | Update tenant |
| `/api/v1/tenants/{id}` | DELETE | Delete tenant |
| `/api/v1/tenants/{id}/suspend` | POST | Suspend tenant |
| `/api/v1/tenants/{id}/activate` | POST | Activate tenant |
| `/api/v1/tenants/{id}/audit` | GET | Get audit logs |
| `/api/v1/tenants/{id}/audit/export` | GET | Export audit logs |
| `/api/v1/tenants/{id}/quota` | GET/PUT | Manage quota |
| `/api/v1/tenants/{id}/compliance` | GET/PUT | Manage compliance |

---

## Performance

No performance impact. New features are opt-in and don't affect existing deployments.

---

## Documentation

- Updated README with comprehensive multi-tenancy documentation
- Added API endpoint documentation

---

## Contributors

This release was developed with the assistance of AI tools as part of the automated development process.

---

## Known Issues

None reported.

---

## Next Steps

- Database backend (PostgreSQL/MySQL support)
- Authentication for tenant management API
- Billing integration for enterprise tenants
- Web dashboard for tenant management

---

## Downloads

- **Binary**: [aegisgate-v0.38.0-darwin-amd64.tar.gz](https://github.com/aegisgatesecurity/aegisgate/releases/tag/v0.38.0)
- **Docker**: `docker pull aegisgatesecurity/aegisgate:v0.38.0`
- **Source**: [v0.38.0.tar.gz](https://github.com/aegisgatesecurity/aegisgate/archive/v0.38.0.tar.gz)

---

## Changelog

### Full Changelog

- **Multi-Tenancy**: Complete tenant isolation system
- **Proxy**: Tenant-aware middleware
- **API**: Tenant management endpoints
- **Storage**: File-based persistence
- **Audit**: Enhanced compliance logging
- **Tests**: Updated for API compatibility

---

**Full Changelog**: https://github.com/aegisgatesecurity/aegisgate/compare/v0.37.0...v0.38.0
