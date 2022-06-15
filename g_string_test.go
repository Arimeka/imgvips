package imgvips_test

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/Arimeka/imgvips"
)

func TestGString(t *testing.T) {
	str := "test string"
	v := imgvips.GString(str)

	_, ok := v.Image()
	if ok {
		t.Fatal("Expected to be not ok")
	}

	result, ok := v.String()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result != str {
		t.Fatalf("Expected return %s, got %s", str, result)
	}

	// Check multiply free
	v.Free()
	v.Free()

	result, ok = v.String()
	if ok {
		t.Fatal("Expected to not be ok")
	}
	if result != "" {
		t.Fatalf("Expected return %s, got %s", "", result)
	}
}

func TestGValue_CopyString(t *testing.T) {
	v := strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(20000)) // nolint:gosec // For testing

	val1 := imgvips.GString(v)

	val2, err := val1.Copy()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	compareStringValsFull(t, v, val1, val2)

	val1.Free()
	result1, ok := val1.String()
	if ok {
		t.Fatal("Expected to not be ok")
	}
	if result1 != "" {
		t.Errorf("Expected val1 contain empty string gValue, got %s", result1)
	}

	result2, ok := val2.String()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result2 != v {
		t.Errorf("Expected val2 contain %s gValue, got %s", v, result2)
	}

	val2.Free()
	result2, ok = val2.String()
	if ok {
		t.Fatal("Expected to not be ok")
	}
	if result2 != "" {
		t.Errorf("Expected val2 contain empty string gValue, got %s", result2)
	}
}

func compareStringValsFull(t *testing.T, v string, val1, val2 *imgvips.GValue) {
	result1, ok := val1.String()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result1 != v {
		t.Errorf("Expected val1 contain %s gValue, got %s", v, result1)
	}
	result2, ok := val2.String()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result2 != v {
		t.Errorf("Expected val2 contain %s gValue, got %s", v, result2)
	}
}

func BenchmarkGString(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expected := strconv.Itoa(i)
		val := imgvips.GString(expected)
		result, ok := val.String()
		if !ok {
			b.Fatal("Expected to be ok")
		}
		if result != expected {
			b.Fatalf("Expected return %s, got %s", expected, result)
		}
		val.Free()
	}
}
