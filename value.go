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

// GValue contains glib gValue and its type
type GValue struct {
	gType  C.GType
	gValue *C.GValue

	free func(val *GValue)
	copy func(val *GValue) (*GValue, error)
	mu   sync.RWMutex
}

// gValue return gValue *C.gValue
func (v *GValue) value() *C.GValue {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return v.gValue
}

// Copy create new instance of *GValue with new *C.gValue and run copy() func
func (v *GValue) Copy() (*GValue, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return v.copy(v)
}

// Free call free func for unref gValue
func (v *GValue) Free() {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.free(v)
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
	if v.gValue == nil {
		return false, true
	}

	val := C.g_value_get_boolean(v.gValue)
	if val != 1 {
		return false, true
	}

	return true, true
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
	if v.gValue == nil {
		return 0, true
	}

	return int(C.g_value_get_int(v.gValue)), true
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
	if v.gValue == nil {
		return 0, true
	}

	return float64(C.g_value_get_double(v.gValue)), true
}

// Double return string gValue, if type is GString.
// If type not match, ok will return false.
// If gValue already freed, gValue will be empty string, ok will be true.
func (v *GValue) String() (value string, ok bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.gType != C.G_TYPE_STRING {
		return "", false
	}
	if v.gValue == nil {
		return "", true
	}

	return C.GoString((*C.char)(unsafe.Pointer(C.g_value_get_string(v.gValue)))), true
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
	if v.gValue == nil {
		return nil, true
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

// GBoolean transform bool gValue to glib gValue
func GBoolean(value bool) *GValue {
	var (
		cBool  C.gboolean
		gValue C.GValue
	)

	if value {
		cBool = C.gboolean(1)
	} else {
		cBool = C.gboolean(0)
	}

	v := &GValue{
		gType:  C.G_TYPE_BOOLEAN,
		gValue: &gValue,
	}
	v.free = func(val *GValue) {
		if val.gValue == nil {
			return
		}

		C.g_value_unset(val.gValue)
		val.gValue = nil
	}
	v.copy = func(val *GValue) (*GValue, error) {
		var (
			gVal C.GValue
		)

		C.g_value_init(&gVal, val.gType)
		C.g_value_copy(val.gValue, &gVal)

		return &GValue{
			gType:  val.gType,
			gValue: &gVal,
			free:   val.free,
			copy:   val.copy,
		}, nil
	}

	C.g_value_init(v.gValue, v.gType)
	C.g_value_set_boolean(v.gValue, cBool)

	return v
}

// GInt transform int gValue to glib gValue
func GInt(value int) *GValue {
	var gValue C.GValue

	v := &GValue{
		gType:  C.G_TYPE_INT,
		gValue: &gValue,
	}
	v.free = func(val *GValue) {
		if val.gValue == nil {
			return
		}

		C.g_value_unset(val.gValue)
		val.gValue = nil
	}
	v.copy = func(val *GValue) (*GValue, error) {
		var (
			gVal C.GValue
		)

		C.g_value_init(&gVal, val.gType)
		C.g_value_copy(val.gValue, &gVal)

		return &GValue{
			gType:  val.gType,
			gValue: &gVal,
			free:   val.free,
			copy:   val.copy,
		}, nil
	}

	C.g_value_init(v.gValue, v.gType)
	C.g_value_set_int(v.gValue, C.gint(value))

	return v
}

// GDouble transform float64 gValue to glib gValue
func GDouble(value float64) *GValue {
	var gValue C.GValue

	v := &GValue{
		gType:  C.G_TYPE_DOUBLE,
		gValue: &gValue,
	}
	v.free = func(val *GValue) {
		if val.gValue == nil {
			return
		}

		C.g_value_unset(val.gValue)
		val.gValue = nil
	}
	v.copy = func(val *GValue) (*GValue, error) {
		var (
			gVal C.GValue
		)

		C.g_value_init(&gVal, val.gType)
		C.g_value_copy(val.gValue, &gVal)

		return &GValue{
			gType:  val.gType,
			gValue: &gVal,
			free:   val.free,
			copy:   val.copy,
		}, nil
	}

	C.g_value_init(v.gValue, v.gType)
	C.g_value_set_double(v.gValue, C.gdouble(value))

	return v
}

// GString transform string gValue to glib gValue
func GString(value string) *GValue {
	var gValue C.GValue

	cStr := C.CString(value)
	defer C.free(unsafe.Pointer(cStr))

	v := &GValue{
		gType:  C.G_TYPE_STRING,
		gValue: &gValue,
	}
	v.free = func(val *GValue) {
		if val.gValue == nil {
			return
		}

		C.g_value_unset(val.gValue)

		val.gValue = nil
	}
	v.copy = func(val *GValue) (*GValue, error) {
		var (
			gVal C.GValue
		)

		str := C.GoString((*C.char)(unsafe.Pointer(C.g_value_get_string(v.gValue))))
		cStr := C.CString(str)
		defer C.free(unsafe.Pointer(cStr))

		C.g_value_init(&gVal, val.gType)
		C.g_value_set_string(&gVal, cStr)

		return &GValue{
			gType:  val.gType,
			gValue: &gVal,
			free:   val.free,
			copy:   val.copy,
		}, nil
	}

	C.g_value_init(v.gValue, v.gType)
	C.g_value_set_string(v.gValue, cStr)

	return v
}

// GVipsImage return gValue, contains new empty VipsImage
func GVipsImage() *GValue {
	value := C.vips_image_new()
	v := GNullVipsImage()

	C.g_value_set_object(v.gValue, C.gpointer(value))

	return v
}

// GNullVipsImage create empty glib object gValue with type for *C.VipsImage
func GNullVipsImage() *GValue {
	var gValue C.GValue

	v := &GValue{
		gType:  C.vips_image_get_type(),
		gValue: &gValue,
	}

	v.free = func(val *GValue) {
		if val.gValue == nil {
			return
		}

		ptr := C.g_value_peek_pointer(val.gValue)
		if ptr != nil {
			C.g_object_unref(ptr)
		}
		C.g_value_unset(val.gValue)

		val.gValue = nil
	}

	v.copy = func(val *GValue) (*GValue, error) {
		var (
			gVal C.GValue
		)

		C.g_value_init(&gVal, val.gType)

		ptr := C.g_value_peek_pointer(val.gValue)
		if ptr == nil {
			return &GValue{
				gType:  val.gType,
				gValue: &gVal,
				free:   val.free,
				copy:   val.copy,
			}, nil
		}

		op, err := NewOperation("copy")
		if err != nil {
			return nil, err
		}
		defer op.Free()

		cIn := C.CString("in")
		defer C.free(unsafe.Pointer(cIn))
		cOut := C.CString("out")
		defer C.free(unsafe.Pointer(cOut))

		C.g_object_set_property((*C.GObject)(unsafe.Pointer(op.operation)), cIn, val.gValue)
		C.g_object_set_property((*C.GObject)(unsafe.Pointer(op.operation)), cOut, &gVal)

		success := C.vips_cache_operation_buildp(&op.operation)

		C.g_object_get_property((*C.GObject)(unsafe.Pointer(op.operation)), cOut, &gVal)

		if success != 0 {
			original := (*C.VipsImage)(C.g_value_get_object(v.gValue))
			ptr := C.g_value_peek_pointer(&gVal)
			if int(C.vips_image_get_width(original)) != 1 || int(C.vips_image_get_height(original)) != 1 || ptr == nil {
				if ptr != nil {
					C.g_object_unref(ptr)
				}
				C.g_value_unset(&gVal)

				return nil, vipsError()
			}

		}

		return &GValue{
			gType:  val.gType,
			gValue: &gVal,
			free:   val.free,
			copy:   val.copy,
		}, nil
	}

	C.g_value_init(v.gValue, C.vips_image_get_type())

	return v
}
