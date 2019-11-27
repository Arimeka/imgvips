package imgvips

/*
#cgo pkg-config: vips
#include "stdlib.h"
#include "vips/vips.h"
*/
import "C"

import (
	"sync"
)

var gIntPool = sync.Pool{
	New: func() interface{} {
		var gValue C.GValue

		v := &GValue{
			gType:  C.G_TYPE_INT,
			gValue: &gValue,
		}

		C.g_value_init(v.gValue, v.gType)

		return v
	},
}

// Int return int gValue, if type is GInt.
// If type not match, ok will return false.
// If gValue already freed, gValue will be 0, ok will be true.
func (v *GValue) Int() (value int, ok bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.gType != C.G_TYPE_INT {
		return 0, false
	}

	return int(C.g_value_get_int(v.gValue)), true
}

// GInt transform int gValue to glib gValue
func GInt(value int) *GValue {
	v := gIntPool.Get().(*GValue)
	v.freed = false
	C.g_value_set_int(v.gValue, C.gint(value))

	if v.free == nil {
		v.free = func(val *GValue) {
			if val.freed {
				return
			}
			C.g_value_reset(val.gValue)
			gIntPool.Put(val)
		}
	}
	if v.copy == nil {
		v.copy = func(val *GValue) (*GValue, error) {
			newVal := gIntPool.Get().(*GValue)
			newVal.freed = false
			if newVal.free == nil {
				newVal.free = val.free
			}
			if newVal.copy == nil {
				newVal.copy = val.copy
			}

			C.g_value_copy(val.gValue, newVal.gValue)

			return newVal, nil
		}
	}

	return v
}
