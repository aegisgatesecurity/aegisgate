# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.1] - 2026-03-11

### Fixed
- Context keys issue in pkg/api/versioning.go - properly defined apiVersionKey and versionInfoKey
- ShouldBlock logic in pkg/scanner/patterns.go - severity threshold corrected
- Dockerfile configuration for production builds

### Changed
- Repository cleanup - removed 150+ temporary files
- Updated .gitignore to exclude development artifacts

---

## [1.0.0] - 2026-03-10

### Added

#### Core Features
- **AI API Proxy** - Full proxy support for OpenAI, Anthropic, Cohere, Azure OpenAI, AWS Bedrock
- **Compliance Frameworks** - OWASP, SOC2, GDPR, HIPAA, PCI-DSS, NIST, ISO 27001
- **ML Anomaly Detection** - Traffic pattern analysis, cost anomalies, threat detection
- **Multi-Provider Support** - Unified interface for multiple AI providers
- **Rate Limiting** - Token bucket algorithm with tier-based limits
- **Connection Limiting** - Concurrent connection management

#### Security
- **TLS/HTTPS Termination** - Secure all traffic
- **JWT Validation** - Token-based authentication
- **API Key Management** - Generate and rotate API keys
- **RBAC** - Role-based access control

#### Observability
- **Metrics** - Prometheus-compatible metrics endpoint
- **Logging** - Structured JSON logging
- **Dashboard** - Web-based admin dashboard

#### Deployment
- **Docker Support** - Full container support
- **Docker Compose** - Local development setup
- **Environment Configuration** - File-based configuration

### Tiers
- **Community** - Free tier with core features
- **Developer** - $29/mo with PostgreSQL, SSO, higher limits
- **Professional** - $99/mo with full compliance, multi-tenancy
- **Enterprise** - Custom pricing with all features

### Documentation
- **README.md** - Project overview and quick start
- **Getting Started Guide** - 5-minute setup guide
- **Configuration Reference** - Complete environment variable documentation
- **Feature Comparison Matrix** - Tier feature comparison
- **Pricing Page** - Public pricing information
- **FAQ** - Frequently asked questions
- **CONTRIBUTING.md** - Contribution guidelines
- **CODE_OF_CONDUCT.md** - Community guidelines
- **SECURITY.md** - Security policy and vulnerability reporting

### Infrastructure
- **GitHub Actions CI/CD** - Automated testing and building
- **Issue Templates** - Structured bug reports and feature requests
- **Makefile** - Build automation

---

## [Unreleased]

### Planned Features

#### Coming Soon
- [ ] Kubernetes deployment manifests
- [ ] Terraform provider
- [ ] Helm charts
- [ ] GraphQL API
- [ ] More compliance frameworks (ISO 42001, FedRAMP)

#### Under Development
- [ ] Grafana dashboard templates
- [ ] SIEM integrations (Splunk, Elastic)
- [ ] Custom provider adapters
- [ ] Browser extension

---

## Upgrade Notes

### From v0.x to v1.0

1. **Configuration Changes**
   - Environment variable prefix changed to `AEGISGATE_`
   - Default ports: HTTP 8080, HTTPS 8443, Metrics 9090

2. **Breaking Changes**
   - API v1 is now stable
   - Some configuration keys renamed for consistency

3. **Migration**
   ```bash
   # Backup your data
   cp -r ./data ./data.backup
   
   # Update configuration
   # Review new environment variables
   
   # Restart AegisGate
   docker-compose down
   docker-compose up -d
   ```

---

## Version History

| Version | Date | Status |
|---------|------|--------|
| 1.0.0 | 2026-03-10 | Current |
| 0.9.0 | 2026-02-01 | Beta |
| 0.8.0 | 2026-01-01 | Alpha |

---

## Deprecation Policy

We will provide at least 3 months notice before removing or significantly changing any feature in a minor or major release.

---

*For older versions, see the [release archives](https://github.com/aegisgatesecurity/aegisgate/releases)*
