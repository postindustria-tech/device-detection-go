param (
    [Parameter(Mandatory=$true)]
    [string]$RepoName
)

./go/build-project.ps1 -RepoName $RepoName
