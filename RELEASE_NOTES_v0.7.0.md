# Release Notes v0.7.0

## 🔒 AegisGate AI Security Gateway - v0.7.0

### Major Feature: HTTPS MITM Interception

This release introduces full HTTPS Man-in-the-Middle (MITM) interception capabilities, enabling comprehensive inspection of encrypted AI API traffic.

#### New Features

**MITM Proxy (pkg/proxy/mitm.go)**
- Full HTTPS interception with dynamic certificate generation
- On-the-fly certificate signing for any domain
- Certificate caching for improved performance
- TLS 1.3 support with modern cipher suites
- Connection pooling and management
- Integrated scanning for both requests and responses
- Configurable timeouts and connection limits
- Upstream proxy support for corporate environments

**Certificate Authority (pkg/tls/ca.go)**
- Dynamic Root CA generation
- Per-domain certificate generation
- Certificate caching with TTL
- Thread-safe operations

**Configuration Enhancements**
- `MaxConns` field added to config (default: 10000)
- Environment variables for MITM control:
  - `AEGISGATE_MITM_ENABLED` - Enable/disable MITM mode
  - `AEGISGATE_MITM_PORT` - MITM proxy port (default: 3128)
  - `AEGISGATE_MITM_SKIP_VERIFY` - Skip TLS verification
  - `AEGISGATE_MITM_UPSTREAM_PROXY` - Chain through corporate proxy

#### New Files
- `pkg/proxy/mitm.go` - MITM proxy implementation
- `pkg/proxy/mitm_test.go` - MITM proxy tests (13 test functions)
- `pkg/tls/ca.go` - Certificate Authority implementation
- `pkg/proxy/certs/ca/` - CA certificate storage

#### Modified Files
- `cmd/aegisgate/main.go` - MITM mode initialization
- `pkg/config/config.go` - Added MaxConns field and validation
- `pkg/config/config_test.go` - Updated tests for new config field

### Use Cases

The MITM capability enables:
1. **Full Traffic Visibility**: Inspect encrypted AI API requests/responses
2. **Comprehensive Security Scanning**: Detect threats in HTTPS traffic
3. **Compliance Enforcement**: Block sensitive data in encrypted streams
4. **Audit & Logging**: Complete audit trail for all AI interactions

### Configuration Example

```bash
# Enable MITM mode
export AEGISGATE_MITM_ENABLED=true
export AEGISGATE_MITM_PORT=3128

# Start AegisGate
./aegisgate
```

### Security Notice

⚠️ **Important**: The MITM CA private key must be secured:
- Set proper file permissions (600)
- Distribute only the public CA certificate to clients
- Never share the private key

### Upgrade Notes

No breaking changes. Existing configurations continue to work. To enable MITM:
1. Set `AEGISGATE_MITM_ENABLED=true`
2. Distribute `./certs/ca/ca.crt` to client trust stores
3. Configure clients to use proxy on port 3128

### Full Changelog

See [GitHub Releases](https://github.com/aegisgatesecurity/aegisgate/releases) for complete changelog.

---

**Previous Release**: v0.6.0 - ISO/IEC 42001 Framework Support
