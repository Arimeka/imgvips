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

var gVipsImagePool = sync.Pool{
	New: func() interface{} {
		var gValue C.GValue

		v := &GValue{
			gType:  C.vips_image_get_type(),
			gValue: &gValue,
		}

		C.g_value_init(v.gValue, v.gType)

		return v
	},
}

// Image return *Image, if type is *C.VipsImage.
// If type not match, ok will return false.
// If gValue already freed, gValue will be nil, ok will be true.
func (v *GValue) Image() (value *Image, ok bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.gType != C.vips_image_get_type() {
		return nil, false
	}
	ptr := C.g_value_peek_pointer(v.gValue)
	if ptr == nil {
		return nil, true
	}

	return &Image{
		image: (*C.VipsImage)(C.g_value_get_object(v.gValue)),
		val:   v,
	}, true
}

// GVipsImage return gValue, contains new empty *C.VipsImage.
//
// Calling Copy() at empty *C.VipsImage will return error.
func GVipsImage() *GValue {
	value := C.vips_image_new()
	v := GNullVipsImage()

	C.g_value_set_object(v.gValue, C.gpointer(value))

	return v
}

// GNullVipsImage create empty glib object gValue with type for *C.VipsImage.
//
// Calling Copy() at empty *C.VipsImage will return error.
func GNullVipsImage() *GValue {
	v := gVipsImagePool.Get().(*GValue)
	v.freed = false

	if v.free == nil {
		v.free = func(val *GValue) {
			if val.freed {
				return
			}
			ptr := C.g_value_peek_pointer(val.gValue)
			if ptr != nil {
				C.g_object_unref(ptr)
			}
			C.g_value_reset(val.gValue)
			gVipsImagePool.Put(val)
		}
	}
	if v.copy == nil {
		v.copy = func(val *GValue) (*GValue, error) {
			newVal := gVipsImagePool.Get().(*GValue)
			newVal.freed = false
			if newVal.free == nil {
				newVal.free = val.free
			}
			if newVal.copy == nil {
				newVal.copy = val.copy
			}

			ptr := C.g_value_peek_pointer(val.gValue)
			if ptr == nil {
				return newVal, nil
			}

			op, err := NewOperation("copy")
			if err != nil {
				return nil, err
			}
			defer op.Free()

			cIn := cStringsCache.get("in")
			cOut := cStringsCache.get("out")

			C.g_object_set_property((*C.GObject)(unsafe.Pointer(op.operation)), cIn, val.gValue)
			C.g_object_set_property((*C.GObject)(unsafe.Pointer(op.operation)), cOut, newVal.gValue)

			newOp := C.vips_cache_operation_build(op.operation)
			if newOp == nil {
				C.g_object_get_property((*C.GObject)(unsafe.Pointer(op.operation)), cOut, newVal.gValue)
				newVal.free(newVal)

				return nil, vipsError()
			}

			C.g_object_unref(C.gpointer(op.operation))
			op.operation = newOp

			C.g_object_get_property((*C.GObject)(unsafe.Pointer(op.operation)), cOut, newVal.gValue)

			return newVal, nil
		}
	}

	return v
}
