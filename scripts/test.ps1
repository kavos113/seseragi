Param(
    [string]$cover
)

if (Test-Path "cover") {
    Remove-Item "cover" -Recurse -Force
}
New-Item -ItemType Directory -Path . -Name "cover" | Out-Null
$coverageDir = Join-Path (Get-Location) "cover"

$getCover = $cover -ne $null -and $cover -eq "cover"

if ($getCover) {
    $coverFile = Join-Path $coverageDir "coverage.out"
    $htmlFile = Join-Path $coverageDir "coverage.html"

    go test ./... -v -covermode=atomic -coverprofile="$coverFile"
    go tool cover -html="$coverFile" -o "$htmlFile"
}
else {
    go test ./... -v
}