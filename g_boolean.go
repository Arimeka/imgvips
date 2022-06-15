package imgvips

/*
#cgo pkg-config: vips
#include "stdlib.h"
#include "vips/vips.h"
*/
import "C"

func newGBoolean() *GValue {
	var gValue C.GValue

	v := &GValue{
		gType:  C.G_TYPE_BOOLEAN,
		gValue: &gValue,
		free: func(val *GValue) {
			if val.freed {
				return
			}
			C.g_value_unset(val.gValue)
			val.gType = C.G_TYPE_NONE
		},
		copy: func(val *GValue) (*GValue, error) {
			newVal := newGBoolean()

			C.g_value_copy(val.gValue, newVal.gValue)

			return newVal, nil
		},
	}

	C.g_value_init(v.gValue, v.gType)

	return v
}

// Boolean return boolean gValue, if type is GBoolean.
// If type not match, ok will return false.
// If gValue already freed, gValue will be false, ok will be true.
func (v *GValue) Boolean() (value, ok bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.gType != C.G_TYPE_BOOLEAN {
		return false, false
	}

	val := C.g_value_get_boolean(v.gValue)
	if val != 1 {
		return false, true
	}

	return true, true
}

// GBoolean transform bool gValue to glib gValue
func GBoolean(value bool) *GValue {
	var (
		cBool C.gboolean
	)

	if value {
		cBool = C.gboolean(1)
	} else {
		cBool = C.gboolean(0)
	}

	v := newGBoolean()
	C.g_value_set_boolean(v.gValue, cBool)

	return v
}
