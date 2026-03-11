with open(r'C:\Users\Administrator\Desktop\Testing\padlock\VERSION', 'w') as f:
    f.write('1.0.0\n')

# Fix core_test.go
content = open(r'C:\Users\Administrator\Desktop\Testing\padlock\pkg\core\core_test.go', 'r').read()

# Replace remaining old tier constants
content = content.replace('TierCore', 'TierCommunity')
content = content.replace('TierEssential', 'TierDeveloper')
content = content.replace('TierPremiumAI', 'TierEnterprise')
content = content.replace('TierCompliance', 'TierEnterprise')

# Fix the expected strings too
content = content.replace('"core"', '"Community"')
content = content.replace('"essential"', '"Developer"')
content = content.replace('"premium-ai"', '"Enterprise"')
content = content.replace('"compliance"', '"Enterprise"')

open(r'C:\Users\Administrator\Desktop\Testing\padlock\pkg\core\core_test.go', 'w').write(content)
print("Fixed core_test.go")

# Fix hipaa.go
content = open(r'C:\Users\Administrator\Desktop\Testing\padlock\pkg\compliance\hipaa\hipaa.go', 'r').read()
content = content.replace('core.TierCompliance', 'core.TierEnterprise')
open(r'C:\Users\Administrator\Desktop\Testing\padlock\pkg\compliance\hipaa\hipaa.go', 'w').write(content)
print("Fixed hipaa.go")

# Fix pci.go
content = open(r'C:\Users\Administrator\Desktop\Testing\padlock\pkg\compliance\pci\pci.go', 'r').read()
content = content.replace('core.TierCompliance', 'core.TierEnterprise')
open(r'C:\Users\Administrator\Desktop\Testing\padlock\pkg\compliance\pci\pci.go', 'w').write(content)
print("Fixed pci.go")

# Fix main.go version to match
content = open(r'C:\Users\Administrator\Desktop\Testing\padlock\cmd\padlock\main.go', 'r').read()
content = content.replace('version   = "0.1.0"', 'version   = "1.0.0"')
open(r'C:\Users\Administrator\Desktop\Testing\padlock\cmd\padlock\main.go', 'w').write(content)
print("Fixed main.go version")

# Fix CHANGELOG.md for v1.0.0
content = open(r'C:\Users\Administrator\Desktop\Testing\padlock\CHANGELOG.md', 'r').read()
if '# Changelog' not in content[:100]:
    new_content = '# Changelog\n\n' + content
    open(r'C:\Users\Administrator\Desktop\Testing\padlock\CHANGELOG.md', 'w').write(new_content)
    content = open(r'C:\Users\Administrator\Desktop\Testing\padlock\CHANGELOG.md', 'r').read()

# Add v1.0.0 entry at top if not present
if '[1.0.0]' not in content:
    changelog_entry = '''## [1.0.0] - 2026-03-10

### Added
- Production release v1.0.0
- 4-tier licensing system (Community/Developer/Professional/Enterprise)
- Comprehensive proxy functionality with ML detection
- Rate limiting and feature gating middleware

### Fixed
- Build system and package conflicts
- Tier references updated to 4-tier system

'''
    # Insert after first line (which should be # Changelog)
    lines = content.split('\n', 1)
    content = lines[0] + '\n' + changelog_entry + lines[1] if len(lines) > 1 else changelog_entry
    open(r'C:\Users\Administrator\Desktop\Testing\padlock\CHANGELOG.md', 'w').write(content)
    print("Updated CHANGELOG.md")

print("\n✅ All fixes applied!")