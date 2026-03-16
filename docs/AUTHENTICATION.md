# Authentication System Documentation

## Overview

The AegisGate PII Detection Platform features a comprehensive authentication framework designed to provide secure, flexible, and scalable access control. The authentication system supports both local (username/password) and OAuth-based authentication providers, with role-based access control (RBAC) and secure session management.

### Key Features

- **Multiple Provider Support**: OAuth 2.0 (Google, Microsoft, GitHub), SAML 2.0, and Local authentication
- **Role-Based Access Control (RBAC)**: Granular permissions with `admin` and `user` roles
- **Secure Session Management**: Encrypted session cookies with configurable timeouts
- **PKCE Support**: Proof Key for Code Exchange for OAuth security
- **CSRF Protection**: Cross-Site Request Forgery prevention
- **Audit Logging**: Complete authentication event logging
- **Zero External Dependencies**: Self-contained implementation

---

## Architecture

### Authentication Flow

```
┌─────────────┐     ┌──────────────┐     ┌───────────────┐
│   Client    │────▶│   AegisGate    │────▶│  Auth Handler │
└─────────────┘     └──────────────┘     └───────┬───────┘
       │                                           │
       │                             ┌───────────┴───────────┐
       │                             ▼                       ▼
       │                    ┌──────────────┐  ┌───────────────┐
       │                    │     Local    │  │     OAuth     │
       │                    │   Provider   │  │   Provider    │
       │                    └──────┬───────┘  └───────┬───────┘
       │                           │                │
       │                           ▼                ▼
       │                    ┌──────────────┐  ┌───────────────┐
       └───────────────────▶│ User Store   │  │ Identity      │
                            │ (In-Memory)  │  │ Provider      │
                            └──────────────┘  └───────────────┘
```

### Components

| Component | Description | Location |
|-----------|-------------|----------|
| `auth.go` | Core authentication logic and user management | `pkg/auth/auth.go` |
| `oauth.go` | OAuth 2.0 provider implementations | `pkg/auth/oauth.go` |
| `session.go` | Session management and token handling | `pkg/auth/session.go` |
| `handlers.go` | HTTP handlers for authentication endpoints | `pkg/auth/handlers.go` |
| `middleware.go` | Authentication middleware for protected routes | `pkg/auth/middleware.go` |
| `local.go` | Local username/password authentication | `pkg/auth/local.go` |

---

## Supported Authentication Providers

### 1. OAuth 2.0 Providers

#### Google
- **Grant Type**: Authorization Code with PKCE
- **Scopes**: `openid`, `email`, `profile`
- **User Info Endpoint**: `https://openidconnect.googleapis.com/v1/userinfo`

#### Microsoft
- **Grant Type**: Authorization Code with PKCE
- **Scopes**: `openid`, `email`, `profile`
- **User Info Endpoint**: `https://graph.microsoft.com/v1.0/me`

#### GitHub
- **Grant Type**: Authorization Code
- **Scopes**: `read:user`, `user:email`
- **User Info Endpoint**: `https://api.github.com/user`

#### Okta
- **Grant Type**: Authorization Code with PKCE
- **Scopes**: `openid`, `email`, `profile`
- **Discovery**: Uses `.well-known/openid-configuration`

#### Auth0
- **Grant Type**: Authorization Code with PKCE
- **Scopes**: `openid`, `email`, `profile`
- **Discovery**: Uses `.well-known/openid-configuration`

### 2. SAML 2.0
- **Supported Bindings**: HTTP Redirect, HTTP POST
- **NameID Formats**: `emailAddress`, `unspecified`
- **Metadata**: Auto-discovery via IdP metadata URL
- **Signature**: SHA-256 with RSA

### 3. Local Authentication
- **Credential Storage**: In-memory with bcrypt hashing
- **Password Requirements**: Configurable minimum length and complexity
- **Rate Limiting**: Built-in brute force protection

---

## Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `AUTH_MODE` | Authentication mode (`local`, `oauth`, `saml`) | `local` | Yes |
| `SECRET_KEY` | Secret key for session encryption | - | Yes |
| `SESSION_TIMEOUT` | Session duration (e.g., `24h`, `30m`) | `24h` | No |
| `SESSION_SECURE` | Secure cookie flag (HTTPS only) | `true` | No |
| `COOKIE_NAME` | Session cookie name | `aegisgate_session` | No |
| `COOKIE_DOMAIN` | Cookie domain scope | `` | No |
| `COOKIE_PATH` | Cookie path scope | `/` | No |
| `CSRF_ENABLED` | Enable CSRF protection | `true` | No |
| `CSRF_TOKEN_NAME` | CSRF token cookie name | `aegisgate_csrf` | No |

### OAuth Configuration

| Variable | Description | Example |
|----------|-------------|---------|
| `GOOGLE_CLIENT_ID` | Google OAuth client ID | `xxx.apps.googleusercontent.com` |
| `GOOGLE_CLIENT_SECRET` | Google OAuth client secret | - |
| `MICROSOFT_CLIENT_ID` | Microsoft OAuth client ID | - |
| `MICROSOFT_CLIENT_SECRET` | Microsoft OAuth client secret | - |
| `GITHUB_CLIENT_ID` | GitHub OAuth client ID | - |
| `GITHUB_CLIENT_SECRET` | GitHub OAuth client secret | - |
| `OKTA_DOMAIN` | Okta organization domain | `dev-xxx.okta.com` |
| `OKTA_CLIENT_ID` | Okta OAuth client ID | - |
| `OKTA_CLIENT_SECRET` | Okta OAuth client secret | - |
| `AUTH0_DOMAIN` | Auth0 tenant domain | `xxx.auth0.com` |
| `AUTH0_CLIENT_ID` | Auth0 client ID | - |
| `AUTH0_CLIENT_SECRET` | Auth0 client secret | - |

### SAML Configuration

| Variable | Description | Example |
|----------|-------------|---------|
| `SAML_IDP_METADATA_URL` | IdP metadata URL | `https://idp.example.com/metadata` |
| `SAML_ENTITY_ID` | SP entity ID | `aegisgate-pii-detection` |
| `SAML_CERT_PATH` | SP certificate path | `./certs/saml.crt` |
| `SAML_KEY_PATH` | SP private key path | `./certs/saml.key` |

### Local Authentication Configuration

| Variable | Description | Default |
|----------|-------------|---------|
| `LOCAL_AUTH_ENABLED` | Enable local authentication | `true` |
| `PASSWORD_MIN_LENGTH` | Minimum password length | `8` |
| `PASSWORD_REQUIRE_UPPERCASE` | Require uppercase letters | `true` |
| `PASSWORD_REQUIRE_LOWERCASE` | Require lowercase letters | `true` |
| `PASSWORD_REQUIRE_NUMBERS` | Require numeric characters | `true` |
| `PASSWORD_REQUIRE_SPECIAL` | Require special characters | `false` |
| `MAX_LOGIN_ATTEMPTS` | Maximum failed login attempts | `5` |
| `LOCKOUT_DURATION` | Account lockout duration | `15m` |

---

## Role-Based Access Control (RBAC)

### User Roles

| Role | Description | Permissions |
|------|-------------|-------------|
| `admin` | Administrator with full access | All operations including user management, settings, and system configuration |
| `user` | Standard user | View data, create scans, view reports, cannot modify system settings |

### Permission Matrix

| Permission | Admin | User |
|------------|-------|------|
| View Dashboard | ✅ | ✅ |
| Run PII Scans | ✅ | ✅ |
| View Reports | ✅ | ✅ |
| Manage Users | ✅ | ❌ |
| Manage Policies | ✅ | ❌ |
| Configure Settings | ✅ | ❌ |
| View Audit Logs | ✅ | ❌ |
| Delete Data | ✅ | ❌ |

### User Properties

```go
type User struct {
    ID       string    // Unique identifier (UUID)
    Email    string    // User's email address
    Name     string    // Display name
    Provider string    // Authentication provider (local, google, etc.)
    Role     string    // User role (admin, user)
    Created  time.Time // Account creation timestamp
    Updated  time.Time // Last update timestamp
}
```

---

## Session Management

### Session Architecture

Sessions are implemented using encrypted HTTP cookies with the following properties:

- **Encryption**: AES-256-GCM with rotating keys
- **Storage**: Stateless (server-side session storage optional via `SESSION_STORE`)
- **Expiration**: Sliding window with configurable timeout
- **Security**: HttpOnly, Secure (configurable), SameSite=Lax

### Session Token Structure

```
session_cookie = base64(
    nonce(12 bytes) ||
    encrypted_data ||
    auth_tag(16 bytes)
)

decrypted_data = {
    user_id: string,
    email: string,
    role: string,
    provider: string,
    issued_at: timestamp,
    expires_at: timestamp,
    csrf_token: string
}
```

### Session Lifecycle

1. **Creation**: Upon successful authentication
2. **Validation**: On every request to protected endpoints
3. **Refresh**: Sliding window renewal on activity
4. **Termination**: On logout or timeout expiration

### Session Timeout Configuration

| Mode | Timeout Setting | Behavior |
|------|-----------------|----------|
| Idle Timeout | `SESSION_TIMEOUT` | Session expires after inactivity |
| Absolute Timeout | `SESSION_MAX_LIFETIME` | Maximum session duration regardless of activity |
| Sliding Window | `SESSION_REFRESH` | Token refreshed on activity before expiration |

---

## Security Features

### PKCE (Proof Key for Code Exchange)

All OAuth 2.0 flows include PKCE protection to prevent authorization code interception attacks:

```
1. Generate code_verifier (random 128 bytes, base64url)
2. Calculate code_challenge = SHA256(code_verifier)
3. Send code_challenge to authorization endpoint
4. Exchange code with code_verifier at token endpoint
5. Server validates code_challenge == SHA256(code_verifier)
```

### CSRF Protection

Cross-Site Request Forgery protection is implemented using double-submit cookie pattern:

1. **CSRF Token**: Generated per session or per request
2. **Cookie Storage**: `aegisgate_csrf` HttpOnly cookie
3. **Header Validation**: `X-CSRF-Token` header required for state-changing requests
4. **Origin Validation**: Strict origin checking for sensitive endpoints

### Secure Cookies

| Attribute | Setting | Description |
|-----------|---------|-------------|
| `HttpOnly` | `true` | Prevents JavaScript access |
| `Secure` | Configurable | HTTPS only transmission |
| `SameSite` | `Lax` | CSRF protection while allowing navigation |
| `Domain` | Configurable | Cookie scope restriction |
| `Path` | `/` | Cookie path scope |

### Header Security

Security headers set on all responses:

```
X-Content-Type-Options: nosniff
X-Frame-Options: DENY
X-XSS-Protection: 1; mode=block
Referrer-Policy: strict-origin-when-cross-origin
Content-Security-Policy: default-src 'self'
```

---

## Provider Setup Guides

### Google OAuth Setup

1. **Create Project**:
   - Go to [Google Cloud Console](https://console.cloud.google.com/)
   - Create a new project or select existing

2. **Enable API**:
   - Navigate to "APIs & Services" > "Credentials"
   - Click "Create Credentials" > "OAuth client ID"

3. **Configure OAuth Consent**:
   - Set application name and support email
   - Add scopes: `openid`, `email`, `profile`

4. **Create Credentials**:
   - Application type: "Web application"
   - Authorized redirect URIs: `https://your-domain.com/auth/callback?provider=google`

5. **Configure Environment**:
   ```bash
   export GOOGLE_CLIENT_ID="your-client-id.apps.googleusercontent.com"
   export GOOGLE_CLIENT_SECRET="your-client-secret"
   ```

### Microsoft OAuth Setup

1. **Register Application**:
   - Go to [Azure Portal](https://portal.azure.com/)
   - Navigate to "Azure Active Directory" > "App registrations"
   - Click "New registration"

2. **Configure Redirect URI**:
   - Platform: "Web"
   - Redirect URI: `https://your-domain.com/auth/callback?provider=microsoft`

3. **Create Secret**:
   - Go to "Certificates & secrets"
   - Add new client secret

4. **Configure Environment**:
   ```bash
   export MICROSOFT_CLIENT_ID="your-app-id"
   export MICROSOFT_CLIENT_SECRET="your-client-secret"
   export MICROSOFT_TENANT="common"  # or specific tenant ID
   ```

### GitHub OAuth Setup

1. **Register OAuth App**:
   - Go to GitHub > Settings > Developer Settings > OAuth Apps
   - Click "New OAuth App"

2. **Configure Application**:
   - Application name: "AegisGate PII Detection"
   - Homepage URL: `https://your-domain.com`
   - Authorization callback URL: `https://your-domain.com/auth/callback?provider=github`

3. **Get Credentials**:
   - Note Client ID and generate Client Secret

4. **Configure Environment**:
   ```bash
   export GITHUB_CLIENT_ID="your-client-id"
   export GITHUB_CLIENT_SECRET="your-client-secret"
   ```

### Okta Setup

1. **Create Application**:
   - Go to Okta Admin Console
   - Applications > Create App Integration
   - Select "OIDC - OpenID Connect" > "Web Application"

2. **Configure Sign-In**:
   - Sign-in redirect URIs: `https://your-domain.com/auth/callback?provider=okta`
   - Sign-out redirect URIs: `https://your-domain.com/auth/logout`

3. **Assign Users**:
   - Assign the application to groups or users

4. **Configure Environment**:
   ```bash
   export OKTA_DOMAIN="dev-xxx.okta.com"
   export OKTA_CLIENT_ID="your-client-id"
   export OKTA_CLIENT_SECRET="your-client-secret"
   ```

### Auth0 Setup

1. **Create Application**:
   - Go to Auth0 Dashboard
   - Applications > Create Application
   - Select "Regular Web Application"

2. **Configure URLs**:
   - Allowed Callback URLs: `https://your-domain.com/auth/callback?provider=auth0`
   - Allowed Logout URLs: `https://your-domain.com`
   - Allowed Web Origins: `https://your-domain.com`

3. **Configure Environment**:
   ```bash
   export AUTH0_DOMAIN="your-tenant.auth0.com"
   export AUTH0_CLIENT_ID="your-client-id"
   export AUTH0_CLIENT_SECRET="your-client-secret"
   ```

### SAML Configuration

1. **Configure Identity Provider**:
   - Provide SP metadata to IdP
   - Obtain IdP metadata URL

2. **Generate SP Credentials** (if not exist):
   ```bash
   openssl req -x509 -newkey rsa:2048 -keyout saml.key -out saml.crt -days 3650 -nodes
   ```

3. **Configure Environment**:
   ```bash
   export SAML_IDP_METADATA_URL="https://idp.example.com/metadata"
   export SAML_ENTITY_ID="aegisgate-pii-detection"
   export SAML_CERT_PATH="./certs/saml.crt"
   export SAML_KEY_PATH="./certs/saml.key"
   ```

---

## Example Configurations

### Development Configuration

```bash
# Authentication mode (local, oauth, saml)
export AUTH_MODE="local"

# Session configuration
export SECRET_KEY="dev-secret-key-minimum-32-characters"
export SESSION_TIMEOUT="24h"
export SESSION_SECURE="false"

# Local user (for development)
export DEFAULT_ADMIN_EMAIL="admin@localhost"
export DEFAULT_ADMIN_PASSWORD="devpassword123"
export DEFAULT_ADMIN_ROLE="admin"
```

### Production OAuth Configuration

```bash
# Use OAuth for authentication
export AUTH_MODE="oauth"

# Session configuration
export SECRET_KEY="$(openssl rand -base64 32)"
export SESSION_TIMEOUT="8h"
export SESSION_MAX_LIFETIME="24h"
export SESSION_SECURE="true"
export COOKIE_NAME="aegisgate_session"
export COOKIE_DOMAIN="your-domain.com"

# Google Provider
export GOOGLE_CLIENT_ID="xxx.apps.googleusercontent.com"
export GOOGLE_CLIENT_SECRET="xxx"

# Microsoft Provider (optional)
export MICROSOFT_CLIENT_ID="xxx"
export MICROSOFT_CLIENT_SECRET="xxx"

# GitHub Provider (optional)
export GITHUB_CLIENT_ID="xxx"
export GITHUB_CLIENT_SECRET="xxx"

# Security
export CSRF_ENABLED="true"
export MAX_LOGIN_ATTEMPTS="5"
export LOCKOUT_DURATION="15m"
```

### Production SAML Configuration

```bash
# Use SAML for enterprise SSO
export AUTH_MODE="saml"

# Session configuration
export SECRET_KEY="$(openssl rand -base64 32)"
export SESSION_TIMEOUT="8h"
export SESSION_SECURE="true"

# SAML configuration
export SAML_IDP_METADATA_URL="https://sso.company.com/metadata"
export SAML_ENTITY_ID="aegisgate-pii-detection"
export SAML_CERT_PATH="/etc/aegisgate/certs/saml.crt"
export SAML_KEY_PATH="/etc/aegisgate/certs/saml.key"

# Security
export CSRF_ENABLED="true"
```

### Docker Compose Configuration

```yaml
version: '3.8'
services:
  aegisgate:
    image: aegisgate:latest
    environment:
      # Core settings
      AEGISGATE_BIND_ADDRESS: ":8080"
      AEGISGATE_UPSTREAM: "http://backend:3000"

      # Authentication
      AUTH_MODE: "oauth"
      SECRET_KEY: "${SECRET_KEY}"
      SESSION_TIMEOUT: "8h"
      SESSION_SECURE: "true"
      COOKIE_DOMAIN: "aegisgate.example.com"

      # Google OAuth
      GOOGLE_CLIENT_ID: "${GOOGLE_CLIENT_ID}"
      GOOGLE_CLIENT_SECRET: "${GOOGLE_CLIENT_SECRET}"

      # Security
      CSRF_ENABLED: "true"
      MAX_LOGIN_ATTEMPTS: "5"
      LOCKOUT_DURATION: "15m"
    ports:
      - "8080:8080"
```

### Kubernetes Configuration

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: aegisgate-auth-secrets
type: Opaque
stringData:
  secret-key: "your-32-character-secret-key-here"
  google-client-id: "xxx.apps.googleusercontent.com"
  google-client-secret: "xxx"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: aegisgate
spec:
  replicas: 3
  selector:
    matchLabels:
      app: aegisgate
  template:
    metadata:
      labels:
        app: aegisgate
    spec:
      containers:
      - name: aegisgate
        image: aegisgate:latest
        env:
        - name: AUTH_MODE
          value: "oauth"
        - name: SECRET_KEY
          valueFrom:
            secretKeyRef:
              name: aegisgate-auth-secrets
              key: secret-key
        - name: GOOGLE_CLIENT_ID
          valueFrom:
            secretKeyRef:
              name: aegisgate-auth-secrets
              key: google-client-id
        - name: GOOGLE_CLIENT_SECRET
          valueFrom:
            secretKeyRef:
              name: aegisgate-auth-secrets
              key: google-client-secret
        - name: SESSION_TIMEOUT
          value: "8h"
        - name: SESSION_SECURE
          value: "true"
        - name: CSRF_ENABLED
          value: "true"
```

---

## API Endpoints

### Authentication Endpoints

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/login` | Local login | No |
| POST | `/api/v1/auth/logout` | Logout (all sessions) | Yes |
| GET | `/api/v1/auth/session` | Get current session info | Yes |
| GET | `/api/v1/auth/providers` | List configured providers | No |
| GET | `/auth/login/{provider}` | Initiate OAuth flow | No |
| GET | `/auth/callback` | OAuth callback handler | No |
| POST | `/api/v1/auth/refresh` | Refresh session token | Yes |

### User Management (Admin Only)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/admin/users` | List all users |
| GET | `/api/v1/admin/users/{id}` | Get user by ID |
| POST | `/api/v1/admin/users` | Create user |
| PUT | `/api/v1/admin/users/{id}` | Update user |
| DELETE | `/api/v1/admin/users/{id}` | Delete user |
| GET | `/api/v1/admin/users/{id}/sessions` | List user sessions |
| DELETE | `/api/v1/admin/users/{id}/sessions` | Revoke all user sessions |

---

## Troubleshooting

### Common Issues

#### Session Not Persisting
- Verify `SECRET_KEY` is at least 32 characters
- Check cookie settings (domain, path, secure)
- Ensure clock sync between client and server
- Check browser console for cookie warnings

#### OAuth Login Failures
- Verify redirect URI matches exactly (including protocol and port)
- Check OAuth credentials are correct
- Ensure provider application is enabled/active
- Review server logs for detailed error messages

#### "Invalid CSRF Token" Errors
- Verify CSRF cookie is being set
- Ensure client is sending `X-CSRF-Token` header
- Check for cross-origin request issues

#### SAML Authentication Errors
- Verify IdP metadata URL is accessible
- Check SP certificate and key match
- Validate clock synchronization (SAML is time-sensitive)
- Review SAML assertion for attribute mapping

#### Rate Limiting Lockout
- Default: 5 failed attempts within 5 minutes
- Lockout duration: 15 minutes
- Clear: Wait for timeout or use admin reset endpoint

### Debug Mode

Enable debug logging for authentication issues:

```bash
export LOG_LEVEL="debug"
export AUTH_DEBUG="true"
```

Debug mode logs:
- Authentication attempts (success/failure)
- Session creation/destruction
- Token validation details
- OAuth/SAML flow steps

---

## Security Best Practices

### Secret Management

1. **Generate Strong Secrets**:
   ```bash
   openssl rand -base64 32
   ```

2. **Rotate Secrets Regularly**:
   - SESSION_SECRET: Every 90 days
   - OAuth credentials: When personnel change
   - SAML certificates: Before expiry

3. **Store Secrets Securely**:
   - Use Kubernetes Secrets or similar
   - Never commit to version control
   - Use environment-specific values

### Production Checklist

- [ ] `SECRET_KEY` is at least 32 characters and randomly generated
- [ ] `SESSION_SECURE` is set to `true` (HTTPS only)
- [ ] `CSRF_ENABLED` is set to `true`
- [ ] OAuth redirect URIs use HTTPS
- [ ] SAML certificates are properly configured
- [ ] Password policies meet organizational requirements
- [ ] Rate limiting is enabled
- [ ] Session timeouts are configured appropriately
- [ ] CORS is properly restricted
- [ ] Security headers are enabled

### Monitoring

Recommended authentication metrics to monitor:

| Metric | Alert Threshold | Description |
|--------|-----------------|-------------|
| Failed Login Rate | > 10% | Indicates brute force or credential issues |
| Active Sessions | > 100/user | Potential account compromise |
| OAuth Failures | > 5% | Provider configuration issues |
| Session Duration | > 24h avg | Users not logging out |

---

## See Also

- [Admin UI Documentation](./ADMIN_UI.md) - Administrative interface guide
- [Security Guide](./SECURITY.md) - Security best practices and features
- [Configuration Guide](./CONFIGURATION.md) - Complete configuration reference
- [Deployment Guide](./DEPLOYMENT_GUIDE.md) - Production deployment instructions
