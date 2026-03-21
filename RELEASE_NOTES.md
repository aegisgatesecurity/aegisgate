# AegisGate v1.0.12 Release Notes

## Release Date
March 18, 2026

## Version
v1.0.12

## Overview
This release introduces the **License Validation Middleware**, enabling AegisGate to integrate with the AegisGate Admin Panel for centralized license management and validation.

---

## New Public Features

### 1. License Validation Middleware

**File:** `pkg/middleware/license.go`

The license validation middleware provides enterprise-grade license checking by integrating with the AegisGate Admin Panel API.

#### Features
- Remote License Validation against Admin Panel
- Automatic caching (5-minute default)
- Fail-Open mode when license service unavailable
- Context integration for tier and rate limits

#### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| ADMIN_PANEL_URL | Admin Panel API URL | http://localhost:8443 |
| LICENSE_KEY | License key to validate | (none - Community tier) |
| LICENSE_PUBLIC_KEY | Public key for verification | (none) |

#### Usage
```go
import "github.com/aegisgatesecurity/aegisgate/pkg/middleware"

handler := middleware.LicenseMiddleware()(yourHandler)
```

### 2. Tier Configuration File

**File:** `config/tier.conf`

Defines tier limits for offline/fallback operation.

### 3. Four-Tier Licensing System

| Tier | Max Servers | Max Users | Rate Limit |
|------|-------------|-----------|------------|
| Community | 1 | 3 | 60/min |
| Developer | 5 | 10 | 600/min |
| Professional | 25 | 50 | 3000/min |
| Enterprise | Unlimited | Unlimited | Unlimited |

*Enterprise pricing is managed through the AegisGate Admin Panel.*

---

## Migration

Existing deployments continue to work. To enable license validation:

```bash
export LICENSE_KEY="AG-xxxx-xxxxxxxxxxxx"
export ADMIN_PANEL_URL="https://admin.yourcompany.com"
```

---

## Links
- GitHub: https://github.com/aegisgatesecurity/aegisgate/releases/tag/v1.0.12
- Docs: https://docs.aegisgate.security