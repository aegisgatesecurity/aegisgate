# AegisGate ML-Enabled Helm Chart

## Installation

### Quick Start
```bash
# Add the repository
helm repo add aegisgate https://aegisgatesecurity.github.io/aegisgate

# Install with default values
helm install aegisgate aegisgate/aegisgate

# Install with ML enabled (default)
helm install aegisgate-ml aegisgate/aegisgate-ml \
  --set ml.enabled=true \
  --set ml.sensitivity=medium
```

### From Source
```bash
# Clone the repository
git clone https://github.com/aegisgatesecurity/aegisgate.git

# Install from local chart
cd aegisgate/deploy/helm/aegisgate-ml
helm install aegisgate-ml . \
  --set ml.enabled=true \
  --set server.upstream=http://my-llm:8080
```

## ML Configuration

The chart includes full ML anomaly detection support:

```yaml
ml:
  enabled: true
  sensitivity: "medium"  # low, medium, high, paranoid
  blockOnCritical: true
  blockOnHigh: false
  sampleRate: 100

mlAdvanced:
  promptInjection:
    enabled: true
    sensitivity: 75
  contentAnalysis:
    enabled: true
  behavioralAnalysis:
    enabled: true
```

## Values Reference

| Parameter | Description | Default |
|-----------|-------------|---------|
| `ml.enabled` | Enable ML anomaly detection | `true` |
| `ml.sensitivity` | Detection sensitivity | `medium` |
| `mlAdvanced.promptInjection.enabled` | Prompt injection detection | `true` |
| `mlAdvanced.contentAnalysis.enabled` | Content analysis | `true` |
| `mlAdvanced.behavioralAnalysis.enabled` | Behavioral analysis | `true` |

## Upgrading

```bash
# Upgrade to new version
helm upgrade aegisgate aegisgate/aegisgate-ml

# Upgrade with new ML settings
helm upgrade aegisgate aegisgate/aegisgate-ml \
  --set ml.sensitivity=high \
  --set ml.blockOnHigh=true
```

## Uninstalling

```bash
helm uninstall aegisgate
```

## Examples

### Production Deployment
```bash
helm install aegisgate-prod aegisgate/aegisgate-ml \
  --set ml.sensitivity=high \
  --set ml.blockOnCritical=true \
  --set ml.blockOnHigh=true \
  --set autoscaling.enabled=true \
  --set autoscaling.minReplicas=3 \
  --set ingress.enabled=true \
  --set ingress.hosts[0].host=aegisgate.example.com
```

### Development/Testing
```bash
helm install aegisgate-dev aegisgate/aegisgate-ml \
  --set ml.sensitivity=paranoid \
  --set ml.logAllAnomalies=true \
  --set logging.level=debug \
  --set replicas=1 \
  --set resources.limits.cpu=1 \
  --set resources.limits.memory=1Gi
```

### With SIEM Integration
```bash
helm install aegisgate aegisgate/aegisgate-ml \
  --set siem.enabled=true \
  --set siem.format=splunk \
  --set siem.endpoint=https://splunk:8088 \
  --set siem.apiKeySecret=my-splunk-key
```
