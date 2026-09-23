#requires -Version 7.2

$ErrorActionPreference = 'Stop'
$root = Join-Path ([System.IO.Path]::GetTempPath()) "tflint-release-download-$([guid]::NewGuid())"
$binaryRoot = Join-Path $root 'binaries'
$fixtureRoot = Join-Path $root 'fixtures'
$previousToken = $env:GH_TOKEN
$previousState = Get-Variable -Name ReleaseDownloadTestState -Scope Global -ErrorAction SilentlyContinue
$global:ReleaseDownloadTestState = @{
    Repository = 'Azure/tflint-ruleset-avm'
    Tag = 'v0.18.0'
    ReleaseId = 12345L
    FixturePaths = @{}
    ApiCalls = [Collections.Generic.List[string]]::new()
    DownloadedIds = [Collections.Generic.List[long]]::new()
}

function Reset-Mock {
    $global:ReleaseDownloadTestState.Release = [pscustomobject]@{
        id = $global:ReleaseDownloadTestState.ReleaseId
        tag_name = $global:ReleaseDownloadTestState.Tag
        draft = $false
        published_at = '2026-09-23T00:00:00Z'
    }
    $global:ReleaseDownloadTestState.Assets = @(
        foreach ($asset in $global:ReleaseDownloadTestState.BaseAssets) {
            [pscustomobject]@{
                id = $asset.id
                name = $asset.name
                size = $asset.size
                state = $asset.state
            }
        }
    )
    $global:ReleaseDownloadTestState.ApiCalls.Clear()
    $global:ReleaseDownloadTestState.DownloadedIds.Clear()
    $global:ReleaseDownloadTestState.FailAssetList = $false
}

function gh {
    if ($args[0] -ne 'api') {
        throw "Unexpected gh command: $args"
    }
    $endpoints = @($args | Where-Object { $_ -like 'repos/*' })
    if ($endpoints.Count -ne 1) {
        throw "Unexpected gh api arguments: $args"
    }
    $endpoint = [string] $endpoints[0]
    $global:ReleaseDownloadTestState.ApiCalls.Add($endpoint)
    $global:LASTEXITCODE = 0

    if ($endpoint -eq "repos/$($global:ReleaseDownloadTestState.Repository)/releases/tags/$($global:ReleaseDownloadTestState.Tag)") {
        return (ConvertTo-Json -InputObject $global:ReleaseDownloadTestState.Release -Compress)
    }
    if ($endpoint -eq "repos/$($global:ReleaseDownloadTestState.Repository)/releases/$($global:ReleaseDownloadTestState.ReleaseId)") {
        return (ConvertTo-Json -InputObject $global:ReleaseDownloadTestState.Release -Compress)
    }
    if ($endpoint -eq "repos/$($global:ReleaseDownloadTestState.Repository)/releases/$($global:ReleaseDownloadTestState.ReleaseId)/assets?per_page=100") {
        if ($args -notcontains '--paginate' -or $args -notcontains '--slurp') {
            throw 'Release asset request must paginate and slurp all pages.'
        }
        if ($global:ReleaseDownloadTestState.FailAssetList) {
            $global:LASTEXITCODE = 1
            return ''
        }
        $pageOne = @($global:ReleaseDownloadTestState.Assets | Select-Object -First 3)
        $pageTwo = @($global:ReleaseDownloadTestState.Assets | Select-Object -Skip 3)
        return (ConvertTo-Json -InputObject @($pageOne, $pageTwo) -Compress -Depth 4)
    }
    throw "Unexpected GitHub API endpoint: $endpoint"
}

function Invoke-WebRequest {
    [CmdletBinding()]
    param(
        [Parameter(Mandatory)]
        [uri] $Uri,

        [Parameter(Mandatory)]
        [string] $OutFile,

        [Parameter(Mandatory)]
        [hashtable] $Headers,

        [string] $Method
    )

    if ($Uri.Scheme -ne 'https' -or $Uri.Host -ne 'api.github.com' -or $Method -ne 'Get' -or
        $Headers.Authorization -ne 'Bearer release-download-test-token' -or
        $Headers.Accept -ne 'application/octet-stream') {
        throw "Unsafe release asset download: $Uri"
    }
    if ($Uri.AbsolutePath -cnotmatch '^/repos/Azure/tflint-ruleset-avm/releases/assets/([1-9][0-9]*)$') {
        throw "Download did not use an asset ID: $Uri"
    }
    $assetId = [long] $Matches[1]
    $global:ReleaseDownloadTestState.DownloadedIds.Add($assetId)
    $path = $global:ReleaseDownloadTestState.FixturePaths[[string] $assetId]
    if (-not $path) {
        throw "Unknown fixture asset ID: $assetId"
    }
    Copy-Item -LiteralPath $path -Destination $OutFile
}

function Invoke-Download {
    param(
        [string] $Directory,
        [switch] $ResolveTag
    )

    $parameters = @{
        Tag = $global:ReleaseDownloadTestState.Tag
        Repository = $global:ReleaseDownloadTestState.Repository
        OutputDirectory = $Directory
    }
    if (-not $ResolveTag) {
        $parameters.ReleaseId = $global:ReleaseDownloadTestState.ReleaseId
    }
    & (Join-Path $PSScriptRoot 'Save-ReleaseAssets.ps1') @parameters
}

function Assert-Fails {
    param(
        [scriptblock] $Action,
        [string] $Expected
    )

    $message = $null
    try {
        & $Action
    }
    catch {
        $message = $_.Exception.Message
    }
    if (-not $message -or $message -notlike "*$Expected*") {
        throw "Expected failure containing '$Expected'; received '$message'."
    }
}

try {
    $env:GH_TOKEN = 'release-download-test-token'
    $targets = Get-Content -LiteralPath (Join-Path $PSScriptRoot '..' 'release' 'targets.json') -Raw |
        ConvertFrom-Json
    foreach ($target in $targets) {
        $directory = Join-Path $binaryRoot "$($target.goos)-$($target.goarch)"
        $null = New-Item -ItemType Directory -Path $directory -Force
        $name = if ($target.goos -eq 'windows') { 'tflint-ruleset-avm.exe' } else { 'tflint-ruleset-avm' }
        "fixture-$($target.goos)-$($target.goarch)" |
            Set-Content -LiteralPath (Join-Path $directory $name) -Encoding utf8NoBOM
    }
    $null = & (Join-Path $PSScriptRoot 'New-ReleaseAssets.ps1') `
        -BinaryDirectory $binaryRoot `
        -OutputDirectory $fixtureRoot

    $files = @(Get-ChildItem -LiteralPath $fixtureRoot -File | Sort-Object Name)
    $id = 1000L
    $global:ReleaseDownloadTestState.BaseAssets = @(
        foreach ($file in $files) {
            $id++
            $global:ReleaseDownloadTestState.FixturePaths[[string] $id] = $file.FullName
            [pscustomobject]@{ id = $id; name = $file.Name; size = $file.Length; state = 'uploaded' }
        }
    )

    Reset-Mock
    $byIdDirectory = Join-Path $root 'by-id'
    Invoke-Download -Directory $byIdDirectory
    & (Join-Path $PSScriptRoot 'Test-ReleaseAssets.ps1') -ReleaseDirectory $byIdDirectory
    if ($global:ReleaseDownloadTestState.DownloadedIds.Count -ne $global:ReleaseDownloadTestState.BaseAssets.Count -or
        @(Compare-Object ($global:ReleaseDownloadTestState.BaseAssets.id | Sort-Object) ($global:ReleaseDownloadTestState.DownloadedIds | Sort-Object)).Count -ne 0) {
        throw 'The download did not include every matching asset ID.'
    }
    if ($global:ReleaseDownloadTestState.ApiCalls -contains "repos/$($global:ReleaseDownloadTestState.Repository)/releases/tags/$($global:ReleaseDownloadTestState.Tag)") {
        throw 'A release event must use its release ID rather than resolving the tag.'
    }

    Reset-Mock
    Invoke-Download -Directory (Join-Path $root 'by-tag') -ResolveTag
    if ($global:ReleaseDownloadTestState.ApiCalls -notcontains "repos/$($global:ReleaseDownloadTestState.Repository)/releases/tags/$($global:ReleaseDownloadTestState.Tag)" -or
        $global:ReleaseDownloadTestState.ApiCalls -notcontains "repos/$($global:ReleaseDownloadTestState.Repository)/releases/$($global:ReleaseDownloadTestState.ReleaseId)") {
        throw 'A manual run must resolve the tag and validate its numeric release ID.'
    }

    Reset-Mock
    $global:ReleaseDownloadTestState.Assets = @($global:ReleaseDownloadTestState.Assets | Where-Object { $_.name -cne 'checksums.txt' })
    Assert-Fails -Expected 'missing checksums.txt' -Action {
        Invoke-Download -Directory (Join-Path $root 'missing-checksums')
    }
    if ($global:ReleaseDownloadTestState.DownloadedIds.Count -ne 0) {
        throw 'Missing checksums.txt must fail before any downloads.'
    }

    Reset-Mock
    $missingZip = ($global:ReleaseDownloadTestState.Assets | Where-Object { $_.name -like '*.zip' } | Select-Object -First 1).name
    $global:ReleaseDownloadTestState.Assets = @($global:ReleaseDownloadTestState.Assets | Where-Object { $_.name -cne $missingZip })
    $missingZipDirectory = Join-Path $root 'missing-zip'
    Invoke-Download -Directory $missingZipDirectory
    Assert-Fails -Expected 'Release files do not match' -Action {
        & (Join-Path $PSScriptRoot 'Test-ReleaseAssets.ps1') -ReleaseDirectory $missingZipDirectory
    }

    Reset-Mock
    ($global:ReleaseDownloadTestState.Assets | Where-Object { $_.name -ceq 'checksums.txt' }).state = 'open'
    Assert-Fails -Expected 'not fully uploaded' -Action {
        Invoke-Download -Directory (Join-Path $root 'incomplete-checksum')
    }

    Reset-Mock
    ($global:ReleaseDownloadTestState.Assets | Where-Object { $_.name -like '*.zip' } | Select-Object -First 1).state = 'open'
    Assert-Fails -Expected 'not fully uploaded' -Action {
        Invoke-Download -Directory (Join-Path $root 'incomplete-zip')
    }

    Reset-Mock
    ($global:ReleaseDownloadTestState.Assets | Where-Object { $_.name -ceq 'checksums.txt' }).size++
    Assert-Fails -Expected 'size mismatch' -Action {
        Invoke-Download -Directory (Join-Path $root 'wrong-size')
    }

    Reset-Mock
    $global:ReleaseDownloadTestState.Release.tag_name = 'v0.19.0'
    Assert-Fails -Expected "expected '$($global:ReleaseDownloadTestState.Tag)'" -Action {
        Invoke-Download -Directory (Join-Path $root 'wrong-tag')
    }

    Reset-Mock
    $global:ReleaseDownloadTestState.Release.id = 54321
    Assert-Fails -Expected 'Release ID mismatch' -Action {
        Invoke-Download -Directory (Join-Path $root 'wrong-id')
    }

    Reset-Mock
    $global:ReleaseDownloadTestState.Release.draft = $true
    Assert-Fails -Expected 'not published' -Action {
        Invoke-Download -Directory (Join-Path $root 'draft')
    }

    Reset-Mock
    $global:ReleaseDownloadTestState.Assets += [pscustomobject]@{
        id = 9999
        name = 'tflint-ruleset-avm_../escape.zip'
        size = 1
        state = 'uploaded'
    }
    Assert-Fails -Expected 'Unsafe release asset name' -Action {
        Invoke-Download -Directory (Join-Path $root 'unsafe-name')
    }

    Reset-Mock
    $global:ReleaseDownloadTestState.Assets += [pscustomobject]@{
        id = 9998
        name = 'checksums.txt'
        size = 1
        state = 'uploaded'
    }
    Assert-Fails -Expected 'Duplicate release asset' -Action {
        Invoke-Download -Directory (Join-Path $root 'duplicate-name')
    }

    Reset-Mock
    $global:ReleaseDownloadTestState.FailAssetList = $true
    Assert-Fails -Expected 'Failed to list release assets' -Action {
        Invoke-Download -Directory (Join-Path $root 'api-error')
    }

    Reset-Mock
    $existingDirectory = Join-Path $root 'existing'
    $null = New-Item -ItemType Directory -Path $existingDirectory
    'stale' | Set-Content -LiteralPath (Join-Path $existingDirectory 'checksums.txt')
    Assert-Fails -Expected 'Release output directory must be empty' -Action {
        Invoke-Download -Directory $existingDirectory
    }

    Reset-Mock
    $env:GH_TOKEN = $null
    Assert-Fails -Expected 'GH_TOKEN is required' -Action {
        Invoke-Download -Directory (Join-Path $root 'no-token')
    }

    Write-Host 'Release download tests passed.'
}
finally {
    $env:GH_TOKEN = $previousToken
    if ($null -eq $previousState) {
        Remove-Variable -Name ReleaseDownloadTestState -Scope Global -ErrorAction SilentlyContinue
    }
    else {
        $global:ReleaseDownloadTestState = $previousState.Value
    }
    if (Test-Path -LiteralPath $root) {
        Remove-Item -LiteralPath $root -Recurse -Force
    }
}
