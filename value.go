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

// NewRawGValue creates GValue from raw glib value
// It can be used to create GValue from something not implemented by the functions below.
// free func called in GValue.Free().
// copy func called in GValue.Copy().
func NewRawGValue(gValue *C.GValue, gType C.GType, freeFn func(*GValue), copyFn func(*GValue) (*GValue, error)) *GValue {
	return &GValue{
		gType:  gType,
		gValue: gValue,
		free:   freeFn,
		copy:   copyFn,
	}
}

// GValue contains glib value and its type
type GValue struct {
	gType  C.GType
	gValue *C.GValue

	free func(val *GValue)
	copy func(val *GValue) (*GValue, error)
	mu   sync.RWMutex
}

// GType return value C.GType
func (v *GValue) GType() C.GType {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return v.gType
}

// GValue return value *C.GValue
func (v *GValue) GValue() *C.GValue {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return v.gValue
}

// Copy create new instance of *GValue with new *C.GValue and run copy() func
func (v *GValue) Copy() (*GValue, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return v.copy(v)
}

// Free call free func for unref value
func (v *GValue) Free() {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.free(v)
}

// Boolean return boolean value, if type is GBoolean.
// If type not match, ok will return false.
// If value already freed, value will be false, ok will be true.
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

// Int return int value, if type is GInt.
// If type not match, ok will return false.
// If value already freed, value will be 0, ok will be true.
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

// Double return float64 value, if type is GDouble.
// If type not match, ok will return false.
// If value already freed, value will be 0, ok will be true.
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

// Double return string value, if type is GString.
// If type not match, ok will return false.
// If value already freed, value will be empty string, ok will be true.
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

// Image return *C.VipsImage value, if type is *C.VipsImage.
// If type not match, ok will return false.
// If value already freed, value will be nil, ok will be true.
func (v *GValue) Image() (value *C.VipsImage, ok bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	if v.gType != C.vips_image_get_type() {
		return nil, false
	}
	if v.gValue == nil {
		return nil, true
	}

	return (*C.VipsImage)(C.g_value_get_object(v.gValue)), true
}

// GBoolean transform bool value to glib value
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

		return NewRawGValue(&gVal, val.gType, val.free, val.copy), nil
	}

	C.g_value_init(v.gValue, v.gType)
	C.g_value_set_boolean(v.gValue, cBool)

	return v
}

// GInt transform int value to glib value
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

		return NewRawGValue(&gVal, val.gType, val.free, val.copy), nil
	}

	C.g_value_init(v.gValue, v.gType)
	C.g_value_set_int(v.gValue, C.gint(value))

	return v
}

// GDouble transform float64 value to glib value
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

		return NewRawGValue(&gVal, val.gType, val.free, val.copy), nil
	}

	C.g_value_init(v.gValue, v.gType)
	C.g_value_set_double(v.gValue, C.gdouble(value))

	return v
}

// GString transform string value to glib value
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
		C.g_value_set_string(&gVal, (*C.gchar)(cStr))

		return NewRawGValue(&gVal, val.gType, val.free, val.copy), nil
	}

	C.g_value_init(v.gValue, v.gType)
	C.g_value_set_string(v.gValue, (*C.gchar)(cStr))

	return v
}

// GVipsImage return GValue, contains new empty VipsImage
func GVipsImage() *GValue {
	value := NewVipsImage()
	v := GNullVipsImage()

	C.g_value_set_object(v.gValue, C.gpointer(value))

	return v
}

// GNullVipsImage create empty glib object value with type for *C.VipsImage
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
			return NewRawGValue(&gVal, val.gType, val.free, val.copy), nil
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
			if ImageHeight(original) != 1 || ImageWidth(original) != 1 || ptr == nil {
				if ptr != nil {
					C.g_object_unref(ptr)
				}
				C.g_value_unset(&gVal)

				return nil, vipsError()
			}

		}

		return NewRawGValue(&gVal, val.gType, val.free, val.copy), nil
	}

	C.g_value_init(v.gValue, C.vips_image_get_type())

	return v
}
