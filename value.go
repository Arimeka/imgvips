package imgvips

/*
#cgo pkg-config: vips
#include "stdlib.h"
#include "vips/vips.h"
*/
import "C"

import (
	"errors"
	"sync"
	"unsafe"
)

var (
	// ErrCopyForbidden returns when GValue forbid copy function
	ErrCopyForbidden = errors.New("copy forbidden for this type")
)

// Value is interface for create own value for operation argument
type Value interface {
	// Free freed data in C
	Free()
	// Ptr return unsafe pointer to *C.GValue
	Ptr() unsafe.Pointer
}

// GValue contains glib gValue and its type
type GValue struct {
	gType  C.GType
	gValue *C.GValue

	free func(val *GValue)
	copy func(val *GValue) (*GValue, error)
	mu   sync.RWMutex
}

// Ptr return unsafe pointer to *C.GValue
func (v *GValue) Ptr() unsafe.Pointer {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return unsafe.Pointer(v.gValue)
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
	if v.gValue == nil {
		return nil, true
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
	v.free = func(val *GValue) {
		if val.gValue == nil {
			return
		}

		ptr := C.g_value_peek_pointer(val.gValue)
		if ptr != nil {
			// VipsBlob can be freed like *C.VipsArea
			C.vips_area_unref((*C.VipsArea)(ptr))
		}
		C.g_value_unset(val.gValue)

		val.gValue = nil
	}
	if len(data) == 0 {
		return v
	}

	size := C.ulong(len(data))
	blob := C.vips_blob_new(nil, unsafe.Pointer(&data[0]), size)

	C.g_value_set_boxed(v.gValue, C.gconstpointer(blob))

	return v
}

// GNullVipsBlob create empty glib object gValue with type for *C.VipsBlob
// Calling Copy() at GValue with type VipsBlob is forbidden.
func GNullVipsBlob() *GValue {
	var gValue C.GValue

	v := &GValue{
		gType:  C.vips_blob_get_type(),
		gValue: &gValue,
	}
	v.free = func(val *GValue) {
		if val.gValue == nil {
			return
		}

		// GNullVipsBlob is used for get result from save_buffer operations,
		// and will be unrefereed with operation
		C.g_value_unset(val.gValue)

		val.gValue = nil
	}
	v.copy = func(val *GValue) (*GValue, error) {
		return nil, ErrCopyForbidden
	}

	C.g_value_init(v.gValue, v.gType)

	return v
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

		newOp := C.vips_cache_operation_build(op.operation)
		if newOp == nil {
			C.g_object_get_property((*C.GObject)(unsafe.Pointer(op.operation)), cOut, &gVal)
			if ptr := C.g_value_peek_pointer(&gVal); ptr != nil {
				C.g_object_unref(ptr)
			}
			C.g_value_unset(&gVal)

			return nil, vipsError()
		}

		C.g_object_unref(C.gpointer(op.operation))
		op.operation = newOp

		C.g_object_get_property((*C.GObject)(unsafe.Pointer(op.operation)), cOut, &gVal)

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
