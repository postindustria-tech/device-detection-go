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
