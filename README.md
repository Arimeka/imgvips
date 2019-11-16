# IMGVips
[![Build Status](https://travis-ci.org/Arimeka/imgvips.svg?branch=master)](https://travis-ci.org/Arimeka/imgvips)
[![Coverage Status](https://coveralls.io/repos/github/Arimeka/imgvips/badge.svg?branch=master)](https://coveralls.io/github/Arimeka/imgvips?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Arimeka/imgvips)](https://goreportcard.com/report/github.com/Arimeka/imgvips)
[![GoDoc](https://godoc.org/github.com/Arimeka/imgvips?status.svg)](https://godoc.org/github.com/Arimeka/imgvips)

Low-level bindings for [libvips](https://github.com/libvips/libvips). For better shooting to the leg.

Implements [VipsOperation](https://libvips.github.io/libvips/API/current/VipsOperation.html) builder, which allows you to call most
of the existing image operations in the library.

# Requirements

* [livips](https://github.com/libvips/libvips) 8+ (a higher version is usually better)
* Go 1.13+ (should work on lower versions, but I did not check)

