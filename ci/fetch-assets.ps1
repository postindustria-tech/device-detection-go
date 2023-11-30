param (
    [Parameter(Mandatory=$true)]
    [string]$RepoName
)

$assets = New-Item -ItemType Directory -Path assets -Force
$deviceDetectionData = "$RepoName/dd/device-detection-cxx/device-detection-data"
$file = "51Degrees-LiteV4.1.hash"

if (!(Test-Path $assets/$file)) {
    Write-Output "Downloading $file"
    Invoke-WebRequest -Uri "https://github.com/51Degrees/device-detection-data/raw/main/51Degrees-LiteV4.1.hash" -OutFile $assets/$file
} else {
    Write-Output "'$file' exists, skipping download"
}

New-Item -ItemType SymbolicLink -Force -Target "$assets/51Degrees-LiteV4.1.hash" -Path "$deviceDetectionData/51Degrees-LiteV4.1.hash"
