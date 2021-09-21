package dd

import (
	"log"
	"strings"
	"testing"
)

// Test function to be executed by a test decorator method
type TestFunc func(manager *ResourceManager, t *testing.T)

// This is a decorator method which take a test function and execute it with
// different performance profile
func testDifferentPerformanceProfiles(testFunc TestFunc, t *testing.T) {
	// TODO: This should be replaced by a function that automatically search for
	// correct file name.
	// Get the data file Path
	dataFilePath := "./device-detection-cxx/device-detection-data/51Degrees-LiteV4.1.hash"
	data := []PerformanceProfile{
		Default,
		LowMemory,
		Balanced,
		// BalancedTemp, // TODO: Enable once error is resolved
		HighPerformance,
		InMemory,
	}

	for _, perf := range data {
		// Setup a resource manager
		manager := NewResourceManager()
		// Create config per performance
		config := NewConfigHash(perf)
		config.SetAllowUnmatched(false)
		err := InitManagerFromFile(
			manager,
			*config,
			"",
			dataFilePath,
		)
		if err != nil {
			log.Fatalf("Failed to initialize resource manager from file \"%s\"\n.",
				dataFilePath)
		}
		// Run the test
		testFunc(manager, t)

		// Free the resource manager
		manager.Free()
		if err != nil {
			log.Fatalln("Failed to free resource manager.")
		}
	}
}

// This tests ResultsHash is created and freed
// correctedly
func TestResultsHashCreateFree(t *testing.T) {
	var testFunc TestFunc = func(manager *ResourceManager, t *testing.T) {
		// Create results hash
		results, err := NewResultsHash(manager, 1, 0)
		if err != nil {
			t.Error("Failed to create results hash object.")
		}

		// Free results hash
		err = results.Free()
		if err != nil {
			t.Error("Failed to free results hash object.")
		}
	}
	// Execute the test with different performance profiles
	testDifferentPerformanceProfiles(testFunc, t)
}

// This tests a detection is performed successfully on a
// given User-Agent
func TestMatchUserAgent(t *testing.T) {
	var testFunc TestFunc = func(manager *ResourceManager, t *testing.T) {
		data := []struct {
			ua       string
			expected string
		}{
			{uaMobile, "true"},
			{uaDesktop, "false"},
		}

		// Create results hash
		// Create results hash
		results, err := NewResultsHash(manager, 1, 0)
		if err != nil {
			t.Error("Failed to create results hash object.")
		}

		// Free results hash at the end
		defer func() {
			err = results.Free()
			if err != nil {
				t.Error("Failed to free results hash object.")
			}
		}()

		for _, testData := range data {
			if err := results.MatchUserAgent(testData.ua); err != nil {
				t.Errorf("Failed to perform detection on User-Agent \"%s\".\n",
					testData.ua)
			}

			// Get values
			value, _, err := results.ValuesString("IsMobile", 50, ",")
			if err != nil {
				t.Errorf("Failed to get values string for User-Agent \"%s\".\n",
					testData.ua)
			}
			// Check if detection is correct
			if !strings.EqualFold(value, testData.expected) {
				t.Errorf("Expected \"%s\" but get \"%s\"\n.",
					testData.expected,
					value,
				)
			}
		}
	}
	// Execute the test with different performance profiles
	testDifferentPerformanceProfiles(testFunc, t)
}

// This tests match metrics are returned in correct
// successfully and in correct format.
func TestMatchMetrics(t *testing.T) {
	var testFunc TestFunc = func(manager *ResourceManager, t *testing.T) {
		// Create results hash
		results, err := NewResultsHash(manager, 1, 0)
		if err != nil {
			t.Error("Failed to create results hash object.")
		}

		// Free results hash at the end
		defer func() {
			err = results.Free()
			if err != nil {
				t.Error("Failed to free results hash object.")
			}
		}()

		if err := results.MatchUserAgent(uaMobile); err != nil {
			t.Errorf("Failed to perform detection on User-Agent \"%s\".\n",
				uaMobile)
		}

		// Check iterations
		if results.Iterations() <= 0 {
			t.Error("Expected iterations to always be biggers than 0.")
		}

		// Check drift
		if results.Drift() < 0 {
			t.Error("Expected drift to always be bigger be than or equal to 0.")
		}

		// Check difference
		if results.Difference() < 0 {
			t.Error("Expected difference to always be bigger than or equal to 0.")
		}

		// Check matched nodes
		if results.MatchedNodes() <= 0 {
			t.Error("Expected matched nodes to always bigger than 0.")
		}

		// Check method
		if results.Method() == None {
			t.Error("Expected method to not be None.")
		}

		// Check User-Agents
		if results.UserAgents() != 1 {
			t.Error("Expected number of used User-Agents to be 1 as there was " +
				"only one input.")
		}

		// Check match User-Agent string
		if strings.TrimSpace(results.UserAgent(0)) == "" {
			t.Error("Expected matched User-Agent not to be an empty or spacey " +
				"string.")
		}
	}
	// Execute the test with different performance profiles
	testDifferentPerformanceProfiles(testFunc, t)
}

// This tests HasValues reports correctly when a
// value is expected to be returned
func TestHasValues(t *testing.T) {
	var testFunc TestFunc = func(manager *ResourceManager, t *testing.T) {
		// Create results hash
		results, err := NewResultsHash(manager, 1, 0)
		if err != nil {
			t.Error("Failed to create results hash object.")
		}

		// Free results hash at the end
		defer func() {
			err = results.Free()
			if err != nil {
				t.Error("Failed to free results hash object.")
			}
		}()

		if err := results.MatchUserAgent(uaMobile); err != nil {
			t.Errorf("Failed to perform detection on User-Agent \"%s\".\n",
				uaMobile)
		}

		// Check results
		propertyName := "IsMobile"
		hasValues, err := results.HasValues(propertyName)
		if err != nil {
			t.Error("Failed to perform 'HasValues' on results.")
		}

		// Test if has values
		if !hasValues {
			t.Errorf("Expected propery \"%s\" to have values.\n", propertyName)
		}
	}
	// Execute the test with different performance profiles
	testDifferentPerformanceProfiles(testFunc, t)
}

// This tests HasValues reports correctly when no
// values is expected. This also test the correct
// message being returned.
func TestHasNoValues(t *testing.T) {
	var testFunc TestFunc = func(manager *ResourceManager, t *testing.T) {
		// Create results hash
		results, err := NewResultsHash(manager, 1, 0)
		if err != nil {
			t.Error("Failed to create results hash object.")
		}

		// Free results hash at the end
		defer func() {
			err = results.Free()
			if err != nil {
				t.Error("Failed to free results hash object.")
			}
		}()

		if err := results.MatchUserAgent(""); err != nil {
			t.Errorf("Failed to perform detection on User-Agent \"%s\".\n",
				uaMobile)
		}

		// Check results
		propertyName := "IsMobile"
		hasValues, err := results.HasValues(propertyName)
		if err != nil {
			t.Error("Failed to perform 'HasValues' on results.")
		}

		// Test if has values
		if hasValues {
			t.Errorf("Expected propery \"%s\" to have no values.\n", propertyName)
		}

		// Test values string
		value, _, err := results.ValuesString(propertyName, 50, ",")
		if err != nil {
			t.Error("Failed to get values string.")
		}

		// Expected no value
		if value != "" {
			t.Error("Expected value to an empty string")
		}

		// Test no values message
		message, err := results.NoValueReasonMessage(propertyName)
		if err != nil {
			t.Error("Failed to get no value message.")
		}

		// Expected certain message
		if !strings.Contains(message, "No matching profiles could be found") {
			t.Error("In correct values message.")
		}
	}
	// Execute the test with different performance profiles
	testDifferentPerformanceProfiles(testFunc, t)
}

// This tests a list of values are returned correctly
// from results of a detection
func TestValues(t *testing.T) {
	// TODO: To be implemented
}

// This tests ValuesString to be returned correctly
// for a property including input.
func TestValuesString(t *testing.T) {
	var testFunc TestFunc = func(manager *ResourceManager, t *testing.T) {
		// Expected output
		data := struct {
			property string
			expected string
		}{"IsMobile", "true"}

		// Create results hash
		results, err := NewResultsHash(manager, 1, 0)
		if err != nil {
			t.Error("Failed to create results hash object.")
		}

		// Free results hash at the end
		defer func() {
			err = results.Free()
			if err != nil {
				t.Error("Failed to free results hash object.")
			}
		}()

		if err := results.MatchUserAgent(uaMobile); err != nil {
			t.Errorf("Failed to perform detection on User-Agent \"%s\".\n",
				uaMobile)
		}

		// Get values
		value, _, err := results.ValuesString(data.property, 50, ",")
		if err != nil {
			t.Errorf("Failed to get values string for User-Agent \"%s\".\n",
				uaMobile)
		}
		// Check if detection is correct
		if !strings.EqualFold(value, data.expected) {
			t.Errorf("Expected \"%s\" but get \"%s\"\n.",
				data.expected,
				value,
			)
		}
	}
	// Execute the test with different performance profiles
	testDifferentPerformanceProfiles(testFunc, t)
}

// This tests ValuesString to be returned correctly
// for a property including input.
func TestValuesStringInsufficientSize(t *testing.T) {
	var testFunc TestFunc = func(manager *ResourceManager, t *testing.T) {
		// Expected output
		data := struct {
			property string
			expected uint64
		}{"IsMobile", uint64(len("true"))}

		// Create results hash
		results, err := NewResultsHash(manager, 1, 0)
		if err != nil {
			t.Error("Failed to create results hash object.")
		}

		// Free results hash at the end
		defer func() {
			err = results.Free()
			if err != nil {
				t.Error("Failed to free results hash object.")
			}
		}()

		if err := results.MatchUserAgent(uaMobile); err != nil {
			t.Errorf("Failed to perform detection on User-Agent \"%s\".\n",
				uaMobile)
		}

		// Get values
		value, actualSize, err := results.ValuesString(data.property, 1, ",")
		if err != nil {
			t.Errorf("Failed to get values string for User-Agent \"%s\".\n",
				uaMobile)
		}

		// Make sure actualSize is correct
		if actualSize != data.expected {
			t.Errorf("Expected actual size to be %d but get %d.\n",
				data.expected,
				actualSize)
		}

		// Make sure value is emtpy
		if value != "" {
			t.Error("Expected value to be an empty string since there wasn't " +
				"enough space allocated.")
		}
	}
	// Execute the test with different performance profiles
	testDifferentPerformanceProfiles(testFunc, t)
}
