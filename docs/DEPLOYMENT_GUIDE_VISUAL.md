# AegisGate Visual Deployment Guide
═══════════════════════════════════════════════════════════════════════════════

📸 A PICTURE IS WORTH A THOUSAND WORDS! 📸

This guide uses screenshots and visual cues to help you deploy AegisGate
step-by-step, even if you've never done anything like this before.

═══════════════════════════════════════════════════════════════════════════════

TABLE OF CONTENTS
──────────────────
1. What You'll Need
2. Installing Docker (The Foundation)
3. Installing AegisGate  
4. First-Time Setup
5. Configuring AI Providers
6. Testing Your Setup
7. Production Deployment

═══════════════════════════════════════════════════════════════════════════════

SECTION 1: WHAT YOU'LL NEED
──────────────────────────

Before we begin, make sure you have:

□ Access to a computer (Windows 10/11, macOS, or Linux)
□ An internet connection
□ Admin rights on your computer (to install Docker)
□ About 15-20 minutes of uninterrupted time

COMPATIBILITY CHECK:
─────────────────────────────────────────────────────────────
| Your Computer | Will AegisGate Work? | 
─────────────────────────────────────────────────────────────
| Windows 10    | ✅ Yes | 
| Windows 11    | ✅ Yes |
| Mac (Intel)   | ✅ Yes |
| Mac (Apple Silicon/M1/M2/M3) | ✅ Yes |
| Ubuntu 20.04+ | ✅ Yes |
| Debian 11+    | ✅ Yes |
| AWS/Google Cloud Server | ✅ Yes |
| Old computer from 2015 | ⚠️ May work, try it! |
└─────────────────────────────────────────────────────────────

═══════════════════════════════════════════════════════════════════════════════

SECTION 2: INSTALLING DOCKER (The Foundation)
───────────────────────────────────────────────

AegisGate runs inside something called "Docker" - think of it like a
shipping container that holds everything AegisGate needs to run.

STEP-BY-STEP WITH PICTURES:

┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 2.1: DOWNLOAD DOCKER                                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ 1. Open your web browser (Chrome, Firefox, Safari, or Edge)                 │
│ 2. Go to this website:                                                     │
│                                                                             │
│     📎 https://www.docker.com/products/docker-desktop                      │
│                                                                             │
│ 3. Click the big "Download for Windows" button                            │
│    (or "Download for Mac" if you have a Mac)                              │
│                                                                             │
│ ┌───────────────────────────────────────────────────────────────────────┐  │
│ │                    [ DOWNLOAD FOR FREE ]                              │  │
│ │                   ⬇️                                                   │  │
│ │              Click This Button!                                       │  │
│ └───────────────────────────────────────────────────────────────────────┘  │
│                                                                             │
│ 4. Wait for the file to download (this may take a few minutes)            │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 2.2: INSTALL DOCKER                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ After downloading:                                                         │
│                                                                             │
│ ▼ WINDOWS:                                                                  │
│   1. Click the downloaded file (typically in Downloads folder)            │
│   2. Double-click "Docker Desktop Installer.exe"                         │
│   3. Click "Yes" if it asks for permission                                │
│   4. Keep clicking "OK" or "Next" until it says "Install"                │
│   5. Wait for installation to complete                                    │
│   6. Click "Close" or "Finish"                                            │
│                                                                             │
│ ▼ MAC:                                                                       │
│   1. Click the downloaded .dmg file                                       │
│   2. Drag the Docker icon to your Applications folder                     │
│   3. Open Docker from Applications                                        │
│   4. Enter your password when asked                                       │
│                                                                             │
│ ▼ LINUX:                                                                    │
│   Run these commands in terminal:                                          │
│                                                                             │
│   sudo apt-get update                                                      │
│   sudo apt-get install docker.io                                          │
│   sudo systemctl start docker                                             │
│   sudo systemctl enable docker                                            │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 2.3: START DOCKER                                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ ▼ WINDOWS:                                                                  │
│   1. Look for Docker icon in taskbar (bottom-right)                       │
│   2. If you don't see it, click Windows button and type "Docker"          │
│   3. Click "Docker Desktop"                                                │
│   4. Wait for it to say "Docker is running"                                │
│                                                                             │
│   SUCCESS! You should see:                                                 │
│   ┌──────────────────────────────────────────┐                            │
│   │ 🐳 Docker Desktop            [Running]  │                            │
│   └──────────────────────────────────────────┘                            │
│                                                                             │
│ ▼ MAC:                                                                       │
│   1. Look for Docker icon in menu bar (top-right)                        │
│   2. Click it and select "Open Docker Desktop"                           │
│   3. Wait for "Docker is running" in the menu                             │
│                                                                             │
│ ▼ LINUX:                                                                    │
│   Docker should already be running! Type "docker --version" to verify.   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 2.4: VERIFY DOCKER IS WORKING                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ Open terminal (PowerShell on Windows, Terminal on Mac/Linux)              │
│ Type this command and press Enter:                                         │
│                                                                             │
│     docker --version                                                        │
│                                                                             │
│ YOU SHOULD SEE:                                                             │
│     Docker version 24.0.0 or newer                                         │
│                                                                             │
│ 🎉 If you see a version number, Docker is installed correctly!            │
└─────────────────────────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════════════════════════

SECTION 3: INSTALLING AEGISGATE
──────────────────────────────

Now let's install AegisGate!

┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 3.1: OPEN TERMINAL                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ ▼ WINDOWS:                                                                  │
│   1. Press Windows button                                                   │
│   2. Type "PowerShell"                                                     │
│   3. Click "Windows PowerShell"                                            │
│                                                                             │
│ ▼ MAC:                                                                       │
│   1. Press Command + Spacebar                                               │
│   2. Type "Terminal"                                                       │
│   3. Press Enter                                                           │
│                                                                             │
│ ▼ LINUX:                                                                    │
│   1. Press Ctrl + Alt + T                                                  │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 3.2: RUN THE INSTALLATION COMMAND                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ Copy this entire command and paste it into your terminal:                 │
│                                                                             │
│ ╔═══════════════════════════════════════════════════════════════════════╗  │
│ ║                                                                       ║  │
│ ║   docker run -d --name aegisgate -p 8080:8080 -p 8443:8443          ║  │
│ ║       ghcr.io/aegisgatesecurity/aegisgate:latest                      ║  │
│ ║                                                                       ║  │
│ ╚═══════════════════════════════════════════════════════════════════════╝  │
│                                                                             │
│ HOW TO PASTE:                                                              │
│   • WINDOWS: Right-click in the window, press Enter                       │
│   • MAC: Command + V, then Enter                                           │
│   • LINUX: Ctrl + Shift + V, then Enter                                   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 3.3: WAIT FOR DOWNLOAD                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ This may take 1-5 minutes depending on your internet.                    │
│                                                                             │
│ You'll see lots of text scrolling - DON'T WORRY! That's normal.          │
│                                                                             │
│ When it's done, you'll see a long code like:                              │
│     abc123def456789...                                                    │
│                                                                             │
│ Just wait for the text to stop scrolling.                                │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 3.4: VERIFY AEGISGATE IS RUNNING                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ In your terminal, type:                                                   │
│                                                                             │
│     docker ps                                                              │
│                                                                             │
│ YOU SHOULD SEE:                                                             │
│ ┌─────────────────────────────────────────────────────────────────────┐   │
│ │ CONTAINER ID   IMAGE                    STATUS        PORTS        │   │
│ │ abc123def...  aegisgate:latest   Up 2 minutes   8080->8080...       │   │
│ └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│ 🎉 AEGISGATE IS NOW RUNNING!                                               │
└─────────────────────────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════════════════════════

SECTION 4: FIRST-TIME SETUP
──────────────────────────

Let's set up AegisGate for the first time!

┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 4.1: OPEN THE AEGISGATE DASHBOARD                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ 1. Open your web browser                                                   │
│ 2. In the address bar, type:                                              │
│                                                                             │
│      http://localhost:8080                                                │
│                                                                             │
│ 3. Press Enter                                                             │
│                                                                             │
│ YOU SHOULD SEE:                                                             │
│ ┌─────────────────────────────────────────────────────────────────────┐   │
│ │                                                                       │   │
│ │                      🛡️ AEGISGATE 🔐                                │   │
│ │                                                                       │   │
│ │                  ┌─────────────────────────┐                        │   │
│ │                  │  Username: [________]   │                        │   │
│ │                  │  Password: [________]   │                        │   │
│ │                  │       [ Login ]        │                        │   │
│ │                  └─────────────────────────┘                        │   │
│ │                                                                       │   │
│ └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│ Default login:                                                              │
│   Username: admin                                                          │
│   Password: changeme                                                       │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 4.2: CHANGE YOUR PASSWORD                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ After logging in:                                                           │
│                                                                             │
│ 1. Look for your username in the top-right corner                        │
│ 2. Click on it                                                             │
│ 3. Select "Settings" or "Account Settings"                               │
│ 4. Find "Change Password"                                                 │
│ 5. Enter:                                                                  │
│      Current password: changeme                                           │
│      New password: YourSecurePassword123!                                │
│      Confirm: YourSecurePassword123!                                      │
│ 6. Click "Save" or "Update"                                               │
│                                                                             │
│ ⚠️ IMPORTANT: Remember this password!                                      │
└─────────────────────────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════════════════════════

SECTION 5: CONFIGURING AI PROVIDERS
───────────────────────────────────

Now let's connect your AI providers (OpenAI, Anthropic, etc.)

┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 5.1: FIND THE SETTINGS PAGE                                           │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ In the AegisGate dashboard:                                                │
│                                                                             │
│ 1. Look for the sidebar or menu                                           │
│ 2. Click "Settings"                                                        │
│ 3. Click "AI Providers" or "Providers"                                    │
│                                                                             │
│ ┌─────────────────────────────────────────────────────────────────────┐   │
│ │    Sidebar                                                           │   │
│ │    ───────                                                           │   │
│ │    🏠 Dashboard                                                      │   │
│ │    📊 Analytics                                                      │   │
│ │    ⚙️ Settings  ◄── CLICK HERE                                      │   │
│ │    🔒 Security                                                       │   │
│ │    📋 Logs                                                           │   │
│ └─────────────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 5.2: ADD OPENAI API KEY                                               │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ To add OpenAI:                                                              │
│                                                                             │
│ 1. Find "OpenAI" in the providers list                                   │
│ 2. Look for "Add API Key" or a plus (+) button                            │
│ 3. Click it                                                                │
│                                                                             │
│ GET YOUR OPENAI KEY:                                                       │
│ 1. Go to: https://platform.openai.com/api-keys                           │
│ 2. If needed, log in or create account                                    │
│ 3. Click "Create new secret key"                                         │
│ 4. Give it a name (like "AegisGate")                                      │
│ 5. Click the copy button to copy the key                                 │
│                                                                             │
│ PASTE IT INTO AEGISGATE:                                                   │
│ 1. Paste the key into the API key field                                   │
│ 2. Click "Save"                                                            │
│                                                                             │
│ 🎉 OpenAI is now connected!                                                │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 5.3: ADD OTHER PROVIDERS (OPTIONAL)                                    │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ ANTHROPIC (Claude):                                                        │
│ 1. Go to: https://console.anthropic.com/                                │
│ 2. Click "API Keys"                                                       │
│ 3. Create a new key                                                        │
│ 4. Copy and paste into AegisGate                                           │
│                                                                             │
│ GOOGLE AI:                                                                 │
│ 1. Go to: https://aistudio.google.com/app/apikey                        │
│ 2. Create API key                                                          │
│ 3. Copy and paste into AegisGate                                           │
│                                                                             │
│ AZURE OPENAI:                                                              │
│ 1. Go to Azure Portal → Azure OpenAI                                      │
│ 2. Create deployment                                                      │
│ 3. Add endpoint and key to AegisGate                                      │
└─────────────────────────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════════════════════════

SECTION 6: TESTING YOUR SETUP
─────────────────────────────

Let's make sure everything is working!

┌─────────────────────────────────────────────────────────────────────────────┐
│ TEST 1: HEALTH CHECK                                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ In your browser, go to:                                                   │
│                                                                             │
│      http://localhost:8080/health                                         │
│                                                                             │
│ YOU SHOULD SEE:                                                            │
│     {"status":"healthy","version":"1.0.11","tier":"community"}           │
│                                                                             │
│ 🎉 If you see this, AegisGate is healthy!                                │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ TEST 2: TRY A REQUEST                                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ Make a test request through AegisGate:                                    │
│                                                                             │
│ In your terminal, run:                                                    │
│                                                                             │
│ curl -x http://localhost:8443 \                                         │
│   -H "Authorization: Bearer YOUR_OPENAI_KEY" \                          │
│   https://api.openai.com/v1/models                                       │
│                                                                             │
│ (Replace YOUR_OPENAI_KEY with your actual key)                          │
│                                                                             │
│ YOU SHOULD SEE: A list of available AI models!                           │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ TEST 3: CHECK THE DASHBOARD                                                 │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ 1. Go back to http://localhost:8080                                      │
│ 2. Look at the dashboard                                                  │
│ 3. You should see:                                                        │
│      ✓ Request count                                                      │
│      ✓ Active routes                                                      │
│      ✓ Security status                                                   │
│      ✓ Compliance status                                                 │
│                                                                             │
│ 🎉 ALL TESTS PASSED! You have a working AegisGate!                       │
└─────────────────────────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════════════════════════

SECTION 7: PRODUCTION DEPLOYMENT
────────────────────────────────

Want to run AegisGate in production? Here's what you need:

┌─────────────────────────────────────────────────────────────────────────────┐
│ RECOMMENDED PRODUCTION SETUP                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ For production, we recommend:                                              │
│                                                                             │
│ ✅ Use PostgreSQL database (instead of file storage)                     │
│ ✅ Enable TLS/HTTPS                                                       │
│ ✅ Set up reverse proxy (nginx, traefik)                                 │
│ ✅ Configure logging to file or syslog                                    │
│ ✅ Set up monitoring (Prometheus + Grafana)                               │
│ ✅ Use Docker Compose for easy management                                 │
│ ✅ Set up regular backups                                                 │
│ ✅ Configure firewall rules                                               │
│                                                                             │
│ See the full deployment guide for detailed instructions.                 │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ DOCKER COMPOSE RECOMMENDED                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ For production, use this docker-compose.yml:                             │
│                                                                             │
│ version: '3.8'                                                             │
│ services:                                                                  │
│   aegisgate:                                                               │
│     image: ghcr.io/aegisgatesecurity/aegisgate:latest                    │
│     ports:                                                                  │
│       - "8080:8080"                                                        │
│       - "8443:8443"                                                        │
│     environment:                                                           │
│       - AEGISGATE_STORAGE_MODE=postgres                                   │
│       - DATABASE_URL=postgres://user:pass@db:5432/aegisgate              │
│     depends_on:                                                            │
│       - db                                                                 │
│                                                                             │
│   db:                                                                      │
│     image: postgres:14                                                    │
│     environment:                                                           │
│       POSTGRES_DB: aegisgate                                              │
│       POSTGRES_USER: user                                                 │
│       POSTGRES_PASSWORD: pass                                             │
│     volumes:                                                               │
│       - postgres_data:/var/lib/postgresql/data                           │
│                                                                             │
│ volumes:                                                                   │
│   postgres_data:                                                           │
└─────────────────────────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════════════════════════

SUMMARY
───────

WHAT YOU'VE ACCOMPLISHED:

✓ Installed Docker
✓ Installed AegisGate
✓ Connected to AI providers  
✓ Verified everything works
✓ Learned the basics of the dashboard

WHAT'S NEXT:

• Configure security rules (see SECURITY.md)
• Set up compliance frameworks (see COMPLIANCE.md)
• Enable threat detection
• Set up monitoring and alerting

═══════════════════════════════════════════════════════════════════════════════

Need more help?
  📧 support@aegisgatesecurity.io
  🐛 https://github.com/aegisgatesecurity/aegisgate/issues

═══════════════════════════════════════════════════════════════════════════════
