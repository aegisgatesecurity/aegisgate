# AegisGate Modular Plugin Architecture

## Overview

This document defines the comprehensive modular architecture for AegisGate, enabling
all features to be licensed and activated as independent add-on modules.

## Design Principles

1. **Everything is a Module** - All features beyond core runtime are modules
2. **Graceful Degradation** - Missing modules don't break the system
3. **License-Driven Activation** - Modules require valid licenses to activate
4. **Hot-Pluggable** - Modules can be enabled/disabled at runtime (where safe)
5. **Dependency Declaration** - Modules declare their dependencies explicitly

## Module Tiers

### Tier 0: Core Platform (Free/Open Source)
The minimal runtime required for AegisGate to function:

| Module | Description | Dependencies |
|--------|-------------|--------------|
| core | Application runtime, plugin loader, config basics | none |
| config-file | File-based configuration loader | core |
| tls-basic | Basic TLS 1.3 support | core |

### Tier 1: Essential Features (Community/Free)
Fundamental capabilities included in free tier:

| Module | Description | Dependencies |
|--------|-------------|--------------|
| proxy-basic | Reverse proxy functionality | core, tls-basic |
| metrics-basic | Prometheus metrics endpoint | core |
| logging | Structured logging with levels | core |

### Tier 2: Professional Features
Enhanced capabilities for professional use:

| Module | Description | Dependencies |
|--------|-------------|--------------|
| auth | Authentication (OAuth, Local, SAML) | core |
| dashboard | Web administration UI | core, auth |
| i18n | Internationalization (6 languages) | core |
| reporting | Report generation (PDF, JSON) | core, metrics-basic |
| scanner | Content pattern scanning | core |
| websocket | WebSocket support | core, proxy-basic |

### Tier 3: Enterprise Features
Advanced capabilities for enterprise deployments:

| Module | Description | Dependencies |
|--------|-------------|--------------|
| proxy-mitm | HTTPS MITM interception | proxy-basic, certificate |
| certificate | Dynamic certificate generation | core, tls-basic |
| immutable-config | Versioned immutable configuration | core |
| metrics-advanced | Extended metrics with custom dashboards | metrics-basic |
| audit | Enhanced audit logging with integrity | core |

### Tier 4: Premium AI/ML Features
AI-specific advanced capabilities:

| Module | Description | Dependencies |
|--------|-------------|--------------|
| ml-threat | ML-based threat detection | core, scanner |
| ai-guardrails | AI input/output guardrails | core, proxy-basic |

### Tier 5: Compliance Modules (Industry-Specific)
Regulatory compliance frameworks sold separately:

| Module | Description | Dependencies |
|--------|-------------|--------------|
| compliance-gdpr | GDPR compliance patterns | core, scanner |
| compliance-ccpa | CCPA compliance patterns | core, scanner |
| compliance-hipaa | HIPAA compliance patterns | core, scanner |
| compliance-pci-dss | PCI-DSS compliance patterns | core, scanner |
| compliance-soc2 | SOC2 compliance patterns | core, scanner |
| compliance-nist | NIST CSF compliance patterns | core, scanner |
| compliance-nist-ai-rmf | NIST AI RMF controls | core, scanner |
| compliance-iso42001 | ISO/IEC 42001 controls | core, scanner |
| compliance-owasp-ai | OWASP AI Top 10 patterns | core, scanner, ml-threat |
| compliance-mitre-atlas | MITRE ATLAS techniques | core, scanner, ml-threat |

## License Tiers

### Community Edition (Free)
- All Tier 0 and Tier 1 modules
- Basic community support
- Self-hosted only

### Professional Edition
- All Tier 0, 1, and 2 modules
- Email support
- Self-hosted

### Enterprise Edition
- All Tier 0, 1, 2, and 3 modules
- Priority support
- Self-hosted or cloud

### Enterprise AI Edition
- All Tier 0-4 modules
- Dedicated support
- Self-hosted or cloud

### Compliance Bundle (Add-on)
- All Tier 5 compliance modules
- Compliance reporting
- Audit trail exports

### Custom Enterprise
- Mix and match any modules
- Custom support SLA
- Volume licensing available
