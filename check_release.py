# Check release.yml for Select-String
with open('.github/workflows/release.yml', 'r') as f:
    for i, line in enumerate(f, 1):
        if 'Select-String' in line or '$version' in line:
            print(f"Line {i}: {line.rstrip()}")
