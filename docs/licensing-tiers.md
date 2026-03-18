# Four-Tier Licensing System

AegisGate implements a comprehensive four-tier licensing model designed to scale with your organization from free testing to enterprise deployment.

## Overview

The licensing system differentiates access based on organizational needs:

| Tier | Level | Price | Target Audience |
|------|-------|-------|-----------------|
| Community | 0 | Free | Testing, Personal |
| Developer | 1 | See pricing page | Individual Developers |
| Professional | 2 | See pricing page | Teams, Businesses |
| Enterprise | 3 | Custom | Large Organizations |

## Tier Details

### Tier 0: Community

The Community tier is designed for testing, learning, and personal projects.

- **Price**: Free
- **Level Value**: `0`
- **Max Servers**: 1
- **Max Users**: 3
- **Rate Limit**: 60 requests/minute

**Features**:
- Basic proxy functionality
- Core modules
- Community support

**Use Cases**:
- Learning AegisGate
- Personal projects
- Testing in development environments

**Example**:
```bash
# No license key needed - Community tier is default
./aegisgate --tier community
```

### Tier 1: Developer

The Developer tier is designed for individual developers and small projects.

- **Price**: See pricing page
- **Level Value**: `1`
- **Max Servers**: 5
- **Max Users**: 10
- **Rate Limit**: 600 requests/minute

**Features**:
- All Community features
- API Access
- Email Support
- Custom Modules

**Use Cases**:
- Individual developers
- Small startup projects
- Production apps with moderate traffic

**Example**:
```bash
# Set license key via environment
export LICENSE_KEY="AG-xxxxxxxxxxxxx"
./aegisgate
```

### Tier 2: Professional

The Professional tier is designed for professional teams and businesses.

- **Price**: See pricing page
- **Level Value**: `2`
- **Max Servers**: 25
- **Max Users**: 50
- **Rate Limit**: 3000 requests/minute

**Features**:
- All Developer features
- Priority Support
- Advanced Analytics
- Load Balancing

**Use Cases**:
- Professional teams
- Mid-size businesses
- High-traffic applications

### Tier 3: Enterprise

The Enterprise tier is designed for large organizations with custom needs.

- **Price**: Contact sales
- **Level Value**: `3`
- **Max Servers**: Unlimited
- **Max Users**: Unlimited
- **Rate Limit**: Unlimited

**Features**:
- All Professional features
- 24/7 Dedicated Support
- Custom Integrations
- SLA Guarantee
- On-premise Deployment
- Custom development

**Use Cases**:
- Large organizations
- Mission-critical applications
- Compliance-heavy industries

## Programmatic Usage

### Setting Tier in Code

```go
import "github.com/aegisgatesecurity/aegisgate/pkg/core"

// Create license validator
validator := middleware.NewLicenseValidator(nil)

// Get the tier
tier := validator.GetTier(ctx)

// Check tier level
if tier >= core.TierProfessional {
    // Enable advanced features
}
```

### Tier Comparison

```go
func checkTierAccess(tier core.Tier, requiredTier core.Tier) bool {
    return tier >= requiredTier
}

// Usage
isProOrHigher := checkTierAccess(userTier, core.TierProfessional)
```

### Available Tier Values

```go
const (
    TierCommunity    Tier = iota  // 0
    TierDeveloper                 // 1
    TierProfessional              // 2
    TierEnterprise                // 3
)
```

## Getting a License

1. **Community**: No registration required
2. **Developer/Professional**: Sign up at [aegisgate.security](https://aegisgate.security)
3. **Enterprise**: Contact sales at [sales@aegisgate.security](mailto:sales@aegisgate.security)

## License Key Format

License keys follow the format: `AG-<ID>-<SIGNATURE>`

Example: `AG-a1b2c3d4-abcdefghijklmnopqrstuvwxy`

## Upgrading

To upgrade your tier:

1. Log into the Admin Panel
2. Navigate to "Licenses"
3. Select "Create License"
4. Choose your desired tier
5. Complete payment (for paid tiers)

## Feature Matrix

| Feature | Community | Developer | Professional | Enterprise |
|---------|-----------|-----------|--------------|------------|
| Basic Proxy | ✅ | ✅ | ✅ | ✅ |
| Core Modules | ✅ | ✅ | ✅ | ✅ |
| API Access | ❌ | ✅ | ✅ | ✅ |
| Custom Modules | ❌ | ✅ | ✅ | ✅ |
| Priority Support | ❌ | ❌ | ✅ | ✅ |
| Advanced Analytics | ❌ | ❌ | ✅ | ✅ |
| Load Balancing | ❌ | ❌ | ✅ | ✅ |
| 24/7 Support | ❌ | ❌ | ❌ | ✅ |
| Custom Integrations | ❌ | ❌ | ❌ | ✅ |
| SLA Guarantee | ❌ | ❌ | ❌ | ✅ |
| On-premise | ❌ | ❌ | ❌ | ✅ |