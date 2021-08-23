package ddonpremise

// #cgo LDFLAGS: -L./device-detection-cxx/build/lib -lfiftyone-hash-c -lfiftyone-device-detection-c -lfiftyone-common-c -lm -latomic
// #include <string.h>
// #include "device-detection-cxx/src/hash/hash.h"
// #include "device-detection-cxx/src/hash/fiftyone.h"
import "C"

// ResultsHash wraps around a pointer to a value of C ResultsHash structure
type ResultsHash struct {
	CPtr *C.ResultsHash
}

/* Metric Getters */

// GetDeviceId returns the unique device id.
func (results *ResultsHash) GetDeviceId() string {
	return ""
}

// GetDeviceIdByIndex returns the unique device id of a result pointed by a index.
func (results *ResultsHash) GetDeviceIdByIndex(index uint32) string {
	return ""
}

// GetIterations returns the number of iterations carried out in order to find
// a match.
func (results *ResultsHash) GetIterations() int {
	return 0
}

// GetDrift returns the maximum drift for a matched substring.
func (results *ResultsHash) GetDrift() int {
	return 0
}

// GetDriftByIndex returns the drift for a matched substring of a result pointed
// by an index.
func (results *ResultsHash) GetDritfByIndex(index uint32) int {
	return 0
}

// GetDifference returns the total difference between the results returned and
// the target User-Agent.
func (results *ResultsHash) GetDifference() int {
	return 0
}

// GetDifferenceByIndex returns the difference between the result pointed by a
// index and the target User-Agent.
func (results *ResultsHash) GetDifferenceByIndex(index uint32) int {
	return 0
}

// GetMatchedNodes returns the number of hash nodes matched within the evidence.
func (results *ResultsHash) GetMatchedNodes() int {
	return 0
}

// GetMethod returns the method used to determine the match result.
func (results *ResultsHash) GetMethod(index uint32) int {
	return 0
}

// GetMethodByIndex returns the method use to determine a match result pointed
// by a index.
func (results *ResultsHash) GetMethodByIndex(index uint32) int {
	return 0
}

// GetTrace returns the trace route in a readable format.
func (results *ResultsHash) GetTrace() string {
	return ""
}

// GetTraceByIndex returns the trace reoute in a readable format of a result
// pointed by a given index.
func (results *ResultsHash) GetTraceByIndex(index uint32) string {
	return ""
}

// GetUserAgents returns number of User-Agents that were used in the results.
func (results *ResultsHash) GetUserAgents() int {
	return 0
}

// GetUserAgent return the user agent of a result pointed by a given index.
func (results *ResultsHash) GetUserAgent(index int) string {
	return ""
}

/* Constructor and Destructor */
// Free free the resource allocated in the C layer.
func (results *ResultsHash) Free() error {
	_, err := C.ResultsHashFree(results.CPtr)
	return err
}

/* Detections */
// MatchUserAgent performs a detection on an User-Agent string. This matches
// the C API fiftyoneDegreesResultsHashFromUserAgent.
func (results *ResultsHash) MatchUserAgent(ua string) error {
	return nil
}

// MatchEvidence performs a detection on evidence encapsulated in a EvidenceHash
// object. This matches the C API fiftyoneDegreesResultsHashFromEvidence.
func (results *ResultsHash) MatchEvidence(evidence EvidenceHash) error {
	return nil
}

// MatchDeviceId performs a detection on a given device id. This matches
// the C API fiftyoneDegreesResultsHashFromDeviceId
func (results *ResultsHash) MatchDeviceId(id string) error {
	return nil
}

// get values functions

// get has values

// get no value message
