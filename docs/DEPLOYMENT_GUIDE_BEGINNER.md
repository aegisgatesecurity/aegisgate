# AegisGate AI Security Gateway - Beginner Deployment Guide

Welcome to AegisGate! This guide will help you get up and running with the AegisGate AI Security Gateway in just a few minutes. We've designed this guide to be friendly and easy to follow, even if you're new to security tools or Go programming.

---

## 1. Introduction - What is AegisGate?

**AegisGate** is a security gateway that sits between your applications and AI services (like OpenAI, Anthropic, Google AI, etc.). Think of it as a security guard that checks all traffic going to and from AI services.

### What Does AegisGate Do?

AegisGate protects your AI applications in several important ways:

- **🔒 Prevents Data Leaks**: Stops sensitive information like passwords, credit card numbers, and personal data from being sent to or received from AI services
- **🛡️ Blocks Attacks**: Detects and stops malicious attempts to manipulate AI systems (like prompt injection attacks)
- **📋 Ensures Compliance**: Helps you meet security standards like GDPR, HIPAA, SOC2, and AI-specific frameworks like OWASP AI Top 10 and MITRE ATLAS
- **📊 Monitors Traffic**: Provides a dashboard to see what's happening with your AI traffic in real-time

### Why Do You Need AegisGate?

If you're building applications that use AI services, you need to ensure:
1. Your customers' private data stays private
2. Bad actors can't trick your AI into doing harmful things
3. You're following security and privacy regulations
4. You can audit and monitor AI interactions

AegisGate handles all of this automatically!

---

## 2. Prerequisites - What Do You Need Before Starting?

Before installing AegisGate, make sure you have the following:

### Minimum Requirements

| Requirement | What You Need | Why |
|-------------|---------------|-----|
| **Computer** | Windows, macOS, or Linux | AegisGate runs on all major operating systems |
| **Internet Connection** | Required | To download dependencies and communicate with AI services |
| **AI Service Access** | An API key from OpenAI, Anthropic, or another AI provider | AegisGate sits in front of your AI service |

### For Docker Installation (Recommended for Beginners)

| Requirement | Minimum Version | How to Check |
|-------------|-----------------|--------------|
| **Docker Desktop** | 4.0 or later | Run `docker --version` in terminal |
| **Docker Compose** | Included with Docker Desktop | Run `docker-compose --version` |

### For Building from Source (Optional)

| Requirement | Minimum Version | How to Check |
|-------------|-----------------|--------------|
| **Go** | 1.23 or later | Run `go version` in terminal |
| **Git** | Any recent version | Run `git --version` in terminal |
| **Make** | Any recent version | Run `make --version` (optional) |

### How to Check If You Have These Installed

Open your terminal or command prompt and run these commands:

```bash
# Check Docker
docker --version

# Check Go (if installing from source)
go version

# Check Git (if installing from source)
git --version
```

If any of these commands fail or show an error, you'll need to install that tool. Don't worry - they're all free and easy to install!

---

## 3. Quick Start - 5 Minute Setup

The fastest way to get AegisGate running is using Docker. Follow these steps:

### Step 1: Create a Configuration File

Create a new file named `aegisgate.yml` in a new folder and paste this simple configuration:

```yaml
# Where your AI service is located
upstream: "https://api.openai.com"

# Network settings
bind_address: "0.0.0.0"
port: 8443

# Enable the dashboard (web interface)
dashboard:
  port: 8080
  enabled: true

# Basic security settings
security:
  rate_limit: 100
  max_body_size: 10485760

# Enable compliance frameworks
compliance:
  frameworks:
    - OWASP_AI_TOP10
    - MITRE_ATLAS
    - NIST_AI_RMF
    - GDPR
```

### Step 2: Start AegisGate with Docker

Run this command in the same folder as your `aegisgate.yml`:

```bash
docker run -d \
  --name aegisgate \
  -p 8443:8443 \
  -p 8080:8080 \
  -v $(pwd)/aegisgate.yml:/config/aegisgate.yml \
  -e AEGISGATE_CONFIG=/config/aegisgate.yml \
  ghcr.io/aegisgatesecurity/aegisgate:latest
```

### Step 3: Verify It's Working

1. Open your web browser and go to: **http://localhost:8080**
2. You should see the AegisGate dashboard!
3. Check the health endpoint: **http://localhost:8080/health**

### Congratulations! 🎉

AegisGate is now running and protecting your AI traffic. Continue to the next sections to learn how to configure it for your specific needs.

---

## 4. Step-by-Step Installation

We offer two ways to install AegisGate. **Option A (Docker)** is recommended for beginners. **Option B (Build from Source)** is for developers who want to customize AegisGate.

---

### Option A: Docker Installation (Easiest for Beginners)

Docker packages everything AegisGate needs into a single container, making installation simple.

#### What is Docker?

Docker is like a "virtual machine" that contains everything needed to run an application. It ensures AegisGate runs the same way regardless of your operating system.

#### Step 1: Install Docker

**Windows:**
1. Download Docker Desktop from https://www.docker.com/products/docker-desktop
2. Run the installer
3. Start Docker Desktop

**macOS:**
1. Download Docker Desktop from https://www.docker.com/products/docker-desktop
2. Move the app to Applications folder
3. Start Docker Desktop

**Linux (Ubuntu):**
```bash
sudo apt-get update
sudo apt-get install docker.io
sudo systemctl start docker
sudo systemctl enable docker
```

#### Step 2: Create Your Configuration

Create a folder for AegisGate, then create a `aegisgate.yml` file inside:

```yaml
# aegisgate.yml - Your first configuration!

# The AI service you want to protect
upstream: "https://api.openai.com"

# Local network settings
bind_address: "0.0.0.0"
port: 8443

# Web dashboard settings
dashboard:
  port: 8080
  enabled: true
  username: "admin"
  password: "change_me_in_production"

# Security settings
security:
  rate_limit: 100        # Requests per minute per user
  max_body_size: 10485760  # 10MB max request size

# Which compliance frameworks to use
compliance:
  frameworks:
    - OWASP_AI_TOP10
    - MITRE_ATLAS
    - NIST_AI_RMF
    - NIST_1500
    - ISO_42001
    - GDPR
    - HIPAA
    - SOC2
```

#### Step 3: Run AegisGate

Open your terminal in the folder containing your `aegisgate.yml` and run:

**Linux/macOS:**
```bash
docker run -d \
  --name aegisgate \
  -p 8443:8443 \
  -p 8080:8080 \
  -p 9090:9090 \
  -v $(pwd)/aegisgate.yml:/config/aegisgate.yml \
  -e AEGISGATE_CONFIG=/config/aegisgate.yml \
  ghcr.io/aegisgatesecurity/aegisgate:latest
```

**Windows (PowerShell):**
```powershell
docker run -d --name aegisgate -p 8443:8443 -p 8080:8080 -p 9090:9090 -v ${PWD}/aegisgate.yml:/config/aegisgate.yml -e AEGISGATE_CONFIG=/config/aegisgate.yml ghcr.io/aegisgatesecurity/aegisgate:latest
```

#### Step 4: Manage AegisGate

```bash
# Check if AegisGate is running
docker ps

# View logs
docker logs aegisgate

# Stop AegisGate
docker stop aegisgate

# Start AegisGate again
docker start aegisgate

# Remove AegisGate container
docker rm aegisgate
```

---

### Option B: Build from Source

If you prefer to build AegisGate yourself, follow these steps. This gives you more control but requires more setup.

#### Step 1: Install Go

**Windows:**
1. Download Go from https://go.dev/dl/
2. Run the installer
3. Restart your terminal

**macOS:**
```bash
# Using Homebrew
brew install go
```

**Linux (Ubuntu):**
```bash
sudo apt-get update
sudo apt-get install golang-go
```

Verify installation:
```bash
go version
# Should show: go version go1.23.x
```

#### Step 2: Download AegisGate Source Code

```bash
# Clone the repository
git clone https://github.com/aegisgatesecurity/aegisgate.git

# Enter the project folder
cd aegisgate
```

#### Step 3: Build AegisGate

```bash
# Build the application
go build -o aegisgate ./cmd/aegisgate

# Verify it built successfully
./aegisgate --version
```

You should see output like: `AegisGate v0.10.1`

#### Step 4: Create Configuration

Create a `aegisgate.yml` file in the project folder (same as shown in Option A).

#### Step 5: Run AegisGate

```bash
# Run with configuration file
./aegisgate --config aegisgate.yml

# Or use environment variables
export AEGISGATE_UPSTREAM="https://api.openai.com"
./aegisgate
```

---

## 5. Configuration - Basic Settings Explained

AegisGate is configured using a YAML file. Here's what each section means:

### Core Settings

```yaml
# Required: The AI service you want to protect
upstream: "https://api.openai.com"    # Can also use: https://api.anthropic.com, etc.

# Required: Network settings
bind_address: "0.0.0.0"   # Listen on all network interfaces
port: 8443                 # Port for the proxy service
```

### Dashboard Settings

```yaml
dashboard:
  port: 8080               # Port for the web dashboard
  enabled: true            # Turn dashboard on/off
  username: "admin"        # Login username
  password: "your_secure_password"  # Login password
```

### Security Settings

```yaml
security:
  rate_limit: 100          # How many requests per minute per user
  max_body_size: 10485760  # Maximum size of requests (in bytes)

  # Enable threat detection
  owasp_ai:
    enabled: true
    action: "block"       # What to do: "block", "log", or "alert"

  atlas:
    enabled: true
    action: "block"
```

### Compliance Frameworks

```yaml
compliance:
  frameworks:
    - OWASP_AI_TOP10       # Detects AI-specific attacks
    - MITRE_ATLAS          # Detects AI attack techniques
    - NIST_AI_RMF          # NIST AI Risk Management Framework
    - NIST_1500            # NIST AI Controls
    - ISO_42001            # AI Management System Standard
    - GDPR                 # European data protection
    - HIPAA                # Healthcare data protection
    - SOC2                 # Security compliance
```

### Environment Variables (Alternative to YAML)

You can also configure AegisGate using environment variables:

```bash
# Core settings
export AEGISGATE_UPSTREAM="https://api.openai.com"
export AEGISGATE_PORT="8443"
export AEGISGATE_LOG_LEVEL="info"

# Dashboard
export AEGISGATE_DASHBOARD_PORT="8080"
export AEGISGATE_DASHBOARD_ENABLED="true"

# Security
export AEGISGATE_RATE_LIMIT="100"
export AEGISGATE_OWASP_AI_ENABLED="true"
export AEGISGATE_ATLAS_ENABLED="true"
```

---

## 6. Verification - How to Confirm It's Working

After installing AegisGate, let's verify everything is working correctly.

### Method 1: Check the Dashboard

1. Open your browser
2. Go to **http://localhost:8080**
3. You should see the AegisGate dashboard with:
   - Traffic statistics
   - Security status
   - Compliance scores

### Method 2: Check Health Endpoints

Run these commands in your terminal:

```bash
# Health check
curl http://localhost:8080/health
# Expected response: {"status":"healthy"}

# Version info
curl http://localhost:8080/version
# Expected response: {"version":"v0.10.1"}

# Readiness check
curl http://localhost:8080/ready
# Expected response: {"status":"ready"}
```

### Method 3: Test the Proxy

Send a test request through AegisGate to verify it's proxying correctly:

```bash
# Test proxy health
curl http://localhost:8443/health

# Test with your AI API (replace with your actual API key)
curl -x http://localhost:8443 \
  -H "Authorization: Bearer YOUR_OPENAI_API_KEY" \
  https://api.openai.com/v1/models
```

### Method 4: Check Logs

```bash
# Docker logs
docker logs aegisgate

# Or if running locally
# Check the terminal where AegisGate is running
```

You should see log messages like:
```
[INFO] Starting AegisGate v0.10.1
[INFO] Dashboard listening on :8080
[INFO] Proxy listening on :8443
[INFO] Security scanner initialized
```

---

## 7. Troubleshooting - Common Issues and Solutions

Don't worry if you run into problems! Here are solutions to the most common issues:

### Issue: "Docker command not found"

**Problem:** You don't have Docker installed or it's not in your PATH.

**Solution:**
- Download and install Docker Desktop from https://www.docker.com/products/docker-desktop
- Restart your computer after installation
- Try again

---

### Issue: "Port already in use"

**Problem:** Another application is using port 8080 or 8443.

**Solution:**
1. Find what's using the port:
   - Windows: `netstat -ano | findstr :8080`
   - Mac/Linux: `lsof -i :8080`
2. Stop the other application, or change AegisGate's port in your config:

```yaml
port: 8444           # Change from 8443 to 8444
dashboard:
  port: 9090         # Change from 8080 to 9090
```

---

### Issue: "Connection refused" or "Cannot connect"

**Problem:** AegisGate isn't running or can't start.

**Solution:**
1. Check if AegisGate is running:
   ```bash
   docker ps
   ```
2. Check the logs for errors:
   ```bash
   docker logs aegisgate
   ```
3. Common fixes:
   - Make sure your `aegisgate.yml` has valid syntax
   - Check that the `upstream` URL is correct
   - Ensure you have internet connectivity

---

### Issue: "Permission denied" when running locally

**Problem:** You don't have permission to run the AegisGate binary.

**Solution:**
- Linux/Mac: `chmod +x aegisgate`
- Or run with sudo: `sudo ./aegisgate --config aegisgate.yml`

---

### Issue: "TLS certificate error"

**Problem:** You're seeing SSL/TLS certificate warnings.

**Solution:**
This is expected in development! For production:
1. Generate proper TLS certificates using Let's Encrypt or your own CA
2. Update your config:

```yaml
tls:
  enabled: true
  cert_file: "/path/to/your/cert.pem"
  key_file: "/path/to/your/key.pem"
```

---

### Issue: "Upstream connection failed"

**Problem:** AegisGate can't reach your AI service.

**Solution:**
1. Check your internet connection
2. Verify the upstream URL is correct
3. Check if your AI API key is valid
4. Try pinging the upstream:
   ```bash
   ping api.openai.com
   ```

---

### Issue: "Out of memory" or high memory usage

**Problem:** AegisGate is using too much memory.

**Solution:**
1. Reduce the rate limit:
   ```yaml
   security:
     rate_limit: 50
   ```
2. Limit concurrent connections:
   ```yaml
   mitm:
     max_connections: 1000
   ```

---

### Issue: Dashboard shows "Not configured"

**Problem:** The dashboard isn't receiving data.

**Solution:**
1. Verify security frameworks are enabled in your config
2. Wait a moment for traffic to flow through
3. Check logs for any errors

---

## 8. Next Steps - Where to Go From Here

Now that AegisGate is running, here's what you can do next:

### Learn More About AegisGate

- 📖 **Read the Full Documentation**: Check out `docs/CONFIGURATION.md` for all configuration options
- 🔒 **Understand Security Features**: See `docs/SECURITY.md` for detailed security documentation
- 📊 **Explore the Dashboard**: Visit http://localhost:8080 and click around!

### Configure Advanced Features

1. **Enable HTTPS Interception (MITM)**
   ```yaml
   mitm:
     enabled: true
     port: 3128
   ```

2. **Add ML-Based Anomaly Detection**
   ```yaml
   ml:
     enabled: true
     anomaly_threshold: 0.85
   ```

3. **Configure Prometheus Metrics**
   ```yaml
   metrics:
     enabled: true
     port: 9090
     path: "/metrics"
   ```

### Set Up Production Deployment

- Use **Docker Compose** for easy management
- Set up **Kubernetes** for scaling (see `deploy/k8s/`)
- Configure **TLS certificates** for secure connections
- Set up **log aggregation** for monitoring
- Enable **alerting** for security events

### Stay Updated

- Check GitHub releases for new versions: https://github.com/aegisgatesecurity/aegisgate/releases
- Review the changelog: `CHANGELOG.md`

### Get Help

- 🐛 **Report Bugs**: https://github.com/aegisgatesecurity/aegisgate/issues
- 💡 **Feature Requests**: https://github.com/aegisgatesecurity/aegisgate/discussions
- 📧 **Commercial Support**: commercial@aegisgatesecurity.io

---

## Quick Reference Card

Here's a handy summary of the most common commands:

```bash
# Start AegisGate (Docker)
docker run -d --name aegisgate -p 8443:8443 -p 8080:8080 -v $(pwd)/aegisgate.yml:/config/aegisgate.yml ghcr.io/aegisgatesecurity/aegisgate:latest

# Check status
curl http://localhost:8080/health

# View logs
docker logs aegisgate

# Stop AegisGate
docker stop aegisgate

# Restart AegisGate
docker restart aegisgate

# Update to latest version
docker pull ghcr.io/aegisgatesecurity/aegisgate:latest
```

---

## Thank You! 🎉

You did it! AegisGate is now protecting your AI applications. Remember:

- AegisGate sits between your apps and AI services
- It blocks attacks and prevents data leaks automatically
- The dashboard shows you what's happening in real-time
- You can customize it to fit your specific needs

If you have any questions or run into issues, don't hesitate to check the troubleshooting section or reach out for help.

**Happy securing! 🔒**
