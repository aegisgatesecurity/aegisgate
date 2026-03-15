# Configuration Reference

Complete configuration guide for AegisGate.

---

## Environment Variables

### Tier Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `AEGISGATE_TIER` | string | `community` | Tier: community, developer, professional, enterprise |
| `AEGISGATE_LICENSE_KEY` | string | - | License key for paid tiers |

### Server Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `AEGISGATE_HTTP_PORT` | int | 8080 | HTTP server port |
| `AEGISGATE_HTTPS_PORT` | int | 8443 | HTTPS server port |
| `AEGISGATE_HOST` | string | 0.0.0.0 | Host to bind to |
| `AEGISGATE_BASE_URL` | string | / | Base URL path |

### Storage Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `AEGISGATE_STORAGE_MODE` | string | file | Storage mode: file, postgres, mysql |
| `AEGISGATE_DATA_DIR` | string | ./data | Data directory for file storage |
| `DATABASE_URL` | string | - | Database connection URL |

### Security Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `AEGISGATE_TLS_ENABLED` | bool | false | Enable TLS |
| `AEGISGATE_TLS_CERT` | string | - | TLS certificate path |
| `AEGISGATE_TLS_KEY` | string | - | TLS private key path |
| `AEGISGATE_JWT_SECRET` | string | - | JWT signing secret |
| `AEGISGATE_API_KEYS_ENABLED` | bool | true | Enable API key authentication |

### Rate Limiting

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `AEGISGATE_RATE_LIMIT_ENABLED` | bool | true | Enable rate limiting |
| `AEGISGATE_RATE_LIMIT_REQUESTS` | int | 200 | Requests per minute |
| `AEGISGATE_RATE_LIMIT_BURST` | int | 50 | Burst requests allowed |

### Logging & Metrics

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `AEGISGATE_LOG_LEVEL` | string | info | Log level: debug, info, warn, error |
| `AEGISGATE_LOG_FORMAT` | string | json | Log format: json, text |
| `AEGISGATE_METRICS_ENABLED` | bool | true | Enable metrics |
| `AEGISGATE_METRICS_PORT` | int | 9090 | Metrics port |

### AI Provider Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `OPENAI_API_KEY` | string | - | OpenAI API key |
| `ANTHROPIC_API_KEY` | string | - | Anthropic API key |
| `COHERE_API_KEY` | string | - | Cohere API key |
| `AZURE_OPENAI_KEY` | string | - | Azure OpenAI API key |
| `AWS_ACCESS_KEY_ID` | string | - | AWS access key |
| `AWS_SECRET_ACCESS_KEY` | string | - | AWS secret key |

### Compliance Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `AEGISGATE_COMPLIANCE_ENABLED` | bool | true | Enable compliance |
| `AEGISGATE_COMPLIANCE_FRAMEWORKS` | string | owasp | Comma-separated frameworks |

---

## Configuration Files

### Community (.env)

```bash
# Tier
AEGISGATE_TIER=community

# Server
AEGISGATE_HTTP_PORT=8080
AEGISGATE_HTTPS_PORT=8443

# Storage
AEGISGATE_STORAGE_MODE=file
AEGISGATE_DATA_DIR=./data

# Security
AEGISGATE_TLS_ENABLED=false
AEGISGATE_JWT_SECRET=your-secret-key

# Rate Limiting
AEGISGATE_RATE_LIMIT_ENABLED=true
AEGISGATE_RATE_LIMIT_REQUESTS=200

# Logging
AEGISGATE_LOG_LEVEL=info
AEGISGATE_METRICS_ENABLED=true
```

### Developer (.env)

```bash
# Tier
AEGISGATE_TIER=developer
AEGISGATE_LICENSE_KEY=dev-xxxxxxxxxxxxx

# Server
AEGISGATE_HTTP_PORT=8080
AEGISGATE_HTTPS_PORT=8443

# Storage (PostgreSQL recommended)
AEGISGATE_STORAGE_MODE=postgres
DATABASE_URL=postgres://user:pass@localhost:5432/aegisgate?sslmode=require

# Security
AEGISGATE_TLS_ENABLED=true
AEGISGATE_TLS_CERT=/etc/aegisgate/certs/cert.pem
AEGISGATE_TLS_KEY=/etc/aegisgate/certs/key.pem
AEGISGATE_JWT_SECRET=your-secret-key
AEGISGATE_MTLS_ENABLED=true

# Rate Limiting
AEGISGATE_RATE_LIMIT_ENABLED=true
AEGISGATE_RATE_LIMIT_REQUESTS=1000

# AI Providers
OPENAI_API_KEY=sk-xxxxx
ANTHROPIC_API_KEY=sk-ant-xxxxx
COHERE_API_KEY=xxxxx

# Compliance
AEGISGATE_COMPLIANCE_ENABLED=true
AEGISGATE_COMPLIANCE_FRAMEWORKS=owasp,nist

# Monitoring
AEGISGATE_METRICS_ENABLED=true
GRAFANA_ENABLED=true
```

### Professional (.env)

```bash
# Tier
AEGISGATE_TIER=professional
AEGISGATE_LICENSE_KEY=pro-xxxxxxxxxxxxx

# Server
AEGISGATE_HTTP_PORT=8080
AEGISGATE_HTTPS_PORT=8443

# Storage (PostgreSQL + Redis)
AEGISGATE_STORAGE_MODE=postgres
DATABASE_URL=postgres://user:pass@postgres:5432/aegisgate?sslmode=require
REDIS_URL=redis://redis:6379

# Security
AEGISGATE_TLS_ENABLED=true
AEGISGATE_MTLS_ENABLED=true
AEGISGATE_PKI_ATTESTATION_ENABLED=true

# Multi-tenancy
AEGISGATE_MULTI_TENANT_ENABLED=true

# Compliance
AEGISGATE_COMPLIANCE_ENABLED=true
AEGISGATE_COMPLIANCE_FRAMEWORKS=owasp,soc2,gdpr,hipaa,pci,nist,iso27001
AEGISGATE_COMPLIANCE_STRICT_MODE=true

# SIEM
SIEM_ENABLED=true
SIEM_ENDPOINT=https://siem.company.com:443
```

---

## Command Line Flags

```bash
# Serve the API
aegisgate serve [flags]

# Flags:
--config string     Config file path (default "config.yaml")
--host string       Host to bind to (default "0.0.0.0")
--port int          Port to bind to (default 8080)
--tls               Enable TLS
--cert string       TLS certificate path
--key string        TLS key path

# Other commands
aegisgate version     Show version
aegisgate health      Check health
aegisgate migrate      Run database migrations
```

---

## Configuration Precedence

Configuration is loaded in this order (later overrides earlier):

1. Default values (hardcoded)
2. Configuration file
3. Environment variables
4. Command line flags

---

## Examples

### Minimal Configuration

```bash
AEGISGATE_TIER=community
```

### Full Configuration

```bash
# Server
AEGISGATE_TIER=developer
AEGISGATE_HTTP_PORT=8080
AEGISGATE_HTTPS_PORT=8443
AEGISGATE_HOST=0.0.0.0

# Security
AEGISGATE_TLS_ENABLED=true
AEGISGATE_TLS_CERT=/etc/aegisgate/certs/server.pem
AEGISGATE_TLS_KEY=/etc/aegisgate/certs/server.key
AEGISGATE_JWT_SECRET=super-secret-jwt-key

# Storage
AEGISGATE_STORAGE_MODE=postgres
DATABASE_URL=postgres://aegisgate:password@db:5432/aegisgate?sslmode=require
REDIS_URL=redis://redis:6379

# Rate Limiting
AEGISGATE_RATE_LIMIT_ENABLED=true
AEGISGATE_RATE_LIMIT_REQUESTS=1000
AEGISGATE_RATE_LIMIT_BURST=200

# Logging
AEGISGATE_LOG_LEVEL=debug
AEGISGATE_LOG_FORMAT=json
AEGISGATE_METRICS_ENABLED=true

# AI Providers
OPENAI_API_KEY=sk-xxxxx
ANTHROPIC_API_KEY=sk-ant-xxxxx

# Compliance
AEGISGATE_COMPLIANCE_ENABLED=true
AEGISGATE_COMPLIANCE_FRAMEWORKS=owasp,nist,soc2
```
