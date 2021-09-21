package dd

import (
	"log"
	"os"
	"strings"
	"testing"
)

// TODO: Move the below to a sub package or common test file.
// test User-Agents
const uaMobile = "Mozilla/5.0 (iPhone; CPU iPhone OS 7_1 like Mac OS X) " +
	"AppleWebKit/537.51.2 (KHTML, like Gecko) Version/7.0 Mobile/11D167 " +
	"Safari/9537.53"

// User-Agent string of Firefox Web browser version 41 on desktop.
const uaDesktop = "Mozilla/5.0 (Windows NT 6.3; WOW64; rv:41.0) " +
	"Gecko/20100101 Firefox/41.0"

// User-Agent string of a MediaHub device.
const uaMediaHub = "Mozilla/5.0 (Linux; Android 4.4.2; X7 Quad Core " +
	"Build/KOT49H) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 " +
	"Chrome/30.0.0.0 Safari/537.36"

// Global manager to be used by each test
var manager *ResourceManager = nil

func TestMain(m *testing.M) {
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

	code := 0
	for _, perf := range data {
		// Setup a resource manager
		manager = NewResourceManager()
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
		code = m.Run()

		// Free the resource manager
		manager.Free()
		if err != nil {
			log.Fatalln("Failed to free resource manager.")
		}

		// Exit if test failed
		if code != 0 {
			os.Exit(code)
		}
	}
	os.Exit(code)
}

// This tests ResultsHash is created and freed
// correctedly
func TestResultsHashCreateFree(t *testing.T) {
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

// This tests a detection is performed successfully on a
// given User-Agent
func TestMatchUserAgent(t *testing.T) {
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

// This tests match metrics are returned in correct
// successfully and in correct format.
func TestMatchMetrics(t *testing.T) {
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

// This tests HasValues reports correctly when a
// value is expected to be returned
func TestHasValues(t *testing.T) {
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

// This tests HasValues reports correctly when no
// values is expected. This also test the correct
// message being returned.
func TestHasNoValues(t *testing.T) {
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

// This tests a list of values are returned correctly
// from results of a detection
func TestValues(t *testing.T) {
	// TODO: To be implemented
}

// This tests ValuesString to be returned correctly
// for a property including input.
func TestValuesString(t *testing.T) {
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

// This tests ValuesString to be returned correctly
// for a property including input.
func TestValuesStringInsufficientSize(t *testing.T) {
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
