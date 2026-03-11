## v0.23.0 - Enterprise Security Gateway for AI Chatbots

**Release Date:** February 27, 2026

---

### 🎉 New Features

#### Integration Test Suite
Comprehensive end-to-end tests ensuring reliability:
- **TestE2EBasicProxyFlow** - Basic request/response flow testing
- **TestE2EMultipleRequests** - Concurrent request handling
- **TestE2EBlockingRequest** - Malicious payload detection (SQL injection, prompt injection, SSN, credit cards)
- **TestE2ELargeRequestBody** - Body size limits validation
- **TestE2ERateLimiting** - Rate limiting enforcement
- **TestE2EHealthCheck** - Health endpoint verification
- **TestE2EStatistics** - Request statistics tracking
- **TestE2EScannerIntegration** - Content scanner integration
- **TestE2EComplianceManagerIntegration** - MITRE ATLAS compliance

#### AI API Mock Fixtures
Test fixtures for development and testing:
- OpenAI Chat Completions
- OpenAI Models
- Anthropic Messages
- Cohere Chat

#### CI/CD Automation
GitHub Actions workflow for:
- Automated testing on push/PR
- Binary builds on multiple platforms
- Security scans (go vet)

#### Documentation
- Comprehensive PROXY_API.md documentation
- Updated README with enterprise focus

---

### ✅ Test Results
- **Integration Tests:** 10/10 PASS
- **Security Scans:** go vet PASS
- **Binary Build:** SUCCESS

---

### 🔒 Security
- MITRE ATLAS framework for prompt injection detection
- NIST AI RMF compliance
- PII detection (SSN, Credit Cards, Emails)
- Credential scanning (API keys, passwords)

---

### 📦 Downloads
- Source code (ZIP)
- Source code (TAR.GZ)

---

### 🙏 Acknowledgments
Thanks to all contributors and the open source community!

---

**Full Changelog:** https://github.com/aegisgatesecurity/aegisgate/compare/v0.22.0...v0.23.0
