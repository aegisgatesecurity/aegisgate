# Check Docker job in test.yml
with open('.github/workflows/test.yml', 'r') as f:
    lines = f.readlines()

print("=== Docker job (starting at line 50) ===")
for i in range(49, len(lines)):
    print(f"{i+1}: {lines[i].rstrip()}")
