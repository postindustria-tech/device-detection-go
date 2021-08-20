package main

/*
To run this program, follow the following steps:
1. Build device-detection-cxx using CMake as described in its README.md
2. Update the linker option '-L' to point to the location of the built static libraries.
3. Naviaget to the parent folder, which mean 'cd ..'
4. Build the executable file using.
  go build -o [executable file name] ./main
5. Run the output file.
  ./[executable file name]
*/

import "fmt"
import "github.com/51Degrees/device-detection-go/ddonpremise"

func Example() {
	// Initialise manager
	manager := ddonpremise.NewResourceManager()
	props := ddonpremise.NewPropertiesRequired("")
	config := ddonpremise.NewConfigHash()
	filePath := "/home/fiftyonedegrees/staff/tung/device-detection-go/ddonpremise/device-detection-cxx/device-detection-data/51Degrees-LiteV4.1.hash"
	err := ddonpremise.FiftyoneDegreesHashInitManagerFromFile(
		manager,
		config,
		props,
		filePath)
	if err != nil {
		panic(err)
	}

	
	// Create results
	results, err := ddonpremise.FiftyoneDegreesResultsHashCreate(
		manager,
		1,
		0)
	if err != nil {
		panic(err)
	}

	// Example User-Agent
	ua := "Mozilla/5.0 (iPhone; CPU iPhone OS 7_1 like Mac OS X) AppleWebKit/537.51.2 (KHTML, like Gecko) Version/7.0 Mobile/11D167 Safari/9537.53"

	// Perform detection
	err = ddonpremise.FiftyoneDegreesResultsHashFromUserAgent(
		results,
		ua)
	if err != nil {
		panic(err)
	}

	// Get the values in string
	value, _, err := ddonpremise.FiftyoneDegreesResultsHashGetValuesString(
		results,
		"IsMobile",
		100,
		",")
	if err != nil {
		panic(err)
	}

	fmt.Println(value)

	// Free allocated resources
	err = ddonpremise.FiftyoneDegreesResultsHashFree(results)
	if err != nil {
		panic(err)
	}

	err = ddonpremise.FiftyoneDegreesResourceManagerFree(manager)
	if err != nil {
		panic(err)
	}

	// Output:
	// True
}
