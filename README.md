# 51Degrees Device Detection Golang API
## Pre-requisites
To build and use 51Degrees device detection for Golang, the following are required:
- Powershell
- CMake 4.10 or above
- Gcc
- Go

## Build and Run
To run device-detection-go onpremise, the core engine binaries are required to
be built first. To do this, make sure the following utilities are available are
installed on your chosen platforms:
- Powershell
- CMake 3.10 or above
- Gcc (Use MinGW-W64 if on Windows)

Then run `prebuild.ps1` as located in `ddonpremise/scripts` folder.

Once the binaries have been successfully built. Examples and test can now be run.

## Test
Tests are located within each package folder.

To run tests, make sure the core binaries are built successfully. Then navigate to folder `ddonpremise` and run:
```
go test
```

## Examples
Examples are located within the `examples` folder.

The examples are testable examples so can be run as go tests. Make sure the core binaries are built successfully.

To run tests on the examples, navigate to folder `examples/onprmeise` and run:
```
go test
```

## Implementation guidelines
Each native resource wrapper that requires to be freed should implement a strategy that force the program to free the resource as soon as the resource is no longer needed.

## Possible strategies
### Use SetFinalizer
**Cons**: This poses a risk of resource not being freed quick enough and cause memory to quickly run out.

### Using defer
**To be considered**: This strictly runs at the end of the function execution so might provide enough assurance on the speed of resource deallocation.

### Using function call back
This can be implemented as:

```
func try(obj *interface{}, func interface{}) {
    func(obj)
    obj.Free()
}
```

This can be used as:
```
try(NewResultsHash(), func(results *ResultsHash) {
    results.MatchUserAgent("User-AGents")
})
```

**Cons**: This does not allow the function to access variables defined in the outside scope.