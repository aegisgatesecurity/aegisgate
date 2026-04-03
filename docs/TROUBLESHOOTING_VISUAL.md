# AegisGate Visual Troubleshooting Guide
═══════════════════════════════════════════════════════════════════════════════

🔧 PROBLEMS? LET'S FIX THEM TOGETHER! 🔧

This guide uses visual examples to help you identify and solve common issues.

═══════════════════════════════════════════════════════════════════════════════

TABLE OF CONTENTS
──────────────────
1. How to Identify the Problem (The First Step!)
2. Common Problems and Solutions
3. Error Messages and What They Mean
4. Network and Connection Issues
5. Performance Problems
6. Getting Help

═══════════════════════════════════════════════════════════════════════════════

SECTION 1: IDENTIFY THE PROBLEM
────────────────────────────────

Let's figure out what's wrong!

┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 1: Check if AegisGate is Running                                      │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ Open your terminal and type:                                               │
│                                                                             │
│     docker ps                                                              │
│                                                                             │
│ LOOK FOR THIS: ✅                                                           │
│ ┌─────────────────────────────────────────────────────────────────────┐   │
│ │ CONTAINER ID   IMAGE                    STATUS        PORTS        │   │
│ │ abc123def456   aegisgate:latest   Up 2 minutes   0.0.0.0:8080...  │   │
│ └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│ IF YOU SEE NOTHING or an error: ❌ AegisGate is not running!               │
│                                                                             │
│ IF YOU SEE THIS: ⚠️                                                         │
│ ┌─────────────────────────────────────────────────────────────────────┐   │
│ │ CONTAINER ID   IMAGE                    STATUS        PORTS        │   │
│ │ abc123def456   aegisgate:latest   Restarting    0.0.0.0:8080...   │   │
│ └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│ AegisGate is crashing and restarting!                                      │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 2: Check the Logs (Tells You What's Wrong!)                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ In terminal, type:                                                         │
│                                                                             │
│     docker logs aegisgate                                                  │
│                                                                             │
│ This shows you what's happening inside AegisGate.                          │
│                                                                             │
│ LOOK FOR ERROR MESSAGES - they'll be in RED or say "ERROR"!                │
│                                                                             │
│ EXAMPLE ERROR:                                                             │
│ ┌─────────────────────────────────────────────────────────────────────┐   │
│ │ 2025-01-15T10:30:00Z ERROR Failed to connect to database           │   │
│ │ 2025-01-15T10:30:00Z ERROR dial tcp: connection refused            │   │
│ └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│ The second line tells you WHAT failed: connection refused = database       │
│ isn't running or wrong address                                            │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ STEP 3: Check If Ports Are Available                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ If you get "port already in use", check what's using your ports:          │
│                                                                             │
│ WINDOWS (PowerShell):                                                      │
│     netstat -ano | findstr :8080                                          │
│                                                                             │
│ MAC/LINUX:                                                                 │
│     lsof -i :8080                                                          │
│                                                                             │
│ YOU'LL SEE:                                                                │
│ ┌─────────────────────────────────────────────────────────────────────┐   │
│ │ COMMAND   PID   USER   FD   TYPE   DEVICE SIZE/OFF NODE NAME      │   │
│ │ nginx    1234   root    6u  IPv4   12345      0t0  TCP *:http     │   │
│ └─────────────────────────────────────────────────────────────────────┘   │
│                                                                             │
│ This shows nginx is using port 8080! Stop it or use a different port.   │
└─────────────────────────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════════════════════════

SECTION 2: COMMON PROBLEMS & SOLUTIONS
──────────────────────────────────────

Find your problem below!

═══════════════════════════════════════════════════════════════════════════════

🔴 PROBLEM: "docker: command not found"
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

WHAT THIS LOOKS LIKE:
┌─────────────────────────────────────────────────────────────────────────────┐
│ PS C:\Users\you> docker run ...                                           │
│ docker: The term 'docker' is not recognized as the name of a cmdlet...   │
│                                                                             │
│     or                                                                     │
│                                                                             │
│ bash: docker: command not found                                             │
└─────────────────────────────────────────────────────────────────────────────┘

CAUSE: Docker is not installed or not in your PATH

SOLUTION:
1. Download Docker Desktop from https://www.docker.com/products/docker-desktop
2. Install it
3. RESTART YOUR COMPUTER
4. Try again

IF ALREADY INSTALLED:
• Windows: Search for "Docker Desktop" and open it
• Mac: Look for Docker icon in menu bar, click it, select "Open Docker Desktop"

═══════════════════════════════════════════════════════════════════════════════

🔴 PROBLEM: "port already in use"
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

WHAT THIS LOOKS LIKE:
┌─────────────────────────────────────────────────────────────────────────────┐
│ Error response from daemon: driver failed programming external            │
│ connectivity on endpoint aegisgate: port 8080 is already in use           │
└─────────────────────────────────────────────────────────────────────────────┘

CAUSE: Another program is using port 8080 or 8443

SOLUTION - OPTION 1: Use Different Ports
docker run -d --name aegisgate -p 9090:8080 -p 9443:8443 ghcr.io/aegisgatesecurity/aegisgate:latest

Then access at: http://localhost:9090 (instead of 8080)

SOLUTION - OPTION 2: Find and Stop the Other Program

WINDOWS:
1. Open PowerShell as Administrator
2. Run: netstat -ano | findstr :8080
3. Note the PID number
4. Run: taskkill /PID [NUMBER] /F

MAC/LINUX:
1. Run: lsof -i :8080
2. Note the PID
3. Kill it: kill -9 [PID]

═══════════════════════════════════════════════════════════════════════════════

🔴 PROBLEM: "network error" or "connection refused"  
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

WHAT THIS LOOKS LIKE:
┌─────────────────────────────────────────────────────────────────────────────┐
│ curl: (7) Failed to connect to localhost port 8080: Connection refused     │
└─────────────────────────────────────────────────────────────────────────────┘

CAUSE: AegisGate is not running or not fully started

SOLUTION:
1. Check if AegisGate is running:
   docker ps

2. If not running, start it:
   docker start aegisgate

3. Wait 30 seconds, then try again

4. Check logs:
   docker logs aegisgate

═══════════════════════════════════════════════════════════════════════════════

🔴 PROBLEM: Container keeps restarting
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

WHAT THIS LOOKS LIKE:
┌─────────────────────────────────────────────────────────────────────────────┐
│ docker ps                                                                  │
│ CONTAINER ID   IMAGE              STATUS                      PORTS       │
│ abc123         aegisgate:latest   Restarting (1) 5 sec ago   0.0.0:8080   │
└─────────────────────────────────────────────────────────────────────────────┘

CAUSE: Something is causing AegisGate to crash on startup

DIAGNOSIS:
1. Check the logs:
   docker logs aegisgate

2. Common causes:
   • Invalid configuration file
   • Missing environment variables
   • Permission issues
   • Port conflicts

SOLUTION:

Check your config file for errors:
• Make sure all paths exist
• Check for typos in variable names
• Verify file permissions

Try running with minimal config:
docker run -d --name aegisgate ghcr.io/aegisgatesecurity/aegisgate:latest
(don't mount any volumes or set env vars first)

═══════════════════════════════════════════════════════════════════════════════

🔴 PROBLEM: Cannot connect to AI provider
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

WHAT THIS LOOKS LIKE:
┌─────────────────────────────────────────────────────────────────────────────┐
│ 2025-01-15T10:30:00Z ERROR upstream connect error                          │
│ 2025-01-15T10:30:00Z ERROR dial tcp: i/o timeout                           │
└─────────────────────────────────────────────────────────────────────────────┘

CAUSE: Cannot reach the AI provider (OpenAI, Anthropic, etc.)

DIAGNOSIS:
1. Check if you can reach the provider from Docker:
   docker exec aegisgate curl -v https://api.openai.com

2. Check your API key:
   echo $OPENAI_API_KEY

3. Verify the API key is correct and has credits

SOLUTIONS:

✓ Check internet connection
✓ Verify firewall allows outbound HTTPS (port 443)
✓ Check API key is valid
✓ Check API key has available credits
✓ Try a different AI provider to test

═══════════════════════════════════════════════════════════════════════════════

🟡 PROBLEM: Slow performance
━━━━━━━━━━━━━━━━━━━━━━━━━━

WHAT THIS LOOKS LIKE:
• Requests taking > 1 second
• Dashboard is slow to load
• High CPU or memory usage

DIAGNOSIS:
1. Check resource usage:
   docker stats

2. You should see:
┌─────────────────────────────────────────────────────────────────────────────┐
│ CONTAINER    CPU %   MEM USAGE / LIMIT     NET I/O           BLOCK I/O     │
│ aegisgate    15.00%  256MiB / 512MiB       1MB / 500KB       0B / 0B      │
└─────────────────────────────────────────────────────────────────────────────┘

SOLUTIONS:

✓ Increase memory limit in docker-compose:
  deploy:
    resources:
      limits:
        memory: 2G

✓ Enable caching for improved performance:
  AEGISGATE_CACHE_ENABLED=true

✓ Use PostgreSQL instead of file storage:
  AEGISGATE_STORAGE_MODE=postgres

✓ Reduce logging verbosity:
  AEGISGATE_LOG_LEVEL=warn

═══════════════════════════════════════════════════════════════════════════════

🟡 PROBLEM: TLS/SSL certificate errors
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

WHAT THIS LOOKS LIKE:
┌─────────────────────────────────────────────────────────────────────────────┐
│ x509: certificate signed by unknown authority                             │
│ or                                                                         │
│ tls: first certificate does not comply with anchor constraints             │
└─────────────────────────────────────────────────────────────────────────────┘

CAUSE: Invalid or self-signed certificate

FOR DEVELOPMENT (self-signed certs):
• This is expected! It's normal for self-signed certs.
• Accept the warning in your browser or use curl -k flag

FOR PRODUCTION:
• Use a valid certificate from Let's Encrypt or your CA
• Make sure cert isn't expired:
  openssl x509 -in cert.pem -dates -noout

═══════════════════════════════════════════════════════════════════════════════

🟡 PROBLEM: Dashboard not loading
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

WHAT THIS LOOKS LIKE:
• Blank page when accessing http://localhost:8080
• JavaScript errors in browser console
• Partial page rendering

DIAGNOSIS:
1. Open browser Developer Tools (F12)
2. Check Console tab for errors
3. Check Network tab for failed requests

SOLUTIONS:

✓ Clear browser cache: Ctrl+Shift+R (Windows) or Cmd+Shift+R (Mac)

✓ Check if frontend build exists:
  docker exec aegisgate ls /app/ui/build

✓ Check logs for frontend errors:
  docker logs aegisgate | grep -i ui

✓ Rebuild the container:
  docker pull ghcr.io/aegisgatesecurity/aegisgate:latest

═══════════════════════════════════════════════════════════════════════════════

🟡 PROBLEM: Compliance checks failing
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

WHAT THIS LOOKS LIKE:
• Compliance violations in dashboard
• Blocked requests due to compliance rules

DIAGNOSIS:
1. Check which frameworks are enabled:
   curl http://localhost:8080/api/v1/compliance/status

2. View specific violations:
   curl http://localhost:8080/api/v1/compliance/violations

SOLUTIONS:

✓ Review what triggers the violations (PII, sensitive data, etc.)

✓ Adjust pattern sensitivity if too strict:
  AEGISGATE_COMPLIANCE_SENSITIVITY=medium

✓ Disable specific frameworks if not needed:
  AEGISGATE_COMPLIANCE_FRAMEWORKS=owasp,mitre_only

✓ For false positives, whitelist specific patterns

═══════════════════════════════════════════════════════════════════════════════

SECTION 3: ERROR MESSAGE DICTIONARY
───────────────────────────────────

Quick reference for common error messages:

┌─────────────────────────────────────────────────────────────────────────────┐
│ ERROR                           │ MEANING                                 │
├─────────────────────────────────┼─────────────────────────────────────────┤
│ connection refused             │ Service not running / port wrong        │
│ connection timeout             │ Network slow / firewall blocking        │
│ certificate unknown           │ Self-signed cert / SSL problem          │
│ permission denied              │ No permission to access file/port      │
│ file not found                 │ Config file path incorrect              │
│ invalid configuration          │ Config file has errors                  │
│ rate limited                   │ Too many requests, wait and retry       │
│ unauthorized                   │ Wrong or missing API key/token          │
│ not found                      │ Route/endpoint doesn't exist            │
│ internal server error          │ Bug in AegisGate - check logs           │
│ out of memory                  │ System resources exhausted              │
│ database error                 │ PostgreSQL connection issue             │
└─────────────────────────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════════════════════════

SECTION 4: NETWORK ISSUES
─────────────────────────

┌─────────────────────────────────────────────────────────────────────────────┐
│ CAN'T REACH LOCALHOST                                                       │
├─────────────────────────────────────────────────────────────────────────────┤
│ Try these URLs:                                                            │
│   http://127.0.0.1:8080                                                    │
│   http://0.0.0.0:8080                                                      │
│                                                                             │
│ On Windows, try:                                                           │
│   ipconfig   (find your IPv4 address)                                      │
│   Then go to http://[YOUR_IP]:8080                                        │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ FIREWALL ISSUES                                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│ If running on a server and can't connect:                                  │
│                                                                             │
│ • Linux: sudo ufw allow 8080/tcp                                          │
│ • AWS: Add inbound rule for port 8080                                      │
│ • GCP: Add firewall rule                                                   │
│ • Azure: Configure NSG                                                     │
└─────────────────────────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════════════════════════

SECTION 5: PERFORMANCE PROBLEMS
────────────────────────────────

┌─────────────────────────────────────────────────────────────────────────────┐
│ HIGH MEMORY USAGE                                                          │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ Check: docker stats                                                         │
│                                                                             │
│ If memory > 80%:                                                            │
│ • Enable rate limiting                                                     │
│ • Reduce log retention                                                     │
│ • Use PostgreSQL instead of file storage                                   │
│ • Increase container memory limit                                          │
│                                                                             │
│ Fix: Run with more memory:                                                 │
│ docker run -m 2g ...                                                       │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ HIGH CPU USAGE                                                             │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ Check: docker stats                                                         │
│                                                                             │
│ If CPU > 80%:                                                              │
│ • Too many concurrent requests                                             │
│ • ML detection too intensive                                               │
│ • Regex patterns too complex                                               │
│                                                                             │
│ Fix:                                                                       │
│ AEGISGATE_RATE_LIMIT_REQUESTS=1000                                        │
│ AEGISGATE_ML_ENABLED=false  (if not needed)                               │
└─────────────────────────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════════════════════════

SECTION 6: GETTING HELP
───────────────────────

Still stuck? We're here to help!

┌─────────────────────────────────────────────────────────────────────────────┐
│ BEFORE CONTACTING SUPPORT                                                  │
├─────────────────────────────────────────────────────────────────────────────┤
│ Please have these ready:                                                   │
│                                                                             │
│ 1. Output of: docker version                                               │
│ 2. Output of: docker ps                                                    │
│ 3. Output of: docker logs aegisgate (last 50 lines)                        │
│ 4. Your configuration (remove sensitive values!)                          │
│ 5. What steps you've already tried                                        │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ CONTACT OPTIONS                                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                             │
│ 📧 EMAIL: support@aegisgatesecurity.io                                     │
│                                                                             │
│ 🐛 GITHUB ISSUES: https://github.com/aegisgatesecurity/aegisgate/issues   │
│                                                                             │
│ 💬 DISCORD: (link coming soon)                                             │
│                                                                             │
│ For security issues: security@aegisgatesecurity.io                        │
└─────────────────────────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════════════════════════

EMERGENCY COMMANDS
──────────────────

Keep these handy for emergencies:

# Check status
docker ps

# Restart AegisGate
docker restart aegisgate

# View last 100 log lines
docker logs --tail 100 aegisgate

# Stop AegisGate
docker stop aegisgate

# Remove and start fresh
docker rm -f aegisgate
docker pull ghcr.io/aegisgatesecurity/aegisgate:latest

# Get inside the container (for debugging)
docker exec -it aegisgate /bin/sh

# Check resource usage in real-time
docker stats

═══════════════════════════════════════════════════════════════════════════════

Don't give up! Most issues have simple fixes.
We're here to help you succeed!

Questions? support@aegisgatesecurity.io

═══════════════════════════════════════════════════════════════════════════════
