Continue = "Stop"
Set-Location "C:\Users\Administrator\Desktop\Testing\aegisgate"

# Fix handlers.go
 = git show d6b3fa3:pkg/auth/handlers.go
 =  -split "
" | Select-Object -First 76
 | Out-File -FilePath "pkg/auth/handlers.go" -Encoding utf8
Write-Host "handlers.go fixed"

# Fix atlas.go  
 = git show d6b3fa3:pkg/compliance/atlas.go
 =  -split "
" | Select-Object -First 81
 | Out-File -FilePath "pkg/compliance/atlas.go" -Encoding utf8
Write-Host "atlas.go fixed"
