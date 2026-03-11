# Check main.go version
with open('cmd/padlock/main.go', 'r') as f:
    content = f.read()

# Check for version line
import re
matches = re.findall(r'const version = "[^"]*"', content)
print("Matches for 'const version =':", matches)

# Also check lines around version
for i, line in enumerate(content.split('\n')):
    if 'version' in line.lower():
        print(f"Line {i+1}: {line}")

# Check actual const block
const_match = re.search(r'const \([\s\S]*?\)', content)
if const_match:
    print("\n=== Const block ===")
    print(const_match.group(0))
