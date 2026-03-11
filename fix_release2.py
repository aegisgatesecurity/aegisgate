# Fix line 62 properly - need to escape version for regex
with open('.github/workflows/release.yml', 'r') as f:
    content = f.read()

# Fix line 62 properly - use escaped version
old_pattern = '''$version = "${{ steps.extract.outputs.version_num }}"
          $pattern = "[`$version]|## `$version"
          if (Select-String -Path CHANGELOG.md -Pattern $pattern -Quiet) {'''

new_pattern = '''$version = "${{ steps.extract.outputs.version_num }}"
          $escapedVersion = [regex]::Escape($version)
          $pattern = "\[$escapedVersion\]|## $escapedVersion"
          if (Select-String -Path CHANGELOG.md -Pattern $pattern -Quiet) {'''

content = content.replace(old_pattern, new_pattern)

with open('.github/workflows/release.yml', 'w') as f:
    f.write(content)

print("Fixed line 62 properly")

# Verify
with open('.github/workflows/release.yml', 'r') as f:
    lines = f.readlines()

print("\n=== Lines 60-70 ===")
for i in range(59, 70):
    print(f"{i+1}: {lines[i].rstrip()}")
