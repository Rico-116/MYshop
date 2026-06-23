param(
  [string]$BaseUrl = "http://localhost:8080",
  [int[]]$HotspotVus = @(10, 20, 50, 100, 200, 400, 800, 1000, 1200, 1500),
  [int[]]$MixVus = @(10, 20, 50, 100, 200, 400, 600, 800),
  [string]$Duration = "60s",
  [string]$WarmupDuration = "20s",
  [int]$ProductId = 1,
  [int]$CategoryId = 1,
  [string]$Keyword = "",
  [int]$PageSize = 10,
  [string]$K6Path = "k6",
  [switch]$HotspotOnly,
  [switch]$MixOnly,
  [switch]$ProbeOnly
)

$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$repoRoot = Split-Path -Parent (Split-Path -Parent $scriptDir)
$k6Script = Join-Path $scriptDir "myshop-report-style-k6.js"
$resultDir = Join-Path $repoRoot "tests\load\results\myshop-report-style"

New-Item -ItemType Directory -Force -Path $resultDir | Out-Null

function Get-MetricValue($metrics, $name, $field) {
  if (-not $metrics.$name) { return $null }
  if (-not $metrics.$name.values) { return $null }

  $property = $metrics.$name.values.PSObject.Properties[$field]
  if (-not $property) { return $null }

  return $property.Value
}

function Format-Number($value, $digits = 2) {
  if ($null -eq $value) { return "" }
  return [Math]::Round([double]$value, $digits)
}

function Format-Percent($value) {
  if ($null -eq $value) { return "" }
  return "{0}%" -f (Format-Number ([double]$value * 100) 3)
}

function Run-One($mode, $vus, $duration, $scene) {
  $summaryFile = Join-Path $resultDir ("summary-{0}-{1}vus.json" -f $mode, $vus)

  Write-Host ""
  Write-Host ("=== {0}: {1} VUs / {2} ===" -f $scene, $vus, $duration)

  $env:BASE_URL = $BaseUrl
  $env:MODE = $mode
  $env:VUS = "$vus"
  $env:DURATION = $duration
  $env:PRODUCT_ID = "$ProductId"
  $env:CATEGORY_ID = "$CategoryId"
  $env:KEYWORD = $Keyword
  $env:PAGE_SIZE = "$PageSize"
  $env:SUMMARY_FILE = $summaryFile

  & $K6Path run $k6Script | Out-Host
  $exitCode = $LASTEXITCODE

  if (-not (Test-Path -LiteralPath $summaryFile)) {
    throw "k6 did not write summary file: $summaryFile"
  }

  $summary = Get-Content -LiteralPath $summaryFile -Raw | ConvertFrom-Json
  $metrics = $summary.metrics

  $httpFail = Get-MetricValue $metrics "http_req_failed" "rate"
  $apiFail = Get-MetricValue $metrics "myshop_api_failed" "rate"
  $qps = Get-MetricValue $metrics "http_reqs" "rate"
  $avg = Get-MetricValue $metrics "http_req_duration" "avg"
  $p95 = Get-MetricValue $metrics "http_req_duration" "p(95)"
  $p99 = Get-MetricValue $metrics "http_req_duration" "p(99)"
  $max = Get-MetricValue $metrics "http_req_duration" "max"

  return [PSCustomObject]@{
    Scene = $scene
    Mode = $mode
    VUS = $vus
    QPS = Format-Number $qps
    HttpFailRate = Format-Percent $httpFail
    ApiFailRate = Format-Percent $apiFail
    AvgRTms = Format-Number $avg
    P95ms = Format-Number $p95
    P99ms = Format-Number $p99
    MaxRTms = Format-Number $max
    ExitCode = $exitCode
    Summary = $summaryFile
  }
}

function Is-Bad-Result($row) {
  $httpFail = [double](($row.HttpFailRate -replace "%", "") -replace "^$", "0")
  $apiFail = [double](($row.ApiFailRate -replace "%", "") -replace "^$", "0")
  $p99 = [double](($row.P99ms -replace "^$", "0"))

  return ($httpFail -gt 5) -or ($apiFail -gt 5) -or ($p99 -gt 5000) -or ($row.ExitCode -ne 0)
}

function ConvertTo-MarkdownTable($rows) {
  $lines = @()
  $lines += "| Scene | VUS | QPS | HTTP Fail | API Fail | Avg RT | P95 | P99 | Max RT |"
  $lines += "|---|---:|---:|---:|---:|---:|---:|---:|---:|"

  foreach ($row in $rows) {
    $lines += "| {0} | {1} | {2} | {3} | {4} | {5}ms | {6}ms | {7}ms | {8}ms |" -f `
      $row.Scene, $row.VUS, $row.QPS, $row.HttpFailRate, $row.ApiFailRate, `
      $row.AvgRTms, $row.P95ms, $row.P99ms, $row.MaxRTms
  }

  return ($lines -join [Environment]::NewLine)
}

$rows = @()

Write-Host ""
Write-Host "=== Warmup ==="
if (-not $MixOnly) {
  [void](Run-One "hotspot" 50 $WarmupDuration "warmup-hotspot")
}
if (-not $HotspotOnly) {
  [void](Run-One "mix" 50 $WarmupDuration "warmup-mix")
}

if ($ProbeOnly) {
  foreach ($mode in @("hotspot", "products", "hot", "category", "mix")) {
    $row = Run-One $mode 50 $Duration ("probe-{0}" -f $mode)
    $rows += $row
  }
} else {
  if (-not $MixOnly) {
    foreach ($vus in $HotspotVus) {
      $row = Run-One "hotspot" $vus $Duration "hotspot-product-detail"
      $rows += $row

      if (Is-Bad-Result $row) {
        Write-Host "Stop hotspot test: this level already has obvious errors or very high P99."
        break
      }
    }
  }

  if (-not $HotspotOnly) {
    foreach ($vus in $MixVus) {
      $row = Run-One "mix" $vus $Duration "public-read-mix"
      $rows += $row

      if (Is-Bad-Result $row) {
        Write-Host "Stop mix test: this level already has obvious errors or very high P99."
        break
      }
    }
  }
}

$csvPath = Join-Path $resultDir "myshop-report-summary.csv"
$mdPath = Join-Path $resultDir "myshop-report-summary.md"

$rows | Export-Csv -LiteralPath $csvPath -NoTypeInformation -Encoding UTF8
ConvertTo-MarkdownTable $rows | Set-Content -LiteralPath $mdPath -Encoding UTF8

Write-Host ""
Write-Host "=== MyShop Report Style Summary ==="
$rows | Format-Table -AutoSize

Write-Host ""
Write-Host ("CSV summary: {0}" -f $csvPath)
Write-Host ("Markdown summary: {0}" -f $mdPath)
