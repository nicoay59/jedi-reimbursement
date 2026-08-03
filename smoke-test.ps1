param(
    [string]$FrontendUrl = "http://localhost:8088",
    [string]$ApiBaseUrl = "http://localhost:8088/api/v1",
    [string]$AdminEmail = "admin@jedi.local",
    [string]$AdminPassword = "Admin123!"
)

$ErrorActionPreference = "Stop"

function Write-Step {
    param([string]$Message)
    Write-Host ""
    Write-Host "==> $Message" -ForegroundColor Cyan
}

function Assert-True {
    param(
        [bool]$Condition,
        [string]$Message
    )

    if (-not $Condition) {
        throw $Message
    }
}

try {
    Write-Step "Memeriksa frontend"
    $frontend = Invoke-WebRequest `
        -Uri $FrontendUrl `
        -Method Get `
        -UseBasicParsing `
        -TimeoutSec 15

    Assert-True `
        ($frontend.StatusCode -eq 200) `
        "Frontend tidak mengembalikan HTTP 200"

    Write-Host "[OK] Frontend dapat diakses" -ForegroundColor Green

    Write-Step "Memeriksa health API"
    $health = Invoke-RestMethod `
        -Uri "$ApiBaseUrl/health" `
        -Method Get `
        -TimeoutSec 15

    Assert-True `
        ($health.success -eq $true) `
        "Health API tidak sehat"

    Assert-True `
        ($health.data.database -eq "connected") `
        "Database tidak terhubung"

    Write-Host "[OK] API dan database sehat" -ForegroundColor Green

    Write-Step "Login administrator"
    $loginBody = @{
        email = $AdminEmail
        password = $AdminPassword
    } | ConvertTo-Json

    $login = Invoke-RestMethod `
        -Uri "$ApiBaseUrl/auth/login" `
        -Method Post `
        -ContentType "application/json" `
        -Body $loginBody `
        -TimeoutSec 15

    $token = $login.data.access_token
    Assert-True `
        (-not [string]::IsNullOrWhiteSpace($token)) `
        "Access token tidak tersedia"

    $headers = @{
        Authorization = "Bearer $token"
    }

    Write-Host "[OK] Login administrator berhasil" -ForegroundColor Green

    Write-Step "Memeriksa profil pengguna"
    $profile = Invoke-RestMethod `
        -Uri "$ApiBaseUrl/auth/me" `
        -Method Get `
        -Headers $headers `
        -TimeoutSec 15

    Assert-True `
        ($profile.data.role -eq "ADMIN") `
        "Role pengguna bukan ADMIN"

    Write-Host "[OK] Profil administrator valid" -ForegroundColor Green

    Write-Step "Memeriksa dashboard"
    $dashboard = Invoke-RestMethod `
        -Uri "$ApiBaseUrl/admin/dashboard" `
        -Method Get `
        -Headers $headers `
        -TimeoutSec 15

    Assert-True `
        ($dashboard.success -eq $true) `
        "Dashboard tidak dapat dimuat"

    Write-Host "[OK] Dashboard administrator dapat dimuat" -ForegroundColor Green

    Write-Step "Memeriksa laporan"
    $reports = Invoke-RestMethod `
        -Uri "$ApiBaseUrl/admin/reports?page=1&limit=10&type=ALL&status=ALL" `
        -Method Get `
        -Headers $headers `
        -TimeoutSec 15

    Assert-True `
        ($reports.success -eq $true) `
        "Laporan tidak dapat dimuat"

    Write-Host "[OK] Laporan administrator dapat dimuat" -ForegroundColor Green

    Write-Host ""
    Write-Host "============================================" -ForegroundColor Green
    Write-Host "SMOKE TEST BERHASIL" -ForegroundColor Green
    Write-Host "Frontend, API, database, login, dashboard," -ForegroundColor Green
    Write-Host "dan laporan berfungsi." -ForegroundColor Green
    Write-Host "============================================" -ForegroundColor Green
    exit 0
}
catch {
    Write-Host ""
    Write-Host "============================================" -ForegroundColor Red
    Write-Host "SMOKE TEST GAGAL" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    Write-Host "============================================" -ForegroundColor Red
    exit 1
}
