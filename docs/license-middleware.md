# License Validation Middleware

The license validation middleware provides enterprise-grade license checking by integrating with the AegisGate Admin Panel API.

## Overview

The middleware validates license keys against the Admin Panel, caches results, and provides fail-open functionality for high availability.

## Features

- **Remote Validation**: Validates license keys against Admin Panel API
- **Smart Caching**: 5-minute cache to reduce API calls
- **Fail-Open Design**: Continues operating if license service is down
- **Context Integration**: Tier and rate limits accessible in request context

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| ADMIN_PANEL_URL | URL of the Admin Panel API | http://localhost:8443 |
| LICENSE_KEY | License key to validate | (none - uses Community tier) |
| LICENSE_PUBLIC_KEY | Public key for signature verification | (none) |

## Usage

### Basic Integration

```go
import (
    "net/http"
    "github.com/aegisgatesecurity/aegisgate/pkg/middleware"
)

func main() {
    mux := http.NewServeMux()
    // Your handlers here
    
    // Wrap with license middleware
    handler := middleware.LicenseMiddleware()(mux)
    
    http.ListenAndServe(":8080", handler)
}
```

### Accessing License Info in Handlers

```go
func myHandler(w http.ResponseWriter, r *http.Request) {
    // Get the tier from the license
    tier := middleware.GetLicenseTierFromContext(r.Context())
    
    // Get the rate limit for this license
    rateLimit := middleware.GetLicenseRateLimitFromContext(r.Context())
    
    // Use the information
    w.Write([]byte(fmt.Sprintf("Tier: %s, Rate Limit: %d/min", tier, rateLimit)))
}
```

### Using the Validator Directly

```go
import "github.com/aegisgatesecurity/aegisgate/pkg/middleware"

func validate() {
    config := &middleware.LicenseConfig{
        AdminPanelURL: "https://admin.example.com",
        LicenseKey: "AG-xxxx-xxxx",
        FailOpen: true,
    }
    
    validator := middleware.NewLicenseValidator(config)
    result, err := validator.Validate(context.Background())
    
    if err != nil {
        // Handle error
    }
    
    if result.Valid {
        fmt.Printf("License valid for tier: %s\n", result.Tier)
    }
}
```

## ValidationResult

The `Validate` method returns a `ValidationResult`:

```go
type ValidationResult struct {
    Valid       bool          // Is the license valid?
    Status      string        // "valid", "expired", "revoked", "fail_open"
    Message     string        // Human-readable message
    Tier        core.Tier     // Tier level (0-3)
    ValidatedAt time.Time     // When validation occurred
    ExpiresAt   time.Time     // License expiration time
    MaxServers  int          // Maximum servers allowed
    MaxUsers    int          // Maximum users allowed
    RateLimit   int          // Requests per minute allowed
}
```

## Tier Levels

| Level | Tier | Description |
|-------|------|-------------|
| 0 | Community | Free tier - 1 server, 3 users, 60/min |
| 1 | Developer | $29/mo - 5 servers, 10 users, 600/min |
| 2 | Professional | $99/mo - 25 servers, 50 users, 3000/min |
| 3 | Enterprise | Custom - Unlimited |

## Error Handling

The middleware handles various error conditions:

- **No License Key**: Returns Community tier (fail-open)
- **License Service Down**: Returns fail-open if enabled
- **Expired License**: Returns error, blocks access
- **Revoked License**: Returns error, blocks access

## Performance

- Results are cached for 5 minutes by default
- Cache can be customized via `CacheDuration` in config
- Clear cache with `validator.ClearCache()`

## Security

- Uses HTTPS for Admin Panel communication (recommended)
- License keys are validated server-side
- Signature verification available via LICENSE_PUBLIC_KEY