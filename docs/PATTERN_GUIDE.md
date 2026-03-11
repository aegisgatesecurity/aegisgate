# AegisGate Pattern Development Guide

## Overview

AegisGate uses regular expression patterns to detect sensitive data in HTTP traffic. 
Patterns are organized into categories with severity levels that determine whether 
the request/response should be blocked.

## Current Pattern Statistics

- **Default Patterns**: 26 patterns (in patterns.go)
- **Extended Patterns**: 18 patterns (in patterns_extended.go)  
- **Total Patterns**: 44 patterns
- **Categories**: 8 (PII, Credential, Financial, Cryptographic, Network, Healthcare, Cloud, Document)

## Pattern Structure

Each pattern follows this structure:

```go
{
    Name:        "PatternName",           // Unique identifier
    Regex:       regexp.MustCompile(`...`), // Go regex (RE2 syntax)
    Severity:    SeverityCritical,         // Info, Low, Medium, High, Critical
    Category:    CategoryCredential,       // One of the category constants
    Description: "Human-readable description",
},
```

## Severity Levels

| Level | Behavior | Use Case |
|-------|----------|----------|
| Info | Log only | Medical codes, non-sensitive IDs |
| Low | Log with context | Internal IPs, MAC addresses |
| Medium | Log with warning | Emails, JWT tokens |
| High | Log with alert | API keys, tokens |
| Critical | Block with HTTP 403 | Credit cards, passwords, private keys |

## Categories

1. **PII** - Personally Identifiable Information (SSN, email, phone, DOB)
2. **Credential** - Authentication tokens and API keys
3. **Financial** - Banking and payment data
4. **Cryptographic** - Private keys and certificates
5. **Network** - IP addresses and network identifiers
6. **Healthcare** - Medical data (NEW)
7. **Cloud** - Cloud provider IDs (NEW)
8. **Document** - Government IDs and vehicle info (NEW)

## How to Add a New Pattern

### Step 1: Choose the Right File

- **patterns.go**: Core patterns that ship with AegisGate
- **patterns_extended.go**: Additional patterns for enhanced coverage

### Step 2: Define the Pattern

Add to the appropriate file:

```go
// Example: Adding a new pattern for OAuth tokens
{
    Name:        "OAuthAccessToken",
    Regex:       regexp.MustCompile(`\b[A-Za-z0-9_]{40}`),
    Severity:    High,
    Category:    CategoryCredential,
    Description: "OAuth 2.0 access token detected",
},
```

### Step 3: Write the Regex

Go uses RE2 regex syntax (not PCRE). Key differences:
- NO lookaheads/lookbehinds: (?=...), (?!...)
- NO backreferences: \1, \2
- YES character classes: [a-z], \d, \w
- YES quantifiers: *, +, ?, {n,m}

**Tools for testing regex:**
- https://regex101.com (select Golang flavor)
- https://golang.org/pkg/regexp/syntax

### Step 4: Add Test Cases

Add to `pkg/scanner/pattern_test.go`:

```go
func TestOAuthAccessTokenPattern(t *testing.T) {
    scanner := New(nil)
    
    // Positive cases (should match)
    positiveCases := []string{
        "token: AbCdEfGhIjKlMnOpQrStUvWxYz1234567890AbCdEfGhIj",
        "Bearer AbCdEfGhIjKlMnOpQrStUvWxYz1234567890AbCdEfGhIj",
    }
    
    for _, tc := range positiveCases {
        findings := scanner.Scan(tc)
        found := false
        for _, f := range findings {
            if f.Pattern.Name == "OAuthAccessToken" {
                found = true
                break
            }
        }
        if !found {
            t.Errorf("OAuthAccessToken pattern should match: %s", tc)
        }
    }
    
    // Negative cases (should NOT match)
    negativeCases := []string{
        "token: short",  // Too short
        "no token here",
    }
    
    for _, tc := range negativeCases {
        findings := scanner.Scan(tc)
        for _, f := range findings {
            if f.Pattern.Name == "OAuthAccessToken" {
                t.Errorf("OAuthAccessToken pattern should NOT match: %s", tc)
            }
        }
    }
}
```

### Step 5: Test and Validate

```bash
# Build the scanner package
cd aegisgate && go build ./pkg/scanner

# Run tests
go test ./pkg/scanner -v

# Run specific pattern test
go test ./pkg/scanner -run TestOAuthAccessTokenPattern -v
```

## Pattern Examples by Type

### 1. Token Patterns (High/Critical Severity)

```go
// Generic API key
{
    Name:        "GenericAPIKey",
    Regex:       regexp.MustCompile(`(?i)\b(api[_-]?key|apikey)\s*[:=]\s*['']?[a-z0-9]{32,}['']?`),
    Severity:    High,
    Category:    CategoryCredential,
    Description: "Generic API key detected",
},

// JWT token
{
    Name:        "JWTToken",
    Regex:       regexp.MustCompile(`\beyJ[a-zA-Z0-9_-]*\.eyJ[a-zA-Z0-9_-]*\.[a-zA-Z0-9_-]*\b`),
    Severity:    Medium,
    Category:    CategoryCredential,
    Description: "JSON Web Token detected",
},
```

### 2. PII Patterns (Low-High Severity)

```go
// Credit Card (Critical - always block)
{
    Name:        "VisaCreditCard",
    Regex:       regexp.MustCompile(`\b4[0-9]{12}(?:[0-9]{3})?\b`),
    Severity:    Critical,
    Category:    CategoryFinancial,
    Description: "Visa credit card detected",
},

// Email (Low - log but don't block)
{
    Name:        "EmailAddress",
    Regex:       regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`),
    Severity:    Low,
    Category:    CategoryPII,
    Description: "Email address detected",
},
```

### 3. Cryptographic Patterns (Critical)

```go
// Private keys (always Critical)
{
    Name:        "RSAPrivateKey",
    Regex:       regexp.MustCompile(`-----BEGIN (?:RSA )?PRIVATE KEY-----`),
    Severity:    Critical,
    Category:    CategoryCryptographic,
    Description: "RSA private key detected",
},
```

## Best Practices

1. **Test regex at regex101.com first** (select Golang flavor)
2. **Use specific patterns** - avoid overly broad regex that causes false positives
3. **Include boundaries** - Use \b for word boundaries when matching standalone values
4. **Case insensitivity** - Use (?i) prefix for case-insensitive matching
5. **Test with real data** - Use actual samples of the data you want to detect
6. **Add negative tests** - Ensure pattern doesn't match benign text
7. **Document severity** - Choose severity based on business impact of exposure

## Performance Considerations

- Go's regex engine is fast but patterns are checked sequentially
- Very complex patterns (>100 chars) may slow scanning
- If adding many patterns, consider the extended patterns file approach
- Test performance with large payloads (>10KB)

## Troubleshooting

### Pattern compiles but doesn't match:
- Check for escape issues (backslashes need to be doubled in backtick strings)
- Verify regex doesn't use PCRE-only features
- Test regex in isolation

### Too many false positives:
- Make pattern more specific
- Add word boundaries (\b)
- Require context keywords (e.g., "password:" instead of just matching the password)

### Pattern causes build error:
- Ensure backticks are used for regex (`...`)
- Check for unescaped backslashes
- Verify all quotes are properly balanced

## Resources

- Go Regex Syntax: https://golang.org/pkg/regexp/syntax
- RE2 Syntax Wiki: https://github.com/google/re2/wiki/Syntax
- Regex101 (Golang): https://regex101.com/?flavor=golang

---

**Total Patterns Available**: 44 across 8 categories
**Last Updated**: Phase 2 - Pattern Enhancement

