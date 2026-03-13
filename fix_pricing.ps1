# Fix pricing table in README

$aegisGatePath = "C:\Users\Administrator\Desktop\Testing\AegisGate"
$readmePath = "$aegisGatePath\README.md"

$content = Get-Content $readmePath -Raw

# Fix the header row - make tier names cleaner
$content = $content -replace '\| Developer \(\$29\/mo\) \|', '| Developer |'
$content = $content -replace '\| Professional \(\$99\/mo\) \|', '| Professional |'
$content = $content -replace '\| Enterprise \(Custom\) \|', '| Enterprise |'

# Fix the note - remove the doubled newlines
$content = $content -replace '### Paid Tiers\r?\n\r?\n> Contact', '### Paid Tiers' + "`n" + '> Contact'

# Fix data rows - show ranges instead of exact numbers (optional - keep for feature comparison)
# These stay as is since they're feature differences, not prices

Set-Content -Path $readmePath -Value $content -NoNewline

Write-Host "=== Pricing Table Check ==="
Select-String -Path $readmePath -Pattern "Requests/min" -Context 0,3
Write-Host "`n=== Paid Tiers Check ==="
Select-String -Path $readmePath -Pattern "Paid Tiers"