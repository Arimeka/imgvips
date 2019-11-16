package imgvips

/*
#include "stdlib.h"
*/
import "C"

import (
	"sync"
	"unsafe"
)

// Argument contains key-value for set it to *C.VipsOperation
type Argument struct {
	name  *C.char
	value *GValue

	mu sync.RWMutex
}

// Name return argument name
func (a *Argument) Name() *C.char {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.name
}

// Value return argument value
func (a *Argument) Value() *GValue {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.value
}

// Free freed argument name and value
func (a *Argument) Free() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.name != nil {
		C.free(unsafe.Pointer(a.name))
		a.name = nil
	}
	if a.value != nil {
		a.value.Free()
		a.value = nil
	}
}
