# Check MockModule in core_test.go
with open('pkg/core/core_test.go', 'r') as f:
    lines = f.readlines()

# Find MockModule definition and NewMockModule
for i, line in enumerate(lines):
    if 'MockModule' in line or 'NewMockModule' in line:
        print(f"{i+1}: {line.rstrip()}")
