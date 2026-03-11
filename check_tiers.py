import os

# Check core_test.go for old tier references
print("=== core_test.go tier references ===")
with open('pkg/core/core_test.go', 'r') as f:
    content = f.read()
    for i, line in enumerate(content.split('\n'), 1):
        if 'TierCore' in line or 'TierEssential' in line or 'TierPremiumAI' in line or 'TierCompliance' in line:
            print(f"Line {i}: {line}")

print("\n=== hipaa.go TierCompliance references ===")
with open('pkg/compliance/hipaa/hipaa.go', 'r') as f:
    content = f.read()
    for i, line in enumerate(content.split('\n'), 1):
        if 'TierCompliance' in line:
            print(f"Line {i}: {line}")

print("\n=== pci.go TierCompliance references ===")
with open('pkg/compliance/pci/pci.go', 'r') as f:
    content = f.read()
    for i, line in enumerate(content.split('\n'), 1):
        if 'TierCompliance' in line:
            print(f"Line {i}: {line}")

# Also check the core module to see what tiers exist
print("\n=== core module tiers ===")
with open('pkg/core/module.go', 'r') as f:
    content = f.read()
    for i, line in enumerate(content.split('\n'), 1):
        if 'Tier' in line and ('=' in line or 'type' in line):
            print(f"Line {i}: {line}")
