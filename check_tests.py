# Check the failing tests in core_test.go
with open('pkg/core/core_test.go', 'r') as f:
    lines = f.readlines()

# TestRegistry_Initialize - around line 117
print("=== TestRegistry_Initialize (lines 115-145) ===")
for i in range(114, 145):
    print(f"{i+1}: {lines[i].rstrip()}")

print("\n=== TestLicenseManager_CoreTierAlwaysFree (lines 206-220) ===")
for i in range(205, 220):
    print(f"{i+1}: {lines[i].rstrip()}")

print("\n=== TestRegistry_EnableDisable (lines 287-315) ===")
for i in range(286, 315):
    print(f"{i+1}: {lines[i].rstrip()}")
