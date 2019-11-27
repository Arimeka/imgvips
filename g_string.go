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

var gStringPool = sync.Pool{
	New: func() interface{} {
		var gValue C.GValue

		v := &GValue{
			gType:  C.G_TYPE_STRING,
			gValue: &gValue,
		}

		C.g_value_init(v.gValue, v.gType)

		return v
	},
}

// String return string gValue, if type is GString.
// If type not match, ok will return false.
// If gValue already freed, gValue will be empty string, ok will be true.
func (v *GValue) String() (value string, ok bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.gType != C.G_TYPE_STRING {
		return "", false
	}

	return C.GoString((*C.char)(unsafe.Pointer(C.g_value_get_string(v.gValue)))), true
}

// GString transform string gValue to glib gValue
func GString(value string) *GValue {
	v := gStringPool.Get().(*GValue)
	v.freed = false

	cStr := C.CString(value)
	defer C.free(unsafe.Pointer(cStr))
	C.g_value_set_string(v.gValue, cStr)

	if v.free == nil {
		v.free = func(val *GValue) {
			if val.freed {
				return
			}
			C.g_value_reset(val.gValue)
			gStringPool.Put(val)
		}
	}
	if v.copy == nil {
		v.copy = func(val *GValue) (*GValue, error) {
			newVal := gStringPool.Get().(*GValue)
			newVal.freed = false
			if newVal.free == nil {
				newVal.free = val.free
			}
			if newVal.copy == nil {
				newVal.copy = val.copy
			}

			str := C.GoString((*C.char)(unsafe.Pointer(C.g_value_get_string(val.gValue))))
			cStr := C.CString(str)
			defer C.free(unsafe.Pointer(cStr))

			C.g_value_set_string(newVal.gValue, cStr)

			return newVal, nil
		}
	}

	return v
}
