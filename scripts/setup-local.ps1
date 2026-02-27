# Deterministic first-run local setup (Windows PowerShell)
$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
Set-Location $Root

Write-Host "==> Checking prerequisites..."
foreach ($cmd in @("docker", "go", "pnpm")) {
    if (-not (Get-Command $cmd -ErrorAction SilentlyContinue)) {
        throw "Missing required command: $cmd"
    }
}

if (-not (Test-Path ".env.local")) {
    Copy-Item ".env.local.example" ".env.local"
    Write-Host "Created .env.local from .env.local.example"
}

Write-Host "==> Starting Docker dependencies..."
docker compose -f infra/docker-compose/docker-compose.yml up -d

Write-Host "==> Waiting for Postgres..."
$ready = $false
for ($i = 0; $i -lt 30; $i++) {
    docker compose -f infra/docker-compose/docker-compose.yml exec -T postgres pg_isready -U replay 2>$null
    if ($LASTEXITCODE -eq 0) { $ready = $true; break }
    Start-Sleep -Seconds 2
}
if (-not $ready) { throw "Postgres did not become ready in time" }

Write-Host "==> Applying API migrations..."
$env:DATABASE_URL = "postgres://replay:replay@localhost:15433/replay?sslmode=disable"
Push-Location apps/api-gateway
go run ./cmd/migrate -dir ./migrations
Pop-Location

Write-Host "==> Installing web dependencies..."
pnpm install --filter @replay/web

Write-Host "==> Setup complete. Next:"
Write-Host "  1. Load env:  Get-Content .env.local | ForEach-Object { if (`$_ -match '^([^#][^=]+)=(.*)$') { Set-Item -Path env:(`$matches[1].Trim()) -Value `$matches[2].Trim() } }"
Write-Host "  2. API:       go run ./apps/api-gateway/cmd/server"
Write-Host "  3. Ingestion: go run ./apps/ingestion/cmd/server  (optional)"
Write-Host "  4. Web:       pnpm --filter @replay/web dev"
Write-Host "  5. Verify:    ./scripts/verify-local.ps1"
