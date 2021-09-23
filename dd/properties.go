package dd

// #include <string.h>
// #include "./device-detection-cxx/src/hash/hash.h"
// #include "./device-detection-cxx/src/hash/fiftyone.h"
import "C"
import (
	"runtime"
	"unsafe"
)

// PropertiesRequired wraps around a pointer to a value of C PropertiesRequired
// structure.
type PropertiesRequired struct {
	CPtr *C.PropertiesRequired
}

// propertiesFinalizer check if C resource has been explicitly
// freed by Free method. Panic if it is not.
func propertiesFinalizer(props *PropertiesRequired) {
	if props.CPtr != nil {
		panic("ERROR: PropertiesRequired should be freed explicitly by " +
			"its Free method.")
	}
}

// NewPropertiesRequired creates a new PropertiesRequired object.
func NewPropertiesRequired(
	properties string) *PropertiesRequired {
	// Create C properties required
	cProps := C.PropertiesDefault
	cProps.string = C.CString(properties)
	props := &PropertiesRequired{&cProps}
	runtime.SetFinalizer(props, propertiesFinalizer)
	return props
}

// Free free the C resource and set CPtr to nil
func (properties *PropertiesRequired) Free() {
	C.free(unsafe.Pointer(properties.CPtr.string))
	properties.CPtr = nil
}

// Properties returns the required properties string
func (properties *PropertiesRequired) Properties() string {
	return C.GoString(properties.CPtr.string)
}
