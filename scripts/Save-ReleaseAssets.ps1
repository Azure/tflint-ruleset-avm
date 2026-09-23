#requires -Version 7.2

[CmdletBinding()]
param(
    [Parameter(Mandatory)]
    [ValidatePattern('^v[0-9]+\.[0-9]+\.[0-9]+$')]
    [string] $Tag,

    [Parameter(Mandatory)]
    [ValidatePattern('^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$')]
    [string] $Repository,

    [Parameter(Mandatory)]
    [string] $OutputDirectory,

    [long] $ReleaseId
)

$ErrorActionPreference = 'Stop'
if ([string]::IsNullOrWhiteSpace($env:GH_TOKEN)) {
    throw 'GH_TOKEN is required to download release assets.'
}
if ($PSBoundParameters.ContainsKey('ReleaseId') -and $ReleaseId -le 0) {
    throw "Invalid release ID: $ReleaseId"
}

function Get-Release {
    param([string] $Endpoint)

    $json = & gh api $Endpoint
    if ($LASTEXITCODE -ne 0) {
        throw "GitHub API request failed: $Endpoint"
    }
    return ($json -join "`n" | ConvertFrom-Json)
}

if (-not $PSBoundParameters.ContainsKey('ReleaseId')) {
    $byTag = Get-Release "repos/$Repository/releases/tags/$Tag"
    if (-not [long]::TryParse([string] $byTag.id, [ref] $ReleaseId) -or $ReleaseId -le 0) {
        throw "Release $Tag has no valid numeric ID."
    }
}

$release = Get-Release "repos/$Repository/releases/$ReleaseId"
$actualId = 0L
if (-not [long]::TryParse([string] $release.id, [ref] $actualId) -or $actualId -ne $ReleaseId) {
    throw "Release ID mismatch for ${Tag}: expected $ReleaseId, got '$($release.id)'."
}
if ($release.tag_name -cne $Tag) {
    throw "Release ID $ReleaseId has tag '$($release.tag_name)', expected '$Tag'."
}
if ($release.draft -or [string]::IsNullOrWhiteSpace([string] $release.published_at)) {
    throw "Release $Tag (ID $ReleaseId) is not published."
}

$pagesJson = & gh api "repos/$Repository/releases/$ReleaseId/assets?per_page=100" --paginate --slurp
if ($LASTEXITCODE -ne 0) {
    throw "Failed to list release assets for $Tag (ID $ReleaseId)."
}
$pages = ConvertFrom-Json -InputObject ($pagesJson -join "`n") -NoEnumerate
if ($pages -isnot [array]) {
    throw "Unexpected asset list response for $Tag (ID $ReleaseId)."
}

$selected = [Collections.Generic.Dictionary[string, object]]::new([StringComparer]::OrdinalIgnoreCase)
foreach ($page in $pages) {
    if ($page -isnot [array]) {
        throw "Unexpected asset page for $Tag (ID $ReleaseId)."
    }
    foreach ($asset in $page) {
        $name = [string] $asset.name
        if ($name -cne 'checksums.txt' -and $name -cnotmatch '^tflint-ruleset-avm_.*\.zip$') {
            continue
        }
        if ($name -cnotmatch '^[A-Za-z0-9._-]+$') {
            throw "Unsafe release asset name: $name"
        }
        if ($selected.ContainsKey($name)) {
            throw "Duplicate release asset: $name"
        }
        $assetId = 0L
        $size = 0L
        if (-not [long]::TryParse([string] $asset.id, [ref] $assetId) -or $assetId -le 0) {
            throw "Invalid ID for release asset $name."
        }
        if (-not [long]::TryParse([string] $asset.size, [ref] $size) -or $size -le 0) {
            throw "Invalid size for release asset $name."
        }
        if ($asset.state -cne 'uploaded') {
            throw "Release asset $name (ID $assetId) is not fully uploaded (state: $($asset.state))."
        }
        $selected.Add($name, [pscustomobject]@{ Name = $name; Id = $assetId; Size = $size })
    }
}
if (-not $selected.ContainsKey('checksums.txt')) {
    throw "Release $Tag (ID $ReleaseId) is missing checksums.txt."
}
if ($selected.Count -eq 1) {
    throw "Release $Tag (ID $ReleaseId) is missing tflint-ruleset-avm_*.zip archives."
}

if (Test-Path -LiteralPath $OutputDirectory) {
    if (-not (Test-Path -LiteralPath $OutputDirectory -PathType Container) -or
        @(Get-ChildItem -LiteralPath $OutputDirectory -Force).Count -ne 0) {
        throw "Release output directory must be empty: $OutputDirectory"
    }
}
else {
    $null = New-Item -ItemType Directory -Path $OutputDirectory
}
$outputRoot = (Resolve-Path -LiteralPath $OutputDirectory).Path
$headers = @{
    Authorization = "Bearer $env:GH_TOKEN"
    Accept = 'application/octet-stream'
}

foreach ($asset in $selected.Values | Sort-Object Name) {
    $path = Join-Path $outputRoot $asset.Name
    $uri = "https://api.github.com/repos/$Repository/releases/assets/$($asset.Id)"
    $null = Invoke-WebRequest -Uri $uri -Method Get -Headers $headers -OutFile $path
    $actualSize = (Get-Item -LiteralPath $path).Length
    if ($actualSize -ne $asset.Size) {
        throw "Release asset $($asset.Name) size mismatch: expected $($asset.Size) bytes, got $actualSize."
    }
}

Write-Host "Downloaded $($selected.Count) release assets for $Tag (ID $ReleaseId)."
