# Security Guide

This document covers the security features, best practices, and configuration for the AegisGate PII Detection Platform.

---

## Overview

AegisGate implements defense-in-depth security across all layers of the application, from transport security to authentication, authorization, and audit logging.

### Security Principles

1. **Zero Trust**: Verify every request, regardless of source
2. **Defense in Depth**: Multiple security layers
3. **Least Privilege**: Minimum necessary permissions
4. **Secure by Default**: Secure configurations out of the box
5. **Zero External Dependencies**: Self-contained security implementation

---

## Authentication Security

### Session Security

**Encrypted Sessions**
- Algorithm: AES-256-GCM
- Key Derivation: HKDF with SHA-256
- Token Structure: Encrypted JSON Web Token format
- Cookie Security: HttpOnly, Secure (configurable), SameSite

**Session Token Structure**
```
base64(
  version(1 byte) ||
  nonce(12 bytes) ||
  ciphertext ||
  tag(16 bytes)
)

Decrypts to:
{
  "user_id": "uuid",
  "email": "user@example.com",
  "role": "user",
  "provider": "google",
  "issued_at": 1708089600,
  "expires_at": 1708176000,
  "csrf_token": "random-token"
}
```

**Session Security Features**

| Feature | Implementation | Purpose |
|---------|---------------|---------|
| Encryption | AES-256-GCM | Confidentiality |
| Nonce | Random 12 bytes | Unique encryption |
| Expiration | Unix timestamps | Time-bound sessions |
| Role Binding | Embedded role | Authorization context |
| CSRF Token | Per-session | CSRF protection |

**Session Configuration**

| Variable | Recommended | Description |
|----------|-------------|-------------|
| `SESSION_TIMEOUT` | 8 hours | Active session duration |
| `SESSION_MAX_LIFETIME` | 24 hours | Absolute maximum session |
| `SESSION_SECURE` | true | HTTPS-only cookies |
| `COOKIE_SAMESITE` | Lax | CSRF protection |
| `COOKIE_HTTPONLY` | true | True* | JavaScript cannot access |

*Always true, not configurable

### PKCE (Proof Key for Code Exchange)

All OAuth flows include PKCE protection:

```
1. Client generates: code_verifier = random(128 bytes)
2. Client sends: code_challenge = SHA256(code_verifier)
3. Server stores: code_challenge with authorization code
4. Client proves: sends code_verifier with code
5. Server validates: SHA256(code_verifier) == code_challenge
```

**Benefits:**
- Prevents authorization code interception
- Mitigates man-in-the-browser attacks
- Required for OAuth 2.0 public clients

### CSRF Protection

Cross-Site Request Forgery protection uses double-submit cookie pattern:

**Mechanism:**
1. Server sets CSRF token in HttpOnly cookie
2. Client reads token and sends in header
3. Server validates token matches

**CSRF Token Properties:**
- 32-byte random value (base64url encoded)
- Same expiration as session
- HttpOnly cookie (not accessible to JavaScript)
- Validated on state-changing requests

**CSRF-Protected Operations:**
- POST, PUT, DELETE requests
- Admin configuration changes
- User management operations
- Authentication endpoints

**Configuration:**
```bash
export CSRF_ENABLED="true"
export CSRF_TOKEN_NAME="aegisgate_csrf"
export CSRF_HEADER="X-CSRF-Token"
```

### Password Security

**Local Authentication**

| Security Measure | Implementation |
|------------------|----------------|
| Hashing | bcrypt (cost factor 12) |
| Salt | Random 16-byte per password |
| Complexity | Configurable requirements |
| Lockout | Rate limiting on failed attempts |
| Reset | Secure token flow |

**bcrypt Details:**
- Algorithm: bcrypt with SHA-256
- Cost: 12 (configurable, minimum 10)
- Salt: 128-bit random
- Output: Modular crypt format

**Password Requirements** (configurable):
```bash
export PASSWORD_MIN_LENGTH="12"
export PASSWORD_REQUIRE_UPPERCASE="true"
export PASSWORD_REQUIRE_LOWERCASE="true"
export PASSWORD_REQUIRE_NUMBERS="true"
export PASSWORD_REQUIRE_SPECIAL="true"
```

**Brute Force Protection:**
```bash
export MAX_LOGIN_ATTEMPTS="5"      # Failed attempts
export LOCKOUT_DURATION="15m"       # Lockout period
export MAX_LOGIN_WINDOW="5m"         # Counting window
```

---

## Authorization Security

### Role-Based Access Control (RBAC)

**Roles:**
| Role | Privileges | Use Case |
|------|------------|----------|
| `admin` | Full access | System administrators |
| `user` | Limited access | Standard users |

**Permission Enforcement:**
- Middleware checks role on protected routes
- API validates role for each operation
- UI elements conditionally rendered
- Backend double-validates all permissions

**Role Storage:**
- Embedded in session token
- Stored in user record (for reference)
- Versioned for audit trail

### API Security

**Authentication Required Endpoints:**
- All `/api/v1/admin/*` routes
- User profile endpoints
- Session management endpoints
- Configuration endpoints

**Public Endpoints:**
- `/health` - Health check
- `/metrics` - Prometheus metrics (if enabled)
- `/auth/providers` - Auth provider list
- `/auth/login/*` - OAuth initiation
- `/auth/callback` - OAuth callback

**Rate Limiting:**
```bash
export AEGISGATE_RATE_LIMIT="100"  # Requests per minute per IP
```

Rate limits by endpoint:
- Login: 10 attempts per minute
- API: Configurable per endpoint
- Admin: 50 requests per minute

---

## Transport Security

### TLS/SSL Configuration

**Recommended TLS Settings:**

```bash
export AEGISGATE_TLS_ENABLED="true"
export AEGISGATE_TLS_CERT_FILE="/etc/aegisgate/certs/server.crt"
export AEGISGATE_TLS_KEY_FILE="/etc/aegisgate/certs/server.key"
export AEGISGATE_TLS_MIN_VERSION="1.2"  # Or 1.3
```

**TLS 1.2 Configuration:**
- Cipher suites: Strong only (no RC4, 3DES)
- Perfect Forward Secrecy: Required
- Certificate validation: Strict

**TLS 1.3 Configuration:**
- Recommended for new deployments
- Simplified handshake
- Improved security

**Certificate Requirements:**
- RSA 2048-bit+ or ECC P-256+
- SHA-256 signature algorithm
- Valid certificate chain
- Not expired

### Security Headers

AegisGate sets security headers on all responses:

| Header | Value | Purpose |
|--------|-------|---------|
| `X-Content-Type-Options` | `nosniff` | Prevent MIME sniffing |
| `X-Frame-Options` | `DENY` | Prevent clickjacking |
| `X-XSS-Protection` | `1; mode=block` | XSS protection |
| `Referrer-Policy` | `strict-origin-when-cross-origin` | Referrer privacy |
| `Content-Security-Policy` | `default-src 'self'` | XSS prevention |

**CSP Directives:**
```
default-src 'self';
script-src 'self' 'unsafe-inline';
style-src 'self' 'unsafe-inline';
img-src 'self' data:;
connect-src 'self';
frame-ancestors 'none';
base-uri 'self';
form-action 'self';
```

### CORS Configuration

Cross-Origin Resource Sharing configuration:

```bash
export CORS_ALLOWED_ORIGINS="https://app.example.com"
export CORS_ALLOWED_METHODS="GET,POST,PUT,DELETE,OPTIONS"
export CORS_ALLOWED_HEADERS="Content-Type,Authorization,X-CSRF-Token"
export CORS_ALLOW_CREDENTIALS="true"
export CORS_MAX_AGE="3600"
```

**Security Notes:**
- Never use `*` with credentials
- Explicitly list allowed origins
- Limit allowed methods
- Validate headers server-side

---

## Audit Logging

### Authentication Events

All authentication events are logged:

| Event | Severity | Details |
|-------|----------|---------|
| `user.login.success` | Info | User ID, IP, provider |
| `user.login.failure` | Warning | Email, IP, reason |
| `user.logout` | Info | User ID, IP |
| `user.session.expired` | Info | User ID, duration |
| `user.session.revoked` | Warning | User ID, admin ID, reason |
| `user.password.changed` | Info | User ID |
| `user.password.reset` | Info | User ID, token expiry |
| `user.locked` | Warning | User ID, attempts |
| `user.unlocked` | Info | User ID, admin ID |

### Event Structure

```json
{
  "timestamp": "2025-02-16T10:30:00Z",
  "level": "info",
  "event": "user.login.success",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "user_email": "user@example.com",
  "ip_address": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "details": {
    "provider": "google",
    "session_id": "sess-xxx",
    "role": "admin",
    "duration_ms": 245
  }
}
```

### Audit Storage

**Retention Policies:**
- Authentication events: 90 days
- Security events: 1 year
- System events: 30 days

**Export Formats:**
- JSON Lines
- CSV
- Syslog (for SIEM integration)

**Compliance:**
- GDPR: Right to access/erasure
- SOC 2: Audit trail requirements
- HIPAA: Access logging

---

## Secret Management

### Required Secrets

| Secret | Generation | Rotation |
|--------|------------|----------|
| `SECRET_KEY` | `openssl rand -base64 32` | 90 days |
| `TLS private key` | RSA 2048+ or ECC P-256 | Before expiry |
| OAuth client secrets | Provider UI | When compromised |
| SAML private key | RSA 2048+ | Before expiry |

### Secure Generation

```bash
# Session secret (32 bytes = 44 base64 chars)
SECRET_KEY="$(openssl rand -base64 32)"

# TLS certificate
openssl req -x509 -newkey rsa:4096 \
  -keyout server.key -out server.crt \
  -days 365 -nodes

# Strong encryption for backups
openssl enc -aes-256-cbc -salt -in secrets.env -out secrets.env.enc
```

### Storage Best Practices

**Kubernetes:**
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: aegisgate-secrets
type: Opaque
data:
  SECRET_KEY: $(echo "secret" | base64)
  GOOGLE_CLIENT_SECRET: $(echo "secret" | base64)
```

**Docker:**
```bash
docker secret create aegisgate_secret_key <(openssl rand -base64 32)
```

**Environment Files:**
```bash
# .env.production (add to .gitignore!)
SECRET_KEY="xxx"
source .env.production
```

### Never Commit Secrets

```bash
# .gitignore
echo ".env" >> .gitignore
echo ".env.*" >> .gitignore
echo "*.pem" >> .gitignore
echo "*.key" >> .gitignore
```

---

## Input Sanitization

### Request Validation

**Size Limits:**
```bash
export AEGISGATE_MAX_BODY_SIZE="10485760"  # 10MB
export AEGISGATE_MAX_HEADER_SIZE="1048576"  # 1MB
```

**Validation Layers:**
1. HTTP server rejects oversized requests
2. Middleware validates content type
3. Handlers validate JSON structure
4. Business logic validates field values

### PII Detection

**Content Scanning:**
- Regex pattern matching
- Heuristic analysis
- Machine learning classification
- Configurable sensitivity

**Pattern Categories:**
| Category | Patterns | Examples |
|----------|----------|----------|
| Financial | Credit cards, IBAN, SSN | 4111-1111-1111-1111 |
| Personal | Name, DOB, Address | John Doe, 01/01/1990 |
| Contact | Email, Phone | user@example.com |
| Healthcare | Medical record numbers | Patient ID patterns |

**Redaction:**
- Configurable replacement characters
- Full or partial redaction
- Audit trail of redacted content

---

## Security Hardening

### Production Checklist

```markdown
SSL/TLS:
- [ ] TLS 1.2+ only
- [ ] Valid certificates
- [ ] HSTS enabled
- [ ] TLS session tickets disabled

Authentication:
- [ ] SECRET_KEY >= 32 characters
- [ ] SESSION_SECURE=true
- [ ] CSRF_ENABLED=true
- [ ] Rate limiting configured
- [ ] Password policy enforced

Session:
- [ ] Reasonable timeout (8h)
- [ ] Max lifetime configured
- [ ] Sliding window enabled
- [ ] HttpOnly cookies
- [ ] SameSite=Lax or Strict

Authorization:
- [ ] Roles enforced
- [ ] Principle of least privilege
- [ ] Admin access limited

Audit:
- [ ] Logging enabled
- [ ] Retention configured
- [ ] Regular log review

Secrets:
- [ ] Secrets not in code
- [ ] Regular rotation
- [ ] Secure storage

General:
- [ ] Security headers set
- [ ] CORS configured properly
- [ ] Dependency scanning
- [ ] Container hardening
```

### Container Security

**Dockerfile Best Practices:**
```dockerfile
# Use distroless image
FROM gcr.io/distroless/static-debian11

# Run as non-root
USER nonroot:nonroot

# Copy only needed files
COPY --chown=nonroot:nonroot ./build/aegisgate /app/

# Read-only root filesystem
readOnlyRootFilesystem: true

# No privileged escalation
allowPrivilegeEscalation: false
```

**Kubernetes Security:**
```yaml
securityContext:
  runAsNonRoot: true
  runAsUser: 65534
  readOnlyRootFilesystem: true
  allowPrivilegeEscalation: false
  capabilities:
    drop:
      - ALL
```

---

## Vulnerability Management

### Dependency Scanning

AegisGate has zero external dependencies, but dependencies may exist in:
- Build tools
- CI/CD pipeline
- Deployment environment

**Scanning Commands:**
```bash
# Go vulnerability database
govulncheck ./...

# Dependency check (if any)
go mod tidy
go list -m all
```

### Security Headers Verification

Test security headers:
```bash
curl -I https://aegisgate.example.com | grep -E "(X-Frame|X-Content|X-XSS|Content-Security)"
```

Expected output:
```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Content-Security-Policy: default-src 'self'
```

### Penetration Testing

**Areas to Test:**
- Authentication bypass
- Session fixation
- CSRF vulnerabilities
- XSS injections
- SQL injection (if database used)
- Path traversal
- IDOR (Insecure Direct Object Reference)

**Tools:**
- OWASP ZAP
- Burp Suite
- Nikto
- Go test fuzzing

---

## Incident Response

### Security Incident Types

| Severity | Examples | Response Time |
|----------|----------|---------------|
| Critical | Data breach, RCE | Immediate |
| High | Auth bypass, Privilege escalation | 1 hour |
| Medium | XSS, CSRF | 24 hours |
| Low | Info disclosure | 7 days |

### Response Procedures

**Authentication Breach:**
1. Revoke all sessions
2. Rotate SECRET_KEY
3. Force password resets
4. Audit user activities
5. Notify affected users

**Session Hijacking:**
1. Identify compromised sessions
2. Revoke specific sessions
3. Review audit logs
4. Implement IP restrictions
5. Monitor for repeat attacks

---

## Compliance

### GDPR Compliance

**Data Minimization:**
- Only collect necessary data
- Retain minimum time required
- Implement data purging

**Right to Erasure:**
```bash
# Delete user and associated data
DELETE /api/v1/admin/users/{id}?hard=true
```

**Data Portability:**
- Export user data in JSON
- Machine-readable format

### SOC 2

**Trust Service Criteria:**
- Security: Access controls, encryption
- Availability: Uptime, monitoring
- Processing Integrity: Application security
- Confidentiality: Access controls
- Privacy: Data handling practices

### HIPAA

**If handling PHI:**
- Access controls
- Audit logging
- Encryption at rest and in transit
- Automatic logoff
- Data integrity controls

---

## Security References

### OWASP Top 10

| OWASP Risk | Mitigation |
|------------|------------|
| A01: Broken Access Control | RBAC, middleware checks |
| A02: Cryptographic Failures | AES-256, bcrypt, HTTPS |
| A03: Injection | Pattern scanning, input validation |
| A05: Security Misconfiguration | Secure defaults, hardening guide |
| A07: ID and Auth Failures | OAuth, PKCE, rate limiting |
| A08: Software and Data Failures | Zero dependencies, no secrets in code |

### Security Resources

- [OWASP Cheat Sheets](https://cheatsheetseries.owasp.org/)
- [CSP Quick Reference](https://content-security-policy.com/)
- [OWASP Authentication Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)

---

## See Also

- [Authentication Guide](./AUTHENTICATION.md) - Authentication system details
- [Admin UI Guide](./ADMIN_UI.md) - Admin panel documentation
- [Configuration Guide](./CONFIGURATION.md) - Configuration reference
- [Deployment Guide](./DEPLOYMENT_GUIDE.md) - Production deployment
