# Download golangci-lint v1.64.0
$ErrorActionPreference = "Continue"
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$zipUrl = "https://github.com/golangci/golangci-lint/releases/download/v1.64.0/golangci-lint-1.64.0-windows-amd64.zip"
$zipPath = "C:\temp\golangci-v1.64.0.zip"
$extractPath = "C:\temp\golangci-v1.64.0"

Write-Host "Downloading golangci-lint v1.64.0..."
Invoke-WebRequest -Uri $zipUrl -OutFile $zipPath -UseBasicParsing

Write-Host "Extracting..."
Expand-Archive -Path $zipPath -DestinationPath $extractPath -Force

Write-Host "Copying to GOBIN..."
Copy-Item "$extractPath\golangci-lint-1.64.0-windows-amd64\golangci-lint.exe" "C:\Users\Administrator\go\bin\golangci-lint.exe" -Force

Write-Host "Verifying..."
& "C:\Users\Administrator\go\bin\golangci-lint.exe" --version
