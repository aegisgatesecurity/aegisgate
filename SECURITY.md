# Security Policy

---

## Reporting Security Vulnerabilities

We take security vulnerabilities seriously. If you discover a security issue, please report it responsibly.

### How to Report

**DO NOT** create a public GitHub issue for security vulnerabilities.

Instead, please report them via:

| Method | Contact |
|--------|---------|
| **Email** | security@aegisgatesecurity.io |
| **Encrypted PGP** | Available on our website |
| **Bug Bounty** | See aegisgatesecurity.io/bugbounty |

### What to Include

When reporting, please include:

1. **Description** - Clear description of the vulnerability
2. **Steps to Reproduce** - Detailed steps to reproduce the issue
3. **Impact** - Assessment of the vulnerability's impact
4. **Affected Versions** - Which versions are affected
5. **Suggested Fix** - If you have any suggestions (optional)

---

## Supported Versions

We actively support and provide security updates for the following versions:

| Version | Supported | Notes |
|---------|:---------:|-------|
| Latest | ✅ | Current stable release |
| v1.x | ✅ | Active support |
| v0.x | ❌ | End of life |

---

## Security Update Process

### Timeline

1. **Report Received** - We acknowledge receipt within 48 hours
2. **Initial Assessment** - We evaluate severity within 7 days
3. **Fix Development** - We develop and test the fix
4. **Release** - We release the security update
5. **Disclosure** - We publish security advisory

### Severity Levels

| Level | Response Time | Example |
|-------|--------------|---------|
| **Critical** | 24 hours | Remote code execution, data breach |
| **High** | 7 days | Privilege escalation, injection |
| **Medium** | 30 days | Information disclosure, bypass |
| **Low** | 90 days | Minor security improvements |

---

## Security Best Practices

### For Users

1. **Keep Updated** - Always run the latest version
2. **Secure Configuration** - Follow our hardening guide
3. **Access Control** - Limit who can access AegisGate
4. **Network Security** - Use TLS, firewall rules
5. **Monitoring** - Enable and review logs regularly

### For Operators

```bash
# Always use TLS in production
AEGISGATE_TLS_ENABLED=true
AEGISGATE_TLS_PORT=8443

# Enable audit logging
AEGISGATE_AUDIT_LOG=true

# Restrict admin access
AEGISGATE_ADMIN_IPS=10.0.0.0/8

# Enable rate limiting
AEGISGATE_RATE_LIMIT=true
```

---

## Security Features

AegisGate includes the following security features:

| Feature | Description |
|---------|-------------|
| **TLS Termination** | Secure all traffic with TLS 1.3 |
| **mTLS Support** | Mutual TLS for service-to-service |
| **PKI Attestation** | Hardware-based identity verification |
| **Secret Rotation** | Automatic credential rotation |
| **Audit Logging** | Comprehensive audit trail |
| **RBAC** | Role-based access control |
| **API Key Management** | Secure key generation and rotation |
| **Rate Limiting** | Prevent abuse and DoS |

---

## Compliance

AegisGate is designed to help you meet compliance requirements:

| Framework | How AegisGate Helps |
|-----------|-------------------|
| **SOC 2** | Audit logging, access control, encryption |
| **HIPAA** | Data encryption, audit trails, access control |
| **PCI-DSS** | TLS, secure defaults, logging |
| **GDPR** | Data encryption, access control, logging |
| **OWASP** | Built-in OWASP Top 10 protection |

---

## Security Research

We welcome security research contributions:

1. **Scope** - What's in scope for security testing
2. **Rules** - Guidelines for responsible disclosure
3. **Recognition** - Credit for valid findings
4. **Rewards** - Bug bounty program (see website)

### In Scope

- Source code on GitHub
- Documentation
- Public API endpoints

### Out of Scope

- Social engineering
- Physical security
- Denial of service (unless critical)

---

## Contact

For security-related inquiries:

- **Email**: security@aegisgatesecurity.io
- **PGP Key**: Available at https://aegisgatesecurity.io/security/pgp
- **Bug Bounty**: https://aegisgatesecurity.io/bugbounty

---

*Last Updated: 2026-03-11*