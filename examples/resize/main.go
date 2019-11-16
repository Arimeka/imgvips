package main

/*
#cgo pkg-config: vips
#include "stdlib.h"
#include "vips/vips.h"
*/
import "C"

import (
	"flag"
	"log"
	"unsafe"

	"github.com/Arimeka/imgvips"
)

const (
	defaultFilename = "./tests/fixtures/img.webp"
	defaultOutput   = "./tests/fixtures/out.png"
	defaultWidth    = 100
)

var (
	inFilename, outFilename string
	width                   int
)

func init() {
	imgvips.VipsDetectMemoryLeak(true)

	flag.StringVar(&inFilename, "inFilename", defaultFilename, "path to input file")
	flag.StringVar(&inFilename, "f", defaultFilename, "path to input file (shorthand)")

	flag.IntVar(&width, "width", defaultWidth, "resize to width")
	flag.IntVar(&width, "w", defaultWidth, "resize to width (shorthand)")

	flag.StringVar(&outFilename, "output", defaultOutput, "path to output file")
	flag.StringVar(&outFilename, "o", defaultOutput, "path to output file (shorthand)")
}

func main() {
	flag.Parse()

	// Load file
	loadedImage := load()
	// Resize
	resizedImage := resize(loadedImage)
	// Save
	save(resizedImage)
}

func load() *imgvips.GValue {
	cFilename := C.CString(inFilename)
	defer C.free(unsafe.Pointer(cFilename))

	// Find image type by inFilename. Package does not implement vips_foreign_find_load, so we call it ourselves.
	cOpName := C.vips_foreign_find_load(cFilename)
	if cOpName == nil {
		log.Fatalf("don't know how to load file %s", inFilename)
	}
	opName := C.GoString(cOpName)

	// It is better to calculate the scaling factor (or shrink) and the type of image before loading the image,
	// so that you can use additional arguments if possible, such as shrink for jpeg and webp (especially for webp).
	op, err := imgvips.NewOperation(opName)
	if err != nil {
		log.Fatalf("operation %s not found: %v", opName, err)
	}
	defer op.Free()

	op.AddInput("filename", imgvips.GString(inFilename))
	out := imgvips.GNullVipsImage()
	op.AddOutput("out", out)

	if err := op.Exec(); err != nil {
		log.Fatalf("load %s return error %v", inFilename, err)
	}

	// op.Free() will destroy out variable, so we make a copy
	result, err := out.Copy()
	if err != nil {
		log.Fatalf("failed copy image %v", err)
	}

	return result
}

func resize(in *imgvips.GValue) *imgvips.GValue {
	image, ok := in.Image()
	if !ok {
		log.Fatal("value is not image")
	}

	op, err := imgvips.NewOperation("resize")
	if err != nil {
		log.Fatalf("operation resize not found: %v", err)
	}
	defer op.Free()

	scale := float64(width) / float64(image.Width())

	op.AddInput("in", in)
	op.AddInput("scale", imgvips.GDouble(scale))
	// Set kernel. Using int kinda working
	// Commented because CI uses vips 8.2.2, which does not have this option
	// op.AddInput("kernel", imgvips.GInt(int(C.VIPS_KERNEL_NEAREST)))
	out := imgvips.GNullVipsImage()
	op.AddOutput("out", out)

	if err := op.Exec(); err != nil {
		log.Fatalf("resize image return error %v", err)
	}

	// op.Free() will destroy out variable, so we make a copy
	result, err := out.Copy()
	if err != nil {
		log.Fatalf("failed copy image %v", err)
	}

	return result
}

func save(in *imgvips.GValue) {
	cFilename := C.CString(outFilename)
	defer C.free(unsafe.Pointer(cFilename))

	// Find image type by inFilename. Package does not implement vips_foreign_find_save, so we call it ourselves.
	cOpName := C.vips_foreign_find_save(cFilename)
	if cOpName == nil {
		log.Fatalf("don't know how to save file %s", outFilename)
	}
	opName := C.GoString(cOpName)

	op, err := imgvips.NewOperation(opName)
	if err != nil {
		log.Fatalf("operation %s not found: %v", opName, err)
	}
	defer op.Free()

	op.AddInput("in", in)
	op.AddInput("filename", imgvips.GString(outFilename))

	if err := op.Exec(); err != nil {
		log.Fatalf("save %s return error %v", outFilename, err)
	}
}
