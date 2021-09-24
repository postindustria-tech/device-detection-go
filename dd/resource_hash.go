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
