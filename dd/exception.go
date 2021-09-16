package dd

// #include <string.h>
// #include "./device-detection-cxx/src/hash/hash.h"
// #include "./device-detection-cxx/src/hash/fiftyone.h"
import "C"

// Exception wraps around a pointer to a value of C Exception structure
type Exception struct {
	CPtr *C.Exception
}

// NewException creates a new Exception object
func NewException() *Exception {
	ce := new(C.Exception)
	e := &Exception{ce}
	e.Clear()
	return e
}

// Clear resets the Exception object
func (e *Exception) Clear() {
	e.CPtr.file = nil
	e.CPtr._func = nil
	e.CPtr.line = C.int(-1)
	e.CPtr.status = C.FIFTYONE_DEGREES_STATUS_NOT_SET
}

// IsOkay check if an exception has been thrown.
func (e *Exception) IsOkay() bool {
	return (e.CPtr == nil ||
		e.CPtr.status == C.FIFTYONE_DEGREES_STATUS_NOT_SET)
}
