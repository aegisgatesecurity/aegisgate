## AegisGate v0.4.0 - NIST AI RMF Framework Integration

### Changes in v0.4.0
- Implemented new NIST AI RMF framework with 18 controls across Govern, Map, Measure, and Manage functions
- Enhanced README with comprehensive NIST AI RMF documentation
- Fixed bug in NIST_AI_RMF constant definition
- Project reaches milestone of 12 core packages and ~25,000 source lines

### New Features
- **NIST AI RMF Framework**: Complete implementation of NIST AI Risk Management Framework
  - Govern (5 controls): Organizational context, Risk management strategy, Supply chain risk management, Stakeholder engagement, Governance and workforce
  - Map (4 controls): AI capabilities and contexts, Stakeholder expectations, System purpose, RiskTHREAT identification
  - Measure (4 controls): Risk measurement, Impact analysis, Likelihood estimation, Risk determination
  - Manage (5 controls): Risk response, Incident management, Continuous monitoring, Risk acceptance, Risk documentation
- Enhanced README documentation covering all NIST AI RMF functions and controls

### Bug Fixes
- Fixed NIST_AI_RMF constant definition to properly expose the framework

### Breaking Changes
None

### Full Changelog
- feat: Implement NIST AI RMF framework with 18 controls
- feat: Add NIST AI RMF documentation to README
- fix: Correct NIST_AI_RMF constant definition

### What's New
This release introduces comprehensive NIST AI Risk Management Framework support, enabling organizations to implement AI governance, identify risks, measure AI system performance, and manage AI-related risks effectively.

### Installation
```bash
# Download release
curl -L https://github.com/aegisgatesecurity/aegisgate/archive/refs/tags/v0.4.0.tar.gz | tar xz
cd aegisgate-0.4.0

# Build
go build -o aegisgate ./cmd/aegisgate

# Run
./aegisgate
```

### Docker
```bash
docker pull aegisgatesecurity/aegisgate:v0.4.0
```

### Project Statistics
- Core Packages: 12
- Source Lines: ~25,000
- Framework: NIST AI RMF (18 controls)

### Checksums
See release assets for checksums.

**Full Changelog**: https://github.com/aegisgatesecurity/aegisgate/commits/v0.4.0
