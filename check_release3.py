# Verify line 286 fix
with open('.github/workflows/release.yml', 'r') as f:
    lines = f.readlines()

print("=== Lines 280-295 ===")
for i in range(279, 295):
    print(f"{i+1}: {lines[i].rstrip()}")
