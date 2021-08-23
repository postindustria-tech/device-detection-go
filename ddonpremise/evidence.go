package ddonpremise

// #cgo LDFLAGS: -L./device-detection-cxx/build/lib -lfiftyone-hash-c -lfiftyone-device-detection-c -lfiftyone-common-c -lm -latomic
// #include <string.h>
// #include "device-detection-cxx/src/hash/hash.h"
// #include "device-detection-cxx/src/hash/fiftyone.h"
import "C"

type EvidenceHash struct {
	CPtr *C.EvidenceKeyValuePairArray
}
