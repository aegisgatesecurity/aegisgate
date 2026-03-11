# Fix .dockerignore - remove 'padlock' line but keep 'padlock-linux'
content = open('.dockerignore', 'r').read()

# Remove just the 'padlock' line (not 'padlock-linux')
lines = content.split('\n')
new_lines = [line for line in lines if line.strip() != 'padlock' or 'linux' in line]
content = '\n'.join(new_lines)

# Verify first 12 lines
print("=== First 12 lines of .dockerignore ===")
for i, line in enumerate(new_lines[:12], 1):
    print(f"{i}: {line}")

open('.dockerignore', 'w').write(content)
print("\nFixed .dockerignore - removed 'padlock' exclusion")