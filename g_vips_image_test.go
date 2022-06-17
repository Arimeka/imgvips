package imgvips_test

import (
	"testing"

	"github.com/Arimeka/imgvips"
)

func TestGVipsImage(t *testing.T) {
	imgvips.VipsCacheSetMaxMem(0)
	imgvips.VipsCacheSetMax(0)

	v, op := generateImage(t)
	defer op.Free()

	_, ok := v.Double()
	if ok {
		t.Fatal("Expected to be not ok")
	}

	result, ok := v.Image()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result == nil {
		t.Fatal("Expected return image, got nil")
	}

	// Check multiply free
	v.Free()
	v.Free()

	result, ok = v.Image()
	if ok {
		t.Fatal("Expected to not be ok")
	}
	if result != nil {
		t.Fatalf("Expected return %v, got %v", nil, result)
	}
}

func TestGNullVipsImage(t *testing.T) {
	imgvips.VipsCacheSetMaxMem(0)
	imgvips.VipsCacheSetMax(0)

	v := imgvips.GNullVipsImage()

	_, ok := v.Int()
	if ok {
		t.Fatal("Expected to be not ok")
	}

	result, ok := v.Image()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result != nil {
		t.Fatalf("Expected return %v, got %v", nil, result)
	}

	// Check multiply free
	v.Free()
	v.Free()

	result, ok = v.Image()
	if ok {
		t.Fatal("Expected to not be ok")
	}
	if result != nil {
		t.Fatalf("Expected return %v, got %v", nil, result)
	}
}

func TestGValue_CopyImage(t *testing.T) {
	imgvips.VipsDetectMemoryLeak(true)
	imgvips.VipsCacheSetMaxMem(0)
	imgvips.VipsCacheSetMax(0)

	val1, op := generateImage(t)
	defer op.Free()

	val2, err := val1.Copy()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	compareImageValsFull(t, val1, val2)

	val1.Free()
	result1, ok := val1.Image()
	if ok {
		t.Fatal("Expected to not be ok")
	}
	if result1 != nil {
		t.Error("Expected val1 to be freed")
	}

	result2, ok := val2.Image()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result2 == nil {
		t.Error("Expected val2 contain image")
	}

	val2.Free()
	result2, ok = val2.Image()
	if ok {
		t.Fatal("Expected to not be ok")
	}
	if result2 != nil {
		t.Error("Expected val2 to be freed")
	}
}

func generateImage(t *testing.T) (*imgvips.GValue, *imgvips.Operation) {
	op, err := imgvips.NewOperation("grey")
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	val1 := imgvips.GNullVipsImage()
	op.AddInput("width", imgvips.GInt(100))
	op.AddInput("height", imgvips.GInt(100))
	op.AddOutput("out", val1)

	if err := op.Exec(); err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	return val1, op
}

func compareImageValsFull(t *testing.T, val1, val2 *imgvips.GValue) {
	result1, ok := val1.Image()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result1 == nil {
		t.Error("Expected val1 contain image")
	}
	result2, ok := val2.Image()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result2 == nil {
		t.Error("Expected val1 contain image")
	}
	if result2.Ptr() == result1.Ptr() {
		t.Errorf("Expected val2 contain %p different from val1 %p", result2.Ptr(), result1.Ptr())
	}
}

func TestGValue_CopyNullImage(t *testing.T) {
	imgvips.VipsDetectMemoryLeak(true)
	imgvips.VipsCacheSetMaxMem(0)
	imgvips.VipsCacheSetMax(0)

	val1 := imgvips.GNullVipsImage()
	val2, err := val1.Copy()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	result1, ok := val1.Image()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result1 != nil {
		t.Error("Expected val1 contain null")
	}
	result2, ok := val2.Image()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result2 != nil {
		t.Error("Expected val1 contain null")
	}
	if result2 != result1 {
		t.Errorf("Expected val2 contain null pointer %p, same in val1 %p", result2, result1)
	}
}

func BenchmarkGNullVipsImage(b *testing.B) {
	imgvips.VipsDetectMemoryLeak(true)
	imgvips.VipsCacheSetMaxMem(0)
	imgvips.VipsCacheSetMax(0)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val := imgvips.GNullVipsImage()
		val.Free()
	}
}

func BenchmarkGValue_CopyVipsImage(b *testing.B) {
	imgvips.VipsDetectMemoryLeak(true)
	imgvips.VipsCacheSetMaxMem(0)
	imgvips.VipsCacheSetMax(0)

	op, err := imgvips.NewOperation("webpload")
	if err != nil {
		b.Fatalf("Unexpected error %v", err)
	}
	defer op.Free()

	out := imgvips.GNullVipsImage()
	op.AddInput("filename", imgvips.GString("./tests/fixtures/img.webp"))
	op.AddInput("scale", imgvips.GDouble(0.1))
	op.AddOutput("out", out)

	if err := op.Exec(); err != nil {
		b.Fatalf("Unexpected error %v", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val, err := out.Copy()
		if err != nil {
			b.Fatalf("Unexpected error %v", err)
		}
		val.Free()
	}
}
