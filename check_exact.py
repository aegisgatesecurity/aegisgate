# Check exact lines
with open('.github/workflows/release.yml', 'r') as f:
    lines = f.readlines()

for i in range(57, 67):
    line = lines[i]
    print(f"Line {i+1}: {repr(line)}")
