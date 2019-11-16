package imgvips

/*
#cgo pkg-config: vips
#include "vips/vips.h"
*/
import "C"

import (
	"errors"
)

var (
	errVipsFailedStart = errors.New("unable to start vips")
)

func init() {
	if success := C.vips_init(C.CString("imgvips")); success != 0 {
		panic(errVipsFailedStart)
	}
}

// VipsDetectMemoryLeak turn on/off memory leak reports
func VipsDetectMemoryLeak(on bool) {
	onVar := int8(0)
	if on {
		onVar = 1
	}

	C.vips_leak_set(C.int(onVar))
}

// VipsErrorFree clear error buffer
func VipsErrorFree() {
	C.vips_error_clear()
}

func vipsError() error {
	s := C.GoString(C.vips_error_buffer())
	C.vips_error_clear()
	C.vips_thread_shutdown()
	return errors.New(s)
}
