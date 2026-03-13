# Update README to remove specific pricing

$aegisGatePath = "C:\Users\Administrator\Desktop\Testing\AegisGate"
$readmePath = "$aegisGatePath\README.md"

$content = Get-Content $readmePath -Raw

# Replace the pricing table header row - remove $29/mo and $99/mo
$content = $content -replace '\| Developer \(\$29/mo\) \| Professional \(\$99/mo\) \| Enterprise \(Custom\) \|', '| Developer | Professional | Enterprise |'

# Replace the pricing row that shows the actual prices
$content = $content -replace '\| \$29/mo \| \$99/mo \| Custom \|', '| Contact Sales | Contact Sales | Custom |'

# Also update the section headers in tiers to not reference specific prices
$content = $content -replace '### Paid Tiers', '### Paid Tiers\n\n> Contact sales@aegisgate.io for pricing information'

# Also update the Free Tier section which mentions "Not required (defaults to Community)"
# that already looks good, no change needed there

# Also update the "What's New" section from v1.0.4 to v1.0.6
$content = $content -replace "### What's New in v1.0.4 - Security Hardening Release", "### What's New in v1.0.6 - Security Hardening Release"

Set-Content -Path $readmePath -Value $content -NoNewline

Write-Host "=== README License Badge ==="
Select-String -Path $readmePath -Pattern "License"

Write-Host "`n=== Pricing Table ==="
Select-String -Path $readmePath -Pattern "Developer \| Professional"

Write-Host "`n=== What's New Section ==="
Select-String -Path $readmePath -Pattern "What's New in"