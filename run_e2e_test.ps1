# Encrypted Channels E2E Test
param(
    [switch]$SkipBuild,
    [switch]$OnlyPush
)

$ErrorActionPreference = "Continue"

Write-Host ""
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host " Encrypted Channels E2E Test" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

# Only push mode
if ($OnlyPush) {
    Write-Host "Push only mode" -ForegroundColor Yellow
    Write-Host ""
    
    try {
        $null = Invoke-WebRequest -Uri "http://localhost:8070/health" -Method GET -TimeoutSec 2
        Write-Host "Server is running" -ForegroundColor Green
    }
    catch {
        Write-Host "Server is not running" -ForegroundColor Red
        Write-Host "Please start: .\devinggo.exe" -ForegroundColor Yellow
        exit 1
    }
    
    Write-Host ""
    Write-Host "Pushing test message..." -ForegroundColor Yellow
    
    $timestamp = (Get-Date).ToString("yyyy-MM-ddTHH:mm:sszzz")
    $testData = @{
        type = "encrypted-test"
        message = "E2E test message"
        amount = 99999.99
        timestamp = $timestamp
        test_id = [guid]::NewGuid().ToString()
    }
    
    $body = @{
        name = "encrypted-message"
        channel = "private-encrypted-secure"
        data = ($testData | ConvertTo-Json -Compress)
    } | ConvertTo-Json
    
    try {
        $response = Invoke-RestMethod -Uri "http://localhost:8070/apps/devinggo-app-id/events" -Method POST -Body $body -ContentType "application/json"
        Write-Host ""
        Write-Host "Push successful!" -ForegroundColor Green
        Write-Host "Message:" -ForegroundColor Cyan
        Write-Host ($testData | ConvertTo-Json) -ForegroundColor White
        Write-Host ""
        Write-Host "Check browser for the message" -ForegroundColor Yellow
    }
    catch {
        Write-Host ""
        Write-Host "Push failed" -ForegroundColor Red
        Write-Host $_.Exception.Message -ForegroundColor Red
    }
    
    exit 0
}

# Full test flow

# Step 1: Build
Write-Host "[1/5] Checking executable..." -ForegroundColor Yellow

if (-not (Test-Path ".\devinggo.exe") -or -not $SkipBuild) {
    Write-Host "Building server..." -ForegroundColor Gray
    $buildOutput = go build -o devinggo.exe .\main.go 2>&1
    
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Build failed" -ForegroundColor Red
        Write-Host $buildOutput -ForegroundColor Red
        exit 1
    }
    Write-Host "Build successful" -ForegroundColor Green
}
else {
    Write-Host "Executable exists" -ForegroundColor Green
}

# Step 2: Cleanup
Write-Host ""
Write-Host "[2/5] Cleaning up..." -ForegroundColor Yellow

$oldProcesses = Get-Process -Name "devinggo" -ErrorAction SilentlyContinue
if ($oldProcesses) {
    $oldProcesses | Stop-Process -Force
    Start-Sleep -Seconds 2
    Write-Host "Cleaned old processes" -ForegroundColor Green
}
else {
    Write-Host "No cleanup needed" -ForegroundColor Green
}

$portsToKill = Get-NetTCPConnection -LocalPort 8070 -ErrorAction SilentlyContinue
if ($portsToKill) {
    foreach ($p in $portsToKill) {
        Stop-Process -Id $p.OwningProcess -Force -ErrorAction SilentlyContinue
    }
    Start-Sleep -Seconds 2
    Write-Host "Released port 8070" -ForegroundColor Green
}

# Step 3: Start server
Write-Host ""
Write-Host "[3/5] Starting server..." -ForegroundColor Yellow

$serverJob = Start-Job -ScriptBlock {
    Set-Location "E:\code\devinggo-light"
    .\devinggo.exe 2>&1
}

Write-Host "Waiting for server..." -ForegroundColor Gray
$retries = 0
$maxRetries = 15
$serverReady = $false

while ($retries -lt $maxRetries) {
    Start-Sleep -Seconds 2
    $retries++
    
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:8070/health" -Method GET -TimeoutSec 2 -ErrorAction Stop
        
        if ($response.StatusCode -eq 200) {
            $serverReady = $true
            Write-Host "Server started" -ForegroundColor Green
            break
        }
    }
    catch {
        Write-Host "." -NoNewline -ForegroundColor Gray
    }
}

if (-not $serverReady) {
    Write-Host ""
    Write-Host "Server failed to start" -ForegroundColor Red
    Write-Host "Server logs:" -ForegroundColor Yellow
    Receive-Job $serverJob | Select-Object -First 30
    Stop-Job $serverJob
    Remove-Job $serverJob
    exit 1
}

# Step 4: Open browser
Write-Host ""
Write-Host "[4/5] Opening test page..." -ForegroundColor Yellow

Start-Process "http://localhost:8070/pusher-test.html"
Write-Host "Browser opened" -ForegroundColor Green

# Step 5: Push messages
Write-Host ""
Write-Host "[5/5] Preparing to push..." -ForegroundColor Yellow
Write-Host ""
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host " Please complete these steps:" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "  1. Click Connect to Pusher" -ForegroundColor White
Write-Host "  2. Find Encrypted Channel section" -ForegroundColor White
Write-Host "  3. Click Subscribe button" -ForegroundColor White
Write-Host "  4. See subscription success" -ForegroundColor White
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Press any key to push test messages..." -ForegroundColor Yellow
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

Write-Host ""
Write-Host "Pushing test messages..." -ForegroundColor Yellow
Write-Host ""

$currentTime = (Get-Date).ToString("yyyy-MM-ddTHH:mm:sszzz")

$testMessages = @(
    @{
        name = "user-notification"
        data = @{
            type = "notification"
            title = "System Notification"
            message = "You have a new message"
            timestamp = $currentTime
        }
    },
    @{
        name = "payment-alert"
        data = @{
            type = "payment"
            message = "Payment received"
            amount = 5000.00
            from = "Zhang San"
            card = "6222 **** **** 1234"
            timestamp = $currentTime
        }
    },
    @{
        name = "security-warning"
        data = @{
            type = "security"
            level = "high"
            message = "Remote location login detected"
            ip = "123.45.67.89"
            location = "Beijing"
            timestamp = $currentTime
        }
    }
)

$successCount = 0
foreach ($msg in $testMessages) {
    $body = @{
        name = $msg.name
        channel = "private-encrypted-secure"
        data = ($msg.data | ConvertTo-Json -Compress)
    } | ConvertTo-Json
    
    try {
        $null = Invoke-RestMethod -Uri "http://localhost:8070/apps/devinggo-app-id/events" -Method POST -Body $body -ContentType "application/json"
        
        Write-Host "[$($msg.name)] Push successful" -ForegroundColor Green
        $successCount++
        Start-Sleep -Milliseconds 500
    }
    catch {
        Write-Host "[$($msg.name)] Push failed" -ForegroundColor Red
        Write-Host $_.Exception.Message -ForegroundColor Red
    }
}

# Summary
Write-Host ""
Write-Host "=====================================" -ForegroundColor Green
Write-Host " Test completed" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Green
Write-Host "  Successful: $successCount / $($testMessages.Count) messages" -ForegroundColor White
Write-Host "=====================================" -ForegroundColor Green
Write-Host ""
Write-Host "Check browser for the messages" -ForegroundColor Yellow
Write-Host "Messages should be automatically decrypted" -ForegroundColor White
Write-Host ""
Write-Host "Tips:" -ForegroundColor Cyan
Write-Host "  - Server running in background" -ForegroundColor Gray
Write-Host "  - Manual push: .\test_encrypted_push.ps1" -ForegroundColor Gray
Write-Host ""
Write-Host "Press any key to stop server and exit..." -ForegroundColor Yellow
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

# Cleanup
Write-Host ""
Write-Host "Cleaning up..." -ForegroundColor Yellow
Stop-Job $serverJob -ErrorAction SilentlyContinue
Remove-Job $serverJob -ErrorAction SilentlyContinue

$processes = Get-Process -Name "devinggo" -ErrorAction SilentlyContinue
if ($processes) {
    $processes | Stop-Process -Force
    Write-Host "Server stopped" -ForegroundColor Green
}

Write-Host ""
Write-Host "Test completed!" -ForegroundColor Green
