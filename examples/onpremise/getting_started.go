package main

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
import "fmt"

func main() {
	var m C.ResourceManager
    ps := C.PropertiesDefault
    ch := C.HashDefaultConfig
    p := C.CString("/home/fiftyonedegrees/staff/tung/device-detection-go/main/device-detection-cxx/device-detection-data/51Degrees-LiteV4.1.hash")
    var e *C.fiftyoneDegreesException
    e = nil
    s, err := C.HashInitManagerFromFile(
        &m,
        &ch,
        &ps,
        p,
        e)
    if err != nil {
       fmt.Errorf("%v", err)
    }

    if s == 1 {
       fmt.Printf("Do something")
    }

    r := C.ResultsHashCreate(&m, 1, 0)
    ua := C.CString("Mozilla/5.0 (iPhone; CPU iPhone OS 7_1 like Mac OS X) AppleWebKit/537.51.2 (KHTML, like Gecko) Version/7.0 Mobile/11D167 Safari/9537.53")
    C.ResultsHashFromUserAgent(
        r,
        ua,
        C.strlen(ua),
        e)

    var b [100]C.char
    C.ResultsHashGetValuesString(
       r,
       C.CString("IsMobile"),
       &b[0],
       100,
       C.CString(","),
       e)
    fmt.Printf("%s\n", C.GoString(&b[0]))

    C.ResultsHashFree(r)
    C.ResourceManagerFree(&m)
}
