# IMGVips
[![Build Status](https://github.com/Arimeka/imgvips/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/Arimeka/imgvips/actions/workflows/test.yml)
[![Coverage Status](https://coveralls.io/repos/github/Arimeka/imgvips/badge.svg?branch=master)](https://coveralls.io/github/Arimeka/imgvips?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Arimeka/imgvips)](https://goreportcard.com/report/github.com/Arimeka/imgvips)
[![Go Reference](https://pkg.go.dev/badge/github.com/Arimeka/imgvips.svg)](https://pkg.go.dev/github.com/Arimeka/imgvips)

Low-level bindings for [libvips](https://github.com/libvips/libvips). For better shooting to the leg.

Implements [VipsOperation](https://libvips.github.io/libvips/API/current/VipsOperation.html) builder, which allows you to call most
of the existing image operations in the library.

# Requirements

* [livips](https://github.com/libvips/libvips) 8+ (a higher version is usually better)
* Go 1.11+ (should work on lower versions, but I did not check)

# Memory leak

To reduce memory leak and avoid possible SIGSEGV is recommended to disable libvips cache and vector calculations,
i.e. use:

```
imgvips.VipsCacheSetMaxMem(0)
imgvips.VipsCacheSetMax(0)
imgvips.VipsVectorSetEnables(false)
```

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
