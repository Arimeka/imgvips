package imgvips_test

import (
	"testing"

	"github.com/Arimeka/imgvips"
)

func TestGBoolean(t *testing.T) {
	v := imgvips.GBoolean(true)
	defer v.Free()

	_, ok := v.Int()
	if ok {
		t.Fatal("Expected to be not ok")
	}
	_, ok = v.Bytes()
	if ok {
		t.Fatal("Expected to be not ok")
	}

	result, ok := v.Boolean()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if !result {
		t.Fatal("Expected return true, got false")
	}

	// Check multiply free
	v.Free()
	v.Free()

	result, ok = v.Boolean()
	if ok {
		t.Fatal("Expected to not be ok")
	}
	if result {
		t.Fatal("Expected return false, got true")
	}

	v = imgvips.GBoolean(false)
	defer v.Free()

	result, ok = v.Boolean()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result {
		t.Fatal("Expected return false, got true")
	}
}

func TestGValue_CopyBoolean(t *testing.T) {
	val1 := imgvips.GBoolean(true)

	val2, err := val1.Copy()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	compareBooleanValsFull(t, val1, val2)

	val1.Free()
	result1, ok := val1.Boolean()
	if ok {
		t.Fatal("Expected to not be ok")
	}
	if result1 {
		t.Error("Expected val1 contain false gValue")
	}

	result2, ok := val2.Boolean()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if !result2 {
		t.Error("Expected val2 contain true gValue")
	}

	val2.Free()
	result2, ok = val2.Boolean()
	if ok {
		t.Fatal("Expected to not be ok")
	}
	if result2 {
		t.Error("Expected val2 contain false gValue")
	}
}

func compareBooleanValsFull(t *testing.T, val1, val2 *imgvips.GValue) {
	result1, ok := val1.Boolean()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if !result1 {
		t.Error("Expected val1 contain true gValue")
	}
	result2, ok := val2.Boolean()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if !result2 {
		t.Error("Expected val2 contain true gValue")
	}
}

func BenchmarkGBoolean(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expected := i%2 == 0
		val := imgvips.GBoolean(expected)
		result, ok := val.Boolean()
		if !ok {
			b.Fatal("Expected to be ok")
		}
		if result != expected {
			b.Fatalf("Expected return %v, got %v", expected, result)
		}
		val.Free()
	}
}
