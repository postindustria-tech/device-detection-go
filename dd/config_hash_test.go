package dd

// #include <string.h>
// #include "./device-detection-cxx/src/hash/hash.h"
// #include "./device-detection-cxx/src/hash/fiftyone.h"
import (
	"testing"
)

func TestNewConfigHashPerformanceProfiles(t *testing.T) {
	var testInputs = []struct {
		perf     PerformanceProfile
		expected PerformanceProfile
	}{
		{Default, Default},
		{LowMemory, LowMemory},
		{BalancedTemp, BalancedTemp},
		{Balanced, Balanced},
		{HighPerformance, HighPerformance},
		{InMemory, InMemory},
		{100, Default},
	}
	for _, inputs := range testInputs {
		config := NewConfigHash(inputs.perf)
		actual := config.PerformanceProfile()
		if actual != inputs.expected {
			t.Errorf("Expected performance \"%d\" but get \"%d\"\n",
				inputs.expected, actual)
		}
	}
}

func TestNewConfigHashSetGetNumeric(t *testing.T) {
	expectedDiff := 4
	expectedDrift := 5
	expectedCon := 6

	config := NewConfigHash(Default)
	config.SetDifference(int32(expectedDiff))
	config.SetDrift(int32(expectedDrift))
	config.SetConcurrency(uint16(expectedCon))

	actualDiff := config.Difference()
	if actualDiff != int32(expectedDiff) {
		t.Errorf("Expected difference \"%d\" but get \"%d\"\n",
			expectedDiff, actualDiff)
	}

	actualDrift := config.Drift()
	if actualDrift != int32(expectedDrift) {
		t.Errorf("Expected difference \"%d\" but get \"%d\"\n",
			expectedDrift, actualDrift)
	}

	actualCon := config.Concurrency()
	if actualCon != uint16(expectedCon) {
		t.Errorf("Expected difference \"%d\" but get \"%d\"\n",
			expectedCon, actualCon)
	}
}

func TestNewConfigHashSetGetBoolean(t *testing.T) {
	// Test cases
	data := []bool{true, false}

	for _, testData := range data {
		config := NewConfigHash(Default)
		config.SetUsePerformanceGraph(testData)
		config.SetUsePredictiveGraph(testData)
		config.SetTraceRoute(testData)
		config.SetUseUpperPrefixHeaders(testData)
		config.SetUpdateMatchedUserAgent(testData)
		config.SetAllowUnmatched(testData)

		actualUsePerf := config.UsePerformanceGraph()
		if actualUsePerf != testData {
			t.Errorf("Expected UsePerformanceGraph to be %t but get \"%t\"\n",
				testData, actualUsePerf)
		}

		actualUsePred := config.UsePredictiveGraph()
		if actualUsePred != testData {
			t.Errorf("Expected UsePredictiveGraph to be %t but get \"%t\"\n",
				testData, actualUsePred)
		}

		actualTraceRoute := config.TraceRoute()
		if actualTraceRoute != testData {
			t.Errorf("Expected TraceRoute to be %t but get \"%t\"\n",
				testData, actualTraceRoute)
		}

		actualUseUpperPrefixHeaders := config.UseUpperPrefixHeaders()
		if actualUseUpperPrefixHeaders != testData {
			t.Errorf("Expected UseUpperPrefixHeaders to be %t but get \"%t\"\n",
				testData, actualUseUpperPrefixHeaders)
		}

		actualUpdateMatchedUserAgent := config.UpdateMatchedUserAgent()
		if actualUpdateMatchedUserAgent != testData {
			t.Errorf("Expected UpdateMatchedUserAgent to be %t but get \"%t\"\n",
				testData, actualUpdateMatchedUserAgent)
		}

		actualAllowUnmatched := config.AllowUnmatched()
		if actualAllowUnmatched != testData {
			t.Errorf("Expected AllowUnmatched to be %t but get \"%t\"\n",
				testData, actualAllowUnmatched)
		}
	}
}

func TestSizeManagerFromFile(t *testing.T) {
	testFunc := func(perf PerformanceProfile, t *testing.T) {
		config := NewConfigHash(perf)
		s, err := config.SizeManagerFromFile("", testDataFilePath)
		if err != nil {
			t.Errorf("Failed to size manager from file")
		}

		if s <= 0 {
			t.Errorf("Expected a positive size")
		}
	}
	performMultiPerfProfiles(testFunc, t)
}

func TestSizeManagerFromMemory(t *testing.T) {
	// TODO: Size the manager from memory
}
