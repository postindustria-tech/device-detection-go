package dd

/*
* Cgo concatenate all C linker and compiler directives so create this file to
* centralize the specifications of these flags.
 */

// #cgo LDFLAGS: -L./device-detection-cxx/build/lib -lfiftyone-hash-c -lfiftyone-device-detection-c -lfiftyone-common-c -lm -latomic
import "C"
