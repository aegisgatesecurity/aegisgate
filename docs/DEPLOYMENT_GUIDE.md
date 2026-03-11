# Deployment Documentation

## Quick Start

### Docker Deployment

Build image:
  docker build -t aegisgate:latest .

Run container:  
  docker run -p 8080:8080 -p 8443:8443 aegisgate:latest

### Kubernetes Deployment

Apply manifests:
  kubectl apply -f kubernetes/

Check deployment:
  kubectl get pods -l app=aegisgate
  kubectl get svc aegisgate-service

## Configuration

AegisGate supports multiple configuration methods:

1. Environment Variables
2. YAML Configuration Files
3. Command-line Flags

## Monitoring

- Health endpoint: /health
- Metrics endpoint: /metrics
- Logs: Structured JSON format

## Security

- All traffic encrypted with TLS 1.3
- Mutual TLS support for service-to-service communication
- Rate limiting by default
- Security headers enforced

---

*For detailed operational procedures, see docs/OPERATIONAL_GUIDE.md*