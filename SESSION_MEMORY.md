# AegisGate - Project Session Memory

**Last Updated:** March 12, 2026  
**Current Version:** 1.0.3  
**Go Version Required:** 1.21+ (go.mod says 1.24)

---

## 📋 Project Overview

AegisGate is an enterprise-grade AI API security platform that acts as a secure proxy gateway between applications and AI providers (OpenAI, Anthropic, Azure, AWS Bedrock, Cohere). It provides comprehensive security, observability, and compliance features.

### Core Capabilities
- **AI API Proxy** - Transparent proxy with full request/response inspection
- **Security Scanning** - Prompt injection, PII detection, malicious payload blocking
- **Observability** - Prometheus metrics, structured logging, dashboard
- **Compliance** - SOC2, GDPR, HIPAA, PCI-DSS, OWASP, NIST, ISO 27001 ready
- **ML Anomaly Detection** - Traffic pattern analysis, cost anomaly detection
- **Multi-tier Access** - Community, Developer, Professional, Enterprise tiers

---

## 🏗️ Current Project Structure

```
AegisGate/
├── cmd/
│   ├── aegisgate/main.go      # Main application entry point
│   ├── debug/main.go          # Debug utilities
│   └── gencerts/main.go       # Certificate generation tool
├── pkg/                       # Core packages (~166KB)
│   ├── api/                   # API server, caching, versioning
│   ├── auth/                  # Authentication, OAuth, session management
│   ├── compliance/            # SOC2, HIPAA, OWASP, PCI, Atlas frameworks
│   ├── proxy/                 # HTTP/2, HTTP/3, mTLS proxy (86KB - largest)
│   ├── ml/                    # ML anomaly detection
│   ├── siem/                  # SIEM integrations (Splunk, Elastic, QRadar)
│   ├── threatintel/           # STIX/TAXII threat feeds
│   ├── tls/                   # TLS/mTLS certificate management
│   ├── webhook/               # Webhook management
│   └── ...                    # Many more packages
├── deploy/
│   ├── docker/                # Dockerfiles and docker-compose
│   ├── helm/aegisgate/        # Helm chart (present but potentially incomplete)
│   └── k8s/                   # Kubernetes manifests
├── .github/workflows/         # GitHub Actions
├── docs/                      # Extensive documentation (27KB)
├── tests/integration/         # Integration tests
├── tests/load/                # Load/benchmark tests
├── VERSION                    # Current version: 1.0.3
├── go.mod                     # Go module definition
├── README.md                  # Main documentation
└── RELEASE_NOTES_v1.0.3.md   # Comprehensive release notes
```

---

## 🚀 Current Release Status

### Version: v1.0.3 (Latest)
**Release Date:** March 2026  
**Status:** ✅ Released and pushed to GitHub

### Recent Git History
```
944379f docs: Add Kubernetes deployment section and release notes for v1.0.3
84829d6 docs: Update README and release notes
3444ff3 fix: use syft-path output variable instead of hardcoded glob pattern
ac57c0c v1.0.3: Add K8s production manifests (HPA, RBAC, PDB, NetworkPolicy)
```

### Git Tags
- `v1.0.1` - Initial community release
- `v1.0.2` - Stable release with core features
- `v1.0.3` - Current - K8s production manifests, workflow fixes

---

## 🔧 GitHub Actions Workflows

### Release Workflow (`.github/workflows/release.yml`)
- **Trigger:** Push of any `v*` tag
- **Builds:** Multi-platform Docker images (amd64, arm64, arm)
- **SBOM Generation:** Uses Syft with correct path variable (`steps.syft.outputs.syft-path`)
- **Artifacts:** Linux/Windows/macOS binaries + SBOMs
- **Provenance:** Uses `actions/attest-build-provenance`

**PRO TIP:** The workflow was failing with "No such file or directory" for Syft. The fix was to use the output variable from the `anchore/sbom-action/download-syft@v1` step instead of hardcoded paths like `./syft_*_linux_amd64/syft`.

```yaml
# CORRECT way to use Syft in GitHub Actions:
- name: Install Syft
  uses: anchore/sbom-action/download-syft@v1
  id: syft

- name: Generate SBOM
  run: |
    ${{ steps.syft.outputs.syft-path }} ./ -o cyclonedx-json=sb.cyclonedx.json
```

---

## ☸️ Kubernetes Deployment

### Included K8s Manifests (`deploy/k8s/`)
| File | Purpose |
|------|---------|
| `namespace.yaml` | Dedicated `aegisgate` namespace |
| `rbac.yaml` | ServiceAccount, ClusterRole, ClusterRoleBinding |
| `deployment.yaml` | Production deployment with replicas |
| `service.yaml` | ClusterIP service |
| `hpa.yaml` | Horizontal Pod Autoscaler |
| `poddisruptionbudget.yaml` | PDB for zero-downtime updates |
| `network-policy.yaml` | Pod-level network security |
| `configmap.yaml` | Configuration management |

### Deploy Commands
```bash
# Direct kubectl
kubectl apply -f deploy/k8s/

# Or use Helm (chart may need verification)
helm install aegisgate ./helm/aegisgate
```

---

## 🐛 Known Issues & Gotchas

### Critical Gotchas

1. **Helm Chart Status:** The README references `helm/aegisgate/` and `helm/aegisgate-ml/` directories. These exist in the tree but may be incomplete or placeholders. **Verify before production use.**

2. **Windows Git Commit Issue:** Git on Windows fails with `git commit -m "message with spaces"` due to path parsing errors. **Workaround:** Use `git commit -F <file>` with a message file.

3. **Go Version Mismatch:** `go.mod` declares `go 1.24.0` but README states "Go 1.21+ required" and GitHub Actions uses Go 1.24. This could cause confusion. Recommend aligning to 1.21 for broader compatibility or 1.24 consistently.

4. **Documentation vs. Reality:** README mentions `k8s/production/` directory but actual path is `deploy/k8s/`. This was documented in the README after the discrepancy was discovered.

5. **HTTP/3 Tests Skipped:** Some HTTP/3 integration tests are skipped in CI. Use `-skip TestHTTP3` for local testing if needed.

6. **Tool Instability:** During development, MCP tools sometimes became unavailable intermittently. If tools fail, retry or use simpler shell commands.

7. **PowerShell Here-String Parsing:** Using PowerShell here-strings (`@'...'@`) with special characters caused "missing terminator" errors. Avoid complex content in here-strings.

### Pro Tips

1. **SBOM Workflow Fix:** Always use `${{ steps.syft.outputs.syft-path }}` instead of hardcoded paths in GitHub Actions for better portability.

2. **Version Consistency:** Keep VERSION file, go.mod, and git tags in sync. The project has `version-check.yml` workflow to validate this.

3. **Docker Multi-platform:** Use `docker buildx` for cross-platform builds (already configured in release.yml).

4. **Helm Chart Version:** When updating Helm charts, ensure chart version in `Chart.yaml` matches the application version to avoid confusion.

---

## 📝 Session Context - Recent Changes

### Session Date: March 12, 2026

**Task Completed:** Fixed Syft SBOM generation in GitHub Actions release workflow and updated README with Kubernetes deployment section.

**Workflow Fix:**
- Changed from hardcoded path `./syft_*_linux_amd64/syft` to `${{ steps.syft.outputs.syft-path }}`
- This fixes the common failure: `syft: No such file or directory (exit code 127)`

**README Updates:**
- Added "Kubernetes Deployment" section with deployment examples
- Added "Release Notes" section summarizing v1.0.3
- Updated Kubernetes paths from `k8s/production/` to correct `deploy/k8s/`

**Files Modified:**
- `.github/workflows/release.yml` - Fixed Syft path issue
- `README.md` - Added K8s and Release Notes sections

---

## ✅ Project "Pristine" Status Assessment

### Is the project ready for public consumption as of v1.0.3?

**Answer: PARTIALLY YES - But with significant caveats**

| Criteria | Status | Notes |
|----------|--------|-------|
| Code Compiles | ✅ Yes | Builds successfully |
| Basic Tests | ⚠️ Partial | Some HTTP/3 tests skipped |
| Documentation | ⚠️ Partial | Some paths incorrect in older docs |
| Kubernetes Manifests | ✅ Yes | Full K8s manifests present |
| Helm Charts | ❓ Unknown | Exists but not verified functional |
| SBOM/CD Pipeline | ✅ Yes | Working release workflow |
| Security Policy | ✅ Yes | SECURITY.md present with disclosure process |
| Release Artifacts | ✅ Yes | Multi-platform binaries and Docker images |

### Recommendations Before Public Release

1. **Verify Helm Chart:** Test `helm install aegisgate ./helm/aegisgate` before advertising Kubernetes support
2. **Fix Documentation Errors:** Update docs that reference non-existent `k8s/production/` path
3. **Align Go Version:** Decide on Go 1.21 vs 1.24 and be consistent
4. **Complete Integration Tests:** Address skipped HTTP/3 tests
5. **Add Terraform Provider:** Mention in roadmap that IaC is "planned" not "in progress"

### What Works Well
- ✅ Core proxy functionality (HTTP/1.1, HTTP/2, HTTP/3, gRPC)
- ✅ Security scanning (prompt injection, PII, content filtering)
- ✅ Compliance frameworks (SOC2, HIPAA, OWASP, PCI)
- ✅ Multi-tier licensing model
- ✅ Docker deployment
- ✅ GitHub Actions CI/CD
- ✅ SBOM generation
- ✅ Basic Kubernetes manifests

### Verdict
The project is suitable for early adopters and evaluation purposes. The core functionality is solid, but the Helm chart and some documentation should be verified before heavy production use in Kubernetes environments.

---

## 🔮 Future Plans (from RELEASE_NOTES_v1.0.3.md)

### v1.1.0 (Planned Q2 2026)
- **Terraform Provider** - Infrastructure as Code support
- **SIEM Integrations** - Splunk, Elastic, QRadar enhancements
- **GraphQL Subscriptions** - Real-time updates
- **Additional Compliance** - ISO 42001, FedRAMP
- **Enhanced ML Features** - Custom ML models

---

## 📞 Key Contacts & Resources

| Resource | Link |
|----------|------|
| GitHub Repository | https://github.com/aegisgatesecurity/aegisgate |
| Security Issues | security@aegisgatesecurity.io |
| Sales | sales@aegisgatesecurity.io |
| Documentation | https://docs.aegisgate.com (coming soon) |
| Discord | https://discord.gg/aegisgate (coming soon) |

---

## 📎 Quick Reference Commands

```bash
# Build from source
make build

# Run locally
./bin/aegisgate -tier community

# Docker
docker-compose up -d

# Kubernetes
kubectl apply -f deploy/k8s/

# Run tests
make test

# Create release (tag push triggers workflow)
git tag v1.0.4 && git push origin v1.0.4
```

---

*This memory file should be updated with each significant development session.*