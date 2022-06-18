package imgvips_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/Arimeka/imgvips"
)

func TestGDouble(t *testing.T) {
	value := rand.New(rand.NewSource(time.Now().UnixNano())).Float64() // nolint:gosec // For testing

	v := imgvips.GDouble(value)

	_, ok := v.String()
	if ok {
		t.Fatal("Expected to be not ok")
	}

	result, ok := v.Double()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result != value {
		t.Fatalf("Expected return %f, got %f", value, result)
	}

	// Check multiply free
	v.Free()
	v.Free()

	result, ok = v.Double()
	if ok {
		t.Fatal("Expected to not be ok")
	}
	if result != 0 {
		t.Fatalf("Expected return %f, got %f", float64(0), result)
	}
}

func TestGValue_CopyDouble(t *testing.T) {
	v := rand.New(rand.NewSource(time.Now().UnixNano())).Float64() // nolint:gosec // For testing

	val1 := imgvips.GDouble(v)

	val2, err := val1.Copy()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	compareDoubleValsFull(t, v, val1, val2)

	val1.Free()
	result1, ok := val1.Double()
	if ok {
		t.Fatal("Expected to not be ok")
	}
	if result1 != 0 {
		t.Errorf("Expected val1 contain %f gValue, got %f", 0.0, result1)
	}

	result2, ok := val2.Double()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result2 != v {
		t.Errorf("Expected val2 contain %f gValue, got %f", v, result2)
	}

	val2.Free()
	result2, ok = val2.Double()
	if ok {
		t.Fatal("Expected to not be ok")
	}
	if result2 != 0 {
		t.Errorf("Expected val2 contain %f gValue, got %f", 0.0, result2)
	}
}

func compareDoubleValsFull(t *testing.T, v float64, val1, val2 *imgvips.GValue) {
	result1, ok := val1.Double()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result1 != v {
		t.Errorf("Expected val1 contain %f gValue, got %f", v, result1)
	}
	result2, ok := val2.Double()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result2 != v {
		t.Errorf("Expected val2 contain %f gValue, got %f", v, result2)
	}
}

func BenchmarkGDouble(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val := imgvips.GDouble(float64(i))
		result, ok := val.Double()
		if !ok {
			b.Fatal("Expected to be ok")
		}
		if result != float64(i) {
			b.Fatalf("Expected return %f, got %f", float64(i), result)
		}
		val.Free()
	}
}
