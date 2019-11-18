package imgvips_test

import (
	"bytes"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/Arimeka/imgvips"
)

func TestGBoolean(t *testing.T) {
	v := imgvips.GBoolean(true)

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
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result {
		t.Fatal("Expected return false, got true")
	}

	v = imgvips.GBoolean(false)
	result, ok = v.Boolean()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result {
		t.Fatal("Expected return false, got true")
	}
}

func TestGInt(t *testing.T) {
	value := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(10000)

	v := imgvips.GInt(value)

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
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result != 0 {
		t.Fatalf("Expected return %d, got %d", 0, result)
	}
}

func TestGDouble(t *testing.T) {
	value := rand.New(rand.NewSource(time.Now().UnixNano())).Float64()

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
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result != 0 {
		t.Fatalf("Expected return %f, got %f", float64(0), result)
	}
}

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
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result != "" {
		t.Fatalf("Expected return %s, got %s", "", result)
	}
}

func TestGBytes(t *testing.T) {
	imgvips.VipsDetectMemoryLeak(true)

	data := []byte("foobar")
	v := imgvips.GVipsBlob(data)

	result, ok := v.Bytes()
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

func TestGVipsImage(t *testing.T) {
	v := imgvips.GVipsImage()

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
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result != nil {
		t.Fatalf("Expected return %v, got %v", nil, result)
	}
}

func TestGNullVipsImage(t *testing.T) {
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
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result != nil {
		t.Fatalf("Expected return %v, got %v", nil, result)
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
	if !ok {
		t.Fatal("Expected to be ok")
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
	if !ok {
		t.Fatal("Expected to be ok")
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

func TestGValue_CopyInt(t *testing.T) {
	v := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(20000)

	val1 := imgvips.GInt(v)

	val2, err := val1.Copy()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	compareIntValsFull(t, v, val1, val2)

	val1.Free()
	result1, ok := val1.Int()
	if !ok {
		t.Fatal("Expected to be ok")
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
	if !ok {
		t.Fatal("Expected to be ok")
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

func TestGValue_CopyDouble(t *testing.T) {
	v := rand.New(rand.NewSource(time.Now().UnixNano())).Float64()

	val1 := imgvips.GDouble(v)

	val2, err := val1.Copy()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	compareDoubleValsFull(t, v, val1, val2)

	val1.Free()
	result1, ok := val1.Double()
	if !ok {
		t.Fatal("Expected to be ok")
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
	if !ok {
		t.Fatal("Expected to be ok")
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

func TestGValue_CopyString(t *testing.T) {
	v := strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Intn(20000))

	val1 := imgvips.GString(v)

	val2, err := val1.Copy()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	compareStringValsFull(t, v, val1, val2)

	val1.Free()
	result1, ok := val1.String()
	if !ok {
		t.Fatal("Expected to be ok")
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
	if !ok {
		t.Fatal("Expected to be ok")
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

func TestGValue_CopyImage(t *testing.T) {
	imgvips.VipsDetectMemoryLeak(true)

	val1, op := generateImage(t)
	defer op.Free()

	val2, err := val1.Copy()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	compareImageValsFull(t, val1, val2)

	val1.Free()
	result1, ok := val1.Image()
	if !ok {
		t.Fatal("Expected to be ok")
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
	if !ok {
		t.Fatal("Expected to be ok")
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
	if result2 == result1 {
		t.Errorf("Expected val2 contain %p different from val1 %p", result2, result1)
	}
}

func TestGValue_CopyNullImage(t *testing.T) {
	imgvips.VipsDetectMemoryLeak(true)

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

func TestGValue_CopyNewImage(t *testing.T) {
	imgvips.VipsDetectMemoryLeak(true)

	val1 := imgvips.GVipsImage()
	val2, err := val1.Copy()
	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}

	compareNewImageVals(t, val1, val2)

	val1.Free()
	result1, ok := val1.Image()
	if !ok {
		t.Fatal("Expected to be ok")
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
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if result2 != nil {
		t.Error("Expected val2 to be freed")
	}
}

func compareNewImageVals(t *testing.T, val1, val2 *imgvips.GValue) {
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
	if result2 == result1 {
		t.Errorf("Expected val2 contain %p different from val1 %p", result2, result1)
	}
}

func TestGValue_CopyBytes(t *testing.T) {
	imgvips.VipsDetectMemoryLeak(true)

	v := imgvips.GNullVipsBlob()

	_, err := v.Copy()
	if err != imgvips.ErrCopyForbidden {
		t.Fatalf("Expected error %v, got %v", imgvips.ErrCopyForbidden, err)
	}
}
