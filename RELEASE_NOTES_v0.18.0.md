# Release Notes - v0.18.0

**Release Date**: February 23, 2026

## Overview

AegisGate v0.18.0 introduces Phase 2 Enterprise Readiness features, adding comprehensive threat intelligence sharing, webhook alerting, and enterprise single sign-on capabilities.

## 🆕 New Features

### Threat Intelligence (pkg/threatintel/)

STIX 2.1 and TAXII 2.1 threat intelligence integration for industry-wide threat sharing.

#### STIX 2.1 Support
- **Object Types**: Indicator, AttackPattern, ThreatActor, Malware, Vulnerability, Tool, Report, IntrusionSet, Campaign, CourseOfAction
- **Observables**: IPv4Address, IPv6Address, DomainName, URL, EmailAddress, MACAddress, File (with hashes)
- **Pattern Language**: STIX pattern generation for indicators
- **Relationships**: Object relationship management (indicates, uses, targets, etc.)
- **Bundle Export**: Full STIX bundle generation with all objects

#### TAXII 2.1 Protocol
- Full TAXII 2.1 client implementation
- Collection discovery and management
- Object retrieval with pagination
- Content-Range header support
- Multiple content types (STIX, JSON, CSV)

#### Export Formats
- **STIX 2.1**: Full STIX bundle export
- **JSON**: Structured threat intel data
- **CSV**: Spreadsheet-compatible export
- **MISP**: Direct MISP import format

#### Integration
- Convert SIEM events to STIX indicators
- MITRE ATT&CK framework mapping
- Automated threat intelligence generation

---

### Webhook Alerting (pkg/webhook/)

Real-time event notification system with enterprise-grade reliability.

#### Authentication Methods
| Method | Description |
|--------|-------------|
| **None** | No authentication (testing only) |
| **Basic** | HTTP Basic authentication |
| **Bearer** | Bearer token authentication |
| **API Key** | Header-based API key |
| **HMAC** | HMAC-SHA256 signature verification |
| **OAuth2** | OAuth2 client credentials flow |

#### Features
- Configurable retry with exponential backoff
- Event filtering by severity, category, and source
- Batch delivery support
- Concurrent webhook delivery with worker pools
- TLS/SSL support with certificate verification
- Delivery status tracking and history

#### Configuration Example
```yaml
webhooks:
  - name: "security-team"
    url: "https://hooks.example.com/security"
    auth_type: "hmac"
    secret: "${WEBHOOK_SECRET}"
    events: ["critical", "high"]
```

---

### Enterprise SSO (pkg/sso/)

Single Sign-On integration for enterprise deployments.

#### SAML 2.0 Features
- SP-initiated Single Sign-On
- IdP-initiated Single Sign-On
- Single Logout (SLO)
- Attribute-based role mapping
- Metadata URL auto-configuration
- Certificate-based signing and validation
- Session management with configurable duration

#### OIDC Features
- Authorization Code Flow (recommended)
- Implicit Flow (SPA support)
- Hybrid Flow
- PKCE support for public clients
- Automatic token refresh
- UserInfo endpoint integration
- Role extraction from claims

#### Configuration Example
```yaml
sso:
  provider: "saml"
  saml:
    entity_id: "https://aegisgate.example.com/saml/metadata"
    idp_metadata_url: "https://idp.example.com/metadata"
    attribute_mapping:
      email: "http://schemas.xmlsoap.org/ws/2005/05/identity/claims/emailaddress"
```

---

## 📦 Package Structure

### pkg/threatintel/
```
types.go       - STIX 2.1 and TAXII 2.1 core types (1400+ lines)
stix.go        - STIX builder and indicator generation (900+ lines)
taxii.go       - TAXII 2.1 client implementation
exporter.go    - Multi-format export manager
threatintel_test.go - Comprehensive test suite
```

### pkg/webhook/
```
types.go       - Webhook configuration and types
manager.go     - Webhook registration and delivery management
sender.go      - HTTP sender with authentication
filter.go      - Event filtering engine
webhook_test.go - Comprehensive test suite
```

### pkg/sso/
```
types.go       - SSO configuration types
saml.go        - SAML 2.0 provider implementation
oidc.go        - OpenID Connect provider implementation
manager.go     - SSO manager with session handling
middleware.go  - HTTP middleware for SSO
sso_test.go    - Comprehensive test suite
```

---

## 🔧 Changes

### Added
- pkg/threatintel/ - STIX 2.1 / TAXII 2.1 threat intelligence
- pkg/webhook/ - Webhook alerting system
- pkg/sso/ - Enterprise SSO (SAML/OIDC)
- golang.org/x/oauth2 dependency for OIDC support

### Changed
- Updated README.md with comprehensive v0.18.0 documentation
- Enhanced auth.Role with AtLeast() method for role hierarchy comparisons

---

## 📊 Statistics

| Metric | Value |
|--------|-------|
| New Packages | 3 |
| New Files | 15 |
| Lines Added | 15,368+ |
| Test Files | 3 |
| Supported SIEM Platforms | 10+ |
| STIX Object Types | 10 |
| SSO Protocols | 2 |

---

## 🔄 Upgrade Guide

### From v0.17.0 to v0.18.0

1. **Pull the latest changes**:
   ```bash
   git pull origin main
   git checkout v0.18.0
   ```

2. **Update dependencies**:
   ```bash
   go mod tidy
   go mod download
   ```

3. **Configure Threat Intelligence** (optional):
   ```yaml
   threat_intel:
     enabled: true
     taxii_server: "https://taxii.example.com"
     collection: "threat-intel"
   ```

4. **Configure Webhooks** (optional):
   ```yaml
   webhooks:
     - name: "security-team"
       url: "https://hooks.example.com/security"
       auth_type: "hmac"
       secret: "${WEBHOOK_SECRET}"
   ```

5. **Configure SSO** (optional):
   See documentation for SAML and OIDC configuration.

---

## 🐛 Known Issues

- Some edge cases in STIX pattern validation may require manual review
- Webhook OAuth2 token refresh is not yet automatic
- SSO session store defaults to memory (Redis recommended for production)

---

## 📚 Documentation

- [README.md](README.md) - Comprehensive project documentation
- [CHANGELOG.md](CHANGELOG.md) - Full release history
- [docs/SIEM_INTEGRATION.md](docs/SIEM_INTEGRATION.md) - SIEM setup guide
- [docs/ARCHITECTURE.md](docs/architecture-mvp.md) - System architecture

---

## 🙏 Contributors

Thanks to all contributors who made this release possible.

---

## 📥 Download

- **Source Code**: [GitHub Releases](https://github.com/aegisgatesecurity/aegisgate/releases/tag/v0.18.0)
- **Docker Image**: `ghcr.io/aegisgatesecurity/aegisgate:v0.18.0`
- **Go Module**: `github.com/aegisgatesecurity/aegisgate v0.18.0`

---

**Previous Release**: [v0.17.0](RELEASE_NOTES_v0.17.0.md) - SIEM Integration
