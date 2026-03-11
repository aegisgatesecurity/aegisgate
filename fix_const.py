import re

# Read the entire file
with open('cmd/padlock/main.go', 'r') as f:
    content = f.read()

# Fix: replace the const block to have var for date
old_pattern = '''// Build info - set during build
const (
	version   = "1.0.0"
	commit    = "dev"
	date      = time.Now().Format(time.RFC3339)
)'''

new_pattern = '''// Build info - set during build
const (
	version   = "1.0.0"
	commit    = "dev"
	// date is set at runtime
)

var date = time.Now().Format(time.RFC3339)'''

content = content.replace(old_pattern, new_pattern)

# Write back
with open('cmd/padlock/main.go', 'w') as f:
    f.write(content)

print("Fixed main.go const/var issue")

# Verify
with open('cmd/padlock/main.go', 'r') as f:
    lines = f.readlines()
    for i, line in enumerate(lines[20:32], 21):
        print(f"{i}: {line.rstrip()}")
