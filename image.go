package imgvips

/*
#cgo pkg-config: vips
#include "vips/vips.h"
*/
import "C"

// Image wrapper around *C.VipsImage
type Image struct {
	image *C.VipsImage

	val *GValue
}

// Width return image width
// Return 0 if image was freed
func (i *Image) Width() int {
	if i.val.gValue == nil {
		return 0
	}

	return int(C.vips_image_get_width(i.image))
}

// Height return image height
// Return 0 if image was freed
func (i *Image) Height() int {
	if i.val.gValue == nil {
		return 0
	}

	return int(C.vips_image_get_height(i.image))
}
