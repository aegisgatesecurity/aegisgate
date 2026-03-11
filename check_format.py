# Check exact format of the version lines
with open('cmd/padlock/main.go', 'r') as f:
    lines = f.readlines()
    for i, line in enumerate(lines[20:30], 21):
        print(f"Line {i}: '{line}'")
