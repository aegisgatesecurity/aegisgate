# 🔐 AegisGate - AI Security Gateway

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Dual_MIT%2FCommercial-blue.svg)](LICENSE)
[![Version](https://img.shields.io/github/v/release/aegisgatesecurity/aegisgate)](https://github.com/aegisgatesecurity/aegisgate/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/aegisgatesecurity/aegisgate/ci.yml?branch=main)](https://github.com/aegisgatesecurity/aegisgate/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/aegisgatesecurity/aegisgate)](https://goreportcard.com/report/github.com/aegisgatesecurity/aegisgate)

**Enterprise-Grade AI/LLM Security Gateway with Real-Time Threat Detection, Phase 2 Architecture Maturity Features, SIEM Integration, PKI Attestation, Threat Intelligence Sharing, and Enterprise SSO**

---

## Table of Contents

- [Overview](#overview)
- [Why AegisGate?](#why-aegisgate)
- [Features](#features)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Installation](#installation)
- [Configuration](#configuration)
- [PKI Attestation (v0.18.3+)](#pki-attestation-v0183)
- [Phase 2 Architecture Maturity Features (v0.19.0-v0.21.0)](#phase-2-architecture-maturity-features-v0190-v0210)
- [SIEM Integration](#siem-integration)
- [Threat Intelligence](#threat-intelligence)
- [Webhook Alerting](#webhook-alerting)
- [Enterprise SSO](#enterprise-sso)
- [Security](#security)
- [Compliance Frameworks](#compliance-frameworks)
- [API Reference](#api-reference)
- [Documentation](#documentation)
- [Development](#development)
- [Deployment](#deployment)
- [Roadmap](#roadmap)
- [Changelog](#changelog)
- [Contributing](#contributing)
- [License](#license)
- [Support](#support)

---

## Overview

AegisGate is a production-ready security gateway designed to protect AI and Large Language Model (LLM) applications. It provides comprehensive threat detection, PKI attestation, SIEM integration, threat intelligence sharing, webhook alerting, enterprise SSO, and compliance validation through a high-performance reverse proxy with MITM capabilities.

### Latest Major Release: v0.21.0 (2026-02-25)

Version 0.21.0 completes Phase 2 Architecture Maturity Features with comprehensive signature verification and hash chain validation systems.

### Why AegisGate?

| Challenge | Solution |
|-----------|----------|
| AI-Specific Threats | Built-in OWASP AI Top 10 and MITRE ATLAS detection |
| Trust Lattice Vulnerabilities | Cryptographic PKI attestation for certificate verification |
| Phase 2 Architecture | Comprehensive feed-level isolation and sandboxing |
| Compliance Complexity | Pre-configured frameworks: HIPAA, PCI-DSS, NIST, ISO 42001 |
| SIEM Integration | Native support for 10+ SIEM platforms |
| Threat Intelligence | STIX 2.1 / TAXII 2.1 support |
| Enterprise SSO | SAML 2.0 and OIDC integration |
| Performance | HTTP/2 support with stream-aware rate limiting |
| Accessibility | WCAG 2.1 AA compliant admin interface |
| Multi-Language | 6 languages: EN, FR, DE, ES, JA, ZH |

### Core Capabilities

- Security Proxy: HTTPS MITM interception with dynamic certificate generation
- Threat Detection: 100+ patterns for AI-specific vulnerabilities
- PKI Attestation: Cryptographic certificate verification and backdoor prevention
- Phase 2 Isolation: Feed-specific trust domains and sandboxing
- SIEM Integration: Native integration with Splunk, QRadar, Sentinel, and more
- Threat Intelligence: STIX 2.1 / TAXII 2.1 for threat intel sharing
- Webhook Alerting: Real-time notifications with HMAC authentication
- Enterprise SSO: SAML 2.0 and OpenID Connect integration
- Compliance: One-click compliance reports for major frameworks
- HTTP/2 Support: Full protocol support with security hardening
- Metrics: Real-time monitoring with atomic counters

---

## Architecture

AegisGate Security Gateway (Reverse Proxy with MITM)
- MITM Proxy: Intercept and Forward
- Threat Detection: 100+ Patterns
- PKI Attestation: Certificate Verification
- Feed-specific Trust Domains: v0.19.0
- Feed-level Sandboxing: v0.20.0
- Digital Signature Verification: v0.21.0
- Hash Chain Validation: v0.21.0
- SIEM Integration: 10+ Platforms
- Threat Intelligence: STIX 2.1
- Enterprise SSO: SAML 2.0, OIDC

---

## Quick Start

### Docker Deployment

docker pull aegisgatesecurity/aegisgate:latest

docker run -d -p 8080:8080 -p 8443:8443 -v ./config:/etc/aegisgate/config --name aegisgate aegisgatesecurity/aegisgate:latest

### Manual Installation

wget https://github.com/aegisgatesecurity/aegisgate/releases/latest/download/aegisgate-linux-amd64

chmod +x aegisgate-linux-amd64

./aegisgate-linux-amd64 --config config/aegisgate.yaml

---

## Installation

### Prerequisites

- Go 1.24+ (for development)
- Docker 20.10+ (for containerized deployment)
- 2GB RAM minimum, 4GB recommended
- 2 CPU cores minimum

### Build from Source

git clone https://github.com/aegisgatesecurity/aegisgate.git

cd aegisgate

go build -o aegisgate ./cmd/aegisgate/main.go

./aegisgate --config config/aegisgate.yaml

---

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| AEGISGATE_HOST | 0.0.0.0 | Gateway host |
| AEGISGATE_HTTP_PORT | 8080 | HTTP port |
| AEGISGATE_HTTPS_PORT | 8443 | HTTPS port |
| AEGISGATE_LOG_LEVEL | info | Logging level |

---

## PKI Attestation (v0.18.3+)

### Overview

PKI (Public Key Infrastructure) attestation provides cryptographic certificate verification for enhanced security. This feature addresses the Trust Lattice security concept and provides comprehensive certificate validation capabilities.

### Features

- Certificate Chain Validation: Full trust path verification
- Revocation Checking: CRL and OCSP support
- Backdoor Detection: Advanced pattern recognition
- Custom Trust Anchors: Configurable root certificate authorities

---

## Phase 2 Architecture Maturity Features (v0.19.0-v0.21.0)

### v0.19.0: Feed-specific Trust Domains

Feed-specific Trust Domains provide complete isolation and policy management per threat feed with enhanced security boundaries.

### v0.20.0: Feed-level Sandboxing

Feed-level Sandboxing implements a resource-qualified container system with quota enforcement and monitoring.

### v0.21.0: Digital Signature Verification & Hash Chain Validation

**Digital Signature Verification:**
- Multi-algorithm support (RSA/ECDSA/Ed25519)
- Key Manager for key management and validation
- VerificationResult and PublicKeyInfo structures
- VerificationStats for statistics tracking

**Hash Chain Validation:**
- HashChain with Merkle tree integration
- SHA256/SHA512 hash algorithms
- HashChainEntry data structure
- MemoryHashStore implementation
- Tamper detection mechanisms

---

## SIEM Integration

### Overview

AegisGate supports integration with 10+ SIEM platforms including Splunk, Elasticsearch, IBM QRadar, Microsoft Sentinel, and more.

### Supported Platforms

- Splunk (JSON HEC): Production Ready
- Elasticsearch (JSON): Production Ready
- IBM QRadar (LEEF, JSON): Production Ready
- Microsoft Sentinel (JSON): Production Ready
- Sumo Logic (JSON): Production Ready
- LogRhythm (Syslog, JSON): Production Ready
- ArcSight (CEF): Production Ready
- AWS CloudWatch (JSON): Production Ready
- AWS Security Hub (JSON): Production Ready
- Generic Syslog (RFC 5424): Production Ready

---

## Threat Intelligence

### Overview

AegisGate supports STIX 2.1 / TAXII 2.1 for industry-wide threat intelligence sharing and collaboration.

### Features

- STIX 2.1 support for standardized threat intelligence
- TAXII 2.1 server and client support
- Feed-specific trust domains for isolation
- Digital signature verification for feed integrity
- Hash chain validation for history integrity

---

## Webhook Alerting

### Overview

Real-time webhook notifications with HMAC authentication for security event alerts.

---

## Enterprise SSO

### Overview

Enterprise Single Sign-On with SAML 2.0 and OpenID Connect integration.

### Features

- SAML 2.0 support for enterprise authentication
- OpenID Connect integration
- Role-based access control
- Session management

---

## Security

### Overview

AegisGate implements defense-in-depth security across all layers of the application.

### Security Features

- Transport Security: TLS 1.2+ with HTTP/2 support
- Authentication: OAuth 2.0 with PKCE, session encryption
- Authorization: RBAC with fine-grained permissions
- Audit Logging: Comprehensive event logging
- Input Validation: Pattern scanning and sanitization
- Security Headers: OWASP compliant headers

---

## Compliance Frameworks

### Overview

AegisGate supports multiple compliance frameworks with pre-configured controls and reporting.

### Supported Frameworks

| Framework | Status | Features |
|-----------|--------|----------|
| HIPAA | Production Ready | PHI protection, audit logging |
| PCI-DSS | Production Ready | Cardholder data protection |
| NIST 800-53 | Production Ready | Federal security controls |
| ISO 42001 | Production Ready | AI management standard |
| SOC 2 | Production Ready | Service organization controls |

---

## Documentation

### Comprehensive Documentation

Full documentation is available in the docs/ directory:

- ARCHITECTURE.md - System architecture overview
- SIEM_INTEGRATION.md - SIEM platform integration guide
- SECURITY.md - Security features and best practices
- DEPLOYMENT_GUIDE.md - Deployment procedures
- COMPLIANCE_REPORTS_GUIDE.md - Compliance framework guide
- ATLAS_FRAMEWORK.md - MITRE ATLAS implementation
- CONFIGURATION.md - Configuration reference

---

## Development

### Getting Started

1. Clone the repository
2. Install Go 1.24+
3. Run tests: go test ./...
4. Build: go build ./cmd/aegisgate

### Project Structure

pkg/ - Core packages
- auth/ - Authentication
- certificate/ - Certificate management
- compliance/ - Compliance frameworks
- config/ - Configuration
- dashboard/ - Admin dashboard
- pkiattest/ - PKI attestation (v0.18.3+)
- proxy/ - MITM proxy (with PKI integration)
- siem/ - SIEM integration
- sso/ - Enterprise SSO
- threatintel/ - Threat intelligence
- trustdomain/ - Feed-specific trust domains (v0.19.0)
- sandbox/ - Feed-level sandboxing (v0.20.0)
- signature_verification/ - Digital signatures (v0.21.0)
- hash_chain/ - Hash chain validation (v0.21.0)

---

## Deployment

### Kubernetes Deployment

kubectl apply -f deploy/kubernetes/

kubectl get pods -l app=aegisgate

### Docker Deployment

docker build -t aegisgate:latest .

docker run -d -p 8080:8080 -p 8443:8443 aegisgate:latest

---

## Roadmap

### Current Release: v0.21.0 (2026-02-25)

- Phase 2 Architecture Maturity Features Complete
- Digital Signature Verification Implemented
- Hash Chain Validation Implemented

### Upcoming Releases

- v0.22.0 (Q3 2026): Integration Testing and Performance Validation
- v0.23.0 (Q4 2026): Advanced ML-based Threat Detection
- v0.24.0 (Q1 2027): Enhanced UI/UX and Multi-language Support

---

## Changelog

See CHANGELOG.md for detailed release history.

### Recent Releases

- v0.21.0 (2026-02-25): Phase 2 Architecture Maturity Features Complete
- v0.20.0 (2026-02-20): Feed-level Sandboxing
- v0.19.0 (2026-02-15): Feed-specific Trust Domains
- v0.18.3 (2026-02-24): PKI Attestation

---

## Contributing

We welcome contributions! Please see CONTRIBUTING.md for details.

---

## License

Dual licensed under:
- MIT License (LICENSE-MIT)
- Commercial License (LICENSE-COMMERCIAL)

---

## Support

### Documentation

- Documentation: https://docs.aegisgate.security
- API Reference: https://docs.aegisgate.security/api

### Community

- GitHub Issues: https://github.com/aegisgatesecurity/aegisgate/issues
- Discussions: https://github.com/aegisgatesecurity/aegisgate/discussions

---

*Last updated: 2026-02-25