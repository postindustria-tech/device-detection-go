package dd

// #include <string.h>
// #include "./device-detection-cxx/src/hash/hash.h"
// #include "./device-detection-cxx/src/hash/fiftyone.h"
import "C"
import (
	"fmt"
	"reflect"
)

// ResultsHash wraps around a pointer to a value of C ResultsHash structure
type ResultsHash struct {
	CPtr *C.ResultsHash
}

/* Constructor and Destructor */
// NewResultsHash create a new ResultsHash object. This matches the C API
// fiftyoneDegreesResultsHashCreate.
func NewResultsHash(
	manager *ResourceManager,
	uaCapacity uint32,
	overridesCapacity uint32) (results ResultsHash, err error) {
	r, err := C.ResultsHashCreate(
		manager.CPtr,
		C.uint(uaCapacity),
		C.uint(overridesCapacity))
	if err != nil {
		return ResultsHash{nil}, err
	}
	return ResultsHash{r}, nil
}

// Free free the resource allocated in the C layer.
func (results *ResultsHash) Free() error {
	_, err := C.ResultsHashFree(results.CPtr)
	return err
}

/* Metric Getters */

// GetDeviceId returns the unique device id. This matches the C API
// fiftyoneDegreesHashGetDeviceIdFromResults.
func (results *ResultsHash) GetDeviceId() (id string, err error) {
	return "", nil
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

// GetHasValues returns whether the last detection returns any matched value
// for a given property index. This matches the C API
// fiftyoneDegreesResultsHashGetHasValues
func (results *ResultsHash) GetHasValues(
	requiredPropertyIndex int) (r bool, err error) {
	return false, nil
}

// GetNoValueReasonMessage returns the no value message of a given property
// index. This matches a combination of the C APIs
// fiftyoneDegreesResultsNoValueReason and
// fiftyoneDegreesResultsHashGetNoValueReasonMessage
func (results *ResultsHash) GetNoValueReasonMessage(
	requiredPropertyIndex int) (message string, err error) {
	return "", nil
}

// GetValues returns a list of values resulted from a detection for a given
// property index.
func (results *ResultsHash) GetValues(
	requiredPropertyIndex int) (values []string, err error) {
	return nil, nil
}

// GetValuesString returns a string of all values resulted from a detection.
// If there are multiple values, they are separated by the specified separator.
func (results *ResultsHash) GetValuesString(
	propertyName string,
	size uint64,
	separator string) (value string, s uint64, err error) {
	// Create slice based on specified size
	buffer := make([]C.char, size)
	e := NewException()
	actualSize, err := C.ResultsHashGetValuesString(
		results.CPtr,
		C.CString(propertyName),
		&buffer[0],
		C.ulong(size),
		C.CString(separator),
		e.CPtr)
	// Check err
	if err != nil {
		return "", uint64(actualSize), err
	}

	// Check exception
	if !e.IsOkay() {
		return "", uint64(actualSize), fmt.Errorf(
			C.GoString(C.ExceptionGetMessage(e.CPtr)))
	}

	return C.GoString(&buffer[0]), uint64(actualSize), nil
}

// GetAvailableProperties returns a list of available properties. The index
// of this property can be used obtain further details about the property.
func (results *ResultsHash) GetAvailableProperties() []string {
	return nil
}

// NewPropertyName returns the name of a property at a given index.
func (results *ResultsHash) GetPropertyName(
	requiredPropertyIndex int) string {
	return ""
}

// GetPropertyType returns the type that a value of a property determined by
// a given proprety index can have.
func (results *ResultsHash) GetPropertyType(
	requiredPropertyIndex int) reflect.Type {
	return nil
}

/* Detections */
// MatchUserAgent performs a detection on an User-Agent string. This matches
// the C API fiftyoneDegreesResultsHashFromUserAgent.
func (results *ResultsHash) MatchUserAgent(ua string) error {
	cUserAgent := C.CString(ua)
	e := NewException()
	_, err := C.ResultsHashFromUserAgent(
		results.CPtr,
		cUserAgent,
		C.strlen(cUserAgent),
		e.CPtr)
	// Check error
	if err != nil {
		return err
	}

	// Check exception
	if !e.IsOkay() {
		return fmt.Errorf(C.GoString(C.ExceptionGetMessage(e.CPtr)))
	}

	return nil
}

// MatchEvidence performs a detection on evidence encapsulated in a EvidenceHash
// object. This matches the C API fiftyoneDegreesResultsHashFromEvidence.
func (results *ResultsHash) MatchEvidence(e Evidence) error {
	return nil
}

// MatchDeviceId performs a detection on a given device id. This matches
// the C API fiftyoneDegreesResultsHashFromDeviceId
func (results *ResultsHash) MatchDeviceId(id string) error {
	return nil
}
