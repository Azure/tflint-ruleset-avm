#requires -Version 7.2

$ErrorActionPreference = 'Stop'
$root = Join-Path ([System.IO.Path]::GetTempPath()) "tflint-ruleset-avm-$([guid]::NewGuid())"
$binaryRoot = Join-Path $root 'binaries'
$releaseRoot = Join-Path $root 'release'

try {
    $targetsPath = Join-Path $PSScriptRoot '..' 'release' 'targets.json'
    $targets = Get-Content -LiteralPath $targetsPath -Raw | ConvertFrom-Json
    foreach ($target in $targets) {
        $directory = Join-Path $binaryRoot "$($target.goos)-$($target.goarch)"
        $null = New-Item -ItemType Directory -Path $directory -Force
        $binaryName = if ($target.goos -eq 'windows') {
            'tflint-ruleset-avm.exe'
        }
        else {
            'tflint-ruleset-avm'
        }
        "fixture-$($target.goos)-$($target.goarch)" |
            Set-Content -LiteralPath (Join-Path $directory $binaryName) -Encoding utf8NoBOM
    }

    & (Join-Path $PSScriptRoot 'New-ReleaseAssets.ps1') `
        -BinaryDirectory $binaryRoot `
        -OutputDirectory $releaseRoot

    $checksumPath = Join-Path $releaseRoot 'checksums.txt'
    $originalChecksums = Get-Content -LiteralPath $checksumPath -Raw
    $tamperedChecksums = $originalChecksums -replace '^[0-9a-f]', '0'
    if ($tamperedChecksums -eq $originalChecksums) {
        $tamperedChecksums = $originalChecksums -replace '^[0-9a-f]', '1'
    }
    $tamperedChecksums | Set-Content -LiteralPath $checksumPath -Encoding utf8NoBOM -NoNewline

    $failedAsExpected = $false
    try {
        & (Join-Path $PSScriptRoot 'Test-ReleaseAssets.ps1') -ReleaseDirectory $releaseRoot
    }
    catch {
        $failedAsExpected = $true
    }
    if (-not $failedAsExpected) {
        throw 'Release validation accepted a tampered checksum.'
    }

    $originalChecksums | Set-Content -LiteralPath $checksumPath -Encoding utf8NoBOM -NoNewline
    'unexpected' | Set-Content -LiteralPath (Join-Path $releaseRoot 'unexpected.txt')
    $failedAsExpected = $false
    try {
        & (Join-Path $PSScriptRoot 'Test-ReleaseAssets.ps1') -ReleaseDirectory $releaseRoot
    }
    catch {
        $failedAsExpected = $true
    }
    if (-not $failedAsExpected) {
        throw 'Release validation accepted an unexpected asset.'
    }

    Write-Host 'Release asset tests passed.'
}
finally {
    if (Test-Path -LiteralPath $root) {
        Remove-Item -LiteralPath $root -Recurse -Force
    }
}
