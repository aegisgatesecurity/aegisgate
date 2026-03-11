# Fix the release.yml PowerShell regex issue
with open('.github/workflows/release.yml', 'r') as f:
    content = f.read()

# Fix line 62 - the problematic pattern
old_pattern = 'if (Select-String -Path CHANGELOG.md -Pattern "\\[\\$version\\]|## \\$version" -Quiet) {'
new_pattern = '$pattern = "[`$version]|## `$version"\n          if (Select-String -Path CHANGELOG.md -Pattern $pattern -Quiet) {'

content = content.replace(old_pattern, new_pattern)

# Also fix line 286 - the regex pattern there might have issues too
old_pattern2 = '$pattern = "(?s)## \\[?$version\\]?[-\\s\\[].*?(?=## \\[|## \\[Unreleased\\]|$)"'
new_pattern2 = '$versionEscaped = [regex]::Escape($version)\n          $pattern = "(?s)## \\[?$versionEscaped\\]?[-\\s\\[].*?(?=## \\[|## \\[Unreleased\\]|$)"'

if old_pattern2 in content:
    content = content.replace(old_pattern2, new_pattern2)
    print("Fixed line 286 pattern")

with open('.github/workflows/release.yml', 'w') as f:
    f.write(content)

print("Fixed release.yml")

# Verify changes
with open('.github/workflows/release.yml', 'r') as f:
    lines = f.readlines()

print("\n=== Lines 60-68 ===")
for i in range(59, 68):
    print(f"{i+1}: {lines[i].rstrip()}")
