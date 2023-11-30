param (
    [Parameter(Mandatory=$true)]
    [string]$RepoName
)

./go/run-unit-tests.ps1 -RepoName $RepoName
