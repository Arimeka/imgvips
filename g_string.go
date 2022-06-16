package imgvips

/*
#cgo pkg-config: vips
#include "stdlib.h"
#include "vips/vips.h"
*/
import "C"

import (
	"unsafe"
)

func newGString() *GValue {
	var gValue C.GValue

	v := &GValue{
		gType:  C.G_TYPE_STRING,
		gValue: &gValue,
		free: func(val *GValue) {
			if val.gValue == nil {
				return
			}
			C.g_value_unset(val.gValue)
			val.gType = C.G_TYPE_NONE
		},
		copy: func(val *GValue) (*GValue, error) {
			newVal := newGString()

			str := C.GoString((*C.char)(unsafe.Pointer(C.g_value_get_string(val.gValue))))
			cStr := C.CString(str)

			C.g_value_set_string(newVal.gValue, cStr)
			C.free(unsafe.Pointer(cStr))

			return newVal, nil
		},
	}

	C.g_value_init(v.gValue, v.gType)

	return v
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
	v := newGString()

	cStr := C.CString(value)
	defer C.free(unsafe.Pointer(cStr))
	C.g_value_set_string(v.gValue, cStr)

	return v
}
