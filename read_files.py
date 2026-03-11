import os

# Read main.go
with open('cmd/padlock/main.go', 'r') as f:
    content = f.read()
    lines = content.split('\n')
    
print("=== main.go version-related lines ===")
for i, line in enumerate(lines):
    if 'version' in line.lower() and 'flag' not in line.lower():
        print(f"Line {i+1}: {line}")

print("\n=== First 50 lines ===")
for i, line in enumerate(lines[:50]):
    print(f"{i+1}: {line}")
