package dd

// #cgo LDFLAGS: -L./device-detection-cxx/build/lib -lfiftyone-hash-c -lfiftyone-device-detection-c -lfiftyone-common-c -lm -latomic
// #include <string.h>
// #include "./device-detection-cxx/src/hash/hash.h"
// #include "./device-detection-cxx/src/hash/fiftyone.h"
import "C"
import (
	"fmt"
	"unsafe"
)

// InitFromFile initializes a resource manager based on a given data file. The
// input properties is a comma separated string list. This matches the C API
// fiftyoneDegreesHashInitManagerFromFile.
func InitManagerFromFile(
	manager *ResourceManager,
	config ConfigHash,
	properties string,
	fileName string) error {

	e := NewException()
	cName := C.CString(fileName)
	_, err := C.HashInitManagerFromFile(
		manager.CPtr,
		config.CPtr,
		NewPropertiesRequired(properties).CPtr,
		cName,
		e.CPtr)
	C.free(unsafe.Pointer(cName))

	// Check err
	if err != nil {
		return err
	}

	// Check exception
	if !e.IsOkay() {
		return fmt.Errorf(C.GoString(C.ExceptionGetMessage(e.CPtr)))
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
	return nil
}
