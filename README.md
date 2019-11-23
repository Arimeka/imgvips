# IMGVips
[![Build Status](https://travis-ci.com/Arimeka/imgvips.svg?branch=master)](https://travis-ci.com/Arimeka/imgvips)
[![Coverage Status](https://coveralls.io/repos/github/Arimeka/imgvips/badge.svg?branch=master)](https://coveralls.io/github/Arimeka/imgvips?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Arimeka/imgvips)](https://goreportcard.com/report/github.com/Arimeka/imgvips)
[![GoDoc](https://godoc.org/github.com/Arimeka/imgvips?status.svg)](https://godoc.org/github.com/Arimeka/imgvips)

Low-level bindings for [libvips](https://github.com/libvips/libvips). For better shooting to the leg.

Implements [VipsOperation](https://libvips.github.io/libvips/API/current/VipsOperation.html) builder, which allows you to call most
of the existing image operations in the library.

# Requirements

* [livips](https://github.com/libvips/libvips) 8+ (a higher version is usually better)
* Go 1.11+ (should work on lower versions, but I did not check)

# Usage

See examples folder.

## Load from filename

```
op, err := imgvips.NewOperation("webpload")
if err != nil {
    panic(err)
}
defer op.Free()

out := imgvips.GNullVipsImage()
op.AddInput("filename", imgvips.GString("path/to/image.webp"))
op.AddOutput("out", out)

if err := op.Exec(); err != nil {
    panic(err)
}
```

## Load from bytes

```
op, err := imgvips.NewOperation("webpload_buffer")
if err != nil {
    panic(err)
}
defer op.Free()

out := imgvips.GNullVipsImage()
op.AddInput("buffer", imgvips.GVipsBlob(data))
op.AddOutput("out", out)

if err := op.Exec(); err != nil {
    panic(err)
}
```

## Save to file

```
op, err := imgvips.NewOperation("jpegsave")
if err != nil {
    panic(err)
}
defer op.Free()

op.AddInput("in", gImage)
op.AddInput("filename", imgvips.GString("image.jpg"))

if err := op.Exec(); err != nil {
    panic(err)
}
```

## Save to bytes

```
op, err := imgvips.NewOperation("jpegsave_buffer")
if err != nil {
    panic(err)
}
defer op.Free()

gData := imgvips.GNullVipsBlob()
op.AddInput("in", gImage)
op.AddOutput("buffer", gData)

if err := op.Exec(); err != nil {
    panic(err)
}
```
