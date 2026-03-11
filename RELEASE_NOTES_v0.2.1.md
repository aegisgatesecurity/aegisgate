# AegisGate v0.2.1 - GUI Administration, TLS Decryption & Compliance Frameworks

## What's New

### GUI Administration Interface
- Electron-based desktop application for Windows, macOS, Linux
- Real-time dashboard with request monitoring
- Traffic inspection with compliance violation highlighting
- Compliance framework management (enable/disable, view violations)
- SSL/TLS certificate generation and management
- Configuration interface for proxy settings

### TLS Decryption Implementation
- HTTPS/HTTP/2 traffic interception capability
- Self-signed certificate generation with configurable CN
- External CA certificate support
- Transparent proxy configuration
- Certificate management via GUI

### Enhanced Compliance Frameworks
- **HIPAA** - Healthcare data protection compliance
- **PCI-DSS** - Payment card industry compliance  
- **SOC2** - Service organization controls compliance
- All frameworks handle nil inputs gracefully
- Comprehensive violation reporting

## Bug Fixes
- Fixed compliance framework nil input handling (CheckRequest, CheckResponse)
- Removed problematic fix.go file causing workflow failures
- Fixed TLS module unused imports causing build failures
- Updated tests with proper assertions

## Technical Details
- Go 1.23 compatible
- All workflows passing (build, CI, test)
- SBOM tracking maintained
- Security-first architecture preserved

## Upgrade Notes
1. GUI can be used separately from core proxy
2. TLS interception is opt-in
3. All compliance frameworks are opt-in

## Full Changelog
https://github.com/aegisgatesecurity/aegisgate/commits/v0.2.1
