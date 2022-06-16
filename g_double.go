package imgvips

/*
#cgo pkg-config: vips
#include "stdlib.h"
#include "vips/vips.h"
*/
import "C"

func newGDouble() *GValue {
	var gValue C.GValue

	v := &GValue{
		gType:  C.G_TYPE_DOUBLE,
		gValue: &gValue,
		free: func(val *GValue) {
			if val.gValue == nil {
				return
			}
			C.g_value_unset(val.gValue)
			val.gType = C.G_TYPE_NONE
		},
		copy: func(val *GValue) (*GValue, error) {
			newVal := newGDouble()

			C.g_value_copy(val.gValue, newVal.gValue)

			return newVal, nil
		},
	}

	C.g_value_init(v.gValue, v.gType)

	return v
}

// Double return float64 gValue, if type is GDouble.
// If type not match, ok will return false.
// If gValue already freed, gValue will be 0, ok will be true.
func (v *GValue) Double() (value float64, ok bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.gType != C.G_TYPE_DOUBLE {
		return 0, false
	}

	return float64(C.g_value_get_double(v.gValue)), true
}

// GDouble transform float64 gValue to glib gValue
func GDouble(value float64) *GValue {
	v := newGDouble()
	C.g_value_set_double(v.gValue, C.gdouble(value))

	return v
}
