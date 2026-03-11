# Check more context around line 62 and 286
with open('.github/workflows/release.yml', 'r') as f:
    lines = f.readlines()

print("=== Lines 55-70 ===")
for i in range(54, 70):
    print(f"{i+1}: {lines[i].rstrip()}")

print("\n=== Lines 278-295 ===")
for i in range(277, 295):
    print(f"{i+1}: {lines[i].rstrip()}")
