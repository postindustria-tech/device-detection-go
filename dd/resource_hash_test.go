package dd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

// Test Init manager from file.
func TestInitManagerFromFile(t *testing.T) {
	// Setup a resource manager
	manager := NewResourceManager()
	// Free the resource manager
	defer manager.Free()
	// Create config per performance
	config := NewConfigHash(Balanced)
	config.SetAllowUnmatched(false)
	err := InitManagerFromFile(
		manager,
		*config,
		"",
		testDataFilePath,
	)
	if err != nil {
		t.Errorf("Failed to initialize resource manager from file \"%s\"\n.",
			testDataFilePath)
	}
}

func TestBadDataFile(t *testing.T) {
	badDataFilePath := "./badData.hash"
	// Create bad data file
	f, err := os.Create(badDataFilePath)
	if err != nil {
		log.Fatalf("ERROR: Failed to create bad data file \"%s\".", badDataFilePath)
	}
	// Make sure the file is removed eventually
	defer func() {
		err := os.Remove(badDataFilePath)
		if err != nil {
			t.Error(err)
		}
	}()
	// Then make sure the file is also closed
	defer f.Close()

	// Create a writer
	w := bufio.NewWriter(f)

	// Write bad data
	fmt.Fprint(w, "Bad data")
	w.Flush()

	// Define main test to be performed with different performance profile
	testFunc := func(perf PerformanceProfile, t *testing.T) {
		// Setup a resource manager
		manager := NewResourceManager()
		// Create config per performance
		config := NewConfigHash(perf)
		config.SetAllowUnmatched(false)
		err := InitManagerFromFile(
			manager,
			*config,
			"",
			badDataFilePath,
		)
		if err == nil {
			t.Errorf("Failed to detect bad data from file \"%s\"\n.",
				badDataFilePath)
		}

		if !strings.Contains(fmt.Sprintf("%s", err), "correct format") {
			t.Errorf("Incorrect error returned. Expected '..incorrect format..'\n")
		}
	}
	performMultiPerfProfiles(testFunc, t)
}
