package ddonpremise

/*
To run this program, follow the following steps:
1. Build device-detection-cxx using CMake as described in its README.md
2. Update the linker option '-L' to point to the location of the built static libraries.
3. Naviaget to the parent folder, which mean 'cd ..'
4. Build the executable file using.
  go build -o [executable file name] ./main
5. Run the output file.
  ./[executable file name]
*/

// #cgo LDFLAGS: -L./device-detection-cxx/build/lib -lfiftyone-hash-c -lfiftyone-device-detection-c -lfiftyone-common-c -lm -latomic
// #include <string.h>
// #include "device-detection-cxx/src/hash/hash.h"
// #include "device-detection-cxx/src/hash/fiftyone.h"
import "C"
import "unsafe"
import "go/types"

type ResultsHash struct {
	cptr *C.ResultsHash
}

type ConfigHash struct {
	cptr *C.ConfigHash
}

type PropertiesRequired struct {
	cptr *C.PropertiesRequired
}

type ResourceManager struct {
	cptr *C.ResourceManager
}

type EvidenceKeyValuePairArray struct {
	cptr *C.EvidenceKeyValuePairArray
}

type ResultsNoValueReason struct {
	reason C.fiftyoneDegreesResultsNoValueReason
}

type Values struct {
	cptr *C.fiftyoneDegreesCollectionItem
}

type DataSetHash struct {
	cptr *C.DataSetHash
}

type ResultHash struct {
	cptr *C.ResultHash
}

type EvidenceKeyValuePair struct {
	cptr *C.EvidenceKeyValuePair
}

type MemoryReader struct {
	cptr *C.MemoryReader
}

type EvidencePrefix int
const (
	httpHeaderString EvidencePrefix = 1 << iota
	httpIpAddresses = 1 << iota
	httpEvidenceServer = 1 << iota
	httpEvidenceQuery = 1 << iota
	httpEvidenceCookie = 1 << iota
	httpEvidenceIgnore = 1 << 7
)

type Property struct {
	cptr *C.Property
}

type PropertiesAvailable struct {
	cptr *C.PropertiesAvailable
}

type Properties struct {
	cptr *C.fiftyoneDegreesCollection
}

func FiftyoneDegreesResultsHashGetValuesString (
	result ResultsHash,
	propertyName string,
	buffer unsafe.Pointer,
	size uint64,
	separator string) (s uint64, err error) {
	return 0, nil
}

func FiftyoneDegreeshashSizeManagerFromFile (
	config ConfigHash,
	properties PropertiesRequired,
	fileName string) (size uint64, err error) {
	return 0, nil
}
	
func FiftyoneDegreeshashSizeManagerFromMemory (
	config ConfigHash,
	properties PropertiesRequired,
	memory unsafe.Pointer,
	size uint64) (s uint64, err error) {
	return 0, nil
}

func FiftyoneDegreesHashInitManagerFromFile (
	manager ResourceManager,
	config ConfigHash,
	properties PropertiesRequired,
	fileName string) error {
	return nil
}

func FiftyoneDegreesHashInitManagerFromMemory (
	manager ResourceManager,
	config ConfigHash,
	properties PropertiesRequired,
	memory unsafe.Pointer,
	size uint64) error {
	return nil
}

func FiftyoneDegreesResultsHashCreate (
	manager ResourceManager,
	uaCapacity uint64,
	overridesCapacity uint64) ResultsHash {
	return ResultsHash{nil}
}

func FiftyoneDegreesResultsHashFree (
	results ResultsHash) {}

func FiftyoneDegreesResultsHashFromUserAgent (
	results ResultsHash,
	userAgent string,
	evidence EvidenceKeyValuePairArray) error {
	return nil
}

func FiftyoneDegreesResultsHashFromEvidence (
	results ResultsHash,
	evidence EvidenceKeyValuePairArray) error {
	return nil
}

func FiftyoneDegreesResultsHashFromDeviceId (
	results ResultsHash,
	deviceId string) error {
	return nil
}

func FiftyoneDegreesResultsHashGetHasValues (
	results ResultsHash,
	requiredPropertyIndex int) (r bool, err error) {
	return false, nil
}

func FiftyoneDegreesResultsNoValueReason (
	results ResultsHash,
	requiredPropertyIndex int) (r ResultsNoValueReason, err error) {
	return ResultsNoValueReason{0}, nil
}

func FiftyoneDegreesResultsHashGetNoValueReasonMessage (
	reason ResultsNoValueReason) string {
	return ""
}

func FiftyoneDegreesResultsHashGetValues (
	results ResultsHash,
	requiredPropertyIndex int) (vs Values, err error) {
	return Values{nil}, nil
}

func FiftyoneDegreesResourceManagerFree (
	manager ResourceManager) {}

func FiftyoneDegreesHashGetDeviceIdFromResults (
	results ResultsHash,
	destination unsafe.Pointer,
	size uint64) (id string, err error) {
	return "", nil
}

func FiftyoneDegreesHashGetDeviceIdFromResult (
	dataSet DataSetHash,
	result ResultHash,
	destination unsafe.Pointer,
	size uint64) (id string, err error) {
	return "", nil
}

func FiftyoneDegreesEvidenceCreate (
	capacity uint32) EvidenceKeyValuePairArray {
	return EvidenceKeyValuePairArray{nil}
}

func FiftyoneDegreesEvidenceFree (
	evidence EvidenceKeyValuePairArray) {}

func FiftyoneDegreesEvidenceAddString (
	evidence EvidenceKeyValuePairArray,
	prefix EvidencePrefix,
	field string,
	originalValue string) EvidenceKeyValuePair {
	return EvidenceKeyValuePair{nil}
}

func FiftyoneDegreesHashReloadManagerFromOriginalFile (
	manager ResourceManager) error {
	return nil
}

func FiftyoneDegreesHashReloadManagerFromFile (
	manager ResourceManager,
	fileName string) error {
	return nil
}

func FiftyoneDegreesHashReloadManagerFromMemory (
	manager ResourceManager,
	source unsafe.Pointer,
	length uint64) error {
	return nil
}

func FiftyoneDegreesFileReadToByteArray (
	fileName string,
	reader MemoryReader) error {
	return nil
}

func FiftyoneDegreesPropertyGetByName (
	dataSet *DataSetHash,
	propertyName string) (p Property, err error) {
	return Property{nil}, nil
}

func FiftyoneDegreesPropertiesGetRequiredPropertyIndexFromName (
	available PropertiesAvailable,
	propertyName string) int {
	return 0
}

func FiftyoneDegreesDataSetHashGet (
	manager ResourceManager) DataSetHash {
	return DataSetHash{nil}
}

func FiftyoneDegreesDataSetHashRelease (
	dataSet DataSetHash) {}

func FiftyoneDegreesPropertyGetValueType (
	properites Properties,
	index uint32) (t types.Type, err error) {
	return nil, nil
}

func FiftyoneDegreesGetResponseHeaders (
	results ResultsHash) (headers map[string]string, err error) {
	return nil, nil
}
