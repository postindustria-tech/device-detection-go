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
