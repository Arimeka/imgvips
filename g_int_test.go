package imgvips_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/Arimeka/imgvips"
)

func TestGInt(t *testing.T) {
	value := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000) // nolint:gosec // For testing

	v := imgvips.GInt(value)
	defer v.Free()

	_, ok := v.Boolean()
	if ok {
		t.Fatal("Expected to be not ok")
	}

	result, ok := v.Int()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result != value {
		t.Fatalf("Expected return %d, got %d", value, result)
	}

	// Check multiply free
	v.Free()
	v.Free()

	result, ok = v.Int()
	if ok {
		t.Fatal("Expected to not be ok")
	}
	if result != 0 {
		t.Fatalf("Expected return %d, got %d", 0, result)
	}
}

func TestGValue_CopyInt(t *testing.T) {
	v := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(20000) // nolint:gosec // For testing

	val1 := imgvips.GInt(v)

	val2, err := val1.Copy()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	compareIntValsFull(t, v, val1, val2)

	val1.Free()
	result1, ok := val1.Int()
	if ok {
		t.Fatal("Expected to not be ok")
	}
	if result1 != 0 {
		t.Errorf("Expected val1 contain %d gValue, got %d", 0, result1)
	}

	result2, ok := val2.Int()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result2 != v {
		t.Errorf("Expected val2 contain %d gValue, got %d", v, result2)
	}

	val2.Free()
	result2, ok = val2.Int()
	if ok {
		t.Fatal("Expected to not be ok")
	}
	if result2 != 0 {
		t.Errorf("Expected val2 contain %d gValue, got %d", 0, result2)
	}
}

func compareIntValsFull(t *testing.T, v int, val1, val2 *imgvips.GValue) {
	result1, ok := val1.Int()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result1 != v {
		t.Errorf("Expected val1 contain %d gValue, got %d", v, result1)
	}
	result2, ok := val2.Int()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result2 != v {
		t.Errorf("Expected val2 contain %d gValue, got %d", v, result2)
	}
}

func BenchmarkGInt(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val := imgvips.GInt(i)
		result, ok := val.Int()
		if !ok {
			b.Fatal("Expected to be ok")
		}
		if result != i {
			b.Fatalf("Expected return %d, got %d", i, result)
		}
		val.Free()
	}
}
