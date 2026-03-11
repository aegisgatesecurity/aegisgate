# Production ML Deployment Runbook

## AegisGate ML Security Gateway - Production Deployment Guide

This runbook provides comprehensive instructions for deploying the AegisGate ML Security Gateway in a production environment with anomaly detection capabilities.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Architecture Overview](#architecture-overview)
3. [Pre-Deployment Checklist](#pre-deployment-checklist)
4. [Environment Configuration](#environment-configuration)
5. [Docker Deployment](#docker-deployment)
6. [Kubernetes Deployment](#kubernetes-deployment)
7. [ML Model Configuration](#ml-model-configuration)
8. [Monitoring & Observability](#monitoring--observability)
9. [Security Considerations](#security-considerations)
10. [Troubleshooting](#troubleshooting)
11. [Incident Response](#incident-response)
12. [Scaling Guidelines](#scaling-guidelines)
13. [Maintenance](#maintenance)

---

## Prerequisites

### System Requirements

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| CPU | 4 cores | 8+ cores |
| RAM | 8 GB | 16 GB |
| Disk | 50 GB SSD | 100 GB SSD |
| Network | 100 Mbps | 1 Gbps |

### Software Requirements

- **Go 1.21+** - For building the application
- **Docker 24.0+** - For containerized deployment
- **Kubernetes 1.28+** - For K8s deployment
- **Prometheus 2.45+** - For metrics collection
- **Grafana 10.0+** - For visualization

### Required Access

- Container registry access
- Kubernetes cluster access
- Prometheus endpoint access
- Grafana admin access
- Cloud provider credentials (if applicable)

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        Load Balancer                            │
└─────────────────────────────┬───────────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────────┐
│                    AegisGate Gateway                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐            │
│  │   Proxy     │  │  ML Middleware │ │  Auth      │            │
│  │   Handler   │──│  - Anomaly   │──│  Handler   │            │
│  │             │  │  - Prompt    │  │            │            │
│  │             │  │    Injection │  │            │            │
│  │             │  │  - Content   │  │            │            │
│  │             │  │    Analysis  │  │            │            │
│  └─────────────┘  └─────────────┘  └─────────────┘            │
└─────────────────────────────┬───────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
┌───────────────┐   ┌───────────────┐   ┌───────────────┐
│   Prometheus  │   │     Logs      │   │  Upstream     │
│   (Metrics)   │   │   (Storage)   │   │  AI APIs      │
└───────────────┘   └───────────────┘   └───────────────┘
        │
        ▼
┌───────────────┐
│    Grafana    │
│  (Dashboards) │
└───────────────┘
```

### ML Detection Pipeline

```
Incoming Request
       │
       ▼
┌──────────────────┐
│ Request Parser   │
└────────┬─────────┘
         │
         ▼
┌──────────────────┐    ┌──────────────────┐
│ Prompt Injection │───▶│ Block/Allow      │
│ Detector         │    │ Decision         │
└────────┬─────────┘    └──────────────────┘
         │
         ▼
┌──────────────────┐
│ Content Analyzer │──▶ PII/Secret Detection
└────────┬─────────┘
         │
         ▼
┌──────────────────┐
│ Behavioral       │──▶ Anomaly Detection
│ Analyzer         │
└────────┬─────────┘
         │
         ▼
    Upstream API
```

---

## Pre-Deployment Checklist

### Infrastructure

- [ ] Load balancer configured
- [ ] SSL/TLS certificates obtained
- [ ] DNS records pointing to gateway
- [ ] Sufficient compute resources allocated
- [ ] Network security groups configured
- [ ] Backup strategy in place

### Monitoring

- [ ] Prometheus endpoint accessible
- [ ] Grafana instance running
- [ ] Alertmanager configured
- [ ] Log aggregation set up

### Security

- [ ] Secrets management configured (Vault, AWS Secrets Manager, etc.)
- [ ] RBAC policies defined
- [ ] Audit logging enabled
- [ ] TLS 1.3 enforced

### Testing

- [ ] Load testing completed
- [ ] Failover testing completed
- [ ] Security scanning completed
- [ ] ML model accuracy validated

---

## Environment Configuration

### Required Environment Variables

```bash
# Core Configuration
AEGISGATE_ENV=production
AEGISGATE_LOG_LEVEL=info
AEGISGATE_PORT=8443

# ML Configuration
ML_ENABLED=true
ML_SENSITIVITY=75
ML_ANOMALY_THRESHOLD=3.0
ML_PROMPT_INJECTION_ENABLED=true
ML_CONTENT_ANALYSIS_ENABLED=true
ML_BEHAVIORAL_ANALYSIS_ENABLED=true

# Prometheus Metrics
METRICS_ENABLED=true
METRICS_PORT=9090

# Upstream Configuration
UPSTREAM_API_URL=https://api.openai.com
UPSTREAM_TIMEOUT=30s

# Authentication
AUTH_ENABLED=true
AUTH_JWT_SECRET=${JWT_SECRET}
AUTH_OAUTH_ENABLED=true

# Security
TLS_CERT_PATH=/etc/aegisgate/certs/tls.crt
TLS_KEY_PATH=/etc/aegisgate/certs/tls.key
```

### Configuration File (aegisgate.yml)

```yaml
server:
  host: 0.0.0.0
  port: 8443
  read_timeout: 30s
  write_timeout: 30s

ml:
  enabled: true
  sensitivity: 75
  anomaly_threshold: 3.0
  prompt_injection:
    enabled: true
    sensitivity: 80
  content_analysis:
    enabled: true
    pii_detection: true
    secret_detection: true
  behavioral_analysis:
    enabled: true
    window_size: 5m
    threshold: 3.0

metrics:
  enabled: true
  port: 9090
  path: /metrics

proxy:
  upstream_url: https://api.openai.com
  timeout: 30s
  retry_count: 3

auth:
  enabled: true
  jwt_secret: ${JWT_SECRET}

logging:
  level: info
  format: json
  output: stdout
```

---

## Docker Deployment

### Quick Start

```bash
# 1. Pull the image
docker pull aegisgate/aegisgate:latest

# 2. Create configuration directory
mkdir -p /etc/aegisgate/config
mkdir -p /etc/aegisgate/certs

# 3. Copy configuration and certificates
cp aegisgate.yml /etc/aegisgate/config/
cp server.crt /etc/aegisgate/certs/
cp server.key /etc/aegisgate/certs/

# 4. Set environment variables
export JWT_SECRET="your-secure-secret-key"

# 5. Run the container
docker run -d \
  --name aegisgate-ml \
  -p 8443:8443 \
  -p 9090:9090 \
  -v /etc/aegisgate/config:/etc/aegisgate/config \
  -v /etc/aegisgate/certs:/etc/aegisgate/certs \
  -e ML_ENABLED=true \
  -e ML_SENSITIVITY=75 \
  -e JWT_SECRET \
  aegisgate/aegisgate:latest
```

### Docker Compose (Recommended)

```yaml
# docker-compose.ml.yml
version: '3.8'

services:
  aegisgate:
    image: aegisgate/aegisgate:latest
    container_name: aegisgate-ml
    ports:
      - "8443:8443"
      - "9090:9090"
    volumes:
      - ./config:/etc/aegisgate/config:ro
      - ./certs:/etc/aegisgate/certs:ro
    environment:
      - ML_ENABLED=true
      - ML_SENSITIVITY=75
      - METRICS_ENABLED=true
      - AEGISGATE_LOG_LEVEL=info
    env_file:
      - .env
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8443/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    networks:
      - aegisgate-net

  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    ports:
      - "9091:9090"
    volumes:
      - ./prometheus/prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - prometheus-data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
    restart: unless-stopped
    networks:
      - aegisgate-net

  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    volumes:
      - ./grafana/provisioning:/etc/grafana/provisioning:ro
      - grafana-data:/var/lib/grafana
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=${GRAFANA_PASSWORD}
      - GF_USERS_ALLOW_SIGN_UP=false
    restart: unless-stopped
    networks:
      - aegisgate-net

volumes:
  prometheus-data:
  grafana-data:

networks:
  aegisgate-net:
    driver: bridge
```

### Starting the Stack

```bash
# Start ML-enabled stack
docker-compose -f docker-compose.ml.yml up -d

# Verify services are running
docker-compose -f docker-compose.ml.yml ps

# View logs
docker-compose -f docker-compose.ml.yml logs -f aegisgate
```

---

## Kubernetes Deployment

### Helm Chart Installation

```bash
# Add Helm repository
helm repo add aegisgate https://aegisgatesecurity.io
helm repo update

# Install with custom values
helm install aegisgate-ml aegisgate/aegisgate-ml \
  --namespace aegisgate \
  --create-namespace \
  --values values-production.yaml
```

### Production Values File

```yaml
replicaCount: 3

image:
  repository: aegisgate/aegisgate
  tag: latest
  pullPolicy: Always

ml:
  enabled: true
  sensitivity: 75
  anomalyThreshold: 3.0
  promptInjection:
    enabled: true
    sensitivity: 80
  contentAnalysis:
    enabled: true
    piiDetection: true
    secretDetection: true
  behavioralAnalysis:
    enabled: true

resources:
  limits:
    cpu: 2000m
    memory: 2Gi
  requests:
    cpu: 1000m
    memory: 1Gi

autoscaling:
  enabled: true
  minReplicas: 3
  maxReplicas: 10
  targetCPUUtilizationPercentage: 70

ingress:
  enabled: true
  className: nginx
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: aegisgate.yourdomain.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: aegisgate-tls
      hosts:
        - aegisgate.yourdomain.com

prometheus:
  enabled: true
  serviceMonitor:
    enabled: true
    interval: 30s

persistence:
  enabled: true
  size: 50Gi

config:
  existingSecret: aegisgate-config

metrics:
  enabled: true
  port: 9090
```

### Verify Deployment

```bash
# Check pod status
kubectl get pods -n aegisgate

# Check ML service
kubectl get svc -n aegisgate

# View ML pod logs
kubectl logs -n aegisgate -l app=aegisgate-ml --tail=100

# Check ML metrics endpoint
kubectl port-forward -n aegisgate svc/aegisgate-ml 9090:9090
curl http://localhost:9090/metrics | grep ml_
```

---

## ML Model Configuration

### Sensitivity Levels

| Level | Score Range | Use Case |
|-------|-------------|----------|
| Low | 0-30 | Testing/Development |
| Medium | 31-60 | Standard Production |
| High | 61-85 | High Security |
| Maximum | 86-100 | Critical Infrastructure |

### Tuning Parameters

```yaml
ml:
  # Overall sensitivity (0-100)
  sensitivity: 75
  
  # Anomaly detection threshold (standard deviations)
  anomaly_threshold: 3.0
  
  prompt_injection:
    sensitivity: 80
    # Higher = more aggressive detection
    
  content_analysis:
    sensitivity: 70
    # PII detection aggressiveness
    
  behavioral_analysis:
    window_size: 5m
    threshold: 3.0
    # Request window and deviation threshold
```

### Pattern Configuration

The ML system includes these default detection patterns:

**Prompt Injection Patterns:**
- Direct instruction overrides (ignore previous, forget instructions)
- Role manipulation (roleplay, act as)
- Jailbreak attempts (DAN, developer mode)
- Prompt extraction attempts
- Hidden token injection
- Base64/Obfuscated content
- Context switching attacks

**Content Analysis:**
- PII detection (SSN, credit cards, emails, phones)
- API key / secret detection
- Private key detection

**Behavioral Analysis:**
- High frequency requests
- Path diversity anomalies
- Data volume anomalies

---

## Monitoring & Observability

### Prometheus Metrics

Key ML metrics to monitor:

```bash
# ML Request throughput
rate(ml_requests_total[5m])

# Anomaly detection rate
rate(ml_anomalies_total[5m])

# Prompt injection blocks
rate(ml_prompt_injection_blocked_total[5m])

# Content violation rate
rate(ml_content_violations_total[5m])

# ML analysis latency
histogram_quantile(0.95, rate(ml_analysis_duration_seconds_bucket[5m]))

# Block rate percentage
sum(rate(ml_blocked_total[5m])) / sum(rate(ml_requests_total[5m])) * 100
```

### Grafana Dashboard

Import the dashboard from: `deploy/docker/grafana/dashboards/aegisgate-ml.json`

**Dashboard Panels:**
- ML Request Throughput
- Total Anomalies Detected
- Prompt Injection Attempts
- Content Violations
- ML Analysis Latency (P50, P95, P99)
- Behavioral Analysis Metrics
- Block Rate Gauge
- Anomalies by Type

### Alert Rules

```yaml
groups:
  - name: aegisgate-ml-alerts
    interval: 30s
    rules:
      - alert: HighAnomalyRate
        expr: rate(ml_anomalies_total[5m]) > 10
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High anomaly detection rate detected"
          
      - alert: MLPipelineLatency
        expr: histogram_quantile(0.95, rate(ml_analysis_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "ML analysis latency above 1 second"
          
      - alert: HighBlockRate
        expr: (sum(rate(ml_blocked_total[5m])) / sum(rate(ml_requests_total[5m]))) > 0.15
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Block rate exceeds 15%"
```

---

## Security Considerations

### TLS Configuration

```yaml
tls:
  enabled: true
  min_version: "1.3"
  cipher_suites:
    - TLS_AES_256_GCM_SHA384
    - TLS_CHACHA20_POLY1305_SHA256
    - TLS_AES_128_GCM_SHA256
```

### Rate Limiting

```yaml
rate_limit:
  enabled: true
  requests_per_minute: 1000
  burst: 100
```

### API Authentication

```yaml
auth:
  enabled: true
  jwt:
    issuer: aegisgate
    expiry: 24h
  oauth:
    enabled: true
    providers:
      - google
      - microsoft
```

### Audit Logging

```yaml
audit:
  enabled: true
  log_all_requests: true
  log_ml_decisions: true
  storage:
    type: elasticsearch
    endpoint: https://elasticsearch:9200
```

---

## Troubleshooting

### Common Issues

#### 1. High False Positive Rate

**Symptoms:** Legitimate requests being blocked

**Solution:**
```bash
# Reduce sensitivity
kubectl set env deployment/aegisgate-ml ML_SENSITIVITY=50 -n aegisgate

# Review blocked requests
kubectl logs -n aegisgate -l app=aegisgate-ml | grep "BLOCKED"
```

#### 2. ML Analysis Latency High

**Symptoms:** Slow request processing

**Solution:**
```bash
# Check resource usage
kubectl top pods -n aegisgate

# Increase resources
kubectl patch deployment aegisgate-ml -n aegisgate -p '{"spec":{"template":{"spec":{"containers":[{"name":"aegisgate","resources":{"limits":{"cpu":"4000m","memory":"4Gi"}}}]}}}}'
```

#### 3. Prometheus Not Scraping ML Metrics

**Symptoms:** No ML metrics in Prometheus

**Solution:**
```bash
# Check service monitor
kubectl get servicemonitor -n aegisgate

# Check target status
curl http://prometheus:9090/api/v1/targets | jq '.data.activeTargets[] | select(.labels.app=="aegisgate-ml")'

# Verify metrics endpoint
kubectl exec -it -n aegisgate <pod-name> -- curl localhost:9090/metrics
```

#### 4. Memory Issues

**Symptoms:** OOMKills, instability

**Solution:**
```bash
# Check memory usage
kubectl top pods -n aegisgate

# Adjust behavioral analysis window
kubectl set env deployment/aegisgate-ml ML_BEHAVIORAL_WINDOW=10m -n aegisgate

# Limit client state cache
kubectl set env deployment/aegisgate-ml ML_MAX_CLIENT_STATES=10000 -n aegisgate
```

### Health Checks

```bash
# Basic health check
curl http://localhost:8443/health

# ML subsystem health
curl http://localhost:8443/health/ml

# Detailed status
curl http://localhost:8443/debug/status
```

### Debug Mode

```bash
# Enable debug logging
kubectl set env deployment/aegisgate-ml AEGISGATE_LOG_LEVEL=debug -n aegisgate

# View detailed ML decisions
kubectl logs -n aegisgate -l app=aegisgate-ml --tail=500 | grep "ML.*DEBUG"
```

---

## Incident Response

### Detection

1. **Monitor Alerts:** Watch for high anomaly rates
2. **Check Dashboard:** Review ML metrics in Grafana
3. **Review Logs:** Examine recent blocking decisions

### Response Steps

```bash
# 1. Scale up if needed
kubectl scale deployment aegisgate-ml --replicas=5 -n aegisgate

# 2. Increase blocking sensitivity
kubectl set env deployment/aegisgate-ml ML_SENSITIVITY=90 -n aegisgate

# 3. Enable strict mode
kubectl set env deployment/aegisgate-ml ML_STRICT_MODE=true -n aegisgate

# 4. Collect evidence
kubectl logs -n aegisgate -l app=aegisgate-ml --since=1h > incident-logs.txt

# 5. Export metrics
curl -o incident-metrics.json http://localhost:9090/api/v1/query_range?query=ml_requests_total&start=$(date -d '1 hour ago' +%s)&end=$(date +%s)&step=60
```

### Recovery

```bash
# 1. Return to normal sensitivity
kubectl set env deployment/aegisgate-ml ML_SENSITIVITY=75 -n aegisgate

# 2. Scale back
kubectl scale deployment aegisgate-ml --replicas=3 -n aegisgate

# 3. Disable strict mode
kubectl set env deployment/aegisgate-ml ML_STRICT_MODE=false -n aegisgate

# 4. Document incident
# Create incident report with timeline and findings
```

---

## Scaling Guidelines

### Horizontal Scaling

```bash
# Manual scaling
kubectl scale deployment aegisgate-ml --replicas=5 -n aegisgate

# Horizontal Pod Autoscaler
kubectl autoscale deployment aegisgate-ml -n aegisgate --cpu-percent=70 --min=3 --max=10
```

### Vertical Scaling

For high-throughput scenarios:

```yaml
resources:
  limits:
    cpu: 4000m
    memory: 4Gi
  requests:
    cpu: 2000m
    memory: 2Gi
```

### ML-Specific Scaling

The ML components are designed to scale horizontally:

- **Prompt Injection Detector:** Stateless, scales freely
- **Content Analyzer:** Stateless, scales freely
- **Behavioral Analyzer:** Uses distributed cache (Redis) for state

For large deployments, enable Redis:

```yaml
ml:
  behavioral_analysis:
    enabled: true
    use_redis: true
    redis_url: redis://redis-master:6379
```

---

## Maintenance

### Regular Tasks

| Task | Frequency | Command |
|------|-----------|---------|
| Log rotation | Daily | Automatic |
| Metrics cleanup | Weekly | Prometheus retention |
| Pattern updates | Monthly | Import new patterns |
| Model retraining | Quarterly | Training pipeline |
| Security patches | As needed | Rolling updates |

### Backup and Recovery

```bash
# Backup configuration
kubectl get configmap aegisgate-config -n aegisgate -o yaml > config-backup.yaml

# Backup ML patterns
kubectl get secret ml-patterns -n aegisgate -o yaml > patterns-backup.yaml

# Restore
kubectl apply -f config-backup.yaml -n aegisgate
```

### Updates

```bash
# Rolling update
kubectl rollout restart deployment/aegisgate-ml -n aegisgate

# Check update status
kubectl rollout status deployment/aegisgate-ml -n aegisgate

# Rollback if needed
kubectl rollout undo deployment/aegisgate-ml -n aegisgate
```

---

## Support

- **Documentation:** https://aegisgatesecurity.io
- **GitHub Issues:** https://github.com/aegisgate/aegisgate/issues
- **Slack:** #aegisgate-support
- **Email:** support@aegisgatesecurity.ioaegisgatesecurity.io

---

## Quick Reference

### Key Commands

```bash
# Deploy
helm install aegisgate-ml aegisgate/aegisgate-ml -n aegisgate

# Check status
kubectl get pods -n aegisgate -l app=aegisgate-ml

# View logs
kubectl logs -n aegisgate -l app=aegisgate-ml -f

# Scale
kubectl scale deployment aegisgate-ml --replicas=5 -n aegisgate

# Update ML sensitivity
kubectl set env deployment/aegisgate-ml ML_SENSITIVITY=80 -n aegisgate

# Access metrics
kubectl port-forward -n aegisgate svc/aegisgate-ml 9090:9090

# Health check
curl http://localhost:8443/health/ml
```

### Ports

| Port | Service | Description |
|------|---------|-------------|
| 8443 | HTTPS | Main API |
| 9090 | Metrics | Prometheus metrics |
| 8080 | gRPC | gRPC API |
| 9091 | Prometheus | Prometheus UI |

---

*Document Version: 1.0*
*Last Updated: 2024*
*AegisGate ML Security Gateway*
