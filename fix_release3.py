# Fix line 63 properly - escape backslashes for Python
with open('.github/workflows/release.yml', 'r') as f:
    content = f.read()

# Fix line 63 - use raw string to properly handle backslashes
old_pattern = '''$escapedVersion = [regex]::Escape($version)
          $pattern = "\[$escapedVersion\]|## $escapedVersion"
          if (Select-String -Path CHANGELOG.md -Pattern $pattern -Quiet) {'''

new_pattern = r'''$escapedVersion = [regex]::Escape($version)
          $pattern = "\[$escapedVersion\]|## $escapedVersion"
          if (Select-String -Path CHANGELOG.md -Pattern $pattern -Quiet) {'''

content = content.replace(old_pattern, new_pattern)

with open('.github/workflows/release.yml', 'w') as f:
    f.write(content)

# Verify
with open('.github/workflows/release.yml', 'r') as f:
    lines = f.readlines()

print("=== Lines 60-70 ===")
for i in range(59, 70):
    print(f"{i+1}: {lines[i].rstrip()}")
