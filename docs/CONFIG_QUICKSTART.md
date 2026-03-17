# AegisGate Configuration Quickstart
═══════════════════════════════════════════════════════════════════════════════

⚡ GET AEGISGATE CONFIGURED IN 5 MINUTES! ⚡

═══════════════════════════════════════════════════════════════════════════════

THIS GUIDE COVERS:
──────────────────
• Minimum configuration to get started
• Connecting AI providers (OpenAI, Anthropic, etc.)
• Basic security settings
• Testing your configuration

═══════════════════════════════════════════════════════════════════════════════

METHOD 1: DOCKER RUN (FASTEST)
───────────────────────────────

The fastest way to start with custom configuration:

Step 1: Create a configuration file

Create a file named "aegisgate.env" with these contents:

# ==================
# AEGISGATE CONFIG 
# ==================

# TIER: community, developer, professional, enterprise
AEGISGATE_TIER=community

# SERVER
AEGISGATE_HTTP_PORT=8080
AEGISGATE_HTTPS_PORT=8443

# STORAGE  
AEGISGATE_STORAGE_MODE=file
AEGISGATE_DATA_DIR=./data

# LOGGING
AEGISGATE_LOG_LEVEL=info
AEGISGATE_LOG_FORMAT=json

# SECURITY
AEGISGATE_TLS_ENABLED=false
AEGISGATE_RATE_LIMIT_ENABLED=true
AEGISGATE_RATE_LIMIT_REQUESTS=200

# COMPLIANCE (which frameworks to use)
AEGISGATE_COMPLIANCE_ENABLED=true
AEGISGATE_COMPLIANCE_FRAMEWORKS=owasp,mitre_atlas,nist_ai_rmf

# YOUR AI PROVIDER KEYS (add your actual keys!)
OPENAI_API_KEY=sk-your-key-here
# ANTHROPIC_API_KEY=sk-ant-your-key-here

Step 2: Run with your configuration

docker run -d \
  --name aegisgate \
  -p 8080:8080 \
  -p 8443:8443 \
  -v $(pwd)/aegisgate.env:/app/config/aegisgate.env \
  --env-file ./aegisgate.env \
  ghcr.io/aegisgatesecurity/aegisgate:latest

Step 3: Verify it works

curl http://localhost:8080/health

Should return: {"status":"healthy","version":"1.0.11","tier":"community"}

═══════════════════════════════════════════════════════════════════════════════

METHOD 2: DOCKER COMPOSE (RECOMMENDED FOR DEVELOPMENT)
────────────────────────────────────────────────────────

Best for local development with persistent data.

Step 1: Create docker-compose.yml

version: '3.8'

services:
  aegisgate:
    image: ghcr.io/aegisgatesecurity/aegisgate:latest
    container_name: aegisgate
    ports:
      - "8080:8080"   # HTTP
      - "8443:8443"   # HTTPS
      - "9090:9090"   # Metrics
    volumes:
      - ./config:/app/config:ro
      - ./data:/app/data
    env_file:
      - ./config/aegisgate.env
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

Step 2: Create your config directory and files

mkdir -p config data

# Create config/aegisgate.env with the environment variables from Method 1

Step 3: Start it up

# Linux/Mac:
docker-compose up -d

# Windows (PowerShell):
docker compose up -d

Step 4: Check if running

docker-compose ps
docker-compose logs -f

═══════════════════════════════════════════════════════════════════════════════

METHOD 3: BUILD FROM SOURCE (FOR DEVELOPERS)
─────────────────────────────────────────────

If you want to modify the code or run natively.

Step 1: Install Go

Download from: https://go.dev/dl/

Verify: go version

Step 2: Clone and build

git clone https://github.com/aegisgatesecurity/aegisgate.git
cd aegisgate
go build -o aegisgate ./cmd/aegisgate

Step 3: Create config and run

# Create aegisgate.env (from Method 1)
# Then run:
./aegisgate serve --config ./config/aegisgate.env

═══════════════════════════════════════════════════════════════════════════════

CONNECTING AI PROVIDERS
────────────────────────

Here's how to connect each major AI provider:

┌─────────────────────────────────────────────────────────────────────────────┐
│ OPENAI (ChatGPT, GPT-4)                                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│ Get API key:                                                               │
│   1. Go to https://platform.openai.com/api-keys                           │
│   2. Create new secret key                                                │
│   3. Copy it                                                               │
│                                                                             │
│ Add to config:                                                             │
│   OPENAI_API_KEY=sk-xxxx...                                               │
│                                                                             │
│ OR in docker-compose.yml environment:                                     │
│   - OPENAI_API_KEY=sk-xxxx...                                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ ANTHROPIC (Claude)                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│ Get API key:                                                               │
│   1. Go to https://console.anthropic.com/                                  │
│   2. Go to API Keys                                                       │
│   3. Create new key                                                       │
│   4. Copy it                                                               │
│                                                                             │
│ Add to config:                                                             │
│   ANTHROPIC_API_KEY=sk-ant-xxxx...                                        │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ GOOGLE AI (Gemini)                                                         │
├─────────────────────────────────────────────────────────────────────────────┤
│ Get API key:                                                               │
│   1. Go to https://aistudio.google.com/app/apikey                        │
│   2. Create API key                                                       │
│   3. Copy it                                                               │
│                                                                             │
│ Add to config:                                                             │
│   GOOGLE_AI_API_KEY=AIzaSyxxxx...                                         │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ AZURE OPENAI                                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│ Get credentials:                                                           │
│   1. Go to Azure Portal                                                   │
│   2. Find your Azure OpenAI resource                                      │
│   3. Get Endpoint and Key                                                 │
│                                                                             │
│ Add to config:                                                             │
│   AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com/         │
│   AZURE_OPENAI_KEY=your-azure-key                                         │
│   AZURE_OPENAI_DEPLOYMENT=gpt-4                                           │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ AWS BEDROCK                                                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│ Get credentials:                                                           │
│   1. Go to AWS Console                                                    │
│   2. Go to Bedrock → Model access                                         │
│   3. Request access to models                                             │
│   4. Get IAM credentials                                                  │
│                                                                             │
│ Add to config:                                                             │
│   AWS_ACCESS_KEY_ID=AKIAxxxx...                                           │
│   AWS_SECRET_ACCESS_KEY=xxxx...                                          │
│   AWS_REGION=us-east-1                                                   │
└─────────────────────────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════════════════════════

SECURITY CONFIGURATION
───────────────────────

Basic security settings to enable:

┌─────────────────────────────────────────────────────────────────────────────┐
│ RATE LIMITING                                                              │
├─────────────────────────────────────────────────────────────────────────────┤
│ Prevent abuse by limiting requests:                                        │
│                                                                             │
│ AEGISGATE_RATE_LIMIT_ENABLED=true                                         │
│ AEGISGATE_RATE_LIMIT_REQUESTS=200      # requests per minute              │
│ AEGISGATE_RATE_LIMIT_BURST=50          # burst capacity                   │
│                                                                             │
│ Or in YAML:                                                                │
│ rate_limit:                                                                │
│   enabled: true                                                            │
│   requests_per_minute: 200                                                │
│   burst: 50                                                                │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ TLS/HTTPS                                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│ Enable TLS for encrypted traffic:                                         │
│                                                                             │
│ AEGISGATE_TLS_ENABLED=true                                                │
│ AEGISGATE_TLS_CERT=/path/to/cert.pem                                      │
│ AEGISGATE_TLS_KEY=/path/to/key.pem                                        │
│                                                                             │
│ Generate self-signed for testing:                                         │
│ openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem \       │
│   -days 365 -nodes -subj "/CN=localhost"                                   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ API KEY AUTHENTICATION                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│ Require API keys for access:                                               │
│                                                                             │
│ AEGISGATE_API_KEYS_ENABLED=true                                           │
│ AEGISGATE_API_KEY=your-secret-key                                          │
│                                                                             │
│ Generate a secure key:                                                     │
│ openssl rand -base64 32                                                    │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ JWT VALIDATION                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│ Validate JWT tokens:                                                       │
│                                                                             │
│ AEGISGATE_JWT_ENABLED=true                                                │
│ AEGISGATE_JWT_SECRET=your-jwt-secret                                      │
│ AEGISGATE_JWT_ISSUER=https://your-auth-server.com                         │
└─────────────────────────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════════════════════════

COMPLIANCE FRAMEWORKS
──────────────────────

Enable compliance frameworks based on your needs:

┌─────────────────────────────────────────────────────────────────────────────┐
│ FRAMEWORK OPTIONS                                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ Community Tier:                                                            │
│   - OWASP_AI_TOP10 (AI security vulnerabilities)                          │
│   - MITRE_ATLAS (AI attack techniques)                                    │
│   - SOC2 (view only)                                                       │
│   - GDPR (view only)                                                       │
│                                                                             │
│ Developer Tier:                                                            │
│   - All Community +                                                        │
│   - NIST_AI_RMF (AI risk management)                                      │
│   - NIST_1500 (AI controls)                                               │
│                                                                             │
│ Professional Tier:                                                         │
│   - All Developer +                                                        │
│   - HIPAA (healthcare)                                                     │
│   - PCI_DSS (payment cards)                                               │
│   - ISO_27001 (security)                                                  │
│                                                                             │
│ Enterprise Tier:                                                          │
│   - All Professional +                                                    │
│   - ISO_42001 (AI management)                                             │
│   - FEDRAMP (US government)                                               │
│                                                                             │
└─────────────────────────────────────────────────────────────────────────────┘

Enable compliance in config:

AEGISGATE_COMPLIANCE_ENABLED=true
AEGISGATE_COMPLIANCE_FRAMEWORKS=owasp,mitre_atlas,nist_ai_rmf

Or as a comma-separated list:
AEGISGATE_COMPLIANCE_FRAMEWORKS=owasp,mitre_atlas,nist_ai_rmf,soc2,gdpr

═══════════════════════════════════════════════════════════════════════════════

COMPLETE EXAMPLE CONFIGS
─────────────────────────

For each tier, see the config files:
  - config/community.env
  - config/developer.env  
  - config/professional.env
  - config/enterprise.env

Minimum to get started (community):

AEGISGATE_TIER=community
AEGISGATE_LOG_LEVEL=info
OPENAI_API_KEY=sk-your-key-here

Full featured (developer):

AEGISGATE_TIER=developer
AEGISGATE_LICENSE_KEY=dev-xxxxxxxxxxxxx
AEGISGATE_HTTP_PORT=8080
AEGISGATE_HTTPS_PORT=8443
AEGISGATE_TLS_ENABLED=true
AEGISGATE_TLS_CERT=/app/certs/server.pem
AEGISGATE_TLS_KEY=/app/certs/server.key
AEGISGATE_RATE_LIMIT_ENABLED=true
AEGISGATE_RATE_LIMIT_REQUESTS=1000
AEGISGATE_LOG_LEVEL=info
AEGISGATE_METRICS_ENABLED=true
AEGISGATE_METRICS_PORT=9090
OPENAI_API_KEY=sk-your-key-here
ANTHROPIC_API_KEY=sk-ant-your-key-here
AEGISGATE_COMPLIANCE_ENABLED=true
AEGISGATE_COMPLIANCE_FRAMEWORKS=owasp,mitre_atlas,nist_ai_rmf

═══════════════════════════════════════════════════════════════════════════════

TESTING YOUR CONFIGURATION
──────────────────────────

After configuring, verify everything works:

1. Check health endpoint:
   curl http://localhost:8080/health

2. Check version:
   curl http://localhost:8080/version

3. List available models:
   curl -x http://localhost:8443 \
     -H "Authorization: Bearer YOUR_OPENAI_KEY" \
     https://api.openai.com/v1/models

4. Check metrics:
   curl http://localhost:9090/metrics

5. Check compliance status:
   curl http://localhost:8080/api/v1/compliance/status

═══════════════════════════════════════════════════════════════════════════════

NEXT STEPS
──────────

✓ Read CONFIGURATION.md for all possible settings
✓ Check SECURITY.md for hardening recommendations  
✓ See ARCHITECTURE.md for system design
✓ Review TROUBLESHOOTING.md if you have issues

Questions? support@aegisgatesecurity.io

═══════════════════════════════════════════════════════════════════════════════
