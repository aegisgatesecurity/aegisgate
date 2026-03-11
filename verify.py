# Verify main.go changes
with open('cmd/padlock/main.go', 'r') as f:
    lines = f.readlines()
    print("=== Lines 20-30 ===")
    for i, line in enumerate(lines[20:30], 21):
        print(f"{i}: {line.rstrip()}")

print("\n=== VERSION ===")
with open('VERSION', 'r') as f:
    print(f.read().strip())
