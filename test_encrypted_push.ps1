# Server Push Test with Pusher Authentication
# Tests Encrypted Channel with server-side encryption
# Server encrypts data using saved shared_secret

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  Encrypted Channel Push Test" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# ===== Configuration =====
$appId = "devinggo-app-id"
$appKey = "devinggo-app-key"
$appSecret = "devinggo-app-secret-change-me"
$authVersion = "1.0"

# ===== Step 1: Check Server =====
Write-Host "[1/5] Checking server..." -ForegroundColor Yellow

try {
    $response = Invoke-WebRequest -Uri "http://localhost:8070/health" -Method GET -TimeoutSec 3
    Write-Host "Server is running" -ForegroundColor Green
}
catch {
    Write-Host "Server is not running" -ForegroundColor Red
    Write-Host "Please start: .\devinggo.exe" -ForegroundColor Yellow
    exit 1
}

# ===== Step 2: User Confirmation =====
Write-Host ""
Write-Host "[2/5] Preparing..." -ForegroundColor Yellow
Write-Host ""
Write-Host "Please confirm:" -ForegroundColor Yellow
Write-Host "  1. Browser opened pusher-test.html" -ForegroundColor White
Write-Host "  2. Connected to Pusher" -ForegroundColor White
Write-Host "  3. Subscribed to 'private-encrypted-secure' channel" -ForegroundColor White
Write-Host ""
Write-Host "Press any key to continue..." -ForegroundColor Cyan
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

# ===== Step 3: Build Message =====
Write-Host ""
Write-Host "[3/5] Building message..." -ForegroundColor Yellow

$timestamp = (Get-Date).ToString("yyyy-MM-ddTHH:mm:sszzz")
$testData = @{
    type = "server-encrypted-test"
    message = "Test encrypted message from server"
    amount = 12345.67
    sensitive = "Secret data: 6222 **** **** 1234"
    timestamp = $timestamp
    test_id = [guid]::NewGuid().ToString()
}

# Use Encrypted Channel - server will encrypt using saved shared_secret
$requestBody = @{
    name = "encrypted-message"
    channels = @("private-encrypted-secure")
    data = ($testData | ConvertTo-Json -Compress)
} | ConvertTo-Json -Compress

Write-Host "Message data:" -ForegroundColor Gray
Write-Host ($testData | ConvertTo-Json) -ForegroundColor Gray

# ===== Step 4: Generate Pusher Signature =====
Write-Host ""
Write-Host "[4/5] Generating Pusher signature..." -ForegroundColor Yellow

# Compute body MD5
$md5 = [System.Security.Cryptography.MD5]::Create()
$bodyBytes = [System.Text.Encoding]::UTF8.GetBytes($requestBody)
$bodyMd5Hash = $md5.ComputeHash($bodyBytes)
$bodyMd5 = [System.BitConverter]::ToString($bodyMd5Hash).Replace("-", "").ToLower()

# Generate timestamp (Unix seconds)
$authTimestamp = [int64](([datetime]::UtcNow) - (Get-Date "1970-01-01")).TotalSeconds

# Build query string (alphabetically sorted)
$queryParams = @{
    auth_key = $appKey
    auth_timestamp = $authTimestamp.ToString()
    auth_version = $authVersion
    body_md5 = $bodyMd5
}

$sortedKeys = $queryParams.Keys | Sort-Object
$queryParts = @()
foreach ($key in $sortedKeys) {
    $queryParts += "$key=$($queryParams[$key])"
}
$queryString = $queryParts -join "&"

# Build string to sign
$method = "POST"
$path = "/apps/$appId/events"
$stringToSign = "$method`n$path`n$queryString"

Write-Host "  Method: $method" -ForegroundColor Gray
Write-Host "  Path: $path" -ForegroundColor Gray
Write-Host "  Body MD5: $bodyMd5" -ForegroundColor Gray
Write-Host "  Timestamp: $authTimestamp" -ForegroundColor Gray
Write-Host "  Query: $queryString" -ForegroundColor Gray

# Compute HMAC-SHA256 signature
$hmacsha = New-Object System.Security.Cryptography.HMACSHA256
$hmacsha.Key = [System.Text.Encoding]::UTF8.GetBytes($appSecret)
$signatureBytes = $hmacsha.ComputeHash([System.Text.Encoding]::UTF8.GetBytes($stringToSign))
$authSignature = [System.BitConverter]::ToString($signatureBytes).Replace("-", "").ToLower()

Write-Host "  Signature: $authSignature" -ForegroundColor Gray

# Build full URL with query parameters
$url = "http://localhost:8070$path`?$queryString&auth_signature=$authSignature"

# ===== Step 5: Push Message =====
Write-Host ""
Write-Host "[5/5] Pushing message..." -ForegroundColor Yellow

try {
    $response = Invoke-RestMethod -Uri $url -Method POST -Body $requestBody -ContentType "application/json"
    
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "  Push Successful!" -ForegroundColor Green
    Write-Host "========================================" -ForegroundColor Green
    Write-Host ""
    Write-Host "Details:" -ForegroundColor Cyan
    Write-Host "  Channel: private-test" -ForegroundColor White
    Write-Host "  Event: server-push-event" -ForegroundColor White
    Write-Host "  Time: $timestamp" -ForegroundColor White
    Write-Host ""
    Write-Host "Check browser console for the message!" -ForegroundColor Yellow
    Write-Host "You should see:" -ForegroundColor Gray
    Write-Host "  📨 收到私有频道事件: server-push-event" -ForegroundColor Gray
    
    if ($response) {
        Write-Host ""
        Write-Host "Server response:" -ForegroundColor Gray
        Write-Host ($response | ConvertTo-Json) -ForegroundColor Gray
    }
}
catch {
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Red
    Write-Host "  Push Failed" -ForegroundColor Red
    Write-Host "========================================" -ForegroundColor Red
    Write-Host ""
    Write-Host "Error:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $responseBody = $reader.ReadToEnd()
        Write-Host ""
        Write-Host "Response body:" -ForegroundColor Red
        Write-Host $responseBody -ForegroundColor Red
    }
    
    Write-Host ""
    Write-Host "Troubleshooting:" -ForegroundColor Yellow
    Write-Host "  1. Check config.yaml pusher.appSecret matches" -ForegroundColor White
    Write-Host "  2. Verify server logs for signature errors" -ForegroundColor White
    Write-Host "  3. Ensure browser subscribed to channel" -ForegroundColor White
    exit 1
}

Write-Host ""
Write-Host "Test completed!" -ForegroundColor Green
Write-Host ""
