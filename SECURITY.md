# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.0.10  | Yes |
| 1.0.3   | Yes (legacy) |
| 1.0.2   | No (critical vulnerability) |
| 1.0.1   | No (critical vulnerability) |
| 1.0.0   | No (critical vulnerability) |
| <1.0.0 | No |

## Reporting a Vulnerability

Please DO NOT open public issues for security vulnerabilities.

Instead, contact us privately:

- **Email**: security@aegisgate.io
- **Subject**: Security Report - [Brief Description]
- **Response Time**: Within 24 hours
- **Fix Timeline**: 30-90 days depending on severity

### What to Include

1. Description of the vulnerability
2. Steps to reproduce
3. Affected versions
4. Potential impact
5. Suggested fix (optional)

### Responsible Disclosure

We follow responsible disclosure practices:
- Acknowledge receipt within 24 hours
- Provide timeline for fix within 7 days
- Notify you when fix is ready
- Request 90 days before public disclosure
- Credit you in release notes (unless anonymous preferred)

## Security Best Practices

### License Security
- Never commit license keys to version control
- Store HMAC secret in secure vault
- Rotate secrets regularly
- Monitor for unauthorized validation attempts

### Deployment Security
- Use mTLS for internal communication
- Enable FIPS 140-2 mode for regulated industries
- Deploy behind WAF/CDN
- Enable comprehensive audit logging

## Compliance

Security practices aligned with:
- SOC 2 Type II
- ISO 27001
- GDPR
- HIPAA
- PCI-DSS
