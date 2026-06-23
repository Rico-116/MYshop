param(
  [string]$BaseUrl = "http://localhost:8080",
  [int[]]$Targets = @(10, 20, 50, 80, 100, 150, 200, 300, 500, 800, 1000, 2000, 3000, 5000, 8000, 10000),
  [string]$Duration = "2m",
  [string]$Endpoint = "mix",
  [int]$ProductId = 1,
  [int]$PreAllocatedVUs = 0,
  [int]$MaxVUs = 0,
  [string]$K6Path = "k6"
)

$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$repoRoot = Split-Path -Parent (Split-Path -Parent $scriptDir)
$k6Script = Join-Path $scriptDir "max-qps-k6.js"
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

foreach ($target in $Targets) {
  $summaryFile = Join-Path $resultDir ("summary-qps-{0}-{1}.json" -f $Endpoint, $target)

  Write-Host ""
  Write-Host "=== Testing target $target QPS for $Duration, endpoint: $Endpoint ==="

  $env:BASE_URL = $BaseUrl
  $env:TARGET_QPS = "$target"
  $env:DURATION = $Duration
  $env:ENDPOINT = $Endpoint
  $env:PRODUCT_ID = "$ProductId"
  $env:SUMMARY_FILE = $summaryFile

  if ($PreAllocatedVUs -gt 0) {
    $env:PRE_ALLOCATED_VUS = "$PreAllocatedVUs"
  } else {
    Remove-Item Env:\PRE_ALLOCATED_VUS -ErrorAction SilentlyContinue
  }

  if ($MaxVUs -gt 0) {
    $env:MAX_VUS = "$MaxVUs"
  } else {
    Remove-Item Env:\MAX_VUS -ErrorAction SilentlyContinue
  }

  & $K6Path run $k6Script | Out-Null
  $exitCode = $LASTEXITCODE

  if (-not (Test-Path $summaryFile)) {
    throw "k6 did not write summary file: $summaryFile"
  }

  $summary = Get-Content -LiteralPath $summaryFile -Raw | ConvertFrom-Json
  $metrics = $summary.metrics

  $failRate = Get-MetricValue $metrics "http_req_failed" "rate"
  $apiFailRate = Get-MetricValue $metrics "myshop_api_failed" "rate"
  $actualRps = Get-MetricValue $metrics "http_reqs" "rate"
  $dropped = Get-MetricValue $metrics "dropped_iterations" "count"
  $p95 = Get-MetricValue $metrics "http_req_duration" "p(95)"
  $p99 = Get-MetricValue $metrics "http_req_duration" "p(99)"
  $p999 = Get-MetricValue $metrics "http_req_duration" "p(99.9)"
  $max = Get-MetricValue $metrics "http_req_duration" "max"

  $pass = ($exitCode -eq 0) -and
    (($failRate -eq $null) -or ([double]$failRate -lt 0.01)) -and
    (($apiFailRate -eq $null) -or ([double]$apiFailRate -lt 0.01)) -and
    (($dropped -eq $null) -or ([double]$dropped -eq 0)) -and
    (($actualRps -ne $null) -and ([double]$actualRps -ge ($target * 0.95))) -and
    (($p95 -eq $null) -or ([double]$p95 -lt 1000)) -and
    (($p99 -eq $null) -or ([double]$p99 -lt 2000))

  $rows += [PSCustomObject]@{
    TargetQPS = $target
    ActualRPS = Format-Number $actualRps
    Dropped = if ($null -eq $dropped) { 0 } else { [int]$dropped }
    FailRate = if ($null -eq $failRate) { "" } else { "{0}%" -f (Format-Number ([double]$failRate * 100) 3) }
    ApiFailRate = if ($null -eq $apiFailRate) { "" } else { "{0}%" -f (Format-Number ([double]$apiFailRate * 100) 3) }
    P95ms = Format-Number $p95
    P99ms = Format-Number $p99
    P999ms = Format-Number $p999
    Maxms = Format-Number $max
    Passed = $pass
  }

  if (-not $pass) {
    Write-Host "Stop: target $target QPS crossed the limit."
    break
  }
}

$csvPath = Join-Path $resultDir ("max-qps-summary-{0}.csv" -f $Endpoint)
$rows | Export-Csv -LiteralPath $csvPath -NoTypeInformation -Encoding UTF8

Write-Host ""
Write-Host "=== Result ==="
$rows | Format-Table -AutoSize

$stable = $rows | Where-Object { $_.Passed -eq $true } | Select-Object -Last 1
if ($stable) {
  Write-Host ""
  Write-Host ("Stable max QPS by current rule: {0} req/s" -f $stable.TargetQPS)
} else {
  Write-Host ""
  Write-Host "No tested QPS target passed the current rule."
}

Write-Host ("CSV summary: {0}" -f $csvPath)
