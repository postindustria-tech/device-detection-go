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

func TestNewConfigHashSetGetBooleanFalse(t *testing.T) {
	config := NewConfigHash(Default)
	config.SetUsePerformanceGraph(false)
	config.SetUsePredictiveGraph(false)
	config.SetTraceRoute(false)
	config.SetUseUpperPrefixHeaders(false)
	config.SetUpdateMatchedUserAgent(false)

	actualUsePerf := config.UsePerformanceGraph()
	if actualUsePerf {
		t.Errorf("Expected UsePerformanceGraph to be false but get \"%t\"\n",
			actualUsePerf)
	}

	actualUsePred := config.UsePredictiveGraph()
	if actualUsePred {
		t.Errorf("Expected UsePredictiveGraph to be false but get \"%t\"\n",
			actualUsePred)
	}

	actualTraceRoute := config.TraceRoute()
	if actualTraceRoute {
		t.Errorf("Expected TraceRoute to be false but get \"%t\"\n",
			actualTraceRoute)
	}

	actualUseUpperPrefixHeaders := config.UseUpperPrefixHeaders()
	if actualUseUpperPrefixHeaders {
		t.Errorf("Expected UseUpperPrefixHeaders to be false but get \"%t\"\n",
			actualUseUpperPrefixHeaders)
	}

	actualUpdateMatchedUserAgent := config.UpdateMatchedUserAgent()
	if actualUpdateMatchedUserAgent {
		t.Errorf("Expected UpdateMatchedUserAgent to be false but get \"%t\"\n",
			actualUpdateMatchedUserAgent)
	}
}

func TestNewConfigHashSetGetBooleanTrue(t *testing.T) {
	config := NewConfigHash(Default)
	config.SetUsePerformanceGraph(true)
	config.SetUsePredictiveGraph(true)
	config.SetTraceRoute(true)
	config.SetUseUpperPrefixHeaders(true)
	config.SetUpdateMatchedUserAgent(true)

	actualUsePerf := config.UsePerformanceGraph()
	if !actualUsePerf {
		t.Errorf("Expected UsePerformanceGraph to be true but get \"%t\"\n",
			actualUsePerf)
	}

	actualUsePred := config.UsePredictiveGraph()
	if !actualUsePred {
		t.Errorf("Expected UsePredictiveGraph to be true but get \"%t\"\n",
			actualUsePred)
	}

	actualTraceRoute := config.TraceRoute()
	if !actualTraceRoute {
		t.Errorf("Expected TraceRoute to be true but get \"%t\"\n",
			actualTraceRoute)
	}

	actualUseUpperPrefixHeaders := config.UseUpperPrefixHeaders()
	if !actualUseUpperPrefixHeaders {
		t.Errorf("Expected UseUpperPrefixHeaders to be true but get \"%t\"\n",
			actualUseUpperPrefixHeaders)
	}

	actualUpdateMatchedUserAgent := config.UpdateMatchedUserAgent()
	if !actualUpdateMatchedUserAgent {
		t.Errorf("Expected UpdateMatchedUserAgent to be true but get \"%t\"\n",
			actualUpdateMatchedUserAgent)
	}
}

func TestSizeManagerFromFile(t *testing.T) {
	// TODO: Size the manager from file
}

func TestSizeManagerFromMemory(t *testing.T) {
	// TODO: Size the manager from memory
}
