
# PowerShell script to fix workflow files
$workflowDir = "C:\Users\Administrator\Desktop\Testing\AegisGate\.github\workflows"

# Process each workflow file
$files = Get-ChildItem -Path $workflowDir -Filter "*.yml" | Where-Object { $_.Name -ne "ci.yml" }

foreach ($file in $files) {
    Write-Host "Processing" $file.Name
    
    # Read file content
    $content = Get-Content -Path $file.FullName -Raw
    
    # Split by the workflow header "name: "
    $parts = $content -split "(?<=
)name: "
    
    if ($parts.Count -gt 2) {
        # Multiple workflows found, keep only the first
        Write-Host "Found duplicate workflows in" $file.Name
        
        # Keep only the first workflow
        $firstWorkflow = "name: " + $parts[1]
        
        # Write back the fixed content
        Set-Content -Path $file.FullName -Value $firstWorkflow -NoNewline
        
        Write-Host "Fixed" $file.Name
    } else {
        Write-Host $file.Name "looks OK (no duplicates)"
    }
}

Write-Host "Workflow file fixing complete"
