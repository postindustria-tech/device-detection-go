# 51Degrees Device Detection Engines

![51Degrees](https://51degrees.com/DesktopModules/FiftyOne/Distributor/Logo.ashx?utm_source=github&utm_medium=repository&utm_content=readme_main&utm_campaign=go-open-source "Data rewards the curious") **Go Device Detection**

## Introduction

This repository contains Go Lite implementation of device detection engine. This provides a light set of device detection features that allow Go users to try and evaluate the capability of 51Degrees device detection solution in Go.

## Supported platforms and architectures

Go Lite implementation is currently supporting the following platforms and architectures:
- Linux 32/64 bit, Intel and ARM processor
- MacOS 64 bit, Intel and ARM processor
- Windows 64bit, Intel

Go version:
- 1.20

## Pre-requisites

### Data File

In order to perform device detection, you will need to use a 51Degrees data file.
A 'lite' file can be found at [device-detection-data](https://github.com/51degrees/device-detection-data). The `ci/fetch-assets.ps1` can be used to download it.
This 'lite' file has a significantly reduced set of properties. To obtain a
file with a more complete set of device properties see the
[51Degrees website](https://51degrees.com/pricing).

### Software

In order to build use device-detection-go the following are required:
- Powershell Core (7 or above)
- A C compiler that support C11 or above (Gcc on Linux, Clang on MacOS and MinGW-x64 on Windows)
- libatomic - which usually come with default Gcc, Clang installation

## Module and Packages

This module name `device-detection-go` at path `github.com/51Degrees/device-detection-go/v4`

This Go Lite version contains only one single package:
- `dd`

## Build and Usage

### Windows
If you are on Windows, make sure the path to the `MinGW-x64` `bin` folder is included in the `PATH`.

### Build steps for all platforms

Import the package as usual and it will get build automatically:

```
import "github.com/51Degrees/device-detection-go/v4/dd"
```

This is achieved by using an amalgamation of the [device-detection-cxx](https://github.com/51degrees/device-detection-cxx) repository, which is built automatically by CGo during `go build`.

The amalgamation is updated regularly by the [Nightly Package Update](https://github.com/51Degrees/device-detection-go/actions/workflows/nightly-package-update.yml) CI workflow.

## Test

Unit tests can be run with `go test` as below if in `device-detection-go` folder:
```
go test ./...
```

## APIs

To view APIs and their descriptions, users can use `go doc` in the package directory.
- First navigate to folder `dd`.
- Then run the below to display all APIs, structures and their descriptions.
```
go doc -all
```

**NOTE**: Not all APIs are available. Ones that are not available include string `TODO: To be implemented` in their descriptions.

## Examples

Examples are included in a separate repository at [device-detection-examples-go](https://github.com/51degrees/device-detection-examples-go). Users can follow its README.md for details and comments in each examples for how to run them.
