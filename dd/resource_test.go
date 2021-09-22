package dd

import "testing"

// Test that resource finalizer will panic if a pointer to
// C resource has not been freed
func TestResourceFinalizer(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Resource finalizer did not panic.")
		}
	}()
	manager := NewResourceManager()
	// Perform finalizer on live resource
	resourceFinalizer(manager)
	manager.Free()
}
