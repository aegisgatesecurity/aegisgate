# Fix main.go to have const version on its own line
with open('cmd/padlock/main.go', 'r') as f:
    content = f.read()

# Replace the const block with separate const statements
old_block = '''// Build info - set during build
const (
	version   = "1.0.0"
	commit    = "dev"
	// date is set at runtime
)

var date = time.Now().Format(time.RFC3339)'''

new_block = '''// Build info - set during build
const version = "1.0.0"
const commit  = "dev"

// date is set at runtime
var date = time.Now().Format(time.RFC3339)'''

content = content.replace(old_block, new_block)

with open('cmd/padlock/main.go', 'w') as f:
    f.write(content)

print("Fixed main.go")

# Verify
with open('cmd/padlock/main.go', 'r') as f:
    lines = f.readlines()

print("\n=== Lines 20-30 ===")
for i in range(19, 30):
    print(f"{i+1}: {lines[i].rstrip()}")
