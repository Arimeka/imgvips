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
	// ErrOperationAlreadyFreed operation already call Free()
	ErrOperationAlreadyFreed = errors.New("operation already freed")
)

// NewOperation initialize new *C.VipsOperation.
// If libvips don't known operation with provided name, function return error.
func NewOperation(name string) (*Operation, error) {
	cStr := C.CString(name)
	defer C.free(unsafe.Pointer(cStr))

	op := C.vips_operation_new(cStr)
	if op == nil {
		return nil, vipsError()
	}

	return &Operation{
		operation: op,
	}, nil
}

// Operation wrapper around *C.VipsOperation
// Contains separates arguments for set to operation and arguments to return from operation.
type Operation struct {
	operation *C.VipsOperation

	inputs  []*Argument
	outputs []*Argument
	mu      sync.Mutex
}

// AddInput adds argument for set to operation.
// After call *Operation.Exec(), all values from input arguments will be freed.
func (op *Operation) AddInput(name string, value *GValue) {
	op.mu.Lock()
	defer op.mu.Unlock()

	op.inputs = append(op.inputs, &Argument{name: C.CString(name), value: value})
}

// AddOutput adds argument for get from operation.
// After call Exec(), all values from output arguments will be updated from operation result.
// This arguments will be freed after call *Operation.Free()
func (op *Operation) AddOutput(name string, value *GValue) {
	op.mu.Lock()
	defer op.mu.Unlock()

	op.outputs = append(op.outputs, &Argument{name: C.CString(name), value: value})
}

// Exec executes operation.
// After execute all input arguments will be freed, all output arguments will be update.
// If operation return error, input arguments will be freed, all output arguments will not be update and not be freed.
func (op *Operation) Exec() error {
	op.mu.Lock()
	defer op.mu.Unlock()

	defer func(args []*Argument) {
		for _, arg := range args {
			arg.Free()
		}
	}(op.inputs)

	if op.operation == nil {
		return ErrOperationAlreadyFreed
	}

	for _, arg := range op.inputs {
		C.g_object_set_property((*C.GObject)(unsafe.Pointer(op.operation)), arg.Name(), arg.Value().GValue())
	}
	for _, arg := range op.outputs {
		C.g_object_set_property((*C.GObject)(unsafe.Pointer(op.operation)), arg.Name(), arg.Value().GValue())
	}

	if success := C.vips_cache_operation_buildp(&op.operation); success != 0 {
		return vipsError()
	}

	for _, arg := range op.outputs {
		C.g_object_get_property((*C.GObject)(unsafe.Pointer(op.operation)), arg.Name(), arg.Value().GValue())
	}

	return nil
}

// Free freed operation outputs, unref operation, and clear vips error
func (op *Operation) Free() {
	op.mu.Lock()
	defer op.mu.Unlock()

	for _, arg := range op.outputs {
		arg.Free()
	}

	if op.operation == nil {
		return
	}

	C.g_object_unref(C.gpointer(op.operation))
	VipsErrorFree()

	op.operation = nil
	op.inputs = nil
	op.outputs = nil
}
