package imgvips_test

import (
	"testing"

	"github.com/Arimeka/imgvips"
)

func TestImage_Sizes(t *testing.T) {
	val := imgvips.GVipsImage()

	img, ok := val.Image()
	if !ok {
		t.Fatal("Expected to be ok")
	}
	if img == nil {
		t.Fatal("Expected return image, got nil")
	}

	if img.Width() != 1 {
		t.Errorf("Expected width to by %d, got %d", 1, img.Width())
	}
	if img.Height() != 1 {
		t.Errorf("Expected height to by %d, got %d", 1, img.Height())
	}

	val.Free()

	if img.Width() != 0 {
		t.Errorf("Expected width to by %d, got %d", 0, img.Width())
	}
	if img.Height() != 0 {
		t.Errorf("Expected height to by %d, got %d", 0, img.Height())
	}
}
