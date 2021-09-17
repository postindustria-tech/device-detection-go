package dd

// #include <string.h>
// #include "./device-detection-cxx/src/common-cxx/fiftyone.h"
// #include "./device-detection-cxx/src/hash/hash.h"
// #include "./device-detection-cxx/src/hash/fiftyone.h"
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

// Maximum bit shift count for unsigned int 32.
const shiftCount32bit = 32

// ResultsHash wraps around a pointer to a value of C ResultsHash structure
type ResultsHash struct {
	CPtr     *C.ResultsHash
	CResults *interface{} // Pointer to a slice holding C results
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
		return ResultsHash{nil, nil}, err
	}

	// Map the items list to Go slice. This is done once so every access to
	// each result won't have to cast and create the slice again.
	var cResults interface{} = (*[1 << shiftCount32bit]C.ResultHash)(
		unsafe.Pointer(r.items))[:r.capacity:r.capacity]
	return ResultsHash{r, &cResults}, nil
}

// Free free the resource allocated in the C layer.
func (results *ResultsHash) Free() error {
	_, err := C.ResultsHashFree(results.CPtr)
	return err
}

/* Metric Getters */

// DeviceId returns the unique device id. This matches the C API
// fiftyoneDegreesHashGetDeviceIdFromResults.
func (results *ResultsHash) DeviceId() (id string, err error) {
	return "", nil
}

// DeviceIdByIndex returns the unique device id of a result pointed by a index.
func (results *ResultsHash) DeviceIdByIndex(index uint32) string {
	return ""
}

// Iterations returns the number of iterations carried out in order to find
// a match.
func (results *ResultsHash) Iterations() int {
	return 0
}

// Iterations returns the number of iterations carried out in order to find
// a match of a result pointed by an index.
func (results *ResultsHash) IterationsByIndex(index uint32) int {
	cResults := (*results.CResults).([]C.ResultHash)
	return int(cResults[index].iterations)
}

// Drift returns the maximum drift for a matched substring.
func (results *ResultsHash) Drift() int {
	return 0
}

// DriftByIndex returns the drift for a matched substring of a result pointed
// by an index.
func (results *ResultsHash) DriftByIndex(index uint32) int {
	cResults := (*results.CResults).([]C.ResultHash)
	return int(cResults[index].drift)
}

// Difference returns the total difference between the results returned and
// the target User-Agent.
func (results *ResultsHash) Difference() int {
	return 0
}

// DifferenceByIndex returns the difference between the result pointed by a
// index and the target User-Agent.
func (results *ResultsHash) DifferenceByIndex(index uint32) int {
	cResults := (*results.CResults).([]C.ResultHash)
	return int(cResults[index].difference)
}

// MatchedNodes returns the number of hash nodes matched within the evidence.
func (results *ResultsHash) MatchedNodes() int {
	return 0
}

// Method returns the method used to determine the match result.
func (results *ResultsHash) Method(index uint32) int {
	return 0
}

// MethodByIndex returns the method use to determine a match result pointed
// by a index.
func (results *ResultsHash) MethodByIndex(index uint32) int {
	return 0
}

// Trace returns the trace route in a readable format.
func (results *ResultsHash) Trace() string {
	return ""
}

// TraceByIndex returns the trace reoute in a readable format of a result
// pointed by a given index.
func (results *ResultsHash) TraceByIndex(index uint32) string {
	return ""
}

// UserAgents returns number of User-Agents that were used in the results.
func (results *ResultsHash) UserAgents() int {
	return int(results.CPtr.count)
}

// UserAgent return the user agent of a result pointed by a given index.
func (results *ResultsHash) UserAgent(index int) string {
	cResults := (*results.CResults).([]C.ResultHash)
	return C.GoString(cResults[index].b.matchedUserAgent)
}

// HasValues returns whether the last detection returns any matched value
// for a given property index. This matches the C API
// fiftyoneDegreesResultsHashGetHasValues
func (results *ResultsHash) HasValues(
	requiredPropertyIndex int) (r bool, err error) {
	e := NewException()
	hasValues, err := C.ResultsHashGetHasValues(
		results.CPtr, C.int(requiredPropertyIndex), e.CPtr)
	if err != nil {
		return false, err
	}

	// Check exception
	if !e.IsOkay() {
		return false, fmt.Errorf(C.GoString(C.ExceptionGetMessage(e.CPtr)))
	}
	return bool(hasValues), nil
}

// NoValueReasonMessage returns the no value message of a given property
// index. This matches a combination of the C APIs
// fiftyoneDegreesResultsNoValueReason and
// fiftyoneDegreesResultsHashGetNoValueReasonMessage
func (results *ResultsHash) NoValueReasonMessage(
	requiredPropertyIndex int) (m string, err error) {
	e := NewException()
	reason, err := C.ResultsHashGetNoValueReason(
		results.CPtr, C.int(requiredPropertyIndex), e.CPtr)
	if err != nil {
		return "", err
	}

	// Check exception
	if !e.IsOkay() {
		return "", fmt.Errorf(C.GoString(C.ExceptionGetMessage(e.CPtr)))
	}

	// Get no value reason message
	message, err := C.ResultsHashGetNoValueReasonMessage(reason)
	if err != nil {
		return "", err
	}

	return C.GoString(message), nil
}

// Values returns a list of values resulted from a detection for a given
// property index.
func (results *ResultsHash) Values(
	requiredPropertyIndex int) (values []string, err error) {
	return nil, nil
}

// ValuesString returns a string of all values resulted from a detection.
// If there are multiple values, they are separated by the specified separator.
func (results *ResultsHash) ValuesString(
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

// AvailableProperties returns a list of available properties. The index
// of this property can be used obtain further details about the property.
func (results *ResultsHash) AvailableProperties() (a []string, err error) {
	dataSet := (*C.DataSetHash)(results.CPtr.b.b.dataSet)
	cAvailable := dataSet.b.b.available
	availableCount := int(cAvailable.count)
	available := make([]string, 0, availableCount)
	// Loop through the C available properties list and add it to the
	// new Go available array.
	for i := 0; i < availableCount; i++ {
		name := C.PropertiesGetNameFromRequiredIndex(cAvailable, C.int(i))
		if err != nil {
			return available, err
		}

		available = append(available, C.GoString(&name.value))
	}
	return available, nil
}

// NewPropertyName returns the name of a property at a given index.
func (results *ResultsHash) PropertyName(
	requiredPropertyIndex int) string {
	return ""
}

// PropertyType returns the type that a value of a property determined by
// a given proprety index can have.
func (results *ResultsHash) PropertyType(
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
