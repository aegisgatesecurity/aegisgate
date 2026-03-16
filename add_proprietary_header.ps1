$header = @"
// SPDX-License-Identifier: MIT
// =========================================================================
// PROPRIETARY - AegisGate Security
// Copyright (c) 2025-2026 AegisGate Security. All rights reserved.
// =========================================================================
//
// This file contains proprietary trade secret information.
// Unauthorized reproduction, distribution, or reverse engineering is prohibited.
// =========================================================================

"@

$basePath = "C:\Users\Administrator\Desktop\Testing\AegisGate"

# Get all .go files (excluding test files and .pb.go)
$goFiles = Get-ChildItem -Path "$basePath\pkg","$basePath\cmd","$basePath\sdk" -Filter "*.go" -Recurse | Where-Object { 
    -not $_.Name.EndsWith("_test.go") -and -not $_.Name.EndsWith(".pb.go") 
}

$updatedCount = 0
foreach ($file in $goFiles) {
    $content = Get-Content $file.FullName -Raw
    if ($content -notmatch "PROPRIETARY - AegisGate Security") {
        $newContent = $header + $content
        Set-Content -Path $file.FullName -Value $newContent -NoNewline
        Write-Host "Updated: $($file.Name)"
        $updatedCount++
    }
}

Write-Host ""
Write-Host "Total files updated: $updatedCount"