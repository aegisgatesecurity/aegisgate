# Troubleshooting Guide

Common issues and their solutions.

---

## Installation Issues

### "command not found: aegisgate" after build

**Problem:** After running `make build`, the `aegisgate` command is not found.

**Solution:**
```bash
# Check if binary exists
ls -la bin/

# Add to PATH
export PATH=$PATH:$(pwd)/bin

# Or use absolute path
./bin/aegisgate --help
```

### Docker build fails

**Problem:** Docker image build fails with errors.

**Solution:**
```bash
# Ensure Docker is running
docker info

# Build with verbose output
docker build -t aegisgate/aegisgate:latest -f deploy/docker/Dockerfile . --verbose

# Check available disk space
docker system df
```

---

## Configuration Issues

### "Invalid license key" error

**Problem:** AegisGate rejects the license key.

**Solutions:**
1. Verify the license key format (should start with `dev-`, `pro-`, or `ent-`)
2. Check for extra whitespace in the key
3. Ensure your tier supports the features you're using
4. Contact support if the issue persists

### "Database connection failed"

**Problem:** Cannot connect to PostgreSQL.

**Solutions:**
```bash
# Check DATABASE_URL format
# Should be: postgres://user:password@host:port/database?sslmode=require

# Test connection manually
psql $DATABASE_URL

# Check if PostgreSQL is running
docker ps | grep postgres

# Check network connectivity from container
docker exec -it aegisgate curl -v postgres:5432
```

---

## Runtime Issues

### High memory usage

**Problem:** AegisGate is using excessive memory.

**Solutions:**
1. Enable rate limiting to prevent abuse
2. Reduce log retention period
3. Use file storage instead of PostgreSQL for development
4. Check for memory leaks with:
   ```bash
   curl http://localhost:9090/debug/pprof/heap
   ```

### Slow performance

**Problem:** API requests are taking too long.

**Solutions:**
1. Enable caching (Professional tier):
   ```bash
   AEGISGATE_CACHE_ENABLED=true
   ```

2. Use PostgreSQL instead of file storage
3. Increase connection pool size:
   ```bash
   DATABASE_MAX_CONNECTIONS=100
   ```

4. Check network latency to AI providers

---

## Connection Issues

### Can't connect to AI provider

**Problem:** Requests to OpenAI/Anthropic are failing.

**Solutions:**
1. Verify API key is set:
   ```bash
   echo $OPENAI_API_KEY
   ```

2. Check network connectivity:
   ```bash
   docker exec aegisgate curl -v https://api.openai.com
   ```

3. Verify API key has sufficient credits
4. Check for rate limiting from provider

### "Too many requests" error

**Problem:** Rate limit exceeded.

**Solutions:**
1. Wait for the rate limit window to reset
2. Upgrade your tier for higher limits
3. Implement exponential backoff in your client
4. Check if someone else is using your API keys

---

## Security Issues

### TLS certificate errors

**Problem:** HTTPS/TLS is not working.

**Solutions:**
1. Verify certificate and key paths are correct
2. Check certificate is not expired:
   ```bash
   openssl x509 -in /path/to/cert.pem -dates -noout
   ```
3. Ensure certificate format is PEM
4. Check key permissions (should be 600)

### Authentication failures

**Problem:** API key or JWT validation failing.

**Solutions:**
1. Verify API key is correct and not expired
2. Check JWT secret matches across services
3. Ensure proper headers are being sent:
   ```
   Authorization: Bearer <token>
   X-API-Key: <api-key>
   ```
4. Check system time is accurate (JWT validation is time-sensitive)

---

## Compliance Issues

### Compliance checks failing

**Problem:** Compliance violations detected.

**Solutions:**
1. Review compliance report:
   ```bash
   curl http://localhost:8080/api/v1/compliance/report
   ```

2. Check which frameworks are enabled
3. Review specific violations and remediate
4. If false positive, adjust pattern sensitivity

---

## Logging & Diagnostics

### Where are the logs?

```bash
# Docker
docker logs aegisgate

# File storage
tail -f ./logs/aegisgate.log

# Systemd
journalctl -u aegisgate -f
```

### Enable debug logging

```bash
AEGISGATE_LOG_LEVEL=debug
```

### Health check

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

---

## Getting Help

### Collect diagnostic information

```bash
# Version info
./bin/aegisgate --version

# Configuration
./bin/aegisgate config --dump

# Runtime stats
curl http://localhost:9090/metrics
```

### Support channels

| Tier | Support |
|------|---------|
| Community | [GitHub Issues](https://github.com/aegisgatesecurity/aegisgate/issues) |
| Developer | Email: support@aegisgatesecurity.ioaegisgate.example.com |
| Professional | Priority email |
| Enterprise | 24/7 dedicated support |

---

## Known Issues

### Windows path separator

On Windows, some configuration paths may need double backslashes:
 Instead```bash
# of:
AEGISGATE_DATA_DIR=C:\aegisgate\data

# Use:
AEGISGATE_DATA_DIR=C:\\aegisgate\\data
# or
AEGISGATE_DATA_DIR=C:/aegisgate/data
```

### IPv6 localhost

Some systems may need explicit IPv4 binding:
```bash
AEGISGATE_HOST=127.0.0.1
```

---

*For issues not listed here, please open a GitHub issue or contact support.*
