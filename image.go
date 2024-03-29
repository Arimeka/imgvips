package imgvips

/*
#cgo pkg-config: vips
#include "vips/vips.h"
*/
import "C"
import (
	"unsafe"
)

// Image wrapper around *C.VipsImage
type Image struct {
	image *C.VipsImage

	val *GValue
}

// Ptr return unsafe pointer to *C.VipsImage
func (i *Image) Ptr() unsafe.Pointer {
	return unsafe.Pointer(i.image)
}

// Width return image width
// Return 0 if image was freed
func (i *Image) Width() int {
	if i.val.wasFreed() {
		return 0
	}

	return int(C.vips_image_get_width(i.image))
}

// Height return image height
// Return 0 if image was freed
func (i *Image) Height() int {
	if i.val.wasFreed() {
		return 0
	}

	return int(C.vips_image_get_height(i.image))
}
