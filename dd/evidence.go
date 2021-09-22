package dd

// #include <string.h>
// #include "./device-detection-cxx/src/common-cxx/fiftyone.h"
import "C"

type EvidencePrefix int

const (
	HttpHeaderString   EvidencePrefix = 1 << iota
	HttpIpAddresses                   = 1 << iota
	HttpEvidenceServer                = 1 << iota
	HttpEvidenceQuery                 = 1 << iota
	HttpEvidenceCookie                = 1 << iota
	HttpEvidenceIgnore                = 1 << 7
)

type Evidence struct {
	CPtr *C.EvidenceKeyValuePairArray
}

// NewEvidenceHash returns a new EvidenceHash object with a given capacity
func NewEvidenceHash(capacity uint32) (evidence *Evidence, err error) {
	// TODO: To be implemented
	return nil, nil
}

// Free frees the evidence resources allocated in the C layer. This matches the
// C API fiftyoneDegreesEvidenceFree
func (evidence *Evidence) Free() error {
	// TODO: To be implemented
	return nil
}

// Add adds a new evidence to the object. This matches the C API
// fiftyoneDegreesEvidenceAddString
func (evidence *Evidence) Add(
	prefix EvidencePrefix,
	key string,
	value string) error {
	// TODO: To be implemented
	return nil
}
