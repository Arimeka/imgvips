package imgvips

/*
#include "stdlib.h"
*/
import "C"

import (
	"sync"
	"unsafe"
)

// Argument contains key-gValue for set it to *C.VipsOperation
type Argument struct {
	cName  *C.char
	gValue *GValue

	mu sync.RWMutex
}

func (a *Argument) name() *C.char {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.cName
}

func (a *Argument) value() *GValue {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.gValue
}

// Free freed argument cName and gValue
func (a *Argument) Free() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.cName != nil {
		C.free(unsafe.Pointer(a.cName))
		a.cName = nil
	}
	if a.gValue != nil {
		a.gValue.Free()
		a.gValue = nil
	}
}
