package dd

// #include <string.h>
// #include "./device-detection-cxx/src/common-cxx/fiftyone.h"
import "C"
import (
	"runtime"
)

// ResourceManagers wraps around a pointer to a value of C ResourceManager
// structure
type ResourceManager struct {
	CPtr *C.ResourceManager
}

// Finalizer function for Resource Manager
func resourceFinalizer(m *ResourceManager) {
	if m.CPtr != nil {
		panic("ERROR: ResourceManager should be explicitly freed " +
			"by its Free method.")
	}
}

// NewResourceManager creats a new object of ResourceManager
func NewResourceManager() *ResourceManager {
	cManager := new(C.ResourceManager)
	manager := &ResourceManager{cManager}
	runtime.SetFinalizer(manager, resourceFinalizer)
	return manager
}

// Free frees the native resources allocated in the C layer for a Resource
// Manager.
func (manager *ResourceManager) Free() error {
	_, err := C.ResourceManagerFree(manager.CPtr)
	// If successfully freed the resource manager. Set the pointer to nil.
	// If not, keep the pointer for future reference.
	if err == nil {
		manager.CPtr = nil
	}
	return err
}
