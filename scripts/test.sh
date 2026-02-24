$testableDir=[
    "manager"
]

if [ -d "cover" ]; then
    rm -rf cover
fi
mkdir cover

$coverageDir=$(dirname "$0")/cover

for dir in "${testableDir[@]}"; do
    echo "----------- Running tests in $dir -----------"
    cd "$dir"

    if [ "$1" == "cover" ]; then
        $coverFile="$coverageDir/$(basename "$dir").coverprofile"
        $htmlFile="$coverageDir/$(basename "$dir").coverage.html"
        go test . -v -covermode=atomic -coverprofile="$coverFile"
        go tool cover -html="$coverFile" -o "$htmlFile"
    else
        go test . -v
    fi
done