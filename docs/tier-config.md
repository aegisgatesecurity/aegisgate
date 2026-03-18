# Tier Configuration

AegisGate uses a four-tier licensing system that defines resource limits for different license levels.

## Overview

The tier system provides scalable access control:
- **Community**: Free tier for testing/personal use
- **Developer**: Individual developers (see pricing)
- **Professional**: Teams and businesses (see pricing)
- **Enterprise**: Large organizations (custom pricing)

## Configuration File

The tier configuration is stored in `config/tier.conf`.

### File Location
```
AegisGate/
└── config/
    └── tier.conf
```

### File Format

```ini
[tiers]
# Tier 0: Community (Free)
community = { max_servers = 1, max_users = 3, rate_limit = 60 }

# Tier 1: Developer
developer = { max_servers = 5, max_users = 10, rate_limit = 600 }

# Tier 2: Professional
professional = { max_servers = 25, max_users = 50, rate_limit = 3000 }

# Tier 3: Enterprise (Custom)
enterprise = { max_servers = -1, max_users = -1, rate_limit = -1 }

[validation]
admin_panel_url = "http://localhost:8443"
fail_open = true
```

## Tier Definitions

### Community (Tier 0)
- **Price**: Free (no license required)
- **Max Servers**: 1
- **Max Users**: 3
- **Rate Limit**: 60 requests/minute
- **Use Case**: Testing, personal projects, learning

### Developer (Tier 1)
- **Price**: See pricing page
- **Max Servers**: 5
- **Max Users**: 10
- **Rate Limit**: 600 requests/minute
- **Use Case**: Individual developers, small projects

### Professional (Tier 2)
- **Price**: See pricing page
- **Max Servers**: 25
- **Max Users**: 50
- **Rate Limit**: 3000 requests/minute
- **Use Case**: Professional teams, businesses

### Enterprise (Tier 3)
- **Price**: Custom (contact sales)
- **Max Servers**: Unlimited (-1)
- **Max Users**: Unlimited (-1)
- **Rate Limit**: Unlimited (-1)
- **Use Case**: Large organizations

## Configuration Values

### Server Limits

| Tier | Value | Meaning |
|------|-------|---------|
| Community | 1 | Single server only |
| Developer | 5 | Up to 5 servers |
| Professional | 25 | Up to 25 servers |
| Enterprise | -1 | Unlimited |

### User Limits

| Tier | Value | Meaning |
|------|-------|---------|
| Community | 3 | 3 users max |
| Developer | 10 | 10 users max |
| Professional | 50 | 50 users max |
| Enterprise | -1 | Unlimited |

### Rate Limits

| Tier | Value | Meaning |
|------|-------|---------|
| Community | 60 | 60 requests/minute |
| Developer | 600 | 600 requests/minute |
| Professional | 3000 | 3000 requests/minute |
| Enterprise | -1 | Unlimited |

## Validation Section

The `[validation]` section controls license validation behavior:

### admin_panel_url
URL of the AegisGate Admin Panel for license validation.
- Default: `http://localhost:8443`
- Example: `https://admin.yourcompany.com`

### fail_open
Determines behavior when the license service is unavailable.
- `true` (default): Allow requests (fail-open)
- `false`: Block requests when service is down

## Environment Variables

The configuration can also be controlled via environment variables:

```bash
# License key (required for non-Community tiers)
export LICENSE_KEY="AG-xxxx-xxxxxxxxxxxx"

# Admin Panel URL
export ADMIN_PANEL_URL="http://localhost:8443"

# Public key (optional)
export LICENSE_PUBLIC_KEY="-----BEGIN PUBLIC KEY-----..."
```

## Reading Configuration in Code

The tier configuration is used by the license middleware:

```go
import "github.com/aegisgatesecurity/aegisgate/pkg/core"

// Get tier from name
tier := core.GetTierByName("developer")

// Get tier display name
name := tier.String() // Returns "Developer"

// Get tier level
level := int(tier) // Returns 1
```

## Programmatic Access

You can define custom tier limits programmatically:

```go
limits := map[core.Tier]TierLimits{
    core.TierCommunity:    {MaxServers: 1, MaxUsers: 3, RateLimit: 60},
    core.TierDeveloper:    {MaxServers: 5, MaxUsers: 10, RateLimit: 600},
    core.TierProfessional: {MaxServers: 25, MaxUsers: 50, RateLimit: 3000},
    core.TierEnterprise:   {MaxServers: -1, MaxUsers: -1, RateLimit: -1},
}
```