# AegisGate v1.0.2 Release Notes

**Release Date:** March 12, 2026  
**Repository:** https://github.com/aegisgatesecurity/aegisgate  
**Documentation:** https://docs.aegisgate.com

---

## Overview

AegisGate v1.0.2 is a critical bug fix release that resolves workflow build failures introduced during the Padlock to AegisGate rebranding. This release ensures all GitHub Actions CI/CD pipelines, Docker builds, and binary artifacts are properly generated.

### Quick Install

```bash
# Using Docker
docker pull aegisgate/aegisgate:latest

# Using Homebrew
brew tap aegisgatesecurity/aegisgate
brew install aegisgate

# From source
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate
make build
```

---

## What's New in v1.0.2

### Build System Improvements

| Change | Description |
|--------|-------------|
| **Makefile Update** | Fixed all build targets to produce `aegisgate` binary instead of `padlock` |
| **Binary Output** | Builds now produce `bin/aegisgate` (Linux/macOS) or `bin/aegisgate.exe` (Windows) |
| **Cross-Platform Builds** | Full support for Linux, macOS, Windows, FreeBSD across amd64 and arm64 architectures |

### GitHub Actions Workflow Fixes

| Workflow | Issue Fixed |
|----------|-------------|
| **release.yml** | Fixed artifact naming from `padlock-*` to `aegisgate-*` |
| **test.yml** | Updated Docker build to use correct artifact names |
| **ci.yml** | Verified working pipeline for Test, Build, and Security jobs |
| **version-check.yml** | Version consistency validation working |
| **version-sync.yml** | Version synchronization automated |

### Documentation Updates

- README.md updated with correct binary names and commands
- All Padlock references replaced with AegisGate throughout
- Makefile help text updated
- Docker deployment documentation refreshed

---

## Detailed Changes

### Code Quality Fixes (v1.0.1 → v1.0.2)

#### Compliance Package (`pkg/compliance/`)

```
✓ Fixed duplicate const blocks in compliance.go causing "undefined: Framework" errors
✓ Resolved type stuttering issues:
    - ComplianceResult → Result
    - ComplianceManager → Manager  
    - ComplianceReport → Report
✓ Added backward compatibility aliases for existing code
✓ Fixed missing Go doc comments on exported symbols:
    - SeverityCritical
    - FrameworkNIST1500
    - PermManagePolicies
    - RoleOperator
    - ProviderMicrosoft
✓ Corrected framework string values for tests:
    - FrameworkNIST1500: "NIST.AI-1.500"
    - FrameworkPCIDSS: "PCI-DSS"
```

#### Authentication Package (`pkg/auth/`)

```
✓ Added missing documentation for exported constants
✓ ProviderMicrosoft: Added proper doc comment
✓ RoleOperator: Fixed doc comment format
✓ PermManagePolicies: Added doc comment
```

#### API Package (`pkg/api/`)

```
✓ Fixed ETagger.Handle method comment format
```

#### Docker & Deployment

```
✓ deploy/docker/Dockerfile: Updated binary name from padlock to aegisgate
✓ Makefile: All targets now produce aegisgate binary
✓ .mlconfig/aegisgate.example.yml: Fixed app_name example
```

---

## Breaking Changes

**None** - This is a pure bug-fix release with full backward compatibility.

---

## Known Issues

| Issue | Status | Workaround |
|-------|--------|------------|
| HTTP/3 tests | Skipped in CI | Use `-skip TestHTTP3` flag |
| Kubernetes manifests | Under development | Use Docker Compose for now |
| Terraform provider | Planned | Not yet available |

---

## Upgrading from v1.0.1

### For Docker Users

```bash
# Pull the latest image
docker pull aegisgate/aegisgate:v1.0.2

# Restart your container
docker-compose down
docker-compose up -d
```

### For Binary Installation

```bash
# Download the appropriate binary from GitHub Releases
# https://github.com/aegisgatesecurity/aegisgate/releases/tag/v1.0.2

# Linux/macOS
chmod +x aegisgate-*-*
mv aegisgate-*-* /usr/local/bin/aegisgate

# Windows
move aegisgate-*.exe C:\Windows\System32\aegisgate.exe
```

### For Source Installation

```bash
git fetch --all --tags
git checkout v1.0.2
make build
./bin/aegisgate serve
```

---

## Comparing Plans

| Feature | Community | Developer | Professional | Enterprise |
|---------|-----------|-----------|--------------|------------|
| **Price** | Free | $29/mo | $99/mo | Custom |
| **AI Providers** | 3 | 8 | All | All |
| **Rate Limiting** | 100/min | 1,000/min | Unlimited | Unlimited |
| **Compliance** | Basic | OWASP | Full | Full |
| **RBAC** | ❌ | ✅ | ✅ | ✅ |
| **SSO/SAML** | ❌ | ✅ | ✅ | ✅ |
| **Multi-tenancy** | ❌ | ❌ | ✅ | ✅ |
| **24/7 Support** | ❌ | Email | Priority | Dedicated |

---

## Community Resources

- **Documentation:** https://docs.aegisgate.com
- **Discord:** https://discord.gg/aegisgate
- **Twitter:** @aegisgatesec
- **GitHub Discussions:** https://github.com/aegisgatesecurity/aegisgate/discussions

---

## Security Notes

For security vulnerabilities, please follow our [SECURITY.md](https://github.com/aegisgatesecurity/aegisgate/blob/main/SECURITY.md) reporting process.

---

## Contributors

Thank you to our contributors who made this release possible:

- @aegisgatesecurity - Core development
- Community bug reporters

---

## Checksums

Example checksum verification:

```bash
# Download and verify
wget https://github.com/aegisgatesecurity/aegisgate/releases/download/v1.0.2/aegisgate-linux-amd64
sha256sum -c aegisgate-linux-amd64.sha256
```

---

## What's Next (v1.1.0)

- Kubernetes Helm charts
- Grafana dashboard templates
- SIEM integration (Splunk, Elastic)
- Additional compliance frameworks (ISO 42001, FedRAMP)
- GraphQL API support

---

**Download AegisGate v1.0.2:** https://github.com/aegisgatesecurity/aegisgate/releases/tag/v1.0.2

*Built with ❤️ by AegisGate Security*