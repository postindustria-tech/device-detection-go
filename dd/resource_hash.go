package dd

// #include <string.h>
// #include "./device-detection-cxx/src/hash/hash.h"
// #include "./device-detection-cxx/src/hash/fiftyone.h"
import "C"
import (
	"fmt"
	"unsafe"
)

// InitFromFile initializes a resource manager based on a given data file path.
// The input properties is a comma separated string list. This matches the C API
// fiftyoneDegreesHashInitManagerFromFile.
func InitManagerFromFile(
	manager *ResourceManager,
	config ConfigHash,
	properties string,
	filePath string) error {

	e := NewException()
	cPath := C.CString(filePath)
	defer C.free(unsafe.Pointer(cPath))
	propsRequired := NewPropertiesRequired(properties)
	defer propsRequired.Free()
	s := C.HashInitManagerFromFile(
		manager.CPtr,
		config.CPtr,
		propsRequired.CPtr,
		cPath,
		e.CPtr)

	// Check exception
	if !e.IsOkay() {
		return fmt.Errorf(C.GoString(C.ExceptionGetMessage(e.CPtr)))
	}

	// Check status code
	if s != C.SUCCESS {
		cMessage := C.StatusGetMessage(s, cPath)
		defer C.MemoryStandardFree(unsafe.Pointer(cMessage))
		return fmt.Errorf(C.GoString(cMessage))
	}

	return nil
}

// InitFromMemory initializes a resource manager based on a given data file
// content held in the memory. The input properties is a comma separated string
// list. This matches the C API fiftyoneDegreesHashInitManagerFromFile.
func InitManagerFromMemory(
	manager *ResourceManager,
	config ConfigHash,
	properties string,
	memory unsafe.Pointer,
	size uint64) error {
	// TODO: To be implemented
	return nil
}
