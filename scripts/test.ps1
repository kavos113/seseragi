Param(
    [string]$cover
)

$testableDir = @(
    "manager",
    "model"
)

if (Test-Path "cover") {
    Remove-Item "cover" -Recurse -Force
}
New-Item -ItemType Directory -Path . -Name "cover" | Out-Null
$coverageDir = Join-Path (Get-Location) "cover"

$getCover = $cover -ne $null -and $cover -eq "cover"

foreach ($dir in $testableDir) {
    Write-Host "-------- Running tests in $dir ------------" -ForegroundColor Cyan
    Push-Location $dir
    try {
        if ($getCover) {
            $coverFile = Join-Path $coverageDir "$dir.coverprofile"
            $htmlFile = Join-Path $coverageDir "$dir.coverage.html"

            go test ./... -v -covermode=atomic -coverprofile="$coverFile"
            go tool cover -html="$coverFile" -o "$htmlFile"
        }
        else {
            go test ./... -v
        }
    }
    finally {
        Pop-Location
    }
}