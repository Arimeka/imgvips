package imgvips

/*
#cgo pkg-config: vips
#include "stdlib.h"
#include "vips/vips.h"
*/
import "C"

import (
	"sync"
	"unsafe"
)

var gVipsBlobPool = sync.Pool{
	New: func() interface{} {
		var gValue C.GValue

		v := &GValue{
			gType:  C.vips_blob_get_type(),
			gValue: &gValue,
			copy: func(val *GValue) (*GValue, error) {
				return nil, ErrCopyForbidden
			},
		}

		C.g_value_init(v.gValue, v.gType)

		return v
	},
}

func freeVipsBlobAreaFn(val *GValue) {
	ptr := C.g_value_peek_pointer(val.gValue)
	if ptr != nil {
		// VipsBlob can be freed like *C.VipsArea
		C.vips_area_unref((*C.VipsArea)(ptr))
	}
	C.g_value_reset(val.gValue)
	val.free = freeVipsBlobNullFn

	gVipsBlobPool.Put(val)
}

func freeVipsBlobNullFn(val *GValue) {
	if val.freed {
		return
	}
	C.g_value_reset(val.gValue)

	gVipsBlobPool.Put(val)
}

// Bytes return bytes slice from GValue.
// It unset gValue after call, so for next call you get nil value
// If type not match, ok will return false.
// If VipsBlob already freed, return nil value, ok will be true.
func (v *GValue) Bytes() (value []byte, ok bool) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if v.gType != C.vips_blob_get_type() {
		return nil, false
	}

	ptr := C.g_value_peek_pointer(v.gValue)
	if ptr == nil {
		return nil, true
	}

	var gSize C.gsize
	result := C.vips_blob_get((*C.VipsBlob)(ptr), &gSize)

	// Copy data from *ptr and return []byte
	// Better create []byte from *ptr address, but it kinda horrible
	value = C.GoBytes(result, (C.int)(gSize))

	return value, true
}

// GVipsBlob create VipsBlob from bytes array.
//
// You must protect bytes array from GC and modification while using the VipsImage loaded from this blob.
//
// VipsBlob is used in load_buffer and save_buffer.
// VipsBlob is a boxed type, so we use g_value_set_boxed instead of g_value_set_object.
//
// Calling Copy() at GValue with type VipsBlob is forbidden.
func GVipsBlob(data []byte) *GValue {
	v := GNullVipsBlob()
	if len(data) == 0 {
		return v
	}

	v.free = freeVipsBlobAreaFn

	size := C.ulong(len(data))
	blob := C.vips_blob_new(nil, unsafe.Pointer(&data[0]), size)

	C.g_value_set_boxed(v.gValue, C.gconstpointer(blob))

	return v
}

// GNullVipsBlob create empty glib object gValue with type for *C.VipsBlob
// Calling Copy() at GValue with type VipsBlob is forbidden.
func GNullVipsBlob() *GValue {
	v := gVipsBlobPool.Get().(*GValue)
	v.freed = false

	if v.free == nil {
		v.free = freeVipsBlobNullFn
	}

	return v
}
