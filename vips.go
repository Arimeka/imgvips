package imgvips

/*
#cgo pkg-config: vips
#include "vips/vips.h"
#include "vips/vector.h"
*/
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

var (
	errVipsFailedStart = errors.New("unable to start vips")
)

// nolint:gochecknoinits // Wanna do init()
func init() {
	// Lock OS thread to current goroutine while initializing libvips
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	name := C.CString("imgvips")
	defer C.free(unsafe.Pointer(name))

	if success := C.vips_init(name); success != 0 {
		panic(errVipsFailedStart)
	}
}

// VipsDetectMemoryLeak turn on/off memory leak reports
func VipsDetectMemoryLeak(on bool) {
	var cBool = C.gboolean(0)
	if on {
		cBool = C.gboolean(1)
	}

	C.vips_leak_set(cBool)
}

// VipsCacheSetMax set maximum number of operation to cache
func VipsCacheSetMax(n int) {
	if n < 0 {
		n = 0
	}
	C.vips_cache_set_max(C.int(n))
}

// VipsCacheSetMaxMem set maximum amount of tracked memory
func VipsCacheSetMaxMem(n int) {
	if n < 0 {
		n = 0
	}
	C.vips_cache_set_max_mem(C.ulong(n))
}

// VipsVectorSetEnables enable fast vector path based on half-float arithmetic
func VipsVectorSetEnables(enabled bool) {
	var cBool = C.gboolean(0)
	if enabled {
		cBool = C.gboolean(1)
	}

	C.vips_vector_set_enabled(cBool)
}

// VipsConcurrencySet set number of threads to use
func VipsConcurrencySet(n int) {
	if n < 0 {
		n = 0
	}
	C.vips_concurrency_set(C.int(n))
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
