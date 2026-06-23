param(
  [string]$BaseUrl = "http://localhost:8080",
  [int[]]$Levels = @(50, 100, 150, 200, 300, 400, 500),
  [string]$Duration = "2m",
  [int]$ProductId = 1,
  [string]$Username = "",
  [string]$Password = "",
  [double]$Sleep = 0.2,
  [string]$K6Path = "k6"
)

$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$repoRoot = Split-Path -Parent (Split-Path -Parent $scriptDir)
$k6Script = Join-Path $scriptDir "max-concurrency-k6.js"
$resultDir = Join-Path $repoRoot "tests\load\results"

New-Item -ItemType Directory -Force -Path $resultDir | Out-Null

function Get-MetricValue($metrics, $name, $field) {
  if (-not $metrics.$name) { return $null }
  if (-not $metrics.$name.values) { return $null }
  return $metrics.$name.values.$field
}

function Format-Number($value, $digits = 2) {
  if ($null -eq $value) { return "" }
  return [Math]::Round([double]$value, $digits)
}

$rows = @()

foreach ($level in $Levels) {
  $summaryFile = Join-Path $resultDir ("summary-vus-{0}.json" -f $level)

  Write-Host ""
  Write-Host "=== Testing $level concurrent users for $Duration ==="

  $env:BASE_URL = $BaseUrl
  $env:VUS = "$level"
  $env:DURATION = $Duration
  $env:PRODUCT_ID = "$ProductId"
  $env:USERNAME = $Username
  $env:PASSWORD = $Password
  $env:SLEEP = "$Sleep"
  $env:SUMMARY_FILE = $summaryFile

  & $K6Path run $k6Script | Out-Null
  $exitCode = $LASTEXITCODE

  if (-not (Test-Path $summaryFile)) {
    throw "k6 did not write summary file: $summaryFile"
  }

  $summary = Get-Content -LiteralPath $summaryFile -Raw | ConvertFrom-Json
  $metrics = $summary.metrics

  $failRate = Get-MetricValue $metrics "http_req_failed" "rate"
  $apiFailRate = Get-MetricValue $metrics "myshop_api_failed" "rate"
  $rps = Get-MetricValue $metrics "http_reqs" "rate"
  $p95 = Get-MetricValue $metrics "http_req_duration" "p(95)"
  $p99 = Get-MetricValue $metrics "http_req_duration" "p(99)"
  $p999 = Get-MetricValue $metrics "http_req_duration" "p(99.9)"
  $max = Get-MetricValue $metrics "http_req_duration" "max"

  $passed = ($exitCode -eq 0) -and
    (($failRate -eq $null) -or ([double]$failRate -lt 0.01)) -and
    (($apiFailRate -eq $null) -or ([double]$apiFailRate -lt 0.01)) -and
    (($p95 -eq $null) -or ([double]$p95 -lt 1000)) -and
    (($p99 -eq $null) -or ([double]$p99 -lt 2000))

  $rows += [PSCustomObject]@{
    VUS = $level
    RPS = Format-Number $rps
    FailRate = if ($null -eq $failRate) { "" } else { "{0}%" -f (Format-Number ([double]$failRate * 100) 3) }
    ApiFailRate = if ($null -eq $apiFailRate) { "" } else { "{0}%" -f (Format-Number ([double]$apiFailRate * 100) 3) }
    P95ms = Format-Number $p95
    P99ms = Format-Number $p99
    P999ms = Format-Number $p999
    Maxms = Format-Number $max
    Passed = $passed
  }

  if (-not $passed) {
    Write-Host "Stop: $level concurrent users crossed the limit."
    break
  }
}

$csvPath = Join-Path $resultDir "max-concurrency-summary.csv"
$rows | Export-Csv -LiteralPath $csvPath -NoTypeInformation -Encoding UTF8

Write-Host ""
Write-Host "=== Result ==="
$rows | Format-Table -AutoSize

$stable = $rows | Where-Object { $_.Passed -eq $true } | Select-Object -Last 1
if ($stable) {
  Write-Host ""
  Write-Host ("Stable max concurrency by current rule: {0} users" -f $stable.VUS)
} else {
  Write-Host ""
  Write-Host "No tested concurrency level passed the current rule."
}

Write-Host ("CSV summary: {0}" -f $csvPath)
