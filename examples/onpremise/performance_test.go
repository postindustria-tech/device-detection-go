package ddonpremise_test

import (
	"fmt"
	"os"
	"log"
	"bufio"
	"sync"
//	"runtime"
	"time"
//	"strings"
	"github.com/51Degrees/device-detection-go/ddonpremise"
)

type report struct {
	uaCount uint32
	uaIsMobile uint32
	uaProcessed uint32
}

func matchUserAgent(
	wg *sync.WaitGroup,
	manager ddonpremise.ResourceManager,
	ua string,
	calibration bool,
	rep report) {
	if !calibration {
		// Create results
		results, err := ddonpremise.FiftyoneDegreesResultsHashCreate(
			manager,
			1,
			0)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println(ua)
		// Perform detection
		err = ddonpremise.FiftyoneDegreesResultsHashFromUserAgent(
			results,
			ua)
		if err != nil {
			log.Fatal(err)
		}

		// Get the value in string
		_, _, err = ddonpremise.FiftyoneDegreesResultsHashGetValuesString(
			results,
			"IsMobile",
			100,
			",")
		if err != nil {
			log.Fatal(err)
		}
		
		// Update report
		// if (strings.Compare("True", value) == 0) {
		// 	rep.uaIsMobile++
		// }
		// rep.uaCount++

		// Free results
		err = ddonpremise.FiftyoneDegreesResultsHashFree(results)
		if err != nil {
			log.Fatal(err)
		}
	}
	// Complete and mark as done
	defer wg.Done()
}

func runTests(
	manager ddonpremise.ResourceManager,
	uaFilePath string,
	calibration bool,
	rep report) {
	// for i := 0; i < 4; i++ {
		// Loop through the User-Agent file
		file, err := os.Open(uaFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		// Create a wait group
		var wg sync.WaitGroup
		for scanner.Scan() {
			// Increase wait group
			wg.Add(1)
			go matchUserAgent(
				&wg,
				manager,
				scanner.Text(),
				calibration, rep)
		}
		// Wait until all goroutines finish
		wg.Wait()

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	// }
}

func run(
	manager ddonpremise.ResourceManager,
	uaFilePath string) {
	// Calibration
	fmt.Println("Calibrating...")
	caliReport := report{0, 0, 0}
	start := time.Now()
	runTests(manager, uaFilePath, true, caliReport)
	end := time.Now()
	caliTime := end.Sub(start)
	fmt.Printf("Calibration time: %d\n", caliTime.Milliseconds())

	// Action
	fmt.Println("Running performance tests...")
	actReport := report{0, 0, 0}
	start = time.Now()
	runTests(manager, uaFilePath, false, actReport)
	end = time.Now()
	actTime := end.Sub(start)
	fmt.Printf("Action time: %d\n", actTime.Milliseconds())

	// Calculate actual performance
	fmt.Printf("Performance: %d\n", actTime.Milliseconds() - caliTime.Milliseconds())

	// Display results
	
}

func runPerformanceExample(
	dataFilePath string,
	uaFilePath string,
	perf ddonpremise.PerformanceProfile) {
	// TODO: Set concurrency to match GOMAXPROCX

	// Properties
	props := ddonpremise.NewPropertiesRequired("IsMobile")

	// TODO: Set the device detection specific parameters to avoid checking for
	// upper case prefixed headers and tracking the matched User-Agent
	// characters.
	
	// Create Resource Manager
	manager := ddonpremise.NewResourceManager()
	// TODO: Create ConfigHash based on specified
	config := ddonpremise.NewConfigHash()
	err := ddonpremise.FiftyoneDegreesHashInitManagerFromFile(
		manager,
		config,
		props,
		dataFilePath)
	if err != nil {
		panic(err)
	}

	// Run the performance tests
	run(manager, uaFilePath)

	// Free the Resource Manager
	err = ddonpremise.FiftyoneDegreesResourceManagerFree(manager)
	if err != nil {
		panic(err)
	}
}

func Example_Performance() {
	// TODO: Determine number of CPUs
	// runtime.NumCPU
	
	// Data file path
	dataFilePath := "../../ddonpremise/device-detection-cxx/device-detection-data/51Degrees-LiteV4.1.hash"
	// User-Agents file path
	uaFilePath := "../../ddonpremise/device-detection-cxx/device-detection-data/20000 User Agents.csv"

	runPerformanceExample(dataFilePath, uaFilePath, ddonpremise.PerfDefault)
	// Output:
	// .*
}
