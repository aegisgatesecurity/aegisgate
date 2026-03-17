# 📸 AegisGate Screenshot Guide

> **Purpose**: This guide documents all screenshots needed for the AegisGate documentation.  
> **Status**: Automated screenshot capture pending - some screenshots require live deployment.  
> **Tools**: Use Browsermcp.browserScreenshot or Chromedevtools.takeScreenshot

---

## 📋 Screenshot Checklist

| # | Screenshot ID | Page/Resource | Description | Priority |
|---|---------------|---------------|-------------|----------|
| 1 | SS-01 | Landing/Dashboard | AegisGate Admin UI - Main Dashboard | 🔴 Critical |
| 2 | SS-02 | Getting Started | Docker Compose up and running terminal | 🔴 Critical |
| 3 | SS-03 | Configuration | Admin UI - AI Provider Configuration | 🔴 Critical |
| 4 | SS-04 | Configuration | Admin UI - Security Policy Settings | 🟡 High |
| 5 | SS-05 | Configuration | Admin UI - Compliance Framework Selection | 🟡 High |
| 6 | SS-06 | Monitoring | Traffic Analytics Dashboard | 🟡 High |
| 7 | SS-07 | Monitoring | Threat Detection Alerts Panel | 🟡 High |
| 8 | SS-08 | Logs | Real-time Request/Response Logs | 🟢 Medium |
| 9 | SS-09 | Compliance | Compliance Report Generator | 🟢 Medium |
| 10 | SS-10 | Users | SSO/SAML User Management | 🟢 Medium |
| 11 | SS-11 | Settings | License Key Activation | 🔴 Critical |
| 12 | SS-12 | Troubleshooting | Health Check Endpoint Response | 🟡 High |
| 13 | SS-13 | API | API Key Rotation Interface | 🟢 Medium |
| 14 | SS-14 | Deployment | Kubernetes Helm Install Success | 🔴 Critical |

---

## 🎯 Screenshot Specifications

### SS-01: Main Dashboard
- **URL**: http://localhost:8080/admin (after login)
- **Capture**: Full viewport, showing sidebar navigation + main content area
- **Annotations needed**: Highlight key areas (sidebar, traffic stats, alerts summary)
- **Dimensions**: 1920x1080 recommended

### SS-02: Docker Running
- **Command shown**: `docker ps` showing aegisgate container running
- **Capture**: Terminal window with green "healthy" status
- **Annotations needed**: Point to port 8080, container status

### SS-03: AI Provider Config
- **URL**: http://localhost:8080/admin/config/providers
- **Capture**: Form with OpenAI, Anthropic, Azure, AWS Bedrock provider fields
- **Annotations needed**: Mark required fields (API Key, Endpoint URL)

### SS-04: Security Policies
- **URL**: http://localhost:8080/admin/config/security
- **Capture**: Toggle switches for threat detection, block mode, rate limiting
- **Annotations needed**: Show enabled vs disabled states

### SS-05: Compliance Frameworks
- **URL**: http://localhost:8080/admin/config/compliance
- **Capture**: Checkboxes for SOC2, HIPAA, PCI-DSS, GDPR, NIST AI RMF
- **Annotations needed**: Multi-select state showing active frameworks

### SS-06: Traffic Analytics
- **URL**: http://localhost:8080/admin/monitoring/analytics
- **Capture**: Charts showing requests/minute, latency, provider distribution
- **Annotations needed**: Time range selector, export button

### SS-07: Threat Alerts
- **URL**: http://localhost:8080/admin/monitoring/alerts
- **Capture**: Table of blocked requests with severity, timestamp, details
- **Annotations needed**: Severity badges (🔴 High, 🟡 Medium, 🟢 Low)

### SS-08: Request Logs
- **URL**: http://localhost:8080/admin/logs/requests
- **Capture**: JSON-formatted log entries with timestamps
- **Annotations needed**: Show PII redaction in action

### SS-09: Compliance Reports
- **URL**: http://localhost:8080/admin/reports/compliance
- **Capture**: Report generation form + generated PDF preview
- **Annotations needed**: Download button, report contents

### SS-10: User Management (SSO)
- **URL**: http://localhost:8080/admin/users
- **Capture**: User table with SSO provider badges
- **Annotations needed**: Role column (Admin, Operator, Viewer)

### SS-11: License Activation
- **URL**: http://localhost:8080/admin/settings/license
- **Capture**: License key input field + activation button
- **Annotations needed**: Success state showing tier/features

### SS-12: Health Check
- **Terminal**: `curl http://localhost:8080/health`
- **Capture**: JSON response showing `{"status":"healthy","version":"1.0.11"}`
- **Annotations needed**: Highlight healthy status

### SS-13: API Key Rotation
- **URL**: http://localhost:8080/admin/settings/api-keys
- **Capture**: List of API keys with rotation status, expiry dates
- **Annotations needed**: "Rotate" button, "Last rotated" column

### SS-14: Kubernetes Deployment
- **Terminal**: `kubectl get pods -n aegisgate`
- **Capture**: All pods Running with ready status
- **Annotations needed**: Pod names, Ready column

---

## 🔧 Automation Commands

### Using Browsermcp (if available)
```bash
# Navigate to page and screenshot
await browserNavigate({ url: "http://localhost:8080/admin" });
await browserScreenshot({ path: "docs/assets/screenshot-01-dashboard.png" });
```

### Using Chromedevtools
```typescript
// Take screenshot of current page
await Chromedevtools.takeScreenshot({ 
  path: "docs/assets/screenshot-01-dashboard.png" 
});
```

### Excalidraw Alternative
If automated screenshots aren't available, create hand-drawn diagrams using:
- Excalidraw.createView for visual representations
- Use ASCII art in markdown as immediate alternative

---

## 📁 File Naming Convention

```
docs/assets/
├── screenshot-01-dashboard.png
├── screenshot-02-docker-running.png
├── screenshot-03-ai-providers.png
├── screenshot-04-security-policies.png
├── screenshot-05-compliance-frameworks.png
├── screenshot-06-analytics.png
├── screenshot-07-alerts.png
├── screenshot-08-logs.png
├── screenshot-09-reports.png
├── screenshot-10-users.png
├── screenshot-11-license.png
├── screenshot-12-health-check.png
├── screenshot-13-api-keys.png
└── screenshot-14-k8s-pods.png
```

---

## ✅ Pre-Capture Checklist

Before capturing screenshots:
- [ ] Docker container running (`docker ps`)
- [ ] Admin UI accessible at http://localhost:8080/admin
- [ ] Default credentials: admin / changeme (or configured)
- [ ] At least one AI provider configured
- [ ] Sample traffic flowing (test request made)

---

*Generated: March 2026 | For AegisGate v1.0.11*
