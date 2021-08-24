package dd

// #cgo LDFLAGS: -L./device-detection-cxx/build/lib -lfiftyone-hash-c -lfiftyone-device-detection-c -lfiftyone-common-c -lm -latomic
// #include <string.h>
// #include "./device-detection-cxx/src/hash/hash.h"
// #include "./device-detection-cxx/src/hash/fiftyone.h"
import "C"
import "unsafe"

// Performance Profile
type PerformanceProfile int

const (
	LowMemory PerformanceProfile = iota
	BalancedTemp
	Balanced
	HighPerformance
	InMemory
)

// ConfigHash wraps around pointer to a value of C ConfigHash structure
type ConfigHash struct {
	CPtr *C.ConfigHash
}

/* Constructor and Destructor */
func NewConfigHash() ConfigHash {
	config := C.HashDefaultConfig
	return ConfigHash{&config}
}

/* Setters */

// SetPerformanceProfile sets the performance strategy of a device detection
// engine.
func (config *ConfigHash) SetPerformanceProfile(perf PerformanceProfile) {}

// SetDifference sets the maximum difference in hash value to allow when
// finding hashes during the processing of HTTP headers.
func (config *ConfigHash) SetDifference(diff int32) {}

// SetDrift sets the maximum drift to allow when matching hashes.
func (config *ConfigHash) SetDrift(drift int32) {}

// SetUsePerformanceGraph sets whether performance optimized graph should be
// used for processing.
func (config *ConfigHash) SetUsePerformanceGraph(use bool) {}

// SetUsePredictiveGraph sets whether predictive optimized graph should be
// used for processing.
func (config *ConfigHash) SetUsePredictiveGraph(use bool) {}

// SetConcurrency sets the expected concurrent requests.
func (config *ConfigHash) SetConcurrency(concurrency uint16) {}

// SetTraceRoute sets whether the route through each graph should be traced
// during processing.
func (config *ConfigHash) SetTraceRoute(trace bool) {}

// SetUseUpperPrefixHeaders set whether or note the HTTP header might be
// prefixed with 'HTTP_'
func (config *ConfigHash) SetUseUpperPrefixHeaders(use bool) {}

/* Getters */

// GetPerformanceProfile get the configured performance profile
func (config *ConfigHash) GetPerformanceProfile() PerformanceProfile {
	return Balanced
}

// GetDifference returns the configured difference value
func (config *ConfigHash) GetDifference() int32 {
	return 0
}

// GetDrift returns the configured drift value
func (config *ConfigHash) GetDrift() int32 {
	return 0
}

// GetUsePerformanceGraph returns whether performance optimized graph should be
// used.
func (config *ConfigHash) GetUsePerformanceGraph() bool {
	return false
}

// GetUsePredictiveGraph returns whether predictive optmized graph should be
// used.
func (config *ConfigHash) GetUsePredictiveGraph() bool {
	return false
}

// GetConcurrency returns the configured concurrency
func (config *ConfigHash) GetConcurrency() uint16 {
	return 0
}

// GetTraceRoute returns whether route through each graph should be traced.
func (config *ConfigHash) GetTraceRoute() bool {
	return false
}

// GetUseUpperPrefixHeaders returns whether HTTP_ upper case prefixes should be
// considered when evaluating HTTP hreaders.
func (config *ConfigHash) GetUseUpperPrefixHeaders() bool {
	return false
}

/* Sizer */

// SizeManagerFromFile returns the size of the memory which will be used for
// holding a resource manager based on a given file. The input properties is a
// comma separated string list. This matches the C API
// fiftyoneDegreeshashSizeManagerFromFile.
func (config *ConfigHash) SizeManagerFromFile(
	properties string,
	fileName string) (size uint64, err error) {
	return 0, nil
}

// SizeManagerFromMemory returns the size of the memory which will be used for
// holding a resource manager based on a given memory. The given memory holds
// a data file content that is intended to initialize the resource amanger. The
// properties is a comma separated string list. This matches the C API
// fiftyoneDegreeshashSizeManagerFromMemory.
func (config *ConfigHash) SizeManagerFromMemory(
	properties string,
	memory unsafe.Pointer,
	size uint64) (s uint64, err error) {
	return 0, nil
}
