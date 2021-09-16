package dd

/*
* Cgo concatenate all C linker and compiler directives so create this file to
* centralize the specifications of these flags.
 */

/*
#cgo CFLAGS: -fcommon
#cgo LDFLAGS: -L./device-detection-cxx/build/lib -lfiftyone-hash-c -lfiftyone-device-detection-c -lfiftyone-common-c -lm
#cgo !darwin LDFLAGS: -latomic
*/
import "C"
