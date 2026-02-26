# Verify local Replay stack health and core API flows
$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
Set-Location $Root

$ApiBase = if ($env:API_BASE_URL) { $env:API_BASE_URL } else { "http://localhost:8080" }
$WebBase = if ($env:WEB_BASE_URL) { $env:WEB_BASE_URL } else { "http://localhost:5173" }
$Results = @()

function Record($name, $ok, $detail) {
    $script:Results += [pscustomobject]@{ Check = $name; Pass = $ok; Detail = $detail }
    $icon = if ($ok) { "PASS" } else { "FAIL" }
    Write-Host "[$icon] $name - $detail"
}

Write-Host "==> Dependency health"
try {
    $h = Invoke-RestMethod -Uri "$ApiBase/health" -TimeoutSec 5
    Record "api-health" ($h.status -eq "ok") ($h | ConvertTo-Json -Compress)
} catch {
    Record "api-health" $false $_.Exception.Message
}

try {
    $ing = Invoke-RestMethod -Uri "http://localhost:8081/health" -TimeoutSec 5
    Record "ingestion-health" ($ing.status -eq "ok") ($ing | ConvertTo-Json -Compress)
} catch {
    Record "ingestion-health" $false "optional: $($_.Exception.Message)"
}

try {
    $web = Invoke-WebRequest -Uri $WebBase -UseBasicParsing -TimeoutSec 5
    Record "web-ui" ($web.StatusCode -eq 200) "status $($web.StatusCode)"
} catch {
    Record "web-ui" $false $_.Exception.Message
}

Write-Host "==> Auth + incident flow"
$email = "verify-$(Get-Date -Format 'yyyyMMddHHmmss')@example.com"
try {
    $reg = Invoke-RestMethod -Method Post -Uri "$ApiBase/v1/auth/register" -ContentType "application/json" -Body (@{
        email = $email; password = "password-12chars"; org_name = "Verify Org"
    } | ConvertTo-Json)
    Record "auth-register" ($null -ne $reg.access_token) "user $($reg.user.id)"

    $headers = @{ Authorization = "Bearer $($reg.access_token)" }
    $proj = Invoke-RestMethod -Method Post -Uri "$ApiBase/v1/projects" -Headers $headers -ContentType "application/json" -Body (@{ name = "verify-project" } | ConvertTo-Json)
    Record "create-project" ($null -ne $proj.id) "project $($proj.id)"

    $now = (Get-Date).ToUniversalTime()
    $body = @{
        window_start = $now.AddHours(-1).ToString("o")
        window_end = $now.ToString("o")
        topic_filters = @("payments.settlement")
    }
    $inc = Invoke-RestMethod -Method Post -Uri "$ApiBase/v1/projects/$($proj.id)/incidents" -Headers $headers -ContentType "application/json" -Body ($body | ConvertTo-Json)
    Record "create-incident" ($inc.status -eq "ready") "status=$($inc.status)"

    $tl = Invoke-RestMethod -Uri "$ApiBase/v1/incidents/$($inc.id)/timeline" -Headers $headers
    $eventCount = @($tl.timeline.events).Count
    Record "timeline" ($eventCount -gt 0) "events=$eventCount"

    $graph = Invoke-RestMethod -Uri "$ApiBase/v1/incidents/$($inc.id)/graph" -Headers $headers
    Record "graph" (@($graph.nodes).Count -gt 0) "nodes=$(@($graph.nodes).Count)"

    $replay = Invoke-RestMethod -Method Post -Uri "$ApiBase/v1/incidents/$($inc.id)/replays" -Headers $headers -ContentType "application/json" -Body (@{ timing_mode = "strict" } | ConvertTo-Json)
    Record "create-replay" ($replay.id) "replay $($replay.id) status=$($replay.status)"
} catch {
    Record "api-flow" $false $_.Exception.Message
}

Write-Host ""
Write-Host "=== Summary ==="
$Results | Format-Table -AutoSize
$failed = @($Results | Where-Object { -not $_.Pass }).Count
if ($failed -gt 0) { exit 1 }
