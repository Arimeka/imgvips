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

var gBooleanPool = sync.Pool{
	New: func() interface{} {
		var gValue C.GValue

		v := &GValue{
			gType:  C.G_TYPE_BOOLEAN,
			gValue: &gValue,
		}

		C.g_value_init(v.gValue, v.gType)

		return v
	},
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

	v := gBooleanPool.Get().(*GValue)
	v.freed = false
	C.g_value_set_boolean(v.gValue, cBool)

	if v.free == nil {
		v.free = func(val *GValue) {
			if val.freed {
				return
			}
			C.g_value_reset(val.gValue)
			gBooleanPool.Put(val)
		}
	}
	if v.copy == nil {
		v.copy = func(val *GValue) (*GValue, error) {
			newVal := gBooleanPool.Get().(*GValue)
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
