package imgvips_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/Arimeka/imgvips"
)

func TestGBytes(t *testing.T) {
	imgvips.VipsDetectMemoryLeak(true)

	v := imgvips.GVipsBlob(nil)
	result, ok := v.Bytes()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if len(result) != 0 {
		t.Fatalf("Expected return data with %d size, got %d size", 0, len(result))
	}
	v.Free()

	data := []byte("foobar")
	v = imgvips.GVipsBlob(data)

	result, ok = v.Bytes()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if len(result) != len(data) {
		t.Fatalf("Expected return data with %d size, got %d size", len(data), len(result))
	}
	if !bytes.Equal(result, data) {
		t.Fatal("Expected result equal to expected data")
	}

	// Check multiply free
	v.Free()
	v.Free()

	result, ok = v.Bytes()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if len(result) != 0 {
		t.Fatalf("Expected return data with %d size, got %d size", 0, len(result))
	}
}

func TestGNullVipsBlob(t *testing.T) {
	v := imgvips.GNullVipsBlob()
	result, ok := v.Bytes()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if len(result) != 0 {
		t.Fatalf("Expected return data with %d size, got %d size", 0, len(result))
	}

	// Check multiply free
	v.Free()
	v.Free()
}

func TestGValue_CopyBytes(t *testing.T) {
	imgvips.VipsDetectMemoryLeak(true)

	v := imgvips.GNullVipsBlob()

	_, err := v.Copy()
	if err != imgvips.ErrCopyForbidden {
		t.Fatalf("Expected error %v, got %v", imgvips.ErrCopyForbidden, err)
	}
}

func BenchmarkGVipsBlob(b *testing.B) {
	imgvips.VipsDetectMemoryLeak(true)

	data, err := ioutil.ReadFile("./tests/fixtures/img.webp")
	if err != nil {
		b.Fatalf("Unexpected error %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val := imgvips.GVipsBlob(data)
		result, ok := val.Bytes()
		if !ok {
			b.Fatal("Expected to be ok")
		}
		if !bytes.Equal(data, result) {
			b.Fatalf("Expected return bytes with len %d, got with len %d", len(data), len(result))
		}
		val.Free()
	}
}

func BenchmarkGNullVipsBlob(b *testing.B) {
	imgvips.VipsDetectMemoryLeak(true)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val := imgvips.GNullVipsBlob()
		val.Free()
	}
}
