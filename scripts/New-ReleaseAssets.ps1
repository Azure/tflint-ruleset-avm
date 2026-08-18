#requires -Version 7.2

[CmdletBinding()]
param(
    [Parameter(Mandatory)]
    [string] $BinaryDirectory,

    [Parameter(Mandatory)]
    [string] $OutputDirectory
)

$ErrorActionPreference = 'Stop'
$pluginName = 'tflint-ruleset-avm'
$targetsPath = Join-Path $PSScriptRoot '..' 'release' 'targets.json'
$targets = Get-Content -LiteralPath $targetsPath -Raw | ConvertFrom-Json
$binaryRoot = (Resolve-Path -LiteralPath $BinaryDirectory).Path

if (-not (Test-Path -LiteralPath $OutputDirectory)) {
    $null = New-Item -ItemType Directory -Path $OutputDirectory
}
$outputRoot = (Resolve-Path -LiteralPath $OutputDirectory).Path

$archives = foreach ($target in $targets) {
    $extension = if ($target.goos -eq 'windows') { '.exe' } else { '' }
    $binaryName = "$pluginName$extension"
    $binaryPath = Join-Path $binaryRoot "$($target.goos)-$($target.goarch)" $binaryName
    if (-not (Test-Path -LiteralPath $binaryPath -PathType Leaf)) {
        throw "Release binary not found: $binaryPath"
    }

    $archiveName = "$pluginName`_$($target.goos)_$($target.goarch).zip"
    $archivePath = Join-Path $outputRoot $archiveName
    if (Test-Path -LiteralPath $archivePath) {
        Remove-Item -LiteralPath $archivePath -Force
    }
    Compress-Archive -LiteralPath $binaryPath -DestinationPath $archivePath -CompressionLevel Optimal
    Get-Item -LiteralPath $archivePath
}

$checksumLines = foreach ($archive in $archives | Sort-Object Name) {
    $hash = (Get-FileHash -LiteralPath $archive.FullName -Algorithm SHA256).Hash.ToLowerInvariant()
    "$hash  $($archive.Name)"
}
$checksumPath = Join-Path $outputRoot 'checksums.txt'
$checksumLines | Set-Content -LiteralPath $checksumPath -Encoding utf8NoBOM

& (Join-Path $PSScriptRoot 'Test-ReleaseAssets.ps1') -ReleaseDirectory $outputRoot
