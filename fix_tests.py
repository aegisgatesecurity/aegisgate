# Fix test failures in core_test.go
with open('pkg/core/core_test.go', 'r') as f:
    content = f.read()

# Fix 1: TestRegistry_Initialize - use TierCommunity instead of TierDeveloper
content = content.replace(
    'func TestRegistry_Initialize(t *testing.T) {\n\tregistry := NewRegistry(nil)\n\tmodule := NewMockModule("test-module", "Test Module", TierDeveloper)',
    'func TestRegistry_Initialize(t *testing.T) {\n\tregistry := NewRegistry(nil)\n\tmodule := NewMockModule("test-module", "Test Module", TierCommunity)'
)

# Fix 2: TestRegistry_StartStop - use TierCommunity instead of TierDeveloper
content = content.replace(
    'func TestRegistry_StartStop(t *testing.T) {\n\tregistry := NewRegistry(nil)\n\tmodule := NewMockModule("test-module", "Test Module", TierDeveloper)',
    'func TestRegistry_StartStop(t *testing.T) {\n\tregistry := NewRegistry(nil)\n\tmodule := NewMockModule("test-module", "Test Module", TierCommunity)'
)

# Fix 3: TestRegistry_EnableDisable - use TierCommunity instead of TierDeveloper
content = content.replace(
    'func TestRegistry_EnableDisable(t *testing.T) {\n\tregistry := NewRegistry(nil)\n\tmodule := NewMockModule("test-module", "Test Module", TierDeveloper)',
    'func TestRegistry_EnableDisable(t *testing.T) {\n\tregistry := NewRegistry(nil)\n\tmodule := NewMockModule("test-module", "Test Module", TierCommunity)'
)

# Fix 4: TestLicenseManager_CoreTierAlwaysFree - remove TierDeveloper check
# The test expects both Community and Developer to be free, but Developer now requires a license
# Change the test to only check Community (which is always free)
old_test = '''func TestLicenseManager_CoreTierAlwaysFree(t *testing.T) {
	lm := NewLicenseManager("")

	// Core and Essential should always be allowed
	if !lm.IsModuleLicensed("any-module", TierCommunity) {
		t.Error("TierCommunity should always be licensed")
	}
	if !lm.IsModuleLicensed("any-module", TierDeveloper) {
		t.Error("TierDeveloper should always be licensed")
	}
}'''

new_test = '''func TestLicenseManager_CoreTierAlwaysFree(t *testing.T) {
	lm := NewLicenseManager("")

	// Community tier should always be allowed (free tier)
	if !lm.IsModuleLicensed("any-module", TierCommunity) {
		t.Error("TierCommunity should always be licensed")
	}
}'''

content = content.replace(old_test, new_test)

with open('pkg/core/core_test.go', 'w') as f:
    f.write(content)

print("Fixed core_test.go")

# Verify changes
with open('pkg/core/core_test.go', 'r') as f:
    lines = f.readlines()

print("\n=== TestRegistry_Initialize (line 118-122) ===")
for i in range(117, 122):
    print(f"{i+1}: {lines[i].rstrip()}")

print("\n=== TestLicenseManager_CoreTierAlwaysFree (line 206-214) ===")
for i in range(205, 214):
    print(f"{i+1}: {lines[i].rstrip()}")
