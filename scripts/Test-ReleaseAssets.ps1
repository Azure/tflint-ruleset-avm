#requires -Version 7.2

[CmdletBinding()]
param(
    [Parameter(Mandatory)]
    [string] $ReleaseDirectory
)

$ErrorActionPreference = 'Stop'
$pluginName = 'tflint-ruleset-avm'
$targetsPath = Join-Path $PSScriptRoot '..' 'release' 'targets.json'
$targets = Get-Content -LiteralPath $targetsPath -Raw | ConvertFrom-Json
$releaseRoot = (Resolve-Path -LiteralPath $ReleaseDirectory).Path

$expectedArchives = @(
    $targets | ForEach-Object {
        "$pluginName`_$($_.goos)_$($_.goarch).zip"
    } | Sort-Object
)
$expectedFiles = @($expectedArchives + 'checksums.txt' | Sort-Object)
$actualFiles = @(
    Get-ChildItem -LiteralPath $releaseRoot -File |
        Select-Object -ExpandProperty Name |
        Sort-Object
)
$differences = @(Compare-Object -ReferenceObject $expectedFiles -DifferenceObject $actualFiles)
if ($differences.Count -ne 0) {
    throw "Release files do not match the expected asset set: $($differences | Out-String)"
}

$checksumPath = Join-Path $releaseRoot 'checksums.txt'
$manifest = @{}
$checksumLines = @(Get-Content -LiteralPath $checksumPath)
foreach ($line in $checksumLines) {
    if ($line -notmatch '^([0-9a-f]{64})  (tflint-ruleset-avm_[a-z0-9]+_[a-z0-9]+\.zip)$') {
        throw "Invalid checksums.txt line: $line"
    }
    if ($manifest.ContainsKey($Matches[2])) {
        throw "Duplicate checksum entry: $($Matches[2])"
    }
    $manifest[$Matches[2]] = $Matches[1]
}

if ($manifest.Count -ne $expectedArchives.Count) {
    throw "checksums.txt contains $($manifest.Count) entries; expected $($expectedArchives.Count)"
}

Add-Type -AssemblyName System.IO.Compression.FileSystem
foreach ($target in $targets) {
    $archiveName = "$pluginName`_$($target.goos)_$($target.goarch).zip"
    if (-not $manifest.ContainsKey($archiveName)) {
        throw "checksums.txt is missing $archiveName"
    }

    $archivePath = Join-Path $releaseRoot $archiveName
    $actualHash = (Get-FileHash -LiteralPath $archivePath -Algorithm SHA256).Hash.ToLowerInvariant()
    if ($manifest[$archiveName] -ne $actualHash) {
        throw "Checksum mismatch for $archiveName"
    }

    $expectedBinary = if ($target.goos -eq 'windows') { "$pluginName.exe" } else { $pluginName }
    $zip = [System.IO.Compression.ZipFile]::OpenRead($archivePath)
    try {
        $entries = @($zip.Entries)
        if ($entries.Count -ne 1 -or $entries[0].FullName -ne $expectedBinary) {
            throw "$archiveName must contain only $expectedBinary"
        }
    }
    finally {
        $zip.Dispose()
    }
}

Write-Host "Validated $($expectedArchives.Count) TFLint release archives and checksums.txt."
