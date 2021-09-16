package dd

// #include <string.h>
// #include "./device-detection-cxx/src/common-cxx/fiftyone.h"
import "C"

// ResourceManagers wraps around a pointer to a value of C ResourceManager
// structure
type ResourceManager struct {
	CPtr *C.ResourceManager
}

// NewResourceManager creats a new object of ResourceManager
func NewResourceManager() *ResourceManager {
	manager := new(C.ResourceManager)
	return &ResourceManager{manager}
}

// Free frees the native resources allocated in the C layer for a Resource
// Manager.
func (manager *ResourceManager) Free() error {
	_, err := C.ResourceManagerFree(manager.CPtr)
	return err
}
