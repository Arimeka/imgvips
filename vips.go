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

// Initialize libvips
//
// By default, libvips cache will be turned off (set to zero), vector arithmetic - turned on.
func Initialize(options ...InitOption) error {
	// Lock OS thread to current goroutine while initializing libvips
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	name := C.CString("imgvips")
	defer C.free(unsafe.Pointer(name))

	if success := C.vips_init(name); success != 0 {
		return errVipsFailedStart
	}

	opts := initOptions{}

	for _, option := range options {
		option.f(&opts)
	}

	vipsDetectMemoryLeak(opts.detectMemoryLeak)
	vipsVectorSetEnables(opts.enableVector)
	vipsCacheSetMax(opts.cacheMax)
	vipsCacheSetMaxMem(opts.cacheMaxMem)
	vipsConcurrencySet(opts.concurrency)

	return nil
}

// InitOption specifies an option for initialize libvips
type InitOption struct {
	f func(*initOptions)
}

type initOptions struct {
	detectMemoryLeak bool
	enableVector     bool
	cacheMax         int
	cacheMaxMem      int
	concurrency      int
}

// VipsDetectMemoryLeak turn on/off memory leak reports
func VipsDetectMemoryLeak(on bool) InitOption {
	return InitOption{func(options *initOptions) {
		options.detectMemoryLeak = true
	}}
}

func vipsDetectMemoryLeak(on bool) {
	var cBool = C.gboolean(0)
	if on {
		cBool = C.gboolean(1)
	}

	C.vips_leak_set(cBool)
}

// VipsCacheSetMax set maximum number of operation to cache
func VipsCacheSetMax(n int) InitOption {
	return InitOption{func(options *initOptions) {
		options.cacheMax = n
	}}
}

func vipsCacheSetMax(n int) {
	if n < 0 {
		n = 0
	}
	C.vips_cache_set_max(C.int(n))
}

// VipsCacheSetMaxMem set maximum amount of tracked memory
func VipsCacheSetMaxMem(n int) InitOption {
	return InitOption{func(options *initOptions) {
		options.cacheMaxMem = n
	}}
}

func vipsCacheSetMaxMem(n int) {
	if n < 0 {
		n = 0
	}
	C.vips_cache_set_max_mem(C.ulong(n))
}

// VipsVectorSetEnables enable fast vector path based on half-float arithmetic
func VipsVectorSetEnables(enabled bool) InitOption {
	return InitOption{func(options *initOptions) {
		options.enableVector = enabled
	}}
}

func vipsVectorSetEnables(enabled bool) {
	var cBool = C.gboolean(0)
	if enabled {
		cBool = C.gboolean(1)
	}

	C.vips_vector_set_enabled(cBool)
}

// VipsConcurrencySet set number of threads to use
func VipsConcurrencySet(n int) InitOption {
	return InitOption{func(options *initOptions) {
		options.concurrency = n
	}}
}

func vipsConcurrencySet(n int) {
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

// GetMem return libvips tracked memory
func GetMem() float64 {
	return float64(C.vips_tracked_get_mem())
}

// GetMemHighwater return libvips tracked memory high-water
func GetMemHighwater() float64 {
	return float64(C.vips_tracked_get_mem_highwater())
}

// GetAllocs return libvips tracked number of active allocations
func GetAllocs() float64 {
	return float64(C.vips_tracked_get_allocs())
}
