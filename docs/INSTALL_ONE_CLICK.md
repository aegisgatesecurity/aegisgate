# AegisGate One-Click Installation Guide
═══════════════════════════════════════════════════════════════════════════════

🛡️ INSTALL AEGISGATE IN 2 MINUTES - NO TECHNICAL EXPERIENCE NEEDED! 🛡️

═══════════════════════════════════════════════════════════════════════════════

TABLE OF CONTENTS
──────────────────
1. What is AegisGate? (Skip if you know!)
2. Choose Your Computer Type
3. Installation (One Command!)
4. How to Verify It's Working
5. How to Use AegisGate
6. Troubleshooting (If Something Goes Wrong)
7. Need Help?

═══════════════════════════════════════════════════════════════════════════════

SECTION 1: WHAT IS AEGISGATE?
─────────────────────────────

AegisGate is a security program that protects your AI applications. Think of it
like a security guard for your AI - it:

✓ Blocks hackers from attacking your AI systems
✓ Keeps your customers' data private  
✓ Helps you follow security rules (like HIPAA, SOC2)
✓ Shows you what's happening with your AI in real-time

YOU DON'T NEED TO UNDERSTAND HOW IT WORKS TO USE IT!

═══════════════════════════════════════════════════════════════════════════════

SECTION 2: CHOOSE YOUR COMPUTER TYPE
────────────────────────────────────

Which type of computer do you have?

► WINDOWS - Use this if your computer has Windows (most laptops/desktops)
  - Look for the Windows button in the bottom-left corner
  - Or type "winver" in the search bar to check

► MAC - Use this if you have a MacBook, iMac, or Mac Mini
  - Look for the Apple icon in the top-left corner

► LINUX - Use this if your computer runs Ubuntu, Debian, or similar
  - Most servers and cloud computers use Linux

═══════════════════════════════════════════════════════════════════════════════

SECTION 3: INSTALLATION - ONE COMMAND!
──────────────────────────────────────

🎯 FOLLOW THESE STEPS EXACTLY - DON'T WORRY ABOUT UNDERSTANDING!

═══════════════════════════════════════════════════════════════════════════════
STEP 3.1: OPEN THE CORRECT APP ON YOUR COMPUTER
═══════════════════════════════════════════════════════════════════════════════

▼ IF YOU HAVE WINDOWS:
   1. Click the Windows button in the bottom-left
   2. Type "PowerShell" in the search box
   3. Click on "Windows PowerShell" (NOT Windows Terminal)
   
   ⚠️ IMPORTANT: A blue or black window will open. This is normal!

▼ IF YOU HAVE MAC:
   1. Click the magnifying glass 🔍 in the top-right corner
   2. Type "Terminal" 
   3. Click on "Terminal" (the black icon)
   
   ⚠️ IMPORTANT: A black window with text will open. This is normal!

▼ IF YOU HAVE LINUX:
   1. Press Ctrl + Alt + T on your keyboard
   2. A terminal window will open

═══════════════════════════════════════════════════════════════════════════════
STEP 3.2: COPY AND PASTE THE INSTALLATION COMMAND
═══════════════════════════════════════════════════════════════════════════════

Copy ONE of these commands depending on your computer:

┌─────────────────────────────────────────────────────────────────────────────┐
│ WINDOWS - Copy this entire line:                                            │
│                                                                             │
│ docker run -d --name aegisgate -p 8080:8080 -p 8443:8443 ghcr.io/aegisgatesecurity/aegisgate:latest                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ MAC - Copy this entire line:                                                 │
│                                                                             │
│ docker run -d --name aegisgate -p 8080:8080 -p 8443:8443 ghcr.io/aegisgatesecurity/aegisgate:latest                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ LINUX - Copy this entire line:                                              │
│                                                                             │
│ docker run -d --name aegisgate -p 8080:8080 -p 8443:8443 ghcr.io/aegisgatesecurity/aegisgate:latest                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘

NOW PASTE IT INTO YOUR TERMINAL WINDOW:
   • WINDOWS: Right-click anywhere in the PowerShell window, then press Enter
   • MAC: Press Command + V, then press Enter
   • LINUX: Press Ctrl + Shift + V, then press Enter

═══════════════════════════════════════════════════════════════════════════════
STEP 3.3: WAIT FOR INSTALLATION TO COMPLETE
═══════════════════════════════════════════════════════════════════════════════

You'll see text scrolling in your terminal. This is normal!

When it's done, you'll see a long string of characters (like a "container ID").

Example of what it will look like:
 CONTAINER ID
 abc123def456...

Don't worry about what this means - just wait until the scrolling stops!

⏱️ THIS TAKES ABOUT 1-2 MINUTES ON A FAST CONNECTION
⏱️ UP TO 5 MINUTES ON A SLOW CONNECTION

═══════════════════════════════════════════════════════════════════════════════

SECTION 4: HOW TO VERIFY IT'S WORKING
────────────────────────────────────

🎉 CONGRATULATIONS! Installation should be complete!

Now let's make sure everything is working:

1. OPEN YOUR WEB BROWSER (Chrome, Firefox, Safari, or Edge)
2. In the address bar at the top, type: http://localhost:8080
3. Press Enter

WHAT YOU SHOULD SEE:
   • A web page with AegisGate logo
   • A login box asking for username/password
   • Or a dashboard showing "healthy" status

IF YOU SEE A LOGIN PAGE:
   • Default username: admin
   • Default password: changeme
   
   ⚠️ IMPORTANT: Change this password for security!

IF YOU SEE AN ERROR:
   Don't panic! Jump to "Section 6: Troubleshooting" below.

═══════════════════════════════════════════════════════════════════════════════

SECTION 5: HOW TO USE AEGISGATE
──────────────────────────────

Now that AegisGate is running, here's what you can do:

┌─────────────────────────────────────────────────────────────────────────────┐
│ ACCESS THE DASHBOARD                                                        │
├─────────────────────────────────────────────────────────────────────────────┤
│ Open your browser and go to:                                                │
│   http://localhost:8080                                                      │
│                                                                             │
│ Username: admin                                                             │
│ Password: changeme (or what you changed it to)                             │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ CONFIGURE YOUR AI PROVIDERS                                                │
├─────────────────────────────────────────────────────────────────────────────┤
│ In the dashboard:                                                            │
│ 1. Look for "Settings" or "Configuration"                                   │
│ 2. Find "AI Providers"                                                       │
│ 3. Add your API keys for:                                                   │
│    • OpenAI (chatgpt)                                                        │
│    • Anthropic (claude)                                                      │
│    • Google AI                                                               │
│    • Azure OpenAI                                                            │
│                                                                             │
│ Look for where to enter: OPENAI_API_KEY=sk-xxxxxxx                         │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ PROTECT YOUR AI TRAFFIC                                                     │
├─────────────────────────────────────────────────────────────────────────────┤
│ Set up routes:                                                               │
│ 1. In dashboard, find "Routes" or "Proxy Configuration"                    │
│ 2. Add your AI endpoint (e.g., api.openai.com)                             │
│ 3. Now route your AI requests through AegisGate instead of directly        │
│                                                                             │
│ Your apps now send requests to:                                             │
│   http://localhost:8080 (instead of api.openai.com)                        │
└─────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│ CHECK THE STATUS                                                            │
├─────────────────────────────────────────────────────────────────────────────┤
│ To check if AegisGate is healthy, open:                                     │
│                                                                             │
│   http://localhost:8080/health                                              │
│                                                                             │
│ You should see:                                                             │
│   {"status":"healthy","version":"1.0.11","tier":"community"}              │
└─────────────────────────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════════════════════════

SECTION 6: TROUBLESHOOTING
──────────────────────────

😟 SOMETHING GOES WRONG? DON'T WORRY - LET'S FIX IT!

═══════════════════════════════════════════════════════════════════════════════
PROBLEM 1: "docker: command not found"
────────────────────────────────────────

WHAT THIS MEANS: You don't have Docker installed.

FIX:
  1. Go to: https://www.docker.com/products/docker-desktop
  2. Click "Download for Windows" (or Mac)
  3. Run the installer
  4. Restart your computer
  5. Try the installation command again

NEED HELP? Watch a video: https://www.youtube.com/watch?v=hQqK55L8o9k

═══════════════════════════════════════════════════════════════════════════════
PROBLEM 2: "port already in use" 
─────────────────────────────────

WHAT THIS MEANS: Another program is using port 8080 or 8443.

FIX:
  • Try changing the port numbers:
  
  WINDOWS (copy and paste):
  docker run -d --name aegisgate -p 9090:8080 -p 9443:8443 ghcr.io/aegisgatesecurity/aegisgate:latest
  
  Then access at: http://localhost:9090

═══════════════════════════════════════════════════════════════════════════════
PROBLEM 3: "network error" or "connection refused"
──────────────────────────────────────────────────

WHAT THIS MEANS: AegisGate isn't fully started yet.

FIX:
  1. Wait 30 more seconds
  2. Try opening http://localhost:8080 again
  3. If still failing, check Docker is running:
     
     WINDOWS: Look for Docker icon in taskbar (bottom right)
     MAC: Look for Docker icon in menu bar (top right)
     LINUX: Run "sudo systemctl status docker"

═══════════════════════════════════════════════════════════════════════════════
PROBLEM 4: "港港港 ERROR" or weird characters
─────────────────────────────────────────────

WHAT THIS MEANS: This is normal! The system is working.

FIX:
  Ignore this and wait 2 minutes, then try http://localhost:8080

═══════════════════════════════════════════════════════════════════════════════
PROBLEM 5: Installation seems stuck (spinning/loading forever)
────────────────────────────────────────────────────────────────

WHAT THIS MEANS: Either slow internet or Docker needs to restart.

FIX:
  1. Press Ctrl + C to cancel
  2. Restart Docker:
     
     WINDOWS: Right-click Docker icon → Restart
     MAC: Right-click Docker icon → Restart
     LINUX: sudo systemctl restart docker
  
  3. Wait 1 minute
  4. Run the installation command again

═══════════════════════════════════════════════════════════════════════════════

SECTION 7: NEED HELP?
─────────────────────

We're here to help! 

📧 EMAIL: support@aegisgatesecurity.io

🐛 REPORT BUGS: https://github.com/aegisgatesecurity/aegisgate/issues

💬 COMMUNITY: Join our Discord! (link coming soon)

📺 VIDEO TUTORIALS: Check YouTube for walkthrough videos

═══════════════════════════════════════════════════════════════════════════════

QUICK REFERENCE CARD
────────────────────

KEEP THESE COMMANDS HANDY:

START AEGISGATE:
  docker start aegisgate

STOP AEGISGATE:
  docker stop aegisgate

CHECK IF RUNNING:
  docker ps

VIEW LOGS:
  docker logs aegisgate

RESTART (if problems):
  docker restart aegisgate

UPDATE TO LATEST VERSION:
  docker pull ghcr.io/aegisgatesecurity/aegisgate:latest
  docker restart aegisgate

═══════════════════════════════════════════════════════════════════════════════

Thank you for choosing AegisGate! 🎉

You now have enterprise-grade AI security running on your computer!

Questions? Don't hesitate to reach out. We're here to help!

═══════════════════════════════════════════════════════════════════════════════
