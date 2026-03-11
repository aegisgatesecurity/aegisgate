# Admin Panel Documentation

## Overview

The AegisGate Admin Panel provides a comprehensive web-based interface for managing users, monitoring system activity, configuring authentication settings, and administering the PII Detection Platform. The admin interface is protected by role-based access control (RBAC) and is only accessible to users with the `admin` role.

### Features

- **User Management**: Create, edit, delete users and manage roles
- **Session Management**: View active sessions and revoke them as needed
- **Statistics Dashboard**: Real-time metrics and monitoring
- **Authentication Configuration**: Configure OAuth providers and security settings
- **Audit Logging**: View authentication and system events
- **Policy Management**: Manage PII detection policies (via API or UI)

---

## Accessing the Admin Panel

### URL and Login

- **Admin Panel URL**: `{AEGISGATE_URL}/admin`
- **Default Access**: Requires authentication with `admin` role
- **Session Requirements**: Valid session with sufficient privileges

### First-Time Setup

1. **Default Admin Creation** (Local Auth Mode):
   ```bash
   # Set environment variables before first run
   export DEFAULT_ADMIN_EMAIL="admin@yourdomain.com"
   export DEFAULT_ADMIN_PASSWORD="SecurePassword123!"
   export DEFAULT_ADMIN_ROLE="admin"
   ```

2. **OAuth First Admin**:
   - First user to authenticate when using OAuth mode gets admin privileges
   - Subsequent users get `user` role by default
   - Promote users: Access user management → Select user → Change role to `admin`

### Navigation

The admin panel consists of the following sections:

```
┌─────────────────────────────────────────────────────────────┐
│  AegisGate Admin Panel     [User] [Logout]           [Help]  │
├─────────────────────────────────────────────────────────────┤
│                                                            │
│  ┌──────────┐  ┌──────────────────────────────────────────┐ │
│  │ Dashboard│  │                                        │ │
│  │──────────│  │   Main Content Area                    │ │
│  │ Users    │  │                                        │ │
│  │ Sessions │  │   • User List                          │ │
│  │ Stats    │  │   • Session Monitor                    │ │
│  │ Audit    │  │   • Statistics                         │ │
│  │ Settings │  │   • Audit Logs                         │ │
│  └──────────┘  └──────────────────────────────────────────┘ │
│                                                            │
└─────────────────────────────────────────────────────────────┘
```

---

## User Management

### User List View

**Path**: `/admin/users`

The user list displays all registered users in the system with the following information:

| Column | Description | Sortable |
|--------|-------------|----------|
| Email | User's email address | ✅ |
| Name | Display name | ✅ |
| Provider | Authentication source (local/google/microsoft/etc.) | ✅ |
| Role | Current role (admin/user) | ✅ |
| Status | Account status (active/locked/pending) | ✅ |
| Last Login | Timestamp of last successful login | ✅ |
| Actions | Edit / Delete / Sessions | - |

**Features**:
- **Search**: Filter users by email or name
- **Pagination**: Configurable items per page (10/25/50/100)
- **Export**: Download user list as CSV or JSON
- **Bulk Actions**: Delete multiple users, change roles

### Creating a User

**Path**: `/admin/users/create`

**Form Fields**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| Email | Email | Yes | User's email address |
| Name | Text | Yes | Display name |
| Role | Select | Yes | `admin` or `user` |
| Provider | Select | Yes | `local` or specific OAuth provider |
| Password | Password | Conditional | Required if Provider = `local` |
| Send Welcome Email | Checkbox | No | Send credentials via email |

**Password Requirements** (Local Provider Only):
- Minimum length: 8 characters (configurable)
- At least one uppercase letter
- At least one lowercase letter
- At least one number
- Special characters optional (configurable)

**Example - Create Local User**:
```bash
# API endpoint
POST /api/v1/admin/users

# Request body
{
  "email": "newuser@example.com",
  "name": "New User",
  "role": "user",
  "provider": "local",
  "password": "SecurePass123!"
}
```

### Editing a User

**Path**: `/admin/users/{id}/edit`

**Editable Fields**:
- Name (all providers)
- Role (all providers)
- Password (local provider only)
- Status (active/locked)

**Actions Available**:
- Reset Password (local users)
- Send Password Reset Email
- Lock Account
- Delete Account

### User Detail View

**Path**: `/admin/users/{id}`

Displays comprehensive user information:

**Profile Section**:
- User ID
- Email
- Name
- Role
- Provider
- Created At
- Updated At
- Last Login

**Sessions Section**:
- Active sessions count
- Session list with IP, User Agent, Created At, Expires At
- Revoke individual sessions
- Revoke all sessions

**Activity Section**:
- Recent login activity
- Failed login attempts
- Actions performed (if audit logging enabled)

### Deleting Users

**Soft Delete vs Hard Delete**:
- **Soft Delete**: Disables account, preserves data, can be restored
- **Hard Delete**: Permanently removes user and all associated data

**Before Deleting**:
- Verify user has transferred ownership of any resources
- Consider impact on audit logs (retain user ID references)
- Revoke all active sessions

**Via Admin Panel**:
1. Navigate to user list
2. Click Delete icon
3. Confirm deletion type (soft/hard)
4. Enter confirmation text

**Via API**:
```bash
# Soft delete
DELETE /api/v1/admin/users/{id}?soft=true

# Hard delete
DELETE /api/v1/admin/users/{id}
```

---

## Session Management

### Active Sessions View

**Path**: `/admin/sessions`

The sessions view provides real-time monitoring of active authentication sessions:

| Column | Description |
|--------|-------------|
| User | Email and name |
| IP Address | Client IP (with geolocation if available) |
| User Agent | Browser and OS information |
| Created | Session start time |
| Expires | Session expiration time |
| Last Activity | Most recent request timestamp |
| Provider | Authentication source |

### Session Actions

**Revoke Session**:
- Terminates specific session
- User will be logged out on next request
- Audit log entry created

**Revoke All User Sessions**:
- Terminates all sessions for a specific user
- Forces re-authentication

**Global Session Management**:
- **Revoke All Sessions**: Emergency logout all users (admin only)
- **Session Timeout Override**: Adjust timeout for active sessions
- **Export Session Data**: For security analysis

### Session Statistics

**Overview Metrics**:
- Total active sessions
- Sessions per user (average, max)
- Sessions by provider
- Sessions by role
- Geolocation distribution
- Session duration statistics

---

## Statistics and Monitoring

### Dashboard Overview

**Path**: `/admin/dashboard`

The admin dashboard provides at-a-glance system health and usage metrics:

**Key Metrics**:
```
┌─────────────────────────────────────────────────────────────┐
│  Active Users:        42   │   Requests Today:    15,234   │
│  System Load:         23%  │   Avg Response:      45ms    │
│  Memory Usage:       156MB  │   Error Rate:        0.2%    │
│  Uptime:          15d 4h     │   Active Sessions:     38    │
├─────────────────────────────────────────────────────────────┤
│  [Authentication Activity Chart]                            │
│  - Logins (last 24h): 127                                   │
│  - Failed Attempts: 12                                      │
│  - New Users: 5                                             │
├─────────────────────────────────────────────────────────────┤
│  [Threat Detection]                                         │
│  - Suspicious IPs: 3 blocked                                │
│  - Lockout Events: 2                                        │
└─────────────────────────────────────────────────────────────┘
```

### Authentication Statistics

**Time Range Options**: Last Hour, Last 24 Hours, Last 7 Days, Last 30 Days, Custom

**Authentication Metrics**:
| Metric | Description |
|--------|-------------|
| Total Logins | Successful authentication attempts |
| Failed Logins | Failed authentication attempts |
| Login Success Rate | Percentage of successful logins |
| New User Registrations | First-time logins via OAuth |
| Unique Active Users | Distinct users with active sessions |
| Sessions Created | New session tokens issued |
| Sessions Terminated | Sessions ended (logout/timeout) |
| Average Session Duration | Mean time from login to logout |

**Charts**:
- Login activity over time (line chart)
- Authentication by provider (pie chart)
- Failed login attempts by reason (bar chart)
- Session duration distribution (histogram)

### System Health

**Memory Stats**:
- Current heap usage
- GC pause times
- Memory allocation rate
- Memory pool utilization

**Request Statistics**:
- Requests per second
- Average response time (P50, P95, P99)
- Error rate breakdown
- Top endpoints by volume

### Provider Statistics

Breakdown of authentication by provider:

| Provider | Users | Active Sessions | Avg Session Duration |
|----------|-------|-----------------|---------------------|
| Local | 15 | 8 | 4.2h |
| Google | 12 | 9 | 6.8h |
| Microsoft | 8 | 6 | 5.1h |
| GitHub | 4 | 3 | 3.9h |

---

## Audit Logging

### Audit Log View

**Path**: `/admin/audit`

The audit log provides a chronological record of system events:

**Filter Options**:
- Date range
- Event type
- User
- Severity (info/warning/error/critical)
- Source (web/API/system)

**Log Entry Structure**:
```json
{
  "timestamp": "2025-02-16T10:30:00Z",
  "event": "user.login.success",
  "user_id": "user-uuid",
  "user_email": "user@example.com",
  "ip_address": "192.168.1.100",
  "user_agent": "Mozilla/5.0...",
  "details": {
    "provider": "google",
    "session_id": "sess-uuid",
    "role": "admin"
  },
  "severity": "info",
  "source": "web"
}
```

### Event Types

#### Authentication Events
| Event | Description |
|-------|-------------|
| `user.login.success` | Successful login |
| `user.login.failure` | Failed login attempt |
| `user.logout` | User-initiated logout |
| `user.session.expired` | Session timeout |
| `user.session.revoked` | Admin or system revoked session |
| `user.password.changed` | Password updated |
| `user.password.reset` | Password reset initiated |

#### User Management Events
| Event | Description |
|-------|-------------|
| `user.created` | New user account created |
| `user.updated` | User details modified |
| `user.deleted` | User account deleted |
| `user.role.changed` | User role modified |
| `user.locked` | Account locked (failed attempts) |
| `user.unlocked` | Account unlocked |

#### System Events
| Event | Description |
|-------|-------------|
| `config.updated` | Configuration changed |
| `provider.enabled` | OAuth provider enabled |
| `provider.disabled` | OAuth provider disabled |

### Exporting Audit Logs

**Export Formats**: JSON Lines, CSV, Syslog

**Retention**: Configurable (default: 90 days for web display)

**API Endpoint**:
```bash
GET /api/v1/admin/audit/export
  ?start=2025-01-01T00:00:00Z
  &end=2025-02-16T23:59:59Z
  &format=json
```

---

## Settings Configuration

### Authentication Settings

**Path**: `/admin/settings/auth`

Configure authentication behavior:

| Setting | Default | Description |
|---------|---------|-------------|
| Session Timeout | 24h | Time before session expires |
| Max Session Lifetime | 7d | Absolute maximum session duration |
| Sliding Window | true | Extend session on activity |
| Secure Cookies | true | HTTPS-only cookies |
| CSRF Protection | true | Enable CSRF token validation |

### Password Policy

**Path**: `/admin/settings/password`

Configure local password requirements:

| Setting | Default | Description |
|---------|---------|-------------|
| Minimum Length | 8 | Minimum password characters |
| Require Uppercase | true | Require A-Z |
| Require Lowercase | true | Require a-z |
| Require Numbers | true | Require 0-9 |
| Require Special | false | Require special characters |
| Max Age | 0 | Force password change period (0=never) |
| History | 5 | Prevent reuse of last N passwords |

### Rate Limiting

**Path**: `/admin/settings/ratelimit`

Configure authentication rate limiting:

| Setting | Default | Description |
|---------|---------|-------------|
| Login Attempts | 5 | Failed attempts before lockout |
| Window | 5m | Time window for counting attempts |
| Lockout Duration | 15m | Account lockout period |
| IP Rate Limit | 20 | Requests per minute per IP |

### OAuth Provider Settings

**Path**: `/admin/settings/providers`

Enable/disable OAuth providers:

- Google
- Microsoft
- GitHub
- Okta
- Auth0

For each provider:
- Enable/Disable toggle
- Client ID display (masked)
- Client Secret update
- Redirect URI configuration
- Scope management

### SAML Configuration

**Path**: `/admin/settings/saml`

Configure SAML 2.0 identity provider:

- **IdP Metadata URL**: Provider discovery URL
- **Entity ID**: Service Provider entity identifier
- **Certificate**: SP certificate management
- **Attributes**: User attribute mapping
  - Email attribute
  - Name attribute
  - Role attribute (optional)

---

## Troubleshooting

### Admin Panel Access Issues

#### "Access Denied - Insufficient Privileges"
**Cause**: User does not have `admin` role
**Solution**:
1. Verify user role in user management
2. If locked out, use API to promote user:
   ```bash
   curl -X PATCH /api/v1/admin/users/{id}/role \
     -H "Authorization: Bearer {admin-token}" \
     -d '{"role": "admin"}'
   ```

#### "Session Expired"
**Cause**: Admin session timeout
**Solution**: Re-authenticate and ensure session has admin role

#### Blank Page / JavaScript Errors
**Cause**: Frontend build issues
**Solutions**:
- Clear browser cache
- Check browser console for errors
- Verify all static assets loaded (CSS/JS)

### User Management Issues

#### Cannot Create User
**Common Causes**:
- Email already in use
- Password does not meet requirements
- Invalid OAuth provider configured

**Password Validation Errors**:
Check `PASSWORD_MIN_LENGTH` and complexity settings:
```bash
export PASSWORD_MIN_LENGTH=8
export PASSWORD_REQUIRE_UPPERCASE=true
export PASSWORD_REQUIRE_LOWERCASE=true
export PASSWORD_REQUIRE_NUMBERS=true
```

#### User Not Appearing in List
**Causes**:
- Soft-deleted (showing active users only)
- Search filter active
- Pagination not showing

#### Cannot Delete User
**Causes**:
- User is the last admin
- User has active sessions
- Foreign key constraints

**Force Delete**:
```bash
# Revoke sessions first
DELETE /api/v1/admin/users/{id}/sessions

# Then delete user
DELETE /api/v1/admin/users/{id}
```

### Session Issues

#### Sessions Not Showing
**Cause**: Sessions are stateless by default
**Solution**: Enable session persistence:
```bash
export SESSION_STORE=redis  # or "memory" for in-memory
export REDIS_URL=redis://localhost:6379
```

#### Cannot Revoke Sessions
**Cause**: Session store unavailable
**Solution**: 
- Check Redis connection
- Verify session store configuration
- Use emergency logout via config reload

### Statistics Not Updating

**Cause 1**: Metrics disabled
```bash
export METRICS_ENABLED=true
```

**Cause 2**: Time range filter incorrect
**Solution**: Adjust dashboard time range

**Cause 3**: Data retention expired
**Solution**: Configure retention policies

---

## API Reference for Admin Operations

### Authentication

All admin API endpoints require:
- Valid authentication token
- `admin` role in session
- CSRF token for state-changing operations

### Admin API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/admin/stats` | System statistics |
| GET | `/api/v1/admin/stats/auth` | Authentication statistics |
| GET | `/api/v1/admin/users` | List all users |
| POST | `/api/v1/admin/users` | Create user |
| GET | `/api/v1/admin/users/{id}` | Get user details |
| PUT | `/api/v1/admin/users/{id}` | Update user |
| DELETE | `/api/v1/admin/users/{id}` | Delete user |
| PATCH | `/api/v1/admin/users/{id}/role` | Change user role |
| POST | `/api/v1/admin/users/{id}/reset-password` | Reset password |
| GET | `/api/v1/admin/users/{id}/sessions` | List user sessions |
| DELETE | `/api/v1/admin/users/{id}/sessions` | Revoke all user sessions |
| GET | `/api/v1/admin/sessions` | List all sessions |
| DELETE | `/api/v1/admin/sessions/{id}` | Revoke session |
| GET | `/api/v1/admin/audit` | Query audit log |
| GET | `/api/v1/admin/audit/export` | Export audit log |
| GET | `/api/v1/admin/settings` | Get all settings |
| PUT | `/api/v1/admin/settings` | Update settings |
| GET | `/api/v1/admin/settings/auth` | Get auth settings |
| PUT | `/api/v1/admin/settings/auth` | Update auth settings |

### Response Format

All admin API responses follow this structure:

```json
{
  "success": true,
  "data": { ... },
  "meta": {
    "page": 1,
    "per_page": 25,
    "total": 42
  }
}
```

Error responses:
```json
{
  "success": false,
  "error": {
    "code": "insufficient_privileges",
    "message": "Admin role required for this operation"
  }
}
```

---

## Best Practices

### User Account Management

1. **Principle of Least Privilege**: Grant minimum necessary access
2. **Regular Review**: Audit user roles quarterly
3. **Offboarding**: Disable accounts immediately on termination
4. **Shared Accounts**: Avoid shared admin accounts
5. **Service Accounts**: Use dedicated accounts for integrations

### Session Security

1. **Monitor Sessions**: Review active sessions weekly
2. **Geolocation Alerts**: Flag logins from unusual locations
3. **Concurrent Sessions**: Limit per user (configurable)
4. **Timeout Policy**: Balance security with usability

### Audit Practices

1. **Regular Review**: Weekly audit log analysis
2. **Alerting**: Set up automated alerts for suspicious activity
3. **Retention**: Maintain logs per compliance requirements
4. **Export**: Archive old logs regularly

### Emergency Procedures

**Account Lockout Emergency**:
```bash
# Disable authentication temporarily
curl -X PUT /api/v1/admin/settings/auth \
  -d '{"enabled": false}'
```

**Mass Session Revocation**:
```bash
# Revoke all non-admin sessions
curl -X DELETE /api/v1/admin/sessions?exclude_role=admin
```

**Configuration Backup**:
- Export settings before major changes
- Version control configuration
- Document changes in change log

---

## See Also

- [Authentication Guide](./AUTHENTICATION.md) - Authentication system details
- [Security Guide](./SECURITY.md) - Security configuration and best practices
- [Configuration Guide](./CONFIGURATION.md) - Complete configuration reference
- [Deployment Guide](./DEPLOYMENT_GUIDE.md) - Production deployment instructions
