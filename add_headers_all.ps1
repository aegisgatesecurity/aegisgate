# Add PROPRIETARY headers to all Go source files
$ErrorActionPreference = "Continue"
$basePath = "C:\Users\Administrator\Desktop\Testing\AegisGate"

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

# Find all .go files excluding test files and generated files
$goFiles = Get-ChildItem -Path $basePath -Include "*.go" -Recurse | Where-Object {
    $_.Name -notlike "*_test.go" -and 
    $_.Name -notlike "*.pb.go" -and
    $_.Name -notlike "fix.py"
}

$updated = 0
$skipped = 0

foreach ($file in $goFiles) {
    $content = Get-Content $file.FullName -Raw -ErrorAction SilentlyContinue
    if (-not $content) { continue }
    
    # Skip if header already exists
    if ($content -match "PROPRIETARY - AegisGate Security") {
        $skipped++
        continue
    }
    
    # Add header
    $newContent = $header + $content
    Set-Content -Path $file.FullName -Value $newContent -NoNewline
    $updated++
    Write-Host "[$updated] Updated: $($file.FullName.Replace($basePath, ''))"
}

Write-Host ""
Write-Host "========================================" 
Write-Host "SUMMARY"
Write-Host "========================================"
Write-Host "Files updated: $updated"
Write-Host "Files skipped (already had header): $skipped"
Write-Host "Total files processed: $($updated + $skipped)"