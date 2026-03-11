import os
import re

# Check core_test.go more thoroughly
print("=== core_test.go - Lines around 80-90 ===")
with open('pkg/core/core_test.go', 'r') as f:
    lines = f.readlines()
    for i in range(78, 95):  # 0-indexed, so lines 79-95 in editor
        if i < len(lines):
            print(f"Line {i+1}: {lines[i].rstrip()}")

print("\n=== core_test.go - Lines around 180-215 ===")
with open('pkg/core/core_test.go', 'r') as f:
    lines = f.readlines()
    for i in range(178, 220):
        if i < len(lines):
            print(f"Line {i+1}: {lines[i].rstrip()}")

print("\n=== hipaa.go - All lines ===")
with open('pkg/compliance/hipaa/hipaa.go', 'r') as f:
    lines = f.readlines()
    for i, line in enumerate(lines[:30], 1):
        print(f"Line {i}: {line.rstrip()}")

print("\n=== pci.go - All lines ===")
with open('pkg/compliance/pci/pci.go', 'r') as f:
    lines = f.readlines()
    for i, line in enumerate(lines[:30], 1):
        print(f"Line {i}: {line.rstrip()}")
