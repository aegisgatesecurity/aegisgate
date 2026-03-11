# Check registry Initialize function
with open('pkg/core/registry.go', 'r') as f:
    content = f.read()

# Find Initialize function - search for it
import re
match = re.search(r'func \(r \*Registry\) Initialize[\s\S]*?^}', content, re.MULTILINE)
if match:
    print("=== Registry.Initialize ===")
    print(match.group(0)[:1500])
