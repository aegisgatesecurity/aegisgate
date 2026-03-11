# Check license.go for IsModuleLicensed
with open('pkg/core/license.go', 'r') as f:
    content = f.read()

# Find IsModuleLicensed function
import re
match = re.search(r'func.*IsModuleLicensed[\s\S]*?^}', content, re.MULTILINE)
if match:
    print("=== IsModuleLicensed ===")
    print(match.group(0))
