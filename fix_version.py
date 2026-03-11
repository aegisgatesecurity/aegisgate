import os
import re

# Fix 1: Update main.go to use const version and version 1.0.0
with open('cmd/padlock/main.go', 'r') as f:
    content = f.read()

# Replace var version = "0.1.0" with const version = "1.0.0"
content = content.replace('var (\n\tversion   = "0.1.0"', 'const (\n\tversion   = "1.0.0"')
content = content.replace('version   = "0.1.0"', 'version   = "1.0.0"')

with open('cmd/padlock/main.go', 'w') as f:
    f.write(content)

print("Fixed main.go version to 1.0.0 with const")

# Fix 2: Update VERSION file
with open('VERSION', 'w') as f:
    f.write('1.0.0\n')

print("Fixed VERSION file to 1.0.0")

# Verify changes
with open('cmd/padlock/main.go', 'r') as f:
    content = f.read()
    if 'const version = "1.0.0"' in content:
        print("Verified: main.go has 'const version = \"1.0.0\"'")
    else:
        # Check for the pattern
        for line in content.split('\n'):
            if 'version' in line and ('1.0.0' in line):
                print(f"Found version line: {line}")

with open('VERSION', 'r') as f:
    print(f"VERSION file now contains: {f.read().strip()}")
