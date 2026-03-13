# AegisGate v1.0.3 Release Notes

**Release Date:** March 12, 2026  
**Repository:** https://github.com/aegisgatesecurity/aegisgate  
**Documentation:** https://docs.aegisgate.com

---

## Overview

AegisGate v1.0.3 is a significant production readiness release adding comprehensive Kubernetes support, enhanced workflow automation, and critical bug fixes. This release introduces production-ready Kubernetes manifests, improved CI/CD pipelines, and important stability improvements.

### Quick Install

```bash
# Using Docker
docker pull aegisgate/aegisgate:v1.0.3
docker run -d -p 8443:8443 aegisgate/aegisgate:v1.0.3

# Using Homebrew
brew upgrade aegisgate

# Binary from GitHub
# Download from https://github.com/aegisgatesecurity/aegisgate/releases/tag/v1.0.3
```

---

## What's New in v1.0.3

### 🚀 Kubernetes Production Manifests

This release adds full Kubernetes production deployment support:

| Manifest | Purpose |
|----------|---------|
| `deploy/k8s/namespace.yaml` | Dedicated namespace for AegisGate |
| `deploy/k8s/rbac.yaml` | Role-based access control for K8s |
| `deploy/k8s/deployment.yaml` | Production deployment with replicas |
| `deploy/k8s/service.yaml` | ClusterIP/LoadBalancer services |
| `deploy/k8s/hpa.yaml` | Horizontal Pod Autoscaler |
| `deploy/k8s/poddisruptionbudget.yaml` | PDB for zero-downtime updates |
| `deploy/k8s/network-policy.yaml` | Network segmentation |
| `deploy/k8s/configmap.yaml` | Configuration management |

#### Kubernetes Deployment Example

```bash
# Quick start with kubectl
kubectl apply -f deploy/k8s/

# Or use Helm (coming v1.1.0)
helm repo add aegisgate https://aegisgatesecurity.github.io/helm-charts
helm install aegisgate aegisgate/aegisgate
```

### 🔧 GitHub Actions Workflow Improvements

#### Release Workflow Fix

The Release workflow now correctly uses the Syft binary path output variable:

```yaml
- name: Install Syft
  uses: anchore/sbom-action/download-syft@v0
  id: syft

- name: Generate SBOM
  run: |
    ${{ steps.syft.outputs.syft-path }} ./ -o cyclonedx-json=sb.cyclonedx.json
```

This fixes the previous failure:
```
./syft_*_linux_amd64/syft: No such file or directory
Error: Process completed with exit code 127
```

#### Automated Version Management

- `version-check.yml`: Validates version consistency across VERSION file, go.mod, and GitHub releases
- `version-sync.yml`: Automated version synchronization

### 📦 Helm Chart Improvements

- Fixed Helm chart version compatibility
- Added production-ready values
- Support for custom TLS certificates
- Integrated health checks and probes

---

## Detailed Changes

### Code Quality Improvements

#### Package Improvements (`pkg/compliance/`)

```
✓ Enhanced compliance module architecture
✓ Improved framework mapping registry
✓ Extended OWASP compliance checks
✓ Added SOC2 framework support
✓ Better error handling and logging
```

#### Documentation Updates

```
✓ README contact info and copyright updated
✓ Getting started guide enhanced
✓ Architecture documentation revised
✓ Contributing guide improved
```

### Build System

| Change | Description |
|--------|-------------|
| **Go Version** | Now requires Go 1.21+ |
| **Cross-Platform** | Full support for Linux, macOS, Windows |
| **Architectures** | amd64, arm64, armv7 |
| **Docker** | Multi-stage builds for minimal images |

---

## Breaking Changes

**None** - v1.0.3 maintains full backward compatibility with v1.0.x.

---

## Known Issues

| Issue | Status | Workaround |
|-------|--------|------------|
| HTTP/3 tests | Skipped in CI | Use `-skip TestHTTP3` for local testing |
| Helm chart | v1.1.0 | Use direct K8s manifests for now |
| Terraform provider | Planned | Not yet available |

---

## Upgrading from v1.0.2

### For Docker Users

```bash
# Pull the latest image
docker pull aegisgate/aegisgate:v1.0.3

# Restart your container
docker-compose down
docker-compose up -d
```

### For Binary Installation

```bash
# Download the appropriate binary from GitHub Releases
# https://github.com/aegisgatesecurity/aegisgate/releases/tag/v1.0.3

# Linux/macOS
chmod +x aegisgate-v1.0.3-*
mv aegisgate-v1.0.3-* /usr/local/bin/aegisgate

# Windows
move aegisgate-v1.0.3-*.exe C:\Windows\System32\aegisgate.exe
```

### For Kubernetes Users

```bash
# Update deployment
kubectl apply -f deploy/k8s/

# Verify pods
kubectl get pods -n aegisgate
kubectl logs -n aegisgate -l app=aegisgate
```

---

## Architecture Highlights

### Multi-Protocol Proxy Support

AegisGate v1.0.3 continues to support:

- **HTTP/1.1** - Legacy compatibility
- **HTTP/2** - Multiplexed connections
- **HTTP/3** - QUIC-based transport
- **gRPC** - Protocol Buffers API
- **WebSocket** - Bidirectional communication
- **mTLS** - Mutual TLS authentication

### Security Layer

| Feature | Implementation |
|---------|----------------|
| TLS Termination | TLS 1.3 with PFS |
| mTLS | Service mesh compatible |
| RBAC | Fine-grained permissions |
| Audit | Immutable JSON logging |
| Secrets | HashiCorp Vault integration |

### AI/ML Integration

- Anomaly detection with configurable thresholds
- Pattern recognition for traffic analysis
- STIX/TAXII threat intelligence feeds
- Behavioral analysis engine

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
| **SLA** | ❌ | 99.5% | 99.9% | 99.99% |

---

## Community Resources

- **Documentation:** https://docs.aegisgate.com
- **Discord:** https://discord.gg/aegisgate
- **Twitter:** @aegisgatesec
- **GitHub Discussions:** https://github.com/aegisgatesecurity/aegisgate/discussions

### Contributing

We welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for setup instructions and coding standards.

---

## Security Notes

For security vulnerabilities, please follow our [SECURITY.md](SECURITY.md) reporting process.

**Reporting:** security@aegisgatesecurity.io  
**Bug Bounty:** https://aegisgatesecurity.io/bugbounty

---

## Contributors

Thank you to our contributors who made this release possible:

- @aegisgatesecurity - Core development
- Kubernetes and DevOps improvements
- Community bug reporters

---

## Checksums

Verify your download:

```bash
# Example for Linux AMD64
sha256sum -c aegisgate-1.0.3-linux-amd64.tar.gz.sha256
```

---

## What's Next (v1.1.0)

Planned for Q2 2026:

- **Terraform Provider** - Infrastructure as Code support
- **SIEM Integrations** - Splunk, Elastic, QRadar
- **GraphQL Subscriptions** - Real-time updates
- **Additional Compliance** - ISO 42001, FedRAMP
- **Enhanced ML Features** - Custom ML models

---

## Changelog Summary

### From v1.0.2 → v1.0.3

```
3444ff3 fix: use syft-path output variable instead of hardcoded glob pattern
ac57c0c v1.0.3: Add K8s production manifests (HPA, RBAC, PDB, NetworkPolicy) and fix Helm chart version
4a537db v1.0.3: Add K8s production manifests (HPA, RBAC, PDB, NetworkPolicy) and fix Helm chart version
0abde57 docs: Update README contact info and copyright
85b982c chore: Sync version to 1.0.2
2d840aa Release v1.0.2
bc03a83 Fix: Add Go doc comments for all exported symbols
743abec Fix: Add Go doc comments to resolve linter errors
044ec1a Fix: Add Go doc comments, remove Padlock rebrand artifacts, update .gitignore
```

---

**Download AegisGate v1.0.3:** https://github.com/aegisgatesecurity/aegisgate/releases/tag/v1.0.3

**Built with ❤️ by AegisGate Security**