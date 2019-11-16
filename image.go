package imgvips

/*
#cgo pkg-config: vips
#include "vips/vips.h"
*/
import "C"

// NewVipsImage create empty *C.VipsImage
func NewVipsImage() *C.VipsImage {
	return C.vips_image_new()
}

// ImageWidth get width from *C.VipsImage
func ImageWidth(img *C.VipsImage) int {
	return int(C.vips_image_get_width(img))
}

// ImageHeight get height from *C.VipsImage
func ImageHeight(img *C.VipsImage) int {
	return int(C.vips_image_get_height(img))
}
