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

	freed bool
	free  func(val *GValue)
	copy  func(val *GValue) (*GValue, error)
	mu    sync.RWMutex
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
	v.freed = true
}

func (v *GValue) wasFreed() bool {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return v.freed
}
