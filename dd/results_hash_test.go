/* *********************************************************************
 * This Original Work is copyright of 51 Degrees Mobile Experts Limited.
 * Copyright 2019 51 Degrees Mobile Experts Limited, 5 Charlotte Close,
 * Caversham, Reading, Berkshire, United Kingdom RG4 7BY.
 *
 * This Original Work is licensed under the European Union Public Licence (EUPL)
 * v.1.2 and is subject to its terms as set out below.
 *
 * If a copy of the EUPL was not distributed with this file, You can obtain
 * one at https://opensource.org/licenses/EUPL-1.2.
 *
 * The 'Compatible Licences' set out in the Appendix to the EUPL (as may be
 * amended by the European Commission) shall be deemed incompatible for
 * the purposes of the Work and the provisions of the compatibility
 * clause in Article 5 of the EUPL shall not apply.
 *
 * If using the Work as, or as part of, a network application, by
 * including the attribution notice(s) required under Article 5 of the EUPL
 * in the end user terms of the application under an appropriate heading,
 * such notice(s) shall fulfill the requirements of that article.
 * ********************************************************************* */

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
	f := func(perf PerformanceProfile, t *testing.T) {
		// Setup a resource manager
		manager := NewResourceManager()
		// Create config per performance
		config := NewConfigHash(perf)
		config.SetAllowUnmatched(false)
		err := InitManagerFromFile(
			manager,
			*config,
			"",
			testDataFilePath,
		)
		if err != nil {
			log.Fatalf("Failed to initialize resource manager from file \"%s\"\n.",
				testDataFilePath)
		}
		// Run the test
		testFunc(manager, t)

		// Free the resource manager
		manager.Free()
	}
	performMultiPerfProfiles(f, t)
}

// This tests ResultsHash is created and freed
// correctedly
func TestResultsHashCreateFree(t *testing.T) {
	var testFunc TestFunc = func(manager *ResourceManager, t *testing.T) {
		// Create results hash
		results := NewResultsHash(manager, 1, 0)

		// Free results hash
		results.Free()
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
		results := NewResultsHash(manager, 1, 0)

		// Free results hash at the end
		defer results.Free()

		for _, testData := range data {
			if err := results.MatchUserAgent(testData.ua); err != nil {
				t.Errorf("Failed to perform detection on User-Agent \"%s\".\n",
					testData.ua)
			}

			// Get values
			value, err := results.ValuesString("IsMobile", ",")
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
		results := NewResultsHash(manager, 1, 0)

		// Free results hash at the end
		defer results.Free()

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
		results := NewResultsHash(manager, 1, 0)

		// Free results hash at the end
		defer results.Free()

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
		results := NewResultsHash(manager, 1, 0)

		// Free results hash at the end
		defer results.Free()

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
		value, err := results.ValuesString(propertyName, ",")
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
// for a property.
func TestValuesString(t *testing.T) {
	var testFunc TestFunc = func(manager *ResourceManager, t *testing.T) {
		// Expected output
		data := struct {
			property string
			expected string
		}{"IsMobile", "true"}

		// Create results hash
		results := NewResultsHash(manager, 1, 0)

		// Free results hash at the end
		defer results.Free()

		if err := results.MatchUserAgent(uaMobile); err != nil {
			t.Errorf("Failed to perform detection on User-Agent \"%s\".\n",
				uaMobile)
		}

		// Get values
		value, err := results.ValuesString(data.property, ",")
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

// This tests valuesStringSwithSize to be returned correctly
// for a property if default size is smaller than actual.
func TestValuesStringWithSize(t *testing.T) {
	var testFunc TestFunc = func(manager *ResourceManager, t *testing.T) {
		// Expected output
		data := struct {
			property string
			expected string
		}{"IsMobile", "true"}
		// Use small size to make sure the function will reattempt with
		// a return actual size
		const testSize = 1
		// Create results hash
		results := NewResultsHash(manager, 1, 0)

		// Free results hash at the end
		defer results.Free()

		if err := results.MatchUserAgent(uaMobile); err != nil {
			t.Errorf("Failed to perform detection on User-Agent \"%s\".\n",
				uaMobile)
		}

		// Get values
		value, err := valuesStringWithSize(results, data.property, testSize, ",")
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

// Test that resource finalizer will panic if a pointer to
// C resource has not been freed
func TestResultsFinalizer(t *testing.T) {
	var testFunc TestFunc = func(manager *ResourceManager, t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Results finalizer did not panic.")
			}
		}()

		// Create results hash
		results := NewResultsHash(manager, 1, 0)

		// Perform finalizer on live results
		resultsFinalizer(results)

		// Free results hash
		results.Free()
	}
	// Execute the test with different performance profiles
	testDifferentPerformanceProfiles(testFunc, t)
}
