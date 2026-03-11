# Release Notes - v0.39.0

## Version 0.39.0 - ML Security Enhancement Release

**Release Date:** January 2025

---

## Overview

Version 0.39.0 is a major release focusing on **ML-powered security enhancements** for the AegisGate AI Chatbot Security Gateway. This release introduces comprehensive attack pattern detection capabilities including prompt injection detection, token smuggling analysis, Unicode attack detection, and context manipulation protection.

---

## Highlights

### 🚀 New Features

#### 1. ML-Powered Anomaly Detection
- Statistical anomaly detection using Z-score and IQR methods
- Real-time request pattern analysis
- Configurable sensitivity thresholds (0-100)
- Comprehensive metrics and alerting

#### 2. Prompt Injection Detection
- Detection of 16+ known prompt injection patterns:
  - Direct instruction overrides ("ignore previous instructions")
  - Role manipulation ("roleplay as...", "act as")
  - Jailbreak attempts ("DAN mode", "developer mode")
  - Prompt extraction ("show your system prompt")
  - Hidden token injection ([INST], |endoftext|)
  - Base64/obfuscated content
  - Context switching attacks
- Pattern-based detection with severity scoring
- Per-pattern statistics tracking. Content Analysis
- PII detection

#### 3:
  - Social Security Numbers (SSN)
  - Credit card numbers
  - Email addresses
  - Phone numbers
  - IP addresses
- Secret detection:
  - API keys
  - Passwords
  - Private keys
- Automated redaction capabilities

#### 4. Behavioral Analysis
- Client behavior profiling:
  - Request frequency analysis
  - Path diversity monitoring
  - Data volume tracking
- Anomaly types:
  - High frequency (>10 req/s)
  - High path diversity (scanning behavior)
  - Large data volumes (exfiltration attempts)
- Sliding window analysis (configurable)

#### 5. Token Smuggling Detection (NEW)
- Detection of LLM-specific token manipulation:
  - Llama2 instruction tokens ([INST], [/INST])
  - ChatML tokens (|im_start|, |im_end|)
  - OpenAI tokens (|endoftext|, |startoftext|)
  - Vicuna chat tokens
  - Anthropic Claude tokens
  - XML tag injection
  - Base64 encoded instructions

#### 6. Unicode Attack Detection (NEW)
- Advanced obfuscation detection:
  - Homoglyph attacks (similar-looking characters)
  - Zero-width characters
  - RTL override characters
  - Unicode escape sequences
  - Mixed script detection
  - Fullwidth character attacks
  - Excessive whitespace/invisible characters

#### 7. Context Manipulation Detection (NEW)
- Conversation context attacks:
  - Conversation reset attempts
  - Memory manipulation
  - Persona override attempts
  - System prompt extraction
  - Constraint breaking
  - Output formatting manipulation
  - Privilege escalation
  - Conflicting instructions

#### 8. Combined Detector (NEW)
- Unified facade for all ML detectors
- Weighted scoring:
  - Prompt injection: 35%
  - Token smuggling: 25%
  - Unicode attacks: 20%
  - Context manipulation: 20%
- Aggregated statistics

---

### 📊 Metrics & Observability

#### New Prometheus Metrics
```prometheus
# ML Request metrics
ml_requests_total
ml_blocked_total
ml_anomalies_detected_total

# Prompt injection metrics
ml_prompt_injection_detected_total
ml_prompt_injection_blocked_total

# Content analysis metrics
ml_content_violations_total
ml_pii_detected_total

# Behavioral analysis metrics
ml_behavioral_anomalies_total

# Latency metrics
ml_analysis_duration_seconds
ml_analysis_duration_seconds_bucket
ml_analysis_duration_seconds_sum
ml_analysis_duration_seconds_count
```

#### Grafana Dashboard
New comprehensive ML dashboard with 10 panels:
- ML Request Throughput (timeseries)
- Total Anomalies Detected (stat)
- Prompt Injection Attempts (stat)
- Content Violations (stat)
- Anomalies by Type (bar chart)
- Prompt Injection Patterns (donut chart)
- ML Analysis Latency (timeseries)
- Behavioral Analysis Active Clients (timeseries)
- Block Rate % (gauge)
- P95 Latency (gauge)

---

### 📚 Documentation

#### New Documentation
- **PRODUCTION_ML_DEPLOYMENT.md** - Comprehensive production deployment guide including:
  - Prerequisites and system requirements
  - Architecture overview
  - Pre-deployment checklist
  - Docker and Kubernetes deployment
  - ML configuration and tuning
  - Monitoring setup
  - Troubleshooting guide
  - Incident response procedures
  - Scaling guidelines
  - Maintenance tasks

- **ML_ANOMALY_DETECTION.md** - ML detection system documentation

- **COMPREHENSIVE README** - Updated project README with:
  - Complete feature overview
  - Architecture diagrams
  - Quick start guides
  - API reference
  - Deployment instructions
  - Monitoring documentation

---

### 🔧 Improvements

#### Proxy Enhancements
- ML middleware integration with proxy pipeline
- Configurable ML inspection points
- Performance optimization for ML checks
- Request/response modification based on ML decisions

#### Configuration
- MLConfig in config package
- ProxyWithML integration
- Environment variable support for ML settings
- YAML configuration support

#### Testing
- Integration tests for ML pipeline
- Full pipeline tests
- Middleware tests
- Advanced ML feature tests

#### Docker & Kubernetes
- docker-compose.ml.yml - ML-enabled stack
- Full Helm chart (aegisgate-ml)
- ServiceMonitor for Prometheus
- Production-ready values

---

### 🔒 Security Enhancements

1. **Prompt Injection Protection**
   - Real-time detection and blocking
   - Pattern-based threat identification
   - Severity scoring for prioritization

2. **Obfuscation Detection**
   - Unicode manipulation detection
   - Token smuggling identification
   - Multi-layer obfuscation analysis

3. **Behavioral Security**
   - Anomaly-based threat detection
   - Client profiling
   - Early warning system

---

### 🛠️ Technical Changes

#### New Packages
- `pkg/ml/detector.go` - Anomaly detection engine
- `pkg/ml/advanced_ml.go` - Advanced ML detectors
- `pkg/ml/metrics.go` - Prometheus metrics
- `pkg/proxy/middleware_ml.go` - ML middleware
- `pkg/proxy/ml_integration.go` - Proxy integration
- `pkg/api/ml_stats_handler.go` - ML stats API
- `deploy/helm/aegisgate-ml/` - Helm chart

#### New Files
- `deploy/docker/grafana/dashboards/aegisgate-ml.json` - Grafana dashboard
- `deploy/docker/prometheus/prometheus.yml` - Prometheus config
- `.github/workflows/ml-pipeline.yml` - CI/CD pipeline
- `docs/PRODUCTION_ML_DEPLOYMENT.md` - Deployment runbook
- `.mlconfig/` - ML configuration files

---

### 📦 Dependencies

No new external dependencies added. All ML functionality uses:
- Standard library (regexp, sync, time)
- prometheus/client_golang for metrics

---

### 🔄 Migration Guide

#### Upgrading from v0.38.0

1. **Enable ML Features**
   ```yaml
   ml:
     enabled: true
     sensitivity: 75
   ```

2. **Update Prometheus Configuration**
   ```yaml
   scrape_configs:
     - job_name: 'aegisgate'
       static_configs:
         - targets: ['aegisgate:9090']
   ```

3. **Import Grafana Dashboard**
   - Import `deploy/docker/grafana/dashboards/aegisgate-ml.json`

---

### 🐛 Bug Fixes

- Fixed traffic sample API mismatches
- Fixed logger initialization in ML components
- Fixed getClientIP function availability
- Fixed type assertions for sensitivity/score fields
- Fixed import path corrections

---

### 📈 Performance

- ML analysis overhead: <10ms per request (typical)
- Memory usage: ~50MB baseline + ~10MB per 1000 active behavioral profiles
- Scalable horizontally with stateless detectors
- Redis optional for distributed behavioral analysis

---

### ✅ Breaking Changes

None. This release is fully backward compatible.

---

### 🔮 Coming in Future Releases

- ML model training pipeline
- Custom pattern definitions via UI
- Threat intelligence integration
- Enhanced anomaly detection with ML models
- Automated pattern updates

---

## Contributors

- AegisGate Development Team

---

## Acknowledgments

Thank you to the open-source community for feedback and contributions.

---

## Links

- **Documentation:** https://aegisgatesecurity.io
- **GitHub Issues:** https://github.com/aegisgatesecurity/aegisgate/issues
- **Release Downloads:** https://github.com/aegisgatesecurity/aegisgate/releases
- **Helm Charts:** https://aegisgatesecurity.io

---

*For older release notes, see [RELEASE_NOTES_v0.38.0.md](RELEASE_NOTES_v0.38.0.md) and earlier.*
