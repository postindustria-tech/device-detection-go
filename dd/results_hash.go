package dd

// #include <string.h>
// #include "./device-detection-cxx/src/common-cxx/fiftyone.h"
// #include "./device-detection-cxx/src/hash/hash.h"
// #include "./device-detection-cxx/src/hash/fiftyone.h"
import "C"
import (
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

// Maximum bit shift count for unsigned int 32.
const shiftCount32bit = 32

// Match method
type MatchMethod int

// Types of method
const (
	None MatchMethod = iota
	Performance
	Combined
	Predictive
)

var cMethodMap = map[C.HashMatchMethod]MatchMethod{
	C.FIFTYONE_DEGREES_HASH_MATCH_METHOD_NONE:        None,
	C.FIFTYONE_DEGREES_HASH_MATCH_METHOD_PERFORMANCE: Performance,
	C.FIFTYONE_DEGREES_HASH_MATCH_METHOD_COMBINED:    Combined,
	C.FIFTYONE_DEGREES_HASH_MATCH_METHOD_PREDICTIVE:  Predictive,
}

// ResultsHash wraps around a pointer to a value of C ResultsHash structure
type ResultsHash struct {
	CPtr     *C.ResultsHash
	CResults *interface{} // Pointer to a slice holding C results
}

/* Constructor and Destructor */

// resultsFinalizer check if C resource has been explicitly
// freed by Free method. Panic if it is not.
func resultsFinalizer(res *ResultsHash) {
	if res.CPtr != nil {
		panic("ERROR: ResultsHash should be freed explicitly by " +
			"its Free method.")
	}
}

// NewResultsHash create a new ResultsHash object. This matches the C API
// fiftyoneDegreesResultsHashCreate.
func NewResultsHash(
	manager *ResourceManager,
	uaCapacity uint32,
	overridesCapacity uint32) *ResultsHash {
	r := C.ResultsHashCreate(
		manager.CPtr,
		C.uint(uaCapacity),
		C.uint(overridesCapacity))

	// Map the items list to Go slice. This is done once so every access to
	// each result won't have to cast and create the slice again.
	var cResults interface{} = (*[1 << shiftCount32bit]C.ResultHash)(
		unsafe.Pointer(r.items))[:r.capacity:r.capacity]
	res := &ResultsHash{r, &cResults}
	runtime.SetFinalizer(res, resultsFinalizer)
	return res
}

// Free free the resource allocated in the C layer.
func (results *ResultsHash) Free() {
	C.ResultsHashFree(results.CPtr)
	results.CPtr = nil
}

/* Metric Getters */

// DeviceId returns the unique device id. This matches the C API
// fiftyoneDegreesHashGetDeviceIdFromResults.
func (results *ResultsHash) DeviceId() (id string, err error) {
	// TODO: To be implemented
	return "", nil
}

// DeviceIdByIndex returns the unique device id of a result pointed by a index.
func (results *ResultsHash) DeviceIdByIndex(index uint32) string {
	// TODO: To be implemented
	return ""
}

// Iterations returns the number of iterations carried out in order to find
// a match.
func (results *ResultsHash) Iterations() int {
	iterations := 0
	cResults := (*results.CResults).([]C.ResultHash)
	for _, cResult := range cResults {
		iterations += int(cResult.iterations)
	}
	return iterations
}

// Iterations returns the number of iterations carried out in order to find
// a match of a result pointed by an index.
func (results *ResultsHash) IterationsByIndex(index uint32) int {
	cResults := (*results.CResults).([]C.ResultHash)
	return int(cResults[index].iterations)
}

// Drift returns the maximum drift for a matched substring.
func (results *ResultsHash) Drift() int {
	drift := 0
	cResults := (*results.CResults).([]C.ResultHash)
	for _, cResult := range cResults {
		drift += int(cResult.drift)
	}
	return drift
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
	difference := 0
	cResults := (*results.CResults).([]C.ResultHash)
	for _, cResult := range cResults {
		difference += int(cResult.difference)
	}
	return difference
}

// DifferenceByIndex returns the difference between the result pointed by a
// index and the target User-Agent.
func (results *ResultsHash) DifferenceByIndex(index uint32) int {
	cResults := (*results.CResults).([]C.ResultHash)
	return int(cResults[index].difference)
}

// MatchedNodes returns the number of hash nodes matched within the evidence.
func (results *ResultsHash) MatchedNodes() int {
	matchedNodes := 0
	cResults := (*results.CResults).([]C.ResultHash)
	for _, cResult := range cResults {
		matchedNodes += int(cResult.matchedNodes)
	}
	return matchedNodes
}

// Method returns the method used to determine the match result.
func (results *ResultsHash) Method() MatchMethod {
	method := results.MethodByIndex(0)
	cResults := (*results.CResults).([]C.ResultHash)
	for _, cResult := range cResults {
		nextMethod := cMethodMap[cResult.method]
		if nextMethod > method {
			method = nextMethod
		}
	}
	return method
}

// MethodByIndex returns the method use to determine a match result pointed
// by a index.
func (results *ResultsHash) MethodByIndex(index uint32) MatchMethod {
	cResults := (*results.CResults).([]C.ResultHash)
	return cMethodMap[cResults[index].method]
}

// Trace returns the trace route in a readable format.
func (results *ResultsHash) Trace() string {
	// TODO: To be implemented
	return ""
}

// TraceByIndex returns the trace reoute in a readable format of a result
// pointed by a given index.
func (results *ResultsHash) TraceByIndex(index uint32) string {
	// TODO: To be implemented
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

// HasValuesByIndex returns whether the last detection returns any matched value
// for a given property index. This matches the C API
// fiftyoneDegreesResultsHashGetHasValues
func (results *ResultsHash) HasValuesByIndex(
	requiredPropertyIndex int) (hasValues bool, err error) {
	exp := NewException()
	r := C.ResultsHashGetHasValues(
		results.CPtr, C.int(requiredPropertyIndex), exp.CPtr)

	// Check exception
	if !exp.IsOkay() {
		return false, fmt.Errorf(C.GoString(C.ExceptionGetMessage(exp.CPtr)))
	}
	return bool(r), nil
}

// HasValuesByIndex returns whether the last detection returns any matched value
// for a given property name.
func (results *ResultsHash) HasValues(
	propertyName string) (hasValues bool, err error) {
	index := results.RequiredPropertyIndexFromName(propertyName)
	v, e := results.HasValuesByIndex(index)
	if e != nil {
		return false, e
	}
	return v, nil
}

// NoValueReasonMessageByIndex returns the no value message of a given property
// index. This matches a combination of the C APIs
// fiftyoneDegreesResultsNoValueReason and
// fiftyoneDegreesResultsHashGetNoValueReasonMessage
func (results *ResultsHash) NoValueReasonMessageByIndex(
	requiredPropertyIndex int) (message string, err error) {
	exp := NewException()
	reason := C.ResultsHashGetNoValueReason(
		results.CPtr, C.int(requiredPropertyIndex), exp.CPtr)

	// Check exception
	if !exp.IsOkay() {
		return "", fmt.Errorf(C.GoString(C.ExceptionGetMessage(exp.CPtr)))
	}

	// Get no value reason message
	m := C.ResultsHashGetNoValueReasonMessage(reason)

	return C.GoString(m), nil
}

// NoValueReasonMessage returns the no value message of a given property
// name.
func (results *ResultsHash) NoValueReasonMessage(
	propertyName string) (message string, err error) {
	index := results.RequiredPropertyIndexFromName(propertyName)
	reason, e := results.NoValueReasonMessageByIndex(index)
	if e != nil {
		return "", e
	}
	return reason, nil
}

// Values returns a list of values resulted from a detection for a given
// property index.
func (results *ResultsHash) Values(
	requiredPropertyIndex int) (values []string, err error) {
	// TODO: To be implemented
	return nil, nil
}

// ValuesString returns a string of all values resulted from a detection.
// If there are multiple values, they are separated by the specified separator.
func (results *ResultsHash) ValuesString(
	propertyName string,
	size uint64,
	separator string) (value string, actual uint64, err error) {
	// Create slice based on specified size
	buffer := make([]C.char, size)
	exp := NewException()
	actualSize := C.ResultsHashGetValuesString(
		results.CPtr,
		C.CString(propertyName),
		&buffer[0],
		C.ulong(size),
		C.CString(separator),
		exp.CPtr)

	// Check exception
	if !exp.IsOkay() {
		return "", uint64(actualSize), fmt.Errorf(
			C.GoString(C.ExceptionGetMessage(exp.CPtr)))
	}

	return C.GoString(&buffer[0]), uint64(actualSize), nil
}

// AvailableProperties returns a list of available properties. The index
// of this property can be used obtain further details about the property.
func (results *ResultsHash) AvailableProperties() []string {
	dataSet := (*C.DataSetHash)(results.CPtr.b.b.dataSet)
	cAvailable := dataSet.b.b.available
	availableCount := int(cAvailable.count)
	a := make([]string, 0, availableCount)
	// Loop through the C available properties list and add it to the
	// new Go available array.
	for i := 0; i < availableCount; i++ {
		name := C.PropertiesGetNameFromRequiredIndex(cAvailable, C.int(i))
		a = append(a, C.GoString(&name.value))
	}
	return a
}

// PropertyName returns the name of a property at a given index.
func (results *ResultsHash) PropertyName(
	requiredPropertyIndex int) string {
	dataSet := (*C.DataSetHash)(results.CPtr.b.b.dataSet)
	cAvailable := dataSet.b.b.available
	n := C.PropertiesGetNameFromRequiredIndex(
		cAvailable, C.int(requiredPropertyIndex))
	return C.GoString(&n.value)
}

// RequiredPropertyIndexFromName returns the required index of a property
// using its name
func (results *ResultsHash) RequiredPropertyIndexFromName(
	propertyName string) int {
	dataSet := (*C.DataSetHash)(results.CPtr.b.b.dataSet)
	cAvailable := dataSet.b.b.available
	cName := C.CString(propertyName)
	i := C.PropertiesGetRequiredPropertyIndexFromName(
		cAvailable, cName)
	C.free(unsafe.Pointer(cName))
	return int(i)
}

// PropertyType returns the type that a value of a property determined by
// a given proprety index can have.
func (results *ResultsHash) PropertyType(
	requiredPropertyIndex int) reflect.Type {
	// TODO: To be implemented
	return nil
}

/* Detections */
// MatchUserAgent performs a detection on an User-Agent string. This matches
// the C API fiftyoneDegreesResultsHashFromUserAgent.
func (results *ResultsHash) MatchUserAgent(ua string) error {
	cUserAgent := C.CString(ua)
	e := NewException()
	C.ResultsHashFromUserAgent(
		results.CPtr,
		cUserAgent,
		C.strlen(cUserAgent),
		e.CPtr)

	// Check exception
	if !e.IsOkay() {
		return fmt.Errorf(C.GoString(C.ExceptionGetMessage(e.CPtr)))
	}

	return nil
}

// MatchEvidence performs a detection on evidence encapsulated in a EvidenceHash
// object. This matches the C API fiftyoneDegreesResultsHashFromEvidence.
func (results *ResultsHash) MatchEvidence(e Evidence) error {
	// TODO: To be implemented
	return nil
}

// MatchDeviceId performs a detection on a given device id. This matches
// the C API fiftyoneDegreesResultsHashFromDeviceId
func (results *ResultsHash) MatchDeviceId(id string) error {
	// TODO: To be implemented
	return nil
}
