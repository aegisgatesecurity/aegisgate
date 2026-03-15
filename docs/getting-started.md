# Getting Started with AegisGate

Welcome to AegisGate! This guide will help you get up and running quickly.

---

## Prerequisites

Before you begin, ensure you have the following:

| Requirement | Minimum | Recommended |
|-------------|---------|-------------|
| **Go** | 1.21 | 1.21+ |
| **Docker** | Latest | Latest |
| **PostgreSQL** | 14+ (optional) | 14+ |
| **RAM** | 2 GB | 4 GB+ |
| **Disk** | 10 GB | 50 GB+ |

---

## Quick Start (5 minutes)

### Option A: Docker Compose (Recommended)

```bash
# 1. Clone the repository
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate

# 2. Start AegisGate
docker-compose -f deploy/docker/docker-compose.yml up -d

# 3. Access the dashboard
# Open http://localhost:3000
```

### Option B: Build from Source

```bash
# 1. Clone the repository
git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate

# 2. Build AegisGate
make build

# 3. Run AegisGate
./bin/aegisgate serve --config ./config/community.env
```

---

## What's Next?

After installation, explore these guides:

| Guide | Description |
|-------|-------------|
| [Configuration](CONFIGURATION.md) | Configure AegisGate for your needs |
| [AI Providers](ai-providers.md) | Connect to OpenAI, Anthropic, etc. |
| [Compliance](compliance.md) | Set up compliance frameworks |
| [Security](security-hardening.md) | Harden your deployment |
| [Deployment](DEPLOYMENT.md) | Deploy to production |

---

## Basic Configuration

### Environment Variables

Create a `.env` file:

```bash
# Tier: community, developer, professional, enterprise
AEGISGATE_TIER=community

# Server
AEGISGATE_HTTP_PORT=8080
AEGISGATE_HTTPS_PORT=8443

# Storage
AEGISGATE_STORAGE_MODE=file
AEGISGATE_DATA_DIR=./data

# Security
AEGISGATE_TLS_ENABLED=false
```

### Starting AegisGate

```bash
# With environment file
./bin/aegisgate serve --config ./config/community.env

# With environment variables
AEGISGATE_TIER=community ./bin/aegisgate serve

# With Docker
docker run -p 8080:8080 -p 8443:8443 aegisgate/aegisgate:latest
```

---

## Verify Installation

### Health Check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "tier": "community"
}
```

### API Endpoint

```bash
curl http://localhost:8080/api/v1/status
```

---

## First Steps

### 1. Configure an AI Provider

```bash
# Add your OpenAI API key
export OPENAI_API_KEY=sk-xxxxx

# Or use configuration file
# Edit config/community.env and add:
# OPENAI_API_KEY=sk-xxxxx
```

### 2. Set Up a Proxy Route

```bash
# Create a route to OpenAI
curl -X POST http://localhost:8080/api/v1/routes \
  -H "Content-Type: application/json" \
  -d '{
    "path": "/openai",
    "target": "https://api.openai.com",
    "provider": "openai"
  }'
```

### 3. Test the Proxy

```bash
# Make a request through AegisGate
curl http://localhost:8080/openai/v1/chat/completions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer sk-xxxxx" \
  -d '{
    "model": "gpt-4",
    "messages": [{"role": "user", "content": "Hello!"}]
  }'
```

---

## Common Tasks

### Enable TLS/HTTPS

```bash
# Edit config
AEGISGATE_TLS_ENABLED=true
AEGISGATE_TLS_CERT=/path/to/cert.pem
AEGISGATE_TLS_KEY=/path/to/key.pem
```

### Set Up PostgreSQL

```bash
# config/professional.env
AEGISGATE_STORAGE_MODE=postgres
DATABASE_URL=postgres://user:pass@localhost:5432/aegisgate?sslmode=require
```

### Enable Metrics

```bash
# Metrics available at http://localhost:9090/metrics
AEGISGATE_METRICS_ENABLED=true
AEGISGATE_METRICS_PORT=9090
```

---

## Troubleshooting

### Container Won't Start

```bash
# Check logs
docker logs aegisgate

# Check port conflicts
netstat -tlnp | grep 8080
```

### Can't Connect to AI Provider

```bash
# Verify network connectivity
docker exec aegisgate curl https://api.openai.com

# Check API key
echo $OPENAI_API_KEY
```

### Performance Issues

```bash
# Check resource usage
docker stats aegisgate

# Review logs
docker logs aegisgate --tail 100
```

---

## Next Steps for Production

1. **Enable TLS** - Use valid certificates
2. **Configure Storage** - Use PostgreSQL
3. **Set Up Monitoring** - Enable metrics and logging
4. **Configure Rate Limits** - Prevent abuse
5. **Set Up Backup** - Regular database backups

See [Production Deployment](DEPLOYMENT.md) for details.

---

## Get Help

| Channel | Link |
|---------|------|
| Discord | https://discord.gg/aegisgate |
| Forum | https://community.aegisgate.example.com |
| GitHub Issues | https://github.com/aegisgatesecurity/aegisgate/issues |

---

*Happy securing! 🔒*
